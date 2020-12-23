package templates

import (
	"fmt"

	"github.com/stretchr/stew/slice"
)

type accessModification string

const (
	privateAccess accessModification = "private"
	publicAccess  accessModification = "public"
)

func newStructProperty(key string, value string) string {
	return fmt.Sprintf(`
		%s %s
	`, key, value)
}

type primitiveFieldTypes string

const (
	stringPrimitive primitiveFieldTypes = "string"
	intPrimitive    primitiveFieldTypes = "int"
)

func isPrimitiveField(field string) bool {
	primitives := []primitiveFieldTypes{stringPrimitive, intPrimitive}
	return slice.Contains(field, primitives)
}

func newStruct(name string, properties map[string]string, access accessModification) string {
	structProperies := make([]string, len(properties))
}
