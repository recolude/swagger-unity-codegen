package unitygen

import (
	"strings"

	"github.com/recolude/swagger-unity-codegen/unitygen/convention"
	"github.com/recolude/swagger-unity-codegen/unitygen/path"
	"github.com/recolude/swagger-unity-codegen/unitygen/security"
)

// Service is a collection of API endpoints (dictated by the tags a route has)
type Service struct {
	name  string
	paths []path.Path
}

// NewService creates a collection of paths
func NewService(name string, paths []path.Path) Service {
	return Service{
		name:  name,
		paths: paths,
	}
}

// Name of the service
func (s Service) Name() string {
	return s.name
}

// Paths are what the service contains
func (s Service) Paths() []path.Path {
	return s.paths
}

// ToCSharp writes out the service as a class with collection of functions that
// correspond to calling different routes
func (s Service) ToCSharp(knownModifiers []security.Auth, serviceConfigName string) string {
	className := convention.TitleCase(s.Name())
	if strings.HasSuffix(className, "Service") == false {
		className += "Service"
	}

	builder := strings.Builder{}
	builder.WriteString("public class ")
	builder.WriteString(className)
	builder.WriteString(" {\n\n\tpublic ")
	builder.WriteString(serviceConfigName)
	builder.WriteString(" Config { get; }\n\n\tpublic ")
	builder.WriteString(className)
	builder.WriteString("(")
	builder.WriteString(serviceConfigName)
	builder.WriteString(" Config) {\n\t\tthis.Config = Config;\n\t}\n\n")

	for _, p := range s.paths {
		builder.WriteString(p.SupportingClasses())
		builder.WriteString("\n")
		builder.WriteString(p.ServiceFunction(knownModifiers))
		builder.WriteString("\n")
	}

	builder.WriteString("}")

	return builder.String()
}
