package schema

import (
	"testing"
)

func TestParseLine_EmptyLine(t *testing.T) {
	lctx := &lineCtx{}

	line := make([]byte, 0)
	err := lctx.parseLine(line)
	if err != nil {
		t.Error("parse line should not return error when noScope present")
	}
	if lctx.scope != noScopeType {
		t.Errorf("expected scope to be noScope, got %s",
			lctx.scope,
		)
	}

	if lctx.prevLineType != emptyLineType {
		t.Errorf("expected prevLineType to be emptyLineType, got %d",
			lctx.prevLineType,
		)
	}
}

func TestParseLine_IndentLine(t *testing.T) {
	parserResponse := ParserResponse{
		definitions: make(map[string]*SchemaScope),
		stages:      make(map[string]*SchemaScope),
	}
	lctx := &lineCtx{
		scopeID:                  "testScope",
		scope:                    schemaDefinitionScopeType,
		transientScopeProperties: make(map[string]string),
		prevLineType:             propertyLineType,
		parserResponse:           parserResponse,
	}

	err := lctx.parseLine([]byte("    name String"))
	if err != nil {
		t.Error("parse line should handle indent")
	}
}

func TestParseLine_EOFCharStageScope(t *testing.T) {
	properties := make(map[string]string)
	properties["TestProperty"] = "String"

	parserResponse := ParserResponse{
		definitions: make(map[string]*SchemaScope),
		stages:      make(map[string]*SchemaScope),
	}

	lctx := &lineCtx{
		scopeID:                  "testScope",
		scope:                    schemaStageScopeType,
		transientScopeProperties: properties,
		prevLineType:             propertyLineType,
		parserResponse:           parserResponse,
	}

	line := []byte("}")
	err := lctx.parseLine(line)
	if err != nil {
		t.Error("EOFChar should not return an error")
	}

	stage := lctx.parserResponse.stages["testScope"]
	if stage == nil || stage.name != "testScope" {
		t.Error("error inserting stage")
	}

	if stage.properties == nil || stage.properties["TestProperty"] != "String" {
		t.Error("error inserting properties into stage")
	}

	if lctx.scope != noScopeType {
		t.Error("error resetting scope type")
	}

	if lctx.prevLineType != symbolLineType {
		t.Error("error updating line type")
	}
}

func TestParseLine_EOFCharDefinitionScope(t *testing.T) {
	properties := make(map[string]string)
	properties["TestProperty"] = "String"

	parserResponse := ParserResponse{
		definitions: make(map[string]*SchemaScope),
		stages:      make(map[string]*SchemaScope),
	}

	lctx := &lineCtx{
		scopeID:                  "testScope",
		scope:                    schemaDefinitionScopeType,
		transientScopeProperties: properties,
		prevLineType:             propertyLineType,
		parserResponse:           parserResponse,
	}

	line := []byte("}")
	err := lctx.parseLine(line)
	if err != nil {
		t.Error("EOFChar should not return an error")
	}

	stage := lctx.parserResponse.definitions["testScope"]
	if stage == nil || stage.name != "testScope" {
		t.Error("error inserting definition")
	}

	if stage.properties == nil || stage.properties["TestProperty"] != "String" {
		t.Error("error inserting properties into definition")
	}

	if lctx.scope != noScopeType {
		t.Error("error resetting scope type")
	}

	if lctx.prevLineType != symbolLineType {
		t.Error("error updating line type")
	}
}

func TestParseLine_NoScopeDefinition(t *testing.T) {
	lctx := &lineCtx{
		scope: noScopeType,
	}
	err := lctx.parseLine([]byte("definition Test"))
	if err != nil {
		t.Error("no scope definition should not return an error")
	}
	if lctx.scope != schemaDefinitionScopeType {
		t.Errorf("expected schemaDefinitionScopeType got %s", lctx.scope)
	}
}

func TestParseLine_ExistingDefinitionScope(t *testing.T) {
	lctx := &lineCtx{
		scope:                    schemaDefinitionScopeType,
		transientScopeProperties: make(map[string]string),
	}
	err := lctx.parseLine([]byte("fileName String"))
	if err != nil {
		t.Error("add property existing scope should not return an error")
	}

	if lctx.transientScopeProperties["fileName"] != "String" {
		t.Errorf("transientScopeProperties expected String got %s", lctx.transientScopeProperties["fileName"])
	}
}

func TestParseLineIndent(t *testing.T) {
	got := parseLineIndent("    cat 123")
	expect := "cat 123"

	if got != expect {
		t.Errorf("parse line indent expected %s got %s", got, expect)
	}
}
