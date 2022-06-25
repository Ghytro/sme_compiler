package ast

import (
	"log"

	"github.com/Ghytro/sme/helpers"
)

type SmeFloat struct {
	SmeFloatingBase
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
	return 4
}

type SmeDouble struct {
	SmeFloatingBase
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
	return 8
}
