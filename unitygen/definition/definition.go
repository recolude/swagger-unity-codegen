package definition

// Definition represents a model of data found inside a swagger file
type Definition interface {
	Name() string
	ToClass() string
	ToVariableType() string
}
