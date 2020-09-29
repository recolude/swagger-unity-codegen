package security_test

import (
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen/security"
	"github.com/stretchr/testify/assert"
)

func Test_APIKey(t *testing.T) {
	// ******************************** ARRANGE *******************************
	route := security.NewAPIKey("someIdentifier", "SomeName", security.Header)

	// ********************************** ACT *********************************
	id := route.Identifier()
	modfier := route.ModifyNetworkRequest()

	// ********************************* ASSERT *******************************
	assert.Equal(t, "someIdentifier", id)
	assert.Equal(t, `unityNetworkReq.SetRequestHeader("SomeName", this.config.Security.SomeIdentifier());`, modfier)
}
