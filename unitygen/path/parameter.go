package path

import (
	"github.com/recolude/swagger-unity-codegen/unitygen/model"
)

type ParameterLocation string

const (
	PathParameterLocation  ParameterLocation = "path"
	QueryParameterLocation ParameterLocation = "query"
	BodyParameterLocation  ParameterLocation = "body"
)

// Parameter represents a variable that should exist somewhere inside a HTTP
// request
type Parameter struct {
	location      ParameterLocation
	name          string
	required      bool
	parameterType model.Property // ex: string, query
}

func NewParameter(location ParameterLocation, name string, required bool, parameterType model.Property) Parameter {
	return Parameter{
		location:      location,
		name:          name,
		required:      required,
		parameterType: parameterType,
	}
}

func (param Parameter) Name() string {
	return param.name
}

func (param Parameter) Location() ParameterLocation {
	return param.location
}

func (param Parameter) Required() bool {
	return param.required
}

func (param Parameter) Schema() model.Property {
	return param.parameterType
}
