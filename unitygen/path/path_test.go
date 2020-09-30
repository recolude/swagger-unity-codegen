package path_test

import (
	"net/http"
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen/definition"
	"github.com/recolude/swagger-unity-codegen/unitygen/path"
	"github.com/recolude/swagger-unity-codegen/unitygen/property"
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
		nil,
		nil,
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
	var unityNetworkReq = new UnityWebRequest(string.Format("{0}/api/v1/dev-keys", this.config.BasePath), UnityWebRequest.kHttpVerbGET);
	unityNetworkReq.SetRequestHeader("X-API-KEY", this.config.Security.CognitoAuth());
	return new DevKeyService_GetDevKeyUnityWebRequest(unityNetworkReq);
}`, functionCode)
}

func Test_ParameterInPath(t *testing.T) {
	// ******************************** ARRANGE *******************************
	route := path.NewPath(
		"/api/v1/users/{userId}",
		"UserService_GetUser",
		http.MethodGet,
		[]string{"UserService"},
		[]path.SecurityMethodReference{
			path.NewSecurityMethodReference("CognitoAuth"),
			path.NewSecurityMethodReference("DevKeyAuth"),
		},
		map[string]path.Response{
			"200":     path.NewResponse("A successful response.", definition.NewObjectReference("#/definitions/v1UserResponse")),
			"default": path.NewResponse("An unexpected error response", definition.NewObjectReference("#/definitions/runtimeError")),
		},
		[]path.Parameter{
			path.NewParameter(path.PathParameterLocation, "userId", true, property.NewString("userId", "")),
		},
	)

	// ********************************** ACT *********************************
	classCode := route.SupportingClasses()
	functionCode := route.ServiceFunction([]security.Auth{
		security.NewAPIKey("CognitoAuth", "CognitoThing", security.Header),
		security.NewAPIKey("DevKeyAuth", "X-API-KEY", security.Header),
	})

	// ********************************* ASSERT *******************************
	assert.Equal(t, `public class UserService_GetUserUnityWebRequest {

	public V1UserResponse success;

	public RuntimeError default;

	public UnityWebRequest UnderlyingRequest{ get; };

	public UserService_GetUserUnityWebRequest(UnityWebRequest req) {
		this.UnderlyingRequest = req;
	}

}`, classCode)

	assert.Equal(t, `public UserService_GetUserUnityWebRequest UserService_GetUser(string userId)
{
	var unityNetworkReq = new UnityWebRequest(string.Format("{0}/api/v1/users/{1}", this.config.BasePath, userId), UnityWebRequest.kHttpVerbGET);
	if (string.IsNullOrEmpty(this.config.Security.CognitoAuth()) == false) {
		unityNetworkReq.SetRequestHeader("CognitoThing", this.config.Security.CognitoAuth());
	}
	if (string.IsNullOrEmpty(this.config.Security.DevKeyAuth()) == false) {
		unityNetworkReq.SetRequestHeader("X-API-KEY", this.config.Security.DevKeyAuth());
	}
	return new UserService_GetUserUnityWebRequest(unityNetworkReq);
}`, functionCode)
}
