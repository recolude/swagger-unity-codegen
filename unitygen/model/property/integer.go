package property

import (
	"strings"

	"github.com/recolude/swagger-unity-codegen/unitygen/convention"
)

type Integer struct {
	PropertyName string
	Format       string
	Description  string
}

func (sp Integer) Name() string {
	return sp.PropertyName
}

func (sp Integer) ToVariableType() string {
	switch sp.Format {
	case "int64":
		return "long"
	default:
		return "int"
	}
}

func (sp Integer) EmptyValue() string {
	switch sp.Format {
	default:
		return "0"
	}
}

func (sp Integer) ClassVariables() string {
	builder := strings.Builder{}

	WriteDescriptionComment(&builder, sp.Description)
	builder.WriteString("	[JsonProperty(\"")
	builder.WriteString(sp.Name())
	builder.WriteString("\")]\n\tpublic ")
	builder.WriteString(sp.ToVariableType())
	builder.WriteString(" ")
	builder.WriteString(convention.TitleCase(sp.Name()))
	builder.WriteString(" { get; set; }\n")
	return builder.String()
}
