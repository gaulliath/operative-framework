package faker

import (
	"regexp"
	"strconv"
	"strings"
)

type FakeFinance interface {
	// CreditCard returns a valid (with valid check digit) card number of one of the given types.
	// If no types are passed, all types in CC_TYPES are used.
	CreditCard(types ...string) string // => "5019-8413-2066-5594"
}

type fakeFinance struct{}

func Finance() FakeFinance {
	return fakeFinance{}
}

// Known credit card types.
const (
	CC_VISA               = "visa"
	CC_MASTERCARD         = "mastercard"
	CC_AMERICAN_EXPRESS   = "american_express"
	CC_DINERS_CLUB        = "diners_club"
	CC_DISCOVER           = "discover"
	CC_MAESTRO            = "maestro"
	CC_SWITCH             = "switch"
	CC_SOLO               = "solo"
	CC_FORBRUGSFORENINGEN = "forbrugsforeningen"
	CC_DANKORT            = "dankort"
	CC_LASER              = "laser"
)

// CC_TYPES holds a list of known credit card types.
var CC_TYPES = []string{
	CC_VISA,
	CC_MASTERCARD,
	CC_AMERICAN_EXPRESS,
	CC_DINERS_CLUB,
	CC_DISCOVER,
	CC_MAESTRO,
	CC_SWITCH,
	CC_SOLO,
	CC_FORBRUGSFORENINGEN,
	CC_DANKORT,
	CC_LASER,
}

var luhnFactors = [...]int{0, 2, 4, 6, 8, 1, 3, 5, 7, 9}

func luhnCheckDigit(s string) string {
	// assuming check digit will be appended as last digit
	odd := (len(s) + 1) & 1

	var sum int
	for i, c := range s {
		if c < '0' || c > '9' {
			panic("invalid number sequence: " + s)
		}
		if i&1 == odd {
			sum += luhnFactors[c-'0']
		} else {
			sum += int(c - '0')
		}
	}

	return strconv.Itoa((10 - (sum % 10)) % 10)
}

var rxDigits = regexp.MustCompile(`\D`)

func (f fakeFinance) CreditCard(types ...string) string {
	var t string
	if len(types) > 0 {
		t = RandomChoice(types)
	} else {
		t = RandomChoice(CC_TYPES)
	}

	tmpl := Numerify(Fetch("credit_card." + t))

	num := rxDigits.ReplaceAllString(tmpl, "")
	luhn := luhnCheckDigit(num)
	res := strings.Replace(tmpl, "L", luhn, 1)

	return res
}
