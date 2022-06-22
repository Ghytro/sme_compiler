package parser

import (
	"bufio"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/Ghytro/stme/helpers"
)

const (
	lpStateUndefined = LineParserStateId(iota - 1)
	lpStateReadingSyntaxVer
	lpStateReadingPackageName
	lpStateReadingStructName
	lpStateReadingStruct
)

type LineParserStateId int
type LineParserFunc func(string, *LineParserState) LineParserStateId

type LineParserState struct {
	stateId            LineParserStateId
	currentPackageNode *AstTreeNode
	currentStructNode  *AstTreeNode
}

func NewLineParserState() *LineParserState {
	return &LineParserState{
		stateId: lpStateReadingSyntaxVer,
	}
}

func ParseReaderContent(reader *bufio.Reader) error {
	ps := NewLineParserState()
	line, err := readLine(reader)
	for err != nil {
		line = []byte(beautifyLine(string(line)))
		if string(line) == "" {
			continue
		}
		parseLine(string(line), ps)
		line, err = readLine(reader)
	}
	if line != nil && err != nil {
		helpers.PrintError(err.Error())
	}
	return nil
}

func readLine(reader *bufio.Reader) ([]byte, error) {
	line := make([]byte, 0)
	var (
		isPrefix bool = true
		linePart []byte
		err      error
	)
	for isPrefix {
		if linePart, isPrefix, err = reader.ReadLine(); err != nil {
			return nil, err
		}
		line = append(line, linePart...)
	}
	return line, nil
}

func beautifyLine(line string) string {
	// removing comments
	commentPos := strings.Index(string(line), "//")
	if commentPos != -1 {
		line = line[commentPos:]
	}
	// removing indents
	line = strings.ReplaceAll(line, "\t", "")
	// removing extra spaces
	for strings.Contains(line, "  ") {
		line = strings.ReplaceAll(line, "  ", " ")
	}
	line = strings.Trim(line, " ")

	// removing extra spaces in assignment and field enumerating
	line = helpers.RemoveExtraSpacesAroundStrings(line, ",", "=", "[", "]")
	return line
}

var parserFuncs = map[LineParserStateId]LineParserFunc{
	lpStateReadingSyntaxVer:   readSyntaxVersion,
	lpStateReadingPackageName: readPackageName,
	lpStateReadingStructName:  readStructName,
	lpStateReadingStruct:      readStruct,
}

// all the methods have the same signature, return value is the new state of parser
func parseLine(line string, ps *LineParserState) {
	ps.stateId = parserFuncs[ps.stateId](line, ps)
}

func printStateConflictError(expectedState LineParserStateId, gotState LineParserStateId) {
	stateMessages := map[LineParserStateId]string{
		lpStateReadingSyntaxVer:   "syntax version declaration",
		lpStateReadingPackageName: "package declaration",
		lpStateReadingStructName:  "struct declaration",
		lpStateReadingStruct:      "struct field declaration",
	}
	helpers.PrintError(
		fmt.Sprintf(
			"expected: %s, but got: %s",
			stateMessages[expectedState],
			stateMessages[gotState],
		),
	)
}

func readSyntaxVersion(line string, ps *LineParserState) LineParserStateId {
	if ps.stateId != lpStateReadingSyntaxVer {
		printStateConflictError(lpStateReadingSyntaxVer, ps.stateId)
	}
	if !strings.HasPrefix(line, "syntax") {
		helpers.PrintError(
			fmt.Sprintf(
				"expected: 'syntax', but got: '%s'",
				strings.Split(line, " ")[0],
			),
		)
	}

	splittedLine := strings.Split(line, " ")
	if len(splittedLine) > 2 {
		helpers.PrintError("version of syntax and the 'syntax' keyword must be separated with exactly one space")
	}
	syntaxVer := splittedLine[1]
	re, err := regexp.Compile(`[0-9]+.[0-9]+.[0-9]+`)
	if err != nil {
		log.Fatal("Debug: ", err)
	}
	if !re.Match([]byte(syntaxVer)) {
		helpers.PrintError("incorrect version format of syntax")
	}
	InitAstTree(syntaxVer)
	return lpStateReadingPackageName
}

func readPackageName(line string, ps *LineParserState) LineParserStateId {
	if ps.stateId != lpStateReadingPackageName {
		printStateConflictError(lpStateReadingPackageName, ps.stateId)
	}
	splittedLine := strings.Split(line, " ")
	if !strings.HasPrefix(line, "package") {
		helpers.PrintError(
			fmt.Sprintf(
				"expected: 'package', but got: '%s'",
				splittedLine[0],
			),
		)
	}
	if len(splittedLine) > 2 {
		helpers.PrintError("package name should not contain spaces")
	}
	packageName := splittedLine[1]
	if packageName == "" {
		helpers.PrintError("package name should be separated from 'package' keyword with exactly one space")
	}
	var err error
	ps.currentPackageNode, err = astTree.AddPackage(packageName)
	if err != nil {
		helpers.PrintError(fmt.Sprintf("package '%s' already exists", packageName))
	}
	return lpStateReadingStructName
}

func readStructName(line string, ps *LineParserState) LineParserStateId {
	if ps.stateId != lpStateReadingStructName {
		printStateConflictError(lpStateReadingStructName, ps.stateId)
	}
	splittedLine := strings.Split(line, " ")
	if !strings.HasPrefix(line, "struct") {
		helpers.PrintError(
			fmt.Sprintf(
				"expected: 'struct', but got: '%s'",
				splittedLine[0],
			),
		)
	}
	structName := splittedLine[1]
	if structName == "" {
		helpers.PrintError("struct name should be separated from 'struct' keyword with exactly one space")
	}
	packageName := ps.currentPackageNode.value.(SmePackage).name
	var err error
	ps.currentStructNode, err = astTree.AddStruct(packageName, structName)
	switch err {
	case nil:
		break
	case errNoSuchPackage:
		helpers.PrintError(
			fmt.Sprintf(
				"no such package: '%s'",
				structName,
			),
		)
	case errStructAlreadyExists:
		helpers.PrintError(
			fmt.Sprintf(
				"struct '%s' already exists in package '%s'",
				structName,
				packageName,
			),
		)
	default:
		helpers.PrintError(
			fmt.Sprintf(
				"an error occured: %s",
				err.Error(),
			),
		)
	}
	return lpStateReadingStruct
}

func readStruct(line string, ps *LineParserState) LineParserStateId {
	if ps.stateId != lpStateReadingStruct {
		printStateConflictError(lpStateReadingStruct, ps.stateId)
	}
	if line == "}" {
		return lpStateReadingStructName
	}
	splittedLine := strings.Split(line, " ")
	isOptional := false
	if splittedLine[0] == "optional" {
		isOptional = true
	}

	return lpStateReadingStruct
}
