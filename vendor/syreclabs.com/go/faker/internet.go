package faker

import (
	"fmt"
	"regexp"
	"strings"
)

type FakeInternet interface {
	Email() string                // => "maritza@farrell.org"
	FreeEmail() string            // => "sven_rice@hotmail.com"
	SafeEmail() string            // => "theron.nikolaus@example.net"
	UserName() string             // => "micah_pfeffer"
	Password(min, max int) string // => "s5CzvVp6Ye"
	DomainName() string           // => "rolfson.info"
	DomainWord() string           // => "heller"
	DomainSuffix() string         // => "net"
	MacAddress() string           // => "15:a9:83:29:76:26"
	IpV4Address() string          // => "121.204.82.227"
	IpV6Address() string          // => "c697:392f:6a0e:bf6d:77e1:714a:10ab:0dbc"
	Url() string                  // => "http://sporerhamill.net/kyla.schmitt"
	Slug() string                 // => "officiis-commodi"
}

type fakeInternet struct{}

func Internet() FakeInternet {
	return fakeInternet{}
}

func (i fakeInternet) Email() string {
	ss := []string{i.UserName(), i.DomainName()}
	return strings.Join(ss, "@")
}

func (i fakeInternet) FreeEmail() string {
	ss := []string{i.UserName(), Fetch("internet.free_email")}
	return strings.Join(ss, "@")
}

func (i fakeInternet) SafeEmail() string {
	ss := []string{
		i.UserName(),
		"example." + RandomChoice([]string{"com", "org", "net"}),
	}
	return strings.Join(ss, "@")
}

var separators = []string{".", "_"}
var rxNonWord = regexp.MustCompile(`\W`)

func sanitizeName(s string) string {
	return rxNonWord.ReplaceAllString(s, "")
}

func (i fakeInternet) UserName() string {
	sep := RandomChoice(separators)
	ss := []string{
		sanitizeName(Fetch("name.first_name")),
		strings.Join([]string{
			sanitizeName(Fetch("name.first_name")),
			sanitizeName(Fetch("name.last_name")),
		}, sep)}
	return strings.ToLower(RandomChoice(ss))
}

func (i fakeInternet) Password(min, max int) string {
	return RandomString(RandomInt(min, max))
}

func (i fakeInternet) DomainName() string {
	ss := []string{i.DomainWord(), i.DomainSuffix()}
	return strings.Join(ss, ".")
}

func (i fakeInternet) DomainWord() string {
	w := strings.Split(Fetch("company.name"), " ")[0]
	return strings.ToLower(sanitizeName(w))
}

func (i fakeInternet) DomainSuffix() string {
	return Fetch("internet.domain_suffix")
}

func (i fakeInternet) MacAddress() string {
	ss := make([]string, 6, 6)
	for i := range ss {
		ss[i] = fmt.Sprintf("%02x", localRand.Int31n(256))
	}
	return strings.Join(ss, ":")
}

func (i fakeInternet) IpV4Address() string {
	ss := make([]string, 4, 4)
	for i := range ss {
		ss[i] = fmt.Sprintf("%d", localRand.Int31n(256))
	}
	return strings.Join(ss, ".")
}

func (i fakeInternet) IpV6Address() string {
	ss := make([]string, 8, 8)
	for i := range ss {
		ss[i] = fmt.Sprintf("%04x", localRand.Int31n(65536))
	}
	return strings.Join(ss, ":")
}

func (i fakeInternet) Url() string {
	return fmt.Sprintf("http://%s/%s", i.DomainName(), i.UserName())
}

func (i fakeInternet) Slug() string {
	sep := "-"
	s := strings.Join(Lorem().Words(2), sep)
	return rxNonWord.ReplaceAllString(s, sep)
}
