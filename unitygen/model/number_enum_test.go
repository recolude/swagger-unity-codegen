package model_test

import (
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen/model"
	"github.com/stretchr/testify/assert"
)

func TestNumberEnum(t *testing.T) {
	// ******************************** ARRANGE *******************************
	enum := model.NewNumberEnum(
		"testEnum",
		[]float64{
			0.125,
			0.25,
			0.5,
			0,
			1,
			1.5,
			2,
			4,
			8,
			122.1109,
			-0.125,
			-0.25,
			-0.5,
			-0,
			-1,
			-1.5,
			-2,
			-4,
			-8,
			-122.1109,
			1.09,
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
	NUMBER_0_DOT_125,
	NUMBER_0_DOT_25,
	NUMBER_0_DOT_5,
	NUMBER_0,
	NUMBER_1,
	NUMBER_1_DOT_5,
	NUMBER_2,
	NUMBER_4,
	NUMBER_8,
	NUMBER_122_DOT_1109,
	NUMBER_NEG_0_DOT_125,
	NUMBER_NEG_0_DOT_25,
	NUMBER_NEG_0_DOT_5,
	NUMBER_0,
	NUMBER_NEG_1,
	NUMBER_NEG_1_DOT_5,
	NUMBER_NEG_2,
	NUMBER_NEG_4,
	NUMBER_NEG_8,
	NUMBER_NEG_122_DOT_1109,
	NUMBER_1_DOT_09
}`, cSharp)
}
