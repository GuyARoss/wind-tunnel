package schema

import template "github.com/GuyARoss/windtunnel/pkg/golang-template"

// GenerationSettings initial settings for template generation
type GenerationSettings struct {
	BuiltinsDir string
}

// GenerationCtx context for stage generation
type GenerationCtx struct {
	BaseStructs map[string]*template.StructTemplate
	Settings    *GenerationSettings
}

func (s *GenerationCtx) generateStage(stageName string, stageProperties map[string]string) (*template.CodeTemplateCtx, error) {
	stageCode := &template.CodeTemplateCtx{}

	in := stageProperties["in"]
	out := stageProperties["out"]
	// @@todo: validate in & out

	builtinErr := stageCode.ApplyBuiltin(s.Settings.BuiltinsDir,
		[]string{"xxx_Pipein", "xxx_Pipeout", "pipeEnv", "writeToPipe", "readFromPipe", "startUpStage"},
		map[string]string{
			"xxx_Pipein":  in,
			"xxx_Pipeout": out,
		},
	)

	if builtinErr != nil {
		return nil, builtinErr
	}

	err := stageCode.ApplyStruct(stageName, map[string]string{
		"codeFile": "string",
	}, template.PrivateAccess)

	// @@ could use a builtin for dis..
	stageCode.Structs[stageName].ApplyFunc("invoke", map[string]string{"input": in}, []string{out}, `
	// @@todo: fill in 
	// - validate input from "input" param
	// - write data to pipe
	// - read data from pipe
	// - marshal pipe data to "output" type
	// - validate output
	// - return output

	return nil
	`)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Generate loads in parsed schema context and generates client code from it
func (s *ParserResponse) Generate(settings *GenerationSettings) (map[string]*template.GeneratedTemplate, error) {
	generationCtx := &GenerationCtx{
		BaseStructs: make(map[string]*template.StructTemplate),
		Settings:    settings,
	}

	for defKey, def := range s.Definitions {
		structTemplate := &template.StructTemplate{
			Name:       defKey,
			Properties: def.Properties,
			Access:     template.PublicAccess,
		}
		structTemplate.ApplyFunc("validate", make(map[string]string), []string{"error"}, `
		//@@todo: fill in
		return nil
		`)
	}

	templates := make(map[string]*template.GeneratedTemplate)

	// @@todo: apply only the structures that are needed for the stage (this may need to be done during code generation)
	// @@performance: make each of these stages into go routines
	for stageName, stageFields := range s.Stages {
		// @@todo: validate that each stage has a code path, if not throw error

		stageGenerationResponse, stageGenerationErr := generationCtx.generateStage(stageName, stageFields.Properties)
		if stageGenerationErr != nil {
			return templates, stageGenerationErr
		}

		codeGeneration, codeGenerationErrr := stageGenerationResponse.Generate()
		if codeGenerationErrr != nil {
			return templates, codeGenerationErrr
		}

		templates[stageName] = codeGeneration
	}

	return templates, nil
}
