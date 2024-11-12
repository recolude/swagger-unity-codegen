package property

import (
	"fmt"
	"strings"

	"github.com/recolude/swagger-unity-codegen/unitygen/convention"
	"github.com/recolude/swagger-unity-codegen/unitygen/model"
)

type Array struct {
	PropertyName string
	Item         model.Property
	Description  string
}

func (sp Array) Name() string {
	return sp.PropertyName
}

func (sp Array) Property() model.Property {
	return sp.Item
}

func (sp Array) ToVariableType() string {
	return fmt.Sprintf("%s[]", sp.Item.ToVariableType())
}

func (sp Array) EmptyValue() string {
	return "null"
}

func (sp Array) ClassVariables() string {
	builder := strings.Builder{}
	WriteDescriptionComment(&builder, sp.Description)
	builder.WriteString("\t[JsonProperty(\"")
	builder.WriteString(sp.PropertyName)
	builder.WriteString("\")]\n\tpublic ")
	builder.WriteString(sp.ToVariableType())
	builder.WriteString(" ")
	builder.WriteString(convention.TitleCase(sp.PropertyName))
	builder.WriteString(" { get; set; }\n")
	return builder.String()
}
