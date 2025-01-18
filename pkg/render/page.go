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

	tp, err := load.GetTemplates()
	if err != nil {
		return "", false, fmt.Errorf("Error parsing templates: %w", err)
	}

	gc, err := load.GetConfig()
	if err != nil {
		return "", false, fmt.Errorf("Error reading global config: %w", err)
	}

	finalSeo := seo.CombineSeo(gc.GlobalSEO, seo.PageSEO(mp.SEO))

	data := PageData{
		Data:         mp,
		GlobalConfig: *gc,
		SEO:          finalSeo,
	}

	var renderedPage bytes.Buffer
	err = tp.ExecuteTemplate(&renderedPage, mp.Mimi.Template, data)
	if err != nil {
		fmt.Printf("Failed to render template: %v", err)
	}

	return renderedPage.String(), false, nil
}
