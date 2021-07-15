package property

import (
	"strings"

	"github.com/recolude/swagger-unity-codegen/unitygen/convention"
	"github.com/recolude/swagger-unity-codegen/unitygen/model"
)

type DefinitionReference struct {
	name       string
	definition model.Definition
}

func NewDefinitionReference(name string, definition model.Definition) DefinitionReference {
	return DefinitionReference{
		name:       name,
		definition: definition,
	}
}

func (dr DefinitionReference) Name() string {
	return dr.name
}

func (dr DefinitionReference) ToVariableType() string {
	return dr.definition.ToVariableType()
}

func (dr DefinitionReference) EmptyValue() string {
	return "null"
}

func (dr DefinitionReference) ClassVariables() string {
	builder := strings.Builder{}
	builder.WriteString("\t[JsonProperty(\"")
	builder.WriteString(dr.name)
	builder.WriteString("\")]\n")

	converter := dr.definition.JsonConverter()
	if converter != "" {
		builder.WriteString("\t[JsonConverter(typeof(")
		builder.WriteString(converter)
		builder.WriteString("))]\n")
	}

	builder.WriteString("\tpublic ")
	builder.WriteString(dr.ToVariableType())
	builder.WriteString(" ")
	builder.WriteString(convention.TitleCase(dr.name))
	builder.WriteString(" { get; private set; }\n")
	return builder.String()
}
