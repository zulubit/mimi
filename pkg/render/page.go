package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"github.com/zulubit/mimi/pkg/load"
	"github.com/zulubit/mimi/pkg/read"
)

type PageData struct {
	Content      template.HTML
	Data         map[string]interface{}
	GlobalConfig read.Config
}

func RenderPage(route string) (string, error) {
	pages, err := load.GetPages()
	if err != nil {
		return "", err
	}

	mp, ok := pages[route]
	if !ok {
		return "", errors.New("page not found in cache")
	}

	// Parse the Markdown file
	pageContent, meta, err := parseMarkdown(mp.Markdown)
	if err != nil {
		return "", fmt.Errorf("Error parsing Markdown file: %w", err)
	}

	gc, err := load.GetConfig()
	if err != nil {
		return "", fmt.Errorf("Error reading global config: %w", err)
	}

	// Prepare data for the page-specific template
	pageData := PageData{
		Content:      template.HTML(pageContent),
		Data:         meta,
		GlobalConfig: *gc,
	}

	// Render the page-specific template
	pageTemplate, err := template.New("page").Parse(string(mp.Template))
	if err != nil {
		return "", fmt.Errorf("Error loading page template: %w", err)
	}

	var pageBuffer bytes.Buffer
	err = pageTemplate.Execute(&pageBuffer, pageData)
	if err != nil {
		return "", fmt.Errorf("Error rendering page-specific template: %w", err)
	}

	// Retrieve the cached layout template
	layoutTemplate, err := load.GetLayoutTemplate()
	if err != nil {
		return "", fmt.Errorf("Error retrieving layout template: %w", err)
	}

	// Render the final page using the layout template
	layoutData := PageData{
		Content:      template.HTML(pageBuffer.String()),
		Data:         meta,
		GlobalConfig: *gc,
	}

	var renderedPage bytes.Buffer
	err = layoutTemplate.Execute(&renderedPage, layoutData)
	if err != nil {
		return "", fmt.Errorf("Error rendering final page with layout: %w", err)
	}

	return renderedPage.String(), nil
}

// parseMarkdown reads and parses a Markdown file into a Page struct
func parseMarkdown(markdown []byte) ([]byte, map[string]interface{}, error) {
	prsr := goldmark.New(goldmark.WithExtensions(meta.Meta))

	// Convert Markdown body to HTML
	var buf bytes.Buffer
	context := parser.NewContext()
	err := prsr.Convert([]byte(markdown), &buf, parser.WithContext(context))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to render Markdown to HTML: %w", err)
	}

	meta := meta.Get(context)

	return buf.Bytes(), meta, nil
}
