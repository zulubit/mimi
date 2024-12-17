package validate

import (
	"errors"
	"strings"

	"github.com/zulubit/mimi/pkg/read"
)

func ValidateRoutes(resources *[]read.Resource) error {
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
		routeMap[r.Route] = true
	}
	return nil
}
