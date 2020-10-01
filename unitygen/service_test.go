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
	code := service.ToCSharp(nil)

	// ********************************* ASSERT *******************************
	assert.Equal(t, `public class TestService {

	public ServiceConfig Config { get; }

	public TestService(ServiceConfig Config) {
		this.Config = Config;
	}

}`, code)
}
