package path_test

import (
	"net/http"
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen/path"
	"github.com/recolude/swagger-unity-codegen/unitygen/security"
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
	classCode := route.SupportingClasses()
	functionCode := route.ServiceFunction([]security.Auth{
		security.NewAPIKey("CognitoAuth", "X-API-KEY", security.Header),
	})

	// ********************************* ASSERT *******************************
	assert.Equal(t, `public class DevKeyService_GetDevKeyUnityWebRequest {

	public UnityWebRequest UnderlyingRequest{ get; };

	public DevKeyService_GetDevKeyUnityWebRequest(UnityWebRequest req) {
		this.UnderlyingRequest = req;
	}

}`, classCode)

	assert.Equal(t, `public DevKeyService_GetDevKeyUnityWebRequest DevKeyService_GetDevKey()
{
	var unityNetworkReq = new UnityWebRequest(this.config.BasePath + "/api/v1/dev-keys", UnityWebRequest.kHttpVerbGET);
	unityNetworkReq.SetRequestHeader("X-API-KEY", this.config.Security.CognitoAuth());
	return new DevKeyService_GetDevKeyUnityWebRequest(unityNetworkReq);
}`, functionCode)
}
