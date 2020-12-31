package template

import (
	"fmt"
	"os"
	"strings"
)

// GeneratedTemplate end product of code generation
type GeneratedTemplate struct {
	Content string
}

// SaveFile saves the content of the template to a file
func (t *GeneratedTemplate) SaveFile(name string) error {
	// @@todo: use flag path here for output
	f, err := os.Create(fmt.Sprintf("./bin/%s.go", name))
	if err != nil {
		return err
	}

	_, writeErr := f.WriteString(t.Content)
	if writeErr != nil {
		return writeErr
	}

	return f.Close()
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

func (t *GeneratedTemplate) generateImports(imports map[string]string) error {
	if len(imports) == 0 {
		return nil
	}

	values := make([]string, len(imports))
	idx := 0
	for _, v := range imports {
		values[idx] = fmt.Sprintf("\"%s\"", v)
	}

	t.append(fmt.Sprintf(`
import (
	%s
)
	`, strings.Join(values, "\n")))

	return nil
}

func (t *GeneratedTemplate) generateStruct(name string, structMapProperties map[string]string) {
	structProperties := make([]string, 0)
	for k, v := range structMapProperties {
		structProperties = append(structProperties, fmt.Sprintf("%s %s", k, v))
	}

	t.append(fmt.Sprintf(`
type %s struct {
	%s
}
	`, name, strings.Join(structProperties, "\n	")))
}

func (t *GeneratedTemplate) generateFunc(temp *FuncTemplate) error {
	output := fmt.Sprintf("%s", strings.Join(temp.outputs, ", "))
	if len(temp.outputs) > 1 {
		output = fmt.Sprintf("(%s)", output)
	}

	inputs := make([]string, 0)
	for k, v := range temp.inputs {
		inputs = append(inputs, fmt.Sprintf("%s %s", k, v))
	}

	if len(temp.receiverType) > 0 {
		t.append(fmt.Sprintf(`
func (r %s) %s(%s) %s {
	%s
}
	`, temp.receiverType, temp.name, strings.Join(inputs, ", "), output, temp.body))
		return nil
	}

	t.append(fmt.Sprintf(`
func %s(%s) %s {
	%s
}
	`, temp.name, strings.Join(inputs, ", "), output, temp.body))

	return nil
}
