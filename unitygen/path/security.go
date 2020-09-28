package path

// SecurityMethodReference is a way to ensure the client properly communicates with
// a specific route
type SecurityMethodReference struct {
	Name     string
	Contents []string
}

func NewSecurityMethodReference(name string) SecurityMethodReference {
	return SecurityMethodReference{
		Name: name,
	}
}
