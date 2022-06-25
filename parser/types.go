package parser

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Ghytro/stme/helpers"
)

var errNoDefaultValue = errors.New("type has no default value")
var errIncorrectType = errors.New("incorrect type specified")
var errListTypeIncorrectFormat = errors.New("incorrect declaration of list")
var errMapTypeIncorrectFormat = errors.New("incorrect declaration of map")
var primitiveTypeStringMapping = map[string]SmeType{
	"int8":  &SmeInt8{},
	"int16": &SmeInt16{},
	"int32": &SmeInt32{},
	"int64": &SmeInt64{},

	"uint8":  &SmeUint8{},
	"uint16": &SmeUint16{},
	"uint32": &SmeUint32{},
	"uint64": &SmeUint64{},

	"float":  &SmeDouble{},
	"double": &SmeDouble{},
	"string": &SmeString{},
	"bool":   &SmeBool{},
	"byte":   &SmeByte{},
	"char":   &SmeChar{},
}

func IsPrimitiveType(typeName string) bool {
	_, err := ParsePrimitiveType(typeName)
	return err == nil
}

func IsParametricType(typeName string) bool {
	return strings.HasPrefix(typeName, "list") || strings.HasPrefix(typeName, "map")
}

func TypeFromString(packageName string, typeName string) (SmeType, error) {
	t, err := ParsePrimitiveType(typeName)
	if err == nil {
		return t, nil
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

var unknownUserTypes = make(map[string]map[string]*AstTreeNode)

func ParseUserDefinedType(packageName string, typeName string) (SmeType, error) {
	node, err := astTree.GetStructNode(packageName, typeName)
	if err != nil {
		node, err = astTree.AddStruct(packageName, typeName)
		if err != nil {
			return nil, err
		}
		unknownUserTypes[packageName][typeName] = node
	}
	var tb SmeTypeBuilder
	return tb.SetType(&UserDefinedStruct{}).SetStructImplNode(node).Done(), nil
}

type SmeType interface {
	Id() uint32
	IsParametric() bool
	SizeOf() uint // Size of a field in the message converted to bytes
	setOptionality()
	setDefaultValue(interface{})
	IsOptional() bool
	DefaultValue() (interface{}, error)
}

type SmeTypeBuilder struct {
	pendingType SmeType
}

func NewSmeTypeBuilder() *SmeTypeBuilder {
	return &SmeTypeBuilder{}
}

func (tb *SmeTypeBuilder) SetType(t SmeType) *SmeTypeBuilder {
	tb.pendingType = t
	return tb
}

func (tb *SmeTypeBuilder) SetOptional() *SmeTypeBuilder {
	tb.pendingType.setOptionality()
	return tb
}

func (tb *SmeTypeBuilder) SetDefaultValue(defaultValue interface{}) *SmeTypeBuilder {
	tb.pendingType.setDefaultValue(defaultValue)
	return tb
}

func (tb *SmeTypeBuilder) SetListValueType(t SmeType) *SmeTypeBuilder {
	tb.pendingType.(*SmeList).valueType = t
	return tb
}

func (tb *SmeTypeBuilder) SetMapKeyType(t SmeType) *SmeTypeBuilder {
	tb.pendingType.(*SmeMap).keyType = t
	return tb
}

func (tb *SmeTypeBuilder) SetMapValueType(t SmeType) *SmeTypeBuilder {
	tb.pendingType.(*SmeMap).valueType = t
	return tb
}

func (tb *SmeTypeBuilder) SetStructImplNode(n *AstTreeNode) *SmeTypeBuilder {
	tb.pendingType.(*UserDefinedStruct).implNode = n
	return tb
}

func (tb *SmeTypeBuilder) Done() SmeType {
	return tb.pendingType
}

func (tb *SmeTypeBuilder) Reset() {
	tb.pendingType = nil
}

type SmeBaseType struct {
	isOptional      bool
	hasDefaultValue bool
}

func (bt *SmeBaseType) setOptionality() {
	bt.isOptional = true
}

func (bt *SmeBaseType) IsOptional() bool {
	return bt.isOptional
}

type SmeInt8 struct {
	SmeBaseType
	defaultValue int8
}

func (i8 *SmeInt8) IsParametric() bool {
	return false
}

func (i8 *SmeInt8) Id() uint32 {
	hash, err := helpers.HashValuesUint32(int8TypeId, i8.isOptional)
	if err != nil {
		log.Fatalf("Debug: error counting hash in SmeInt8.Id(), %s", err)
	}
	return hash
}

func (i8 *SmeInt8) SizeOf() uint {
	return 4 + 1
}

func (i8 *SmeInt8) setDefaultValue(v interface{}) {
	i8.hasDefaultValue = true
	i8.defaultValue = v.(int8)
}

func (i8 *SmeInt8) DefaultValue() (interface{}, error) {
	if i8.hasDefaultValue {
		return i8.defaultValue, nil
	}
	return nil, errNoDefaultValue
}

type SmeInt16 struct {
	SmeBaseType
	defaultValue int16
}

func (i16 *SmeInt16) IsParametric() bool {
	return false
}

func (i16 *SmeInt16) Id() uint32 {
	hash, err := helpers.HashValuesUint32(int16TypeId, i16.isOptional)
	if err != nil {
		log.Fatalf("Debug: error counting hash in SmeInt16.Id(), %s", err)
	}
	return hash
}

func (i16 *SmeInt16) SizeOf() uint {
	return 4 + 2
}

func (i16 *SmeInt16) setDefaultValue(v interface{}) {
	i16.hasDefaultValue = true
	i16.defaultValue = v.(int16)
}

func (i16 *SmeInt16) DefaultValue() (interface{}, error) {
	if i16.hasDefaultValue {
		return i16.defaultValue, nil
	}
	return nil, errNoDefaultValue
}

type SmeInt32 struct {
	SmeBaseType
	defaultValue int32
}

func (i32 *SmeInt32) IsParametric() bool {
	return false
}

func (i32 *SmeInt32) Id() uint32 {
	hash, err := helpers.HashValuesUint32(int32TypeId, i32.isOptional)
	if err != nil {
		log.Fatalf("Debug: error counting hash in SmeInt32.Id(), %s", err)
	}
	return hash
}

func (i32 *SmeInt32) SizeOf() uint {
	return 4 + 4
}

func (i32 *SmeInt32) setDefaultValue(v interface{}) {
	i32.hasDefaultValue = true
	i32.defaultValue = v.(int32)
}

func (i32 *SmeInt32) DefaultValue() (interface{}, error) {
	if i32.hasDefaultValue {
		return i32.defaultValue, nil
	}
	return nil, errNoDefaultValue
}

type SmeInt64 struct {
	SmeBaseType
	defaultValue int64
}

func (i64 *SmeInt64) IsParametric() bool {
	return false
}

func (i64 *SmeInt64) Id() uint32 {
	hash, err := helpers.HashValuesUint32(int64TypeId, i64.isOptional)
	if err != nil {
		log.Fatalf("Debug: error counting hash in SmeInt64.Id(), %s", err)
	}
	return hash
}

func (i64 *SmeInt64) SizeOf() uint {
	return 4 + 8
}

func (i64 *SmeInt64) setDefaultValue(v interface{}) {
	i64.hasDefaultValue = true
	i64.defaultValue = v.(int64)
}

func (i64 *SmeInt64) DefaultValue() (interface{}, error) {
	if i64.hasDefaultValue {
		return i64.defaultValue, nil
	}
	return nil, errNoDefaultValue
}

type SmeUint8 struct {
	SmeBaseType
	defaultValue uint8
}

func (ui8 *SmeUint8) IsParametric() bool {
	return false
}

func (ui8 *SmeUint8) Id() uint32 {
	hash, err := helpers.HashValuesUint32(uint8TypeId, ui8.isOptional)
	if err != nil {
		log.Fatalf("Debug: error counting hash in SmeUint8.Id(), %s", err)
	}
	return hash
}

func (ui8 *SmeUint8) SizeOf() uint {
	return 4 + 1
}

func (ui8 *SmeUint8) setDefaultValue(v interface{}) {
	ui8.hasDefaultValue = true
	ui8.defaultValue = v.(uint8)
}

func (ui8 *SmeUint8) DefaultValue() (interface{}, error) {
	if ui8.hasDefaultValue {
		return ui8.defaultValue, nil
	}
	return nil, errNoDefaultValue
}

type SmeUint16 struct {
	SmeBaseType
	defaultValue uint16
}

func (ui16 *SmeUint16) IsParametric() bool {
	return false
}

func (ui16 *SmeUint16) Id() uint32 {
	hash, err := helpers.HashValuesUint32(uint16TypeId, ui16.isOptional)
	if err != nil {
		log.Fatalf("Debug: error counting hash in SmeUint16.Id(), %s", err)
	}
	return hash
}

func (ui16 *SmeUint16) SizeOf() uint {
	return 4 + 2
}

func (ui16 *SmeUint16) setDefaultValue(v interface{}) {
	ui16.hasDefaultValue = true
	ui16.defaultValue = v.(uint16)
}

func (ui16 *SmeUint16) DefaultValue() (interface{}, error) {
	if ui16.hasDefaultValue {
		return ui16.defaultValue, nil
	}
	return nil, errNoDefaultValue
}

type SmeUint32 struct {
	SmeBaseType
	defaultValue uint32
}

func (ui32 *SmeUint32) IsParametric() bool {
	return false
}

func (ui32 *SmeUint32) Id() uint32 {
	hash, err := helpers.HashValuesUint32(uint32TypeId, ui32.isOptional)
	if err != nil {
		log.Fatalf("Debug: error counting hash in SmeUint32.Id(), %s", err)
	}
	return hash
}

func (ui32 *SmeUint32) SizeOf() uint {
	return 4 + 4
}

func (ui32 *SmeUint32) setDefaultValue(v interface{}) {
	ui32.hasDefaultValue = true
	ui32.defaultValue = v.(uint32)
}

func (ui32 *SmeUint32) DefaultValue() (interface{}, error) {
	if ui32.hasDefaultValue {
		return ui32.defaultValue, nil
	}
	return nil, errNoDefaultValue
}

type SmeUint64 struct {
	SmeBaseType
	defaultValue uint64
}

func (ui64 *SmeUint64) IsParametric() bool {
	return false
}

func (ui64 *SmeUint64) Id() uint32 {
	hash, err := helpers.HashValuesUint32(uint64TypeId, ui64.isOptional)
	if err != nil {
		log.Fatalf("Debug: error counting hash in SmeUint64.Id(), %s", err)
	}
	return hash
}

func (ui64 *SmeUint64) SizeOf() uint {
	return 4 + 8
}

func (ui64 *SmeUint64) setDefaultValue(v interface{}) {
	ui64.hasDefaultValue = true
	ui64.defaultValue = v.(uint64)
}

func (ui64 *SmeUint64) DefaultValue() (interface{}, error) {
	if ui64.hasDefaultValue {
		return ui64.defaultValue, nil
	}
	return nil, errNoDefaultValue
}

type SmeFloat struct {
	SmeBaseType
	defaultValue float32
}

func (f *SmeFloat) IsParametric() bool {
	return false
}

func (f *SmeFloat) Id() uint32 {
	hash, err := helpers.HashValuesUint32(floatTypeId, f.isOptional)
	if err != nil {
		log.Fatalf("Debug: error counting hash in SmeFloat.Id(), %s", err)
	}
	return hash
}

func (f *SmeFloat) SizeOf() uint {
	return 4 + 4
}

func (f *SmeFloat) setDefaultValue(v interface{}) {
	f.hasDefaultValue = true
	f.defaultValue = v.(float32)
}

func (f *SmeFloat) DefaultValue() (interface{}, error) {
	if f.hasDefaultValue {
		return f.defaultValue, nil
	}
	return nil, errNoDefaultValue
}

type SmeDouble struct {
	SmeBaseType
	defaultValue float64
}

func (d *SmeDouble) IsParametric() bool {
	return false
}

func (d *SmeDouble) Id() uint32 {
	hash, err := helpers.HashValuesUint32(doubleTypeId, d.isOptional)
	if err != nil {
		log.Fatalf("Debug: error counting hash in SmeDouble.Id(), %s", err)
	}
	return hash
}

func (d *SmeDouble) SizeOf() uint {
	return 4 + 8
}

func (d *SmeDouble) setDefaultValue(v interface{}) {
	d.hasDefaultValue = true
	d.defaultValue = v.(float64)
}

func (d *SmeDouble) DefaultValue() (interface{}, error) {
	if d.hasDefaultValue {
		return d.defaultValue, nil
	}
	return nil, errNoDefaultValue
}

type SmeString struct {
	SmeBaseType
	defaultValue string
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
	return 4 + 4 // type id + length of string
}

func (s *SmeString) setDefaultValue(v interface{}) {
	s.hasDefaultValue = true
	s.defaultValue = v.(string)
}

func (s *SmeString) DefaultValue() (interface{}, error) {
	if s.hasDefaultValue {
		return s.defaultValue, nil
	}
	return nil, errNoDefaultValue
}

type SmeChar struct {
	SmeBaseType
	defaultValue byte
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
	return 4 + 1
}

func (c *SmeChar) setDefaultValue(v interface{}) {
	c.hasDefaultValue = true
	c.defaultValue = v.(byte)
}

func (c *SmeChar) DefaultValue() (interface{}, error) {
	if c.hasDefaultValue {
		return c.defaultValue, nil
	}
	return nil, errNoDefaultValue
}

type SmeBool struct {
	SmeBaseType
	defaultValue bool
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
	return 4 + 1
}

func (b *SmeBool) setDefaultValue(v interface{}) {
	b.hasDefaultValue = true
	b.defaultValue = v.(bool)
}

func (b *SmeBool) DefaultValue() (interface{}, error) {
	if b.hasDefaultValue {
		return b.defaultValue, nil
	}
	return nil, errNoDefaultValue
}

type SmeByte struct {
	SmeBaseType
	defaultValue byte
}

func (b *SmeByte) IsParametric() bool {
	return false
}

func (b *SmeByte) Id() uint32 {
	hash, err := helpers.HashValuesUint32(byteTypeId, b.isOptional)
	if err != nil {
		log.Fatalf("Debug: error counting hash in SmeByte.Id(), %s", err)
	}
	return hash
}

func (b *SmeByte) SizeOf() uint {
	return 4 + 1
}

func (b *SmeByte) setDefaultValue(v interface{}) {
	b.hasDefaultValue = true
	b.defaultValue = v.(byte)
}

func (b *SmeByte) DefaultValue() (interface{}, error) {
	if b.hasDefaultValue {
		return b.defaultValue, nil
	}
	return nil, errNoDefaultValue
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
	return 4 + 4 // id of listType + size of list
}

func (l *SmeList) ValueType() SmeType {
	return l.valueType
}

func (l *SmeList) setDefaultValue(v interface{}) {
	l.hasDefaultValue = true
	l.defaultValue = v
}

func (l *SmeList) DefaultValue() (interface{}, error) {
	if l.hasDefaultValue {
		return l.defaultValue, nil
	}
	return nil, errNoDefaultValue
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

func (m *SmeMap) setDefaultValue(v interface{}) {
	m.hasDefaultValue = true
	m.defaultValue = v
}

func (m *SmeMap) DefaultValue() (interface{}, error) {
	if m.hasDefaultValue {
		return m.defaultValue, nil
	}
	return nil, errNoDefaultValue
}

func (m *SmeMap) KeyType() SmeType {
	return m.keyType
}

func (m *SmeMap) ValueType() SmeType {
	return m.valueType
}

type UserDefinedStruct struct {
	SmeBaseType
	implNode *AstTreeNode
}

func (uds *UserDefinedStruct) IsParametric() bool {
	return false
}

func (uds *UserDefinedStruct) IsOptional() bool {
	return uds.isOptional
}

func (uds *UserDefinedStruct) SizeOf() uint {
	// maybe deprecate sizeof
	return 0
}

func (uds *UserDefinedStruct) setDefaultValue(v interface{}) {
}

func (uds *UserDefinedStruct) DefaultValue() (interface{}, error) {
	return nil, nil
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
