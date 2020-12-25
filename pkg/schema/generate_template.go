package schema

import template "github.com/GuyARoss/windtunnel/pkg/golang-template"

type SchemaTemplate struct {
	StageCodePaths map[string]string
}

// Generate generates a new code template for the schema
func (s *SchemaTemplate) Generate(schemaParser *ParserResponse) (map[string]*template.CodeTemplate, error) {
	// @@ validate that the pre, main & post input/ outputs line up
	templates := make(map[string]*template.CodeTemplate)

	// @lazy load in structs, apply the ones that will be needed for the stage

	for stageName, stageFields := range schemaParser.Stages {
		// @@ validate that each stage has a code path, if not throw error

		stageGenerationResponse, stageGenerationErr := generateStage(stageName, stageFields.Properties)
		if stageGenerationErr != nil {
			return templates, stageGenerationErr
		}

		templates[stageName] = stageGenerationResponse
	}

	return templates, nil
}

func generateStage(stageName string, stageProperties map[string]string) (*template.CodeTemplate, error) {
	stageCode := &template.CodeTemplate{}
	err := stageCode.ApplyStruct(stageName, make(map[string]string, 0), template.PrivateAccess)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

/* @@todo
type Stage1In struct{}

func (r *Stage1In) validate() error {
	// ensure that the fields are correct
}

type Stage1Out struct{}

func (r *Stage1Out) validate() error {

}


type Stage struct {}

func stageManager(data []byte) (*Stage1Out, error) {
	// marshal data from input
	// call stage1 pre, stage, post
}

func (s *Stage1) stage() (*Stage1In, error) {
	// ~validate input
	// call the builtin running the content
	// validate output + marshal it
}

func (s *Stage1) prestage() (*Stage1In, error) {
	// ~ validate input
	// call the builtin running the content
	// ~ validate output + marshal it
}

func (s *Stage1) poststage() (*Stage1Out, error) {
	// ~validate input
	// call the builtin running the content
	// ~validate output is what we are expecting + marshal it
}
*/

// high level
// grpc in request
// call stages
// 	pre, validate
//  main, validate
// 	post, validate
// grpc out request

/*
	for k, v := range schemaParser.Definitions {
		applyErr := s.codeTemplate.ApplyStruct(k, v.Properties, template.PublicAccess)
		if applyErr != nil {
			return applyErr
		}
		validateOutput := []string{"error"}
		s.codeTemplate.ApplyFunc("validate", make(map[string]string), validateOutput, k, `
		return nil
		`)
	}

*/
