/*
Package faker is a library for generating fake data such as names, addresses, and phone numbers.

It is a (mostly) API-compatible port of Ruby Faker gem (https://github.com/stympy/faker) to Go.
*/
package faker // import "syreclabs.com/go/faker"

import (
	"fmt"
	"math/rand"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"syreclabs.com/go/faker/locales"
)

const (
	digits           = "0123456789"
	uLetters         = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	dLetters         = "abcdefghijklmnopqrstuvwxyz"
	letters          = uLetters + dLetters
	digitsAndLetters = digits + letters
)

var (
	aryDigits   = strings.Split(digits, "")
	aryULetters = strings.Split(uLetters, "")
	aryDLetters = strings.Split(dLetters, "")
	aryLetters  = strings.Split(letters, "")
)

// Generate locales
//go:generate go run cmd/generate.go yaml/de-AT.yml locales/de-at.go
//go:generate go run cmd/generate.go yaml/de-CH.yml locales/de-ch.go
//go:generate go run cmd/generate.go yaml/de.yml locales/de.go
//go:generate go run cmd/generate.go yaml/en-au-ocker.yml locales/en-au-ocker.go
//go:generate go run cmd/generate.go yaml/en-AU.yml locales/en-au.go
//go:generate go run cmd/generate.go yaml/en-BORK.yml locales/en-bork.go
//go:generate go run cmd/generate.go yaml/en-CA.yml locales/en-ca.go
//go:generate go run cmd/generate.go yaml/en-GB.yml locales/en-gb.go
//go:generate go run cmd/generate.go yaml/en-IND.yml locales/en-ind.go
//go:generate go run cmd/generate.go yaml/en-NEP.yml locales/en-nep.go
//go:generate go run cmd/generate.go yaml/en-US.yml locales/en-us.go
//go:generate go run cmd/generate.go yaml/en.yml locales/en.go
//go:generate go run cmd/generate.go yaml/es.yml locales/es.go
//go:generate go run cmd/generate.go yaml/fa.yml locales/fa.go
//go:generate go run cmd/generate.go yaml/fr.yml locales/fr.go
//go:generate go run cmd/generate.go yaml/it.yml locales/it.go
//go:generate go run cmd/generate.go yaml/ja.yml locales/ja.go
//go:generate go run cmd/generate.go yaml/ko.yml locales/ko.go
//go:generate go run cmd/generate.go yaml/nb-NO.yml locales/nb-no.go
//go:generate go run cmd/generate.go yaml/nl.yml locales/nl.go
//go:generate go run cmd/generate.go yaml/pl.yml locales/pl.go
//go:generate go run cmd/generate.go yaml/pt-BR.yml locales/pt-br.go
//go:generate go run cmd/generate.go yaml/ru.yml locales/ru.go
//go:generate go run cmd/generate.go yaml/sk.yml locales/sk.go
//go:generate go run cmd/generate.go yaml/sv.yml locales/sv.go
//go:generate go run cmd/generate.go yaml/vi.yml locales/vi.go
//go:generate go run cmd/generate.go yaml/zh-CN.yml locales/zh-cn.go
//go:generate go run cmd/generate.go yaml/zh-TW.yml locales/zh-tw.go

// Locale holds the default locale.
var Locale = locales.En

// RandomInt returns random int in [min, max] range.
func RandomInt(min, max int) int {
	if max <= min {
		// degenerate case, return min
		return min
	}
	return min + localRand.Intn(max-min+1)
}

// RandomInt64 returns random int64 in [min, max] range.
func RandomInt64(min, max int64) int64 {
	if max <= min {
		// degenerate case, return min
		return min
	}
	return min + localRand.Int63n(max-min+1)
}

// RandomString returns a random alphanumeric string with length n.
func RandomString(n int) string {
	bytes := make([]byte, n)
	localRand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = digitsAndLetters[b%byte(len(digitsAndLetters))]
	}
	return string(bytes)
}

// RandomRepeat returns a new string consisting of random number of copies of the string s.
func RandomRepeat(s string, min, max int) string {
	return strings.Repeat(s, RandomInt(min, max))
}

// RandomChoice returns random string from slice of strings.
func RandomChoice(ss []string) string {
	return ss[localRand.Intn(len(ss))]
}

func includesString(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}

// Numerify replaces pattern like '##-###' with randomly generated digits.
func Numerify(s string) string {
	first := true
	for _, sm := range regexp.MustCompile(`#`).FindAllStringSubmatch(s, -1) {
		if first {
			// make sure result does not start with zero
			s = strings.Replace(s, sm[0], RandomChoice(aryDigits[1:]), 1)
			first = false
		} else {
			s = strings.Replace(s, sm[0], RandomChoice(aryDigits), 1)
		}
	}
	return s
}

// Letterify replaces pattern like '??? ??' with randomly generated uppercase letters
func Letterify(s string) string {
	for _, sm := range regexp.MustCompile(`\?`).FindAllStringSubmatch(s, -1) {
		s = strings.Replace(s, sm[0], RandomChoice(aryULetters), 1)
	}
	return s
}

// NumerifyAndLetterify both numerifies and letterifies s.
func NumerifyAndLetterify(s string) string {
	return Letterify(Numerify(s))
}

// Regexify attempts to generate a string that would match given regular expression.
// It does not handle ., *, unbounded ranges such as {1,},
// extensions such as (?=), character classes, some abbreviations for character classes,
// and nested parentheses.
func Regexify(s string) (string, error) {
	// ditch the anchors
	res := regexp.MustCompile(`^\/?\^?`).ReplaceAllString(s, "")
	res = regexp.MustCompile(`\$?\/?$`).ReplaceAllString(res, "")

	// all {2} become {2,2} and ? become {0,1}
	res = regexp.MustCompile(`\{(\d+)\}`).ReplaceAllString(res, `{$1,$1}`)
	res = regexp.MustCompile(`\?`).ReplaceAllString(res, `{0,1}`)

	// [12]{1,2} becomes [12] or [12][12]
	for _, sm := range regexp.MustCompile(`(\[[^\]]+\])\{(\d+),(\d+)\}`).FindAllStringSubmatch(res, -1) {
		min, _ := strconv.Atoi(sm[2])
		max, _ := strconv.Atoi(sm[3])
		res = strings.Replace(res, sm[0], RandomRepeat(sm[1], min, max), 1)
	}

	// (12|34){1,2} becomes (12|34) or (12|34)(12|34)
	for _, sm := range regexp.MustCompile(`(\([^\)]+\))\{(\d+),(\d+)\}`).FindAllStringSubmatch(res, -1) {
		min, _ := strconv.Atoi(sm[2])
		max, _ := strconv.Atoi(sm[3])
		res = strings.Replace(res, sm[0], RandomRepeat(sm[1], min, max), 1)
	}

	// A{1,2} becomes A or AA or \d{3} becomes \d\d\d
	for _, sm := range regexp.MustCompile(`(\\?.)\{(\d+),(\d+)\}`).FindAllStringSubmatch(res, -1) {
		min, _ := strconv.Atoi(sm[2])
		max, _ := strconv.Atoi(sm[3])
		res = strings.Replace(res, sm[0], RandomRepeat(sm[1], min, max), 1)
	}

	// (this|that) becomes 'this' or 'that'
	for _, sm := range regexp.MustCompile(`\((.*?)\)`).FindAllStringSubmatch(res, -1) {
		res = strings.Replace(res, sm[0], RandomChoice(strings.Split(sm[1], "|")), 1)
	}

	// all A-Z inside of [] become C (or X, or whatever)
	for _, sm := range regexp.MustCompile(`\[([^\]]+)\]`).FindAllStringSubmatch(res, -1) {
		cls := sm[1]
		// find and replace all ranges in character class cls
		for _, subsm := range regexp.MustCompile(`(\w\-\w)`).FindAllStringSubmatch(cls, -1) {
			rng := strings.Split(subsm[1], "-")
			repl := string(RandomInt(int(rng[0][0]), int(rng[1][0])))
			cls = strings.Replace(cls, subsm[0], repl, 1)
		}
		res = strings.Replace(res, sm[1], cls, 1)
	}

	// all [ABC] become B (or A or C)
	for _, sm := range regexp.MustCompile(`\[([^\]]+)\]`).FindAllStringSubmatch(res, -1) {
		res = strings.Replace(res, sm[0], RandomChoice(strings.Split(sm[1], "")), 1)
	}

	// all \d become random digits
	res = regexp.MustCompile(`\\d`).ReplaceAllStringFunc(res, func(s string) string {
		return RandomChoice(aryDigits)
	})

	// all \w become random letters
	res = regexp.MustCompile(`\\d`).ReplaceAllStringFunc(res, func(s string) string {
		return RandomChoice(aryLetters)
	})

	return res, nil
}

func localeValueAt(path string, locale map[string]interface{}) (interface{}, bool) {
	var val interface{} = locale
	for _, key := range strings.Split(path, ".") {
		v, ok := val.(map[string]interface{})
		if !ok {
			// all nodes are expected to have map[string]interface{} type
			panic(fmt.Sprintf("%v: invalid value type %v", path, reflect.TypeOf(val)))
		}
		val, ok = v[key]
		if !ok {
			// given path does not exists in given locale
			return nil, false
		}
	}
	return val, true
}

func valueAt(path string) interface{} {
	val, ok := localeValueAt(path, Locale)
	if !ok {
		// path does not exist in given locale, fallback to En
		val, ok = localeValueAt(path, locales.En)
		if !ok {
			// not in En either, give up
			panic(fmt.Sprintf("%v: invalid path", path))
		}
	}
	return val
}

// Fetch returns a value at given key path in default locale. If key path holds an array,
// it returns random array element. If value looks like a regex, it attempts to regexify it.
func Fetch(path string) string {
	var res string

	switch val := valueAt(path).(type) {
	case [][]string:
		// slice of string slices - select random element and join
		choices := make([]string, len(val))
		for i, slice := range val {
			choices[i] = RandomChoice(slice)
		}
		res = strings.Join(choices, " ")
	case []string:
		// slice of strings - select random element
		res = RandomChoice(val)
	case string:
		// plain string
		res = val
	default:
		// not supported
		panic(fmt.Sprintf("%v: invalid value type %v", path, reflect.TypeOf(val)))
	}

	// recursively substitute #{...} value references
	for _, sm := range regexp.MustCompile(`#\{([A-Za-z]+\.[^\}]+)\}`).FindAllStringSubmatch(res, -1) {
		path := strings.ToLower(sm[1])
		res = strings.Replace(res, sm[0], Fetch(path), 1)
	}

	// if res looks like regex, regexify
	if strings.HasPrefix(res, "/") && strings.HasSuffix(res, "/") {
		res, err := Regexify(res)
		if err != nil {
			panic(fmt.Sprintf("failed to regexify %v: %v", res, err))
		}
		return res
	}

	return res
}

// from https://github.com/golang/go/blob/go1.10.3/src/math/rand/rand.go#L371
type lockedSource struct {
	lk  sync.Mutex
	src rand.Source64
}

func (r *lockedSource) Int63() (n int64) {
	r.lk.Lock()
	n = r.src.Int63()
	r.lk.Unlock()
	return
}

func (r *lockedSource) Uint64() (n uint64) {
	r.lk.Lock()
	n = r.src.Uint64()
	r.lk.Unlock()
	return
}

func (r *lockedSource) Seed(seed int64) {
	r.lk.Lock()
	r.src.Seed(seed)
	r.lk.Unlock()
}

// lockedReadRand provides Rand.Read that is safe for concurrent use.
type lockedReadRand struct {
	lk sync.Mutex
	*rand.Rand
}

func (lr *lockedReadRand) Read(p []byte) (n int, err error) {
	lr.lk.Lock()
	n, err = lr.Rand.Read(p)
	lr.lk.Unlock()
	return
}

var (
	localSource = &lockedSource{src: rand.NewSource(time.Now().UTC().UnixNano()).(rand.Source64)}
	localRand   = &lockedReadRand{Rand: rand.New(localSource)}
)

// Seed uses the provided seed value to initialize the random source to a deterministic state.
func Seed(seed int64) {
	localRand.Seed(seed)
}
