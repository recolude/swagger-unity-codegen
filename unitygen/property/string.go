package property

import (
	"fmt"
	"strings"

	"github.com/recolude/swagger-unity-codegen/unitygen/convention"
)

type String struct {
	name   string
	format string
}

func NewString(name string, format string) String {
	return String{
		name:   name,
		format: format,
	}
}

func (sp String) Name() string {
	return sp.name
}

func (sp String) ToVariableType() string {
	switch sp.format {
	case "date-time":
		return "System.DateTime"

	default:
		return "string"
	}
}

func (sp String) EmptyValue() string {
	return "null"
}

func (sp String) ClassVariables() string {

	switch sp.format {
	case "date-time":
		builder := strings.Builder{}
		builder.WriteString("\t[SerializeField]\n")
		fmt.Fprintf(&builder, "\tprivate string %s;\n\n", sp.Name())
		fmt.Fprintf(&builder, "\tpublic System.DateTime %s { get => System.DateTime.Parse(%s); }\n", convention.TitleCase(sp.Name()), sp.Name())
		return builder.String()

	default:
		return fmt.Sprintf("\tpublic %s %s;\n", sp.ToVariableType(), sp.Name())
	}

}
