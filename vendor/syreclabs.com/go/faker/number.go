package faker

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"unsafe"
)

type FakeNumber interface {
	Number(digits int) string            // => "43202"
	NumberInt(digits int) int            // => 213
	NumberInt32(digits int) int32        // => 92938
	NumberInt64(digits int) int64        // => 1689541633257139096
	Decimal(precision, scale int) string // => "879420.60"
	Digit() string                       // => "7"
	Hexadecimal(digits int) string       // => "e7f3"
	Between(min, max int) string         // => "-47"
	Positive(max int) string             // => "3"
	Negative(min int) string             // => "-16"
}

type fakeNumber struct{}

func Number() FakeNumber {
	return fakeNumber{}
}

func (n fakeNumber) Number(digits int) string {
	if digits <= 0 {
		panic("invalid digits value")
	}
	dd := make([]string, digits, digits)
	for i := range dd {
		dd[i] = n.Digit()
	}
	return strings.Join(dd, "")
}

const (
	maxDigitsInt32 = 10
	maxDigitsInt64 = 19
)

func nonZeroDigit() string {
	return fmt.Sprintf("%d", localRand.Int31n(9)+1)
}

func numberPattern(digits, maxDigits int) string {
	switch digits {
	case 1:
		return nonZeroDigit()
	case maxDigits:
		return fmt.Sprintf(`1\d{%d}`, digits-1)
	default:
		return fmt.Sprintf(`%s\d{%d}`, nonZeroDigit(), digits-1)
	}
}

func (n fakeNumber) NumberInt32(digits int) int32 {
	if digits <= 0 || digits > maxDigitsInt32 {
		panic("invalid digits value")
	}
	pat := numberPattern(digits, maxDigitsInt32)
	num, err := Regexify(pat)
	if err != nil {
		panic(fmt.Sprintf("error regexifying %v: %v", pat, err))
	}
	res, err := strconv.ParseInt(num, 10, 32)
	if err != nil {
		panic(fmt.Sprintf("error parsing %v as int32: %v", num, err))
	}
	return int32(res)
}

func (n fakeNumber) NumberInt64(digits int) int64 {
	if digits <= 0 || digits > maxDigitsInt64 {
		panic("invalid digits value")
	}
	pat := numberPattern(digits, maxDigitsInt64)
	num, err := Regexify(pat)
	if err != nil {
		panic(fmt.Sprintf("error regexifying %v: %v", pat, err))
	}
	res, err := strconv.ParseInt(num, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("error parsing %v as int64: %v", num, err))
	}
	return int64(res)
}

func (n fakeNumber) NumberInt(digits int) int {
	if unsafe.Sizeof(int(0)) == unsafe.Sizeof(int64(0)) {
		return int(n.NumberInt64(digits))
	}
	return int(n.NumberInt32(digits))
}

func (n fakeNumber) Decimal(precision, scale int) string {
	if precision <= 0 || scale < 0 || precision < scale {
		panic("invalid precision or scale values")
	}
	s := n.Number(precision)
	ii := s[:precision-scale]
	if len(ii) == 0 {
		ii = "0"
	}
	ff := s[precision-scale:]
	if len(ff) == 0 {
		ff = "0"
	}

	return ii + "." + ff
}

func (n fakeNumber) Digit() string {
	return fmt.Sprintf("%d", localRand.Int31n(10))
}

func (n fakeNumber) Hexadecimal(digits int) string {
	if digits <= 0 {
		panic("invalid digits value")
	}
	bytes := make([]byte, (digits+1)/2)
	localRand.Read(bytes)
	return hex.EncodeToString(bytes)[:digits]
}

func (n fakeNumber) Between(min, max int) string {
	if min > max {
		panic("invalid range")
	}
	return fmt.Sprintf("%d", RandomInt(min, max))
}

func (n fakeNumber) Positive(max int) string {
	if max < 0 {
		panic("invalid max value")
	}
	return n.Between(0, max)
}

func (n fakeNumber) Negative(min int) string {
	if min > 0 {
		panic("invalid max value")
	}
	return n.Between(min, 0)
}
