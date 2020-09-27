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

func (p Parser) interpretDefinition(path []string, name string, obj *gabs.Container) (definition.Definition, error) {
	newPath := append(path, name)
	if obj == nil {
		return nil, InvalidSpecError{Path: newPath, Reason: "Definition contains no contents"}
	}
	properties := make([]property.Property, 0)

	for key, val := range obj.Path("properties").ChildrenMap() {
		prop, err := p.interpretObjectPropertyDefinition(append(newPath, "properties"), key, val)
		if err != nil {
			return nil, err
		}
		properties = append(properties, prop)
	}

	return definition.NewObject(name, properties), nil
}

// Parse reads through the input stream and constructs an understanding of the
// API our Unity3D client needs to interact with
func (p Parser) Parse(in io.Reader) (Spec, error) {

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
		def, err := p.interpretDefinition([]string{"definitions"}, key, val)
		if err != nil {
			return Spec{}, err
		}
		definitions = append(definitions, def)
	}

	return NewSpec(info, definitions), nil
}
