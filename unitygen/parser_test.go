package unitygen_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen"
	"github.com/recolude/swagger-unity-codegen/unitygen/path"
	"github.com/stretchr/testify/assert"
)

func TestReadDefinition(t *testing.T) {
	// ******************************** ARRANGE *******************************
	swaggerDotJSON := `{
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
						"400": {
							"schema": {
								"type": "file"
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
				"default": "V_UNKNOWN",
				"enum": [
				  "V_UNKNOWN",
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

	[SerializeField]
	private string createdAt;

	public System.DateTime CreatedAt { get => System.DateTime.Parse(createdAt); }

	public string description;

	public string name;

	public V1Permission[] permissions;

	public V1Project project;

	public string token;

}`, spec.Definitions[0].ToCSharp())

		assert.Equal(t,
			`public enum V1EnumVisibility {
	V_UNKNOWN = 0,
	V_PUBLIC = 1,
	V_PRIVATE = 2
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

			if assert.Len(t, spec.Services[0].Paths()[0].Responses(), 3) {
				if assert.NotNil(t, spec.Services[0].Paths()[0].Responses()["200"]) {
					assert.Equal(t, "V1DevKeyResponse", spec.Services[0].Paths()[0].Responses()["200"].VariableType())
				}
				if assert.NotNil(t, spec.Services[0].Paths()[0].Responses()["default"]) {
					assert.Equal(t, "RuntimeError", spec.Services[0].Paths()[0].Responses()["default"].VariableType())
				}

				if assert.NotNil(t, spec.Services[0].Paths()[0].Responses()["400"]) {
					assert.Equal(t, "byte[]", spec.Services[0].Paths()[0].Responses()["400"].VariableType())
				}
			}

			if assert.Len(t, spec.Services[0].Paths()[1].Responses(), 2) {
				if assert.NotNil(t, spec.Services[0].Paths()[1].Responses()["200"]) {
					assert.Equal(t, "V1DevKeyResponse", spec.Services[0].Paths()[1].Responses()["200"].VariableType())
				}
				if assert.NotNil(t, spec.Services[0].Paths()[1].Responses()["default"]) {
					assert.Equal(t, "RuntimeError", spec.Services[0].Paths()[1].Responses()["default"].VariableType())
				}
			}

			assert.Len(t, spec.Services[0].Paths()[0].Parameters(), 0)
			assert.Len(t, spec.Services[0].Paths()[1].Parameters(), 0)
		}

	}
}

func Test_ReadParameters(t *testing.T) {
	// ******************************** ARRANGE *******************************
	swaggerDotJSON := `{
		"info": {
			"title": "Read Params",
			"version": "1.0.1"
		},
		"paths": {
			"/api/v1/echo": {
				"get": {
				  "security": [
					{
					  "ApiKeyAuth": [],
					  "CognitoAuth": [],
					  "DevKeyAuth": []
					}
				  ],
				  "description": "Echos your message back to you",
				  "tags": [
					"EchoService"
				  ],
				  "operationId": "EchoService_EchoString",
				  "parameters": [
					{
					  "type": "string",
					  "name": "value",
					  "in": "query"
					},
					{
					"type": "integer",
					"format": "int32",
					"name": "grantId",
					"in": "path",
					"required": true
					}
				  ],
				  "responses": {
					"200": {
					  "description": "A successful response.",
					  "schema": {
						"$ref": "#/definitions/v1Echo"
					  },
					  "examples": {
						"application/json": {
						  "value": "your message"
						}
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
					  "ApiKeyAuth": [],
					  "CognitoAuth": [],
					  "DevKeyAuth": []
					}
				  ],
				  "description": "Echos your message back to you",

				  "operationId": "EchoService_EchoString2",
				  "parameters": [
					{
					  "name": "body",
					  "in": "body",
					  "required": true,
					  "schema": {
						"$ref": "#/definitions/v1Echo"
					  }
					}
				  ],
				  "responses": {
					"200": {
					  "description": "A successful response.",
					  "schema": {
						"$ref": "#/definitions/v1Echo"
					  },
					  "examples": {
						"application/json": {
						  "value": "your message"
						}
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
		}
	}`
	// ********************************** ACT *********************************
	spec, err := unitygen.Parser{}.ParseJSON(strings.NewReader(swaggerDotJSON))

	// ********************************* ASSERT *******************************
	if assert.NoError(t, err) == false {
		return
	}
	assert.Equal(t, "Read Params", spec.Info.Title)
	assert.Equal(t, "1.0.1", spec.Info.Version)

	if assert.Len(t, spec.Services, 2) {
		assert.Equal(t, "Default", spec.Services[0].Name())
		if assert.Len(t, spec.Services[0].Paths(), 1) {
			if assert.Len(t, spec.Services[0].Paths()[0].Parameters(), 1) {
				assert.Equal(t, "body", spec.Services[0].Paths()[0].Parameters()[0].Name())
				assert.Equal(t, path.BodyParameterLocation, spec.Services[0].Paths()[0].Parameters()[0].Location())
				assert.Equal(t, true, spec.Services[0].Paths()[0].Parameters()[0].Required())
				assert.Equal(t, "V1Echo", spec.Services[0].Paths()[0].Parameters()[0].Schema().ToVariableType())
			}
		}

		assert.Equal(t, "EchoService", spec.Services[1].Name())
		if assert.Len(t, spec.Services[1].Paths(), 1) {
			if assert.Len(t, spec.Services[1].Paths()[0].Parameters(), 2) {
				assert.Equal(t, "value", spec.Services[1].Paths()[0].Parameters()[0].Name())
				assert.Equal(t, path.QueryParameterLocation, spec.Services[1].Paths()[0].Parameters()[0].Location())
				assert.Equal(t, false, spec.Services[1].Paths()[0].Parameters()[0].Required())
				assert.Equal(t, "string", spec.Services[1].Paths()[0].Parameters()[0].Schema().ToVariableType())

				assert.Equal(t, "grantId", spec.Services[1].Paths()[0].Parameters()[1].Name())
				assert.Equal(t, path.PathParameterLocation, spec.Services[1].Paths()[0].Parameters()[1].Location())
				assert.Equal(t, true, spec.Services[1].Paths()[0].Parameters()[1].Required())
				assert.Equal(t, "int", spec.Services[1].Paths()[0].Parameters()[1].Schema().ToVariableType())
			}
		}

	}
}

func Test_ReadNestedObjectPropertyDefinition(t *testing.T) {
	// ******************************** ARRANGE *******************************
	swaggerDotJSON := `{
			"info": {
				"title": "Analytic Service",
				"version": "1.0.0"
			},
			"paths": {
			},
			"definitions": {
				"AggMetadataQuery": {
					"type": "object",
					"properties": {
					  "query": {
						"type": "object",
						"properties": {
						  "field": {
							"type": "string"
						  },
						  "modifier": {
							"$ref": "#/definitions/AggModifier"
						  },
						  "minDate": {
							"type": "integer",
							"format": "int32"
						  },
						  "maxDate": {
							"type": "integer",
							"format": "int32"
						  },
						  "onEntity": {
							"$ref": "#/definitions/RecordingEntity"
						  }
						},
						"required": ["field", "modifier"]
					  }
					}
				  }
			}
		}`
	// ********************************** ACT *********************************
	spec, err := unitygen.Parser{}.ParseJSON(strings.NewReader(swaggerDotJSON))

	// ********************************* ASSERT *******************************
	if assert.NoError(t, err) == false {
		return
	}
	assert.Equal(t, "Analytic Service", spec.Info.Title)
	assert.Equal(t, "1.0.0", spec.Info.Version)
	if assert.Len(t, spec.Definitions, 1) {
		def := spec.Definitions[0]
		assert.Equal(t, "AggMetadataQuery", def.Name())
		assert.Equal(t, "AggMetadataQuery", def.ToVariableType())
		assert.Equal(t, `[System.Serializable]
public class AggMetadataQuery {

	[System.Serializable]
public class Query {

	public string field;

	public int maxDate;

	public int minDate;

	public AggModifier modifier;

	public RecordingEntity onEntity;

}
	public Query query;
}`, def.ToCSharp())
	}
}

func Test_ReadNumberEnums(t *testing.T) {
	// ******************************** ARRANGE *******************************
	swaggerDotJSON := `{
			"info": {
				"title": "Analytic Service",
				"version": "1.0.0"
			},
			"paths": {
			},
			"definitions": {
				"BinSize": {
					"type": "number",
					"enum": [0.125, 0.25, 0.5, 1, 2, 4, 8]
				}
			}
		}`
	// ********************************** ACT *********************************
	spec, err := unitygen.Parser{}.ParseJSON(strings.NewReader(swaggerDotJSON))

	// ********************************* ASSERT *******************************
	if assert.NoError(t, err) == false {
		return
	}
	assert.Equal(t, "Analytic Service", spec.Info.Title)
	assert.Equal(t, "1.0.0", spec.Info.Version)
	if assert.Len(t, spec.Definitions, 1) {
		def := spec.Definitions[0]
		assert.Equal(t, "BinSize", def.Name())
		assert.Equal(t, "BinSize", def.ToVariableType())
		assert.Equal(t, `public enum BinSize {
	NUMBER_0_DOT_125,
	NUMBER_0_DOT_25,
	NUMBER_0_DOT_5,
	NUMBER_1,
	NUMBER_2,
	NUMBER_4,
	NUMBER_8
}`, def.ToCSharp())
	}
}
