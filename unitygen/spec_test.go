package unitygen_test

import (
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen"
	"github.com/recolude/swagger-unity-codegen/unitygen/security"
	"github.com/stretchr/testify/assert"
)

func TestSpec_ServiceConfig_NoSecurityDefinitions(t *testing.T) {
	// ******************************** ARRANGE *******************************
	service := unitygen.Spec{}

	// ********************************** ACT *********************************
	code := service.ServiceConfig("ServiceConfig", "Server/Config")

	// ********************************* ASSERT *******************************
	assert.Equal(t, `public interface Config {

	// The base URL to which the endpoint paths are appended
	string BasePath { get; set; }

}

[System.Serializable]
[CreateAssetMenu(menuName = "Server/Config", fileName = "ServiceConfig")]
public class ServiceConfig: ScriptableObject, Config {

	// The base URL to which the endpoint paths are appended
	[SerializeField]
	public string BasePath { get; set; }

	public ServiceConfig(string basePath) {
		this.BasePath = basePath;
	}

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
	code := service.ServiceConfig("RecoludeConfig", "Recolude/Config")

	// ********************************* ASSERT *******************************
	assert.Equal(t, `public interface Config {

	// The base URL to which the endpoint paths are appended
	string BasePath { get; set; }

	// AnotherIdentifier is a API Key "DIF-KEY" found in a request's body
	string AnotherIdentifier { get; set; }

	// SomeIdentifier is a API Key "DA-KEY" found in a request's header
	string SomeIdentifier { get; set; }

}

[System.Serializable]
[CreateAssetMenu(menuName = "Recolude/Config", fileName = "RecoludeConfig")]
public class RecoludeConfig: ScriptableObject, Config {

	// The base URL to which the endpoint paths are appended
	[SerializeField]
	public string BasePath { get; set; }

	// AnotherIdentifier is a API Key "DIF-KEY" found in a request's body
	[SerializeField]
	public string AnotherIdentifier { get; set; }

	// SomeIdentifier is a API Key "DA-KEY" found in a request's header
	[SerializeField]
	public string SomeIdentifier { get; set; }

	public RecoludeConfig(string basePath) {
		this.BasePath = basePath;
	}

}`, code)
}
