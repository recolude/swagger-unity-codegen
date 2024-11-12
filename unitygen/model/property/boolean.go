package property

import (
	"fmt"
	"strings"

	"github.com/recolude/swagger-unity-codegen/unitygen/convention"
)

type Boolean struct {
	PropertyName string
	Description  string
}

func (sp Boolean) Name() string {
	return sp.PropertyName
}

func (sp Boolean) ToVariableType() string {
	return "bool"
}

func (sp Boolean) EmptyValue() string {
	return "false"
}

func (sp Boolean) ClassVariables() string {
	builder := strings.Builder{}

	WriteDescriptionComment(&builder, sp.Description)
	fmt.Fprintf(&builder, "\t[JsonProperty(\"%s\")]\n\tpublic %s %s { get; set; }\n", sp.Name(), sp.ToVariableType(), convention.TitleCase(sp.Name()))
	return builder.String()
}
