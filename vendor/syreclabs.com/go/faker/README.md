# faker

Faker is a library for generating fake data such as names, addresses, and phone numbers.

It is a (mostly) API-compatible port of Ruby Faker gem (https://github.com/stympy/faker) to Go.

[![GoDoc](https://godoc.org/syreclabs.com/go/faker?status.svg)](https://godoc.org/syreclabs.com/go/faker)
[![Build Status](https://travis-ci.org/dmgk/faker.svg?branch=master)](https://travis-ci.org/dmgk/faker)
[![Coverage Status](https://coveralls.io/repos/github/dmgk/faker/badge.svg)](https://coveralls.io/github/dmgk/faker)

To install

    go get -u syreclabs.com/go/faker

## Usage

### Address
-------------------
```go
faker.Address().City()                        // => "North Dessie"
faker.Address().StreetName()                  // => "Buckridge Lakes"
faker.Address().StreetAddress()               // => "586 Sylvester Turnpike"
faker.Address().SecondaryAddress()            // => "Apt. 411"
faker.Address().BuildingNumber()              // => "754"
faker.Address().Postcode()                    // => "31340"
faker.Address().PostcodeByState("IN")         // => "46511"
faker.Address().ZipCode()                     // ZipCode is an alias for Postcode.
faker.Address().ZipCodeByState("IN")          // ZipCodeByState is an alias for PostcodeByState.
faker.Address().TimeZone()                    // => "Asia/Taipei"
faker.Address().CityPrefix()                  // => "East"
faker.Address().CitySuffix()                  // => "town"
faker.Address().StreetSuffix()                // => "Square"
faker.Address().State()                       // => "Maryland"
faker.Address().StateAbbr()                   // => "IL"
faker.Address().Country()                     // => "Uruguay"
faker.Address().CountryCode()                 // => "JP"
faker.Address().Latitude()                    // => (float32) -38.811367
faker.Address().Longitude()                   // => (float32) 89.2171
faker.Address().String()                      // => "6071 Heaney Island Suite 553, Ebbaville Texas 37307"
```

### App
-------------------
```go
faker.App().Name()    // => "Alphazap"
faker.App().Version() // => "2.6.0"
faker.App().Author()  // => "Dorian Shields"
faker.App().String()  // => "Tempsoft 4.51"
```

### Avatar
-------------------
```go
faker.Avatar().Url("jpg", 100, 200) // => "http://robohash.org/NX34rZw7s0VFzgWY.jpg?size=100x200"
faker.Avatar().String()             // => "http://robohash.org/XRWjFigoImqdeDuA.png?size=300x300"
```

### Bitcoin
-------------------
```go
faker.Bitcoin().Address() // => "1GpEKM5UvD4XDLMirpNLoDnRVrGutogMj2"
faker.Bitcoin().String()  // String is an alias for Address.
```

### Business
-------------------
```go
faker.Business().CreditCardNumber()     // => "1234-2121-1221-1211"
faker.Business().CreditCardExpiryDate() // => "2015-11-11"
faker.Business().CreditCardType()       // => "mastercard"
```

### Code
-------------------
```go
faker.Code().Isbn10() // => "048931033-8"
faker.Code().Isbn13() // => "391668236072-1"
faker.Code().Ean13()  // => "7742864258656"
faker.Code().Ean8()   // => "03079010"
faker.Code().Rut()    // => "14371602-3"
faker.Code().Abn()    // => "57914951376"
```

### Commerce
-------------------
```go
faker.Commerce().Color()       // => "lime"
faker.Commerce().Department()  // => "Electronics, Health & Baby"
faker.Commerce().ProductName() // => "Ergonomic Granite Shoes"
faker.Commerce().Price()       // => (float32) 97.79
```

### Company
-------------------
```go
faker.Company().Name()        // => "Aufderhar LLC"
faker.Company().Suffix()      // => "Inc"
faker.Company().CatchPhrase() // => "Universal logistical artificial intelligence"
faker.Company().Bs()          // => "engage distributed applications"
faker.Company().Ein()         // => "58-6520513"
faker.Company().DunsNumber()  // => "16-708-2968"
faker.Company().Logo()        // => "http://www.biz-logo.com/examples/015.gif"
faker.Company().String()      // String is an alias for Name.
```

### Date
-------------------
```go
// Between returns random time in [from, to] interval, with second resolution.
faker.Date().Between(from, to time.Time) time.Time

// Forward returns random time in [time.Now(), time.Now() + duration] interval, with second resolution.
faker.Date().Forward(duration time.Duration) time.Time

// Backward returns random time in [time.Now() - duration, time.Now()] interval, with second resolution.
faker.Date().Backward(duration time.Duration) time.Time

// Birthday returns random time so that age of the person born at that moment would be between minAge and maxAge years.
faker.Date().Birthday(minAge, maxAge int) time.Time
```

### Finance
-------------------
```go
// CreditCard returns a valid (with valid check digit) card number of one of the given types.
// If no types are passed, all types in CC_TYPES are used.
faker.Finance().CreditCard(faker.CC_VISA) // => "4190418835414"
```

### Hacker
-------------------
```go
faker.Hacker().SaySomethingSmart() // => "If we connect the bus, we can get to the XML microchip through the digital TCP sensor!"
faker.Hacker().Abbreviation()      // => "HTTP"
faker.Hacker().Adjective()         // => "cross-platform"
faker.Hacker().Noun()              // => "interface"
faker.Hacker().Verb()              // => "bypass"
faker.Hacker().IngVerb()           // => "parsing"
faker.Hacker().Phrases() []string  // => []string{
                                   //        "If we bypass the program, we can get to the AGP protocol through the optical SDD alarm!",
                                   //        "We need to calculate the back-end XML microchip!",
                                   //        "Try to generate the GB bus, maybe it will hack the neural panel!",
                                   //        "You can't navigate the transmitter without synthesizing the optical SMS bus!",
                                   //        "Use the optical THX application, then you can override the mobile port!",
                                   //        "The CSS monitor is down, quantify the multi-byte bus so we can calculate the XSS bandwidth!",
                                   //        "Connecting the card won't do anything, we need to back up the multi-byte RSS card!",
                                   //        "I'll reboot the primary SMTP feed, that should monitor the XML protocol!`",
                                   //    }
```

### Internet
-------------------
```go
faker.Internet().Email()                // => "maritza@farrell.org"
faker.Internet().FreeEmail()            // => "sven_rice@hotmail.com"
faker.Internet().SafeEmail()            // => "theron.nikolaus@example.net"
faker.Internet().UserName()             // => "micah_pfeffer"
faker.Internet().Password(8, 14)        // => "s5CzvVp6Ye"
faker.Internet().DomainName()           // => "rolfson.info"
faker.Internet().DomainWord()           // => "heller"
faker.Internet().DomainSuffix()         // => "net"
faker.Internet().MacAddress()           // => "15:a9:83:29:76:26"
faker.Internet().IpV4Address()          // => "121.204.82.227"
faker.Internet().IpV6Address()          // => "c697:392f:6a0e:bf6d:77e1:714a:10ab:0dbc"
faker.Internet().Url()                  // => "http://sporerhamill.net/kyla.schmitt"
faker.Internet().Slug()                 // => "officiis-commodi"
```

### Lorem
-------------------
```go
faker.Lorem().Character()    // => "c"
faker.Lorem().Characters(17) // => "wqFyJIrXYfVP7cL9M"
faker.Lorem().Word()         // => "veritatis"
faker.Lorem().Words(3)       // => []string{"omnis", "libero", "neque"}
faker.Lorem().Sentence(3)    // => "Necessitatibus sit autem."

// Sentences returns a slice of "num" sentences, 3 to 11 words each.
faker.Lorem().Sentences(num int) []string

// Paragraph returns a random text of "sentences" sentences length.
faker.Lorem().Paragraph(sentences int)

// Paragraphs returns a slice of "num" paragraphs, 3 to 11 sentences each.
faker.Lorem().Paragraphs(num int) []string

// String returns a random sentence 3 to 11 words in length.
faker.Lorem().String()
```

### Name
-------------------
```go
faker.Name().Name()      // => "Natasha Hartmann"
faker.Name().FirstName() // => "Carolina"
faker.Name().LastName()  // => "Kohler"
faker.Name().Prefix()    // => "Dr."
faker.Name().Suffix()    // => "Jr."
faker.Name().Title()     // => "Chief Functionality Orchestrator"
faker.Name().String()    // String is an alias for Name.
```

### Number
-------------------
```go
faker.Number().Number(5)          // => "43202"
faker.Number().NumberInt(3)       // => 213
faker.Number().NumberInt32(5)     // => 92938
faker.Number().NumberInt64(19)    // => 1689541633257139096
faker.Number().Decimal(8, 2)      // => "879420.60"
faker.Number().Digit()            // => "7"
faker.Number().Hexadecimal(4)     // => "e7f3"
faker.Number().Between(-100, 100) // => "-47"
faker.Number().Positive(100)      // => "3"
faker.Number().Negative(-100)     // => "-16"
```

### PhoneNumber
-------------------
```go
faker.PhoneNumber().PhoneNumber()       // => "1-599-267-6597 x537"
faker.PhoneNumber().CellPhone()         // => "+49-131-0003060"
faker.PhoneNumber().AreaCode()          // => "903"
faker.PhoneNumber().ExchangeCode()      // => "574"
faker.PhoneNumber().SubscriberNumber(4) // => "1512"
faker.PhoneNumber().String()            // String is an alias for PhoneNumber.
```

### Team
-------------------
```go
faker.Team().Name()     // => "Colorado cats"
faker.Team().Creature() // => "cats"
faker.Team().State()    // => "Oregon"
faker.Team().String()   // String is an alias for Name.
```

### Time
-------------------
```go
// Between returns random time in [from, to] interval, with second resolution.
faker.Time().Between(from, to time.Time) time.Time

// Forward returns random time in [time.Now(), time.Now() + duration] interval, with second resolution.
faker.Time().Forward(duration time.Duration) time.Time

// Backward returns random time in [time.Now() - duration, time.Now()] interval, with second resolution.
faker.Time().Backward(duration time.Duration) time.Time

// Birthday returns random time so that age of the person born at that moment would be between minAge and maxAge years.
faker.Time().Birthday(minAge, maxAge int) time.Time
```
