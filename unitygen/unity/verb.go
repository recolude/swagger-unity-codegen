package unity

import (
	"fmt"
	"net/http"
)

// ToUnityHTTPVerb takes how golang represents a HTTP verb (http.MethodGet) and
// translates it for unity (UnityWebRequest.kHttpVerbGET).
func ToUnityHTTPVerb(httpMethod string) string {
	switch httpMethod {
	case http.MethodGet:
		return "UnityWebRequest.kHttpVerbGET"
	case http.MethodPut:
		return "UnityWebRequest.kHttpVerbPUT"
	case http.MethodPost:
		return "UnityWebRequest.kHttpVerbPOST"
	case http.MethodDelete:
		return "UnityWebRequest.kHttpVerbDELETE"
	case http.MethodHead:
		return "UnityWebRequest.kHttpVerbHEAD"
	}
	panic(fmt.Sprintf("unknown verb \"%s\"", httpMethod))
}
