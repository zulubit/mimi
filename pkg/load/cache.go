package load

import (
	"fmt"
	"html/template"

	"github.com/zulubit/mimi/pkg/read"
	"github.com/zulubit/mimi/pkg/seo"
	"gopkg.in/yaml.v3"
)

// PageStack now holds the raw template, parsed template, markdown, and parsed metadata
type PageStack struct {
	PageData PageData
	Parsed   *template.Template // Precompiled template
	Seo      seo.SEO
}

type Mimi struct {
	Route       string `yaml:"route"`
	Type        string `yaml:"type"`
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Layout      string `yaml:"layout"`
	Template    string `yaml:"template"`
}

type SEO struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
}

// Combine known structs and a flexible map for unknown fields
type PageData struct {
	Mimi Mimi                   `yaml:"mimi"`
	SEO  SEO                    `yaml:"seo"`
	Meta map[string]interface{} `yaml:"meta"` // Catch-all for other fields
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
	var pcc []PageData
	seenRoutes := make(map[string]struct{}) // Track seen routes to detect conflicts

	for _, p := range *rc {
		// Parse Markdown and get metadata
		pd, err := parseYaml(p)
		if err != nil {
			return fmt.Errorf("failed to parse Markdown: %w", err)
		}

		// Check for route conflicts
		if _, exists := seenRoutes[pd.Mimi.Route]; exists {
			return fmt.Errorf("route conflict detected: %q is defined multiple times", pd.Mimi.Route)
		}
		seenRoutes[pd.Mimi.Route] = struct{}{}

		// Append the page configuration for later validation
		pcc = append(pcc, *pd)

		// Parse the template
		tp, err := read.ReadTemplate(pd.Mimi.Template)
		if err != nil {
			return fmt.Errorf("failed to read template %q: %w", pd.Mimi.Template, err)
		}

		// Precompile the template
		parsedTemplate, err := template.New("page-" + pd.Mimi.Route).Parse(string(tp))
		if err != nil {
			return fmt.Errorf("failed to parse template for route %q: %w", pd.Mimi.Route, err)
		}

		// TODO: merge seo

		// Build the PageStack
		currStack := PageStack{
			PageData: *pd,
			Parsed:   parsedTemplate,
		}

		c[Route(pd.Mimi.Route)] = currStack
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

func parseYaml(p []byte) (*PageData, error) {
	var data PageData
	err := yaml.Unmarshal(p, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}
	return &data, nil
}
