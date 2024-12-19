package render

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/zulubit/mimi/pkg/read"
	"gopkg.in/yaml.v3"
)

// Page represents the parsed content and frontmatter of a Markdown file
type Page struct {
	Content template.HTML          // Rendered Markdown content
	Data    map[string]interface{} // Frontmatter fields as key-value pairs
}

func RenderPage(pageConfigPath string, globalConfig *read.Config) (string, error) {
	// Load and parse the page configuration
	pageConfig, seo, err := loadPageConfig(pageConfigPath)
	if err != nil {
		return "", fmt.Errorf("Error loading page configuration: %w", err)
	}

	// Extract paths from the page configuration
	markdownPath, ok := pageConfig["markdown"].(string)
	if !ok {
		return "", fmt.Errorf("'markdown' field missing or invalid in page config")
	}
	layoutPath, ok := pageConfig["layout"].(string)
	if !ok {
		return "", fmt.Errorf("'layout' field missing or invalid in page config")
	}
	templatePath, ok := pageConfig["template"].(string)
	if !ok {
		return "", fmt.Errorf("'template' field missing or invalid in page config")
	}

	// Parse the Markdown file
	page, err := parseMarkdown(markdownPath)
	if err != nil {
		return "", fmt.Errorf("Error parsing Markdown file: %w", err)
	}

	// Parse the page-specific template
	pageTemplate, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("Error loading page template: %w", err)
	}

	// Render the page content
	var pageContent bytes.Buffer
	err = pageTemplate.Execute(&pageContent, page)
	if err != nil {
		return "", fmt.Errorf("Error rendering page content: %w", err)
	}

	// Embed rendered content into the page data
	page.Data["RenderedContent"] = template.HTML(pageContent.String())
	page.Data["seo"] = seo
	page.Data["globalConfig"] = globalConfig // Add global config to the page data

	// Parse the layout template
	layoutTemplate, err := template.ParseFiles(layoutPath)
	if err != nil {
		return "", fmt.Errorf("Error loading layout template: %w", err)
	}

	// Render the final page with layout
	var renderedPage bytes.Buffer
	err = layoutTemplate.Execute(&renderedPage, page.Data)
	if err != nil {
		return "", fmt.Errorf("Error rendering final page: %w", err)
	}

	return renderedPage.String(), nil
}

// parseMarkdown reads and parses a Markdown file into a Page struct
func parseMarkdown(markdownPath string) (Page, error) {
	var page Page

	// Read the Markdown file
	content, err := os.ReadFile(markdownPath)
	if err != nil {
		return page, fmt.Errorf("failed to read Markdown file: %w", err)
	}

	// Split frontmatter and body
	parts := strings.SplitN(string(content), "---", 3)
	if len(parts) < 3 {
		return page, fmt.Errorf("invalid Markdown format in file: %s", markdownPath)
	}

	// Parse frontmatter into Data
	var frontmatter map[string]interface{}
	err = yaml.Unmarshal([]byte(parts[1]), &frontmatter)
	if err != nil {
		return page, fmt.Errorf("failed to parse frontmatter: %w", err)
	}
	page.Data = frontmatter

	// Convert Markdown body to HTML
	var buf bytes.Buffer
	err = goldmark.New().Convert([]byte(parts[2]), &buf)
	if err != nil {
		return page, fmt.Errorf("failed to render Markdown to HTML: %w", err)
	}
	page.Content = template.HTML(buf.String()) // Mark as safe HTML

	return page, nil
}

// loadPageConfig reads and parses the page's JSON configuration file
func loadPageConfig(configPath string) (map[string]interface{}, map[string]string, error) {
	// Read the JSON file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load page config: %w", err)
	}

	// Parse the JSON data
	var pageConfig map[string]interface{}
	err = json.Unmarshal(data, &pageConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse page config: %w", err)
	}

	// Extract SEO metadata
	seo := make(map[string]string)
	if rawSEO, ok := pageConfig["seo"].(map[string]interface{}); ok {
		for key, value := range rawSEO {
			seo[key] = fmt.Sprintf("%v", value)
		}
	}

	return pageConfig, seo, nil
}
