package parser

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Ghytro/sme/ast"
	"github.com/Ghytro/sme/helpers"
)

const (
	lpStateUndefined = LineParserStateId(iota - 1)
	lpStateReadingSyntaxVer
	lpStateReadingPackageName
	lpStateReadingStructName
	lpStateReadingStruct
)

type LineParserStateId int
type LineParserFunc func(string, *LineParserState) (LineParserStateId, error)

type LineParserState struct {
	stateId            LineParserStateId
	lineNumber         int
	currentPackageNode *ast.AstPackageNode
	currentStructNode  *ast.AstStructNode
}

func NewLineParserState() *LineParserState {
	return &LineParserState{
		stateId:    lpStateReadingSyntaxVer,
		lineNumber: 1,
	}
}

func ParseFileContent(file *os.File) error {
	ps := NewLineParserState()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = beautifyLine(string(line))
		if line == "" {
			continue
		}
		if err := parseLine(string(line), ps); err != nil {
			return err
		}
	}
	return nil
}

func beautifyLine(line string) string {
	// removing comments
	commentPos := strings.Index(string(line), "//")
	if commentPos != -1 {
		line = line[commentPos:]
	}
	//removing extra tabulations in beginning and ending
	line = strings.Trim(line, " \t")
	return line
}

var parserFuncs = map[LineParserStateId]LineParserFunc{
	lpStateReadingSyntaxVer:   readSyntaxVersion,
	lpStateReadingPackageName: readPackageName,
	lpStateReadingStructName:  readStructName,
	lpStateReadingStruct:      readStruct,
}

// all the methods have the same signature, return value is the new state of parser
func parseLine(line string, ps *LineParserState) error {
	var err error
	ps.stateId, err = parserFuncs[ps.stateId](line, ps)
	if err != nil {
		return err
	}
	ps.lineNumber++
	return nil
}

type LineParserStateConflictErr struct {
	expectedState LineParserStateId
	gotState      LineParserStateId
}

func newLineParserStateConflictErr(
	expected LineParserStateId,
	got LineParserStateId) *LineParserStateConflictErr {
	return &LineParserStateConflictErr{expected, got}
}

func (lpsce *LineParserStateConflictErr) Error() string {
	stateMessages := map[LineParserStateId]string{
		lpStateReadingSyntaxVer:   "syntax version declaration",
		lpStateReadingPackageName: "package declaration",
		lpStateReadingStructName:  "struct declaration",
		lpStateReadingStruct:      "struct field declaration",
	}
	return fmt.Sprintf(
		"expected: %s, but got: %s",
		stateMessages[lpsce.expectedState],
		stateMessages[lpsce.gotState],
	)
}

type ExpectedSyntaxErr struct {
	SyntaxErr
}

func newExpectedSyntaxErr(line int, got string) (ese *ExpectedSyntaxErr) {
	ese.line = line
	ese.column = 0
	ese.description = fmt.Sprintf("expected: 'syntax' keyword, got: %s", got)
	return ese
}

type IncorrectSyntaxVerErr struct {
	SyntaxErr
}

func newIncorrectSyntaxVerErr(line int, got string) (isve *IncorrectSyntaxVerErr) {
	isve.line = line
	isve.column = len("syntax") + 1
	isve.description = fmt.Sprintf("incorrect syntax version specified: %s", got)
	return isve
}

func readSyntaxVersion(line string, ps *LineParserState) (LineParserStateId, error) {
	if ps.stateId != lpStateReadingSyntaxVer {
		return lpStateUndefined, newLineParserStateConflictErr(lpStateReadingSyntaxVer, ps.stateId)
	}
	if !strings.HasPrefix(line, "syntax") {
		return lpStateUndefined, newExpectedSyntaxErr(ps.lineNumber, strings.Split(line, " ")[0])
	}
	versionOffset := len("syntax")
	for helpers.EqualsAny(line[versionOffset], ' ', '\t') && versionOffset < len(line) {
		versionOffset++
	}
	if versionOffset == len(line) {
		return lpStateUndefined, newExpectedSyntaxErr(ps.lineNumber, "end of line")
	}
	syntaxVer := line[versionOffset:]
	syntaxVerCorrect, err := helpers.MatchString(`[0-9]+.[0-9]+.[0-9]+`, syntaxVer)
	if err != nil {
		log.Fatal("Debug: ", err)
	}
	if !syntaxVerCorrect {
		return lpStateUndefined, newIncorrectSyntaxVerErr(ps.lineNumber, syntaxVer)
	}
	ast.InitAstTree(syntaxVer)
	return lpStateReadingPackageName, nil
}

type ExpectedPackageKwErr struct {
	SyntaxErr
}

func newExpectedPackageKwErr(line int, got string) *ExpectedPackageKwErr {
	epke := new(ExpectedPackageKwErr)
	epke.line = line
	epke.column = 0
	epke.description = fmt.Sprintf("expected 'package' keyword, got: %s", got)
	return epke
}

type IncorrectPackageNameErr struct {
	SyntaxErr
}

func newIncorrectPackageNameErr(line int, got string) *IncorrectPackageNameErr {
	ipne := new(IncorrectPackageNameErr)
	ipne.line = line
	ipne.column = len("package")
	ipne.description = fmt.Sprintf("incorrect format of package name: %s", got)
	return ipne
}

func readPackageName(line string, ps *LineParserState) (LineParserStateId, error) {
	if ps.stateId != lpStateReadingPackageName {
		return lpStateUndefined, newLineParserStateConflictErr(lpStateReadingPackageName, ps.stateId)
	}
	if !strings.HasPrefix(line, "package") {
		return lpStateUndefined, newExpectedPackageKwErr(ps.lineNumber, strings.Split(line, " ")[0])
	}
	packageNameOffset := len("package")
	for helpers.EqualsAny(line[packageNameOffset], ' ', '\t') {
		packageNameOffset++
	}
	packageName := line[packageNameOffset:]
	packageNameCorrect, err := helpers.MatchString(`[A-Za-z][A-Za-z0-9_]*`, packageName)
	if err != nil {
		helpers.PrintError("debug: incorrect regular expression at readPackageName")
	}
	if !packageNameCorrect {
		return lpStateUndefined, newIncorrectPackageNameErr(ps.lineNumber, packageName)
	}
	ps.currentPackageNode, _ = ast.AddPackage(packageName)
	return lpStateReadingStructName, nil
}

type ExpectedStructKwErr struct {
	SyntaxErr
}

func newExpectedStructKwErr(line int, got string) *ExpectedStructKwErr {
	eske := new(ExpectedStructKwErr)
	eske.line = line
	eske.column = 0
	eske.description = fmt.Sprintf("expected 'package' keyword, got: %s", got)
	return eske
}

type ExpectedOpeningCurlyBraceErr struct {
	SyntaxErr
}

func newExpectedOpeningCurlyBraceErr(line int, column int) *ExpectedOpeningCurlyBraceErr {
	eocbe := new(ExpectedOpeningCurlyBraceErr)
	eocbe.line = line
	eocbe.column = column
	eocbe.description = "expected opening curly brace"
	return eocbe
}

type NoStructNameErr struct {
	SyntaxErr
}

func newNoStructNameErr(line int, column int) *NoStructNameErr {
	nsne := new(NoStructNameErr)
	nsne.line = line
	nsne.column = column
	nsne.description = "expected struct name"
	return nsne
}

type NoSuchPackageErr struct {
	SyntaxErr
}

func newNoSuchPackageErr(line int, column int, packageName string) *NoSuchPackageErr {
	nspe := new(NoSuchPackageErr)
	nspe.line = line
	nspe.column = column
	nspe.description = fmt.Sprintf("no such package: %s", packageName)
	return nspe
}

type StructAlreadyExistsErr struct {
	SyntaxErr
}

func newStructAlreadyExistsErr(line int, column int, structName string) *StructAlreadyExistsErr {
	saee := new(StructAlreadyExistsErr)
	saee.line = line
	saee.column = column
	saee.description = fmt.Sprintf("struct already exists: %s", structName)
	return saee
}

func readStructName(line string, ps *LineParserState) (LineParserStateId, error) {
	if ps.stateId != lpStateReadingStructName {
		return lpStateUndefined, newLineParserStateConflictErr(lpStateReadingStructName, ps.stateId)
	}
	if !strings.HasPrefix(line, "struct") {
		return lpStateUndefined, newExpectedStructKwErr(ps.lineNumber, strings.Split(line, " ")[0])
	}
	idx := len("struct")
	for idx < len(line) && helpers.EqualsAny(line[idx], ' ', '\t') {
		idx++
	}
	var structNameBuff strings.Builder
	for idx < len(line) && helpers.IsAllowedStructChar(line[idx]) {
		structNameBuff.WriteByte(line[idx])
		idx++
	}
	for idx < len(line) && helpers.EqualsAny(line[idx], ' ', '\t') {
		idx++
	}
	isCorrectStructName, err := helpers.MatchString(`[A-za-z][A-za-z0-9_]*`, structNameBuff.String())
	if err != nil {
		helpers.PrintError("debug: incorrect regex at readStructName")
	}
	if !isCorrectStructName {
		return lpStateUndefined, newSyntaxError(ps.lineNumber, len("struct "), fmt.Sprintf("incorrect name of struct: %s", structNameBuff.String()))
	}
	if idx == len(line) || line[idx] != '{' {
		return lpStateUndefined, newExpectedOpeningCurlyBraceErr(ps.lineNumber, idx)
	}
	structName := structNameBuff.String()
	if structName == "" {
		return lpStateUndefined, newNoStructNameErr(ps.lineNumber, idx)
	}
	packageName := ps.currentPackageNode.GetName()
	ps.currentStructNode, err = ast.AddStruct(packageName, structName)
	switch err {
	case nil:
		break
	case ast.ErrNoSuchPackage:
		return lpStateUndefined, newNoSuchPackageErr(ps.lineNumber, idx, packageName)
	default:
		helpers.PrintError(
			fmt.Sprintf(
				"an error occured: %s",
				err.Error(),
			),
		)
	}
	return lpStateReadingStruct, nil
}

func readStruct(line string, ps *LineParserState) (LineParserStateId, error) {
	if ps.stateId != lpStateReadingStruct {
		return lpStateUndefined, newLineParserStateConflictErr(lpStateReadingStruct, ps.stateId)
	}
	if line == "}" {
		return lpStateReadingStructName, nil
	}
	declData, err := parseFieldDeclarations(line, ps.lineNumber)
	if err != nil {
		return lpStateUndefined, err
	}

	// parse the type of fields
	packageName := ps.currentPackageNode.GetName()
	structName := ps.currentStructNode.GetName()
	for _, f := range declData.Fields {
		fieldSmeType, err := ast.TypeFromString(
			packageName,
			declData.FieldsType,
			declData.IsOptional,
			f.HasDefaultValue,
			f.DefaultValue,
		)
		if err != nil {
			return lpStateUndefined, err
		}
		_, err = ast.AddStructField(packageName, structName, f.Name, fieldSmeType)
		if err != nil {
			return lpStateUndefined, err
		}
	}
	return lpStateReadingStruct, nil
}

type fieldData struct {
	Name            string
	DefaultValue    interface{}
	HasDefaultValue bool
}

type fieldDeclData struct {
	IsOptional bool
	FieldsType string
	Fields     []fieldData
}

type SyntaxErr struct {
	line        int
	column      int
	description string
}

func newSyntaxError(line int, column int, desc string) *SyntaxErr {
	return &SyntaxErr{line, column, desc}
}

func (se *SyntaxErr) Error() string {
	return fmt.Sprintf(
		"syntax error at %d:%d - %s",
		se.line,
		se.column,
		se.description,
	)
}

func parseFieldDeclarations(line string, lineNumber int) (result fieldDeclData, err error) {
	const (
		stateReadingTypeName = iota
		stateReadingFieldName
		stateReadingDefaultValue
	)
	idx := 0
	state := stateReadingTypeName
	if strings.HasPrefix(line, "optional") {
		result.IsOptional = true
		idx = len("optional")
	}
	var (
		buffer       strings.Builder
		pendingField fieldData
	)
	for idx < len(line) {
		switch state {
		case stateReadingTypeName:
			for idx < len(line) && helpers.EqualsAny(line[idx], ' ', '\t') {
				idx++
			}
			for idx < len(line) && !helpers.EqualsAny(line[idx], ' ', '\t') {
				buffer.WriteByte(line[idx])
				idx++
			}
			if strings.HasPrefix((buffer.String()), "map[") {
				for idx < len(line) && line[idx] != ']' {
					buffer.WriteByte(line[idx])
					idx++
				}
				if idx == len(line) {
					return fieldDeclData{}, newSyntaxError(lineNumber, idx, "expected closing bracket, but got: end of line")
				}
				buffer.WriteByte(']')
				typeName := strings.ReplaceAll(buffer.String(), " ", "")
				m, err := helpers.MatchString(`map\[.*,.*\]`, typeName)
				if err != nil {
					helpers.PrintError("debug: incorrect regex at parseFieldDeclarations")
				}
				if !m {
					return fieldDeclData{}, newSyntaxError(lineNumber, idx-len(typeName), "incorrect declaration of a map")
				}
				buffer.Reset()
				buffer.WriteString(typeName)
			} else if !strings.HasPrefix(buffer.String(), "list[") && !ast.IsPrimitiveTypeName(buffer.String()) {
				correctTypeName, err := helpers.MatchString(`[A-Za-z]?[A-Za-z0-9_](.[A-Za-z]?[A-Za-z0-9_])?`, buffer.String())
				if err != nil {
					helpers.PrintError("debug: unable to compile regexp at parseFieldDeclarations")
				}
				if !correctTypeName {
					return fieldDeclData{}, newSyntaxError(lineNumber, 0, fmt.Sprintf("incorrect type name: %s", buffer.String()))
				}
			}
			result.FieldsType = buffer.String()
			state = stateReadingFieldName
		case stateReadingFieldName:
			for helpers.EqualsAny(line[idx], ' ', '\t') {
				idx++
			}
			for idx < len(line) && !helpers.EqualsAny(line[idx], ' ', '\t', ',', '=') {
				buffer.WriteByte(line[idx])
				idx++
			}
			pendingField.Name = buffer.String()
			if pendingField.Name[0] >= '0' && pendingField.Name[0] <= '9' {
				return fieldDeclData{},
					newSyntaxError(
						lineNumber,
						idx-len(pendingField.Name),
						"field name can not start with a number",
					)
			}
			for idx < len(line) && !helpers.EqualsAny(line[idx], '=', ',') {
				idx++
			}
			if idx < len(line) && line[idx] == '=' {
				state = stateReadingDefaultValue
			} else {
				idx++
				result.Fields = append(result.Fields, pendingField)
				pendingField = fieldData{}
			}
		case stateReadingDefaultValue:
			for helpers.EqualsAny(line[idx], ' ', '\t') {
				idx++
			}
			if result.FieldsType == "string" {
				for idx < len(line) && !helpers.EqualsAny(line[idx], '"') {
					idx++
				}
				if idx == len(line) {
					return fieldDeclData{}, newSyntaxError(lineNumber, idx, "expected opening quotes in string default value declaration, but got: end of line")
				}
				idx++
				for idx < len(line) && !helpers.EqualsAny(line[idx], '"') {
					buffer.WriteByte(line[idx])
					idx++
				}
				if idx == len(line) {
					return fieldDeclData{},
						newSyntaxError(
							lineNumber,
							idx,
							"expected closing quotes, but got: end of line",
						)
				}
				idx += 2
			} else {
				for helpers.EqualsAny(line[idx], ' ', '\t') {
					idx++
				}
				for idx < len(line) && !helpers.EqualsAny(line[idx], ' ', '\t', ',') {
					buffer.WriteByte(line[idx])
					idx++
				}
				for idx < len(line) && !helpers.EqualsAny(line[idx], ',') {
					idx++
				}
				idx++
				if buffer.String() == "" {
					return fieldDeclData{}, newSyntaxError(lineNumber, idx, "expected default value declaration for value, but got: end of line")
				}
			}
			if !result.IsOptional && buffer.String() == "null" {
				return fieldDeclData{}, newSyntaxError(lineNumber, idx, "non-optional types cannot hold null as default value")
			}
			pendingField.HasDefaultValue = true
			if buffer.String() == "null" {
				pendingField.DefaultValue = nil
			} else {
				pendingField.DefaultValue = buffer.String()
			}
			result.Fields = append(result.Fields, pendingField)
			pendingField = fieldData{}
			state = stateReadingFieldName
		}
		buffer.Reset()
	}
	return result, nil
}
