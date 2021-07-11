package path_test

import (
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen/path"
	"github.com/stretchr/testify/assert"
)

func Test_NumberResponse(t *testing.T) {
	// ARRANGE ================================================================
	desciption := "A bunch of cool cats"
	defResp := path.NewNumberResponse(desciption)

	// ACT ====================================================================
	desc := defResp.Description()
	interpret := defResp.Interpret("somethin", "download")

	// ASSERT =================================================================
	assert.Equal(t, desciption, desc)
	assert.Equal(t, "somethin = float.Parse(download.text);", interpret)
	assert.Equal(t, "float", defResp.VariableType())
}
