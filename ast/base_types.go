package ast

import "strconv"

type SmeType interface {
	Id() uint32
	IsParametric() bool
	SizeOf() uint // Size of a field in the message converted to bytes
	SetOptionality()
	SetDefaultValue(string) error
	IsOptional() bool
	DefaultValue() (string, error)
}

type SmeBaseType struct {
	isOptional      bool
	hasDefaultValue bool
	defaultValue    string
}

func (bt *SmeBaseType) SetOptionality() {
	bt.isOptional = true
}

func (bt *SmeBaseType) IsOptional() bool {
	return bt.isOptional
}

func (bt *SmeBaseType) Id() uint32 {
	return ^uint32(0)
}

func (bt *SmeBaseType) IsParametric() bool {
	return false
}

func (bt *SmeBaseType) SizeOf() uint {
	return 0
}

func (bt *SmeBaseType) SetDefaultValue(v string) error {
	bt.hasDefaultValue = true
	bt.defaultValue = v
	return nil
}

func (bt *SmeBaseType) DefaultValue() (string, error) {
	if bt.hasDefaultValue {
		return bt.defaultValue, nil
	}
	return "", errNoDefaultValue
}

type SmeIntegerBase struct {
	SmeBaseType
}

func (ib *SmeIntegerBase) IsUnsigned() bool {
	return false
}

func (ib *SmeIntegerBase) SetDefaultValue(v string) error {
	if _, err := strconv.ParseUint(v, 10, int(ib.SizeOf())); err != nil {
		if ib.IsUnsigned() {
			if _, err := strconv.ParseInt(v, 10, int(ib.SizeOf())); err != nil {
				return err
			}
		}
	}
	ib.hasDefaultValue = true
	ib.defaultValue = v
	return nil
}

type SmeFloatingBase struct {
	SmeBaseType
}

func (fb *SmeFloatingBase) SetDefaultValue(v string) error {
	if _, err := strconv.ParseFloat(v, int(fb.SizeOf())); err != nil {
		return err
	}
	fb.hasDefaultValue = true
	fb.defaultValue = v
	return nil
}
