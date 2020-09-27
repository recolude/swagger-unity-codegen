package property

import "fmt"

type Array struct {
	name string
	prop Property
}

func NewArray(name string, prop Property) Array {
	return Array{
		name: name,
		prop: prop,
	}
}

func (sp Array) Name() string {
	return sp.name
}

func (sp Array) ToVariableType() string {
	return fmt.Sprintf("%s[]", sp.prop.ToVariableType())
}
