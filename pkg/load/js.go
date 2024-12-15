package load

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/evanw/esbuild/pkg/api"
)

func TriggerBuild(inputDir, outputDir string) error {
	// Ensure the output directory exists
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Define the entry point (main.js) and output file
	entryPoint := filepath.Join(inputDir, "main.js")
	outFile := filepath.Join(outputDir, "bundle.min.js")

	// Check if main.js exists
	if _, err := os.Stat(entryPoint); os.IsNotExist(err) {
		return fmt.Errorf("entry file %s does not exist", entryPoint)
	}

	// Bundle the JavaScript files using esbuild
	result := api.Build(api.BuildOptions{
		EntryPoints:       []string{entryPoint},
		Outfile:           outFile,
		Bundle:            true,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		Loader: map[string]api.Loader{
			".scss": api.LoaderCSS, // Use CSS loader for SCSS files
		},
		MinifySyntax: true,
		Sourcemap:    api.SourceMapLinked,
		Write:        true,
		LogLevel:     api.LogLevelInfo,
	})

	// Check for errors during the build process
	if len(result.Errors) > 0 {
		return fmt.Errorf("build failed: %v", result.Errors)
	}

	return nil
}
