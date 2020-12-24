package template

import "testing"

func TestApplyStruct_Public(t *testing.T) {
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

func TestApplyStruct_Private(t *testing.T) {
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

func TestApplyFunc_Reciver(t *testing.T) {
	template := &CodeTemplate{}

	expected := `
	func (r *Test) TestFunc(cat string, toast string) (string,error) {
		// do some code
	}
	`

	inputs := make(map[string]string)
	inputs["cat"] = "string"
	inputs["toast"] = "string"

	output := []string{"string", "error"}

	template.ApplyFunc("TestFunc", inputs, output, "*Test", "// do some code")

	if expected != template.Content {
		t.Errorf("\n expected: \n %s got: \n %s", expected, template.Content)
	}
}

func TestApplyFunc_NoReciver(t *testing.T) {
	template := &CodeTemplate{}

	expected := `
	func TestFunc(cat string, toast string) (string,error) {
		// do some code
	}
	`

	inputs := make(map[string]string)
	inputs["cat"] = "string"
	inputs["toast"] = "string"

	output := []string{"string", "error"}

	template.ApplyFunc("TestFunc", inputs, output, "", "// do some code")

	if expected != template.Content {
		t.Errorf("\n expected: \n %s got: \n %s", expected, template.Content)
	}
}
