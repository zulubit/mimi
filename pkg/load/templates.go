package load

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
)

var templateCache *template.Template

type BlockContext struct {
	FullData  interface{}
	BlockData interface{}
}

// Function to build the template cache
func buildTemplateCache() error {
	// Create a new template object
	tmpl := template.New("")

	// TODO: on the topic of blocks, figure out how to define their fields.
	// TODO: move template func collection elsewhere for readability
	// Register the custom function first
	tmpl.Funcs(map[string]interface{}{
		"DynamicBlocks": func(blocks []map[string]interface{}, data interface{}) (template.HTML, error) {
			buf := bytes.NewBuffer([]byte{})
			for _, block := range blocks {
				// Ensure that each block has a name and that the template exists
				name, ok := block["name"].(string)
				if !ok {
					return "", fmt.Errorf("block does not have a valid name")
				}

				// Create the BlockContext struct
				blockContext := BlockContext{
					FullData:  data,
					BlockData: block,
				}

				// Execute each block template with the data
				err := tmpl.ExecuteTemplate(buf, name, blockContext)
				if err != nil {
					return "", err
				}
			}
			return template.HTML(buf.String()), nil
		},
		"Dump": func(data interface{}) (template.HTML, error) {
			// Convert the data to a pretty-printed JSON string
			buf := bytes.NewBuffer([]byte{})
			encoder := json.NewEncoder(buf)
			encoder.SetIndent("", "  ") // Pretty print with indentation
			err := encoder.Encode(data)
			if err != nil {
				return "", err
			}
			// Return the data wrapped in a <pre> tag
			return template.HTML(fmt.Sprintf("<pre>%s</pre>", buf.String())), nil
		},
	})

	// Parse the templates from the directories
	var err error
	tmpl, err = tmpl.ParseGlob("sitedata/theme/*.html")
	if err != nil {
		return err
	}

	tmpl, err = tmpl.ParseGlob("sitedata/theme/templates/*.html")
	if err != nil {
		return err
	}

	tmpl, err = tmpl.ParseGlob("sitedata/theme/blocks/*.html")
	if err != nil {
		return err
	}

	// Log the available templates for debugging
	for _, t := range tmpl.Templates() {
		fmt.Println("Available template:", t.Name())
	}

	// Set the global template cache
	templateCache = tmpl

	return nil
}

func GetTemplates() (*template.Template, error) {
	if templateCache == nil {
		err := buildTemplateCache()
		if err != nil {
			return nil, err
		}
	}
	return templateCache, nil
}
