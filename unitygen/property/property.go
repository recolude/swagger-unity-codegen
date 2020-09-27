package property

type Property interface {
	Name() string
	ToVariableType() string
}
