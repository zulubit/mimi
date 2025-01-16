package load

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"github.com/zulubit/mimi/pkg/read"
	"github.com/zulubit/mimi/pkg/seo"
	"github.com/zulubit/mimi/pkg/validate"
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

	c := make(PageCache)
	var pcc []read.Page
	seenRoutes := make(map[string]struct{}) // Track seen routes to detect conflicts

	for _, p := range *rc {
		// Parse Markdown and get metadata
		content, meta, pageConfig, err := parseMarkdown(p)
		if err != nil {
			return fmt.Errorf("failed to parse Markdown: %w", err)
		}

		// Check for route conflicts
		if _, exists := seenRoutes[pageConfig.Route]; exists {
			return fmt.Errorf("route conflict detected: %q is defined multiple times", pageConfig.Route)
		}
		seenRoutes[pageConfig.Route] = struct{}{}

		// Append the page configuration for later validation
		pcc = append(pcc, *pageConfig)

		// Parse the template
		tp, err := read.ReadTemplate(pageConfig.Template)
		if err != nil {
			return fmt.Errorf("failed to read template %q: %w", pageConfig.Template, err)
		}

		// Precompile the template
		parsedTemplate, err := template.New("page-" + pageConfig.Route).Parse(string(tp))
		if err != nil {
			return fmt.Errorf("failed to parse template for route %q: %w", pageConfig.Route, err)
		}

		// Build the PageStack
		currStack := PageStack{
			Config:   *pageConfig,
			Template: tp,
			Parsed:   parsedTemplate,
			Markdown: content,
			Meta:     meta,
		}

		c[Route(pageConfig.Route)] = currStack
	}

	// Validate routes post-build for possible conflicts
	err = validate.ValidateRoutes(&pcc)
	if err != nil {
		return fmt.Errorf("route validation failed: %w", err)
	}

	// Set the global page cache
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
func parseMarkdown(markdown []byte) ([]byte, map[string]interface{}, *read.Page, error) {
	prsr := goldmark.New(goldmark.WithExtensions(meta.Meta))

	// Convert Markdown body to HTML
	var buf bytes.Buffer
	context := parser.NewContext()
	err := prsr.Convert(markdown, &buf, parser.WithContext(context))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to render Markdown to HTML: %w", err)
	}

	metaData := meta.Get(context)

	// Define required keys for Page
	requiredKeys := map[string]*string{
		"mimi-route":    new(string),
		"mimi-classes":  new(string),
		"mimi-type":     new(string),
		"mimi-layout":   new(string),
		"mimi-template": new(string),
	}

	// Define SEO keys for PageSEO
	seoKeys := map[string]*string{
		"mimi-title":       new(string),
		"mimi-description": new(string),
	}

	// Extract metadata values for Page
	for key, ptr := range requiredKeys {
		if val, ok := metaData[key].(string); ok {
			*ptr = val
		} else {
			*ptr = "" // Default to empty string if missing
		}
	}

	// Extract metadata values for SEO
	for key, ptr := range seoKeys {
		if val, ok := metaData[key].(string); ok {
			*ptr = val
		} else {
			*ptr = "" // Default to empty string if missing
		}
	}

	// Extract Keywords for SEO
	var keywords []string
	if kw, ok := metaData["mimi-keywords"].([]interface{}); ok {
		for _, v := range kw {
			if str, ok := v.(string); ok {
				keywords = append(keywords, str)
			}
		}
	}

	// Extract ExtraSEO as []template.HTML
	var extraSEO []template.HTML
	if ex, ok := metaData["mimi-extraseo"].([]interface{}); ok {
		for _, v := range ex {
			if str, ok := v.(string); ok {
				extraSEO = append(extraSEO, template.HTML(str))
			}
		}
	}

	// Populate the Page struct
	page := &read.Page{
		Route:    *requiredKeys["mimi-route"],
		Class:    *requiredKeys["mimi-classes"],
		Type:     *requiredKeys["mimi-type"],
		Layout:   *requiredKeys["mimi-layout"],
		Template: *requiredKeys["mimi-template"],
		Markdown: string(markdown),
		SEO: seo.PageSEO{
			Title:       *seoKeys["mimi-title"],
			Description: *seoKeys["mimi-description"],
			Keywords:    keywords,
			Extra:       extraSEO,
		},
	}

	return buf.Bytes(), metaData, page, nil
}
