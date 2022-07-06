package ast

import (
	"errors"
	"log"

	"github.com/Ghytro/sme/helpers"
)

var errNoDefaultValue = errors.New("type has no default value")
var errIncorrectType = errors.New("incorrect type specified")
var errListTypeIncorrectFormat = errors.New("incorrect declaration of list")
var errMapTypeIncorrectFormat = errors.New("incorrect declaration of map")
var errIncorrectDefaultValue = errors.New("incorrect default value for the type")
var errNotParametricType = errors.New("given type is not parametric")

type SmeString struct {
	SmeBaseType
}

func (s *SmeString) IsParametric() bool {
	return false
}

func (s *SmeString) Id() uint32 {
	hash, err := helpers.HashValuesUint32(stringTypeId, s.isOptional)
	if err != nil {
		log.Fatalf("Debug: error counting hash in SmeString.Id(), %s", err)
	}
	return hash
}

func (s *SmeString) SizeOf() uint {
	return 4 // int with length of string
}

type SmeChar struct {
	SmeBaseType
}

func (c *SmeChar) IsParametric() bool {
	return false
}

func (c *SmeChar) Id() uint32 {
	hash, err := helpers.HashValuesUint32(charTypeId, c.isOptional)
	if err != nil {
		log.Fatalf("Debug: error counting hash in SmeChar.Id(), %s", err)
	}
	return hash
}

func (c *SmeChar) SizeOf() uint {
	return 1
}

func (c *SmeChar) SetDefaultValue(v string) error {
	if len(v) != 1 {
		return errIncorrectDefaultValue
	}
	c.hasDefaultValue = true
	c.defaultValue = v[:1]
	return nil
}

type SmeBool struct {
	SmeBaseType
}

func (b *SmeBool) IsParametric() bool {
	return false
}

func (b *SmeBool) Id() uint32 {
	hash, err := helpers.HashValuesUint32(boolTypeId, b.isOptional)
	if err != nil {
		log.Fatalf("Debug: error counting hash in SmeBool.Id(), %s", err)
	}
	return hash
}

func (b *SmeBool) SizeOf() uint {
	return 1
}

func (b *SmeBool) SetDefaultValue(v string) error {
	if v == "true" || v == "1" {
		b.hasDefaultValue = true
		b.defaultValue = "true"
		return nil
	} else if v == "false" || v == "0" {
		b.hasDefaultValue = true
		b.defaultValue = "false"
		return nil
	}
	return errIncorrectDefaultValue
}

type SmeList struct {
	SmeBaseType
	valueType SmeType
}

func (l *SmeList) IsParametric() bool {
	return true
}

func (l *SmeList) Id() uint32 {
	hash, err := helpers.HashValuesUint32(listTypeId, l.valueType, l.isOptional)
	if err != nil {
		log.Fatalf("Debug: error counting hash in SmeList.Id(), %s", err)
	}
	return hash
}

func (l *SmeList) SizeOf() uint {
	return 4 //size of list
}

func (l *SmeList) ValueType() SmeType {
	return l.valueType
}

func (l *SmeList) SetValueType(t SmeType) {
	l.valueType = t
}

type SmeMap struct {
	SmeBaseType
	keyType   SmeType
	valueType SmeType
}

func (m *SmeMap) IsParametric() bool {
	return true
}

func (m *SmeMap) Id() uint32 {
	hash, err := helpers.HashValuesUint32(mapTypeId, m.keyType, m.valueType, m.isOptional)
	if err != nil {
		log.Fatalf("Debug: error counting hash in SmeMap.Id(), %s", err)
	}
	return hash
}

func (m *SmeMap) SizeOf() uint {
	return 4 + 4 // id of maptype + size of map
}

func (m *SmeMap) SetKeyType(t SmeType) {
	m.keyType = t
}

func (m *SmeMap) KeyType() SmeType {
	return m.keyType
}

func (m *SmeMap) SetValueType(t SmeType) {
	m.valueType = t
}

func (m *SmeMap) ValueType() SmeType {
	return m.valueType
}

type UserDefinedStruct struct {
	SmeBaseType
	implNode *AstStructNode
}

func (uds *UserDefinedStruct) IsParametric() bool {
	return false
}

func (uds *UserDefinedStruct) Id() uint32 {
	return GetStructId(uds.implNode)
}

func (uds *UserDefinedStruct) SetImplNode(n *AstStructNode) {
	uds.implNode = n
}

const (
	// primitive types
	int8TypeId = uint32(iota)
	int16TypeId
	int32TypeId
	int64TypeId
	uint8TypeId
	uint16TypeId
	uint32TypeId
	uint64TypeId
	floatTypeId
	doubleTypeId
	stringTypeId
	charTypeId
	boolTypeId
	byteTypeId

	// parametrised types
	listTypeId
	mapTypeId
)
