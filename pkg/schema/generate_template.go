package schema

import template "github.com/GuyARoss/windtunnel/pkg/golang-template"

type SchemaTemplate struct {
	StageCodePaths map[string]string
	codeTemplate   *template.CodeTemplate
}

func (s *SchemaTemplate) Generate(schemaParser *ParserResponse) error {
	for k, v := range schemaParser.Definitions {
		applyErr := s.codeTemplate.ApplyStruct(k, v.Properties, template.PublicAccess)
		if applyErr != nil {
			return applyErr
		}
		validateOutput := []string{"error"}
		s.codeTemplate.ApplyFunc("validate", make(map[string]string), validateOutput, propertyValue)	
	}

	// for k, v := range schemaParser.Stages {
	// 	output[len(output)-1] = newStruct(k, v.Properties, publicAccess)
	// }

	return nil
}

func (s *SchemaTemplate) generateStage(stageName string, stageProperties map[string]string) error {
	err := s.codeTemplate.ApplyStruct(stageName, make(map[string]string, 0), template.PrivateAccess)

	if err != nil {
		return err
	}
	
	return nil
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
