package faker

import "fmt"

type FakeApp interface {
	Name() string    // => "Alphazap"
	Version() string // => "2.6.0"
	Author() string  // => "Dorian Shields"
	String() string  // => "Tempsoft 4.51"
}

type fakeApp struct{}

func App() FakeApp {
	return fakeApp{}
}

func (a fakeApp) Name() string {
	return Fetch("app.name")
}

func (a fakeApp) Version() string {
	return Numerify(Fetch("app.version"))
}

func (a fakeApp) Author() string {
	return Fetch("app.author")
}

func (a fakeApp) String() string {
	return fmt.Sprintf("%v %v", a.Name(), a.Version())
}
