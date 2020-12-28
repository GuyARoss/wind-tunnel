package schema

import template "github.com/GuyARoss/windtunnel/pkg/golang-template"

type GenerationSettings struct {
	StageCodePaths []string
}

type GenerationCtx struct {
	BaseStructs map[string]*template.StructTemplate
	Settings    *GenerationSettings
}

func (s *GenerationCtx) generateStage(stageName string, stageProperties map[string]string) (*template.CodeTemplateCtx, error) {
	stageCode := &template.CodeTemplateCtx{}

	err := stageCode.ApplyStruct(stageName, map[string]string{
		"codeFile": "string"
	}, template.PrivateAccess)

	in := stageProperties["in"]
	out := stageProperties["out"]

	stageCode.Structs[stageName].ApplyFunc("invoke", map[string]string{in}, []string{out}, `
	// @@todo: fill in 
	return nil
	`)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Generate loads in parsed schema context and generates client code from it
func (s *ParserResponse) Generate(settings *GenerationSettings) (map[string]*template.CodeTemplateCtx, error) {
	generationCtx := &GenerationCtx{
		BaseStructs: make(map[string]*template.StructTemplate),
		Settings: settings
	}
	
	for defKey, def := range s.Definitions {
		// @@todo: add the validation functions to each of these
		structTemplate := &template.StructTemplate{
			name: defKey,
			properties: def.Properties,
			access: template.PublicAccess,
		}
		structTemplate.ApplyFunc("validate", make(map[string]string), []string{"error"}, `
		//@@todo: fill in
		return nil
		`)
	}
	
	templates := make(map[string]*template.CodeTemplateCtx)
	
	// @@todo: apply only the structures that are needed for the stage (this may need to be done during code generation)
	// @@performance: make each of these stages into go routines
	for stageName, stageFields := range s.Stages {
		// @@todo: validate that each stage has a code path, if not throw error

		stageGenerationResponse, stageGenerationErr := settings.generateStage(stageName, stageFields.Properties)
		if stageGenerationErr != nil {
			return templates, stageGenerationErr
		}

		templates[stageName] = stageGenerationResponse
	}

	return templates, nil
}
