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
