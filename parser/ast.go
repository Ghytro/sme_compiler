package parser

import (
	"errors"
	"log"

	"github.com/Ghytro/stme/helpers"
)

var astTree *AstTree
var errPackageAlreadyExists = errors.New("package with this name already exists in ast tree")
var errNoSuchPackage = errors.New("no such package declared in AST tree")
var errNoSuchStruct = errors.New("no such struct declared in this package")
var errStructAlreadyExists = errors.New("struct with this name is already declared in this package")
var errFieldAlreadyExists = errors.New("field with this name was already declared in this struct")

type AstTreeNode struct {
	value    interface{}
	children []*AstTreeNode
}

type SmeModule struct {
	syntaxVer        string
	cppNamespaceName *string
	goPackageName    *string
}

type SmePackage struct {
	name string
}

type SmeStruct struct {
	name string
}

type SmeStructField struct {
	fieldType SmeType
	name      string
}

type AstTree struct {
	root *AstTreeNode
}

func InitAstTree(syntaxVer string) {
	if astTree != nil {
		log.Fatal("Debug: an attempt to initialize AST tree twice while running a program")
	}
	astTree = &AstTree{
		root: &AstTreeNode{
			value: SmeModule{syntaxVer: syntaxVer},
		},
	}
}

// returns tree node that contains added package
// if the package exists returns a node with existing package
func (t *AstTree) AddPackage(packageName string) *AstTreeNode {
	for _, c := range t.root.children {
		if c.value.(*SmePackage).name == packageName {
			return c
		}
	}
	newPackageNode := &AstTreeNode{value: SmePackage{packageName}}
	t.root.children = append(
		t.root.children,
		newPackageNode,
	)
	return newPackageNode
}

// returns tree node that contains added struct
// if the struct exists in this package returns a node with existing struct
func (t *AstTree) AddStruct(packageName string, structName string) (*AstTreeNode, error) {
	packageNode := new(AstTreeNode)
	for _, c := range t.root.children {
		if c.value.(SmePackage).name == packageName {
			packageNode = c
			break
		}
	}
	if packageNode == nil {
		return nil, errNoSuchPackage
	}
	for _, c := range packageNode.children {
		if c.value.(SmeStruct).name == structName {
			return c, nil
		}
	}
	newStructNode := &AstTreeNode{value: SmeStruct{name: structName}}
	packageNode.children = append(
		packageNode.children,
		newStructNode,
	)
	return newStructNode, nil
}

func (t *AstTree) AddStructField(packageName string, structName string, fieldName string, fieldType SmeType) (*AstTreeNode, error) {
	packageNode := new(AstTreeNode)
	for _, c := range t.root.children {
		if c.value.(SmePackage).name == packageName {
			packageNode = c
			break
		}
	}
	if packageNode == nil {
		return nil, errNoSuchPackage
	}

	structNode := new(AstTreeNode)
	for _, c := range packageNode.children {
		if c.value.(SmeStruct).name == structName {
			structNode = c
			break
		}
	}
	if structNode == nil {
		return nil, errNoSuchStruct
	}
	for _, c := range structNode.children {
		if c.value.(SmeStructField).name == fieldName {
			return c, nil
		}
	}

	newFieldNode := &AstTreeNode{value: SmeStructField{name: fieldName, fieldType: fieldType}}
	structNode.children = append(
		structNode.children,
		newFieldNode,
	)

	return newFieldNode, nil
}

func (t *AstTree) GetStructId(n *AstTreeNode) uint32 {
	result := uint32(0)

	var recF func(*AstTreeNode, *AstTreeNode)
	recF = func(parent *AstTreeNode, current *AstTreeNode) {
		if current == n {
			packageName := parent.value.(SmePackage).name
			structName := current.value.(SmeStruct).name
			hash, err := helpers.HashValuesUint32(packageName, structName)
			if err != nil {
				log.Fatal("Debug: error in counting hash in GetStructId")
			}
			result = hash
			return
		}
		for _, c := range current.children {
			if result == 0 {
				recF(current, c)
			} else {
				break
			}
		}
	}
	recF(nil, t.root)
	return result
}
