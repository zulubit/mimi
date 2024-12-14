package main

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/zulubit/mimi/pkg/router"
)

func main() {
	// Set up the router
	r := router.SetupRouter()

	// Logging middleware
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	// Start the server
	err := http.ListenAndServe(":8080", loggedRouter)
	panic(err)
}
