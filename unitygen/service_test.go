package unitygen_test

import (
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen"
	"github.com/stretchr/testify/assert"
)

func TestServiceNoPaths(t *testing.T) {
	// ******************************** ARRANGE *******************************
	service := unitygen.NewService("test", nil)

	// ********************************** ACT *********************************
	code := service.ToCSharp(nil, "ServiceConfig")

	// ********************************* ASSERT *******************************
	assert.Equal(t, `public class TestService {

	public ServiceConfig Config { get; }

	public TestService(ServiceConfig Config) {
		this.Config = Config;
	}

}`, code)
}

func TestService_DoesntAppend2ndServiceToName(t *testing.T) {
	// ******************************** ARRANGE *******************************
	service := unitygen.NewService("testService", nil)

	// ********************************** ACT *********************************
	code := service.ToCSharp(nil, "RecoludeConfig")

	// ********************************* ASSERT *******************************
	assert.Equal(t, `public class TestService {

	public RecoludeConfig Config { get; }

	public TestService(RecoludeConfig Config) {
		this.Config = Config;
	}

}`, code)
}
