package faker

import (
	"fmt"
	"reflect"
)

type FakeAddress interface {
	City() string                        // => "North Dessie"
	StreetName() string                  // => "Buckridge Lakes"
	StreetAddress() string               // => "586 Sylvester Turnpike"
	SecondaryAddress() string            // => "Apt. 411"
	BuildingNumber() string              // => "754"
	Postcode() string                    // => "31340"
	PostcodeByState(state string) string // => "46511"
	ZipCode() string                     // ZipCode is an alias for Postcode.
	ZipCodeByState(state string) string  // ZipCodeByState is an alias for PostcodeByState.
	TimeZone() string                    // => "Asia/Taipei"
	CityPrefix() string                  // => "East"
	CitySuffix() string                  // => "town"
	StreetSuffix() string                // => "Square"
	State() string                       // => "Maryland"
	StateAbbr() string                   // => "IL"
	Country() string                     // => "Uruguay"
	CountryCode() string                 // => "JP"
	Latitude() float32                   // => -38.811367
	Longitude() float32                  // => 89.2171
	String() string                      // => "6071 Heaney Island Suite 553, Ebbaville Texas 37307"
}

type fakeAddress struct{}

func Address() FakeAddress {
	return fakeAddress{}
}

func (a fakeAddress) City() string {
	return Fetch("address.city")
}

func (a fakeAddress) StreetName() string {
	return Fetch("address.street_name")
}

func (a fakeAddress) StreetAddress() string {
	return Numerify(Fetch("address.street_address"))
}

func (a fakeAddress) SecondaryAddress() string {
	return Numerify(Fetch("address.secondary_address"))
}

func (a fakeAddress) BuildingNumber() string {
	return NumerifyAndLetterify(Fetch("address.building_number"))
}

func (a fakeAddress) Postcode() string {
	return NumerifyAndLetterify(Fetch("address.postcode"))
}

func (a fakeAddress) PostcodeByState(state string) string {
	// postcode_by_state can be either a map[string] or a slice (as in default En locale)
	switch pbs := valueAt("address.postcode_by_state").(type) {
	case map[string]interface{}:
		_, ok := pbs[state]
		if ok {
			return NumerifyAndLetterify(Fetch("address.postcode_by_state." + state))
		}
		panic(fmt.Sprintf("invalid state: %v", state))
	case []string:
		return NumerifyAndLetterify(Fetch("address.postcode_by_state"))
	default:
		panic(fmt.Sprintf("invalid postcode_by_state value type: %v", reflect.TypeOf(pbs)))
	}
}

func (a fakeAddress) ZipCode() string {
	return a.Postcode()
}

func (a fakeAddress) ZipCodeByState(state string) string {
	return a.PostcodeByState(state)
}

func (a fakeAddress) TimeZone() string {
	return Fetch("address.time_zone")
}

func (a fakeAddress) CityPrefix() string {
	return Fetch("address.city_prefix")
}

func (a fakeAddress) CitySuffix() string {
	return Fetch("address.city_suffix")
}

func (a fakeAddress) StreetSuffix() string {
	return Fetch("address.street_suffix")
}

func (a fakeAddress) State() string {
	return Fetch("address.state")
}

func (a fakeAddress) StateAbbr() string {
	return Fetch("address.state_abbr")
}

func (a fakeAddress) Country() string {
	return Fetch("address.country")
}

func (a fakeAddress) CountryCode() string {
	return Fetch("address.country_code")
}

func (a fakeAddress) Latitude() float32 {
	return (localRand.Float32() * 180.0) - 90.0
}

func (a fakeAddress) Longitude() float32 {
	return (localRand.Float32() * 360.0) - 180.0
}

func (a fakeAddress) String() string {
	return fmt.Sprintf("%v %v, %v %v %v", a.StreetAddress(), a.SecondaryAddress(), a.City(), a.State(), a.Postcode())
}
