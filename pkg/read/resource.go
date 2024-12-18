package read

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// GlobalMeta defines metadata common across the site
type GlobalMeta struct {
	CreatedBy string `json:"createdBy"`
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
	Type     string                 `json:"type"`
	Template string                 `json:"template,omitempty"`
	Class    string                 `json:"class"`
	Body     string                 `json:"body"`
	Data     map[string]interface{} `json:"data"`
	Children []DataItem             `json:"children"`
}

// Resource defines the overall structure of a page
type Resource struct {
	Route   string                 `json:"route"`
	Class   string                 `json:"class"`
	Name    string                 `json:"Name"`
	Type    string                 `json:"type"`
	Group   string                 `json:"group"`
	Meta    map[string]interface{} `json:"meta"`
	SEO     SEO                    `json:"seo"`
	Content []DataItem             `json:"content"`
}

func ReadResources(dirPath string) (*[]Resource, error) {
	var resources []Resource

	// Walk through the directory and its subdirectories
	err := filepath.WalkDir(dirPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing %s: %v", path, err)
		}
		fmt.Println(path)
		// Check if the current file is not a directory and has a .json extension
		if !d.IsDir() && filepath.Ext(d.Name()) == ".json" {

			rawJSON, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %v", path, err)
			}

			page, err := ParseResource(rawJSON)
			if err != nil {
				return fmt.Errorf("failed to parse file %s: %v", path, err)
			}

			resources = append(resources, *page)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &resources, nil
}

// Parse parses the raw JSON and returns a Page struct or an error
func ParseResource(rawJSON []byte) (*Resource, error) {
	var page Resource
	err := json.Unmarshal(rawJSON, &page)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	// Validate the page struct
	if err := validateResource(&page); err != nil {
		return nil, fmt.Errorf("validation failed: %v", err)
	}

	return &page, nil
}

// validateResource checks for required fields and other constraints
func validateResource(page *Resource) error {
	if page.SEO.Description == "" {
		return errors.New("seo.description is required")
	}
	for _, item := range page.Content {
		if item.Type == "" {
			return errors.New("data item type is required")
		}
	}
	return nil
}
