package path

import (
	"fmt"

	"github.com/recolude/swagger-unity-codegen/unitygen/model/property"
)

// ArrayResponse is an array of type T as a response.
type ArrayResponse struct {
	array       property.Array
	description string
}

// NewArrayResponse creates a new array response
func NewArrayResponse(description string, array property.Array) ArrayResponse {
	return ArrayResponse{
		array:       array,
		description: description,
	}
}

// Description of the response
func (resp ArrayResponse) Description() string {
	return resp.description
}

func (resp ArrayResponse) Interpret(variableName string, downloadHandlerVariableName string) string {
	return fmt.Sprintf("%s = JsonConvert.DeserializeObject<%s>(%s.text);", variableName, resp.VariableType(), downloadHandlerVariableName)
}

func (resp ArrayResponse) VariableType() string {
	return resp.array.ToVariableType()
}
