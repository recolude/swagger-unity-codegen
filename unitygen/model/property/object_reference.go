package property

import (
	"path/filepath"
	"strings"

	"github.com/recolude/swagger-unity-codegen/unitygen/convention"
)

type ObjectReference struct {
	name          string
	referencePath string
}

func NewObjectReference(name string, referencePath string) ObjectReference {
	return ObjectReference{
		name:          name,
		referencePath: referencePath,
	}
}

func (orp ObjectReference) Name() string {
	return orp.name
}

func (orp ObjectReference) ToVariableType() string {
	return convention.TitleCase(filepath.Base(orp.referencePath))
}

func (orp ObjectReference) EmptyValue() string {
	return "null"
}

func (orp ObjectReference) ClassVariables() string {
	builder := strings.Builder{}
	builder.WriteString("\t[JsonProperty(\"")
	builder.WriteString(orp.name)
	builder.WriteString("\")]\n\tpublic ")
	builder.WriteString(orp.ToVariableType())
	builder.WriteString(" ")
	builder.WriteString(convention.TitleCase(orp.name))
	builder.WriteString(" { get; private set; }\n")
	return builder.String()
}
