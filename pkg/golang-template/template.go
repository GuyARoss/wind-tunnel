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

// CodeTemplateCtx holds the template context throughout the lifecycle
type CodeTemplateCtx struct {
	structs map[string][]string
	funcs   map[string]*funcTemplate

	BuiltinDir string
}

// ApplyStruct creates a new struct within the code template
func (t *CodeTemplateCtx) ApplyStruct(name string, properties map[string]string, access accessModification) error {
	if t.structs[name] != nil {
		// @@todo raise already exists error
	}

	structProperties := make([]string, 0)
	for propertyKey, propertyValue := range properties {
		value := propertyValue
		if isPrimitiveType(value) {
			value = strings.ToLower(propertyValue)
		}
		structProperties = append(structProperties, newStructProperty(propertyKey, value, access))
	}

	name = access.formatToAccessType(name)
	t.structs[name] = structProperties

	return nil
}

// ApplyFunc creates a new func within the code template
// note: body is not validated
func (t *CodeTemplateCtx) ApplyFunc(name string, inputs map[string]string, output []string, receiver string, body string) error {
	seralizedInputs := make([]string, 0)
	for k, v := range inputs {
		// @@todo validate that the values exist in scope
		seralizedInputs = append(seralizedInputs, fmt.Sprintf("%s %s", k, v))
	}

	t.funcs[name] = &funcTemplate{
		name:             name,
		body:             body,
		receiverType:     receiver,
		seralizedInputs:  seralizedInputs,
		seralizedOutputs: strings.Join(output, ","),
	}

	return nil
}

// LoadBuiltin applies builtin + dependencies to the code template
func (ctx *builtinCtx) LoadBuiltin(
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

	// @@ apply source maps + validate imports

	return nil
}
