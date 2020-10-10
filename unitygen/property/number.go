package property

import "fmt"

type Number struct {
	name   string
	format string
}

func NewNumber(name string, format string) Number {
	return Number{
		name:   name,
		format: format,
	}
}

func (sp Number) Name() string {
	return sp.name
}

func (sp Number) ToVariableType() string {
	if sp.format == "" {
		return "float"
	}

	if sp.format == "int32" {
		return "int"
	}

	return sp.format
}

func (sp Number) EmptyValue() string {
	if sp.format == "" {
		return "0f"
	}

	if sp.format == "int32" {
		return "0"
	}

	return "0f"
}

func (sp Number) ClassVariables() string {
	return fmt.Sprintf("\tpublic %s %s;\n", sp.ToVariableType(), sp.Name())
}
