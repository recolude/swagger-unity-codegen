package path

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/recolude/swagger-unity-codegen/unitygen/convention"
	"github.com/recolude/swagger-unity-codegen/unitygen/security"
	"github.com/recolude/swagger-unity-codegen/unitygen/unity"
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

func (p Path) Parameters() []Parameter {
	return p.parameters
}

func (p Path) OperationID() string {
	return p.operationID
}

func (p Path) Route() string {
	return p.route
}

func (p Path) Method() string {
	return p.httpMethod
}

func (p Path) SecurityReferences() []SecurityMethodReference {
	return p.security
}

func (p Path) Responses() map[string]Response {
	return p.responses
}

// Tags are what different tags a route is associated with
func (p Path) Tags() []string {
	return p.tags
}

func (p Path) reqPathName() string {
	return fmt.Sprintf("%sUnityWebRequest", p.operationID)
}

func (p Path) respVariableName(k string) string {

	switch k {
	case "200":
		return "success"

	case "401":
		return "unauthorized"
	case "403":
		return "forbidden"
	case "404":
		return "notFound"

	case "500":
		return "internalServerError"
	case "501":
		return "notImplemented"
	case "502":
		return "badGateway"
	case "503":
		return "serviceUnavailable"
	case "504":
		return "gatewayTimeout"

	case "default":
		return "fallbackResponse"

	}

	panic("unkown response key: " + k)
}

func (p Path) renderConditionalResponseCast(code string, resp Response) string {
	if resp.schema == nil {
		panic(fmt.Sprintf("code %s has nil response schema", code))
	}

	if parsed, err := strconv.Atoi(code); err == nil {
		return fmt.Sprintf("if (UnderlyingRequest.responseCode == %d) {\n\t\t\t%s = JsonUtility.FromJson<%s>(UnderlyingRequest.downloadHandler.text);\n\t\t}", parsed, p.respVariableName(code), resp.schema.ToVariableType())
	}
	return fmt.Sprintf("fallbackResponse = JsonUtility.FromJson<%s>(UnderlyingRequest.downloadHandler.text);", resp.schema.ToVariableType())
}

func (p Path) renderHandleResponse() string {
	keys := make([]string, 0)
	for k, resp := range p.responses {
		if resp.schema != nil {
			keys = append(keys, k)
		}
	}
	builder := strings.Builder{}
	sort.Strings(keys)
	if len(keys) > 0 {
		builder.WriteString("\t\t")
	}

	if len(keys) == 1 {
		builder.WriteString(p.renderConditionalResponseCast(keys[0], p.responses[keys[0]]))
		builder.WriteString("\n")
	} else if len(keys) > 1 {
		for keyIndex := range keys {
			if keys[keyIndex] == "default" {
				builder.WriteString("{\n\t\t\t")
				builder.WriteString(p.renderConditionalResponseCast(keys[keyIndex], p.responses[keys[keyIndex]]))
				builder.WriteString("\n\t\t}")
			} else {
				builder.WriteString(p.renderConditionalResponseCast(keys[keyIndex], p.responses[keys[keyIndex]]))
			}
			if keyIndex < len(keys)-1 {
				builder.WriteString(" else ")
			}
		}
		builder.WriteString("\n")
	}
	return builder.String()
}

// SupportingClasses will write out different helper classes in C# to assist
// in network requests
func (p Path) SupportingClasses() string {
	builder := strings.Builder{}

	fmt.Fprintf(&builder, "public class %s {\n\n", p.reqPathName())

	// Outline all portential responses
	keys := make([]string, 0)
	for k := range p.responses {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		// Some responses are defined as an empty body, making them nil!
		if p.responses[k].schema != nil {
			fmt.Fprintf(&builder, "\tpublic %s %s;\n\n", p.responses[k].schema.ToVariableType(), p.respVariableName(k))
		}
	}

	// underlying network request
	fmt.Fprint(&builder, "\tpublic UnityWebRequest UnderlyingRequest{ get; }\n\n")

	// constructor
	fmt.Fprintf(&builder, "\tpublic %s(UnityWebRequest req) {\n\t\tthis.UnderlyingRequest = req;\n\t}\n\n", p.reqPathName())

	// Function that will actually execute the request
	builder.WriteString("\tpublic IEnumerator Run() {\n\t\tyield return this.UnderlyingRequest.SendWebRequest();\n")
	builder.WriteString(p.renderHandleResponse())
	builder.WriteString("\t}\n\n}")
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
	routeReplacements := "this.Config.BasePath"

	// Do all path parameters first
	for _, param := range p.parameters {
		if param.location == PathParameterLocation {
			finalRoute = strings.Replace(finalRoute, "{"+param.name+"}", fmt.Sprintf("{%d}", paramsInURL+1), 1)
			routeReplacements += ", "
			if param.parameterType.ToVariableType() == "string" {
				routeReplacements += fmt.Sprintf("UnityWebRequest.EscapeURL(%s)", param.name)
			} else {
				routeReplacements += param.name
			}
			paramsInURL++
		}
	}

	// Then do query parameters next
	firstQuery := true
	for _, param := range p.parameters {
		if param.location == QueryParameterLocation {
			if firstQuery {
				finalRoute += "?"
				firstQuery = false
			} else {
				finalRoute += "&"
			}
			finalRoute += fmt.Sprintf("%s={%d}", param.name, paramsInURL+1)

			routeReplacements += ", "
			if param.parameterType.ToVariableType() == "string" {
				routeReplacements += fmt.Sprintf("UnityWebRequest.EscapeURL(%s)", param.name)
			} else {
				routeReplacements += param.name
			}
			paramsInURL++
		}
	}

	return fmt.Sprintf("string.Format(\"{0}%s\", %s)", finalRoute, routeReplacements)
}

// ServiceFunction generates C# code that is used to make network requests
func (p Path) ServiceFunction(knownModifiers []security.Auth) string {
	builder := strings.Builder{}

	fmt.Fprintf(&builder, "public %s %s(%s)\n{\n", p.reqPathName(), p.operationID, p.serviceFunctionParameters())
	fmt.Fprintf(&builder, "\tvar unityNetworkReq = new UnityWebRequest(%s, %s);\n", p.serviceFunctionNetReqURL(), unity.ToUnityHTTPVerb(p.httpMethod))

	if len(p.responses) > 0 {
		builder.WriteString("\tunityNetworkReq.downloadHandler = new DownloadHandlerBuffer();\n")
	}

	if len(p.security) == 1 {
		fmt.Fprintf(&builder, "\t%s\n", p.guard(p.security[0], knownModifiers).ModifyNetworkRequest())
	} else if len(p.security) > 1 {
		for _, sec := range p.security {
			fmt.Fprintf(&builder, "\tif (string.IsNullOrEmpty(this.Config.%s) == false) {\n", convention.TitleCase(sec.Identifier))
			fmt.Fprintf(&builder, "\t\t%s\n", p.guard(sec, knownModifiers).ModifyNetworkRequest())
			builder.WriteString("\t}\n")
		}
	}
	fmt.Fprintf(&builder, "\treturn new %s(unityNetworkReq);\n}", p.reqPathName())

	return builder.String()
}
