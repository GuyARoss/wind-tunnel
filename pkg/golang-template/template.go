package template

import (
	"fmt"
	"strings"

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

type primitiveFieldTypes string

const (
	stringPrimitive primitiveFieldTypes = "string"
	intPrimitive    primitiveFieldTypes = "int"
)

func isPrimitiveType(field string) bool {
	primitives := []primitiveFieldTypes{stringPrimitive, intPrimitive}
	return slice.Contains(primitives, strings.ToLower(field))
}

// CodeTemplate holds the template context throughout the lifecycle
type CodeTemplate struct {
	structs map[string]map[string]string

	Content string
}

func (t *CodeTemplate) append(data string) {
	if len(t.Content) == 0 {
		t.Content += fmt.Sprintf("%s", data)
		return
	}

	t.Content += fmt.Sprintf("%s \n", data)
}

// ApplyStruct creates a new struct within the code template
func (t *CodeTemplate) ApplyStruct(name string, properties map[string]string, access accessModification) error {
	structProperties := make([]string, 0)
	for propertyKey, propertyValue := range properties {
		value := propertyValue
		if isPrimitiveType(value) {
			value = strings.ToLower(propertyValue)
		}
		structProperties = append(structProperties, newStructProperty(propertyKey, value, access))
	}

	name = access.formatToAccessType(name)

	// @@todo validate that the struct doesn't already exist
	t.append(fmt.Sprintf(`
	type %s struct {
	%s
	}
	`, name, strings.Join(structProperties, "\n	")))

	return nil
}

// ApplyFunc creates a new func within the code template
// note: body is not validated
func (t *CodeTemplate) ApplyFunc(name string, inputs map[string]string, output []string, receiver string, body string) error {
	seralizedInputs := make([]string, 0)
	for k, v := range inputs {
		// @@todo validate that the values exist in scope
		seralizedInputs = append(seralizedInputs, fmt.Sprintf("%s %s", k, v))
	}

	if len(receiver) > 0 {
		t.append(fmt.Sprintf(`
	func (r %s) %s(%s) (%s) {
		%s
	}
	`, receiver, name, strings.Join(seralizedInputs, ", "), strings.Join(output, ","), body))
		return nil
	}

	t.append(fmt.Sprintf(`
	func %s(%s) (%s) {
		%s
	}
	`, name, strings.Join(seralizedInputs, ", "), strings.Join(output, ","), body))

	return nil
}
