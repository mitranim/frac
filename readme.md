## Overview

Missing feature of the Go standard library: parsing and formatting integers as fractional numeric strings, without any rounding or bignums, by using a fixed fraction size. Supports arbitrary radixes from 2 to 36.

For example:

```
"123"     <- frac 2, radix 10 -> 123_00
"123.45"  <- frac 2, radix 10 -> 123_45
"123.456" <- frac 2, radix 10 -> <error>
"123"     <- frac 3, radix 10 -> 123_000
"123.456" <- frac 3, radix 10 -> 123_456
```

Performance on 64-bit machines is comparable to `strconv` and shouldn't be your bottleneck.

See API docs at https://pkg.go.dev/github.com/mitranim/frac.

## Why

* You use integers for money.
* You deal with external APIs that use decimal strings for money.
* You want to avoid rounding errors.
* You don't want to deal with "big decimal" libraries.

Then `frac` is for you!

## Usage

Basic usage:

```golang
import "github.com/mitranim/frac"

func main() {
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

func assert(ok bool) {if !ok {panic("unreachable")}}
```

Implementing a monetary type:

```golang
import "github.com/mitranim/frac"

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
```

The resulting type `Cents` is an integer, but when decoding and encoding text, it's represented as a fractional with 2 decimal points.

## Known Limitations

* The code is too assembly-like. Kinda like the standard library.

* No special support for unsigned integers.

* When formatting, fractional precision is limited to `64`. (Imagine allocating gigabytes of memory for `0.0...01`.)

## License

https://unlicense.org

## Misc

I'm receptive to suggestions. If this library _almost_ satisfies you but needs changes, open an issue or chat me up. Contacts: https://mitranim.com/#contacts
