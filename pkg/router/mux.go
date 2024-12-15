package router

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zulubit/mimi/pkg/handle"
	"github.com/zulubit/mimi/pkg/load"
)

// SetupRouter initializes the mux router and defines the routes
func SetupRouter() *mux.Router {

	// Create a new router
	r := mux.NewRouter()
	r.StrictSlash(true)

	// API v1 routes
	api := r.PathPrefix("/api/v1").Subrouter()

	// Health check route
	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		health := map[string]string{
			"status": "healthy",
			"db":     "connected",
		}
		json.NewEncoder(w).Encode(health)
	}).Methods("GET")

	// Build route to trigger JavaScript bundling
	r.HandleFunc("/build", func(w http.ResponseWriter, r *http.Request) {
		err := load.TriggerBuild("./sitedata/theme/", "./static/")
		if err != nil {
			http.Error(w, "Build failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Build successful."))
	}).Methods("GET")

	// Serve static files
	staticDir := "./static/"
	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))),
	).Methods("GET")

	// Catch-all route for resource handling
	// Define this last to ensure it acts as a fallback for undefined routes
	r.PathPrefix("/").HandlerFunc(handle.GetResource).Methods("GET")

	return r
}
