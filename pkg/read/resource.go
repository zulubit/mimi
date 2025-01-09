package read

import (
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

func ReadResources(dirPath string) (*[][]byte, error) {
	var rawMarkdowns [][]byte

	// Walk through the directory and its subdirectories
	err := filepath.WalkDir(dirPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing %s: %v", path, err)
		}
		fmt.Println(path)
		// Check if the current file is not a directory and has a .json extension
		if !d.IsDir() && filepath.Ext(d.Name()) == ".md" {

			rawMd, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %v", path, err)
			}

			rawMarkdowns = append(rawMarkdowns, rawMd)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &rawMarkdowns, nil
}
