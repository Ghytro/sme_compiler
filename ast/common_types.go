package ast

import (
	"errors"
	"log"
	"regexp"
	"strings"

	"github.com/Ghytro/sme/helpers"
)

var errNoDefaultValue = errors.New("type has no default value")
var errIncorrectType = errors.New("incorrect type specified")
var errListTypeIncorrectFormat = errors.New("incorrect declaration of list")
var errMapTypeIncorrectFormat = errors.New("incorrect declaration of map")
var errIncorrectDefaultValue = errors.New("incorrect default value for the type")

// type pool contains pointers to already parsed type
// to not allocate the memory twice if the type was
// obtained multiple times while program execution

type smeTypePool struct {
	requiredTypes *requiredTypesNode
	optionalTypes *optionalTypesNode
}

var errNoSuchType = errors.New("no such type added to pool")

func (tp *smeTypePool) getType(typeName string, isOptional bool, hasDefaultValue bool, defaultValue interface{}) (SmeType, error) {
	if isOptional {
		if hasDefaultValue {
			if v, ok := tp.optionalTypes.defaultValueTypes[typeName][defaultValue]; !ok {
				return nil, errNoSuchType
			} else {
				return v, nil
			}
		}
		if v, ok := tp.optionalTypes.noDefaultValueTypes[typeName]; !ok {
			return nil, errNoSuchType
		} else {
			return v, nil
		}
	}
	if hasDefaultValue {
		if v, ok := tp.requiredTypes.defaultValueTypes[typeName][defaultValue.(string)]; !ok {
			return nil, errNoSuchType
		} else {
			return v, nil
		}
	}
	if v, ok := tp.requiredTypes.noDefaultValueTypes[typeName]; !ok {
		return nil, errNoSuchType
	} else {
		return v, nil
	}
}

func (tp *smeTypePool) addType(typeName string, t SmeType) {
	if t.IsOptional() {
		if v, err := t.DefaultValue(); err == nil {
			tp.optionalTypes.defaultValueTypes[typeName][v] = t
		} else {
			tp.optionalTypes.noDefaultValueTypes[typeName] = t
		}
	} else {
		if v, err := t.DefaultValue(); err == nil {
			tp.requiredTypes.defaultValueTypes[typeName][v] = t
		} else {
			tp.requiredTypes.noDefaultValueTypes[typeName] = t
		}
	}
}

func newSmeTypePool() *smeTypePool {
	result := new(smeTypePool)
	result.requiredTypes = newRequiredTypesNode()
	result.optionalTypes = newOptionalTypesNode()
	return result
}

type requiredTypesNode struct {
	noDefaultValueTypes noDefaultValueTypes
	defaultValueTypes   requiredDefaultValueTypes
}

func newRequiredTypesNode() *requiredTypesNode {
	result := new(requiredTypesNode)
	result.noDefaultValueTypes = makeNoDefaultValueTypes(false)
	result.defaultValueTypes = make(requiredDefaultValueTypes)
	return result
}

type noDefaultValueTypes map[string]SmeType

func makeNoDefaultValueTypes(isOptional bool) noDefaultValueTypes {
	result := make(noDefaultValueTypes)
	result = map[string]SmeType{
		"int8":   &SmeInt8{},
		"int16":  &SmeInt16{},
		"int32":  &SmeInt32{},
		"int64":  &SmeInt64{},
		"uint8":  &SmeUint8{},
		"uint16": &SmeUint16{},
		"uint32": &SmeUint32{},
		"uint64": &SmeUint64{},
		"float":  &SmeFloat{},
		"double": &SmeDouble{},
		"string": &SmeString{},
		"bool":   &SmeBool{},
		"char":   &SmeChar{},
	}
	if isOptional {
		for k := range result {
			result[k].SetOptionality()
		}
	}
	return result
}

type requiredDefaultValueTypes map[string]map[string]SmeType

type optionalTypesNode struct {
	noDefaultValueTypes noDefaultValueTypes
	defaultValueTypes   optionalDefaultValueTypes
}

func newOptionalTypesNode() *optionalTypesNode {
	result := new(optionalTypesNode)
	result.noDefaultValueTypes = makeNoDefaultValueTypes(true)
	result.defaultValueTypes = make(optionalDefaultValueTypes)
	return result
}

type optionalDefaultValueTypes map[string]map[interface{}]SmeType

var typePool = newSmeTypePool()

func IsPrimitiveTypeName(typeName string) bool {
	re, err := regexp.Compile(`u?int(8|16|32|64)|float|double|string|bool|char`)
	if err != nil {
		helpers.PrintError("debug: error compiling regex at isPrimitiveTypeName")
	}
	return re.Match([]byte(typeName))
}

func IsParametricTypeName(typeName string) bool {
	return strings.HasPrefix(typeName, "map") || strings.HasPrefix(typeName, "list")
}

func unwrapTypeName(packageName, typeName string) (string, error) {
	if IsPrimitiveTypeName(typeName) {
		return typeName
	}
	if IsParametricTypeName(typeName) {
		reslut, err := unwrapParametricTypeName(packageName, typeName)
		if err != nil {
			return "", err
		}
	}

	splittedTypeName := strings.Split(typeName, ".")
	if len(splittedTypeName) == 1 {
		if !IsParametricTypeName(typeName) {

		}
	}
}

func unwrapParametricTypeName(packageName, typeName string) string {

}

func TypeFromString(packageName, typeName string, isOptional bool, hasDefaultValue string, defaultValue interface{}) (SmeType, error) {
	typeName = unwrapTypeName(packageName, typeName)
	t, err := typePool.getType()

}

// func TypeFromString(packageName, typeName string, isOptional bool, defaultValue *string) (SmeType, error) {
// 	if isOptional {
// 		if t1, ok := typePool[typeName]; ok {
// 			if t2, ok := t1[*defaultValue]; ok {
// 				return t2, nil
// 			}
// 			newType := reflect.New(reflect.ValueOf(t1[""]).Elem().Type()).Interface().(SmeType)
// 			if err := newType.SetDefaultValue(*defaultValue); err != nil {
// 				return nil, err
// 			}
// 			t1[*defaultValue] = newType
// 			return newType, nil
// 		}
// 	} else {
// 		if t1, ok := optionalTypePool[typeName]; ok {
// 			for val, t := range t1 {
// 				if val == nil && defaultValue == nil {
// 					return t, nil
// 				} else if val != nil {
// 					if *val == *defaultValue {
// 						return t, nil
// 					}
// 				}
// 			}
// 			newType := reflect.New(reflect.ValueOf(t1[nil]).Elem().Type()).Interface().(SmeType)
// 			newType.SetOptionality()
// 			if defaultValue != nil {
// 				if err := newType.SetDefaultValue(*defaultValue); err != nil {
// 					return nil, err
// 				}
// 			}
// 			t1[defaultValue] = newType
// 		}
// 	}

// 	t, err := AddParametricType(packageName, typeName)
// 	if err == nil {
// 		return t, nil
// 	}
// 	t, err = AddUserDefinedType(packageName, typeName)
// 	if err == nil {
// 		return t, nil
// 	}
// 	return nil, errIncorrectType
// }

// var errIncorrectDefaultValue = errors.New("given default value does not match a type")

// func AddParametricType(packageName string, typeName string) (SmeType, error) {
// 	if strings.HasPrefix(typeName, "list") {
// 		var paramStrType string
// 		_, err := fmt.Sscanf(typeName, "list[%s]", &paramStrType)
// 		if err != nil {
// 			return nil, errListTypeIncorrectFormat
// 		}
// 		valueType, err := TypeFromString(packageName, paramStrType, "")
// 		if err != nil {
// 			return nil, err
// 		}
// 		resultType := &SmeList{}
// 		resultType.SetValueType(valueType)
// 		typePool[typeName] = map[string]SmeType{"": resultType}
// 		return resultType, nil
// 	}
// 	if strings.HasPrefix(typeName, "map") {
// 		paramStrType := [2]string{}
// 		paramSmeType := [2]SmeType{}
// 		_, err := fmt.Sscanf(
// 			typeName,
// 			"map[%s,%s]",
// 			&paramStrType[0],
// 			&paramStrType[1],
// 		)
// 		if err != nil {
// 			return nil, errMapTypeIncorrectFormat
// 		}
// 		for i, t := range paramStrType {
// 			paramSmeType[i], err = TypeFromString(packageName, t, "")
// 			if err != nil {
// 				return nil, err
// 			}
// 		}
// 		resultType := &SmeMap{}
// 		resultType.SetKeyType(paramSmeType[0])
// 		resultType.SetValueType(paramSmeType[1])
// 		typePool[typeName] = map[string]SmeType{"": resultType}
// 		return resultType, nil
// 	}
// 	return nil, errIncorrectType
// }

// func AddUserDefinedType(packageName string, typeName string) (SmeType, error) {
// 	splittedTypeName := strings.Split(typeName, ".")
// 	if len(splittedTypeName) == 2 {
// 		packageName, typeName = splittedTypeName[0], splittedTypeName[1]
// 	} else if len(splittedTypeName) != 1 {
// 		return nil, errIncorrectType
// 	}
// 	node, err := GetStructNode(packageName, typeName)
// 	if err != nil {
// 		node, err = AddStruct(packageName, typeName)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	if len(splittedTypeName) == 2 {
// 		typeName = packageName + "." + typeName
// 	}
// 	strct := &UserDefinedStruct{}
// 	strct.SetImplNode(node)
// 	typePool[typeName] = map[string]SmeType{"": strct}
// 	return strct, nil
// }

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
