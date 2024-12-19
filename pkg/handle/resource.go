package handle

import (
	"fmt"
	"net/http"

	"github.com/zulubit/mimi/pkg/read"
	"github.com/zulubit/mimi/pkg/render"
)

func GetResource(w http.ResponseWriter, r *http.Request) {
	// Load global configuration
	globalConfig, err := read.ReadConfig()
	if err != nil {
		http.Error(w, "Error loading configuration", http.StatusInternalServerError)
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	// Render the page
	renderedPage, err := render.RenderPage("content/page.json", globalConfig)
	if err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		fmt.Printf("Error rendering page: %v\n", err)
		return
	}

	// Write the rendered HTML
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(renderedPage)); err != nil {
		fmt.Printf("Error writing response: %v\n", err)
	}
}
