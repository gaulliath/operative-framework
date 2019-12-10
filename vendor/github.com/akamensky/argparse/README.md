# Golang argparse

[![GoDoc](https://godoc.org/github.com/akamensky/argparse?status.svg)](https://godoc.org/github.com/akamensky/argparse) [![Go Report Card](https://goreportcard.com/badge/github.com/akamensky/argparse)](https://goreportcard.com/report/github.com/akamensky/argparse) [![Coverage Status](https://coveralls.io/repos/github/akamensky/argparse/badge.svg?branch=update-travis-yml)](https://coveralls.io/github/akamensky/argparse?branch=update-travis-yml) [![Build Status](https://travis-ci.org/akamensky/argparse.svg?branch=master)](https://travis-ci.org/akamensky/argparse)

Let's be honest -- Go's standard command line arguments parser `flag` terribly sucks. 
It cannot come anywhere close to the Python's `argparse` module. This is why this project exists.

The goal of this project is to bring ease of use and flexibility of `argparse` to Go. 
Which is where the name of this package comes from.

#### Installation

To install and start using argparse simply do:

```
$ go get -u -v github.com/akamensky/argparse
```

You are good to go to write your first command line tool!
See Usage and Examples sections for information how you can use it

#### Usage

To start using argparse in Go see above instructions on how to install.
From here on you can start writing your first program.
Please check out examples from `examples/` directory to see how to use it in various ways.

Here is basic example of print command (from `examples/print/` directory):
```go
package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	"os"
)

func main() {
	// Create new parser object
	parser := argparse.NewParser("print", "Prints provided string to stdout")
	// Create string flag
	s := parser.String("s", "string", &argparse.Options{Required: true, Help: "String to print"})
	// Parse input
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
	}
	// Finally print the collected string
	fmt.Println(*s)
}
```

#### Basic options

Create your parser instance and pass it program name and program description.
Program name if empty will be taken from `os.Args[0]` (which is okay in most cases).
Description can be as long as you wish and will be used in `--help` output
```go
parser := argparse.NewParser("progname", "Description of my awesome program. It can be as long as I wish it to be")
```

String will allow you to get a string from arguments, such as `$ progname --string "String content"`
```go
var myString *string = parser.String("s", "string", ...)
```

Flag will tell you if a simple flag was set on command line (true is set, false is not).
For example `$ progname --force`
```go
var myFlag *bool = parser.Flag("f", "force", ...)
```

List allows to collect multiple values into the slice of strings by repeating same flag multiple times.
Such as `$ progname --host hostname1 --host hostname2 -H hostname3`
```go
var myList *[]string = parser.List("H", "hostname", ...)
```

Selector works same as a string, except that it will only allow specific values.
For example like this `$ progname --debug-level WARN`
```go
var mySelector *string = parser.Selector("d", "debug-level", []string{"INFO", "DEBUG", "WARN"}, ...)
```

File will validate that file exists and will attempt to open it with provided privileges.
To be used like this `$ progname --log-file /path/to/file.log`
```go
var myLogFile *os.File = parser.File("l", "log-file", os.O_RDWR, 0600, ...)
```

You can implement sub-commands in your CLI using `parser.NewCommand()` or go even deeper with `command.NewCommand()`.
Since parser inherits from command, every command supports exactly same options as parser itself,
thus allowing to add arguments specific to that command or more global arguments added on parser itself!

#### Caveats

There are a few caveats (or more like design choices) to know about:
* Shorthand arguments MUST be a single character. Shorthand arguments are prepended with single dash `"-"`
* If not convenient shorthand argument can be completely skipped by passing empty string `""` as first argument
* Shorthand arguments ONLY for `parser.Flag()` can be combined into single argument same as `ps -aux` or `rm -rf`
* Long arguments must be specified and cannot be empty. They are prepended with double dash `"--"`
* You cannot define two same arguments. Only first one will be used. For example doing `parser.Flag("t", "test", nil)` followed by `parser.String("t", "test2", nil)` will not work as second `String` argument will be ignored (note that both have `"t"` as shorthand argument). However since it is case-sensitive library, you can work arounf it by capitalizing one of the arguments
* There is a pre-defined argument for `-h|--help`, so from above attempting to define any argument using `h` as shorthand will fail
* `parser.Parse()` returns error in case of something going wrong, but it is not expected to cover ALL cases
* Any arguments that left un-parsed will be regarded as error


#### Contributing

Can you write in Go? Then this projects needs your help!

Take a look at open issues, specially the ones tagged as `help-wanted`.
If you have any improvements to offer, please open an issue first to ensure this improvement is discussed.

There are following tasks to be done:
* Add more examples
* Improve code quality (it is messy right now and could use a major revamp to improve gocyclo report)
* Add more argument options (such as numbers parsing)
* Improve test coverage
* Write a wiki for this project

However note that the logic outlined in method comments must be preserved 
as the the library must stick with backward compatibility promise!

#### Acknowledgments

Thanks to Python developers for making a great `argparse` which inspired this package to match for greatness of Go
