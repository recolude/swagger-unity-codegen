package path_test

import (
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen/model"
	"github.com/recolude/swagger-unity-codegen/unitygen/path"
	"github.com/stretchr/testify/assert"
)

func Test_Response(t *testing.T) {
	// ARRANGE ================================================================
	desc := "test"
	def := model.NewStringEnum("woo", []string{"a", "b"})

	// ACT ====================================================================
	resp := path.NewResponse(desc, def)

	// ASSERT =================================================================
	assert.Equal(t, desc, resp.Description())
	assert.Equal(t, def, resp.Schema())
}
