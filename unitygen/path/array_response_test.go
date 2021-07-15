package path_test

import (
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen/model/property"
	"github.com/recolude/swagger-unity-codegen/unitygen/path"
	"github.com/stretchr/testify/assert"
)

func Test_ArrayResponse(t *testing.T) {
	// ARRANGE ================================================================
	desciption := "A bunch of cool cats"
	schema := property.NewArray("cats", property.NewInteger("cats", "int"))
	defResp := path.NewArrayResponse(desciption, schema)

	// ACT ====================================================================
	desc := defResp.Description()
	interpret := defResp.Interpret("somethin", "download")

	// ASSERT =================================================================
	assert.Equal(t, desciption, desc)
	assert.Equal(t, "somethin = JsonConvert.DeserializeObject<int[]>(download.text);", interpret)
	assert.Equal(t, "int[]", defResp.VariableType())
}
