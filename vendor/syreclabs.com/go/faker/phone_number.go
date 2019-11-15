package faker

type FakePhoneNumber interface {
	PhoneNumber() string                // => "1-599-267-6597 x537"
	CellPhone() string                  // => "+49-131-0003060"
	AreaCode() string                   // => "903"
	ExchangeCode() string               // => "574"
	SubscriberNumber(digits int) string // => "1512"
	String() string                     // String is an alias for PhoneNumber.
}

type fakePhoneNumber struct{}

func PhoneNumber() FakePhoneNumber {
	return fakePhoneNumber{}
}

func (p fakePhoneNumber) PhoneNumber() string {
	return Numerify(Fetch("phone_number.formats"))
}

func (p fakePhoneNumber) CellPhone() string {
	return Numerify(Fetch("cell_phone.formats"))
}

func (p fakePhoneNumber) AreaCode() string {
	var res string
	defer func() {
		if err := recover(); err != nil {
			res = ""
		}
	}()
	res = Fetch("phone_number.area_code")
	return res
}

func (p fakePhoneNumber) ExchangeCode() string {
	var res string
	defer func() {
		if err := recover(); err != nil {
			res = ""
		}
	}()
	res = Fetch("phone_number.exchange_code")
	return res
}

func (p fakePhoneNumber) SubscriberNumber(digits int) string {
	return Number().Number(digits)
}

func (p fakePhoneNumber) String() string {
	return p.PhoneNumber()
}
