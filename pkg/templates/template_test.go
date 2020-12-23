package templates

import "testing"

func TestNewStruct_Public(t *testing.T) {
	expected := `
	type Test123 struct {
		Property1 string
	}
	`

	properties := make(map[string]string, 3)
	properties["Property1"] = "String"

	got := newStruct("test123", properties, publicAccess)
	if expected != got {
		t.Errorf("got %s expected %s", got, expected)
	}
}

func TestNewStruct_Private(t *testing.T) {
	expected := `
	type test123 struct {
		Property1 string
	}
	`

	properties := make(map[string]string, 3)
	properties["Property1"] = "String"

	got := newStruct("test123", properties, privateAccess)
	if expected != got {
		t.Errorf("got %s expected %s", got, expected)
	}
}
