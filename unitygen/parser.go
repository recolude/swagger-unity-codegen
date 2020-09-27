package unitygen

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/Jeffail/gabs/v2"
	"github.com/recolude/swagger-unity-codegen/unitygen/definition"
	"github.com/recolude/swagger-unity-codegen/unitygen/property"
)

// Parser reads through a file and interprets a swagger definition
type Parser struct{}

func (p Parser) interpretArrayProperty(path []string, name string, obj *gabs.Container) (property.Array, error) {
	items := obj.Path("items")
	if items == nil {
		return property.Array{}, InvalidSpecError{Path: path, Reason: "Unable to find array type (missing items property)"}
	}

	prop, err := p.interpretObjectPropertyDefinition(append(path, "items"), "", items)
	if err != nil {
		return property.Array{}, err
	}

	return property.NewArray(name, prop), nil
}

func (p Parser) interpretStringProperty(path []string, name string, obj *gabs.Container) (property.String, error) {
	format := ""
	formatInSpec := obj.Path("format").Data()
	if formatInSpec != nil {
		format = formatInSpec.(string)
	}
	return property.NewString(name, format), nil
}

func (p Parser) interpretObjectPropertyDefinition(path []string, name string, obj *gabs.Container) (property.Property, error) {

	objRef, ok := obj.Path("$ref").Data().(string)
	if ok {
		return property.NewObjectReference(name, objRef), nil
	}

	propType, ok := obj.Path("type").Data().(string)
	if !ok {
		return nil, InvalidSpecError{Path: append(path, name), Reason: "Property type not found on definition"}
	}

	switch propType {

	case "string":
		return p.interpretStringProperty(path, name, obj)

	case "array":
		return p.interpretArrayProperty(path, name, obj)

	default:
		return nil, InvalidSpecError{Path: append(path, name), Reason: fmt.Sprintf("unknown property type \"%s\"", propType)}
	}
}

func (p Parser) interpretObjectDefinition(path []string, name string, obj *gabs.Container) (definition.Object, error) {
	newPath := append(path, name)
	if obj == nil {
		return definition.Object{}, InvalidSpecError{Path: newPath, Reason: "Definition contains no contents"}
	}
	properties := make([]property.Property, 0)

	for key, val := range obj.Path("properties").ChildrenMap() {
		prop, err := p.interpretObjectPropertyDefinition(append(newPath, "properties"), key, val)
		if err != nil {
			return definition.Object{}, err
		}
		properties = append(properties, prop)
	}

	return definition.NewObject(name, properties), nil
}

func (p Parser) interpretStringDefinition(path []string, name string, obj *gabs.Container) (definition.Definition, error) {

	enum := obj.Path("enum")
	if enum == nil {
		return nil, InvalidSpecError{Path: append(path, name), Reason: "Unimplemented string case"}
	}

	children := enum.Children()
	parsedValues := make([]string, len(children))
	for i, child := range enum.Children() {
		parsedValues[i] = child.Data().(string)
	}

	return definition.NewEnum(name, parsedValues), nil
}

// ParseJSON reads through the input stream and constructs an understanding of
// the API our Unity3D client needs to interact with
func (p Parser) ParseJSON(in io.Reader) (Spec, error) {

	entireIn, err := ioutil.ReadAll(in)
	if err != nil {
		return Spec{}, err
	}

	jsonParsed, err := gabs.ParseJSON(entireIn)
	if err != nil {
		return Spec{}, err
	}

	info := SpecInfo{
		Title:   jsonParsed.Path("info.title").Data().(string),
		Version: jsonParsed.Path("info.version").Data().(string),
	}

	definitions := make([]definition.Definition, 0)
	for key, val := range jsonParsed.Path("definitions").ChildrenMap() {

		definitionType, ok := val.Path("type").Data().(string)
		if !ok {
			return Spec{}, InvalidSpecError{Path: []string{"definitions", key}, Reason: "Definition type not found on definition"}
		}

		var def definition.Definition
		switch definitionType {
		case "object":
			def, err = p.interpretObjectDefinition([]string{"definitions"}, key, val)
			break

		case "string":
			def, err = p.interpretStringDefinition([]string{"definitions"}, key, val)
			break

		default:
			return Spec{}, InvalidSpecError{Path: []string{"definitions", key, "type"}, Reason: fmt.Sprintf("Unknown definition type \"%s\"", definitionType)}
		}

		if err != nil {
			return Spec{}, err
		}
		definitions = append(definitions, def)
	}

	return NewSpec(info, definitions), nil
}
