package property

import (
	"fmt"
	"strings"

	"github.com/recolude/swagger-unity-codegen/unitygen/convention"
)

type String struct {
	PropertyName string
	Format       string
	Description  string
}

func (sp String) Name() string {
	return sp.PropertyName
}

func (sp String) ToVariableType() string {
	switch sp.Format {
	case "date-time":
		return "System.DateTime"

	case "binary":
		return "byte[]"

	default:
		return "string"
	}
}

func (sp String) EmptyValue() string {
	return "null"
}

func (sp String) ClassVariables() string {
	builder := strings.Builder{}

	WriteDescriptionComment(&builder, sp.Description)

	fmt.Fprintf(&builder, "\t[JsonProperty(\"%s\")]\n", sp.Name())

	switch sp.Format {
	case "date-time":
		fmt.Fprintf(&builder, "\tpublic string %s;\n\n", convention.CamelCase(sp.Name()))
		fmt.Fprintf(&builder, "\tpublic System.DateTime %s { get => System.DateTime.Parse(%s); }\n", convention.TitleCase(sp.Name()), convention.CamelCase(sp.Name()))

	case "binary":
		fmt.Fprintf(&builder, "\tpublic byte[] %s { get; set; }\n", convention.TitleCase(sp.Name()))

	default:
		fmt.Fprintf(&builder, "\tpublic string %s { get; set; }\n", convention.TitleCase(sp.Name()))
	}

	return builder.String()
}
