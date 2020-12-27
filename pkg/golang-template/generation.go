package template

import (
	"fmt"
	"strings"
)

// GeneratedTemplate end product of code generation
type GeneratedTemplate struct {
	Content string
}

func (t *GeneratedTemplate) append(data string) {
	if len(t.Content) == 0 {
		t.Content += fmt.Sprintf("%s", data)
		return
	}

	t.Content += fmt.Sprintf("%s \n", data)
}

func (t *GeneratedTemplate) generateModuleName(name string) {
	t.Content += fmt.Sprintf("package %s", name)
}

func (t *GeneratedTemplate) generateImports(imports map[string]string) {
	values := make([]string, len(imports))
	idx := 0
	for _, v := range imports {
		values[idx] = fmt.Sprintf("\"%s\"", v)
	}

	t.Content += fmt.Sprintf(`
	import (
		%s
	)
	`, strings.Join(values, "\n"))
}

func (t *GeneratedTemplate) generateStruct(name string, structProperties []string) {
	t.append(fmt.Sprintf(`
	type %s struct {
		%s
	}
	`, name, strings.Join(structProperties, "\n	")))
}

func (t *GeneratedTemplate) generateFunc(temp *funcTemplate) error {
	if len(temp.receiverType) > 0 {
		t.append(fmt.Sprintf(`
	func (r %s) %s(%s) (%s) {
		%s
	}
	`, temp.receiverType, temp.name, strings.Join(temp.seralizedInputs, ", "), temp.seralizedOutputs, temp.body))
		return nil
	}

	t.append(fmt.Sprintf(`
	func %s(%s) (%s) {
		%s
	}
	`, temp.name, strings.Join(temp.seralizedInputs, ", "), temp.seralizedOutputs, temp.body))

	return nil
}
