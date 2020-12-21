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
}

type supportedPropertyType string

const (
	stringPropertyType supportedPropertyType = "String"
	intPropertyType    supportedPropertyType = "Int"
)

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
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		err := ctx.parseLine(scanner.Bytes())
		if err != nil {
			return nil, err
		}
	}

	schemaError := validateDefinitionStageMatch(ctx)
	return &ctx.parserResponse, schemaError
}

func validateDefinitionStageMatch(ctx *lineCtx) error {
	// ensure that all definitions that are defined within stages exist
	for stage, stageValues := range ctx.parserResponse.stages {
		for _, definitionName := range stageValues.properties {
			if ctx.parserResponse.definitions[definitionName] == nil {
				return fmt.Errorf("definition (%s) missing for stage %s", definitionName, stage)
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

	strLine := string(line)
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

	scopePartionErr := validateScopeParition(linePartitions)

	if ctx.scope == noScopeType {
		noScopeErr := validateNoScope(linePartitions)
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
		definitionScopeErr := validateDefinitionScope(propertyType)
		if definitionScopeErr != nil {
			return definitionScopeErr
		}
	}

	ctx.prevLineType = propertyLineType
	ctx.transientScopeProperties[propertyName] = propertyType

	return nil
}

func validateDefinitionScope(propertyType string) error {
	validProperties := []supportedPropertyType{stringPropertyType, intPropertyType}
	if slice.Contains(validProperties, propertyType) {
		return nil
	}

	return fmt.Errorf("Unsupported property %s", propertyType)
}

func validateNoScope(linePartitions []string) error {
	newScope := linePartitions[0] // @validate that it is not empty

	validScopes := []scopeType{schemaDefinitionScopeType, schemaStageScopeType}

	for _, validScope := range validScopes {
		// attempt to find a valid scope
		if string(validScope) == newScope {
			return nil
		}
	}

	return errors.New("invalid scope type")
}

func validateScopeParition(linePartitions []string) error {
	// @@ validate property
	if len(linePartitions) == 2 {
		return nil
	}

	return errors.New("invalid scope line partition size")
}
