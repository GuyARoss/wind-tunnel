package template

import (
	"testing"
)

func TestApplyStruct_Public(t *testing.T) {
	template := &CodeTemplateCtx{
		Structs: make(map[string]*StructTemplate),
	}

	properties := make(map[string]string, 3)
	properties["Property1"] = "String"

	template.ApplyStruct("test123", properties, PublicAccess)

	expected := "string"
	got := template.Structs["Test123"].Properties["Property1"]
	if got != expected {
		t.Errorf("expected %s got %s", expected, got)
	}
}

func TestApplyStruct_Private(t *testing.T) {
	template := &CodeTemplateCtx{
		Structs: make(map[string]*StructTemplate),
	}

	properties := make(map[string]string, 3)
	properties["Property1"] = "String"

	template.ApplyStruct("test123", properties, PrivateAccess)

	expected := "string"
	got := template.Structs["test123"].Properties["Property1"]
	if got != expected {
		t.Errorf("expected %s got %s", expected, got)
	}
}

func TestApplyFunc_Reciver(t *testing.T) {
	template := &CodeTemplateCtx{
		Funcs: make(map[string]*FuncTemplate),
	}

	inputs := make(map[string]string)
	inputs["cat"] = "string"

	output := []string{"string", "error"}

	template.ApplyFunc("TestFunc", inputs, output, "*Test", "// do some code")

	if template.Funcs["TestFunc"].name != "TestFunc" {
		t.Error("error building test func")
	}
}

func TestApplyFunc_NoReciver(t *testing.T) {
	template := &CodeTemplateCtx{
		Funcs: make(map[string]*FuncTemplate),
	}

	inputs := make(map[string]string)
	inputs["cat"] = "string"

	output := []string{"string", "error"}

	template.ApplyFunc("TestFunc", inputs, output, "", "// do some code")

	if template.Funcs["TestFunc"].name != "TestFunc" {
		t.Error("error building test func")
	}
}
