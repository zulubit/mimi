package load

import (
	"fmt"

	"github.com/zulubit/mimi/pkg/read"
	"github.com/zulubit/mimi/pkg/seo"
	"gopkg.in/yaml.v3"
)

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
	Mimi   Mimi                     `yaml:"mimi"`
	SEO    seo.PageSEO              `yaml:"seo"`
	Meta   map[string]interface{}   `yaml:"meta"` // Catch-all for other fields
	Blocks []map[string]interface{} `yaml:"blocks"`
}

type Route string

type PageCache map[Route]*PageData

var config *read.Config
var pages PageCache

func BuildConfigCache() error {
	rc, err := read.ReadConfig()
	if err != nil {
		return err
	}

	config = rc

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

		// TODO: merge seo

		// Build the PageStack
		currStack := pd

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

func parseYaml(p []byte) (*PageData, error) {
	var data PageData
	err := yaml.Unmarshal(p, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}
	return &data, nil
}
