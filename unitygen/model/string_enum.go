package model

import (
	"fmt"
	"strings"

	"github.com/recolude/swagger-unity-codegen/unitygen/convention"
)

// StringEnum is a c# enum
type StringEnum struct {
	name   string
	values []string
}

// NewStringEnum creates a new c# enum
func NewStringEnum(name string, values []string) StringEnum {
	return StringEnum{
		name:   name,
		values: values,
	}
}

// ToVariableType generates a identifier for the definition
func (e StringEnum) ToVariableType() string {
	return convention.TitleCase(e.Name())
}

// Name returns the enums name
func (e StringEnum) Name() string {
	return e.name
}

// ToCSharp generates a c# enum for unity
func (e StringEnum) ToCSharp() string {
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
