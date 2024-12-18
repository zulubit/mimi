package render

import (
	"bytes"
	"encoding/json"
	"html/template"
	"strings"

	"github.com/zulubit/mimi/pkg/read"
)

func PrepareTemplate(config *read.Config, resource *read.Resource) (string, error) {

	pageBody, err := buildContentString(resource.Content)
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
    <meta name="viewport" content="width=device-width, initial-scale=1.0">` + metaDescription(resource) + metaKeywords(resource) + titleTag(config, resource) + metaGlobalSeo(config) + metaExtra(resource) + `
    <script type="module" src="/static/bundle.min.js"></script>
	<link rel="stylesheet" href="/static/bundle.min.css">` + configHead.String() + `
</head>
	<body class="` + resource.Type + " " + resource.Class + `">
` +
		pageBody.String() +
		`
    ` + configBody.String() + `
</body>
</html>
`, nil
}

func titleTag(config *read.Config, resource *read.Resource) string {
	if resource.SEO.Title == "" {
		return `<title>` + config.Seo.Title + " - " + resource.Name + `</title>`
	}

	return `<title>` + resource.SEO.Title + `</title>`
}

func metaDescription(resource *read.Resource) string {
	if resource.SEO.Description != "" {
		return `<meta name="description" content="` + resource.SEO.Description + `">`
	}
	return ""
}

func metaKeywords(resource *read.Resource) string {
	if len(resource.SEO.Keywords) > 0 {
		return `<meta name="keywords" content="` + strings.Join(resource.SEO.Keywords, ", ") + `">`
	}
	return ""
}

func metaGlobalSeo(config *read.Config) string {
	globalString := ""
	if len(config.Seo.Global) > 0 {
		for _, g := range config.Seo.Global {
			globalString += g + " "
		}
	}
	return globalString
}

func metaExtra(resource *read.Resource) string {
	extras := ""
	if len(resource.SEO.Extra) > 0 {
		for _, e := range resource.SEO.Extra {
			extras += e + " "
		}
	}
	return extras
}

func buildConfigStrings(conf *read.Config) (*strings.Builder, *strings.Builder, error) {
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

func buildContentString(content []read.DataItem) (*strings.Builder, error) {
	var contentBuilder strings.Builder

	for _, c := range content {
		if c.Renderable {
			// Marshal the map to JSON
			rawJSON, err := json.Marshal(c.Data)
			if err != nil {
				return nil, err
			}

			// Minify the JSON
			var minifiedJSON bytes.Buffer
			err = json.Compact(&minifiedJSON, rawJSON)
			if err != nil {
				return nil, err
			}

			// Build the content
			contentBuilder.WriteString("\n<")
			contentBuilder.WriteString(c.Template)
			contentBuilder.WriteString(">")

			// Add JSON as a script tag
			contentBuilder.WriteString(`<script type="application/json">`)
			contentBuilder.WriteString(minifiedJSON.String())
			contentBuilder.WriteString(`</script>`)

			// Handle children here
			if len(c.Children) > 0 {
				bs, err := buildContentString(c.Children)
				if err != nil {
					return nil, err
				}
				contentBuilder.WriteString(bs.String())
			}

			contentBuilder.WriteString(`</`)
			contentBuilder.WriteString(c.Template)
			contentBuilder.WriteString(">")
		}
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

// TODO: we will have 3 data types: resources, fragments and meta. Fragments should be wrapped in a div, span or section.
// fragments are just reuable templates that are meant to be set once and carry over all the pages (central editing of headers and footer)
