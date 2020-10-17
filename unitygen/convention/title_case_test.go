package convention_test

import (
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen/convention"
	"github.com/stretchr/testify/assert"
)

func TestTitleCase(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		"simple":              {input: "abcDef", want: "AbcDef"},
		"already capitalized": {input: "AbcDef", want: "AbcDef"},
		"empty string":        {input: "", want: ""},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.want, convention.TitleCase(tc.input))
		})
	}
}
