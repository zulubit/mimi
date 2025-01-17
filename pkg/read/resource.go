package read

import (
	"fmt"
	"os"
	"path/filepath"
)

func ReadResources(dirPath string) (*[][]byte, error) {
	var rawYaml [][]byte

	// Walk through the directory and its subdirectories
	err := filepath.WalkDir(dirPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing %s: %v", path, err)
		}
		fmt.Println(path)
		// Check if the current file is not a directory and has a .json extension
		if !d.IsDir() && filepath.Ext(d.Name()) == ".yaml" {

			rawMd, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %v", path, err)
			}

			rawYaml = append(rawYaml, rawMd)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &rawYaml, nil
}
