package security

type Auth interface {
	ToCSharp() string
}
