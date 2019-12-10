package faker

import "fmt"

type FakeAvatar interface {
	Url(format string, width, height int) string // => "http://robohash.org/NX34rZw7s0VFzgWY.jpg?size=100x200"
	String() string                              // => "http://robohash.org/XRWjFigoImqdeDuA.png?size=300x300"
}

type fakeAvatar struct{}

func Avatar() FakeAvatar {
	return fakeAvatar{}
}

func (a fakeAvatar) Url(format string, w, h int) string {
	return fmt.Sprintf("http://robohash.org/%s.%s?size=%dx%d", Lorem().Characters(16), format, w, h)
}

func (a fakeAvatar) String() string {
	return a.Url("png", 300, 300)
}
