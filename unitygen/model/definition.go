package model

// Definition represents a model of data found inside a swagger file
type Definition interface {
	Name() string
	ToCSharp() string
	ToVariableType() string

	// JsonConverter is the class name responsible for helping interpret a
	// json string and converting it to the specific definition. An empty
	// string means no special class is required.
	JsonConverter() string
}
