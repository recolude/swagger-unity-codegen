package unitygen

import (
	"sort"

	"github.com/recolude/swagger-unity-codegen/unitygen/definition"
)

// Spec is the overall interpretted swagger file
type Spec struct {
	Info        SpecInfo
	Definitions []definition.Definition
}

func NewSpec(info SpecInfo, definitions []definition.Definition) Spec {
	sort.Sort(sortByDefinitionName(definitions))
	return Spec{
		Info:        info,
		Definitions: definitions,
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
