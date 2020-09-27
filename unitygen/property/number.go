package property

type Number struct {
	name   string
	format string
}

func NewNumber(name string, format string) Number {
	return Number{
		name:   name,
		format: format,
	}
}

func (sp Number) Name() string {
	return sp.name
}

func (sp Number) ToVariableType() string {
	if sp.format == "" {
		return "float"
	}
	return sp.format
}
