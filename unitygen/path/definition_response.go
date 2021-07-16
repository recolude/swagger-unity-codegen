package path

import (
	"errors"
	"fmt"

	"github.com/recolude/swagger-unity-codegen/unitygen/model"
)

// DefinitionResponse is a type of response that expects an object definition.
type DefinitionResponse struct {
	description string
	schema      model.Definition
}

// NewDefinitionResponse creates a new definition response
func NewDefinitionResponse(description string, schema model.Definition) DefinitionResponse {
	return DefinitionResponse{
		description: description,
		schema:      schema,
	}
}

func (resp DefinitionResponse) Description() string {
	return resp.description
}

func (resp DefinitionResponse) Interpret(variableName string, downloadHandlerVariableName string) string {
	if resp.schema == nil {
		panic(errors.New("can not build a response interpretation from nil definition"))
	}
	return fmt.Sprintf("%s = JsonConvert.DeserializeObject<%s>(%s.text);", variableName, resp.schema.ToVariableType(), downloadHandlerVariableName)
}

func (resp DefinitionResponse) VariableType() string {
	if resp.schema == nil {
		panic(errors.New("can not build a response variable type from nil definition"))
	}
	return resp.schema.ToVariableType()
}
