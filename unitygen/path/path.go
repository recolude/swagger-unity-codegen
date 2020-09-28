package path

import (
	"fmt"
	"strings"
)

// Path represents an HTTP endpoint that our unity client can ping
type Path struct {
	route       string
	httpMethod  string
	tags        []string
	operationID string
	security    []SecurityMethodReference
}

// NewPath creates a new path
func NewPath(route, operationID, method string, tags []string, security []SecurityMethodReference) Path {
	return Path{
		route:       route,
		httpMethod:  method,
		operationID: operationID,
		security:    security,
		tags:        tags,
	}
}

// SupportingClasses will write out different helper classes in C# to assist
// in network requests
func (p Path) SupportingClasses() string {
	builder := strings.Builder{}

	fmt.Fprintf(&builder, "public class %sUnityWebRequest {\n\n", p.operationID)
	fmt.Fprint(&builder, "\tprivate UnityEngine.Networking.UnityWebRequest UnderlyingRequest{ get; };\n\n")
	fmt.Fprint(&builder, "}")

	return builder.String()
}
