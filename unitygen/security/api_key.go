package security

// APIKeyAuth is a guard on route that requires a specific API key to be
// present somewhere in the request
type APIKeyAuth struct {
	Name string
	In   string
}
