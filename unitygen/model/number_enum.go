package model

import (
	"fmt"
	"math"
	"strings"

	"github.com/recolude/swagger-unity-codegen/unitygen/convention"
)

// NumberEnum is a c# enum that translates to some number when used in a web
// request
type NumberEnum struct {
	name   string
	values []float64
}

// NewStringEnum creates a new c# enum
func NewNumberEnum(name string, values []float64) NumberEnum {
	return NumberEnum{
		name:   name,
		values: values,
	}
}

// ToVariableType generates a identifier for the definition
func (e NumberEnum) ToVariableType() string {
	return convention.TitleCase(e.Name())
}

// Name returns the enums name
func (e NumberEnum) Name() string {
	return e.name
}

func floatToEnumMember(x float64) string {
	i := x

	sb := strings.Builder{}

	sb.WriteString("NUMBER_")

	if i < 0 {
		i = math.Abs(i)
		sb.WriteString("NEG_")
	}

	wholeValue := math.Floor(i)
	sb.WriteString(fmt.Sprintf("%d", int(wholeValue)))

	remaining := i - wholeValue

	if remaining <= 0 {
		return sb.String()
	}

	sb.WriteString("_DOT_")

	for remaining > 0.000001 {
		remaining *= 10
		wholeValue := math.Floor(remaining)
		sb.WriteString(fmt.Sprintf("%d", int(wholeValue)))
		remaining -= wholeValue
	}

	return sb.String()
}

// ToCSharp generates a c# enum for unity
func (e NumberEnum) ToCSharp() string {
	var enumBuilder strings.Builder

	enumBuilder.WriteString("public enum ")
	enumBuilder.WriteString(e.ToVariableType())
	enumBuilder.WriteString(" {\n")
	for i, prop := range e.values {
		enumBuilder.WriteString(fmt.Sprintf("\t%s", floatToEnumMember(prop)))
		if i < len(e.values)-1 {
			enumBuilder.WriteString(",\n")
		}
	}
	enumBuilder.WriteString("\n}")
	return enumBuilder.String()
}
