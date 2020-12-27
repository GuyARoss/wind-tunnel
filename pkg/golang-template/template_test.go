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
	func (r *Test) TestFunc(cat string) (string,error) {
		// do some code
	}
	`

	inputs := make(map[string]string)
	inputs["cat"] = "string"

	output := []string{"string", "error"}

	template.ApplyFunc("TestFunc", inputs, output, "*Test", "// do some code")

	if expected != template.Content {
		t.Errorf("\n expected: \n %s got: \n %s", expected, template.Content)
	}
}

func TestApplyFunc_NoReciver(t *testing.T) {
	template := &CodeTemplate{}

	expected := `
	func TestFunc(cat string) (string,error) {
		// do some code
	}
	`

	inputs := make(map[string]string)
	inputs["cat"] = "string"

	output := []string{"string", "error"}

	template.ApplyFunc("TestFunc", inputs, output, "", "// do some code")

	if expected != template.Content {
		t.Errorf("\n expected: \n %s got: \n %s", expected, template.Content)
	}
}

func TestParseBuiltinLine_EOD(t *testing.T) {
	bctx := &builtinCtx{
		scope:     "test",
		sourceMap: make(map[string]string),
	}

	err := bctx.parseBuiltinLine([]byte("}"))
	if err != nil {
		t.Error("eod parse should not throw err")
	}

	if bctx.sourceMap["test"] != "}" {
		t.Errorf("expected } got %s", bctx.sourceMap["test"])
	}
}

func TestParseBuiltinLine_IsDependecy(t *testing.T) {
	bctx := &builtinCtx{
		scope:                "",
		sourceMap:            make(map[string]string),
		requiredDependencies: []string{"myFunc"},
	}

	lineContent := "func myFunc(lst []string) error"

	err := bctx.parseBuiltinLine([]byte(lineContent))
	if err != nil {
		t.Error("eod parse should not throw err")
	}

	if bctx.sourceMap["myFunc"] != lineContent {
		t.Errorf("line content does not match, expected %s got %s", lineContent, bctx.sourceMap["myFunc"])
	}
}

func TestLinearStrContains(t *testing.T) {
	in := "doesthiswork"
	resp := linearStrContains(in, "this")

	if resp != true {
		t.Error("expected to assert true")
	}
}

func TestParseImportLine(t *testing.T) {
	expected := &importLine{
		name: "test",
		path: "github.com/GuyARoss/test",
	}

	out := parseImportLine("\"github.com/GuyARoss/test\"")

	if out.name != expected.name && out.path != expected.path {
		t.Errorf("expected %s got %s", expected, out)
	}
}
