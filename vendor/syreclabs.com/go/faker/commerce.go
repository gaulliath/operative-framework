package faker

import (
	"math"
	"strings"
)

type FakeCommerce interface {
	Color() string       // => "lime"
	Department() string  // => "Electronics, Health & Baby"
	ProductName() string // => "Ergonomic Granite Shoes"
	Price() float32      // => 97.79
}

type fakeCommerce struct{}

func Commerce() FakeCommerce {
	return fakeCommerce{}
}

func (c fakeCommerce) Color() string {
	return Fetch("commerce.color")
}

func (c fakeCommerce) Department() string {
	n := RandomInt(1, 3)

	deps := make([]string, n)
	for i := range deps {
		d := Fetch("commerce.department")
		for includesString(deps, d) {
			d = Fetch("commerce.department")
		}
		deps[i] = d
	}

	if n > 1 {
		sep := Fetch("separator")
		res := strings.Join([]string{
			strings.Join(deps[0:len(deps)-1], ", "),
			deps[len(deps)-1],
		}, sep)
		return res
	}
	return deps[0]
}

func (c fakeCommerce) ProductName() string {
	words := []string{
		Fetch("commerce.product_name.adjective"),
		Fetch("commerce.product_name.material"),
		Fetch("commerce.product_name.product"),
	}
	return strings.Join(words, " ")
}

func (c fakeCommerce) Price() float32 {
	return float32(math.Floor(localRand.Float64()*10000.0)/100.0 + 0.01)
}
