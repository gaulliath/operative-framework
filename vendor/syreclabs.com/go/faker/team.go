package faker

type FakeTeam interface {
	Name() string     // => "Colorado cats"
	Creature() string // => "cats"
	State() string    // => "Oregon"
	String() string   // String is an alias for Name.
}

type fakeTeam struct{}

func Team() FakeTeam {
	return fakeTeam{}
}

func (t fakeTeam) Name() string {
	return Fetch("team.name")
}

func (t fakeTeam) Creature() string {
	return Fetch("team.creature")
}

func (t fakeTeam) State() string {
	return Fetch("address.state")
}

func (t fakeTeam) String() string {
	return t.Name()
}
