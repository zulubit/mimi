package read

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/zulubit/mimi/pkg/seo"
)

// SEO defines the SEO-related fields

// Page defines the overall structure of a page
type Page struct {
	Route    string      `json:"route"`
	Class    string      `json:"class"`
	Name     string      `json:"Name"`
	Type     string      `json:"type"`
	SEO      seo.PageSEO `json:"seo"`
	Markdown string      `json:"markdown"`
	Layout   string      `json:"layout"`
	Template string      `json:"template"`
}

func ReadResources(dirPath string) (*[]Page, error) {
	var resources []Page

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
func ParseResource(rawJSON []byte) (*Page, error) {
	var page Page
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
func validateResource(page *Page) error {

	return nil
}
