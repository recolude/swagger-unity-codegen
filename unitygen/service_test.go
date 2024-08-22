package unitygen_test

import (
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen"
	"github.com/recolude/swagger-unity-codegen/unitygen/path"
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

func TestServiceInterfaceNoPaths(t *testing.T) {
	// ******************************** ARRANGE *******************************
	service := unitygen.NewService("test", nil)

	// ********************************** ACT *********************************
	code := service.Interface()

	// ********************************* ASSERT *******************************
	assert.Equal(t, `public interface ITestService {
}`, code)
}

func TestServiceInterfaceOnePath(t *testing.T) {
	// ******************************** ARRANGE *******************************
	service := unitygen.NewService("test", []path.Path{
		path.NewPath("/get", "doThing", "GET", nil, nil, nil, nil),
	})

	// ********************************** ACT *********************************
	code := service.Interface()

	// ********************************* ASSERT *******************************
	assert.Equal(t, `public interface ITestService {
	public DoThingUnityWebRequest DoThing(DoThingRequestParams requestParams);
}`, code)
}
