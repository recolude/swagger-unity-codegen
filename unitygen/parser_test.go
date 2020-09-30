package unitygen_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen"
	"github.com/stretchr/testify/assert"
)

func TestReadDefinition(t *testing.T) {
	// ******************************** ARRANGE *******************************
	var swaggerDotJSON = `{
		"info": {
			"title": "Recolude Service",
			"version": "1.0"
		},
		"paths": {
			"/api/v1/dev-keys": {
				"get": {
					"security": [
						{
							"CognitoAuth": []
						}
					],	
					"tags": [
						"DevKeyService"
					],
					"operationId": "DevKeyService_GetDevKey",
					"responses": {
						"200": {
							"description": "A successful response.",
							"schema": {
								"$ref": "#/definitions/v1DevKeyResponse"
							}
						},
						"default": {
							"description": "An unexpected error response",
							"schema": {
								"$ref": "#/definitions/runtimeError"
							}
						}
					}
				},
				"post": {
					"security": [
						{
							"CognitoAuth": [],
							"DevKeyAuth": []
						}
					],
					"tags": [
						"DevKeyService"
					],
					"operationId": "DevKeyService_CreateDevKey",
					"responses": {
						"200": {
							"description": "A successful response.",
							"schema": {
								"$ref": "#/definitions/v1DevKeyResponse"
							}
						},
						"default": {
							"description": "An unexpected error response",
							"schema": {
								"$ref": "#/definitions/runtimeError"
							}
						}
					}
				}
			}
		},
		"definitions": {
			"v1ApiKey": {
				"type": "object",
				"properties": {
					"createdAt": {
						"type": "string",
						"format": "date-time"
					},
					"description": {
						"type": "string"
					},
					"name": {
						"type": "string"
					},
					"permissions": {
						"type": "array",
						"items": {
							"$ref": "#/definitions/v1Permission"
						}
					},
					"project": {
						"$ref": "#/definitions/v1Project"
					},
					"token": {
						"type": "string"
					}
				}
			},
			"v1EnumVisibility": {
				"type": "string",
				"default": "V_UNKOWN",
				"enum": [
				  "V_UNKOWN",
				  "V_PUBLIC",
				  "V_PRIVATE"
				]
			},
			"v1ListLicensesRequest": {
				"type": "object",
				"properties": {
				  "limit": {
					"type": "integer",
					"format": "int32"
				  },
				  "time": {
					"type": "number",
					"format": "float"
				  },
				  "duration": {
					"type": "number",
					"format": "double"
				  }
				}
			}
		},
		"securityDefinitions": {
		  "ApiKeyAuth": {
			"type": "apiKey",
			"name": "X-API-KEY",
			"in": "header"
		  },
		  "CognitoAuth": {
			"type": "apiKey",
			"name": "Authorization",
			"in": "header"
		  },
		  "DevKeyAuth": {
			"type": "apiKey",
			"name": "X-DEV-KEY",
			"in": "header"
		  }
		}
	}`
	// ********************************** ACT *********************************
	spec, err := unitygen.Parser{}.ParseJSON(strings.NewReader(swaggerDotJSON))

	// ********************************* ASSERT *******************************
	if assert.NoError(t, err) == false {
		return
	}
	assert.Equal(t, "Recolude Service", spec.Info.Title)
	assert.Equal(t, "1.0", spec.Info.Version)
	if assert.Len(t, spec.Definitions, 3) {
		assert.Equal(t, "v1ApiKey", spec.Definitions[0].Name())
		assert.Equal(t, "v1EnumVisibility", spec.Definitions[1].Name())
		assert.Equal(t,
			`[System.Serializable]
public class V1ApiKey {

	public System.DateTime createdAt;

	public string description;

	public string name;

	public V1Permission[] permissions;

	public V1Project project;

	public string token;

}`, spec.Definitions[0].ToCSharp())

		assert.Equal(t,
			`public enum V1EnumVisibility {
	V_UNKOWN,
	V_PUBLIC,
	V_PRIVATE
}`, spec.Definitions[1].ToCSharp())

		assert.Equal(t,
			`[System.Serializable]
public class V1ListLicensesRequest {

	public double duration;

	public int limit;

	public float time;

}`, spec.Definitions[2].ToCSharp())
	}

	if assert.Len(t, spec.AuthDefinitions, 3) {
		if assert.NotNil(t, spec.AuthDefinitions[0]) {
			assert.Equal(t, "ApiKeyAuth", spec.AuthDefinitions[0].Identifier())
		}
		if assert.NotNil(t, spec.AuthDefinitions[1]) {
			assert.Equal(t, "CognitoAuth", spec.AuthDefinitions[1].Identifier())
		}
		if assert.NotNil(t, spec.AuthDefinitions[2]) {
			assert.Equal(t, "DevKeyAuth", spec.AuthDefinitions[2].Identifier())
		}
	}

	if assert.Len(t, spec.Services, 1) {
		assert.Equal(t, "DevKeyService", spec.Services[0].Name())
		if assert.Len(t, spec.Services[0].Paths(), 2) {
			assert.Equal(t, "/api/v1/dev-keys", spec.Services[0].Paths()[0].Route())
			assert.Equal(t, "/api/v1/dev-keys", spec.Services[0].Paths()[1].Route())
			assert.Equal(t, http.MethodGet, spec.Services[0].Paths()[0].Method())
			assert.Equal(t, http.MethodPost, spec.Services[0].Paths()[1].Method())

			// Security References
			if assert.Len(t, spec.Services[0].Paths()[0].SecurityReferences(), 1) {
				assert.Equal(t, "CognitoAuth", spec.Services[0].Paths()[0].SecurityReferences()[0].Identifier)
			}
			if assert.Len(t, spec.Services[0].Paths()[1].SecurityReferences(), 2) {
				assert.Equal(t, "CognitoAuth", spec.Services[0].Paths()[1].SecurityReferences()[0].Identifier)
				assert.Equal(t, "DevKeyAuth", spec.Services[0].Paths()[1].SecurityReferences()[1].Identifier)
			}

			// Operation ID
			assert.Equal(t, "DevKeyService_GetDevKey", spec.Services[0].Paths()[0].OperationID())
			assert.Equal(t, "DevKeyService_CreateDevKey", spec.Services[0].Paths()[1].OperationID())

			if assert.Len(t, spec.Services[0].Paths()[0].Responses(), 2) {
				if assert.NotNil(t, spec.Services[0].Paths()[0].Responses()["200"].Schema()) {
					assert.Equal(t, "V1DevKeyResponse", spec.Services[0].Paths()[0].Responses()["200"].Schema().ToVariableType())
				}
				if assert.NotNil(t, spec.Services[0].Paths()[0].Responses()["default"].Schema()) {
					assert.Equal(t, "RuntimeError", spec.Services[0].Paths()[0].Responses()["default"].Schema().ToVariableType())
				}
			}

			if assert.Len(t, spec.Services[0].Paths()[1].Responses(), 2) {
				if assert.NotNil(t, spec.Services[0].Paths()[1].Responses()["200"].Schema()) {
					assert.Equal(t, "V1DevKeyResponse", spec.Services[0].Paths()[1].Responses()["200"].Schema().ToVariableType())
				}
				if assert.NotNil(t, spec.Services[0].Paths()[1].Responses()["default"].Schema()) {
					assert.Equal(t, "RuntimeError", spec.Services[0].Paths()[1].Responses()["default"].Schema().ToVariableType())
				}
			}

			// TODO: Parameters
		}

	}

}
