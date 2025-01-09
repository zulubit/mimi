package load

import (
	"bytes"
	"fmt"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"github.com/zulubit/mimi/pkg/read"
	"github.com/zulubit/mimi/pkg/validate"
	"html/template"
)

// PageStack now holds the raw template, parsed template, markdown, and parsed metadata
type PageStack struct {
	Config   read.Page
	Template []byte             // Raw template
	Parsed   *template.Template // Precompiled template
	Markdown []byte
	Meta     map[string]interface{} // Parsed metadata
}

type Route string

type PageCache map[Route]PageStack

var config *read.Config
var pages PageCache
var layoutTemplate *template.Template

func BuildConfigCache() error {
	rc, err := read.ReadConfig()
	if err != nil {
		return err
	}

	config = rc

	// Load layout template
	layout, err := template.ParseFiles("sitedata/theme/layout.html")
	if err != nil {
		return err
	}
	layoutTemplate = layout

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
		// Read Markdown file
		md, err := read.ReadMarkdown(p.Markdown)
		if err != nil {
			return err
		}

		// Parse the template
		tp, err := read.ReadTemplate(p.Template)
		if err != nil {
			return err
		}

		// Parse Markdown and get metadata
		content, meta, err := parseMarkdown(md)
		if err != nil {
			return err
		}

		// Precompile the template
		parsedTemplate, err := template.New("page-" + p.Route).Parse(string(tp))
		if err != nil {
			return err
		}

		currStack := PageStack{
			Config:   p,
			Template: tp,
			Parsed:   parsedTemplate,
			Markdown: content,
			Meta:     meta,
		}

		c[Route(p.Route)] = currStack
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

// parseMarkdown reads and parses a Markdown file into HTML and extracts metadata
func parseMarkdown(markdown []byte) ([]byte, map[string]interface{}, error) {
	prsr := goldmark.New(goldmark.WithExtensions(meta.Meta))

	// Convert Markdown body to HTML
	var buf bytes.Buffer
	context := parser.NewContext()
	err := prsr.Convert(markdown, &buf, parser.WithContext(context))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to render Markdown to HTML: %w", err)
	}

	meta := meta.Get(context)

	return buf.Bytes(), meta, nil
}
