package path

import "github.com/recolude/swagger-unity-codegen/unitygen/definition"

// Response is one of potentially many responses a single path can recieve
type Response struct {
	description string
	schema      definition.Definition
}

func NewResponse(description string, schema definition.Definition) Response {
	return Response{
		description: description,
		schema:      schema,
	}
}
