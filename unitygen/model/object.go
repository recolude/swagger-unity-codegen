package model

import (
	"fmt"
	"sort"
	"strings"

	"github.com/recolude/swagger-unity-codegen/unitygen/convention"
)

// Object is a collection of properties
type Object struct {
	ObjectName             string
	objectToTakeProperties *Object
	properties             []Property
	children               []*Object
	discriminator          string
	inherits               *Object
}

// NewObject creates a new object
func NewObject(name string, properties []Property) Object {
	sort.Sort(sortByPropName(properties))
	return Object{
		ObjectName: name,
		properties: properties,
	}
}

// NewAllOfObject is an object that has all properties of another object. It's not
// proper inheritance, just composition. For inheritance, look at discriminator
// object.
func NewAllOfObject(name string, objectToTakeProperties Object, extraProperties []Property) Object {
	sort.Sort(sortByPropName(extraProperties))
	return Object{
		ObjectName:             name,
		properties:             extraProperties,
		objectToTakeProperties: &objectToTakeProperties,
	}
}

func NewDiscriminatorObject(name string, properties []Property, discriminator string) Object {
	sort.Sort(sortByPropName(properties))
	return Object{
		ObjectName:    name,
		properties:    properties,
		discriminator: discriminator,
	}
}

type sortByPropName []Property

func (a sortByPropName) Len() int           { return len(a) }
func (a sortByPropName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortByPropName) Less(i, j int) bool { return a[i].Name() < a[j].Name() }

func (od Object) HasDiscriminator() bool {
	return od.discriminator != ""
}

func (od *Object) SetWhatToInherit(obj *Object) {
	od.inherits = obj
}

func (od *Object) AddChild(child *Object) {
	od.children = append(od.children, child)
}

// Name is the name of the definition
func (od Object) Name() string {
	return od.ObjectName
}

func (od Object) Properties() []Property {
	if od.objectToTakeProperties != nil {
		return append(od.objectToTakeProperties.Properties(), od.properties...)
	}
	return od.properties
}

// SetAllOfObject makes this object take on the properties of the object passed
// in, satisfying the "allOf" modifier in swagger 2.0
func (od *Object) SetAllOfObject(objectToTakeProperties *Object) {
	od.objectToTakeProperties = objectToTakeProperties
}

// ToCSharp generates a c# class for unity
func (od Object) ToCSharp() string {
	var classBuilder strings.Builder

	classBuilder.WriteString("[System.Serializable]\n")

	if od.HasDiscriminator() {
		classBuilder.WriteString("[JsonConverter(typeof(JsonSubtypes), \"")
		classBuilder.WriteString(od.discriminator)
		classBuilder.WriteString("\")]\n")

		for _, child := range od.children {
			if child == nil {
				panic(fmt.Errorf("%s was elected as parent class and provided a nil child", od.ToVariableType()))
			}
			classBuilder.WriteString("[JsonSubtypes.KnownSubType(typeof(")
			classBuilder.WriteString(child.ToVariableType())
			classBuilder.WriteString("), \"")
			classBuilder.WriteString(child.ToVariableType())
			classBuilder.WriteString("\")]\n")
		}
	}

	classBuilder.WriteString("public class ")
	classBuilder.WriteString(od.ToVariableType())

	if od.inherits != nil {
		classBuilder.WriteString(" : ")
		classBuilder.WriteString(od.inherits.ToVariableType())
	}

	classBuilder.WriteString(" {\n\n")

	// List out any "any of" properties
	if od.objectToTakeProperties != nil {
		for _, prop := range od.objectToTakeProperties.Properties() {
			classBuilder.WriteString(prop.ClassVariables())
			classBuilder.WriteString("\n")
		}
	}

	for _, prop := range od.properties {
		classBuilder.WriteString(prop.ClassVariables())
		classBuilder.WriteString("\n")
	}
	classBuilder.WriteString("}")
	return classBuilder.String()
}

// ToVariableType generates a identifier for the definition
func (od Object) ToVariableType() string {
	return convention.TitleCase(od.Name())
}

func (od Object) JsonConverter() string {
	return ""
}
