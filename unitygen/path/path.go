package path

import (
	"errors"
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
	summary     string
	description string
	httpMethod  string
	tags        []string
	operationID string
	security    []SecurityMethodReference

	// mapping of HTTP Status Codes to responses
	responses map[string]Response

	parameters []Parameter
}

// NewPath creates a new path
func NewPath(route, summary, description, operationID, method string, tags []string, security []SecurityMethodReference, responses map[string]Response, parameters []Parameter) Path {
	bodyFound := false
	for _, p := range parameters {
		if p.location == BodyParameterLocation {
			if bodyFound {
				panic(errors.New("can not have multiple body parameters for a single path"))
			}
			bodyFound = true
		}
	}

	return Path{
		route:       route,
		summary:     strings.TrimSpace(summary),
		description: strings.TrimSpace(description),
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

func (p Path) unityWebReqPathName() string {
	return fmt.Sprintf("%sRequest", p.OperationIDFunctionName())
}

func (p Path) requestParamClassName() string {
	return fmt.Sprintf("%sRequestParams", p.OperationIDFunctionName())
}

func (p Path) respVariableName(k string) string {
	switch k {
	case "200":
		return "Success"

	case "400":
		return "BadRequest"
	case "401":
		return "Unauthorized"
	case "403":
		return "Forbidden"
	case "404":
		return "NotFound"

	case "500":
		return "InternalServerError"
	case "501":
		return "NotImplemented"
	case "502":
		return "BadGateway"
	case "503":
		return "ServiceUnavailable"
	case "504":
		return "GatewayTimeout"

	case "default":
		return "FallbackResponse"

	}

	panic("unknown response key: " + k)
}

func (p Path) queryParamCount() int {
	count := 0
	for _, param := range p.parameters {
		if param.location == QueryParameterLocation {
			count++
		}
	}
	return count
}

func (p Path) bodyParam() *Parameter {
	for _, param := range p.parameters {
		if param.location == BodyParameterLocation {
			return &param
		}
	}
	return nil
}

// RequestParamClass a class to act as a container for all parameters associated
// for making a specific web request
func (p Path) RequestParamClass() string {
	if len(p.parameters) == 0 {
		return ""
	}

	builder := strings.Builder{}

	fmt.Fprintf(&builder, "public class %s\n{\n", p.requestParamClassName())

	for _, param := range p.parameters {
		privateVarName := convention.CamelCase(param.name)
		propertyName := convention.TitleCase(param.name)
		fmt.Fprintf(&builder, "\tprivate bool %sSet = false;\n", privateVarName)
		fmt.Fprintf(&builder, "\tprivate %s %s;\n", param.parameterType.ToVariableType(), privateVarName)
		fmt.Fprintf(
			&builder,
			"\tpublic %s %s { get { return %s; } set { %sSet = true; %s = value; } }\n",
			param.parameterType.ToVariableType(),
			propertyName,
			privateVarName,
			privateVarName,
			privateVarName,
		)
		fmt.Fprintf(&builder, "\tpublic void Unset%s() { %s = %s; %sSet = false; }\n\n", propertyName, privateVarName, param.parameterType.EmptyValue(), privateVarName)
	}

	builder.WriteString("\tpublic UnityWebRequest BuildUnityWebRequest(string baseURL)\n\t{\n")

	fmt.Fprintf(&builder, "\t\tvar finalPath = baseURL + \"%s\";\n", p.route)

	// Do all path parameters first
	for _, param := range p.parameters {
		if param.location == PathParameterLocation {
			privateVarName := convention.CamelCase(param.name)
			fmt.Fprintf(&builder, "\t\tfinalPath = finalPath.Replace(\"{%s}\", %sSet ? UnityWebRequest.EscapeURL(%s.ToString()) : \"\");\n", param.name, privateVarName, privateVarName)
		}
	}

	if p.queryParamCount() > 0 {
		builder.WriteString("\t\tvar queryAdded = false;\n\n")
	}

	// Build out the final url by appending set query params to the url.
	for _, param := range p.parameters {
		if param.location == QueryParameterLocation {
			privateVarName := convention.CamelCase(param.name)
			fmt.Fprintf(&builder, "\t\tif (%sSet) {\n", privateVarName)
			fmt.Fprintf(&builder, "\t\t\tfinalPath += (queryAdded ? \"&\" : \"?\") + \"%s=\";\n", param.name)
			builder.WriteString("\t\t\tqueryAdded = true;\n")
			fmt.Fprintf(&builder, "\t\t\tfinalPath += UnityWebRequest.EscapeURL(%s.ToString());\n\t\t}\n\n", privateVarName)
		}
	}

	// Build the unity request object
	fmt.Fprintf(&builder, "\t\tvar unityWebReq = new UnityWebRequest(finalPath, %s);\n", unity.ToUnityHTTPVerb(p.httpMethod))

	// Set the body of the request
	bodyParam := p.bodyParam()
	if bodyParam != nil {
		fmt.Fprintf(&builder, "\t\tvar unityRawUploadHandler = new UploadHandlerRaw(Encoding.Unicode.GetBytes(JsonConvert.SerializeObject(%s)));\n", convention.CamelCase(bodyParam.name))
		builder.WriteString("\t\tunityRawUploadHandler.contentType = \"application/json\";\n")
		builder.WriteString("\t\tunityWebReq.uploadHandler = unityRawUploadHandler;\n")
	}

	// Return result
	builder.WriteString("\t\treturn unityWebReq;\n\t}\n}")

	return builder.String()
}

// UnityWebRequest is the CSharp code that handles asynchronous web requests
func (p Path) UnityWebRequest() string {
	builder := strings.Builder{}

	fmt.Fprintf(&builder, "public class %s : I%s {\n\n", p.unityWebReqPathName(), p.unityWebReqPathName())

	// Outline all portential responses
	keys := make([]string, 0)
	for k := range p.responses {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		// Some responses are defined as an empty body, making them nil!
		if p.responses[k] != nil {
			desc := strings.TrimSpace(p.responses[k].Description())
			lines := strings.Split(strings.ReplaceAll(desc, "\r\n", "\n"), "\n")

			validLines := false
			for _, line := range lines {
				if line != "" {
					validLines = true
				}
			}

			if validLines {
				fmt.Fprint(&builder, "\t/// <summary>\n")
				for _, line := range lines {
					if line != "" {
						fmt.Fprintf(&builder, "\t/// %s\n", line)
					}
				}
				fmt.Fprint(&builder, "\t/// </summary>\n")
			}

			fmt.Fprintf(&builder, "\tpublic %s %s { get; private set; }\n\n", p.responses[k].VariableType(), p.respVariableName(k))
		}
	}

	// underlying network request
	fmt.Fprint(&builder, "\tpublic UnityWebRequest UnderlyingRequest{ get; }\n\n")

	// constructor
	fmt.Fprintf(&builder, "\tpublic %s(UnityWebRequest req) {\n\t\tthis.UnderlyingRequest = req;\n\t}\n\n", p.unityWebReqPathName())

	// Function that will actually execute the request
	builder.WriteString("\tpublic IEnumerator Run() {\n\t\tyield return this.UnderlyingRequest.SendWebRequest();\n")
	if len(p.responses) > 0 {
		builder.WriteString("\t\tInterpret(this.UnderlyingRequest);\n")
	}
	builder.WriteString("\t}\n\n")

	if len(p.responses) > 0 {
		builder.WriteString("\tpublic void Interpret(UnityWebRequest req) {\n")
		builder.WriteString(p.renderHandleResponse())
		builder.WriteString("\t}\n\n")
	}

	builder.WriteString("}")

	return builder.String()
}

func (p Path) renderConditionalResponseCast(code string, resp Response) string {
	if parsed, err := strconv.Atoi(code); err == nil {
		if resp != nil {
			return fmt.Sprintf("if (req.responseCode == %d) {\n\t\t\t%s\n\t\t}", parsed, resp.Interpret(p.respVariableName(code), "req.downloadHandler"))
		} else {
			return fmt.Sprintf("if (req.responseCode == %d) {\n\t\t\t%s\n\t\t}", parsed, "// No expected response. Do nothing!")
		}
	}
	return resp.Interpret("fallbackResponse", "req.downloadHandler")
}

func (p Path) renderHandleResponse() string {
	codes := make([]string, 0)
	for k := range p.responses {
		codes = append(codes, k)
	}
	builder := strings.Builder{}
	sort.Strings(codes)
	if len(codes) > 0 {
		builder.WriteString("\t\t")
	}

	if len(codes) == 1 {
		builder.WriteString(p.renderConditionalResponseCast(codes[0], p.responses[codes[0]]))
		builder.WriteString("\n")
	} else if len(codes) > 1 {
		for codeIndex, code := range codes {
			if code == "default" {
				builder.WriteString("{\n\t\t\t")
				builder.WriteString(p.renderConditionalResponseCast(code, p.responses[codes[codeIndex]]))
				builder.WriteString("\n\t\t}")
			} else {
				builder.WriteString(p.renderConditionalResponseCast(code, p.responses[codes[codeIndex]]))
			}
			if codeIndex < len(codes)-1 {
				builder.WriteString(" else ")
			}
		}
		builder.WriteString("\n")
	}
	return builder.String()
}

func (p Path) Interface() string {
	builder := strings.Builder{}
	fmt.Fprintf(&builder, "public interface I%s : IWebRequest", p.unityWebReqPathName())

	// Outline all portential responses
	keys := make([]string, 0)
	for k := range p.responses {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		// Some responses are defined as an empty body, making them nil!
		if p.responses[k] != nil {

			fmt.Fprintf(&builder, ", I%sResponse<%s>", p.respVariableName(k), p.responses[k].VariableType())
		}
	}

	builder.WriteString(" { }")

	return builder.String()
}

// SupportingClasses will write out different helper classes in C# to assist
// in network requests
func (p Path) SupportingClasses() string {
	builder := strings.Builder{}
	builder.WriteString(p.Interface())
	builder.WriteString("\n")
	builder.WriteString(p.UnityWebRequest())
	builder.WriteString("\n")
	builder.WriteString(p.RequestParamClass())
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
		fmt.Fprintf(&sb, "%s %s", param.parameterType.ToVariableType(), convention.CamelCase(param.name))
		if i < len(p.parameters)-1 {
			sb.WriteString(", ")
		}
	}

	return sb.String()
}

func (p Path) functionOveride(knownModifiers []security.Auth) string {
	builder := strings.Builder{}

	builder.WriteString(p.serviceFunctionSummary())
	fmt.Fprintf(&builder, "public I%s %s(%s)\n{\n", p.unityWebReqPathName(), p.OperationIDFunctionName(), p.serviceFunctionParameters())
	fmt.Fprintf(&builder, "\treturn %s(new %s() {\n", p.OperationIDFunctionName(), p.requestParamClassName())
	// fmt.Fprintf(&builder, "\tvar unityNetworkReq = new UnityWebRequest(%s, %s);\n", p.serviceFunctionNetReqURL(), unity.ToUnityHTTPVerb(p.httpMethod))

	for _, param := range p.parameters {
		privateVarName := convention.CamelCase(param.name)
		propertyName := convention.TitleCase(param.name)
		fmt.Fprintf(&builder, "\t\t%s=%s,\n", propertyName, privateVarName)
	}
	builder.WriteString("\t});\n}")

	return builder.String()
}

func (p Path) OperationIDFunctionName() string {
	return convention.ClassName(p.operationID)
}

func (p Path) ServiceInterfaceFunction(extra string) string {
	return fmt.Sprintf("public %s I%s %s(%s requestParams);", extra, p.unityWebReqPathName(), p.OperationIDFunctionName(), p.requestParamClassName())
}

func (p Path) serviceFunctionSummary() string {
	hasSummary := p.summary != ""
	hasDescription := p.description != ""
	if !hasSummary && !hasDescription {
		return ""
	}

	builder := strings.Builder{}
	builder.WriteString("/// <summary>\n")
	if hasSummary && hasDescription {
		fmt.Fprintf(&builder, "/// %s\n", p.summary)
		builder.WriteString("/// \n")
		fmt.Fprintf(&builder, "/// %s\n", p.description)
	} else {
		var text = p.summary
		if hasDescription {
			text = p.description
		}
		fmt.Fprintf(&builder, "/// %s\n", text)

	}

	builder.WriteString("/// </summary>\n")

	return builder.String()
}

// ServiceFunction generates C# code that is used to make network requests
func (p Path) ServiceFunction(knownModifiers []security.Auth) string {
	builder := strings.Builder{}

	builder.WriteString(p.serviceFunctionSummary())

	if len(p.parameters) > 0 {
		fmt.Fprintf(&builder, "public override I%s %s(%s requestParams)\n{\n", p.unityWebReqPathName(), p.OperationIDFunctionName(), p.requestParamClassName())
		builder.WriteString("\tvar unityNetworkReq = requestParams.BuildUnityWebRequest(this.Config.BasePath);\n")
	} else {
		fmt.Fprintf(&builder, "public I%s %s()\n{\n", p.unityWebReqPathName(), p.OperationIDFunctionName())
		fmt.Fprintf(&builder, "\tvar unityNetworkReq = new UnityWebRequest(string.Format(\"{0}%s\", this.Config.BasePath), %s);\n", p.route, unity.ToUnityHTTPVerb(p.httpMethod))
	}

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
	fmt.Fprintf(&builder, "\treturn new %s(unityNetworkReq);\n}", p.unityWebReqPathName())

	if len(p.parameters) > 0 {
		fmt.Fprintf(&builder, "\n\n%s", p.functionOveride(knownModifiers))
	}

	return builder.String()
}
