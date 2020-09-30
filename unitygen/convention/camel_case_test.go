package convention_test

import (
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen/convention"
	"github.com/stretchr/testify/assert"
)

func TestCamelCase(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		"simple":            {input: "AbcDef", want: "abcDef"},
		"already lowercase": {input: "abcDef", want: "abcDef"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.want, convention.CamelCase(tc.input))
		})
	}
}
