package load

import (
	"html/template"

	"github.com/zulubit/mimi/pkg/read"
	"github.com/zulubit/mimi/pkg/validate"
)

// TODO: cache is now not a thing, we need to recache all the page configs and redo the functions that take it in

var config *read.Config
var pages PageCache
var layoutTemplate *template.Template

type PageStack struct {
	Config   read.Page
	Template []byte
	Markdown []byte
}

type PageCache map[string]PageStack

func BuildConfigCache() error {

	rc, err := read.ReadConfig()
	if err != nil {
		return err
	}

	config = rc

	layout, err := template.ParseFiles("templates/layout.html")
	if err != nil {
		return err
	}
	layoutTemplate = layout

	return nil

	return nil
}

func GetConfig() (*read.Config, error) {

	if config == nil {
		err := BuildConfigCache()
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}

func BuildPageCache() error {

	rc, err := read.ReadResources("./content")
	if err != nil {
		return err
	}

	err = validate.ValidateRoutes(rc)
	if err != nil {
		return err
	}

	c := make(PageCache)

	for _, p := range *rc {
		md, err := read.ReadMarkdown(p.Markdown)
		if err != nil {
			return err
		}

		tp, err := read.ReadTemplate(p.Template)
		if err != nil {
			return err
		}

		currStack := PageStack{
			Config:   p,
			Markdown: md,
			Template: tp,
		}

		c[p.Route] = currStack

	}

	pages = c

	return nil
}

func GetPages() (PageCache, error) {
	if pages == nil {
		err := BuildPageCache()
		if err != nil {
			return nil, err
		}
	}

	return pages, nil
}

func GetLayoutTemplate() (*template.Template, error) {
	if layoutTemplate == nil {
		err := BuildConfigCache() // Ensure layout is cached
		if err != nil {
			return nil, err
		}
	}
	return layoutTemplate, nil
}
