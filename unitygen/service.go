package unitygen

import (
	"fmt"
	"strings"

	"github.com/recolude/swagger-unity-codegen/unitygen/convention"
	"github.com/recolude/swagger-unity-codegen/unitygen/path"
	"github.com/recolude/swagger-unity-codegen/unitygen/security"
)

// Service is a collection of API endpoints (dictated by the tags a route has)
type Service struct {
	name        string
	description string
	paths       []path.Path
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

func (s Service) Interface() string {
	builder := strings.Builder{}

	builder.WriteString("public interface ")
	builder.WriteString(s.InterfaceName())

	builder.WriteString(` {
`)

	for _, p := range s.paths {
		builder.WriteString("\t")
		builder.WriteString(p.ServiceInterfaceFunction(""))
		builder.WriteString("\n")
	}

	builder.WriteString("}\n")

	return builder.String()
}

func (s Service) ServiceProvider(menuName string) string {
	className := s.ClassName()

	builder := &strings.Builder{}
	nameWithSpaces := convention.AddSpace(className)
	fmt.Fprintf(builder, "[CreateAssetMenu(fileName = \"%s Provider\", menuName = \"%s/%s Provider\")]\n", nameWithSpaces, menuName, nameWithSpaces)
	fmt.Fprintf(builder, "public class %sProvider: ServiceProvider<%sObject> { }\n\n", className, className)

	return builder.String()
}

func (s Service) ScriptableObject() string {
	builder := strings.Builder{}

	fmt.Fprintf(&builder, "public abstract class %sObject : ScriptableObject, %s {\n", s.ClassName(), s.InterfaceName())

	for _, p := range s.paths {
		builder.WriteString("\t")
		builder.WriteString(p.ServiceInterfaceFunction("abstract"))
		builder.WriteString("\n")
	}

	builder.WriteString("}\n")

	return builder.String()
}

func (s Service) InterfaceName() string {
	return "I" + s.ClassName()
}

func (s Service) ClassName() string {
	className := convention.TitleCase(s.Name())
	if !strings.HasSuffix(className, "Service") {
		className += "Service"
	}
	return className
}

func (s Service) Comments() string {
	if s.description == "" {
		return ""
	}

	builder := strings.Builder{}
	builder.WriteString("/// <summary>\n/// ")
	builder.WriteString(s.description)
	builder.WriteString("\n/// </summary>\n")
	return builder.String()
}

// ToCSharp writes out the service as a class with collection of functions that
// correspond to calling different routes
func (s Service) ToCSharp(knownModifiers []security.Auth, serviceConfigName string, menuName string) string {
	className := s.ClassName()

	builder := &strings.Builder{}

	for _, p := range s.paths {
		builder.WriteString(p.SupportingClasses())
		builder.WriteString("\n")
	}

	builder.WriteString(s.Interface())
	builder.WriteString(s.ScriptableObject())

	builder.WriteString(s.Comments()) // [CreateAssetMenu(fileName = "SemanticLabelService", menuName = "DAT/API/Semantic Label Service")]
	nameWithSpaces := convention.AddSpace(className)
	fmt.Fprintf(builder, "[CreateAssetMenu(fileName = \"%s\", menuName = \"%s/%s\")]\n", nameWithSpaces, menuName, nameWithSpaces)
	fmt.Fprintf(builder, "public class %s: %sObject {\n\n", className, className)
	fmt.Fprintf(builder, "\t[SerializeField]\n\tprivate %s Config;\n", serviceConfigName)

	for _, p := range s.paths {
		builder.WriteString(p.ServiceFunction(knownModifiers))
		builder.WriteString("\n")
	}

	builder.WriteString("}")

	return builder.String()
}
