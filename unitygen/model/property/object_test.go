package property_test

import (
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen/model"
	"github.com/recolude/swagger-unity-codegen/unitygen/model/property"
	"github.com/stretchr/testify/assert"
)

func Test_Objec(t *testing.T) {
	// ******************************** ARRANGE *******************************
	anonObjPropName := "some-name"
	ref := property.NewObject(anonObjPropName, model.NewObject(anonObjPropName, []model.Property{
		property.NewBoolean("my bool"),
	}))

	// ********************************** ACT *********************************
	name := ref.Name()
	varType := ref.ToVariableType()
	nullVal := ref.EmptyValue()
	csharp := ref.ClassVariables()

	// ********************************* ASSERT *******************************
	assert.Equal(t, anonObjPropName, name)
	assert.Equal(t, "SomeName", varType)
	assert.Equal(t, "null", nullVal)
	assert.Equal(t, `	[System.Serializable]
public class SomeName {

	[JsonProperty("my bool")]
	public bool MyBool { get; private set; }

}
	[JsonProperty("some-name")]
	public SomeName SomeName { get; private set; }
`, csharp)
}
