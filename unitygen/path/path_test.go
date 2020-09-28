package path_test

import (
	"net/http"
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen/path"
	"github.com/stretchr/testify/assert"
)

func Test_SimpleGet(t *testing.T) {
	// ******************************** ARRANGE *******************************
	route := path.NewPath(
		"/api/v1/dev-keys",
		"DevKeyService_GetDevKey",
		http.MethodGet,
		[]string{"DevKeyService"},
		[]path.SecurityMethodReference{path.NewSecurityMethodReference("CognitoAuth")},
	)

	// ********************************** ACT *********************************
	out := route.SupportingClasses()

	// ********************************* ASSERT *******************************
	assert.Equal(t, `public class DevKeyService_GetDevKeyUnityWebRequest {

	private UnityEngine.Networking.UnityWebRequest UnderlyingRequest{ get; };

}`, out)
}
