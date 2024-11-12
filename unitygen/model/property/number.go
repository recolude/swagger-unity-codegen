package property

import (
	"fmt"
	"strings"

	"github.com/recolude/swagger-unity-codegen/unitygen/convention"
)

type Number struct {
	PropertyName string
	Format       string
	Description  string
}

func (sp Number) Name() string {
	return sp.PropertyName
}

func (sp Number) ToVariableType() string {

	switch sp.Format {
	case "int32":
		return "int"

	case "double":
		return "double"

	default:
		return "float"
	}

}

func (sp Number) EmptyValue() string {
	if sp.Format == "" {
		return "0f"
	}

	if sp.Format == "int32" {
		return "0"
	}

	return "0f"
}

func (sp Number) ClassVariables() string {
	builder := strings.Builder{}
	WriteDescriptionComment(&builder, sp.Description)
	fmt.Fprintf(&builder, "\t[JsonProperty(\"%s\")]\n", sp.Name())
	fmt.Fprintf(&builder, "\tpublic %s %s { get; set; }\n", sp.ToVariableType(), convention.TitleCase(sp.Name()))
	return builder.String()
}
