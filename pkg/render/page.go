package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"

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

	gc, err := load.GetConfig()
	if err != nil {
		return "", fmt.Errorf("Error reading global config: %w", err)
	}

	// Use precompiled template and parsed markdown/meta
	pageData := PageData{
		Content:      template.HTML(mp.Markdown), // Already parsed Markdown
		Data:         mp.Meta,                    // Metadata
		GlobalConfig: *gc,                        // Global config
	}

	// Render the page using the precompiled page-specific template
	var pageBuffer bytes.Buffer
	err = mp.Parsed.Execute(&pageBuffer, pageData)
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
		Data:         mp.Meta,
		GlobalConfig: *gc,
	}

	var renderedPage bytes.Buffer
	err = layoutTemplate.Execute(&renderedPage, layoutData)
	if err != nil {
		return "", fmt.Errorf("Error rendering final page with layout: %w", err)
	}

	return renderedPage.String(), nil
}
