package convention_test

import (
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen/convention"
	"github.com/stretchr/testify/assert"
)

func TestAddSpce(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		"simple":            {input: "AbcDef", want: "Abc Def"},
		"already lowercase": {input: "abcDef", want: "Abc Def"},
		"snake_case":        {input: "abc_def", want: "Abc Def"},
		"spaces":            {input: " abc def", want: "Abc Def"},
		"spaces2":           {input: " Abc def", want: "Abc Def"},
		"empty string":      {input: "", want: ""},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.want, convention.AddSpace(tc.input))
		})
	}
}
