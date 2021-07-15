package unitygen

import (
	"fmt"
	"io"
	"io/ioutil"
	"sort"
	"strings"

	"github.com/Jeffail/gabs/v2"
	"github.com/recolude/swagger-unity-codegen/unitygen/convention"
	"github.com/recolude/swagger-unity-codegen/unitygen/model"
	"github.com/recolude/swagger-unity-codegen/unitygen/model/property"
	"github.com/recolude/swagger-unity-codegen/unitygen/path"
	"github.com/recolude/swagger-unity-codegen/unitygen/security"
)

// Parser reads through a file and interprets a swagger definition
type Parser struct {
	workingDefinitions map[string]*model.DefinitionWrapper
}

func NewParser() *Parser {
	return &Parser{
		workingDefinitions: make(map[string]*model.DefinitionWrapper),
	}
}

func (p *Parser) interpretArrayProperty(path []string, objectName, propertyName string, obj *gabs.Container) (property.Array, error) {
	items := obj.Path("items")
	if items == nil {
		return property.Array{}, InvalidSpecError{Path: path, Reason: "Unable to find array type (missing items property)"}
	}

	prop, err := p.interpretObjectDefinitionProperty(append(path, "items"), objectName, propertyName, items)
	if err != nil {
		return property.Array{}, err
	}

	return property.NewArray(propertyName, prop), nil
}

func (p *Parser) interpretStringProperty(path []string, name string, obj *gabs.Container) (property.String, error) {
	format := ""
	formatInSpec := obj.Path("format").Data()
	if formatInSpec != nil {
		format = formatInSpec.(string)
	}
	return property.NewString(name, format), nil
}

func (p *Parser) interpretIntProperty(path []string, name string, obj *gabs.Container) (property.Integer, error) {
	format := ""
	formatInSpec := obj.Path("format").Data()
	if formatInSpec != nil {
		format = formatInSpec.(string)
	}
	return property.NewInteger(name, format), nil
}

func (p *Parser) interpretNumberProperty(path []string, name string, obj *gabs.Container) (property.Number, error) {
	format := ""
	formatInSpec := obj.Path("format").Data()
	if formatInSpec != nil {
		format = formatInSpec.(string)
	}
	return property.NewNumber(name, format), nil
}

func (p *Parser) interpretBooleanProperty(name string) (property.Boolean, error) {
	return property.NewBoolean(name), nil
}

func (p *Parser) interpretNestedObjectProperty(path []string, objectName, propertyName string, obj *gabs.Container) (model.Property, error) {
	def, err := p.interpretObjectDefinition(path, fmt.Sprintf("%s%s", objectName, convention.ClassName(propertyName)), obj)
	return property.NewObject(propertyName, def), err
}

func (p *Parser) interpretObjectDefinitionProperty(path []string, objectName, propertyName string, obj *gabs.Container) (model.Property, error) {
	objRefUrl, ok := obj.Path("$ref").Data().(string)

	if ok {
		_, hasDef := p.workingDefinitions[objRefUrl]

		if hasDef == false {
			p.workingDefinitions[objRefUrl] = model.NewDefinitionWrapper(nil)
		}

		return property.NewDefinitionReference(propertyName, p.workingDefinitions[objRefUrl]), nil
	}

	propType, ok := obj.Path("type").Data().(string)
	if !ok {
		return nil, InvalidSpecError{Path: append(path, propertyName), Reason: "Property type not found on definition"}
	}

	switch propType {

	case "string":
		return p.interpretStringProperty(path, propertyName, obj)

	case "array":
		return p.interpretArrayProperty(path, objectName, propertyName, obj)

	case "integer":
		return p.interpretIntProperty(path, propertyName, obj)

	case "number":
		return p.interpretNumberProperty(path, propertyName, obj)

	case "boolean":
		return p.interpretBooleanProperty(propertyName)

	case "object":
		return p.interpretNestedObjectProperty(path, objectName, propertyName, obj)

	default:
		return nil, InvalidSpecError{Path: append(path, propertyName), Reason: fmt.Sprintf("unknown property type \"%s\"", propType)}
	}
}

func (p *Parser) interpretObjectDefinition(path []string, objectName string, obj *gabs.Container) (model.Object, error) {
	newPath := append(path, objectName)
	if obj == nil {
		return model.Object{}, InvalidSpecError{Path: newPath, Reason: "Definition contains no contents"}
	}
	properties := make([]model.Property, 0)

	for propertyName, val := range obj.Path("properties").ChildrenMap() {
		prop, err := p.interpretObjectDefinitionProperty(append(newPath, "properties"), objectName, propertyName, val)
		if err != nil {
			return model.Object{}, err
		}
		properties = append(properties, prop)
	}

	return model.NewObject(objectName, properties), nil
}

func (p *Parser) interpretStringDefinition(path []string, name string, obj *gabs.Container) (model.Definition, error) {
	enum := obj.Path("enum")
	if enum == nil {
		return nil, InvalidSpecError{Path: append(path, name), Reason: "Unimplemented string case"}
	}

	children := enum.Children()
	parsedValues := make([]string, len(children))
	for i, child := range enum.Children() {
		parsedValues[i] = child.Data().(string)
	}

	return model.NewStringEnum(name, parsedValues), nil
}

func (p *Parser) interpretNumberDefinition(path []string, name string, obj *gabs.Container) (model.Definition, error) {
	enum := obj.Path("enum")
	if enum == nil {
		return nil, InvalidSpecError{Path: append(path, name), Reason: "Unimplemented number case"}
	}

	children := enum.Children()
	parsedValues := make([]float64, len(children))
	for i, child := range enum.Children() {
		parsedValues[i] = child.Data().(float64)
	}

	return model.NewNumberEnum(name, parsedValues), nil
}

func (p *Parser) parseDefinitions(obj *gabs.Container) ([]model.Definition, error) {
	definitions := make([]model.Definition, 0)
	var err error
	for key, val := range obj.Path("definitions").ChildrenMap() {

		definitionType, ok := val.Path("type").Data().(string)
		if !ok {
			return nil, InvalidSpecError{Path: []string{"definitions", key}, Reason: "Definition type not found on definition"}
		}

		var def model.Definition
		switch definitionType {
		case "object":
			def, err = p.interpretObjectDefinition([]string{"definitions"}, key, val)
			break

		case "string":
			def, err = p.interpretStringDefinition([]string{"definitions"}, key, val)
			break

		case "number":
			def, err = p.interpretNumberDefinition([]string{"definitions"}, key, val)
			break

		default:
			return nil, InvalidSpecError{Path: []string{"definitions", key, "type"}, Reason: fmt.Sprintf("Unknown definition type \"%s\"", definitionType)}
		}

		if err != nil {
			return nil, err
		}
		if wrapper, ok := p.workingDefinitions["#/definitions/"+key]; ok {
			wrapper.UpdateDefinition(def)
		} else {
			p.workingDefinitions["#/definitions/"+key] = model.NewDefinitionWrapper(def)
		}
		definitions = append(definitions, def)
	}

	return definitions, nil
}

func (p *Parser) interpretAPIKeyDefinition(path []string, name string, obj *gabs.Container) (security.Auth, error) {
	keyPath := append(path, name)
	if obj == nil {
		return nil, InvalidSpecError{Path: keyPath, Reason: "Definition contains no contents"}
	}

	apikeyName, ok := obj.Path("name").Data().(string)
	if !ok || apikeyName == "" {
		return nil, InvalidSpecError{Path: keyPath, Reason: "No name found for key"}
	}

	foundIn, ok := obj.Path("in").Data().(string)
	if !ok || foundIn == "" {
		return nil, InvalidSpecError{Path: keyPath, Reason: "No destination for API Key found"}
	}

	var keyLoc security.APIKeyLocation
	switch foundIn {
	case "header":
		keyLoc = security.Header
		break

	default:
		return nil, InvalidSpecError{Path: keyPath, Reason: fmt.Sprintf("Unimplemnted key location: \"%s\"", foundIn)}
	}

	return security.NewAPIKey(name, apikeyName, keyLoc), nil
}

func (p *Parser) parseSecurityDefinitions(obj *gabs.Container) ([]security.Auth, error) {
	definitions := make([]security.Auth, 0)
	var err error
	for key, val := range obj.Path("securityDefinitions").ChildrenMap() {

		definitionType, ok := val.Path("type").Data().(string)
		if !ok {
			return nil, InvalidSpecError{Path: []string{"securityDefinitions", key}, Reason: "Definition type not found on definition"}
		}

		var def security.Auth
		switch definitionType {
		case "apiKey":
			def, err = p.interpretAPIKeyDefinition([]string{"securityDefinitions"}, key, val)
			break

		default:
			return nil, InvalidSpecError{Path: []string{"securityDefinitions", key, "type"}, Reason: fmt.Sprintf("Unknown security type \"%s\"", definitionType)}
		}

		if err != nil {
			return nil, err
		}
		definitions = append(definitions, def)
	}

	return definitions, nil
}

func (p *Parser) interpretPathPameterProperty(path []string, name string, obj *gabs.Container) (model.Property, error) {
	schemaNode := obj.Path("schema")
	if schemaNode != nil {
		refNode := schemaNode.Path("$ref")
		if refNode == nil {
			return nil, InvalidSpecError{Path: append(path, "schema"), Reason: "Expected $ref"}
		}
		refName, hasRef := refNode.Data().(string)
		if hasRef {
			_, hasDef := p.workingDefinitions[refName]

			if hasDef == false {
				p.workingDefinitions[refName] = model.NewDefinitionWrapper(nil)
			}

			return property.NewDefinitionReference(name, p.workingDefinitions[refName]), nil
		}

		return nil, InvalidSpecError{Path: append(path, "schema"), Reason: "Expected $ref to be string"}
	}

	propType, ok := obj.Path("type").Data().(string)
	if !ok {
		return nil, InvalidSpecError{Path: append(path, name), Reason: "Property type not found on definition"}
	}

	switch propType {

	case "string":
		return p.interpretStringProperty(path, name, obj)

	case "array":
		return p.interpretArrayProperty(path, "", name, obj)

	case "integer":
		return p.interpretIntProperty(path, name, obj)

	case "number":
		return p.interpretNumberProperty(path, name, obj)

	case "boolean":
		return p.interpretBooleanProperty(name)

	default:
		return nil, InvalidSpecError{Path: append(path, name), Reason: fmt.Sprintf("unknown property type \"%s\"", propType)}
	}
}

func (p *Parser) parsePaths(url string, routeObj *gabs.Container) ([]path.Path, error) {
	paths := make([]path.Path, 0)
	for verb, verbObj := range routeObj.ChildrenMap() {
		tagsInJSON := make([]string, 0)
		for _, child := range verbObj.Path("tags").Children() {
			tagsInJSON = append(tagsInJSON, child.Data().(string))
		}

		securityReferences := make([]path.SecurityMethodReference, 0)
		for _, child := range verbObj.Path("security").Children() {
			for key := range child.ChildrenMap() {
				securityReferences = append(securityReferences, path.NewSecurityMethodReference(key))
			}
		}
		sort.Sort(sortBySecurityReferenceIdentifier(securityReferences))

		operationID, ok := verbObj.Path("operationId").Data().(string)
		if !ok {
			return nil, InvalidSpecError{Path: []string{"paths", url, verb}, Reason: "unable to locate operation ID"}
		}

		responses := make(map[string]path.Response)
		for code, respJSON := range verbObj.Path("responses").ChildrenMap() {
			schemaJSON := respJSON.Path("schema")
			description := ""
			descriptionNode := respJSON.Path("description")
			if descriptionNode != nil {
				description = descriptionNode.Data().(string)
			}

			if schemaJSON != nil {
				refNode := schemaJSON.Path("$ref")
				if refNode != nil {
					responses[code] = path.NewDefinitionResponse(
						description,
						model.NewDefinitionReference(refNode.Data().(string)),
					)
					continue
				}

				typeNode := schemaJSON.Path("type")
				if typeNode != nil {
					typeValue := typeNode.Data().(string)
					switch typeValue {
					case "file":
						responses[code] = path.NewFileResponse(description)
						break

					case "number":
						responses[code] = path.NewNumberResponse(description)
						break

					case "array":
						prop, err := p.interpretArrayProperty([]string{"paths", url, verb, code}, "", code, schemaJSON)
						if err != nil {
							return nil, err
						}
						responses[code] = path.NewArrayResponse(description, prop)
						break

					default:
						return nil, InvalidSpecError{Path: []string{"paths", url, verb, code}, Reason: "unable to interpret response schema: " + typeValue}
					}
					continue
				}

				return nil, InvalidSpecError{Path: []string{"paths", url, verb, code}, Reason: "unable to interpret response"}
			} else {
				responses[code] = nil
			}
		}

		parameters := make([]path.Parameter, 0)
		for paramIndex, param := range verbObj.Path("parameters").Children() {
			required, ok := param.Path("required").Data().(bool)
			if !ok {
				required = false
			}

			paramName := param.Path("name").Data().(string)

			paramProperty, err := p.interpretPathPameterProperty(
				[]string{url, verb, "parameters", fmt.Sprintf("[%d]", paramIndex)},
				paramName,
				param,
			)
			if err != nil {
				return nil, err
			}

			parameters = append(parameters, path.NewParameter(
				path.ParameterLocation(param.Path("in").Data().(string)),
				paramName,
				required,
				paramProperty,
			))
		}

		paths = append(paths,
			path.NewPath(
				url,
				operationID,
				strings.ToUpper(verb),
				tagsInJSON,
				securityReferences,
				responses,
				parameters,
			),
		)
	}

	sort.Sort(sortByPathMethod(paths))
	return paths, nil
}

type sortBySecurityReferenceIdentifier []path.SecurityMethodReference

func (a sortBySecurityReferenceIdentifier) Len() int      { return len(a) }
func (a sortBySecurityReferenceIdentifier) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a sortBySecurityReferenceIdentifier) Less(i, j int) bool {
	return a[i].Identifier < a[j].Identifier
}

type sortByPathMethod []path.Path

func (a sortByPathMethod) Len() int           { return len(a) }
func (a sortByPathMethod) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortByPathMethod) Less(i, j int) bool { return a[i].Method() < a[j].Method() }

func (p *Parser) serviceForPath(s Service, pa path.Path) bool {
	for _, t := range pa.Tags() {
		if t == s.name {
			return true
		}
	}
	return false
}

type sortByServiceName []Service

func (a sortByServiceName) Len() int           { return len(a) }
func (a sortByServiceName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortByServiceName) Less(i, j int) bool { return a[i].name < a[j].name }

func (p *Parser) parseServices(obj *gabs.Container) ([]Service, error) {
	services := make([]Service, 0)
	defaultServiceIndex := -1

	for key, val := range obj.Path("paths").ChildrenMap() {
		paths, err := p.parsePaths(key, val)
		if err != nil {
			return nil, err
		}
		for _, foundPath := range paths {

			// Have a default service for paths without any tags
			if len(foundPath.Tags()) == 0 {
				if defaultServiceIndex == -1 {
					defaultServiceIndex = len(services)
					services = append(services, Service{name: "Default"})
				}
				services[defaultServiceIndex].paths = append(services[defaultServiceIndex].paths, foundPath)
			}

			foundService := false
			for serviceIndex, service := range services {
				if p.serviceForPath(service, foundPath) {
					foundService = true
					services[serviceIndex].paths = append(services[serviceIndex].paths, foundPath)
				}
			}

			if foundService == false {
				for _, tag := range foundPath.Tags() {
					services = append(services, Service{name: tag})
					services[len(services)-1].paths = append(services[len(services)-1].paths, foundPath)
				}
			}
		}
	}

	sort.Sort(sortByServiceName(services))

	return services, nil
}

// ParseJSON reads through the input stream and constructs an understanding of
// the API our Unity3D client needs to interact with
func (p *Parser) ParseJSON(in io.Reader) (Spec, error) {
	entireIn, err := ioutil.ReadAll(in)
	if err != nil {
		return Spec{}, err
	}

	jsonParsed, err := gabs.ParseJSON(entireIn)
	if err != nil {
		return Spec{}, err
	}

	var info SpecInfo
	infoNode := jsonParsed.Path("info")
	if infoNode != nil {
		info = SpecInfo{
			Title:   jsonParsed.Path("info.title").Data().(string),
			Version: jsonParsed.Path("info.version").Data().(string),
		}
	}

	parsedDefinitions, err := p.parseDefinitions(jsonParsed)
	if err != nil {
		return Spec{}, err
	}

	parsedSecurityDefinitions, err := p.parseSecurityDefinitions(jsonParsed)
	if err != nil {
		return Spec{}, err
	}

	parsedServices, err := p.parseServices(jsonParsed)
	if err != nil {
		return Spec{}, err
	}

	return NewSpec(info, parsedDefinitions, parsedSecurityDefinitions, parsedServices), nil
}
