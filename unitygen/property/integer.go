package property

type Integer struct {
	name   string
	format string
}

func NewInteger(name string, format string) Integer {
	return Integer{
		name:   name,
		format: format,
	}
}

func (sp Integer) Name() string {
	return sp.name
}

func (sp Integer) ToVariableType() string {
	switch sp.format {

	default:
		return "int"
	}
}

func (sp Integer) EmptyValue() string {
	switch sp.format {

	default:
		return "0"
	}
}
