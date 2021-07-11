package path

// Response is one of potentially many responses a single path can receive
type Response interface {
	Description() string
	Interpret(variableName string, downloadHandlerName string) string
	VariableType() string
}
