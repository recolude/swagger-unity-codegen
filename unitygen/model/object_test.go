package model_test

import (
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen/model"
	"github.com/recolude/swagger-unity-codegen/unitygen/model/property"
	"github.com/stretchr/testify/assert"
)

func TestObject(t *testing.T) {
	// ******************************** ARRANGE *******************************
	obj := model.NewObject(
		"testObj",
		[]model.Property{
			property.NewInteger("num", ""),
		},
	)

	// ********************************** ACT *********************************
	varType := obj.ToVariableType()
	name := obj.Name()
	cSharp := obj.ToCSharp()

	// ********************************* ASSERT *******************************
	assert.Equal(t, "TestObj", varType)
	assert.Equal(t, "testObj", name)
	assert.Equal(t, `[System.Serializable]
public class TestObj {

	public int num;

}`, cSharp)
}

func TestObject_DatesCorrectly(t *testing.T) {
	// ******************************** ARRANGE *******************************
	obj := model.NewObject(
		"testObj",
		[]model.Property{
			property.NewString("date", "date-time"),
		},
	)

	// ********************************** ACT *********************************
	varType := obj.ToVariableType()
	name := obj.Name()
	cSharp := obj.ToCSharp()

	// ********************************* ASSERT *******************************
	assert.Equal(t, "TestObj", varType)
	assert.Equal(t, "testObj", name)
	assert.Equal(t, `[System.Serializable]
public class TestObj {

	[SerializeField]
	private string date;

	public System.DateTime Date { get => System.DateTime.Parse(date); }

}`, cSharp)
}
