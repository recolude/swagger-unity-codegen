package main

import (
	"github.com/recolude/swagger-unity-codegen/unitygen"
	"github.com/recolude/swagger-unity-codegen/unitygen/model"
	"github.com/recolude/swagger-unity-codegen/unitygen/model/property"
)

func filterSpecForTags(spec unitygen.Spec, tags []string) unitygen.Spec {
	if tags == nil {
		return spec
	}
	filteredServices := make([]unitygen.Service, 0)
	for _, service := range spec.Services {
		for _, tag := range tags {
			if service.Name() == tag {
				filteredServices = append(filteredServices, service)
			}
		}
	}
	spec.Services = filteredServices
	return spec
}

func alreadyRecursed(inquestion string, alreadyAdded []string) bool {
	for _, q := range alreadyAdded {
		if q == inquestion {
			return true
		}
	}
	return false
}

func findReferencePropRecurse(inQuestion model.Property, defs []model.Definition, alreadyAdded []string) []string {
	finalReferences := alreadyAdded

	objectReferenceDefinition, ok := inQuestion.(property.DefinitionReference)
	if ok {
		if alreadyRecursed(objectReferenceDefinition.ToVariableType(), finalReferences) {
			return finalReferences
		}
		for _, def := range defs {
			if def.ToVariableType() == objectReferenceDefinition.ToVariableType() {
				finalReferences = findReferenceRecurse(def, defs, finalReferences)
			}
		}
	}

	arrayDefinition, ok := inQuestion.(property.Array)
	if ok {
		return findReferencePropRecurse(arrayDefinition.Property(), defs, finalReferences)
	}

	objectDefinition, ok := inQuestion.(property.Object)
	if ok {
		finalReferences = findReferenceRecurse(objectDefinition.Object(), defs, finalReferences)
	}

	return finalReferences
}

func findReferenceRecurse(inQuestion model.Definition, defs []model.Definition, alreadyFound []string) []string {
	finalReferences := alreadyFound

	objectDefinition, ok := inQuestion.(model.Object)
	if ok {
		if alreadyRecursed(objectDefinition.ToVariableType(), finalReferences) {
			return finalReferences
		}
		finalReferences = append(finalReferences, objectDefinition.ToVariableType())
		for _, prop := range objectDefinition.Properties() {
			finalReferences = findReferencePropRecurse(prop, defs, finalReferences)
		}
	}

	stringEnumDefinition, ok := inQuestion.(model.StringEnum)
	if ok {
		finalReferences = append(finalReferences, stringEnumDefinition.ToVariableType())
	}

	numEnumDefinition, ok := inQuestion.(model.NumberEnum)
	if ok {
		finalReferences = append(finalReferences, numEnumDefinition.ToVariableType())
	}

	return finalReferences
}

func buildReferenceMapping(defs []model.Definition) map[string][]string {
	refMapping := make(map[string][]string)
	for _, def := range defs {
		refMapping[def.ToVariableType()] = findReferenceRecurse(def, defs, nil)
	}
	return refMapping
}

func filterSpecForUnusedDefinitions(spec unitygen.Spec) unitygen.Spec {
	thingsToKeep := make(map[string]bool)
	for _, def := range spec.Definitions {
		thingsToKeep[def.ToVariableType()] = false
	}

	referenceMapping := buildReferenceMapping(spec.Definitions)

	for _, service := range spec.Services {
		for _, path := range service.Paths() {
			for _, param := range path.Parameters() {
				if param.Schema() != nil {
					if _, ok := thingsToKeep[param.Schema().ToVariableType()]; ok {
						thingsToKeep[param.Schema().ToVariableType()] = true

						for _, reference := range referenceMapping[param.Schema().ToVariableType()] {
							thingsToKeep[reference] = true
						}
					}
				}
			}

			for _, resp := range path.Responses() {
				if resp != nil {
					if _, ok := thingsToKeep[resp.VariableType()]; ok {
						thingsToKeep[resp.VariableType()] = true
						for _, reference := range referenceMapping[resp.VariableType()] {
							thingsToKeep[reference] = true
						}
					}
				}
			}
		}
	}

	filteredDefinitions := make([]model.Definition, 0)
	for _, def := range spec.Definitions {
		if thingsToKeep[def.ToVariableType()] {
			filteredDefinitions = append(filteredDefinitions, def)
		}
	}

	spec.Definitions = filteredDefinitions
	return spec
}
