package template

import "testing"

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
