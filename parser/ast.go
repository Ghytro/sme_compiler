package parser

import (
	"errors"
	"log"

	"github.com/Ghytro/stme/helpers"
)

var astTree *AstTree
var errPackageAlreadyExists = errors.New("package with this name already exists in ast tree")
var errNoSuchPackage = errors.New("no such package declared in AST tree")
var errStructAlreadyExists = errors.New("struct with this name is already declared in this package")

type AstTree struct {
	root *AstTreeNode
}

func InitAstTree(syntaxVer string) {
	if astTree != nil {
		log.Fatal("Debug: an attempt to initialize AST tree twice while running a program")
	}
	astTree = &AstTree{
		root: &AstTreeNode{
			value: SmeModule{syntaxVer},
		},
	}
}

func (t *AstTree) AddPackage(packageName string) error {
	if t.root.children == nil {
		t.root.children = make([]*AstTreeNode, 0)
	}
	for _, c := range t.root.children {
		if c.value.(*SmePackage).name == packageName {
			return errPackageAlreadyExists
		}
	}
	t.root.children = append(
		t.root.children,
		&AstTreeNode{
			value: SmePackage{packageName},
		},
	)
	return nil
}

func (t *AstTree) AddStruct(packageName string, structName string) error {
	packageNode := new(AstTreeNode)
	for _, c := range t.root.children {
		if c.value.(SmePackage).name == packageName {
			packageNode = c
			break
		}
	}
	if packageNode == nil {
		return errNoSuchPackage
	}
	for _, c := range packageNode.children {
		if c.value.(SmeStruct).name == structName {
			return errStructAlreadyExists
		}
	}
	packageNode.children = append(
		packageNode.children,
		&AstTreeNode{
			value: SmeStruct{
				name: structName,
			},
		},
	)
	return nil
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

type AstTreeNode struct {
	value    interface{}
	children []*AstTreeNode
}

type SmeModule struct {
	syntaxVer string
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
