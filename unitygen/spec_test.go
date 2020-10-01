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
	assert.Equal(t, `[System.Serializable]
[CreateAssetMenu(menuName = "Server/Config", fileName = "ServiceConfig")]
public class ServiceConfig {

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
	assert.Equal(t, `[System.Serializable]
[CreateAssetMenu(menuName = "Recolude/Config", fileName = "RecoludeConfig")]
public class RecoludeConfig {

	public string BasePath { get; set; }

	// AnotherIdentifier is a API Key "DIF-KEY" found in a request's body
	public string AnotherIdentifier { get; set; }

	// SomeIdentifier is a API Key "DA-KEY" found in a request's header
	public string SomeIdentifier { get; set; }

	public RecoludeConfig(string basePath) {
		this.BasePath = basePath;
	}

}`, code)
}
