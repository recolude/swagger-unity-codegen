package definition

import (
	"fmt"
	"sort"
	"strings"

	"github.com/recolude/swagger-unity-codegen/unitygen/convention"
	"github.com/recolude/swagger-unity-codegen/unitygen/property"
)

// ObjectDefinition is a collection of properties
type Object struct {
	ObjectName string
	Properties []property.Property
}

func NewObject(name string, properties []property.Property) Object {
	sort.Sort(sortByPropName(properties))
	return Object{
		ObjectName: name,
		Properties: properties,
	}
}

type sortByPropName []property.Property

func (a sortByPropName) Len() int           { return len(a) }
func (a sortByPropName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortByPropName) Less(i, j int) bool { return a[i].Name() < a[j].Name() }

// Name is the name of the definition
func (od Object) Name() string {
	return od.ObjectName
}

// ToClass generates a c# class for unity
func (od Object) ToClass() string {
	var classBuilder strings.Builder

	classBuilder.WriteString("[System.Serializable]\npublic class ")
	classBuilder.WriteString(od.ToVariableType())
	classBuilder.WriteString(" {\n\n")
	for _, prop := range od.Properties {
		classBuilder.WriteString(fmt.Sprintf("\tpublic %s %s;\n", prop.ToVariableType(), prop.Name()))
		classBuilder.WriteString("\n")
	}
	classBuilder.WriteString("}")
	return classBuilder.String()
}

// ToVariableType generates a identifier for the definition
func (od Object) ToVariableType() string {
	return convention.TitleCase(od.Name())
}
