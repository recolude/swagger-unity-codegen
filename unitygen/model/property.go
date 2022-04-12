package model

// Property represents a single variable within a c# class definition
type Property interface {
	// Name of the property (generally a c# variable name)
	Name() string

	// ToVariableType returns the name of the variable type that exists in c# (ie: float, int, string)
	ToVariableType() string

	// EmptyValue is the value that represents the property has yet to be set.
	EmptyValue() string

	// What gets written to the c# class definition.
	ClassVariables() string
}
