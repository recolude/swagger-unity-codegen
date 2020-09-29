package path

// SecurityMethodReference is a way to ensure the client properly communicates with
// a specific route
type SecurityMethodReference struct {
	Identifier string
	Contents   []string
}

func NewSecurityMethodReference(name string) SecurityMethodReference {
	return SecurityMethodReference{
		Identifier: name,
	}
}
