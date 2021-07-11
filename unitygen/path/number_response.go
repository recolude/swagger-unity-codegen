package path

import "fmt"

// NumberResponse is a type of response that expects a floating point number.
type NumberResponse struct {
	description string
}

func NewNumberResponse(description string) NumberResponse {
	return NumberResponse{
		description: description,
	}
}

func (resp NumberResponse) Description() string {
	return resp.description
}

func (resp NumberResponse) Interpret(variableName string, downloadHandlerVariableName string) string {
	return fmt.Sprintf("%s = float.Parse(%s.text);", variableName, downloadHandlerVariableName)
}

func (resp NumberResponse) VariableType() string {
	return "float"
}
