package path

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/recolude/swagger-unity-codegen/unitygen/security"
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

func (p Path) reqPathName() string {
	return fmt.Sprintf("%sUnityWebRequest", p.operationID)
}

func (p Path) toUnityHTTPVerb() string {
	switch p.httpMethod {
	case http.MethodGet:
		return "UnityWebRequest.kHttpVerbGET"
	}
	panic(fmt.Sprintf("unknown verb \"%s\"", p.httpMethod))
}

// SupportingClasses will write out different helper classes in C# to assist
// in network requests
func (p Path) SupportingClasses() string {
	builder := strings.Builder{}

	fmt.Fprintf(&builder, "public class %s {\n\n", p.reqPathName())

	// underlying network request
	fmt.Fprint(&builder, "\tpublic UnityWebRequest UnderlyingRequest{ get; };\n\n")

	// constructor
	fmt.Fprintf(&builder, "\tpublic %s(UnityWebRequest req) {\n\t\tthis.UnderlyingRequest = req;\n\t}\n\n", p.reqPathName())
	fmt.Fprint(&builder, "}")

	return builder.String()
}

func (p Path) guard(reference SecurityMethodReference, knownModifiers []security.Auth) security.Auth {
	for _, g := range knownModifiers {
		if g.Identifier() == reference.Identifier {
			return g
		}
	}
	panic("no known modifier matches reference " + reference.Identifier)
}

// ServiceFunction generates C# code that is used to make network requests
func (p Path) ServiceFunction(knownModifiers []security.Auth) string {
	builder := strings.Builder{}

	fmt.Fprintf(&builder, "public %s %s()\n{\n", p.reqPathName(), p.operationID)
	fmt.Fprintf(&builder, "\tvar unityNetworkReq = new UnityWebRequest(this.config.BasePath + \"%s\", %s);\n", p.route, p.toUnityHTTPVerb())
	if len(p.security) == 1 {
		fmt.Fprintf(&builder, "\t%s\n", p.guard(p.security[0], knownModifiers).ModifyNetworkRequest())
	}
	fmt.Fprintf(&builder, "\treturn new %s(unityNetworkReq);\n}", p.reqPathName())

	return builder.String()
}
