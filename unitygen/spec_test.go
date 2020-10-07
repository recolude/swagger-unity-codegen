package unitygen_test

import (
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen"
	"github.com/recolude/swagger-unity-codegen/unitygen/security"
	"github.com/stretchr/testify/assert"
)

func TestSpec_ServiceConfig_NoScriptableObject(t *testing.T) {
	// ******************************** ARRANGE *******************************
	service := unitygen.Spec{}

	// ********************************** ACT *********************************
	code := service.ServiceConfig("ServiceConfig", "Server/Config", false)

	// ********************************* ASSERT *******************************
	assert.Equal(t, `public interface Config {

	// The base URL to which the endpoint paths are appended
	string BasePath { get; }

}`, code)
}

func TestSpec_ServiceConfig_NoSecurityDefinitions(t *testing.T) {
	// ******************************** ARRANGE *******************************
	service := unitygen.Spec{}

	// ********************************** ACT *********************************
	code := service.ServiceConfig("ServiceConfig", "Server/Config", true)

	// ********************************* ASSERT *******************************
	assert.Equal(t, `public interface Config {

	// The base URL to which the endpoint paths are appended
	string BasePath { get; }

}

#if UNITY_EDITOR
[UnityEditor.CustomEditor(typeof(ServiceConfig))]
public class ServiceConfigEditor : UnityEditor.Editor
{

	public override void OnInspectorGUI()
	{
		if (target == null)
		{
			return;
		}

		var castedTarget = (ServiceConfig)target;

		UnityEditor.EditorGUILayout.Space();
		UnityEditor.EditorGUILayout.LabelField("The base URL to which the endpoint paths are appended");
		var newBasePath = UnityEditor.EditorGUILayout.TextField("BasePath", castedTarget.BasePath);
		if (newBasePath != castedTarget.BasePath) {
			castedTarget.BasePath = newBasePath;
			UnityEditor.EditorUtility.SetDirty(target);
		}

	}

}
#endif

[System.Serializable]
[CreateAssetMenu(menuName = "Server/Config", fileName = "ServiceConfig")]
public class ServiceConfig: ScriptableObject, Config {

	[SerializeField]
	private string basePath;

	// The base URL to which the endpoint paths are appended
	public string BasePath { get { return basePath; } set { basePath = value; } }

}`, code)
}

func TestSpec_ServiceConfig_MultipleSecurityDefinitions(t *testing.T) {
	// ******************************** ARRANGE *******************************
	service := unitygen.Spec{
		AuthDefinitions: []security.Auth{
			security.NewAPIKey("anotherIdentifier", "DIF-KEY", security.Body),
			security.NewAPIKey("SomeIdentifier", "DA-KEY", security.Header),
		},
	}

	// ********************************** ACT *********************************
	code := service.ServiceConfig("RecoludeConfig", "Recolude/Config", true)

	// ********************************* ASSERT *******************************
	assert.Equal(t, `public interface Config {

	// The base URL to which the endpoint paths are appended
	string BasePath { get; }

	// AnotherIdentifier is a API Key 'DIF-KEY' found in a request's body
	string AnotherIdentifier { get; }

	// SomeIdentifier is a API Key 'DA-KEY' found in a request's header
	string SomeIdentifier { get; }

}

#if UNITY_EDITOR
[UnityEditor.CustomEditor(typeof(RecoludeConfig))]
public class RecoludeConfigEditor : UnityEditor.Editor
{

	public override void OnInspectorGUI()
	{
		if (target == null)
		{
			return;
		}

		var castedTarget = (RecoludeConfig)target;

		UnityEditor.EditorGUILayout.Space();
		UnityEditor.EditorGUILayout.LabelField("The base URL to which the endpoint paths are appended");
		var newBasePath = UnityEditor.EditorGUILayout.TextField("BasePath", castedTarget.BasePath);
		if (newBasePath != castedTarget.BasePath) {
			castedTarget.BasePath = newBasePath;
			UnityEditor.EditorUtility.SetDirty(target);
		}

		UnityEditor.EditorGUILayout.Space();
		UnityEditor.EditorGUILayout.LabelField("AnotherIdentifier is a API Key 'DIF-KEY' found in a request's body");
		var newAnotherIdentifier = UnityEditor.EditorGUILayout.TextField("AnotherIdentifier", castedTarget.AnotherIdentifier);
		if (newAnotherIdentifier != castedTarget.AnotherIdentifier) {
			castedTarget.AnotherIdentifier = newAnotherIdentifier;
			UnityEditor.EditorUtility.SetDirty(target);
		}

		UnityEditor.EditorGUILayout.Space();
		UnityEditor.EditorGUILayout.LabelField("SomeIdentifier is a API Key 'DA-KEY' found in a request's header");
		var newSomeIdentifier = UnityEditor.EditorGUILayout.TextField("SomeIdentifier", castedTarget.SomeIdentifier);
		if (newSomeIdentifier != castedTarget.SomeIdentifier) {
			castedTarget.SomeIdentifier = newSomeIdentifier;
			UnityEditor.EditorUtility.SetDirty(target);
		}

	}

}
#endif

[System.Serializable]
[CreateAssetMenu(menuName = "Recolude/Config", fileName = "RecoludeConfig")]
public class RecoludeConfig: ScriptableObject, Config {

	[SerializeField]
	private string basePath;

	// The base URL to which the endpoint paths are appended
	public string BasePath { get { return basePath; } set { basePath = value; } }

	[SerializeField]
	private string anotherIdentifier;

	// AnotherIdentifier is a API Key 'DIF-KEY' found in a request's body
	public string AnotherIdentifier { get { return anotherIdentifier; } set { anotherIdentifier = value; } }

	[SerializeField]
	private string someIdentifier;

	// SomeIdentifier is a API Key 'DA-KEY' found in a request's header
	public string SomeIdentifier { get { return someIdentifier; } set { someIdentifier = value; } }

}`, code)
}
