package schema

import template "github.com/GuyARoss/windtunnel/pkg/golang-template"

type SchemaTemplate struct {
}

func (schema *SchemaTemplate) Generate(schemaParser *ParserResponse) (string, error) {
	codeTemplate := &template.CodeTemplate{}

	for k, v := range schemaParser.Definitions {
		applyErr := codeTemplate.ApplyStruct(k, v.Properties, template.PublicAccess)
		if applyErr != nil {
			return "", applyErr
		}
	}

	// for k, v := range schemaParser.Stages {
	// 	output[len(output)-1] = newStruct(k, v.Properties, publicAccess)
	// }

	return "", nil
}

/* @@todo
// @@ do this for all of em
type Stage1 struct {}

func (r *Stage1In) validate() (error) {
	// ensure that the fields are correct
}

func (s *Stage1) x() (Stage1In, error) {
	// call the builtin running the content
	// validate output is what we are expecting + marshal it
}
*/
