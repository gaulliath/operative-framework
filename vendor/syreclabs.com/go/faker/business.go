package faker

type FakeBusiness interface {
	CreditCardNumber() string     // => "1234-2121-1221-1211"
	CreditCardExpiryDate() string // => "2015-11-11"
	CreditCardType() string       // => "mastercard"
}

type fakeBusiness struct{}

func Business() FakeBusiness {
	return fakeBusiness{}
}

func (b fakeBusiness) CreditCardNumber() string {
	return Fetch("business.credit_card_numbers")
}

func (b fakeBusiness) CreditCardExpiryDate() string {
	return Fetch("business.credit_card_expiry_dates")
}

func (b fakeBusiness) CreditCardType() string {
	return Fetch("business.credit_card_types")
}
