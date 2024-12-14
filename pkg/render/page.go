package render

import (
	"bytes"
	"encoding/json"
	"html/template"
	"strings"
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

func PrepareTemplate(seo Seo, content []Content) string {
	var contentBuilder strings.Builder

	for _, c := range content {
		// Escape JSON
		escapedData, err := escapeJSON(c.Data)
		if err != nil {
			escapedData = "{}" // Fallback to empty JSON if escaping fails
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

	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="description" content="{{.Seo.Description}}">
    <meta name="keywords" content="{{.Seo.Keywords}}">
    <title>{{.Seo.Title}}</title>
    <script type="module" src="/static/bundle.min.js"></script>
	<link rel="stylesheet" href="/static/bundle.min.css">
</head>
<body>
` +
		contentBuilder.String() +
		`
    <!-- Optionally include custom elements here -->
</body>
</html>
`
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
func RenderPage(seo Seo, content []Content) (string, error) {
	t := PrepareTemplate(seo, content)
	// Parse the colocated template
	tmpl, err := template.New("page").Parse(t)
	if err != nil {
		return "", err

	}

	data := PageData{
		Seo:     seo,
		Content: []Content{},
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
