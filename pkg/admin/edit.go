package admin

// TODO: these methods need to be changed to actually load resources lol

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

type PageConfig struct {
	Route    string            `json:"route"`
	Type     string            `json:"type"`
	Name     string            `json:"name"`
	Markdown string            `json:"markdown"`
	Layout   string            `json:"layout"`
	Template string            `json:"template"`
	SEO      map[string]string `json:"seo"`
}

type AdminData struct {
	Config   string // JSON configuration as a string
	Markdown string // Markdown content
	Template string // HTML template content
}

func ServeAdminDashboard(w http.ResponseWriter, r *http.Request) {
	// Extract the 'name' query parameter from the URL
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Missing 'name' parameter", http.StatusBadRequest)
		return
	}

	// Find the page configuration based on the name
	configPath, err := findPageConfig(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Load the page configuration
	configData, err := os.ReadFile(configPath)
	if err != nil {
		http.Error(w, "Failed to load page configuration", http.StatusInternalServerError)
		return
	}

	// Parse the JSON configuration
	var pageConfig PageConfig
	if err := json.Unmarshal(configData, &pageConfig); err != nil {
		http.Error(w, "Failed to parse page configuration", http.StatusInternalServerError)
		return
	}

	// Load the Markdown content
	markdownContent, err := os.ReadFile(pageConfig.Markdown)
	if err != nil {
		http.Error(w, "Failed to load Markdown file", http.StatusInternalServerError)
		return
	}

	// Load the template content
	templateContent, err := os.ReadFile(pageConfig.Template)
	if err != nil {
		http.Error(w, "Failed to load template file", http.StatusInternalServerError)
		return
	}

	// Prepare admin data for embedding
	adminData := AdminData{
		Config:   string(configData),
		Markdown: string(markdownContent),
		Template: string(templateContent),
	}

	// Parse and execute the admin dashboard template
	adminTemplatePath := "templates/admin_dashboard.html"
	tmpl, err := template.ParseFiles(adminTemplatePath)
	if err != nil {
		http.Error(w, "Failed to load admin dashboard template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, adminData); err != nil {
		http.Error(w, "Failed to render admin dashboard", http.StatusInternalServerError)
	}
}

func findPageConfig(name string) (string, error) {
	contentDir := "content"
	files, err := filepath.Glob(filepath.Join(contentDir, "*.json"))
	if err != nil {
		return "", err
	}

	// Debug: Log the list of files found
	fmt.Printf("Searching for page with name '%s' in files: %v\n", name, files)

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			// Debug: Log file read error
			fmt.Printf("Failed to read file '%s': %v\n", file, err)
			continue
		}

		var pageConfig PageConfig
		if err := json.Unmarshal(data, &pageConfig); err != nil {
			// Debug: Log JSON parsing error
			fmt.Printf("Failed to parse JSON file '%s': %v\n", file, err)
			continue
		}

		// Debug: Log the parsed Name field
		fmt.Printf("Checking file '%s': Name='%s'\n", file, pageConfig.Name)

		if pageConfig.Name == name {
			// Debug: Log the match
			fmt.Printf("Match found: '%s'\n", file)
			return file, nil
		}
	}

	return "", fmt.Errorf("Page with name '%s' not found", name)
}
