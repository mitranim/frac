package frac_test

import "github.com/mitranim/frac"

func ExampleParseDec() {
	num, err := frac.ParseDec(`-123`, 2)
	assert(err == nil && num == -123_00)

	num, err = frac.ParseDec(`-123.00`, 2)
	assert(err == nil && num == -123_00)

	num, err = frac.ParseDec(`-123.45`, 2)
	assert(err == nil && num == -123_45)

	// Exponent exceeds allotted precision. Conversion is impossible.
	num, err = frac.ParseDec(`-123.456`, 2)
	assert(err != nil && num == 0)
}

func assert(ok bool) {
	if !ok {
		panic("unreachable")
	}
}

func ExampleCents() {
	var _ Cents
}

type Cents int64

func (self *Cents) UnmarshalText(input []byte) error {
	num, err := frac.UnmarshalDec(input, 2)
	if err != nil {
		return err
	}
	*self = Cents(num)
	return nil
}

func (self Cents) MarshalText() ([]byte, error) {
	return frac.AppendDec(nil, int64(self), 2)
}
