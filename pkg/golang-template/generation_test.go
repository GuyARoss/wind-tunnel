package template

import "testing"

func TestGenerateStruct(t *testing.T) {
	temp := &GeneratedTemplate{Content: ""}
	expected := `
	type Test123 struct {
		Property1 string
	}
	`

	temp.generateStruct("Test123", []string{
		"Property1 string",
	})

	if temp.Content != expected {
		t.Errorf("expected %s got %s", expected, temp.Content)
	}
}

func TestGenerateFunc_Reciver(t *testing.T) {
	temp := &GeneratedTemplate{Content: ""}
	expected := `
	func (r *Test) TestFunc(cat string) (string,error) {
		// do some code
	}
	`

	temp.generateFunc(&funcTemplate{
		name:             "TestFunc",
		body:             "// do some code",
		receiverType:     "*Test",
		seralizedInputs:  []string{"cat string"},
		seralizedOutputs: "string,error",
	})

	if expected != temp.Content {
		t.Errorf("expected %s got %s", expected, temp.Content)
	}
}

func TestGenerateFunc_NoReciver(t *testing.T) {
	temp := &GeneratedTemplate{Content: ""}
	expected := `
	func TestFunc(cat string) (string,error) {
		// do some code
	}
	`

	temp.generateFunc(&funcTemplate{
		name:             "TestFunc",
		body:             "// do some code",
		receiverType:     "",
		seralizedInputs:  []string{"cat string"},
		seralizedOutputs: "string,error",
	})

	if expected != temp.Content {
		t.Errorf("expected %s got %s", expected, temp.Content)
	}
}
