package security

import (
	"fmt"

	"github.com/recolude/swagger-unity-codegen/unitygen/convention"
)

// APIKeyLocation represents a location inside a request in where an API Key can be found
type APIKeyLocation string

const (
	// Header indicates the API Key should be found inside a header in the request
	Header APIKeyLocation = "header"

	// Body indicates the API Key should be found inside the body of the request
	Body = "body"
)

// APIKeyAuth is a guard on route that requires a specific API key to be
// present somewhere in the request
type APIKeyAuth struct {
	// How the swagger file refers to they key
	identifier string

	// what you will actually find in something like an HTTP header
	key string

	// Where the key should be located in the HTTP message (header, body, etc.)
	loc APIKeyLocation
}

// NewAPIKey creates a new API key
func NewAPIKey(identifier string, key string, location APIKeyLocation) APIKeyAuth {
	return APIKeyAuth{
		identifier: identifier,
		key:        key,
		loc:        location,
	}
}

// Identifier returns a unique string that represents how the swagger file
// refers to the API key
func (key APIKeyAuth) Identifier() string {
	return key.identifier
}

func (key APIKeyAuth) name() string {
	return key.key
}

// ModifyNetworkRequest generates C# code that appends this API Key to a
// specific network request
func (key APIKeyAuth) ModifyNetworkRequest() string {
	return fmt.Sprintf("unityNetworkReq.SetRequestHeader(\"%s\", this.Config.%s);", key.name(), convention.TitleCase(key.Identifier()))
}

func (key APIKeyAuth) String() string {
	return fmt.Sprintf(
		"%s is a API Key '%s' found in a request's %s",
		convention.TitleCase(key.Identifier()),
		key.key,
		key.loc,
	)
}
