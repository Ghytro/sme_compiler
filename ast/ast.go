package ast

import (
	"errors"
	"log"

	"github.com/Ghytro/sme/helpers"
)

var astTree *AstTree
var ErrPackageAlreadyExists = errors.New("package with this name already exists in ast tree")
var ErrNoSuchPackage = errors.New("no such package declared in AST tree")
var ErrNoSuchStruct = errors.New("no such struct declared in this package")
var ErrStructAlreadyExists = errors.New("struct with this name is already declared in this package")
var ErrFieldAlreadyExists = errors.New("field with this name was already declared in this struct")

type AstModuleNode struct {
	syntaxVer        string
	cppNamespaceName *string
	goPackageName    *string

	children []*AstPackageNode
}

type AstPackageNode struct {
	name string

	children []*AstStructNode
}

func (pn AstPackageNode) GetName() string {
	return pn.name
}

type AstStructNode struct {
	name string

	children []*AstStructFieldNode
}

func (sn AstStructNode) GetName() string {
	return sn.name
}

type AstStructFieldNode struct {
	fieldType SmeType
	name      string
}

func (sn AstStructFieldNode) GetName() string {
	return sn.name
}

func (sn AstStructFieldNode) GetFieldType() SmeType {
	return sn.fieldType
}

type AstTree struct {
	root *AstModuleNode
}

func InitAstTree(syntaxVer string) {
	if astTree != nil {
		log.Fatal("Debug: an attempt to initialize AST tree twice while running a program")
	}
	astTree = &AstTree{
		root: &AstModuleNode{
			syntaxVer: syntaxVer,
		},
	}
}

// returns tree node that contains added package
// if the package exists returns a node with existing package
func AddPackage(packageName string) (*AstPackageNode, error) {
	for _, c := range astTree.root.children {
		if c.name == packageName {
			return c, ErrPackageAlreadyExists
		}
	}
	newPackageNode := &AstPackageNode{name: packageName}
	astTree.root.children = append(
		astTree.root.children,
		newPackageNode,
	)
	return newPackageNode, nil
}

// returns tree node that contains added struct
func AddStruct(packageName string, structName string) (*AstStructNode, error) {
	packageNode := new(AstPackageNode)
	for _, c := range astTree.root.children {
		if c.name == packageName {
			packageNode = c
			break
		}
	}
	if packageNode == nil {
		return nil, ErrNoSuchPackage
	}
	for _, c := range packageNode.children {
		if c.name == structName {
			return nil, ErrStructAlreadyExists
		}
	}
	newStructNode := &AstStructNode{name: structName}
	packageNode.children = append(
		packageNode.children,
		newStructNode,
	)
	return newStructNode, nil
}

func GetStructNode(packageName string, structName string) (*AstStructNode, error) {
	packageNode := new(AstPackageNode)
	for _, c := range astTree.root.children {
		if c.name == packageName {
			packageNode = c
			break
		}
	}
	if packageNode == nil {
		return nil, ErrNoSuchPackage
	}
	for _, c := range packageNode.children {
		if c.name == structName {
			return c, nil
		}
	}
	return nil, ErrNoSuchStruct
}

func AddStructField(
	packageName string,
	structName string,
	fieldName string,
	fieldType SmeType) (*AstStructFieldNode, error) {
	packageNode := new(AstPackageNode)
	for _, c := range astTree.root.children {
		if c.name == packageName {
			packageNode = c
			break
		}
	}
	if packageNode == nil {
		return nil, ErrNoSuchPackage
	}

	structNode := new(AstStructNode)
	for _, c := range packageNode.children {
		if c.name == structName {
			structNode = c
			break
		}
	}
	if structNode == nil {
		return nil, ErrNoSuchStruct
	}
	for _, c := range structNode.children {
		if c.name == fieldName {
			return nil, ErrFieldAlreadyExists
		}
	}

	newFieldNode := &AstStructFieldNode{name: fieldName, fieldType: fieldType}
	structNode.children = append(
		structNode.children,
		newFieldNode,
	)

	return newFieldNode, nil
}

func GetStructId(n *AstStructNode) uint32 {
	if n == nil {
		return 0
	}

	for _, pNode := range astTree.root.children {
		for _, sNode := range pNode.children {
			if sNode == n {
				hash, err := helpers.HashValuesUint32(pNode.name, sNode.name)
				if err != nil {
					return 0
				}
				return hash
			}
		}
	}

	return 0
}
