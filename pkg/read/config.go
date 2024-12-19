package read

import (
	"encoding/json"
	"html/template"
	"os"
)

// Config represents the application's configuration
type Config struct {
	Settings struct {
		DateFormat      string `json:"dateFormat"`
		PaginationLimit int    `json:"paginationLimit"`
		Language        string `json:"language"`
	} `json:"settings"`
	SEO struct {
		Title  string          `json:"title"`
		Global []template.HTML `json:"global"`
	} `json:"seo"`
	Inserts struct {
		Head      []Insert `json:"head"`
		EndOfBody []Insert `json:"endOfBody"`
	} `json:"inserts"`
}

// Insert represents an HTML script or tag to be inserted into the page
type Insert struct {
	Tag    string        `json:"tag"`
	Script template.HTML `json:"script"`
}

// GetConfig reads and parses the configuration file
func ReadConfig() (*Config, error) {
	var config Config

	// Read the configuration file
	data, err := os.ReadFile("config.json")
	if err != nil {
		return nil, err
	}

	// Parse the JSON data into the Config struct
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
