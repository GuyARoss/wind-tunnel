package schema

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/stretchr/stew/slice"
)

type parseSymbol string

const (
	eofParseSymbol parseSymbol = "}"
	sofParseSymbol parseSymbol = "{"
)

type lineType int32

const (
	emptyLineType lineType = 0
	schemaLineType
	propertyLineType
	symbolLineType
)

type scopeType string

const (
	noScopeType               scopeType = ""
	schemaDefinitionScopeType scopeType = "definition"
	schemaStageScopeType      scopeType = "stage"
)

type SchemaScope struct {
	Name       string
	Properties map[string]string
}

type ParserResponse struct {
	Definitions map[string]*SchemaScope
	Stages      map[string]*SchemaScope
}

type lineCtx struct {
	prevLineType             lineType
	scope                    scopeType
	scopeID                  string
	transientScopeProperties map[string]string
	parserResponse           *ParserResponse
	lineIndex                int
	fileName                 string
}

type supportedPropertyType string

const (
	stringPropertyType supportedPropertyType = "String"
	intPropertyType    supportedPropertyType = "Int"
)

func (ctx *lineCtx) createLineErr(errorStr string) error {
	updatedMsg := fmt.Sprintf("schema compilation error: '%s' at line %d of %s", errorStr, ctx.lineIndex, ctx.fileName)
	return errors.New(updatedMsg)
}

// ParseFile uses prexisting parser response to apply to schema file
func (r *ParserResponse) ParseFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	ctx := &lineCtx{
		prevLineType:             emptyLineType,
		scope:                    noScopeType,
		transientScopeProperties: make(map[string]string),
		parserResponse:           r,
		lineIndex:                0,
		fileName:                 path,
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		err := ctx.parseLine(scanner.Bytes())
		if err != nil {
			return err
		}

		ctx.lineIndex++
	}

	return nil
}

// ParseFile parse file objects to the schema profile
func ParseFile(path string) (*ParserResponse, error) {
	resp := &ParserResponse{
		Definitions: make(map[string]*SchemaScope),
		Stages:      make(map[string]*SchemaScope),
	}

	return resp, resp.ParseFile(path)
}

// ValidateDefinitionStageMatch validates that definitions can be linked to stages
func (r *ParserResponse) ValidateDefinitionStageMatch() error {
	// ensure that all definitions that are defined within stages exist
	for stage, stageValues := range r.Stages {
		fmt.Println(stage, stageValues.Properties)
		for _, definitionName := range stageValues.Properties {
			if definitionName != "None" && r.Definitions[definitionName] == nil {
				return fmt.Errorf("definition `%s` missing for stage `%s`", definitionName, stage)
			}
		}
	}
	return nil
}

func (ctx *lineCtx) parseLine(line []byte) error {
	if len(line) == 0 {
		// empty line detected
		ctx.prevLineType = emptyLineType
		ctx.scope = noScopeType
		ctx.transientScopeProperties = make(map[string]string)

		return nil
	}

	strLine := parseLineIndent(string(line))
	linePartitions := strings.Split(strLine, " ")

	if linePartitions[0] == string(eofParseSymbol) {
		completedScope := &SchemaScope{
			Name:       ctx.scopeID,
			Properties: ctx.transientScopeProperties,
		}

		switch ctx.scope {
		case schemaDefinitionScopeType:
			{
				ctx.parserResponse.Definitions[ctx.scopeID] = completedScope
				break
			}
		case schemaStageScopeType:
			{
				ctx.parserResponse.Stages[ctx.scopeID] = completedScope
				break
			}
		}

		ctx.scopeID = ""
		ctx.scope = noScopeType
		ctx.transientScopeProperties = make(map[string]string)
		ctx.prevLineType = symbolLineType

		return nil
	}

	scopePartionErr := ctx.validateScopeParition(linePartitions)

	if ctx.scope == noScopeType {
		noScopeErr := ctx.validateNoScope(linePartitions)
		if noScopeErr != nil {
			return noScopeErr
		}

		ctx.scope = scopeType(linePartitions[0])
		ctx.scopeID = linePartitions[1] // @@ validate dis
		ctx.prevLineType = schemaLineType
		return nil
	}

	if scopePartionErr != nil {
		return scopePartionErr
	}

	propertyName := linePartitions[0]
	propertyType := linePartitions[1]

	if ctx.scope == schemaDefinitionScopeType {
		definitionScopeErr := ctx.validateDefinitionScope(propertyType)
		if definitionScopeErr != nil {
			return definitionScopeErr
		}
	}

	ctx.prevLineType = propertyLineType
	ctx.transientScopeProperties[propertyName] = propertyType

	return nil
}

func (ctx *lineCtx) validateDefinitionScope(propertyType string) error {
	validProperties := []supportedPropertyType{stringPropertyType, intPropertyType}
	if slice.Contains(validProperties, propertyType) {
		return nil
	}

	return ctx.createLineErr(fmt.Sprintf("unsupported property `%s`", propertyType))
}

func (ctx *lineCtx) validateStageProperties() error {
	requiredFields := []string{"prestage", "in", "out", "poststage"}

	for stage, stageValues := range ctx.parserResponse.Stages {
		countFound := 0
		for propertyKey := range stageValues.Properties {
			if slice.Contains(requiredFields, propertyKey) {
				countFound++
			}
		}

		if countFound != len(requiredFields) {
			return fmt.Errorf("stage `%s` missing required stages", stage)
		}
	}

	return nil
}

func (ctx *lineCtx) validateNoScope(linePartitions []string) error {
	newScope := linePartitions[0] // @validate that it is not empty

	validScopes := []scopeType{schemaDefinitionScopeType, schemaStageScopeType}

	for _, validScope := range validScopes {
		// attempt to find a valid scope
		if string(validScope) == newScope {
			return nil
		}
	}

	return ctx.createLineErr(fmt.Sprintf("invalid scope type `%s`", newScope))
}

func (ctx *lineCtx) validateScopeParition(linePartitions []string) error {
	// @@ validate property
	if len(linePartitions) == 2 {
		return nil
	}

	return ctx.createLineErr("invalid line formatting")
}

func parseLineIndent(input string) string {
	stop := 0
	for idx, char := range input {
		if char != ' ' {
			stop = idx
			break
		}
	}

	return input[stop:]
}
