package faker

import (
	"fmt"
	"strconv"
)

type FakeCode interface {
	Isbn10() string // => "048931033-8"
	Isbn13() string // => "391668236072-1"
	Ean13() string  // => "7742864258656"
	Ean8() string   // => "03079010"
	Rut() string    // => "14371602-3"
	Abn() string    // => "57914951376"
}

type fakeCode struct{}

func Code() FakeCode {
	return fakeCode{}
}

func (c fakeCode) Isbn10() string {
	val, err := Regexify(`\d{9}`)
	if err != nil {
		panic(err)
	}

	var sum int
	for i, v := range val {
		n, err := strconv.Atoi(string(v))
		if err != nil {
			panic(err)
		}
		sum += (i + 1) * n
	}
	rem := sum % 11

	if rem == 10 {
		return fmt.Sprintf("%s-X", val)
	}
	return fmt.Sprintf("%s-%d", val, rem)
}

func ean13() (ean string, checkDigit string) {
	ean, err := Regexify(`\d{12}`)
	if err != nil {
		panic(err)
	}

	var sum int
	for i, d := range ean {
		n, err := strconv.Atoi(string(d))
		if err != nil {
			panic(err)
		}
		if i%2 == 0 {
			sum += n
		} else {
			sum += n * 3
		}
	}
	rem := sum % 10
	checkDigit = strconv.Itoa((10 - rem) % 10)
	return
}

func (c fakeCode) Isbn13() string {
	ean, checkDigit := ean13()
	return fmt.Sprintf("%s-%s", ean, checkDigit)
}

func (c fakeCode) Ean13() string {
	ean, checkDigit := ean13()
	return fmt.Sprintf("%s%s", ean, checkDigit)
}

func (c fakeCode) Ean8() string {
	ean, err := Regexify(`\d{7}`)
	if err != nil {
		panic(err)
	}

	var sum int
	for i, d := range ean {
		n, err := strconv.Atoi(string(d))
		if err != nil {
			panic(err)
		}
		if i%2 == 0 {
			sum += n * 3
		} else {
			sum += n
		}
	}
	rem := sum % 10

	return fmt.Sprintf("%s%d", ean, (10-rem)%10)
}

var rutFactors = [...]int{3, 2, 7, 6, 5, 4, 3, 2}

func (c fakeCode) Rut() string {
	rut, err := Regexify(`\d{8}`)
	if err != nil {
		panic(err)
	}

	var sum int
	for i, d := range rut {
		n, err := strconv.Atoi(string(d))
		if err != nil {
			panic(err)
		}
		sum += n * rutFactors[i]
	}
	rem := 11 - (sum % 11)

	var checkDigit string
	switch rem {
	case 11:
		checkDigit = "0"
	case 10:
		checkDigit = "K"
	default:
		checkDigit = strconv.Itoa(rem)
	}

	return fmt.Sprintf("%s-%s", rut, checkDigit)
}

var abnFactors = [...]int{3, 5, 7, 9, 11, 13, 15, 17, 19}

func (c fakeCode) Abn() string {
	acn, err := Regexify(`\d{9}`)
	if err != nil {
		panic(err)
	}

	var sum int
	for i, d := range acn {
		n, err := strconv.Atoi(string(d))
		if err != nil {
			panic(err)
		}
		sum += n * abnFactors[i]
	}
	rem := (sum/89+1)*89 - sum
	checkDigits := rem + 10

	return fmt.Sprintf("%d%s", checkDigits, acn)
}
