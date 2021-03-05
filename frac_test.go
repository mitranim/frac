package frac

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"testing"
)

const (
	benchSrcInt   = `-01230`
	benchSrcFrac  = `-01230.04560`
	benchNumInt   = -1234
	benchNumFrac  = -123456
	benchNumFloat = -123.456
)

var (
	maxInt64 = strconv.FormatInt(math.MaxInt64, 10)
	minInt64 = strconv.FormatInt(math.MinInt64, 10)
)

func BenchmarkParseDecInt(b *testing.B) {
	for range counter(b.N) {
		_, err := ParseDec(benchSrcInt, 4)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkStrconvParseInt(b *testing.B) {
	for range counter(b.N) {
		_, err := strconv.ParseInt(benchSrcInt, 10, 64)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseDecFrac(b *testing.B) {
	for range counter(b.N) {
		_, err := ParseDec(benchSrcFrac, 4)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkStrconvParseFloat(b *testing.B) {
	for range counter(b.N) {
		_, err := strconv.ParseFloat(benchSrcFrac, 64)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkFormatDecInt(b *testing.B) {
	for range counter(b.N) {
		_, err := FormatDec(benchNumInt, 4)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkStrconvFormatInt(b *testing.B) {
	for range counter(b.N) {
		_ = strconv.FormatInt(benchNumInt, 10)
	}
}

func BenchmarkFormatDecFrac(b *testing.B) {
	for range counter(b.N) {
		_, err := FormatDec(benchNumFrac, 4)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkStrconvFormatFloat(b *testing.B) {
	for range counter(b.N) {
		_ = strconv.FormatFloat(benchNumFloat, 'f', -1, 64)
	}
}

func TestParseRadix(*testing.T) {
	const frac = 0
	const val = 1

	for radix := uint(2); radix < 36; radix++ {
		testParse(`1`, frac, radix, val)
	}

	testParseErr(`0`, `unsupported radix`, 0, frac)
	testParseErr(`0`, `unsupported radix`, 1, frac)
	testParseErr(`0`, `unsupported radix`, 37, frac)
}

func TestParseDec(t *testing.T) {
	t.Run(`invalid`, func(*testing.T) {
		testParseErrDec(``, `empty input`, 0, 1, 2)
		testParseErrDec(` `, `non-digit character`, 0, 1, 2)
		testParseErrDec(` 12`, `non-digit character`, 0, 1, 2)
		testParseErrDec(`12 `, `non-digit character`, 0, 1, 2)
		testParseErrDec(`12.`, `unexpected end of input`, 0, 1, 2)
		testParseErrDec(`.12`, `non-digit character`, 0, 1, 2)
		testParseErrDec(`+`, `unexpected end of input`, 0, 1, 2)
		testParseErrDec(`-`, `unexpected end of input`, 0, 1, 2)
		testParseErrDec(`+.12`, `non-digit character`, 0, 1, 2)
		testParseErrDec(`-.12`, `non-digit character`, 0, 1, 2)
		testParseErrDec(`--12`, `non-digit character`, 0, 1, 2)
		testParseErrDec(`++12`, `non-digit character`, 0, 1, 2)
		testParseErrDec(`-+12`, `non-digit character`, 0, 1, 2)
		testParseErrDec(`+-12`, `non-digit character`, 0, 1, 2)
	})

	t.Run(`int`, func(*testing.T) {
		testParseDec(`0`, 0, 0)
		testParseDec(`0`, 1, 0)
		testParseDec(`0`, 2, 0)
		testParseDec(`01230`, 2, 1_230_00)
		testParseDec(`+01230`, 2, 1_230_00)
		testParseDec(`-01230`, 2, -1_230_00)
		testParseDec(maxInt64, 0, math.MaxInt64)
		testParseDec(minInt64, 0, math.MinInt64)

		testParseErrDec(maxInt64+`0`, `overflow`, 0, 1, 2, 3)
		testParseErrDec(minInt64+`0`, `underflow`, 0, 1, 2, 3)
	})

	t.Run(`frac`, func(*testing.T) {
		testParseErrDec(`01230.04560`, `exponent exceeds`, 0, 1, 2, 3)
		testParseErrDec(`+01230.04560`, `exponent exceeds`, 0, 1, 2, 3)
		testParseErrDec(`-01230.04560`, `exponent exceeds`, 0, 1, 2, 3)

		testParseDec(`01230.04560`, 4, 1_230_0456)
		testParseDec(`+01230.04560`, 4, 1_230_0456)
		testParseDec(`-01230.04560`, 4, -1_230_0456)

		testParseDec(maxInt64+`.000000000000000`, 0, math.MaxInt64)
		testParseDec(minInt64+`.000000000000000`, 0, math.MinInt64)
		testParseDec(`0.0000000000000000000000000000000000000000000009223372036854775807`, 64, math.MaxInt64)
		testParseDec(`-0.0000000000000000000000000000000000000000000009223372036854775808`, 64, math.MinInt64)
		testParseDec(`0.000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000009223372036854775807`, 111, math.MaxInt64)
		testParseDec(`-0.000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000009223372036854775808`, 111, math.MinInt64)

		testParseErrDec(maxInt64+`.1`, `overflow`, 1, 2, 3)
		testParseErrDec(minInt64+`.1`, `underflow`, 1, 2, 3)
	})
}

func TestParseBin(t *testing.T) {
	t.Run(`invalid`, func(*testing.T) {
		testParseErrBin(`2`, `non-digit character`, 0)
	})

	t.Run(`int`, func(*testing.T) {
		testParseBin(`0`, 0, 0)
		testParseBin(`0`, 1, 0)
		testParseBin(`0`, 2, 0)
		testParseBin(`1`, 2, 0b1_00)
		testParseBin(`01`, 2, 0b01_00)
		testParseBin(`10`, 2, 0b10_00)
		testParseBin(`010101`, 2, 0b010101_00)
		testParseBin(`+010101`, 2, 0b010101_00)
		testParseBin(`-1010101`, 2, -0b1010101_00)
	})

	t.Run(`frac`, func(*testing.T) {
		testParseErrBin(`010101.01010`, `exponent exceeds`, 0, 1, 2, 3)
		testParseErrBin(`+010101.01010`, `exponent exceeds`, 0, 1, 2, 3)
		testParseErrBin(`-010101.01010`, `exponent exceeds`, 0, 1, 2, 3)

		testParseBin(`010101.01010`, 4, 0b_010_101_0101)
		testParseBin(`+010101.01010`, 4, 0b_010_101_0101)
		testParseBin(`-010101.01010`, 4, -0b_010_101_0101)
	})
}

func TestParseHex(t *testing.T) {
	t.Run(`invalid`, func(*testing.T) {
		testParseErrHex(`g`, `non-digit character`, 0)
	})

	t.Run(`int`, func(*testing.T) {
		testParseHex(`0`, 0, 0)
		testParseHex(`0`, 1, 0)
		testParseHex(`0`, 2, 0)
		testParseHex(`1`, 2, 0x1_00)
		testParseHex(`f`, 2, 0xf_00)
		testParseHex(`0f`, 2, 0x0f_00)
		testParseHex(`f0`, 2, 0xf0_00)
		testParseHex(`0f0f0f`, 2, 0x0f0f0f_00)
		testParseHex(`+0f0f0f`, 2, 0x0f0f0f_00)
		testParseHex(`-f0f0f0f`, 2, -0xf0f0f0f_00)
	})

	t.Run(`frac`, func(*testing.T) {
		testParseErrHex(`0f0f0f.0f0f0`, `exponent exceeds`, 0, 1, 2, 3)
		testParseErrHex(`+0f0f0f.0f0f0`, `exponent exceeds`, 0, 1, 2, 3)
		testParseErrHex(`-0f0f0f.0f0f0`, `exponent exceeds`, 0, 1, 2, 3)

		testParseHex(`0f0f0f.0f0f0`, 4, 0x_0f0_f0f_0f0f)
		testParseHex(`+0f0f0f.0f0f0`, 4, 0x_0f0_f0f_0f0f)
		testParseHex(`-0f0f0f.0f0f0`, 4, -0x_0f0_f0f_0f0f)
	})
}

func TestFormatRadix(*testing.T) {
	const frac = 0
	const val = 1

	for radix := uint(2); radix < 36; radix++ {
		testFormat(val, frac, radix, `1`)
	}

	testFormatErr(val, 0, frac, `unsupported radix`)
	testFormatErr(val, 1, frac, `unsupported radix`)
	testFormatErr(val, 37, frac, `unsupported radix`)
}

func TestFormatDec(*testing.T) {
	testFormatDec(0, 0, `0`)
	testFormatDec(0, 1, `0`)
	testFormatDec(0, 2, `0`)
	testFormatDec(0, 3, `0`)
	testFormatDec(1, 0, `1`)
	testFormatDec(1, 1, `0.1`)
	testFormatDec(1, 2, `0.01`)
	testFormatDec(1, 3, `0.001`)
	testFormatDec(123456, 0, `123456`)
	testFormatDec(123456, 1, `12345.6`)
	testFormatDec(123456, 2, `1234.56`)
	testFormatDec(123456, 3, `123.456`)
	testFormatDec(123456, 4, `12.3456`)
	testFormatDec(123456, 5, `1.23456`)
	testFormatDec(123456, 6, `0.123456`)
	testFormatDec(123456, 7, `0.0123456`)
	testFormatDec(123000, 0, `123000`)
	testFormatDec(123000, 1, `12300`)
	testFormatDec(123000, 2, `1230`)
	testFormatDec(123000, 3, `123`)
	testFormatDec(123000, 4, `12.3`)
	testFormatDec(123000, 5, `1.23`)
	testFormatDec(123000, 6, `0.123`)
	testFormatDec(123000, 7, `0.0123`)
	testFormatDec(math.MaxInt64, 0, `9223372036854775807`)
	testFormatDec(math.MinInt64, 0, `-9223372036854775808`)
	testFormatDec(math.MaxInt64, 64, `0.0000000000000000000000000000000000000000000009223372036854775807`)
	testFormatDec(math.MinInt64, 64, `-0.0000000000000000000000000000000000000000000009223372036854775808`)
	testFormatErrDec(1, 65, `exceeds limit`)
	testFormatErrDec(-1, 65, `exceeds limit`)
}

func TestFormatBin(*testing.T) {
	testFormatBin(0b0, 0, `0`)
	testFormatBin(0b0, 1, `0`)
	testFormatBin(0b0, 2, `0`)
	testFormatBin(0b0, 3, `0`)
	testFormatBin(0b1, 0, `1`)
	testFormatBin(0b1, 1, `0.1`)
	testFormatBin(0b1, 2, `0.01`)
	testFormatBin(0b1, 3, `0.001`)
	testFormatBin(0b1010101, 0, `1010101`)
	testFormatBin(0b1010101, 1, `101010.1`)
	testFormatBin(0b1010101, 2, `10101.01`)
	testFormatBin(0b1010101, 3, `1010.101`)
	testFormatBin(0b1010101, 4, `101.0101`)
	testFormatBin(0b1010101, 5, `10.10101`)
	testFormatBin(0b1010101, 6, `1.010101`)
	testFormatBin(0b1010101, 7, `0.1010101`)
	testFormatBin(0b1010101, 8, `0.01010101`)
	testFormatBin(0b111000, 0, `111000`)
	testFormatBin(0b111000, 1, `11100`)
	testFormatBin(0b111000, 2, `1110`)
	testFormatBin(0b111000, 3, `111`)
	testFormatBin(0b111000, 4, `11.1`)
	testFormatBin(0b111000, 5, `1.11`)
	testFormatBin(0b111000, 6, `0.111`)
	testFormatBin(0b111000, 7, `0.0111`)
	testFormatBin(math.MaxInt64, 0, `111111111111111111111111111111111111111111111111111111111111111`)
	testFormatBin(math.MinInt64, 0, `-1000000000000000000000000000000000000000000000000000000000000000`)
	testFormatBin(math.MaxInt64, 64, `0.0111111111111111111111111111111111111111111111111111111111111111`)
	testFormatBin(math.MinInt64, 64, `-0.1`)
	testFormatErrBin(1, 65, `exceeds limit`)
	testFormatErrBin(-1, 65, `exceeds limit`)
}

func TestFormatHex(*testing.T) {
	testFormatHex(0x0, 0, `0`)
	testFormatHex(0x0, 1, `0`)
	testFormatHex(0x0, 2, `0`)
	testFormatHex(0x0, 3, `0`)
	testFormatHex(0xf, 0, `f`)
	testFormatHex(0xf, 1, `0.f`)
	testFormatHex(0xf, 2, `0.0f`)
	testFormatHex(0xf, 3, `0.00f`)
	testFormatHex(0xf0f0f0f, 0, `f0f0f0f`)
	testFormatHex(0xf0f0f0f, 1, `f0f0f0.f`)
	testFormatHex(0xf0f0f0f, 2, `f0f0f.0f`)
	testFormatHex(0xf0f0f0f, 3, `f0f0.f0f`)
	testFormatHex(0xf0f0f0f, 4, `f0f.0f0f`)
	testFormatHex(0xf0f0f0f, 5, `f0.f0f0f`)
	testFormatHex(0xf0f0f0f, 6, `f.0f0f0f`)
	testFormatHex(0xf0f0f0f, 7, `0.f0f0f0f`)
	testFormatHex(0xf0f0f0f, 8, `0.0f0f0f0f`)
	testFormatHex(0xfff000, 0, `fff000`)
	testFormatHex(0xfff000, 1, `fff00`)
	testFormatHex(0xfff000, 2, `fff0`)
	testFormatHex(0xfff000, 3, `fff`)
	testFormatHex(0xfff000, 4, `ff.f`)
	testFormatHex(0xfff000, 5, `f.ff`)
	testFormatHex(0xfff000, 6, `0.fff`)
	testFormatHex(0xfff000, 7, `0.0fff`)
	testFormatHex(math.MaxInt64, 0, `7fffffffffffffff`)
	testFormatHex(math.MinInt64, 0, `-8000000000000000`)
	testFormatHex(math.MaxInt64, 64, `0.0000000000000000000000000000000000000000000000007fffffffffffffff`)
	testFormatHex(math.MinInt64, 64, `-0.0000000000000000000000000000000000000000000000008`)
	testFormatErrHex(1, 65, `exceeds limit`)
	testFormatErrHex(-1, 65, `exceeds limit`)
}

func testParseBin(src string, frac uint, exp int64) { testParse(src, frac, 2, exp) }
func testParseDec(src string, frac uint, exp int64) { testParse(src, frac, 10, exp) }
func testParseHex(src string, frac uint, exp int64) { testParse(src, frac, 16, exp) }

func testParseErrBin(src string, msg string, fracs ...uint) {
	testParseErr(src, msg, 2, fracs...)
}

func testParseErrDec(src string, msg string, fracs ...uint) {
	testParseErr(src, msg, 10, fracs...)
}

func testParseErrHex(src string, msg string, fracs ...uint) {
	testParseErr(src, msg, 16, fracs...)
}

func testFormatBin(num int64, frac uint, exp string) { testFormat(num, frac, 2, exp) }
func testFormatDec(num int64, frac uint, exp string) { testFormat(num, frac, 10, exp) }
func testFormatHex(num int64, frac uint, exp string) { testFormat(num, frac, 16, exp) }

func testFormatErrBin(num int64, frac uint, msg string) { testFormatErr(num, 2, frac, msg) }
func testFormatErrDec(num int64, frac uint, msg string) { testFormatErr(num, 10, frac, msg) }
func testFormatErrHex(num int64, frac uint, msg string) { testFormatErr(num, 16, frac, msg) }

func testParse(src string, frac uint, radix uint, exp int64) {
	act, err := Parse(src, frac, radix)
	if err != nil {
		panic(fmt.Errorf(`failed to parse %q (frac %v, radix %v): %+v`, src, frac, radix, err))
	}
	if exp != act {
		panic(fmt.Errorf(`expected to parse %q (frac %v, radix %v) into %v, got %v`, src, frac, radix, exp, act))
	}
}

func testParseErr(src string, msg string, radix uint, fracs ...uint) {
	for _, frac := range fracs {
		res, err := Parse(src, frac, radix)
		if err == nil {
			panic(fmt.Errorf(`expected parsing %q (frac %v, radix %v) to fail; instead got %v`, src, frac, radix, res))
		}
		if !strings.Contains(err.Error(), msg) {
			panic(fmt.Errorf(`expected error from parsing %q (frac %v, radix %v) to contain %q, got %q`, src, frac, radix, msg, err))
		}
	}
}

func testFormat(num int64, frac uint, radix uint, exp string) {
	act, err := Format(num, frac, radix)
	if err != nil {
		panic(fmt.Errorf(`failed to format %v (frac %v, radix %v): %+v`, num, frac, radix, err))
	}
	if exp != act {
		panic(fmt.Errorf(`expected to format %v (frac %v, radix %v) into %q, got %q`, num, frac, radix, exp, act))
	}
}

func testFormatErr(num int64, radix uint, frac uint, msg string) {
	res, err := Format(num, frac, radix)
	if err == nil {
		panic(fmt.Errorf(`expected parsing %q (frac %v, radix %v) to fail; instead got %q`, num, frac, radix, res))
	}
	if !strings.Contains(err.Error(), msg) {
		panic(fmt.Errorf(`expected error from parsing %q (frac %v, radix %v) to contain %q, got %q`, num, frac, radix, msg, err))
	}
}

func counter(count int) []struct{} { return make([]struct{}, count) }
