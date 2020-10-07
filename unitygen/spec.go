package unitygen

import (
	"fmt"
	"sort"
	"strings"

	"github.com/recolude/swagger-unity-codegen/unitygen/convention"
	"github.com/recolude/swagger-unity-codegen/unitygen/definition"
	"github.com/recolude/swagger-unity-codegen/unitygen/security"
)

// Spec is the overall interpretted swagger file
type Spec struct {
	Info            SpecInfo
	Definitions     []definition.Definition
	AuthDefinitions []security.Auth
	Services        []Service
}

func NewSpec(info SpecInfo, definitions []definition.Definition, authDefinitions []security.Auth, services []Service) Spec {
	sort.Sort(sortByDefinitionName(definitions))
	sort.Sort(sortBySecurityIdentifier(authDefinitions))
	return Spec{
		Info:            info,
		Definitions:     definitions,
		AuthDefinitions: authDefinitions,
		Services:        services,
	}
}

type sortByDefinitionName []definition.Definition

func (a sortByDefinitionName) Len() int           { return len(a) }
func (a sortByDefinitionName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortByDefinitionName) Less(i, j int) bool { return a[i].Name() < a[j].Name() }

type sortBySecurityIdentifier []security.Auth

func (a sortBySecurityIdentifier) Len() int           { return len(a) }
func (a sortBySecurityIdentifier) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortBySecurityIdentifier) Less(i, j int) bool { return a[i].Identifier() < a[j].Identifier() }

// SpecInfo is general info about the spec itself
type SpecInfo struct {
	Title   string
	Version string
}

func (s Spec) basePathDescription() string {
	return "The base URL to which the endpoint paths are appended"
}

func (s Spec) renderInterfaceBody() string {
	builder := strings.Builder{}

	builder.WriteString("\t// ")
	builder.WriteString(s.basePathDescription())
	builder.WriteString("\n")
	builder.WriteString("\tstring BasePath { get; }\n\n")
	for _, authGuard := range s.AuthDefinitions {
		fmt.Fprintf(&builder, "\t// %s\n", authGuard.String())
		fmt.Fprintf(&builder, "\tstring %s { get; }\n\n", convention.TitleCase(authGuard.Identifier()))
	}
	return builder.String()
}

func (s Spec) renderScriptableObjectBody() string {
	builder := strings.Builder{}

	builder.WriteString("\t[SerializeField]\n")
	builder.WriteString("\tprivate string basePath;\n\n")

	builder.WriteString("\t// ")
	builder.WriteString(s.basePathDescription())
	builder.WriteString("\n")
	builder.WriteString("\tpublic string BasePath { get { return basePath; } set { basePath = value; } }\n\n")
	for _, authGuard := range s.AuthDefinitions {

		privateVarName := convention.CamelCase(authGuard.Identifier())

		builder.WriteString("\t[SerializeField]\n")
		fmt.Fprintf(&builder, "\tprivate string %s;\n\n", privateVarName)

		fmt.Fprintf(&builder, "\t// %s\n", authGuard.String())
		fmt.Fprintf(&builder, "\tpublic string %s { get { return %s; } set { %s = value; } }\n\n", convention.TitleCase(authGuard.Identifier()), privateVarName, privateVarName)
	}
	return builder.String()
}

// ServiceConfig prints out a c# class with variables to be used for all requests
func (s Spec) ServiceConfig(configName, menuName string, includeScriptableObject bool) string {
	properClassName := convention.TitleCase(configName)
	builder := strings.Builder{}

	// Interface
	builder.WriteString("public interface Config {\n\n")
	builder.WriteString(s.renderInterfaceBody())
	builder.WriteString("}")

	// Editor Config Code
	if includeScriptableObject {
		builder.WriteString("\n\n#if UNITY_EDITOR\n[UnityEditor.CustomEditor(typeof(")
		builder.WriteString(properClassName)
		builder.WriteString("))]\npublic class ")
		builder.WriteString(properClassName)
		builder.WriteString("Editor : UnityEditor.Editor\n{\n\n\tpublic override void OnInspectorGUI()\n\t{\n")
		builder.WriteString("\t\tif (target == null)\n\t\t{\n\t\t\treturn;\n\t\t}\n\n\t\tvar castedTarget = (")
		builder.WriteString(properClassName)
		builder.WriteString(")target;\n\n\t\tUnityEditor.EditorGUILayout.Space();\n")
		fmt.Fprintf(&builder, "\t\tUnityEditor.EditorGUILayout.LabelField(\"%s\");\n", s.basePathDescription())
		fmt.Fprintf(&builder, "\t\tvar newBasePath = UnityEditor.EditorGUILayout.TextField(\"BasePath\", castedTarget.BasePath);\n")
		fmt.Fprintf(&builder, "\t\tif (newBasePath != castedTarget.BasePath) {\n")
		fmt.Fprintf(&builder, "\t\t\tcastedTarget.BasePath = newBasePath;\n")
		builder.WriteString("\t\t\tUnityEditor.EditorUtility.SetDirty(target);\n\t\t}\n\n")
		for _, authGuard := range s.AuthDefinitions {
			propertyName := convention.TitleCase(authGuard.Identifier())
			privateVarName := "new" + propertyName
			fmt.Fprintf(&builder, "\t\tUnityEditor.EditorGUILayout.Space();\n")
			fmt.Fprintf(&builder, "\t\tUnityEditor.EditorGUILayout.LabelField(\"%s\");\n", authGuard.String())
			fmt.Fprintf(&builder, "\t\tvar %s = UnityEditor.EditorGUILayout.TextField(\"%s\", castedTarget.%s);\n", privateVarName, propertyName, propertyName)
			fmt.Fprintf(&builder, "\t\tif (%s != castedTarget.%s) {\n", privateVarName, propertyName)
			fmt.Fprintf(&builder, "\t\t\tcastedTarget.%s = %s;\n", propertyName, privateVarName)
			builder.WriteString("\t\t\tUnityEditor.EditorUtility.SetDirty(target);\n\t\t}\n\n")
		}
		builder.WriteString("\t}\n\n}\n#endif\n\n")

		// Scriptable Object
		builder.WriteString("[System.Serializable]\n")
		fmt.Fprintf(&builder, "[CreateAssetMenu(menuName = \"%s\", fileName = \"%s\")]\n", menuName, properClassName)
		fmt.Fprintf(&builder, "public class %s: ScriptableObject, Config {\n\n", properClassName)
		builder.WriteString(s.renderScriptableObjectBody())

		// Scriptable objects can't have constructors, need to decide on alternative method of instantiating them in code
		// fmt.Fprintf(&builder, "\tpublic %s(string basePath) {\n", properClassName)
		// builder.WriteString("\t\tthis.BasePath = basePath;\n")
		// builder.WriteString("\t}\n")

		builder.WriteString("}")
	}

	return builder.String()
}
