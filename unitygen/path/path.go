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

	// mapping of HTTP Status Codes to responses
	responses map[string]Response

	parameters []Parameter
}

// NewPath creates a new path
func NewPath(route, operationID, method string, tags []string, security []SecurityMethodReference, responses map[string]Response, parameters []Parameter) Path {
	return Path{
		route:       route,
		httpMethod:  method,
		operationID: operationID,
		security:    security,
		tags:        tags,
		responses:   responses,
		parameters:  parameters,
	}
}

func (p Path) reqPathName() string {
	return fmt.Sprintf("%sUnityWebRequest", p.operationID)
}

func (p Path) toUnityHTTPVerb() string {
	switch p.httpMethod {
	case http.MethodGet:
		return "UnityWebRequest.kHttpVerbGET"
	case http.MethodPut:
		return "UnityWebRequest.kHttpVerbPUT"
	case http.MethodPost:
		return "UnityWebRequest.kHttpVerbPOST"
	case http.MethodDelete:
		return "UnityWebRequest.kHttpVerbDELETE"
	case http.MethodHead:
		return "UnityWebRequest.kHttpVerbHEAD"
	}
	panic(fmt.Sprintf("unknown verb \"%s\"", p.httpMethod))
}

func (p Path) respVariableName(k string) string {
	if k == "200" {
		return "success"
	}

	if k == "default" {
		return "default"
	}

	panic("unkown response key: " + k)
}

// SupportingClasses will write out different helper classes in C# to assist
// in network requests
func (p Path) SupportingClasses() string {
	builder := strings.Builder{}

	fmt.Fprintf(&builder, "public class %s {\n\n", p.reqPathName())

	// Outline all portential responses
	for respKey, resp := range p.responses {
		fmt.Fprintf(&builder, "\tpublic %s %s;\n\n", resp.schema.ToVariableType(), p.respVariableName(respKey))
	}

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

func (p Path) serviceFunctionParameters() string {
	if len(p.parameters) == 0 {
		return ""
	}

	sb := strings.Builder{}
	for i, param := range p.parameters {
		fmt.Fprintf(&sb, "%s %s", param.parameterType.ToVariableType(), param.parameterType.Name())
		if i < len(p.parameters)-1 {
			sb.WriteString(", ")
		}
	}

	return sb.String()
}

func (p Path) serviceFunctionNetReqURL() string {
	paramsInURL := 0
	finalRoute := p.route
	routeReplacements := "this.config.BasePath"
	for _, param := range p.parameters {
		if param.location == PathParameterLocation {
			finalRoute = strings.Replace(finalRoute, "{"+param.name+"}", fmt.Sprintf("{%d}", paramsInURL+1), 1)
			routeReplacements = routeReplacements + ", " + param.name
			paramsInURL++
		}
	}
	return fmt.Sprintf("string.Format(\"{0}%s\", %s)", finalRoute, routeReplacements)
}

// ServiceFunction generates C# code that is used to make network requests
func (p Path) ServiceFunction(knownModifiers []security.Auth) string {
	builder := strings.Builder{}

	fmt.Fprintf(&builder, "public %s %s(%s)\n{\n", p.reqPathName(), p.operationID, p.serviceFunctionParameters())
	fmt.Fprintf(&builder, "\tvar unityNetworkReq = new UnityWebRequest(%s, %s);\n", p.serviceFunctionNetReqURL(), p.toUnityHTTPVerb())
	if len(p.security) == 1 {
		fmt.Fprintf(&builder, "\t%s\n", p.guard(p.security[0], knownModifiers).ModifyNetworkRequest())
	} else if len(p.security) > 1 {
		for _, sec := range p.security {
			fmt.Fprintf(&builder, "\tif (string.IsNullOrEmpty(this.config.Security.%s()) == false) {\n", sec.Identifier)
			fmt.Fprintf(&builder, "\t\t%s\n", p.guard(sec, knownModifiers).ModifyNetworkRequest())
			builder.WriteString("\t}\n")
		}
	}
	fmt.Fprintf(&builder, "\treturn new %s(unityNetworkReq);\n}", p.reqPathName())

	return builder.String()
}
