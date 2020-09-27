package property

import (
	"github.com/recolude/swagger-unity-codegen/unitygen/convention"
	"path/filepath"
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
