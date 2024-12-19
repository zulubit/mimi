package admin

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

type PageInfo struct {
	Name  string `json:"name"`
	Route string `json:"route"`
}

func ServeAdminHome(w http.ResponseWriter, r *http.Request) {
	// Define the directory containing the page configurations
	contentDir := "content"

	// Find all JSON files in the content directory
	files, err := filepath.Glob(filepath.Join(contentDir, "*.json"))
	if err != nil {
		http.Error(w, "Failed to load page configurations", http.StatusInternalServerError)
		return
	}

	// Parse each JSON file to extract page information
	var pages []PageInfo
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			http.Error(w, "Failed to read page configuration", http.StatusInternalServerError)
			return
		}

		var pageConfig map[string]interface{}
		if err := json.Unmarshal(data, &pageConfig); err != nil {
			http.Error(w, "Failed to parse page configuration", http.StatusInternalServerError)
			return
		}

		// Extract the name and route fields
		name, _ := pageConfig["name"].(string)
		pages = append(pages, PageInfo{Name: name})
	}

	// Serve the Admin Home Page
	adminTemplatePath := "templates/admin_home.html"
	tmpl, err := template.ParseFiles(adminTemplatePath)
	if err != nil {
		http.Error(w, "Failed to load admin home template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, pages); err != nil {
		http.Error(w, "Failed to render admin home", http.StatusInternalServerError)
	}
}
