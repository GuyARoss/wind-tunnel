package templates

import "testing"

func TestNewStruct_Public(t *testing.T) {
	expected := `
	struct Test123 {
		Property1 String
		Property2 String
		Property3 String
	}
	`

	properties := make(map[string]string, 3)

	got := newStruct("test123", properties, publicAccess)
	if expected != got {
		t.Errorf("got %s expected %s", got, expected)
	}
}

func TestNewStruct_Private(t *testing.T) {
	expected := `
	struct test123 {
		Property1 string
		Property2 string
		Property3 string
	}
	`

	properties := make(map[string]string, 3)

	got := newStruct("test123", properties, privateAccess)
	if expected != got {
		t.Errorf("got %s expected %s", got, expected)
	}
}
