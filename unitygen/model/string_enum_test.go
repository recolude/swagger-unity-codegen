package model_test

import (
	"testing"

	"github.com/recolude/swagger-unity-codegen/unitygen/model"
	"github.com/stretchr/testify/assert"
)

func TestStringEnum(t *testing.T) {
	// ******************************** ARRANGE *******************************
	enum := model.NewStringEnum(
		"testEnum",
		[]string{
			"A",
			"b",
			"e-c-d",
			"CDF",
		},
	)

	// ********************************** ACT *********************************
	varType := enum.ToVariableType()
	name := enum.Name()
	cSharp := enum.ToCSharp()

	// ********************************* ASSERT *******************************
	assert.Equal(t, "TestEnum", varType)
	assert.Equal(t, "testEnum", name)
	assert.Equal(t, `public enum TestEnum {
	A = 0,
	B = 1,
	ECD = 2,
	CDF = 3
}
public class TestEnumJsonConverter : JsonConverter {
	public override void WriteJson(JsonWriter w, object val, JsonSerializer s) {
		TestEnum castedVal = (TestEnum)val;
		switch (castedVal) {
			case TestEnum.A:
				w.WriteValue("A");
				break;
			case TestEnum.B:
				w.WriteValue("b");
				break;
			case TestEnum.ECD:
				w.WriteValue("e-c-d");
				break;
			case TestEnum.CDF:
				w.WriteValue("CDF");
				break;
			default:
				throw new System.Exception("Unknown value. Living on the dangerous side editing generated code?");
		}
	}

	public override object ReadJson(JsonReader r, System.Type t, object existingValue, JsonSerializer s) {
		var enumString = (string)r.Value;
		switch (enumString) {
			case "A":
				return TestEnum.A;
			case "b":
				return TestEnum.B;
			case "e-c-d":
				return TestEnum.ECD;
			case "CDF":
				return TestEnum.CDF;
			default:
				throw new System.Exception("Unknown value. Perhaps you need to regenerate this code?");
		}
	}

	public override bool CanConvert(System.Type objectType) {
		return objectType == typeof(string);
	}
}`, cSharp)
}
