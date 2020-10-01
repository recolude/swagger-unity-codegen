package security

// Auth is a specific guard desired on one or more routes in an API
type Auth interface {
	// Identifier is the identifier used throughout a swagger file to refer to a
	// specific type of authentication desired for a specific route
	Identifier() string

	ModifyNetworkRequest() string

	String() string
}
