package convention_test

import (
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen/convention"
	"github.com/stretchr/testify/assert"
)

func TestClassName(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		"simple":              {input: "abcDef", want: "AbcDef"},
		"snake_case":          {input: "abc_def", want: "Abc_Def"},
		"snake_case weird":    {input: "abc__def", want: "Abc__Def"},
		"already capitalized": {input: "AbcDef", want: "AbcDef"},
		"kebab weird":         {input: "-_abc-def", want: "AbcDef"},
		"empty string":        {input: "", want: ""},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.want, convention.ClassName(tc.input))
		})
	}
}
