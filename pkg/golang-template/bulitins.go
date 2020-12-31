package template

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func linearStrContains(line string, matchTo string) bool {
	matchCharIdx := 0
	matchSize := 1
	if len(matchTo) > 1 {
		matchSize = len(matchTo) - 1
	}

	for _, c := range line {
		if string(c) == string(matchTo[matchCharIdx]) {
			matchCharIdx++
		}

		// fmt.Printf("%d %d %s %s \n", matchSize, matchCharIdx, line, matchTo)
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
	fmt.Printf("importLine :: %s \n", line)
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
	requiredDependencies []*BuiltinRequirement
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
	}

	if linearStrContains(lineStr, string(importScopeType)) {
		fmt.Printf("importScope:set: %s \n", lineStr)
		ctx.scope = importScopeType
		return nil
	}

	// @@todo: check if single line import
	if ctx.scope == importScopeType {
		fmt.Printf("importScope:: %s \n", lineStr)
		if linearStrContains(lineStr, string(endImportBlockChar)) {
			fmt.Printf("dis gets hit 4 some reason? %s \n", lineStr)
			ctx.scope = nonScopeType
			return nil
		}

		importLine := parseImportLine(lineStr)
		ctx.imports[importLine.name] = importLine.path
		fmt.Printf("imports:: %s \n", ctx.imports)

		return nil
	}

	for _, rd := range ctx.requiredDependencies {
		if linearStrContains(lineStr, rd.Name) && linearStrContains(lineStr, string(rd.Type)) {
			fmt.Printf("requiredDep:: %s :: %s \n", rd.Name, lineStr)

			ctx.scope = builtinScopeType(rd.Name)
		}
	}

	ctx.sourceMap[string(ctx.scope)] += fmt.Sprintf("%s \n", lineStr)

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
