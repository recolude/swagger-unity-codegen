package definition

import (
	"fmt"
	"strings"

	"github.com/recolude/swagger-unity-codegen/unitygen/convention"
)

// Enum is a c# enum
type Enum struct {
	name   string
	values []string
}

// NewEnum creates a new c# enum
func NewEnum(name string, values []string) Enum {
	return Enum{
		name:   name,
		values: values,
	}
}

// ToVariableType generates a identifier for the definition
func (e Enum) ToVariableType() string {
	return convention.TitleCase(e.Name())
}

// Name returns the enums name
func (e Enum) Name() string {
	return e.name
}

// ToCSharp generates a c# enum for unity
func (e Enum) ToCSharp() string {
	var enumBuilder strings.Builder

	enumBuilder.WriteString("public enum ")
	enumBuilder.WriteString(e.ToVariableType())
	enumBuilder.WriteString(" {\n")
	for i, prop := range e.values {
		enumBuilder.WriteString(fmt.Sprintf("\t%s = %d", prop, i))
		if i < len(e.values)-1 {
			enumBuilder.WriteString(",\n")
		}
	}
	enumBuilder.WriteString("\n}")
	return enumBuilder.String()
}
