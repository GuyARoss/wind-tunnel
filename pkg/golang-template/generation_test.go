package template

import "testing"

func TestGenerateStruct(t *testing.T) {
	temp := &GeneratedTemplate{Content: ""}
	expected := `
type Test123 struct {
	Property1 string
}
	`

	temp.generateStruct("Test123", map[string]string{
		"Property1": "string",
	})

	if temp.Content != expected {
		t.Errorf("expected %s got %s", expected, temp.Content)
	}
}

func TestGenerateFunc_Reciver(t *testing.T) {
	temp := &GeneratedTemplate{Content: ""}
	expected := `
func (r *Test) TestFunc(cat string) (string, error) {
	// do some code
}
	`

	temp.generateFunc(&FuncTemplate{
		name:         "TestFunc",
		body:         "// do some code",
		receiverType: "*Test",
		inputs: map[string]string{
			"cat": "string",
		},
		outputs: []string{"string", "error"},
	})

	if expected != temp.Content {
		t.Errorf("expected %s got %s", expected, temp.Content)
	}
}

func TestGenerateFunc_NoReciver(t *testing.T) {
	temp := &GeneratedTemplate{Content: ""}
	expected := `
func TestFunc(cat string) (string, error) {
	// do some code
}
	`

	temp.generateFunc(&FuncTemplate{
		name:         "TestFunc",
		body:         "// do some code",
		receiverType: "",
		inputs: map[string]string{
			"cat": "string",
		},
		outputs: []string{"string", "error"},
	})

	if expected != temp.Content {
		t.Errorf("expected %s got %s", expected, temp.Content)
	}
}

func TestGenerateImports(t *testing.T) {
	temp := &GeneratedTemplate{Content: ""}
	expected := `
import (
	"thing"
)
	`

	imports := make(map[string]string)
	imports["roast"] = "thing"

	temp.generateImports(imports)

	if expected != temp.Content {
		t.Errorf("expected %s got %s", expected, temp.Content)
	}
}
