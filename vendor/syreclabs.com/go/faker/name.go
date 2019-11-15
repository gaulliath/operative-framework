package faker

type FakeName interface {
	Name() string      // => "Natasha Hartmann"
	FirstName() string // => "Carolina"
	LastName() string  // => "Kohler"
	Prefix() string    // => "Dr."
	Suffix() string    // => "Jr."
	Title() string     // => "Chief Functionality Orchestrator"
	String() string    // String is an alias for Name.
}

type fakeName struct{}

func Name() FakeName {
	return fakeName{}
}

func (n fakeName) Name() string {
	return Fetch("name.name")
}

func (n fakeName) FirstName() string {
	return Fetch("name.first_name")
}

func (n fakeName) LastName() string {
	return Fetch("name.last_name")
}

func (n fakeName) Prefix() string {
	return Fetch("name.prefix")
}

func (n fakeName) Suffix() string {
	return Fetch("name.suffix")
}

func (n fakeName) Title() string {
	return Fetch("name.title.descriptor") + " " + Fetch("name.title.level") + " " + Fetch("name.title.job")
}

func (n fakeName) String() string {
	return n.Name()
}
