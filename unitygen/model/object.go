package model

import (
	"sort"
	"strings"

	"github.com/recolude/swagger-unity-codegen/unitygen/convention"
)

// Object is a collection of properties
type Object struct {
	ObjectName string
	Properties []Property
}

// NewObject creates a new object
func NewObject(name string, properties []Property) Object {
	sort.Sort(sortByPropName(properties))
	return Object{
		ObjectName: name,
		Properties: properties,
	}
}

type sortByPropName []Property

func (a sortByPropName) Len() int           { return len(a) }
func (a sortByPropName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortByPropName) Less(i, j int) bool { return a[i].Name() < a[j].Name() }

// Name is the name of the definition
func (od Object) Name() string {
	return od.ObjectName
}

// ToCSharp generates a c# class for unity
func (od Object) ToCSharp() string {
	var classBuilder strings.Builder

	classBuilder.WriteString("[System.Serializable]\npublic class ")
	classBuilder.WriteString(od.ToVariableType())
	classBuilder.WriteString(" {\n\n")
	for _, prop := range od.Properties {
		// classBuilder.WriteString()
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
