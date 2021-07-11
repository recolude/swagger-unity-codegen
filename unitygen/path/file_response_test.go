package path_test

import (
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen/path"
	"github.com/stretchr/testify/assert"
)

func Test_FileResponse(t *testing.T) {
	// ARRANGE ================================================================
	desciption := "A bunch of cool cats"
	defResp := path.NewFileResponse(desciption)

	// ACT ====================================================================
	desc := defResp.Description()
	interpret := defResp.Interpret("somethin", "download")

	// ASSERT =================================================================
	assert.Equal(t, desciption, desc)
	assert.Equal(t, "somethin = download.data;", interpret)
	assert.Equal(t, "byte[]", defResp.VariableType())
}
