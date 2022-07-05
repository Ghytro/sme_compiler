package ast

import (
	"log"

	"github.com/Ghytro/sme/helpers"
)

type SmeInt8 struct {
	SmeIntegerBase
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
	return 1
}

func (i8 *SmeInt8) IsUnsigned() bool {
	return false
}

type SmeInt16 struct {
	SmeIntegerBase
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
	return 2
}

func (i16 *SmeInt16) IsUnsigned() bool {
	return false
}

type SmeInt32 struct {
	SmeIntegerBase
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
	return 4
}

func (i32 *SmeInt32) IsUnsigned() bool {
	return false
}

type SmeInt64 struct {
	SmeIntegerBase
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
	return 8
}

func (i64 *SmeInt64) IsUnsigned() bool {
	return false
}

type SmeUint8 struct {
	SmeIntegerBase
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

func (ui8 *SmeUint8) IsUnsigned() bool {
	return true
}

type SmeUint16 struct {
	SmeIntegerBase
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

func (ui16 *SmeUint16) IsUnsigned() bool {
	return true
}

type SmeUint32 struct {
	SmeIntegerBase
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

func (ui32 *SmeUint32) IsUnsigned() bool {
	return true
}

type SmeUint64 struct {
	SmeIntegerBase
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

func (ui64 *SmeUint64) IsUnsigned() bool {
	return true
}
