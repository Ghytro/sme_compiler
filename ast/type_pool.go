package ast

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/Ghytro/sme/helpers"
)

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

func newSmeTypeByName(typeName string, isOptional, hasDefaulValue bool, defaultValue interface{}) (SmeType, error) {
	var baseType SmeType
	if IsPrimitiveTypeName(typeName) {
		switch typeName {
		case "uint8":
			baseType = &SmeUint8{}
		case "int8":
			baseType = &SmeInt8{}
		case "uint16":
			baseType = &SmeUint16{}
		case "int16":
			baseType = &SmeInt16{}
		case "uint32":
			baseType = &SmeUint32{}
		case "int32":
			baseType = &SmeInt32{}
		case "uint64":
			baseType = &SmeUint64{}
		case "int64":
			baseType = &SmeInt64{}
		case "float":
			baseType = &SmeFloat{}
		case "double":
			baseType = &SmeDouble{}
		case "bool":
			baseType = &SmeBool{}
		case "char":
			baseType = &SmeChar{}
		case "string":
			baseType = &SmeString{}
		}
		if isOptional {
			baseType.SetOptionality()
		}
		if hasDefaulValue {
			switch v := defaultValue.(type) {
			case string:
				if v != "" {
					baseType.SetDefaultValue(v)
				}
			}
		}
		return baseType, nil
	}
	if IsParametricTypeName(typeName) {
		if strings.HasPrefix(typeName, "list") {
			baseType = &SmeList{}
			var valueTypeName string
			_, err := fmt.Sscanf(typeName, "list[%s]", &valueTypeName)
			if err != nil {
				return nil, err
			}
			valueType, err := TypeFromString("", valueTypeName, false, false, nil)
			if err != nil {
				return nil, err
			}
			baseType.(*SmeList).SetValueType(valueType)
			return baseType, nil
		}
		if strings.HasPrefix(typeName, "map") {
			baseType = &SmeMap{}
			var keyTypeName, valueTypeName string
			_, err := fmt.Sscanf(typeName, "map[%s,%s]", &keyTypeName, &valueTypeName)
			if err != nil {
				return nil, err
			}
			keyType, err := TypeFromString("", keyTypeName, false, false, nil)
			if err != nil {
				return nil, err
			}
			valueType, err := TypeFromString("", valueTypeName, false, false, nil)
			if err != nil {
				return nil, err
			}
			baseType.(*SmeMap).SetKeyType(keyType)
			baseType.(*SmeMap).SetValueType(valueType)
			return baseType, nil
		}
	}
	baseType = &UserDefinedStruct{}
	splittedTypeName := strings.Split(typeName, ".")
	packageName, structName := splittedTypeName[0], splittedTypeName[1]
	node, err := AddStruct(packageName, structName)
	if err != nil {
		if err == ErrStructAlreadyExists {
			node, err = GetStructNode(packageName, structName)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	baseType.(*UserDefinedStruct).SetImplNode(node)
	return baseType, nil
}

// only unwrapped type name should be passed
func (tp *smeTypePool) addType(typeName string, isOptional, hasDefaultValue bool, defaultValue interface{}) (SmeType, error) {
	t, err := newSmeTypeByName(typeName, isOptional, hasDefaultValue, defaultValue)
	if err != nil {
		return nil, err
	}
	if isOptional {
		if hasDefaultValue {
			tp.optionalTypes.defaultValueTypes[typeName][defaultValue] = t
		} else {
			tp.optionalTypes.noDefaultValueTypes[typeName] = t
		}
	} else {
		if hasDefaultValue {
			tp.requiredTypes.defaultValueTypes[typeName][defaultValue.(string)] = t
		} else {
			tp.requiredTypes.noDefaultValueTypes[typeName] = t
		}
	}
	return t, nil
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
	result := map[string]SmeType{
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
		return typeName, nil
	}
	if IsParametricTypeName(typeName) {
		return unwrapParametricTypeName(packageName, typeName)
	}

	splittedTypeName := strings.Split(typeName, ".")
	if len(splittedTypeName) == 1 {
		typeName = packageName + "." + typeName
	}
	return typeName, nil
}

func unwrapParametricTypeName(packageName, typeName string) (string, error) {
	if strings.HasPrefix(typeName, "map") {
		var keyType, valueType string
		_, err := fmt.Sscanf(typeName, "map[%s,%s]", &keyType, &valueType)
		if err != nil {
			return "", err
		}
		unwrappedKey, err := unwrapTypeName(packageName, keyType)
		if err != nil {
			return "", err
		}
		unwrappedValue, err := unwrapTypeName(packageName, valueType)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("map[%s,%s]", unwrappedKey, unwrappedValue), nil
	}
	if strings.HasPrefix(typeName, "list") {
		var valueType string
		_, err := fmt.Sscanf(typeName, "list[%s]", &valueType)
		if err != nil {
			return "", err
		}
		unwrappedValue, err := unwrapTypeName(packageName, valueType)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("list[%s]", unwrappedValue), nil
	}
	return "", errNotParametricType
}

func TypeFromString(packageName, typeName string, isOptional bool, hasDefaultValue bool, defaultValue interface{}) (SmeType, error) {
	typeName, err := unwrapTypeName(packageName, typeName)
	if err != nil {
		return nil, err
	}
	t, err := typePool.getType(typeName, isOptional, hasDefaultValue, defaultValue)
	if err != nil {
		t, err = typePool.addType(typeName, isOptional, hasDefaultValue, defaultValue)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}
