package faker

import "fmt"

type FakeCompany interface {
	Name() string        // => "Aufderhar LLC"
	Suffix() string      // => "Inc"
	CatchPhrase() string // => "Universal logistical artificial intelligence"
	Bs() string          // => "engage distributed applications"
	Ein() string         // => "58-6520513"
	DunsNumber() string  // => "16-708-2968"
	Logo() string        // => "http://www.biz-logo.com/examples/015.gif"
	String() string      // String is an alias for Name.
}

type fakeCompany struct{}

func Company() FakeCompany {
	return fakeCompany{}
}

func (c fakeCompany) Name() string {
	return Fetch("company.name")
}

func (c fakeCompany) Suffix() string {
	return Fetch("company.suffix")
}

func (c fakeCompany) CatchPhrase() string {
	return Fetch("company.buzzwords")
}

func (c fakeCompany) Bs() string {
	return Fetch("company.bs")
}

func (c fakeCompany) Ein() string {
	return Numerify("##-#######")
}

func (c fakeCompany) DunsNumber() string {
	return Numerify("##-###-####")
}

func (c fakeCompany) Logo() string {
	n := RandomInt(1, 77)
	return fmt.Sprintf("http://www.biz-logo.com/examples/%03d.gif", n)
}

func (c fakeCompany) String() string {
	return c.Name()
}
