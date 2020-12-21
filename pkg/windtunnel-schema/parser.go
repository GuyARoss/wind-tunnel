package schema

import (
	"strings"

	"github.com/stretchr/stew/slice"
)

type SchemaParser interface {
	// Parse(f []byte) (ParserResponse, error)
}

type WindTunnelSchemaParser struct{}

func (p *WindTunnelSchemaParser) Parse() (*ParserResponse, error) {
	return nil, nil
}

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
	if !slice.Contains(validProperties, propertyType) {
		// @@ throw unsupported property type
	}

	return nil
}

func validateStageScope() error {
	return nil
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

	// @@ can't find scope throw error dat da scope type is invalid
	return nil
}

func validateScopeParition(linePartitions []string) error {
	// @@ validate property
	if len(linePartitions) != 2 {
		// @@ throw invalid length size err
	}

	return nil
}
