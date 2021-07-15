package property_test

import (
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen/model"
	"github.com/recolude/swagger-unity-codegen/unitygen/model/property"
	"github.com/stretchr/testify/assert"
)

func Test_DefinitionReference(t *testing.T) {
	// ******************************** ARRANGE *******************************
	def := model.NewObject("v1Permission", nil)
	ref := property.NewDefinitionReference("my-permissions", def)

	// ********************************** ACT *********************************
	name := ref.Name()
	varType := ref.ToVariableType()
	nullVal := ref.EmptyValue()
	cSharp := ref.ClassVariables()

	// ********************************* ASSERT *******************************
	assert.Equal(t, "my-permissions", name)
	assert.Equal(t, "V1Permission", varType)
	assert.Equal(t, "null", nullVal)
	assert.Equal(t, `	[JsonProperty("my-permissions")]
	public V1Permission MyPermissions { get; private set; }
`, cSharp)
}

func Test_DefinitionReference_ShowsJSONConverter(t *testing.T) {
	// ******************************** ARRANGE *******************************
	def := model.NewStringEnum("cool cats", []string{"mortimer", "cookie"})
	ref := property.NewDefinitionReference("awesome cat", def)

	// ********************************** ACT *********************************
	name := ref.Name()
	varType := ref.ToVariableType()
	nullVal := ref.EmptyValue()
	cSharp := ref.ClassVariables()

	// ********************************* ASSERT *******************************
	assert.Equal(t, "awesome cat", name)
	assert.Equal(t, "CoolCats", varType)
	assert.Equal(t, "null", nullVal)
	assert.Equal(t, `	[JsonProperty("awesome cat")]
	[JsonConverter(typeof(CoolCatsJsonConverter))]
	public CoolCats AwesomeCat { get; private set; }
`, cSharp)
}
