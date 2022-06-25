package ast

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/Ghytro/sme/helpers"
)

var errNoDefaultValue = errors.New("type has no default value")
var errIncorrectType = errors.New("incorrect type specified")
var errListTypeIncorrectFormat = errors.New("incorrect declaration of list")
var errMapTypeIncorrectFormat = errors.New("incorrect declaration of map")

// type pool contains pointers to already parsed type
// to not allocate the memory twice if the type was
// obtained multiple times while program execution
var typePool = map[string]map[string]SmeType{
	"int8": map[string]SmeType{
		"": &SmeInt8{},
	},
	"int16": map[string]SmeType{
		"": &SmeInt16{},
	},
	"int32": map[string]SmeType{
		"": &SmeInt32{},
	},
	"int64": map[string]SmeType{
		"": &SmeInt64{},
	},

	"uint8": map[string]SmeType{
		"": &SmeUint8{},
	},
	"uint16": map[string]SmeType{
		"": &SmeUint16{},
	},
	"uint32": map[string]SmeType{
		"": &SmeUint32{},
	},
	"uint64": map[string]SmeType{
		"": &SmeUint64{},
	},

	"float": map[string]SmeType{
		"": &SmeFloat{},
	},
	"double": map[string]SmeType{
		"": &SmeDouble{},
	},
	"string": map[string]SmeType{
		"": &SmeString{},
	},
	"bool": map[string]SmeType{
		"": &SmeBool{},
	},
	"char": map[string]SmeType{
		"": &SmeChar{},
	},
}

func IsPrimitiveType(typeName string) bool {
	_, err := ParsePrimitiveType(typeName)
	return err == nil
}

func IsParametricType(typeName string) bool {
	return strings.HasPrefix(typeName, "list") || strings.HasPrefix(typeName, "map")
}

func TypeFromString(packageName, typeName string, defaultValue string) (SmeType, error) {
	if t1, ok := typePool[typeName]; ok { // all the primitive types work here so no need in parsing
		if t2, ok := t1[defaultValue]; ok {
			return t2, nil
		}
		newType := reflect.New(reflect.ValueOf(t1[""]).Elem().Type()).Interface().(SmeType)
		if err := newType.SetDefaultValue(defaultValue); err != nil {
			return nil, err
		}
		t1[defaultValue] = newType
		return newType, nil
	}

	t, err = ParseParametricType(packageName, typeName)
	if err == nil {
		return t, nil
	}
	t, err = ParseUserDefinedType(packageName, typeName)
	if err == nil {
		return t, nil
	}
	return nil, errIncorrectType
}

func intBitSize(t SmeType) int {
	switch t.Id() {
	case uint8TypeId, int8TypeId:
		return 8
	case uint16TypeId, int16TypeId:
		return 16
	case uint32TypeId, int32TypeId:
		return 32
	case uint64TypeId, int64TypeId:
		return 64
	}
	return 0
}

var errNotAnIntType = errors.New("not an int type")

func IntFromString(t SmeType, value string) (interface{}, error) {
	bitSz := intBitSize(t)
	if bitSz == 0 {
		return nil, errNotAnIntType
	}
	parsed, err := strconv.ParseInt(value, 10, bitSz)
	if err != nil {
		return nil, err
	}
	switch t.Id() {
	case uint8TypeId:
		return uint8(parsed), nil
	case int8TypeId:
		return int8(parsed), nil
	case uint16TypeId:
		return uint16(parsed), nil
	case int16TypeId:
		return int16(parsed), nil
	case uint32TypeId:
		return uint32(parsed), nil
	case int32TypeId:
		return int32(parsed), nil
	case uint64TypeId:
		return uint64(parsed), nil
	case int64TypeId:
		return int64(parsed), nil
	}
	return nil, errNotAnIntType
}

func floatBitSize(t SmeType) int {
	switch t.Id() {
	case floatTypeId:
		return 32
	case doubleTypeId:
		return 64
	}
	return 0
}

var errNotFloatType = errors.New("not a float type")

func FloatFromString(t SmeType, value string) (interface{}, error) {
	bitSz := floatBitSize(t)
	if bitSz == 0 {
		return nil, errNotFloatType
	}
	parsed, err := strconv.ParseFloat(value, bitSz)
	if err != nil {
		return nil, err
	}
	switch t.Id() {
	case floatTypeId:
		return float32(parsed), nil
	case doubleTypeId:
		return float64(parsed), nil
	}
	return nil, errNotFloatType
}

var errNotBoolType = errors.New("not a bool value")

func BoolFromString(value string) (bool, error) {
	if value == "true" || value == "1" {
		return true, nil
	} else if value == "false" || value == "0" {
		return false, nil
	}
	return false, errNotBoolType
}

var errIncorrectDefaultValue = errors.New("incorrect data type for default value")
var errUnknownDefaultValueType = errors.New("unknown default value type")

func ParseDefaultValue(t SmeType, value string) (interface{}, error) {
	if t.IsOptional() && value == "null" {
		return nil, nil
	}
	if t.Id() == stringTypeId {
		return value, nil
	}
	if parsedInt, err := IntFromString(t, value); err == nil {
		return parsedInt, nil
	} else if err != errNotAnIntType {
		return nil, errIncorrectDefaultValue
	}
	if parsedFloat, err := FloatFromString(t, value); err == nil {
		return parsedFloat, nil
	} else if err != errNotFloatType {
		return nil, errIncorrectDefaultValue
	}
	if parsedBool, err := BoolFromString(value); err != nil {
		return parsedBool, nil
	} else if err != errNotBoolType {
		return nil, errIncorrectDefaultValue
	}
	return nil, errUnknownDefaultValueType
}

func ParsePrimitiveType(typeName string) (SmeType, error) {
	if smeType, ok := primitiveTypeStringMapping[typeName]; ok {
		return smeType, nil
	}
	return nil, errIncorrectType
}

func ParseParametricType(packageName string, typeName string) (SmeType, error) {
	var (
		tb  SmeTypeBuilder
		err error
	)
	if strings.HasPrefix(typeName, "list") {
		tb.SetType(&SmeList{})
		var (
			paramStrType string
			paramSmeType SmeType
		)
		_, err = fmt.Sscanf(typeName, "list[%s]", &paramStrType)
		if err != nil {
			return nil, errListTypeIncorrectFormat
		}
		paramSmeType, err = ParsePrimitiveType(paramStrType)
		if err != nil {
			paramSmeType, err = ParseParametricType(packageName, paramStrType)
			if err != nil {

			}
		}
		return tb.SetListValueType(paramSmeType).Done(), nil
	}
	if strings.HasPrefix(typeName, "map") {
		tb.SetType(&SmeMap{})
		var (
			paramStrType []string  = make([]string, 2)
			paramSmeType []SmeType = make([]SmeType, 2)
		)
		_, err = fmt.Sscanf(
			typeName,
			"map[%s,%s]",
			&paramStrType[0],
			&paramStrType[1],
		)
		if err != nil {
			return nil, errMapTypeIncorrectFormat
		}
		for i, t := range paramStrType {
			paramSmeType[i], err = ParsePrimitiveType(t)
			if err != nil {
				paramSmeType[i], err = ParseParametricType(packageName, t)
			}
			if err != nil {
				paramSmeType[i], err = ParseUserDefinedType(packageName, t)
			}
			if err != nil {
				return nil, err
			}
		}
		return tb.SetMapKeyType(paramSmeType[0]).SetMapValueType(paramSmeType[1]).Done(), nil
	}
	return nil, errIncorrectType
}

func ParseUserDefinedType(packageName string, typeName string) (SmeType, error) {
	node, err := astTree.GetStructNode(packageName, typeName)
	if err != nil {
		node, err = astTree.AddStruct(packageName, typeName)
		if err != nil {
			return nil, err
		}
		unknownUserTypes[packageName][typeName] = node
	}

}

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
	valueType    SmeType
	defaultValue interface{}
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

type SmeMap struct {
	SmeBaseType
	keyType      SmeType
	valueType    SmeType
	defaultValue interface{}
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

func (m *SmeMap) KeyType() SmeType {
	return m.keyType
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
	return astTree.GetStructId(uds.implNode)
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

	userDefinedStructId
)
