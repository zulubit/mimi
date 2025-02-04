package handle

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/zulubit/mimi/pkg/render"
)

func GetResource(w http.ResponseWriter, r *http.Request) {

	route := r.URL.Path

	renderedPage, notFound, err := render.RenderPage(route)

	// TODO: Implement a default 404 html page and take 404.html from theme if available.
	if notFound {
		health := map[string]string{
			"error": "not found",
		}

		w.WriteHeader(http.StatusNotFound)

		json.NewEncoder(w).Encode(health)
		return
	}

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
