package property

type String struct {
	name   string
	format string
}

func NewString(name string, format string) String {
	return String{
		name:   name,
		format: format,
	}
}

func (sp String) Name() string {
	return sp.name
}

func (sp String) ToVariableType() string {
	switch sp.format {
	case "date-time":
		return "System.DateTime"

	default:
		return "string"
	}
}
