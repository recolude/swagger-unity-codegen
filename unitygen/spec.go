package unitygen

import (
	"sort"

	"github.com/recolude/swagger-unity-codegen/unitygen/definition"
	"github.com/recolude/swagger-unity-codegen/unitygen/path"
	"github.com/recolude/swagger-unity-codegen/unitygen/security"
)

// Spec is the overall interpretted swagger file
type Spec struct {
	Info            SpecInfo
	Definitions     []definition.Definition
	AuthDefinitions []security.Auth
	Paths           []path.Path
}

func NewSpec(info SpecInfo, definitions []definition.Definition, authDefinitions []security.Auth, paths []path.Path) Spec {
	sort.Sort(sortByDefinitionName(definitions))
	return Spec{
		Info:            info,
		Definitions:     definitions,
		AuthDefinitions: authDefinitions,
		Paths:           paths,
	}
}

type sortByDefinitionName []definition.Definition

func (a sortByDefinitionName) Len() int           { return len(a) }
func (a sortByDefinitionName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortByDefinitionName) Less(i, j int) bool { return a[i].Name() < a[j].Name() }

// SpecInfo is general info about the spec itself
type SpecInfo struct {
	Title   string
	Version string
}
