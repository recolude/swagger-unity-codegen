package unitygen

import "github.com/recolude/swagger-unity-codegen/unitygen/path"

// Service is a collection of API endpoints (dictated by the tags a route has)
type Service struct {
	name  string
	paths []path.Path
}

// Name of the service
func (s Service) Name() string {
	return s.name
}

// Paths are what the service contains
func (s Service) Paths() []path.Path {
	return s.paths
}
