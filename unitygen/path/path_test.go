package path_test

import (
	"net/http"
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen/model"
	"github.com/recolude/swagger-unity-codegen/unitygen/model/property"
	"github.com/recolude/swagger-unity-codegen/unitygen/path"
	"github.com/recolude/swagger-unity-codegen/unitygen/security"
	"github.com/stretchr/testify/assert"
)

func Test_PanicsWithMultipleBodyParameters(t *testing.T) {
	assert.PanicsWithError(t, "can not have multiple body parameters for a single path", func() {
		path.NewPath(
			"/api/v1/dev-keys",
			"DevKeyService_GetDevKey",
			http.MethodGet,
			[]string{"DevKeyService"},
			[]path.SecurityMethodReference{path.NewSecurityMethodReference("CognitoAuth")},
			nil,
			[]path.Parameter{
				path.NewParameter(path.BodyParameterLocation, "1", true, nil),
				path.NewParameter(path.BodyParameterLocation, "2", true, nil),
			},
		)
	})
}

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
			"200":     path.NewDefinitionResponse("A successful response.", model.NewObjectReference("#/definitions/v1UserResponse")),
			"default": path.NewDefinitionResponse("An unexpected error response", model.NewObjectReference("#/definitions/runtimeError")),
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

	// A successful response.
	public V1UserResponse success;

	// An unexpected error response
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
			"200": path.NewDefinitionResponse("", model.NewObjectReference("#/definitions/v1UserResponse")),
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
			"default": path.NewDefinitionResponse("An unexpected error response", model.NewObjectReference("#/definitions/runtimeError")),
		},
		[]path.Parameter{
			path.NewParameter(path.PathParameterLocation, "userId", true, property.NewString("userId", "")),
		},
	)

	// ********************************** ACT *********************************
	classCode := route.UnityWebRequest()

	// ********************************* ASSERT *******************************
	assert.Equal(t, `public class UserService_GetUserUnityWebRequest {

	// An unexpected error response
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
			"200":     path.NewDefinitionResponse("A successful response.", model.NewObjectReference("#/definitions/v1UserResponse")),
			"401":     path.NewDefinitionResponse("Weird Unauthorized response.", model.NewObjectReference("#/definitions/v1Unauthorized")),
			"default": path.NewDefinitionResponse("An unexpected error response", model.NewObjectReference("#/definitions/runtimeError")),
		},
		[]path.Parameter{
			path.NewParameter(path.PathParameterLocation, "userId", true, property.NewString("userId", "")),
		},
	)

	// ********************************** ACT *********************************
	classCode := route.UnityWebRequest()

	// ********************************* ASSERT *******************************
	assert.Equal(t, `public class UserService_GetUserUnityWebRequest {

	// A successful response.
	public V1UserResponse success;

	// Weird Unauthorized response.
	public V1Unauthorized unauthorized;

	// An unexpected error response
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
	params := []path.Parameter{path.NewParameter(path.PathParameterLocation, "userId", true, property.NewString("userId", ""))}
	opID := "UserService_GetUser"
	urlRotue := "/api/v1/users/{userId}"
	method := http.MethodGet
	security := []path.SecurityMethodReference{
		path.NewSecurityMethodReference("CognitoAuth"),
		path.NewSecurityMethodReference("DevKeyAuth"),
	}
	tags := []string{"UserService"}

	route := path.NewPath(
		urlRotue,
		opID,
		method,
		tags,
		security,
		map[string]path.Response{
			"200":     path.NewDefinitionResponse("A successful response.", model.NewObjectReference("#/definitions/v1UserResponse")),
			"401":     nil,
			"501":     path.NewFileResponse("some file"),
			"default": path.NewDefinitionResponse("An unexpected error response", model.NewObjectReference("#/definitions/runtimeError")),
		},
		params,
	)

	// ********************************** ACT *********************************
	classCode := route.UnityWebRequest()

	// ********************************* ASSERT *******************************
	assert.Equal(t, params, route.Parameters())
	assert.Equal(t, urlRotue, route.Route())
	assert.Equal(t, opID, route.OperationID())
	assert.Equal(t, method, route.Method())
	assert.Equal(t, security, route.SecurityReferences())
	assert.Equal(t, tags, route.Tags())
	assert.Equal(t, `public class UserService_GetUserUnityWebRequest {

	// A successful response.
	public V1UserResponse success;

	// some file
	public byte[] notImplemented;

	// An unexpected error response
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
		} else if (req.responseCode == 501) {
			notImplemented = req.downloadHandler.data;
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
			"default": path.NewDefinitionResponse("An unexpected error response", model.NewObjectReference("#/definitions/runtimeError")),
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

func Test_DealsWithMultipleQueryParamsAndBody(t *testing.T) {
	// ******************************** ARRANGE *******************************
	route := path.NewPath(
		"/api/v1/users/{userId}/{user-name}",
		"getUser",
		http.MethodGet,
		[]string{"UserService"},
		nil,
		map[string]path.Response{
			"default": path.NewDefinitionResponse("An unexpected error response", model.NewObjectReference("#/definitions/runtimeError")),
		},
		[]path.Parameter{
			path.NewParameter(path.PathParameterLocation, "userId", true, property.NewString("userId", "")),
			path.NewParameter(path.PathParameterLocation, "user-name", true, property.NewString("userId", "")),
			path.NewParameter(path.QueryParameterLocation, "diffId", true, property.NewString("diffId", "")),
			path.NewParameter(path.QueryParameterLocation, "anotherId", true, property.NewInteger("anotherId", "")),
			path.NewParameter(path.BodyParameterLocation, "query", true, property.NewObjectReference("test", "#/definitions/Query")),
		},
	)

	// ********************************** ACT *********************************
	functionCode := route.ServiceFunction(nil)
	supportingClasses := route.SupportingClasses()

	// ********************************* ASSERT *******************************
	assert.Equal(t, `public GetUserUnityWebRequest GetUser(GetUserRequestParams requestParams)
{
	var unityNetworkReq = requestParams.BuildUnityWebRequest(this.Config.BasePath);
	unityNetworkReq.downloadHandler = new DownloadHandlerBuffer();
	return new GetUserUnityWebRequest(unityNetworkReq);
}

public GetUserUnityWebRequest GetUser(string userId, string userName, string diffId, int anotherId, Query query)
{
	return GetUser(new GetUserRequestParams() {
		UserId=userId,
		UserName=userName,
		DiffId=diffId,
		AnotherId=anotherId,
		Query=query,
	});
}`, functionCode)

	assert.Equal(t, `public class GetUserUnityWebRequest {

	// An unexpected error response
	public RuntimeError fallbackResponse;

	public UnityWebRequest UnderlyingRequest{ get; }

	public GetUserUnityWebRequest(UnityWebRequest req) {
		this.UnderlyingRequest = req;
	}

	public IEnumerator Run() {
		yield return this.UnderlyingRequest.SendWebRequest();
		Interpret(this.UnderlyingRequest);
	}

	public void Interpret(UnityWebRequest req) {
		fallbackResponse = JsonUtility.FromJson<RuntimeError>(req.downloadHandler.text);
	}

}
public class GetUserRequestParams
{
	private bool userIdSet = false;
	private string userId;
	public string UserId { get { return userId; } set { userIdSet = true; userId = value; } }
	public void UnsetUserId() { userId = null; userIdSet = false; }

	private bool userNameSet = false;
	private string userName;
	public string UserName { get { return userName; } set { userNameSet = true; userName = value; } }
	public void UnsetUserName() { userName = null; userNameSet = false; }

	private bool diffIdSet = false;
	private string diffId;
	public string DiffId { get { return diffId; } set { diffIdSet = true; diffId = value; } }
	public void UnsetDiffId() { diffId = null; diffIdSet = false; }

	private bool anotherIdSet = false;
	private int anotherId;
	public int AnotherId { get { return anotherId; } set { anotherIdSet = true; anotherId = value; } }
	public void UnsetAnotherId() { anotherId = 0; anotherIdSet = false; }

	private bool querySet = false;
	private Query query;
	public Query Query { get { return query; } set { querySet = true; query = value; } }
	public void UnsetQuery() { query = null; querySet = false; }

	public UnityWebRequest BuildUnityWebRequest(string baseURL)
	{
		var finalPath = baseURL + "/api/v1/users/{userId}/{user-name}";
		finalPath = finalPath.Replace("{userId}", userIdSet ? UnityWebRequest.EscapeURL(userId.ToString()) : "");
		finalPath = finalPath.Replace("{user-name}", userNameSet ? UnityWebRequest.EscapeURL(userName.ToString()) : "");
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

		var unityWebReq = new UnityWebRequest(finalPath, UnityWebRequest.kHttpVerbGET);
		var unityRawUploadHandler = new UploadHandlerRaw(Encoding.Unicode.GetBytes(JsonConvert.SerializeObject(query)));
		unityRawUploadHandler.contentType = "application/json";
		unityWebReq.uploadHandler = unityRawUploadHandler;
		return unityWebReq;
	}
}`, supportingClasses)
}
