package model_test

import (
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen/model"
	"github.com/recolude/swagger-unity-codegen/unitygen/model/property"
	"github.com/stretchr/testify/assert"
)

func TestObject(t *testing.T) {
	// ******************************** ARRANGE *******************************
	prop := property.NewInteger("num", "")
	obj := model.NewObject(
		"testObj",
		[]model.Property{prop},
	)

	// ********************************** ACT *********************************
	varType := obj.ToVariableType()
	name := obj.Name()
	cSharp := obj.ToCSharp()
	converter := obj.JsonConverter()
	properties := obj.Properties()

	// ********************************* ASSERT *******************************
	if assert.Len(t, properties, 1) {
		assert.Equal(t, prop, properties[0])
	}
	assert.Equal(t, "", converter)
	assert.Equal(t, "TestObj", varType)
	assert.Equal(t, "testObj", name)
	assert.Equal(t, `[System.Serializable]
public class TestObj {

	[JsonProperty("num")]
	public int Num { get; private set; }

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

	[JsonProperty("date")]
	public string date;

	public System.DateTime Date { get => System.DateTime.Parse(date); }

}`, cSharp)
}

func TestAllOfObject(t *testing.T) {
	// ******************************** ARRANGE *******************************
	p1 := property.NewInteger("num", "")
	p2 := property.NewString("date", "date-time")
	p3 := property.NewBoolean("anotha one")
	p4 := property.NewString("and anotha one", "")

	obj := model.NewObject(
		"testObj",
		[]model.Property{p1, p2},
	)

	allOf := model.NewAllOfObject(
		"CompositionThingy",
		obj,
		[]model.Property{p3, p4},
	)

	// ********************************** ACT *********************************
	varType := allOf.ToVariableType()
	name := allOf.Name()
	converter := allOf.JsonConverter()
	cSharp := allOf.ToCSharp()
	properties := allOf.Properties()

	// ********************************* ASSERT *******************************
	if assert.Len(t, properties, 4) {
		assert.Equal(t, p1, properties[1])
		assert.Equal(t, p2, properties[0])
		assert.Equal(t, p3, properties[3])
		assert.Equal(t, p4, properties[2])
	}
	assert.Equal(t, "", converter)
	assert.Equal(t, "CompositionThingy", varType)
	assert.Equal(t, "CompositionThingy", name)
	assert.Equal(t, `[System.Serializable]
public class CompositionThingy {

	[JsonProperty("date")]
	public string date;

	public System.DateTime Date { get => System.DateTime.Parse(date); }

	[JsonProperty("num")]
	public int Num { get; private set; }

	[JsonProperty("and anotha one")]
	public string AndAnothaOne { get; private set; }

	[JsonProperty("anotha one")]
	public bool AnothaOne { get; private set; }

}`, cSharp)
}

func Test_CanSetAllOfObject(t *testing.T) {
	// ******************************** ARRANGE *******************************
	p1 := property.NewInteger("num", "")
	p2 := property.NewString("date", "date-time")
	p3 := property.NewBoolean("anotha one")
	p4 := property.NewString("and anotha one", "")

	obj := model.NewObject(
		"testObj",
		[]model.Property{p1, p2},
	)

	allOf := model.NewObject(
		"CompositionThingy",
		[]model.Property{p3, p4},
	)
	allOf.SetAllOfObject(&obj)

	// ********************************** ACT *********************************
	varType := allOf.ToVariableType()
	name := allOf.Name()
	converter := allOf.JsonConverter()
	cSharp := allOf.ToCSharp()
	properties := allOf.Properties()

	// ********************************* ASSERT *******************************
	if assert.Len(t, properties, 4) {
		assert.Equal(t, p1, properties[1])
		assert.Equal(t, p2, properties[0])
		assert.Equal(t, p3, properties[3])
		assert.Equal(t, p4, properties[2])
	}
	assert.Equal(t, "", converter)
	assert.Equal(t, "CompositionThingy", varType)
	assert.Equal(t, "CompositionThingy", name)
	assert.Equal(t, `[System.Serializable]
public class CompositionThingy {

	[JsonProperty("date")]
	public string date;

	public System.DateTime Date { get => System.DateTime.Parse(date); }

	[JsonProperty("num")]
	public int Num { get; private set; }

	[JsonProperty("and anotha one")]
	public string AndAnothaOne { get; private set; }

	[JsonProperty("anotha one")]
	public bool AnothaOne { get; private set; }

}`, cSharp)
}

func TestObject_CanSetChildAndParent(t *testing.T) {
	// ******************************** ARRANGE *******************************
	p0 := property.NewString("disc", "")
	p1 := property.NewInteger("num", "")
	p2 := property.NewString("date", "date-time")
	p3 := property.NewBoolean("anotha one")
	p4 := property.NewString("and anotha one", "")

	parentObj := model.NewDiscriminatorObject(
		"Parent",
		[]model.Property{p0, p1, p2},
		"disc",
	)

	childObj := model.NewObject(
		"Child",
		[]model.Property{p3, p4},
	)
	childObj.SetWhatToInherit(&parentObj)
	parentObj.AddChild(&childObj)

	// ********************************** ACT *********************************
	varType := childObj.ToVariableType()
	name := childObj.Name()
	converter := childObj.JsonConverter()
	childCSharp := childObj.ToCSharp()
	properties := childObj.Properties()

	parentCSharp := parentObj.ToCSharp()

	// ********************************* ASSERT *******************************
	if assert.Len(t, properties, 2) {
		assert.Equal(t, p3, properties[1])
		assert.Equal(t, p4, properties[0])
	}
	assert.Equal(t, "", converter)
	assert.Equal(t, "Child", varType)
	assert.Equal(t, "Child", name)
	assert.Equal(t, `[System.Serializable]
public class Child : Parent {

	[JsonProperty("and anotha one")]
	public string AndAnothaOne { get; private set; }

	[JsonProperty("anotha one")]
	public bool AnothaOne { get; private set; }

}`, childCSharp)

	assert.Equal(t, `[System.Serializable]
[JsonConverter(typeof(JsonSubtypes), "disc")]
[JsonSubtypes.KnownSubType(typeof(Child), "Child")]
public class Parent {

	[JsonProperty("date")]
	public string date;

	public System.DateTime Date { get => System.DateTime.Parse(date); }

	[JsonProperty("disc")]
	public string Disc { get; private set; }

	[JsonProperty("num")]
	public int Num { get; private set; }

}`, parentCSharp)
}

func TestObject_ParentPanicsWithNilChild(t *testing.T) {
	// ******************************** ARRANGE *******************************
	p0 := property.NewString("disc", "")

	parentObj := model.NewDiscriminatorObject(
		"Parent",
		[]model.Property{p0},
		"disc",
	)

	parentObj.AddChild(nil)

	// ******************************* ACT/ASSERT *****************************
	assert.PanicsWithError(t, "Parent was elected as parent class and provided a nil child", func() {
		parentObj.ToCSharp()
	})
}
