package schema

import (
	"bufio"
	"errors"
	"fmt"
	"io"
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
	name       string
	properties map[string]string
}

type ParserResponse struct {
	definitions map[string]*SchemaScope
	stages      map[string]*SchemaScope
}

type lineCtx struct {
	prevLineType             lineType
	scope                    scopeType
	scopeID                  string
	transientScopeProperties map[string]string
	parserResponse           ParserResponse
	lineIndex                int
}

type supportedPropertyType string

const (
	stringPropertyType supportedPropertyType = "String"
	intPropertyType    supportedPropertyType = "Int"
)

func (ctx *lineCtx) createLineErr(errorStr string) error {
	updatedMsg := fmt.Sprintf("schema compilation error: '%s' near our around line %d", errorStr, ctx.lineIndex)
	return errors.New(updatedMsg)
}

// ParseFile parse file objects to the schema profile
func ParseFile(file io.Reader) (*ParserResponse, error) {
	ctx := &lineCtx{
		prevLineType:             emptyLineType,
		scope:                    noScopeType,
		transientScopeProperties: make(map[string]string),
		parserResponse: ParserResponse{
			definitions: make(map[string]*SchemaScope),
			stages:      make(map[string]*SchemaScope),
		},
		lineIndex: 0,
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		err := ctx.parseLine(scanner.Bytes())
		if err != nil {
			return nil, err
		}

		ctx.lineIndex++
	}

	schemaError := ctx.validateDefinitionStageMatch()
	return &ctx.parserResponse, schemaError
}

func (ctx *lineCtx) validateDefinitionStageMatch() error {
	// ensure that all definitions that are defined within stages exist
	for stage, stageValues := range ctx.parserResponse.stages {
		for _, definitionName := range stageValues.properties {
			if ctx.parserResponse.definitions[definitionName] == nil {
				return ctx.createLineErr(
					fmt.Sprintf("definition (%s) missing for stage %s", definitionName, stage),
				)
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
	// strLine := string(line)
	linePartitions := strings.Split(strLine, " ")

	if linePartitions[0] == string(eofParseSymbol) {
		completedScope := &SchemaScope{
			name:       ctx.scopeID,
			properties: ctx.transientScopeProperties,
		}

		switch ctx.scope {
		case schemaDefinitionScopeType:
			{
				ctx.parserResponse.definitions[ctx.scopeID] = completedScope
				break
			}
		case schemaStageScopeType:
			{
				ctx.parserResponse.stages[ctx.scopeID] = completedScope
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

	return ctx.createLineErr(fmt.Sprintf("Unsupported property %s", propertyType))
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

	return ctx.createLineErr(fmt.Sprintf("invalid scope type '%s'", newScope))
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
