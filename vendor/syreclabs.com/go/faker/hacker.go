package faker

import (
	"fmt"
	"unicode"
)

type FakeHacker interface {
	SaySomethingSmart() string // => "If we connect the bus, we can get to the XML microchip through the digital TCP sensor!"
	Abbreviation() string      // => "HTTP"
	Adjective() string         // => "cross-platform"
	Noun() string              // => "interface"
	Verb() string              // => "bypass"
	IngVerb() string           // => "parsing"
	Phrases() []string         /* =>
	"If we bypass the program, we can get to the AGP protocol through the optical SDD alarm!",
	"We need to calculate the back-end XML microchip!",
	"Try to generate the GB bus, maybe it will hack the neural panel!",
	"You can't navigate the transmitter without synthesizing the optical SMS bus!",
	"Use the optical THX application, then you can override the mobile port!",
	"The CSS monitor is down, quantify the multi-byte bus so we can calculate the XSS bandwidth!",
	"Connecting the card won't do anything, we need to back up the multi-byte RSS card!",
	"I'll reboot the primary SMTP feed, that should monitor the XML protocol!`"
	*/
}

type fakeHacker struct{}

func Hacker() FakeHacker {
	return fakeHacker{}
}

func (h fakeHacker) SaySomethingSmart() string {
	return RandomChoice(h.Phrases())
}

func (h fakeHacker) Abbreviation() string {
	return Fetch("hacker.abbreviation")
}

func (h fakeHacker) Adjective() string {
	return Fetch("hacker.adjective")
}

func (h fakeHacker) Noun() string {
	return Fetch("hacker.noun")
}

func (h fakeHacker) Verb() string {
	return Fetch("hacker.verb")
}

func (h fakeHacker) IngVerb() string {
	return Fetch("hacker.ingverb")
}

func capitalize(s string) string {
	return string(unicode.ToTitle(rune(s[0]))) + s[1:]
}

// Example:
//	package main
//
//	import (
//		"fmt"
//		"syreclabs.com/go/faker"
//	)
//
//	func main() {
//		for _, s := range (faker.Hacker{}.Phrases()) {
//			fmt.Println(s)
//		}
//	}
//
// Output:
// 	If we bypass the program, we can get to the AGP protocol through the optical SDD alarm!
// 	We need to calculate the back-end XML microchip!
// 	Try to generate the GB bus, maybe it will hack the neural panel!
// 	You can't navigate the transmitter without synthesizing the optical SMS bus!
// 	Use the optical THX application, then you can override the mobile port!
// 	The CSS monitor is down, quantify the multi-byte bus so we can calculate the XSS bandwidth!
// 	Connecting the card won't do anything, we need to back up the multi-byte RSS card!
// 	I'll reboot the primary SMTP feed, that should monitor the XML protocol!
func (h fakeHacker) Phrases() []string {
	return []string{
		fmt.Sprintf("If we %s the %s, we can get to the %s %s through the %s %s %s!",
			h.Verb(), h.Noun(), h.Abbreviation(), h.Noun(), h.Adjective(), h.Abbreviation(), h.Noun()),
		fmt.Sprintf("We need to %s the %s %s %s!",
			h.Verb(), h.Adjective(), h.Abbreviation(), h.Noun()),
		fmt.Sprintf("Try to %s the %s %s, maybe it will %s the %s %s!",
			h.Verb(), h.Abbreviation(), h.Noun(), h.Verb(), h.Adjective(), h.Noun()),
		fmt.Sprintf("You can't %s the %s without %s the %s %s %s!",
			h.Verb(), h.Noun(), h.IngVerb(), h.Adjective(), h.Abbreviation(), h.Noun()),
		fmt.Sprintf("Use the %s %s %s, then you can %s the %s %s!",
			h.Adjective(), h.Abbreviation(), h.Noun(), h.Verb(), h.Adjective(), h.Noun()),
		fmt.Sprintf("The %s %s is down, %s the %s %s so we can %s the %s %s!",
			h.Abbreviation(), h.Noun(), h.Verb(), h.Adjective(), h.Noun(), h.Verb(), h.Abbreviation(), h.Noun()),
		capitalize(fmt.Sprintf("%s the %s won't do anything, we need to %s the %s %s %s!",
			h.IngVerb(), h.Noun(), h.Verb(), h.Adjective(), h.Abbreviation(), h.Noun())),
		fmt.Sprintf("I'll %s the %s %s %s, that should %s the %s %s!",
			h.Verb(), h.Adjective(), h.Abbreviation(), h.Noun(), h.Noun(), h.Abbreviation(), h.Noun()),
	}
}
