package templates

import (
	"fmt"
	"strings"

	"github.com/stretchr/stew/slice"
)

type accessModification string

const (
	privateAccess accessModification = "private"
	publicAccess  accessModification = "public"
)

func (access accessModification) formatToAccessType(value string) string {
	switch access {
	case privateAccess:
		return strings.ToLower(value)
	case publicAccess:
		return strings.Title(value)
	default:
		// @@ invalid
		return ""
	}
}

func newStructProperty(key string, value string, access accessModification) string {
	if access == publicAccess {
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

func newStruct(name string, properties map[string]string, access accessModification) string {
	structProperies := make([]string, 0)
	for propertyKey, propertyValue := range properties {
		value := propertyValue
		if isPrimitiveType(value) {
			value = strings.ToLower(propertyValue)
		}
		structProperies = append(structProperies, newStructProperty(propertyKey, value, access))
	}

	name = access.formatToAccessType(name)

	return fmt.Sprintf(`
	type %s struct {
	%s
	}
	`, name, strings.Join(structProperies, "\n	"))
}
