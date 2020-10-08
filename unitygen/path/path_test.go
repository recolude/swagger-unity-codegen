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
	classCode := route.UnityWebRequest()
	functionCode := route.ServiceFunction([]security.Auth{
		security.NewAPIKey("CognitoAuth", "X-API-KEY", security.Header),
	})
	requestParamsCode := route.RequestParamClass()

	// ********************************* ASSERT *******************************
	assert.Equal(t, `public class DevKeyService_GetDevKeyUnityWebRequest {

	public UnityWebRequest UnderlyingRequest{ get; }

	public DevKeyService_GetDevKeyUnityWebRequest(UnityWebRequest req) {
		this.UnderlyingRequest = req;
	}

	public IEnumerator Run() {
		yield return this.UnderlyingRequest.SendWebRequest();
	}

}`, classCode)

	assert.Equal(t, `public DevKeyService_GetDevKeyUnityWebRequest DevKeyService_GetDevKey()
{
	var unityNetworkReq = new UnityWebRequest(string.Format("{0}/api/v1/dev-keys", this.Config.BasePath), UnityWebRequest.kHttpVerbGET);
	unityNetworkReq.SetRequestHeader("X-API-KEY", this.Config.CognitoAuth);
	return new DevKeyService_GetDevKeyUnityWebRequest(unityNetworkReq);
}`, functionCode)

	assert.Equal(t, "", requestParamsCode)
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
	classCode := route.UnityWebRequest()
	functionCode := route.ServiceFunction([]security.Auth{
		security.NewAPIKey("CognitoAuth", "CognitoThing", security.Header),
		security.NewAPIKey("DevKeyAuth", "X-API-KEY", security.Header),
	})

	// ********************************* ASSERT *******************************
	assert.Equal(t, `public class UserService_GetUserUnityWebRequest {

	public V1UserResponse success;

	public RuntimeError fallbackResponse;

	public UnityWebRequest UnderlyingRequest{ get; }

	public UserService_GetUserUnityWebRequest(UnityWebRequest req) {
		this.UnderlyingRequest = req;
	}

	public IEnumerator Run() {
		yield return this.UnderlyingRequest.SendWebRequest();
		Interpret(this.UnderlyingRequest);
	}

	public void Interpret(UnityWebRequest req) {
		if (req.responseCode == 200) {
			success = JsonUtility.FromJson<V1UserResponse>(req.downloadHandler.text);
		} else {
			fallbackResponse = JsonUtility.FromJson<RuntimeError>(req.downloadHandler.text);
		}
	}

}`, classCode)

	assert.Equal(t, `public UserService_GetUserUnityWebRequest UserService_GetUser(UserService_GetUserRequestParams requestParams)
{
	var unityNetworkReq = requestParams.BuildUnityWebRequest(this.Config.BasePath);
	unityNetworkReq.downloadHandler = new DownloadHandlerBuffer();
	if (string.IsNullOrEmpty(this.Config.CognitoAuth) == false) {
		unityNetworkReq.SetRequestHeader("CognitoThing", this.Config.CognitoAuth);
	}
	if (string.IsNullOrEmpty(this.Config.DevKeyAuth) == false) {
		unityNetworkReq.SetRequestHeader("X-API-KEY", this.Config.DevKeyAuth);
	}
	return new UserService_GetUserUnityWebRequest(unityNetworkReq);
}

public UserService_GetUserUnityWebRequest UserService_GetUser(string userId)
{
	return UserService_GetUser(new UserService_GetUserRequestParams() {
		UserId=userId,
	});
}`, functionCode)
}

func Test_AcknowledgesSingleResponses(t *testing.T) {
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
			"200": path.NewResponse("A successful response.", definition.NewObjectReference("#/definitions/v1UserResponse")),
		},
		[]path.Parameter{
			path.NewParameter(path.PathParameterLocation, "userId", true, property.NewString("userId", "")),
		},
	)

	// ********************************** ACT *********************************
	classCode := route.UnityWebRequest()

	// ********************************* ASSERT *******************************
	assert.Equal(t, `public class UserService_GetUserUnityWebRequest {

	public V1UserResponse success;

	public UnityWebRequest UnderlyingRequest{ get; }

	public UserService_GetUserUnityWebRequest(UnityWebRequest req) {
		this.UnderlyingRequest = req;
	}

	public IEnumerator Run() {
		yield return this.UnderlyingRequest.SendWebRequest();
		Interpret(this.UnderlyingRequest);
	}

	public void Interpret(UnityWebRequest req) {
		if (req.responseCode == 200) {
			success = JsonUtility.FromJson<V1UserResponse>(req.downloadHandler.text);
		}
	}

}`, classCode)
}

func Test_AcknowledgesDefaultResponses(t *testing.T) {
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
			"default": path.NewResponse("An unexpected error response", definition.NewObjectReference("#/definitions/runtimeError")),
		},
		[]path.Parameter{
			path.NewParameter(path.PathParameterLocation, "userId", true, property.NewString("userId", "")),
		},
	)

	// ********************************** ACT *********************************
	classCode := route.UnityWebRequest()

	// ********************************* ASSERT *******************************
	assert.Equal(t, `public class UserService_GetUserUnityWebRequest {

	public RuntimeError fallbackResponse;

	public UnityWebRequest UnderlyingRequest{ get; }

	public UserService_GetUserUnityWebRequest(UnityWebRequest req) {
		this.UnderlyingRequest = req;
	}

	public IEnumerator Run() {
		yield return this.UnderlyingRequest.SendWebRequest();
		Interpret(this.UnderlyingRequest);
	}

	public void Interpret(UnityWebRequest req) {
		fallbackResponse = JsonUtility.FromJson<RuntimeError>(req.downloadHandler.text);
	}

}`, classCode)
}

func Test_ThreeParametersInPath(t *testing.T) {
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
			"401":     path.NewResponse("Weird Unauthorized response.", definition.NewObjectReference("#/definitions/v1Unauthorized")),
			"default": path.NewResponse("An unexpected error response", definition.NewObjectReference("#/definitions/runtimeError")),
		},
		[]path.Parameter{
			path.NewParameter(path.PathParameterLocation, "userId", true, property.NewString("userId", "")),
		},
	)

	// ********************************** ACT *********************************
	classCode := route.UnityWebRequest()

	// ********************************* ASSERT *******************************
	assert.Equal(t, `public class UserService_GetUserUnityWebRequest {

	public V1UserResponse success;

	public V1Unauthorized unauthorized;

	public RuntimeError fallbackResponse;

	public UnityWebRequest UnderlyingRequest{ get; }

	public UserService_GetUserUnityWebRequest(UnityWebRequest req) {
		this.UnderlyingRequest = req;
	}

	public IEnumerator Run() {
		yield return this.UnderlyingRequest.SendWebRequest();
		Interpret(this.UnderlyingRequest);
	}

	public void Interpret(UnityWebRequest req) {
		if (req.responseCode == 200) {
			success = JsonUtility.FromJson<V1UserResponse>(req.downloadHandler.text);
		} else if (req.responseCode == 401) {
			unauthorized = JsonUtility.FromJson<V1Unauthorized>(req.downloadHandler.text);
		} else {
			fallbackResponse = JsonUtility.FromJson<RuntimeError>(req.downloadHandler.text);
		}
	}

}`, classCode)

}

func Test_HandlesNilResponseDefinitions(t *testing.T) {
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
			"401":     path.NewResponse("Weird Unauthorized response.", nil),
			"default": path.NewResponse("An unexpected error response", definition.NewObjectReference("#/definitions/runtimeError")),
		},
		[]path.Parameter{
			path.NewParameter(path.PathParameterLocation, "userId", true, property.NewString("userId", "")),
		},
	)

	// ********************************** ACT *********************************
	classCode := route.UnityWebRequest()

	// ********************************* ASSERT *******************************
	assert.Equal(t, `public class UserService_GetUserUnityWebRequest {

	public V1UserResponse success;

	public RuntimeError fallbackResponse;

	public UnityWebRequest UnderlyingRequest{ get; }

	public UserService_GetUserUnityWebRequest(UnityWebRequest req) {
		this.UnderlyingRequest = req;
	}

	public IEnumerator Run() {
		yield return this.UnderlyingRequest.SendWebRequest();
		Interpret(this.UnderlyingRequest);
	}

	public void Interpret(UnityWebRequest req) {
		if (req.responseCode == 200) {
			success = JsonUtility.FromJson<V1UserResponse>(req.downloadHandler.text);
		} else {
			fallbackResponse = JsonUtility.FromJson<RuntimeError>(req.downloadHandler.text);
		}
	}

}`, classCode)

}

func Test_DealsWithQueryParams(t *testing.T) {
	// ******************************** ARRANGE *******************************
	route := path.NewPath(
		"/api/v1/users/{userId}",
		"UserService_GetUser",
		http.MethodGet,
		[]string{"UserService"},
		nil,
		map[string]path.Response{
			"default": path.NewResponse("An unexpected error response", definition.NewObjectReference("#/definitions/runtimeError")),
		},
		[]path.Parameter{
			path.NewParameter(path.PathParameterLocation, "userId", true, property.NewString("userId", "")),
			path.NewParameter(path.QueryParameterLocation, "diffId", true, property.NewString("diffId", "")),
		},
	)

	// ********************************** ACT *********************************
	functionCode := route.ServiceFunction(nil)

	// ********************************* ASSERT *******************************
	assert.Equal(t, `public UserService_GetUserUnityWebRequest UserService_GetUser(UserService_GetUserRequestParams requestParams)
{
	var unityNetworkReq = requestParams.BuildUnityWebRequest(this.Config.BasePath);
	unityNetworkReq.downloadHandler = new DownloadHandlerBuffer();
	return new UserService_GetUserUnityWebRequest(unityNetworkReq);
}

public UserService_GetUserUnityWebRequest UserService_GetUser(string userId, string diffId)
{
	return UserService_GetUser(new UserService_GetUserRequestParams() {
		UserId=userId,
		DiffId=diffId,
	});
}`, functionCode)
}

func Test_DealsWithMultipleQueryParams(t *testing.T) {
	// ******************************** ARRANGE *******************************
	route := path.NewPath(
		"/api/v1/users/{userId}",
		"UserService_GetUser",
		http.MethodGet,
		[]string{"UserService"},
		nil,
		map[string]path.Response{
			"default": path.NewResponse("An unexpected error response", definition.NewObjectReference("#/definitions/runtimeError")),
		},
		[]path.Parameter{
			path.NewParameter(path.PathParameterLocation, "userId", true, property.NewString("userId", "")),
			path.NewParameter(path.QueryParameterLocation, "diffId", true, property.NewString("diffId", "")),
			path.NewParameter(path.QueryParameterLocation, "anotherId", true, property.NewInteger("anotherId", "")),
		},
	)

	// ********************************** ACT *********************************
	functionCode := route.ServiceFunction(nil)
	requestParamsClass := route.RequestParamClass()

	// ********************************* ASSERT *******************************
	assert.Equal(t, `public UserService_GetUserUnityWebRequest UserService_GetUser(UserService_GetUserRequestParams requestParams)
{
	var unityNetworkReq = requestParams.BuildUnityWebRequest(this.Config.BasePath);
	unityNetworkReq.downloadHandler = new DownloadHandlerBuffer();
	return new UserService_GetUserUnityWebRequest(unityNetworkReq);
}

public UserService_GetUserUnityWebRequest UserService_GetUser(string userId, string diffId, int anotherId)
{
	return UserService_GetUser(new UserService_GetUserRequestParams() {
		UserId=userId,
		DiffId=diffId,
		AnotherId=anotherId,
	});
}`, functionCode)

	assert.Equal(t, `public class UserService_GetUserRequestParams
{
	private bool userIdSet = false;
	private string userId;
	public string UserId { get { return userId; } set { userIdSet = true; userId = value; } }
	public void UnsetUserId() { userId = null; userIdSet = false; }

	private bool diffIdSet = false;
	private string diffId;
	public string DiffId { get { return diffId; } set { diffIdSet = true; diffId = value; } }
	public void UnsetDiffId() { diffId = null; diffIdSet = false; }

	private bool anotherIdSet = false;
	private int anotherId;
	public int AnotherId { get { return anotherId; } set { anotherIdSet = true; anotherId = value; } }
	public void UnsetAnotherId() { anotherId = 0; anotherIdSet = false; }

	public UnityWebRequest BuildUnityWebRequest(string baseURL)
	{
		var finalPath = string.Format("{0}/api/v1/users/{userId}", baseURL);
		finalPath = finalPath.Replace("{userId}", userIdSet ? UnityWebRequest.EscapeURL(userId.ToString()) : "");
		var queryAdded = false;

		if (diffIdSet) {
			finalPath += (queryAdded ? "&" : "?") + "diffId=";
			queryAdded = true;
			finalPath += UnityWebRequest.EscapeURL(diffId.ToString());
		}

		if (anotherIdSet) {
			finalPath += (queryAdded ? "&" : "?") + "anotherId=";
			queryAdded = true;
			finalPath += UnityWebRequest.EscapeURL(anotherId.ToString());
		}

		return new UnityWebRequest(finalPath, UnityWebRequest.kHttpVerbGET);
	}
}`, requestParamsClass)

}
