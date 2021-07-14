package property

import (
	"fmt"
	"strings"

	"github.com/recolude/swagger-unity-codegen/unitygen/convention"
)

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
	builder := strings.Builder{}
	fmt.Fprintf(&builder, "\t[JsonProperty(\"%s\")]\n", sp.Name())
	fmt.Fprintf(&builder, "\tpublic %s %s { get; private set; }\n", sp.ToVariableType(), convention.TitleCase(sp.Name()))
	return builder.String()
}
