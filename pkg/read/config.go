package read

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/zulubit/mimi/pkg/dejson"
)

// Config represents the structure of the configuration file
type Config struct {
	Settings struct {
		Timezone        string `json:"timezone"`
		DateFormat      string `json:"dateFormat"`
		PaginationLimit int    `json:"paginationLimit"`
		DefaultPostType string `json:"defaultPostType"`
		Language        string `json:"language"`
		URL             string `json:"url"`
	} `json:"settings"`
	Seo       dejson.SEO `jsodn:"seo"`
	PostTypes []struct {
		Name      string `json:"name"`
		Directory string `json:"directory"`
	} `json:"postTypes"`
	Inserts struct {
		Head      []Insert `json:"head"`
		EndOfBody []Insert `json:"endOfBody"`
	} `json:"inserts"`
}

// Insert represents a script or tag to be injected
type Insert struct {
	Tag    string `json:"tag"`
	Script string `json:"script"`
}

// ReadConfig reads and parses the configuration file
func ReadConfig(filepath string) (*Config, error) {
	// Read the file contents
	contents, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	// Parse the JSON into the Config struct
	var config Config
	err = json.Unmarshal(contents, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %v", err)
	}

	return &config, nil
}
