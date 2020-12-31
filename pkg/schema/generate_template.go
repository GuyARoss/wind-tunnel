package schema

import (
	"fmt"

	template "github.com/GuyARoss/windtunnel/pkg/golang-template"
)

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
	stageCode := &template.CodeTemplateCtx{
		Structs:    s.BaseStructs,
		Funcs:      make(map[string]*template.FuncTemplate),
		Imports:    make(map[string]string),
		Builtins:   make(map[string]string),
		BuiltinDir: s.Settings.BuiltinsDir,
	}

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

	updatedStageName, err := stageCode.ApplyStruct(stageName, map[string]string{
		"codeFile": "string",
	}, template.PrivateAccess)

	if err != nil {
		return nil, err
	}

	pointerIn := fmt.Sprintf("*%s", in)
	// @@fyi could use a builtin for dis..
	stageCode.Structs[*updatedStageName].ApplyFunc("invoke", map[string]string{"input": pointerIn}, []string{out}, `
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

	return stageCode, nil
}

// Generate loads in parsed schema context and generates client code from it
func (s *ParserResponse) Generate(settings *GenerationSettings) (map[string]*template.GeneratedTemplate, error) {
	generationCtx := &GenerationCtx{
		BaseStructs: make(map[string]*template.StructTemplate),
		Settings:    settings,
	}

	for defKey, def := range s.Definitions {
		structTemplate, structCreateErr := template.CreateStructTemplate(defKey, def.Properties, template.PublicAccess)
		if structCreateErr != nil {
			return nil, structCreateErr
		}
		structTemplate.ApplyFunc("validate", make(map[string]string), []string{"error"}, `
//@@todo: fill in
return nil
		`)

		generationCtx.BaseStructs[defKey] = structTemplate
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
