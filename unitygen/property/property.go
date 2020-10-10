package property

type Property interface {
	Name() string
	ToVariableType() string

	// EmptyValue is the value that represents the property has yet to be set.
	EmptyValue() string

	ClassVariables() string
}
