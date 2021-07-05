package property

import (
	"fmt"

	"github.com/recolude/swagger-unity-codegen/unitygen/convention"
	"github.com/recolude/swagger-unity-codegen/unitygen/model"
)

type Object struct {
	name string
	obj  model.Object
}

func NewObject(name string, obj model.Object) Object {
	return Object{
		name: name,
		obj:  obj,
	}
}

// Name of the property (generally a c# variable name)
func (op Object) Name() string {
	return op.name
}

// ToVariableType returns the name of the variable type that exists in c# (ie: float, int,s tring)
func (op Object) ToVariableType() string {
	return op.obj.ToVariableType()
}

// EmptyValue is the value that represents the property has yet to be set.
func (op Object) EmptyValue() string {
	return "null"
}

// What gets written to the c# class definition.
func (op Object) ClassVariables() string {
	return fmt.Sprintf("\t%s\n\tpublic %s %s;", op.obj.ToCSharp(), op.ToVariableType(), convention.CamelCase(op.name))
}
