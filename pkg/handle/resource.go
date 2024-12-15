package handle

import (
	"fmt"
	"net/http"

	"github.com/zulubit/mimi/pkg/load"
	"github.com/zulubit/mimi/pkg/read"
	"github.com/zulubit/mimi/pkg/render"
)

func GetResource(w http.ResponseWriter, r *http.Request) {
	route := r.URL.Path

	resources, err := load.GetResources()
	if err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		fmt.Printf("Error reading resources: %v\n", err)
		return
	}

	resource, found := findResourceWithRoute(resources, route)

	if !found {
		http.NotFound(w, r)
		fmt.Printf("Page not found: %s\n", route)
		return
	}

	config, err := load.GetConfig()
	if err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		fmt.Printf("Error loading congig: %v\n", err)
		return
	}

	// Render the page using RenderPage
	renderedPage, err := render.PrepareTemplate(config, resource)
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

func findResourceWithRoute(resources load.ResoruceMap, route string) (*read.Resource, bool) {
	cleanRoute := route

	if len(route) > 0 && route[len(route)-1] != '/' {
		cleanRoute = route + "/" // Strip trailing slash
	}

	for _, r := range *resources {
		fmt.Println(r.Route)
		nr := r.Route
		if len(r.Route) > 0 && r.Route[len(r.Route)-1] != '/' {
			nr = r.Route + "/" // Strip trailing slash
		}
		if nr == cleanRoute {
			return &r, true
		}
	}

	return nil, false
}
