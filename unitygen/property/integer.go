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
	case "int32":
		return "int32"

	default:
		return "int"
	}
}
