package template

import (
	"testing"
)

func TestApplyStruct_Public(t *testing.T) {
	template := &CodeTemplateCtx{
		structs: make(map[string][]string),
	}

	properties := make(map[string]string, 3)
	properties["Property1"] = "String"

	template.ApplyStruct("test123", properties, PublicAccess)

	expected := "	Property1 string"
	got := template.structs["Test123"][0]
	if got != expected {
		t.Errorf("expected %s got %s", expected, got)
	}
}

func TestApplyStruct_Private(t *testing.T) {
	template := &CodeTemplateCtx{
		structs: make(map[string][]string),
	}

	properties := make(map[string]string, 3)
	properties["Property1"] = "String"

	template.ApplyStruct("test123", properties, PrivateAccess)

	expected := "	Property1 string"
	got := template.structs["test123"][0]
	if got != expected {
		t.Errorf("expected %s got %s", expected, got)
	}
}

func TestApplyFunc_Reciver(t *testing.T) {
	template := &CodeTemplateCtx{
		funcs: make(map[string]*funcTemplate),
	}

	inputs := make(map[string]string)
	inputs["cat"] = "string"

	output := []string{"string", "error"}

	template.ApplyFunc("TestFunc", inputs, output, "*Test", "// do some code")

	if template.funcs["TestFunc"].name != "TestFunc" {
		t.Error("error building test func")
	}
}

func TestApplyFunc_NoReciver(t *testing.T) {
	template := &CodeTemplateCtx{
		funcs: make(map[string]*funcTemplate),
	}

	inputs := make(map[string]string)
	inputs["cat"] = "string"

	output := []string{"string", "error"}

	template.ApplyFunc("TestFunc", inputs, output, "", "// do some code")

	if template.funcs["TestFunc"].name != "TestFunc" {
		t.Error("error building test func")
	}
}
