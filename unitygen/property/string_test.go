package property_test

import (
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen/property"
	"github.com/stretchr/testify/assert"
)

func Test_StringDefaultsToFloat(t *testing.T) {
	// ******************************** ARRANGE *******************************
	ref := property.NewString("someName", "")

	// ********************************** ACT *********************************
	name := ref.Name()
	varType := ref.ToVariableType()
	nullVal := ref.EmptyValue()
	classVars := ref.ClassVariables()

	// ********************************* ASSERT *******************************
	assert.Equal(t, "someName", name)
	assert.Equal(t, "string", varType)
	assert.Equal(t, "null", nullVal)
	assert.Equal(t, "	public string someName;\n", classVars)
}

func Test_StringInterpretsDate(t *testing.T) {
	// ******************************** ARRANGE *******************************
	ref := property.NewString("someName", "date-time")

	// ********************************** ACT *********************************
	name := ref.Name()
	varType := ref.ToVariableType()
	nullVal := ref.EmptyValue()
	classVars := ref.ClassVariables()

	// ********************************* ASSERT *******************************
	assert.Equal(t, "someName", name)
	assert.Equal(t, "System.DateTime", varType)
	assert.Equal(t, "null", nullVal)
	assert.Equal(t, `	[SerializeField]
	private string someName;

	public System.DateTime SomeName { get => System.DateTime.Parse(someName); }
`, classVars)
}
