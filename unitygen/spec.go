package unitygen

import (
	"fmt"
	"sort"
	"strings"

	"github.com/recolude/swagger-unity-codegen/unitygen/convention"
	"github.com/recolude/swagger-unity-codegen/unitygen/definition"
	"github.com/recolude/swagger-unity-codegen/unitygen/security"
)

// Spec is the overall interpretted swagger file
type Spec struct {
	Info            SpecInfo
	Definitions     []definition.Definition
	AuthDefinitions []security.Auth
	Services        []Service
}

func NewSpec(info SpecInfo, definitions []definition.Definition, authDefinitions []security.Auth, services []Service) Spec {
	sort.Sort(sortByDefinitionName(definitions))
	sort.Sort(sortBySecurityIdentifier(authDefinitions))
	return Spec{
		Info:            info,
		Definitions:     definitions,
		AuthDefinitions: authDefinitions,
		Services:        services,
	}
}

type sortByDefinitionName []definition.Definition

func (a sortByDefinitionName) Len() int           { return len(a) }
func (a sortByDefinitionName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortByDefinitionName) Less(i, j int) bool { return a[i].Name() < a[j].Name() }

type sortBySecurityIdentifier []security.Auth

func (a sortBySecurityIdentifier) Len() int           { return len(a) }
func (a sortBySecurityIdentifier) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortBySecurityIdentifier) Less(i, j int) bool { return a[i].Identifier() < a[j].Identifier() }

// SpecInfo is general info about the spec itself
type SpecInfo struct {
	Title   string
	Version string
}

func (s Spec) renderConfigParams(serialize bool) string {
	builder := strings.Builder{}

	if serialize {
		builder.WriteString("\t[SerializeField]\n")
	}
	builder.WriteString("\tpublic string BasePath { get; set; }\n\n")
	for _, authGuard := range s.AuthDefinitions {
		fmt.Fprintf(&builder, "\t// %s\n", authGuard.String())
		if serialize {
			builder.WriteString("\t[SerializeField]\n")
		}
		fmt.Fprintf(&builder, "\tpublic string %s { get; set; }\n\n", convention.TitleCase(authGuard.Identifier()))
	}
	return builder.String()
}

// ServiceConfig prints out a c# class with variables to be used for all requests
func (s Spec) ServiceConfig(configName, menuName string) string {
	properClassName := convention.TitleCase(configName)
	builder := strings.Builder{}

	// Interface
	builder.WriteString("public interface Config {\n\n")
	builder.WriteString(s.renderConfigParams(false))
	builder.WriteString("}\n\n")

	// Scriptable Object
	builder.WriteString("[System.Serializable]\n")
	fmt.Fprintf(&builder, "[CreateAssetMenu(menuName = \"%s\", fileName = \"%s\")]\n", menuName, properClassName)
	fmt.Fprintf(&builder, "public class %s: ScriptableObject, Config {\n\n", properClassName)
	builder.WriteString(s.renderConfigParams(true))
	fmt.Fprintf(&builder, "\tpublic %s(string basePath) {\n", properClassName)
	builder.WriteString("\t\tthis.BasePath = basePath;\n")
	builder.WriteString("\t}\n")
	builder.WriteString("\n}")
	return builder.String()
}
