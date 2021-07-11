package path

import "fmt"

// FileResponse is a type of response that expects a binary file.
type FileResponse struct {
	description string
}

func NewFileResponse(description string) FileResponse {
	return FileResponse{
		description: description,
	}
}

func (resp FileResponse) Description() string {
	return resp.description
}

func (resp FileResponse) Interpret(variableName string, downloadHandlerVariableName string) string {
	return fmt.Sprintf("%s = %s.data;", variableName, downloadHandlerVariableName)
}

func (resp FileResponse) VariableType() string {
	return "byte[]"
}
