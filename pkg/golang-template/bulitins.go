package template

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/stretchr/stew/slice"
)

func linearStrContains(line string, matchTo string) bool {
	matchCharIdx := 0
	matchSize := len(matchTo) - 1

	for _, c := range line {
		if string(c) != string(matchTo[matchCharIdx]) {
			matchCharIdx = 0
			continue
		}

		matchCharIdx++
		if matchCharIdx == matchSize {
			return true
		}
	}

	return matchSize == matchCharIdx
}

type importLine struct {
	name string
	path string
}

func parseImportLine(line string) *importLine {
	path := strings.Split(line, "\"")[1]
	filePaths := strings.Split(path, "/")

	return &importLine{
		name: filePaths[len(filePaths)-1],
		path: path,
	}
}

type builtinScopeType string

const (
	nonScopeType    builtinScopeType = ""
	importScopeType builtinScopeType = "import"
)

type builtinCtx struct {
	requiredDependencies []string
	sourceMap            map[string]string
	scope                builtinScopeType
	imports              map[string]string
}

func (ctx *builtinCtx) parseBuiltinLine(
	line []byte,
) error {
	lineStr := string(line)
	if ctx.scope != nonScopeType {
		if lineStr == string(endCodeBlockChar) {
			ctx.sourceMap[string(ctx.scope)] += lineStr
			ctx.scope = nonScopeType
		}

		// not end of def yet, so pass
		return nil
	}

	// @@ check if single line import
	if ctx.scope == importScopeType {
		if slice.Contains(lineStr, endImportBlockChar) {
			ctx.scope = nonScopeType
			return nil
		}

		importLine := parseImportLine(string(line))
		ctx.imports[importLine.name] = importLine.path
	}

	for _, rd := range ctx.requiredDependencies {
		if linearStrContains(lineStr, rd) {
			ctx.scope = builtinScopeType(rd)
			ctx.sourceMap[rd] += lineStr

			return nil
		}
	}

	return nil
}

func (ctx *builtinCtx) loadBuiltinFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ctx.parseBuiltinLine(scanner.Bytes())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
