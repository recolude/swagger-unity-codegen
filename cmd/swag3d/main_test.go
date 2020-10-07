package main

import (
	"os"
	"strings"
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen"
	"github.com/recolude/swagger-unity-codegen/unitygen/definition"
	"github.com/recolude/swagger-unity-codegen/unitygen/path"
	"github.com/recolude/swagger-unity-codegen/unitygen/property"
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
		[]definition.Definition{
			definition.NewEnum("SomeEnum", []string{"OneEnum", "TwoEnum"}),
			definition.NewEnum("FurtherRemovedEnum", []string{"OneEnum", "TwoEnum"}),
			definition.NewObject("ToBeRemoved", nil),
			definition.NewObject("Basic", nil),
			definition.NewObject("X", nil),
			definition.NewObject("Y", nil),
			definition.NewObject("A", []property.Property{property.NewObjectReference("a", "#/definitions/B")}),
			definition.NewObject("B", []property.Property{property.NewObjectReference("a", "#/definitions/C")}),
			definition.NewObject("C", []property.Property{
				property.NewObjectReference("a", "#/definitions/FurtherRemovedEnum"),
				property.NewArray("test", property.NewObjectReference("a", "#/definitions/D")),
			}),
			definition.NewObject("D", nil),
			definition.NewObject("RecurseL", []property.Property{property.NewObjectReference("a", "#/definitions/RecurseR")}),
			definition.NewObject("RecurseR", []property.Property{property.NewObjectReference("a", "#/definitions/RecurseL")}),
			definition.NewObject("RecurseM", []property.Property{property.NewObjectReference("a", "#/definitions/RecurseM")}),
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
							"200": path.NewResponse("", definition.NewObjectReference("#/definitions/A")),
							"400": path.NewResponse("", definition.NewObjectReference("#/definitions/SomeEnum")),
							"500": path.NewResponse("", definition.NewObjectReference("#/definitions/RecurseM")),
							"501": path.NewResponse("", definition.NewObjectReference("#/definitions/RecurseL")),
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
							"200": path.NewResponse("", definition.NewObjectReference("#/definitions/Basic")),
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
	err := app.Run([]string{"swag3d", "--file", "swagger.json", "generate"})

	// ********************************* ASSERT *******************************
	assert.NoError(t, err)
	assert.Equal(t, "", errOut.String())
	assert.Equal(t, `// This code was generated by: 
// https://github.com/recolude/swagger-unity-codegen
// Issues and PRs welcome :)

using UnityEngine;
using UnityEngine.Networking;
using System.Collections;

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
	err := app.Run([]string{"swag3d", "--file", "swagger.json", "generate", "--namespace", "example"})

	// ********************************* ASSERT *******************************
	assert.NoError(t, err)
	assert.Equal(t, "", errOut.String())
	assert.Equal(t, `// This code was generated by: 
// https://github.com/recolude/swagger-unity-codegen
// Issues and PRs welcome :)

using UnityEngine;
using UnityEngine.Networking;
using System.Collections;

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
	err := app.Run([]string{"swag3d", "--file", "swagger.json", "generate", "--out", "."})

	// ********************************* ASSERT *******************************
	assert.NoError(t, err)
	AssertFileExists(t, appFS, "Definitions.cs")
	AssertFileExists(t, appFS, "Services.cs")
	AssertFileExists(t, appFS, "ServiceConfig.cs")
}
