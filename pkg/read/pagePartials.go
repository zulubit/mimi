package read

import (
	"os"
)

// GetConfig reads and parses the configuration file
func ReadMarkdown(path string) ([]byte, error) {
	// Read the configuration file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetConfig reads and parses the configuration file
func ReadTemplate(path string) ([]byte, error) {
	// Read the configuration file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return data, nil
}
