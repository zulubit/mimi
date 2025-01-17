package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"

	"github.com/zulubit/mimi/pkg/load"
	"github.com/zulubit/mimi/pkg/read"
	"github.com/zulubit/mimi/pkg/seo"
)

type PageNotFound bool

type PageData struct {
	Content      template.HTML
	Data         interface{}
	GlobalConfig read.Config
	SEO          seo.SEO
}

func RenderPage(route string) (string, PageNotFound, error) {
	pages, err := load.GetPages()
	if err != nil {
		return "", false, err
	}

	mp, ok := pages[load.Route(route)]
	if !ok {
		return "", true, errors.New("page not found in cache")
	}

	gc, err := load.GetConfig()
	if err != nil {
		return "", false, fmt.Errorf("Error reading global config: %w", err)
	}

	// Render the page using the precompiled page-specific template
	var pageBuffer bytes.Buffer
	err = mp.Parsed.Execute(&pageBuffer, mp.PageData)
	if err != nil {
		return "", false, fmt.Errorf("Error rendering page-specific template: %w", err)
	}

	// Retrieve the cached layout template
	layoutTemplate, err := load.GetLayoutTemplate()
	if err != nil {
		return "", false, fmt.Errorf("Error retrieving layout template: %w", err)
	}

	finalSeo := seo.CombineSeo(gc.GlobalSEO, seo.PageSEO(mp.Seo))

	// Render the final page using the layout template
	layoutData := PageData{
		Content:      template.HTML(pageBuffer.String()),
		Data:         mp.PageData,
		GlobalConfig: *gc,
		SEO:          finalSeo,
	}

	var renderedPage bytes.Buffer
	err = layoutTemplate.Execute(&renderedPage, layoutData)
	if err != nil {
		return "", false, fmt.Errorf("Error rendering final page with layout: %w", err)
	}

	return renderedPage.String(), false, nil
}
