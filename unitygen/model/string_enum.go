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
	varType := e.ToVariableType()

	// Write Actuall Enum
	enumBuilder.WriteString("public enum ")
	enumBuilder.WriteString(varType)
	enumBuilder.WriteString(" {\n")
	for i, prop := range e.values {
		enumBuilder.WriteString(fmt.Sprintf("\t%s = %d", convention.ClassName(prop), i))
		if i < len(e.values)-1 {
			enumBuilder.WriteString(",\n")
		}
	}
	enumBuilder.WriteString("\n}\n")

	// Write JSON converter for enum
	fmt.Fprintf(&enumBuilder, "public class %sJsonConverter : JsonConverter {\n", varType)

	// WriteJSON function
	enumBuilder.WriteString("\tpublic override void WriteJson(JsonWriter w, object val, JsonSerializer s) {\n")
	fmt.Fprintf(&enumBuilder, "\t\t%s castedVal = (%s)val;\n", varType, varType)
	enumBuilder.WriteString("\t\tswitch (castedVal) {\n")
	for _, prop := range e.values {
		fmt.Fprintf(&enumBuilder, "\t\t\tcase %s.%s:\n", varType, convention.ClassName(prop))
		fmt.Fprintf(&enumBuilder, "\t\t\t\tw.WriteValue(\"%s\");\n", prop)
		enumBuilder.WriteString("\t\t\t\tbreak;\n")
	}
	enumBuilder.WriteString("\t\t\tdefault:\n")
	enumBuilder.WriteString("\t\t\t\tthrow new System.Exception(\"Unknown value. Living on the dangerous side editing generated code?\");\n")
	enumBuilder.WriteString("\t\t}\n\t}\n\n")

	// ReadJSON function
	enumBuilder.WriteString("\tpublic override object ReadJson(JsonReader r, System.Type t, object existingValue, JsonSerializer s) {\n")
	enumBuilder.WriteString("\t\tvar enumString = (string)r.Value;\n")
	enumBuilder.WriteString("\t\tswitch (enumString) {\n")
	for _, prop := range e.values {
		fmt.Fprintf(&enumBuilder, "\t\t\tcase \"%s\":\n", prop)
		fmt.Fprintf(&enumBuilder, "\t\t\t\treturn %s.%s;\n", varType, convention.ClassName(prop))
	}
	enumBuilder.WriteString("\t\t\tdefault:\n")
	enumBuilder.WriteString("\t\t\t\tthrow new System.Exception(\"Unknown value. Perhaps you need to regenerate this code?\");\n")
	enumBuilder.WriteString("\t\t}\n\t}\n\n")

	// Can Convert function
	enumBuilder.WriteString("\tpublic override bool CanConvert(System.Type objectType) {\n\t\treturn objectType == typeof(string);\n\t}")

	// Close out converter class
	enumBuilder.WriteString("\n}")

	return enumBuilder.String()
}
