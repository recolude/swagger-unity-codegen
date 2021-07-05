package property

import (
	"fmt"

	"github.com/recolude/swagger-unity-codegen/unitygen/model"
)

type Array struct {
	name string
	prop model.Property
}

func NewArray(name string, prop model.Property) Array {
	return Array{
		name: name,
		prop: prop,
	}
}

func (sp Array) Name() string {
	return sp.name
}

func (sp Array) Property() model.Property {
	return sp.prop
}

func (sp Array) ToVariableType() string {
	return fmt.Sprintf("%s[]", sp.prop.ToVariableType())
}

func (sp Array) EmptyValue() string {
	return "null"
}

func (sp Array) ClassVariables() string {
	return fmt.Sprintf("\tpublic %s %s;\n", sp.ToVariableType(), sp.Name())
}
