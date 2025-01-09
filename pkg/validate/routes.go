package validate

import (
	"errors"
	"strings"

	"github.com/zulubit/mimi/pkg/read"
)

func ValidateRoutes(resources *[]read.Page) error {
	routeMap := make(map[string]bool)
	for _, r := range *resources {
		// Check if the route starts with "/"
		if !strings.HasPrefix(r.Route, "/") {
			return errors.New("route must start with '/'")
		}

		// Check for duplicate routes
		if routeMap[r.Route] {
			return errors.New("duplicate route")
		}

		if strings.HasPrefix(r.Route, "/mimi-api/v") || strings.HasPrefix(r.Route, "/mimi-admin") || strings.HasPrefix(r.Route, "/mimi-services") {
			return errors.New("route collides with api/v1(2,3) or /mimi-admin/ or /mimi-services routes")
		}

		routeMap[r.Route] = true
	}
	return nil
}
