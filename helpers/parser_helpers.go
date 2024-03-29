package helpers

import (
	"crypto"
	_ "crypto/md5"
	"errors"
	"fmt"
	"hash/fnv"
	"reflect"
	"regexp"
)

func Uint32Sum32Hash(ints ...uint32) (uint32, error) {
	if len(ints) == 0 {
		return 0, errors.New("no ints passed to make hash of")
	}
	byteArr := make([]byte, len(ints)*32)
	for i := range byteArr {
		byteArr[i] = byte((ints[i/4] >> (i % 4 * 8)) & 0xFF)
	}
	hash := fnv.New32a()
	if _, err := hash.Write(byteArr); err != nil {
		return 0, err
	}
	return hash.Sum32(), nil
}

func HashValuesUint32(values ...interface{}) (uint32, error) {
	digester := crypto.MD5.New()
	for _, v := range values {
		fmt.Fprint(digester, reflect.TypeOf(v))
		fmt.Fprint(digester, v)
	}
	hash := fnv.New32a()
	if _, err := hash.Write(digester.Sum(nil)); err != nil {
		return 0, err
	}
	return hash.Sum32(), nil
}

func EqualsAny(c byte, chars ...byte) bool {
	for _, _c := range chars {
		if _c == c {
			return true
		}
	}
	return false
}

func IsAllowedStructChar(c byte) bool {
	return c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c == '_'
}

func MatchString(regex string, s string) (bool, error) {
	re, err := regexp.Compile(regex)
	if err != nil {
		return false, err
	}
	idx := re.FindStringIndex(s)
	if idx == nil {
		return false, nil
	}
	return idx[0] == 0 && idx[1] == len(s), nil
}
