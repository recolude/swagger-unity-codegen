package unitygen_test

import (
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
			}
		}
	}`
	// ********************************** ACT *********************************
	spec, err := unitygen.Parser{}.Parse(strings.NewReader(swaggerDotJSON))

	// ********************************* ASSERT *******************************
	assert.NoError(t, err)
	assert.Equal(t, "Recolude Service", spec.Info.Title)
	assert.Equal(t, "1.0", spec.Info.Version)
	if assert.Len(t, spec.Definitions, 2) {
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

}`,
			spec.Definitions[0].ToClass())
	}
}
