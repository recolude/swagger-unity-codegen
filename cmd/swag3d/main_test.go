package main

import (
	"os"
	"strings"
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen"
	"github.com/recolude/swagger-unity-codegen/unitygen/model"
	"github.com/recolude/swagger-unity-codegen/unitygen/model/property"
	"github.com/recolude/swagger-unity-codegen/unitygen/path"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func AssertFileExists(t *testing.T, fs afero.Fs, filename string) bool {
	info, err := fs.Stat(filename)
	if os.IsNotExist(err) {
		t.Errorf("expected file \"%s\" to exist but doesn't", filename)
		return false
	}
	if info.IsDir() {
		t.Errorf("expected \"%s\" to to be file but is instead directory", filename)
	}
	return true
}

func TestFilterServiceByTags_DoesNothingWithNoTags(t *testing.T) {
	// ******************************** ARRANGE *******************************
	spec := unitygen.NewSpec(unitygen.SpecInfo{}, nil, nil, []unitygen.Service{
		unitygen.NewService("A", nil),
		unitygen.NewService("B", nil),
	})

	// ********************************** ACT *********************************
	out := filterSpecForTags(spec, nil)

	// ********************************* ASSERT *******************************
	if assert.Len(t, out.Services, 2) {
		assert.Equal(t, out.Services[0].Name(), "A")
		assert.Equal(t, out.Services[1].Name(), "B")
	}
}

func TestFilterServiceByTags_Filters(t *testing.T) {
	// ******************************** ARRANGE *******************************
	spec := unitygen.NewSpec(unitygen.SpecInfo{}, nil, nil, []unitygen.Service{
		unitygen.NewService("A", nil),
		unitygen.NewService("B", nil),
	})

	// ********************************** ACT *********************************
	out := filterSpecForTags(spec, []string{"A"})

	// ********************************* ASSERT *******************************
	if assert.Len(t, out.Services, 1) {
		assert.Equal(t, out.Services[0].Name(), "A")
	}
}

func TestFilterUnusedDefinitions(t *testing.T) {
	// ******************************** ARRANGE *******************************
	spec := unitygen.NewSpec(
		unitygen.SpecInfo{},
		[]model.Definition{
			model.NewStringEnum("SomeEnum", []string{"OneEnum", "TwoEnum"}),
			model.NewStringEnum("FurtherRemovedEnum", []string{"OneEnum", "TwoEnum"}),
			model.NewObject("ToBeRemoved", nil),
			model.NewObject("Basic", nil),
			model.NewObject("X", nil),
			model.NewObject("Y", nil),
			model.NewObject("A", []model.Property{property.NewObjectReference("a", "#/definitions/B")}),
			model.NewObject("B", []model.Property{property.NewObjectReference("a", "#/definitions/C")}),
			model.NewObject("C", []model.Property{
				property.NewObjectReference("a", "#/definitions/FurtherRemovedEnum"),
				property.NewArray("test", property.NewObjectReference("a", "#/definitions/D")),
			}),
			model.NewObject("D", nil),
			model.NewObject("RecurseL", []model.Property{property.NewObjectReference("a", "#/definitions/RecurseR")}),
			model.NewObject("RecurseR", []model.Property{property.NewObjectReference("a", "#/definitions/RecurseL")}),
			model.NewObject("RecurseM", []model.Property{property.NewObjectReference("a", "#/definitions/RecurseM")}),
		},
		nil,
		[]unitygen.Service{
			unitygen.NewService(
				"A",
				[]path.Path{
					path.NewPath(
						"aaerg",
						"",
						"",
						nil,
						nil,
						map[string]path.Response{
							"200": path.NewDefinitionResponse("", model.NewObjectReference("#/definitions/A")),
							"400": path.NewDefinitionResponse("", model.NewObjectReference("#/definitions/SomeEnum")),
							"500": path.NewDefinitionResponse("", model.NewObjectReference("#/definitions/RecurseM")),
							"501": path.NewDefinitionResponse("", model.NewObjectReference("#/definitions/RecurseL")),
						},
						[]path.Parameter{
							path.NewParameter(
								path.BodyParameterLocation,
								"",
								false,
								property.NewObjectReference("a", "#/definitions/X"),
							),
						},
					),
				},
			),
			unitygen.NewService(
				"B",
				[]path.Path{
					path.NewPath(
						"aaerg",
						"",
						"",
						nil,
						nil,
						map[string]path.Response{
							"200": path.NewDefinitionResponse("", model.NewObjectReference("#/definitions/Basic")),
						},
						[]path.Parameter{
							path.NewParameter(
								path.BodyParameterLocation,
								"",
								false,
								property.NewObjectReference("a", "#/definitions/Y"),
							),
						},
					),
				},
			),
		},
	)

	// ********************************** ACT *********************************
	out := filterSpecForUnusedDefinitions(spec)

	// ********************************* ASSERT *******************************
	assert.Len(t, out.Services, 2)
	assert.Len(t, out.Definitions, 12)
}

func TestNoNamespace(t *testing.T) {
	// ******************************** ARRANGE *******************************
	appFS := afero.NewMemMapFs()
	afero.WriteFile(appFS, "swagger.json", []byte("{ }"), os.ModePerm)

	out := strings.Builder{}
	errOut := strings.Builder{}
	app := buildApp(appFS, &out, &errOut)

	// ********************************** ACT *********************************
	err := app.Run([]string{"swag3d", "generate", "--file", "swagger.json"})

	// ********************************* ASSERT *******************************
	assert.NoError(t, err)
	assert.Equal(t, "", errOut.String())
	assert.Equal(t, `// This code was generated by: 
// https://github.com/recolude/swagger-unity-codegen
// Issues and PRs welcome :)

using UnityEngine;
using UnityEngine.Networking;
using System.Collections;
using System.Text;
using Newtonsoft.Json;

#region Definitions

#endregion

#region Services

public interface Config {

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

}

#endregion

`, out.String())
}

func TestWithNamespace(t *testing.T) {
	// ******************************** ARRANGE *******************************
	appFS := afero.NewMemMapFs()
	afero.WriteFile(appFS, "swagger.json", []byte("{ }"), os.ModePerm)

	out := strings.Builder{}
	errOut := strings.Builder{}
	app := buildApp(appFS, &out, &errOut)

	// ********************************** ACT *********************************
	err := app.Run([]string{"swag3d", "generate", "--file", "swagger.json", "--namespace", "example"})

	// ********************************* ASSERT *******************************
	assert.NoError(t, err)
	assert.Equal(t, "", errOut.String())
	assert.Equal(t, `// This code was generated by: 
// https://github.com/recolude/swagger-unity-codegen
// Issues and PRs welcome :)

using UnityEngine;
using UnityEngine.Networking;
using System.Collections;
using System.Text;
using Newtonsoft.Json;

namespace Example {

#region Definitions

#endregion

#region Services

public interface Config {

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

}

#endregion

}`, out.String())
}

func TestSpecifyingOutWritesMultipleFiles(t *testing.T) {
	// ******************************** ARRANGE *******************************
	appFS := afero.NewMemMapFs()
	afero.WriteFile(appFS, "swagger.json", []byte("{ }"), os.ModePerm)

	out := strings.Builder{}
	errOut := strings.Builder{}
	app := buildApp(appFS, &out, &errOut)

	// ********************************** ACT *********************************
	err := app.Run([]string{"swag3d", "generate", "--file", "swagger.json", "--out", "."})

	// ********************************* ASSERT *******************************
	assert.NoError(t, err)
	AssertFileExists(t, appFS, "Definitions.cs")
	AssertFileExists(t, appFS, "Services.cs")
	AssertFileExists(t, appFS, "ServiceConfig.cs")
}

func TestErrorsWithNoFileToReadFrom(t *testing.T) {
	// ******************************** ARRANGE *******************************
	appFS := afero.NewMemMapFs()
	// afero.WriteFile(appFS, "swagger.json", []byte("{ }"), os.ModePerm)

	out := strings.Builder{}
	errOut := strings.Builder{}
	app := buildApp(appFS, &out, &errOut)

	// ********************************** ACT *********************************
	err := app.Run([]string{"swag3d", "generate", "--file", "swagger.json", "--out", "."})

	// ********************************* ASSERT *******************************
	assert.EqualError(t, err, "Error opening swagger file: open swagger.json: file does not exist")
}
