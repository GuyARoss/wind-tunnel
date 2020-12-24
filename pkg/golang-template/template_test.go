package template

import "testing"

func TestNewStruct_Public(t *testing.T) {
	template := &CodeTemplate{}

	expected := `
	type Test123 struct {
		Property1 string
	}
	`

	properties := make(map[string]string, 3)
	properties["Property1"] = "String"

	template.ApplyStruct("test123", properties, PublicAccess)
	if expected != template.Content {
		t.Errorf("got %s expected %s", template.Content, expected)
	}
}

func TestNewStruct_Private(t *testing.T) {
	template := &CodeTemplate{}

	expected := `
	type test123 struct {
		Property1 string
	}
	`

	properties := make(map[string]string, 3)
	properties["Property1"] = "String"

	template.ApplyStruct("test123", properties, PrivateAccess)
	if expected != string(template.Content) {
		t.Errorf("got %s expected %s", template.Content, expected)
	}
}
