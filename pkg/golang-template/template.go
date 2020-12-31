package template

import (
	"fmt"
	"strings"

	"github.com/GuyARoss/windtunnel/pkg/utilities"
	"github.com/stretchr/stew/slice"
)

type accessModification string

const (
	// PrivateAccess access modifier of private
	PrivateAccess accessModification = "private"

	// PublicAccess access modifier of public
	PublicAccess accessModification = "public"
)

func (access accessModification) formatToAccessType(value string) string {
	switch access {
	case PrivateAccess:
		return strings.ToLower(value)
	case PublicAccess:
		return strings.Title(value)
	default:
		// @@ invalid
		return ""
	}
}

func newStructKey(key string, value string, access accessModification) string {
	if access == PublicAccess {
		key = strings.Title(key)
	}

	return key
}

type codeChars string

const (
	endCodeBlockChar     codeChars = "}"
	endImportBlockChar   codeChars = ")"
	startImportBlockChar codeChars = "("
)

type primitiveFieldTypes string

const (
	stringPrimitive primitiveFieldTypes = "string"
	intPrimitive    primitiveFieldTypes = "int"
	boolPrimitive   primitiveFieldTypes = "bool"
)

func isPrimitiveType(field string) bool {
	primitives := []primitiveFieldTypes{stringPrimitive, intPrimitive}
	return slice.Contains(primitives, strings.ToLower(field))
}

type FuncTemplate struct {
	name         string
	body         string
	receiverType string
	inputs       map[string]string
	outputs      []string
}

func CreateFuncTemplate(name string, inputs map[string]string, outputs []string, receiver string, body string) (*FuncTemplate, error) {
	// seralizedInputs := make([]string, 0)
	// for k, v := range inputs {
	// 	// @@todo validate that the values exist in scope
	// 	seralizedInputs = append(seralizedInputs, fmt.Sprintf("%s %s", k, v))
	// }

	return &FuncTemplate{
		name:         name,
		body:         body,
		receiverType: receiver,
		inputs:       inputs,
		outputs:      outputs,
	}, nil
}

func CreateStructTemplate(name string, properties map[string]string, access accessModification) (*StructTemplate, error) {
	modifiedProperties := make(map[string]string)

	for propertyKey, propertyValue := range properties {
		value := propertyValue
		if isPrimitiveType(value) {
			value = strings.ToLower(propertyValue)
		}
		structKey := newStructKey(propertyKey, value, access)
		modifiedProperties[structKey] = value
	}

	name = access.formatToAccessType(name)
	return &StructTemplate{
		Name:       name,
		Properties: modifiedProperties,
		Access:     access,
		Funcs:      make(map[string]*FuncTemplate),
	}, nil
}

type StructTemplate struct {
	Name       string
	Properties map[string]string
	Access     accessModification
	Funcs      map[string]*FuncTemplate
}

func (r *StructTemplate) ApplyFunc(name string, inputs map[string]string, output []string, body string) error {
	reciever := fmt.Sprintf("*%s", r.Name)

	temp, err := CreateFuncTemplate(name, inputs, output, reciever, body)
	if err != nil {
		return err
	}

	r.Funcs[name] = temp
	return nil
}

// CodeTemplateCtx holds the template context throughout the lifecycle
type CodeTemplateCtx struct {
	Structs  map[string]*StructTemplate
	Funcs    map[string]*FuncTemplate
	Imports  map[string]string
	Builtins map[string]string

	BuiltinDir string
}

// ApplyStruct creates a new struct within the code template
func (t *CodeTemplateCtx) ApplyStruct(name string, properties map[string]string, access accessModification) (*string, error) {
	if t.Structs[name] != nil {
		// @@todo: raise already exists error
	}

	newStruct, err := CreateStructTemplate(name, properties, access)
	if err != nil {
		return nil, err
	}

	t.Structs[newStruct.Name] = newStruct
	return &newStruct.Name, nil
}

// ApplyFunc creates a new func within the code template
// note: body is not validated
func (t *CodeTemplateCtx) ApplyFunc(name string, inputs map[string]string, output []string, receiver string, body string) error {
	fnTemplate, err := CreateFuncTemplate(name, inputs, output, receiver, body)
	if err != nil {
		return err
	}

	t.Funcs[name] = fnTemplate

	return nil
}

type BuiltinRequirementType string

const (
	StructBuiltinRequirement BuiltinRequirementType = "struct"
	FuncBuiltinRequirement   BuiltinRequirementType = "func"
)

type BuiltinRequirement struct {
	Name string
	Type BuiltinRequirementType
}

// ApplyBuiltin applies builtin + their dependencies to the code template
func (t *CodeTemplateCtx) ApplyBuiltin(
	builtinsDir string,
	requirements []*BuiltinRequirement,
	changeMap map[string]string,
) error {
	// @@ we could run part of this process once to load all of the builtins-
	// @@ then again to apply them to the template
	files := utilities.FindFiles(builtinsDir, ".go")
	fmt.Println(files)

	bctx := &builtinCtx{
		requiredDependencies: requirements,
		scope:                nonScopeType,
		sourceMap:            make(map[string]string),
		imports:              make(map[string]string),
	}

	for _, file := range files {
		bctx.loadBuiltinFile(file)
	}

	fmt.Printf("imports: %s \n", bctx.imports)
	// apply source maps + validate imports
	for k, v := range bctx.imports {
		t.Imports[k] = v
	}

	for k, v := range bctx.sourceMap {
		// @@todo: apply changeMap
		t.Builtins[k] = v
	}

	return nil
}

// Generate performs the code generation process for the code template
func (t *CodeTemplateCtx) Generate() (*GeneratedTemplate, error) {
	gtemp := &GeneratedTemplate{
		Content: "",
	}
	gtemp.generateModuleName("main")
	gtemp.generateImports(t.Imports)

	for structKey, data := range t.Structs {
		gtemp.generateStruct(structKey, data.Properties)

		for _, fn := range data.Funcs {
			gtemp.generateFunc(fn)
		}
	}

	for _, funcTemplate := range t.Funcs {
		gtemp.generateFunc(funcTemplate)
	}

	fmt.Println(t.Builtins)
	for _, bn := range t.Builtins {
		gtemp.append(bn)
	}

	return gtemp, nil
}
