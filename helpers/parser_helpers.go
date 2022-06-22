package helpers

import (
	"crypto"
	_ "crypto/md5"
	"errors"
	"fmt"
	"hash/fnv"
	"reflect"
	"strings"
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

func RemoveExtraSpacesAroundStrings(modified string, strs ...string) string {
	for _, s := range strs {
		modified = strings.ReplaceAll(
			strings.ReplaceAll(
				modified,
				" "+s,
				s,
			),
			s+" ",
			s,
		)
	}
	return modified
}
