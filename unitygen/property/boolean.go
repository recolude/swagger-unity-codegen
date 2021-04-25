package property

import "fmt"

type Boolean struct {
	name string
}

func NewBoolean(name string) Boolean {
	return Boolean{
		name: name,
	}
}

func (sp Boolean) Name() string {
	return sp.name
}

func (sp Boolean) ToVariableType() string {
	return "bool"
}

func (sp Boolean) EmptyValue() string {
	return "false"
}

func (sp Boolean) ClassVariables() string {
	return fmt.Sprintf("\tpublic %s %s;\n", sp.ToVariableType(), sp.Name())
}
