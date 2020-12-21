package schema

type SchemaDefinition struct {
	name              string
	properties        map[string]string
	fileNameReference string
}

type SchemaStage struct {
	preStageDefinition  SchemaDefinition
	inputDefinition     SchemaDefinition
	outputDefinition    SchemaDefinition
	postStageDefinition SchemaDefinition
}

type ParserResponse struct {
	definitions map[string]SchemaDefinition
	stages      map[string]SchemaStage
}

type SchemaParser interface {
	// Parse(f []byte) (ParserResponse, error)
}

type WindTunnelSchemaParser struct{}

func (p *WindTunnelSchemaParser) Parse() (*ParserResponse, error) {
	return nil, nil
}

type lineType int32

const (
	emptyLineType lineType = 0
	schemaTypeType
)

type scopeType string

const (
	noScopeType               scopeType = ""
	schemaDefinitionScopeType scopeType = "definition"
	schemaStageScopeType      scopeType = "stage"
)

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

		return nil
	}

	// strLine := string(line)

	if ctx.scope == noScopeType {
		// dis line should now indentify scope
	}

	// dis should have properties

	return nil
}
