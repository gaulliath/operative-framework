package faker

type FakeTime interface {
	FakeDate
}

type fakeTime struct {
	fakeDate
}

func Time() FakeTime {
	return fakeTime{}
}
