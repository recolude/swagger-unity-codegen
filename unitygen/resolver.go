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
	var allOfObjIndex int

	for i, def := range definitions {
		if def.Name() == aof.objToChangeName {
			objToChangeIndex = i
			switch v := def.(type) {
			case model.Object:
				objToChange = &v
			case *model.Object:
				objToChange = v
			}
		}

		if def.Name() == aof.allOfObjName {
			allOfObjIndex = i

			switch v := def.(type) {
			case model.Object:
				allOfObj = &v
			case *model.Object:
				allOfObj = v
			default:
				return fmt.Errorf("definition `%s` references non-object `%s` in 'allOf'", aof.objToChangeName, aof.allOfObjName)
			}
		}
	}

	if objToChange == nil {
		return errors.New("could not find object to resolve for. This is most likely a bug with the generator and not a problem with the swagger file. Please open up a bug report and supply your json.")
	}

	if allOfObj == nil {
		return fmt.Errorf("definition `%s` that `%s` references in 'allOf' was not found in the swagger file", aof.allOfObjName, aof.objToChangeName)
	}

	if allOfObj.HasDiscriminator() {
		objToChange.SetWhatToInherit(allOfObj)
		allOfObj.AddChild(objToChange)
		definitions[allOfObjIndex] = *allOfObj
	} else {
		objToChange.SetAllOfObject(allOfObj)
	}
	definitions[objToChangeIndex] = *objToChange
	return nil
}
