package schema

import template "github.com/GuyARoss/windtunnel/pkg/golang-template"

type SchemaTemplate struct {
	StageCodePaths map[string]string
}

// Generate generates a new code template for the schema
func (s *SchemaTemplate) Generate(schemaParser *ParserResponse) (map[string]*template.CodeTemplateCtx, error) {
	// @@ validate that the pre, main & post input/ outputs line up
	templates := make(map[string]*template.CodeTemplateCtx)

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

func generateStage(stageName string, stageProperties map[string]string) (*template.CodeTemplateCtx, error) {
	stageCode := &template.CodeTemplateCtx{}
	err := stageCode.ApplyStruct(stageName, make(map[string]string, 0), template.PrivateAccess)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

/* @@todo
type Stage1In struct{}

func (r *Stage1In) validate() error {}

type Stage1Out struct{}

func (r *Stage1Out) validate() error {}


type Stage struct {
	codeFile
}

func stageManager(data []byte) (*Stage1Out, error) {
	// marshal data from input
	// call stage
	// return stage
}

func (s *Stage1) stage() (*Stage1In, error) {
	// ~validate input
	// call the builtin running the content
	// validate output + marshal it
}

*/

// high level
// grpc in request
// call stage
//  main, validate
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
