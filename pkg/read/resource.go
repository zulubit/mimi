package read

import (
	"encoding/json"
	"errors"
	"fmt"
)

// GlobalMeta defines metadata common across the site
type GlobalMeta struct {
	CreatedBy string `json:"createdBy"`
	Type      string `json:"type"`
}

// SEO defines the SEO-related fields
type SEO struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Keywords    []string `json:"keywords"`
	Extra       []string `json:"extra"`
}

// DataItem represents a single item in the `data` array
type DataItem struct {
	Type       string                 `json:"type"`
	Renderable bool                   `json:"renderable"`
	Template   string                 `json:"template,omitempty"`
	Data       map[string]interface{} `json:"data"`
	Internal   map[string]interface{} `json:"internal"`
}

// Page defines the overall structure of a page
type Page struct {
	GlobalMeta GlobalMeta `json:"global_meta"`
	SEO        SEO        `json:"seo"`
	Data       []DataItem `json:"data"`
}

// Parse parses the raw JSON and returns a Page struct or an error
func ParsePage(rawJSON []byte) (*Page, error) {
	var page Page
	err := json.Unmarshal(rawJSON, &page)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	// Validate the page struct
	if err := validatePage(&page); err != nil {
		return nil, fmt.Errorf("validation failed: %v", err)
	}

	return &page, nil
}

// validatePage checks for required fields and other constraints
func validatePage(page *Page) error {
	if page.GlobalMeta.CreatedBy == "" {
		return errors.New("global_meta.createdBy is required")
	}
	if page.GlobalMeta.Type == "" {
		return errors.New("global_meta.type is required")
	}
	if page.SEO.Title == "" {
		return errors.New("seo.title is required")
	}
	if page.SEO.Description == "" {
		return errors.New("seo.description is required")
	}
	for _, item := range page.Data {
		if item.Type == "" {
			return errors.New("data item type is required")
		}
		if item.Renderable && item.Template == "" {
			return errors.New("renderable data item must have a template")
		}
	}
	return nil
}
