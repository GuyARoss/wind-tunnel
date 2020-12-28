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

func newStructProperty(key string, value string, access accessModification) string {
	if access == PublicAccess {
		key = strings.Title(key)
	}

	return fmt.Sprintf("	%s %s", key, value)
}

type codeChars string

const (
	endCodeBlockChar   codeChars = "}"
	endImportBlockChar codeChars = ")"
)

type primitiveFieldTypes string

const (
	stringPrimitive primitiveFieldTypes = "string"
	intPrimitive    primitiveFieldTypes = "int"
)

func isPrimitiveType(field string) bool {
	primitives := []primitiveFieldTypes{stringPrimitive, intPrimitive}
	return slice.Contains(primitives, strings.ToLower(field))
}

type funcTemplate struct {
	name             string
	body             string
	receiverType     string
	seralizedInputs  []string
	seralizedOutputs string
}

func createFuncTemplate(name string, inputs map[string]string, output []string, receiver string, body string) (*funcTemplate, error) {
	seralizedInputs := make([]string, 0)
	for k, v := range inputs {
		// @@todo validate that the values exist in scope
		seralizedInputs = append(seralizedInputs, fmt.Sprintf("%s %s", k, v))
	}

	return &funcTemplate{
		name:             name,
		body:             body,
		receiverType:     receiver,
		seralizedInputs:  seralizedInputs,
		seralizedOutputs: strings.Join(output, ","),
	}, nil
}

type structTemplate struct {
	name       string
	properties map[string]string
	access     accessModification
	funcs      map[string]*funcTemplate
}

func (r *structTemplate) applyFunc(name string, inputs map[string]string, output []string, body string) error {
	reciever := fmt.Sprintf("*%s", r.name)

	temp, err := createFuncTemplate(name, inputs, output, reciever, body)
	if err != nil {
		return err
	}

	r.funcs[name] = temp
	return nil
}

// CodeTemplateCtx holds the template context throughout the lifecycle
type CodeTemplateCtx struct {
	structs  map[string]*structTemplate
	funcs    map[string]*funcTemplate
	imports  map[string]string
	builtins map[string]string

	BuiltinDir string
}

// ApplyStruct creates a new struct within the code template
func (t *CodeTemplateCtx) ApplyStruct(name string, properties map[string]string, access accessModification) error {
	if t.structs[name] != nil {
		// @@todo: raise already exists error
	}

	for propertyKey, propertyValue := range properties {
		value := propertyValue
		if isPrimitiveType(value) {
			value = strings.ToLower(propertyValue)
		}
		properties[propertyKey] = newStructProperty(propertyKey, value, access)
	}

	name = access.formatToAccessType(name)
	t.structs[name] = &structTemplate{
		name:       name,
		properties: properties,
		access:     access,
	}

	return nil
}

// ApplyFunc creates a new func within the code template
// note: body is not validated
func (t *CodeTemplateCtx) ApplyFunc(name string, inputs map[string]string, output []string, receiver string, body string) error {
	fnTemplate, err := createFuncTemplate(name, inputs, output, receiver, body)
	if err != nil {
		return err
	}

	t.funcs[name] = fnTemplate

	return nil
}

// ApplyBuiltin applies builtin + their dependencies to the code template
func (t *CodeTemplateCtx) ApplyBuiltin(
	builtinsDir string,
	requiredDependencies []string,
	changeMap map[string]string,
) error {
	files := utilities.FindFiles(builtinsDir, ".go")

	bctx := &builtinCtx{
		requiredDependencies: requiredDependencies,
		sourceMap:            make(map[string]string),
		scope:                nonScopeType,
		imports:              make(map[string]string),
	}

	for _, file := range files {
		bctx.loadBuiltinFile(file)
	}

	// apply source maps + validate imports
	for k, v := range bctx.imports {
		t.imports[k] = v
	}

	for k, v := range bctx.sourceMap {
		t.builtins[k] = v
	}

	return nil
}

// Generate performs the code generation process for the code template
func (t *CodeTemplateCtx) Generate() (*GeneratedTemplate, error) {
	gtemp := &GeneratedTemplate{
		Content: "",
	}

	for structKey, data := range t.structs {
		gtemp.generateStruct(structKey, data.properties)
	}

	for _, funcTemplate := range t.funcs {
		gtemp.generateFunc(funcTemplate)
	}

	gtemp.generateImports(t.imports)

	return gtemp, nil
}
