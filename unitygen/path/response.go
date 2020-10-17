package path

import "github.com/recolude/swagger-unity-codegen/unitygen/definition"

// Response is one of potentially many responses a single path can receive
type Response struct {
	description string
	schema      definition.Definition
}

// NewResponse creates a new response
func NewResponse(description string, schema definition.Definition) Response {
	return Response{
		description: description,
		schema:      schema,
	}
}

// Schema returns how the response should be structured
func (resp Response) Schema() definition.Definition {
	return resp.schema
}
