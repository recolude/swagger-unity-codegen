package model

import (
	"path/filepath"

	"github.com/recolude/swagger-unity-codegen/unitygen/convention"
)

// ObjectReference points to another object for it's definition
type ObjectReference struct {
	ref string
}

// NewObjectReference creates a new object reference
func NewObjectReference(ref string) ObjectReference {
	return ObjectReference{
		ref: ref,
	}
}

// Name is the name of the definition
func (or ObjectReference) Name() string {
	return or.ref
}

func (or ObjectReference) ToCSharp() string {
	// What a smell, am I rite?
	// Maybe all this should do is create a class that inherits the reference?
	// Doens't make sense where request body's schema is direct reference to this tho...
	// But I guess it wouldn't actuall be *that* bad.
	panic("unimplemented")
}

// ToVariableType generates a identifier for the definition
func (or ObjectReference) ToVariableType() string {
	return convention.TitleCase(filepath.Base(or.ref))
}
