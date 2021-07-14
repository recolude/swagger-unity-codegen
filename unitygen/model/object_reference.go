package model

import (
	"errors"
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
	//
	// But honestly the only time this is used is as a refernce for a response.
	// I've at this point refactored response as an internface. Maybe at this point
	// I just make a ObjectReferenceResponse type?
	panic(errors.New("unimplemented"))
}

// ToVariableType generates a identifier for the definition
func (or ObjectReference) ToVariableType() string {
	return convention.TitleCase(filepath.Base(or.ref))
}
