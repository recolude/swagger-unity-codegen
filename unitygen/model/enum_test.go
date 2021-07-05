package model_test

import (
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen/model"
	"github.com/stretchr/testify/assert"
)

func TestEnum(t *testing.T) {
	// ******************************** ARRANGE *******************************
	enum := model.NewEnum(
		"testEnum",
		[]string{
			"A",
			"b",
			"CDF",
		},
	)

	// ********************************** ACT *********************************
	varType := enum.ToVariableType()
	name := enum.Name()
	cSharp := enum.ToCSharp()

	// ********************************* ASSERT *******************************
	assert.Equal(t, "TestEnum", varType)
	assert.Equal(t, "testEnum", name)
	assert.Equal(t, `public enum TestEnum {
	A = 0,
	b = 1,
	CDF = 2
}`, cSharp)
}
