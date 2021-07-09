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
	str := route.String()

	// ********************************* ASSERT *******************************
	assert.Equal(t, "someIdentifier", id)
	assert.Equal(t, `unityNetworkReq.SetRequestHeader("SomeName", this.Config.SomeIdentifier);`, modfier)
	assert.Equal(t, "SomeIdentifier is a API Key 'SomeName' found in a request's header", str)
}
