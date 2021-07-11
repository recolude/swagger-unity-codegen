package path_test

import (
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen/model"
	"github.com/recolude/swagger-unity-codegen/unitygen/path"
	"github.com/stretchr/testify/assert"
)

func Test_DefinitionResponse(t *testing.T) {
	// ARRANGE ================================================================
	desciption := "A bunch of cool cats"
	schema := model.NewStringEnum("cool cats", []string{"cookie", "mortimer"})
	defResp := path.NewDefinitionResponse(desciption, schema)

	// ACT ====================================================================
	desc := defResp.Description()
	interpret := defResp.Interpret("somethin", "download")

	// ASSERT =================================================================
	assert.Equal(t, desciption, desc)
	assert.Equal(t, "somethin = JsonUtility.FromJson<CoolCats>(download.text);", interpret)
	assert.Equal(t, "CoolCats", defResp.VariableType())
}

func Test_DefinitionResponse_PanicsWithNilSchema(t *testing.T) {
	// ARRANGE ================================================================
	desciption := "A bunch of cool cats"
	defResp := path.NewDefinitionResponse(desciption, nil)

	// ACT ====================================================================
	assert.PanicsWithError(t, "can not build a response interpretation from nil definition", func() {
		defResp.Interpret("a", "b")
	})
	assert.PanicsWithError(t, "can not build a response variable type from nil definition", func() {
		defResp.VariableType()
	})
}
