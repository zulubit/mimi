package handle

import (
	"fmt"
	"net/http"

	"github.com/zulubit/mimi/pkg/render"
)

func GetResource(w http.ResponseWriter, r *http.Request) {

	// Render the page
	renderedPage, err := render.RenderPage("/")
	if err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		fmt.Printf("Error rendering page: %v\n", err)
		return
	}

	//TODO: find the right page or 404

	// Write the rendered HTML
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(renderedPage)); err != nil {
		fmt.Printf("Error writing response: %v\n", err)
	}
}
