package unity_test

import (
	"net/http"
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen/unity"
	"github.com/stretchr/testify/assert"
)

func TestToUnityHTTPVerb(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		"get":    {input: http.MethodGet, want: "UnityWebRequest.kHttpVerbGET"},
		"put":    {input: http.MethodPut, want: "UnityWebRequest.kHttpVerbPUT"},
		"post":   {input: http.MethodPost, want: "UnityWebRequest.kHttpVerbPOST"},
		"delete": {input: http.MethodDelete, want: "UnityWebRequest.kHttpVerbDELETE"},
		"head":   {input: http.MethodHead, want: "UnityWebRequest.kHttpVerbHEAD"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := unity.ToUnityHTTPVerb(tc.input)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestPanicAtUnkown(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	unity.ToUnityHTTPVerb("unknown")
}
