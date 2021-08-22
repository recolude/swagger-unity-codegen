package unitygen

import (
	"errors"
	"fmt"

	"github.com/recolude/swagger-unity-codegen/unitygen/model"
)

// resolvers get called after the definitions have had their initial
// passthrough. At this point, all definitions have been recognized, and
// connections between them can be resovled.
type resolver interface {
	Resolve(definitions []model.Definition) error
}

type allOfResolver struct {
	objToChangeName string
	allOfObjName    string
}

func (aof allOfResolver) Resolve(definitions []model.Definition) error {
	var objToChange *model.Object
	var objToChangeIndex int
	var allOfObj *model.Object

	for i, def := range definitions {
		if def.Name() == aof.objToChangeName {
			objToChangeIndex = i
			casted, ok := def.(model.Object)
			if ok {
				objToChange = &casted
			}
		}

		if def.Name() == aof.allOfObjName {
			casted, ok := def.(model.Object)
			if ok {
				allOfObj = &casted
			}
		}
	}

	if objToChange == nil {
		return errors.New("could not find object to resolve for")
	}

	if allOfObj == nil {
		return fmt.Errorf("definition `%s` that `%s` referenced in 'allOf' was not found in the swagger file", aof.allOfObjName, aof.objToChangeName)
	}

	objToChange.SetAllOfObject(allOfObj)
	definitions[objToChangeIndex] = objToChange
	return nil
}
