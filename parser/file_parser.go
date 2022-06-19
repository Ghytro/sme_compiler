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
	lpStateUndefined = iota - 1
	lpStateReadingSyntaxVer
	lpStateReadingPackageName
	lpStateReadingStructName
	lpStateReadingStruct
)

type LineParserStateId int

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

	// removing extra spaces
	for strings.Contains(line, "  ") {
		line = strings.ReplaceAll(line, "  ", " ")
	}
	line = strings.TrimRight(line, " ")

	// removing extra spaces in assignment and field enumerating
	line = strings.ReplaceAll(
		strings.ReplaceAll(
			strings.ReplaceAll(
				strings.ReplaceAll(
					line,
					" =",
					"=",
				),
				"= ",
				"=",
			),
			" ,",
			",",
		),
		", ",
		",",
	)
	return line
}

func parseLine(line string, ps *LineParserState) {
	switch ps.stateId {

	case lpStateReadingSyntaxVer:
		readSyntaxVersion(line, ps)

	case lpStateReadingPackageName:
		readPackageName(line, ps)

	case lpStateReadingStructName:
		readStructName(line, ps)

	case lpStateReadingStruct:
		readStruct(line, ps)
	}
}

func checkZeroIndentation(line string) bool {
	return !(strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t"))
}

func readSyntaxVersion(line string, ps *LineParserState) {
	if !checkZeroIndentation(line) {
		helpers.PrintError("wrong intendation for specifying syntax version")
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
	re, err := regexp.Compile(`\d*.\d*.\d*`)
	if err != nil {
		log.Fatal("Debug: ", err)
	}
	if !re.Match([]byte(syntaxVer)) {
		helpers.PrintError("incorrect version format of syntax")
	}
	InitAstTree(syntaxVer)
	ps.isReadingSyntax = false
	ps.isReadingPackageName = true
}

func readPackageName(line string, ps *LineParserState) {
	if !checkZeroIndentation(line) {
		helpers.PrintError("wrong intendation for specifying package name")
	}
	if !strings.HasPrefix(line, "package") {
		helpers.PrintError(
			fmt.Sprintf(
				"expected: 'package', but got: '%s'",
				strings.Split(line, " ")[0],
			),
		)
	}

	splittedLine := strings.Split(line, " ")
	if len(splittedLine) > 2 {
		helpers.PrintError("package name should not contain spaces")
	}
	packageName := splittedLine[1]
	ps.currentPackageNode = astTree.AddPackage(packageName)
	ps.isReadingPackageName = false
	ps.isReadingStructName = true
}

func readStruct(line string, ps *LineParserState) {

}
