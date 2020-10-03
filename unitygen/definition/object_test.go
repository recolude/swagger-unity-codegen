package definition_test

import (
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen/definition"
	"github.com/recolude/swagger-unity-codegen/unitygen/property"
	"github.com/stretchr/testify/assert"
)

func TestObject(t *testing.T) {
	// ******************************** ARRANGE *******************************
	obj := definition.NewObject(
		"testObj",
		[]property.Property{
			property.NewInteger("num", ""),
		},
	)

	// ********************************** ACT *********************************
	varType := obj.ToVariableType()
	name := obj.Name()
	cSharp := obj.ToCSharp()

	// ********************************* ASSERT *******************************
	assert.Equal(t, "TestObj", varType)
	assert.Equal(t, "testObj", name)
	assert.Equal(t, `[System.Serializable]
public class TestObj {

	public int num;

}`, cSharp)
}