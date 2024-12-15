package render

import (
	"encoding/json"
	"html/template"
	"strings"

	"github.com/zulubit/mimi/pkg/read"
)

type PageData struct {
	Seo     Seo
	Content []Content
}

type Seo struct {
	Title       string
	Description string
	Keywords    string
}

type Content struct {
	Element string // Represents the main element type (e.g., <div>)
	Data    string // Raw HTML or JSON data
	Index   int
}

func PrepareTemplate(config read.Config, seo Seo, content []Content) (string, error) {

	pageBody, err := buildContetString(content)
	if err != nil {
		return "", err
	}

	configHead, configBody, err := buildConfigStrings(config)
	if err != nil {
		return "", err
	}

	return `<!DOCTYPE html>
<html lang="` + config.Settings.Language + `">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="description" content="` + seo.Description + `">
    <meta name="keywords" content="` + seo.Keywords + `">
    <title>` + seo.Title + `</title>
    <script type="module" src="/static/bundle.min.js"></script>
	<link rel="stylesheet" href="/static/bundle.min.css">` + configHead.String() + `
</head>
<body>
` +
		pageBody.String() +
		`
    <!-- Optionally include custom elements here -->` + configBody.String() + `
</body>
</html>
`, nil
}

func buildConfigStrings(conf read.Config) (*strings.Builder, *strings.Builder, error) {
	var contentBuilder strings.Builder

	for _, c := range conf.Inserts.Head {
		// Build the content
		contentBuilder.WriteString("\n")
		contentBuilder.WriteString(c.Script)
	}

	var contentBuilderBody strings.Builder

	for _, c := range conf.Inserts.EndOfBody {
		// Build the content
		contentBuilderBody.WriteString("\n")
		contentBuilderBody.WriteString(c.Script)
	}

	return &contentBuilder, &contentBuilderBody, nil
}

func buildContetString(content []Content) (*strings.Builder, error) {
	var contentBuilder strings.Builder

	for _, c := range content {
		// Escape JSON
		escapedData, err := escapeJSON(c.Data)
		if err != nil {
			return nil, err
		}

		// Build the content
		contentBuilder.WriteString("\n<")
		contentBuilder.WriteString(c.Element)
		contentBuilder.WriteString(` mimi-data="`)
		contentBuilder.WriteString(escapedData)
		contentBuilder.WriteString(`"></`)
		contentBuilder.WriteString(c.Element)
		contentBuilder.WriteString(">")
	}

	return &contentBuilder, nil
}

// escapeJSON ensures that JSON is safe for embedding in HTML attributes
func escapeJSON(input string) (string, error) {
	// Parse the input as JSON to ensure it's valid
	var parsed interface{}
	err := json.Unmarshal([]byte(input), &parsed)
	if err != nil {
		return "", err
	}

	// Marshal it back to a string, escaping special HTML characters
	escaped, err := json.Marshal(parsed)
	if err != nil {
		return "", err
	}

	// Convert JSON bytes to a string and ensure it's HTML-safe
	return template.HTMLEscapeString(string(escaped)), nil
}

// RenderPage renders an HTML page using the provided SEO and content data
func RenderPage(conf read.Config, seo Seo, content []Content) (string, error) {
	t, err := PrepareTemplate(conf, seo, content)
	if err != nil {
		return "", err
	}

	return t, nil
}
