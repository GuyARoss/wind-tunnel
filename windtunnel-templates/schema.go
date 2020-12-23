package templates

import (
	"github.com/GuyARoss/windtunnel/pkg/schema"
)

type SchemaTemplate struct {
}

func (template *SchemaTemplate) Generate(schemaParser *schema.ParserResponse) string {
	output := make([]string, len(schemaParser.Definitions)+len(schemaParser.Stages))

	for k, v := range schemaParser.Definitions {
		output[len(output)-1] = newStruct(k, v.Properties, publicAccess)
	}

	for k, v := range schemaParser.Stages {
		output[len(output)-1] = newStruct(k, v.Properties, publicAccess)
	}

	return ``
}
