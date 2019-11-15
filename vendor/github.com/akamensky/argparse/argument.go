package argparse

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type arg struct {
	result   interface{} // Pointer to the resulting value
	opts     *Options    // Options
	sname    string      // Short name (in Parser will start with "-"
	lname    string      // Long name (in Parser will start with "--"
	size     int         // Size defines how many args after match will need to be consumed
	unique   bool        // Specifies whether flag should be present only ones
	parsed   bool        // Specifies whether flag has been parsed already
	fileFlag int         // File mode to open file with
	filePerm os.FileMode // File permissions to set a file
	selector *[]string   // Used in Selector type to allow to choose only one from list of options
	parent   *Command    // Used to get access to specific Command
}

// Arg interface provides exporting of arg structure, while exposing it
type Arg interface {
	GetOpts() *Options
	GetSname() string
	GetLname() string
}

func (o arg) GetOpts() *Options {
	return o.opts
}

func (o arg) GetSname() string {
	return o.sname
}

func (o arg) GetLname() string {
	return o.lname
}

type help struct{}

func (o *arg) check(argument string) bool {
	// Shortcut to showing help
	if argument == "-h" || argument == "--help" {
		helpText := o.parent.Help(nil)
		fmt.Print(helpText)
		os.Exit(0)
	}

	// Check for long name only if not empty
	if o.lname != "" {
		// If argument begins with "--" and next is not "-" then it is a long name
		if len(argument) > 2 && strings.HasPrefix(argument, "--") && argument[2] != '-' {
			if argument[2:] == o.lname {
				return true
			}
		}
	}
	// Check for short name only if not empty
	if o.sname != "" {
		// If argument begins with "-" and next is not "-" then it is a short name
		if len(argument) > 1 && strings.HasPrefix(argument, "-") && argument[1] != '-' {
			switch o.result.(type) {
			case *bool:
				// For flags we allow multiple shorthand in one
				if strings.Contains(argument[1:], o.sname) {
					return true
				}
			default:
				// For all other types it must be separate argument
				if argument[1:] == o.sname {
					return true
				}
			}
		}
	}

	return false
}

func (o *arg) reduce(position int, args *[]string) {
	argument := (*args)[position]
	// Check for long name only if not empty
	if o.lname != "" {
		// If argument begins with "--" and next is not "-" then it is a long name
		if len(argument) > 2 && strings.HasPrefix(argument, "--") && argument[2] != '-' {
			if argument[2:] == o.lname {
				for i := position; i < position+o.size; i++ {
					(*args)[i] = ""
				}
			}
		}
	}
	// Check for short name only if not empty
	if o.sname != "" {
		// If argument begins with "-" and next is not "-" then it is a short name
		if len(argument) > 1 && strings.HasPrefix(argument, "-") && argument[1] != '-' {
			switch o.result.(type) {
			case *bool:
				// For flags we allow multiple shorthand in one
				if strings.Contains(argument[1:], o.sname) {
					(*args)[position] = strings.Replace(argument, o.sname, "", -1)
					if (*args)[position] == "-" {
						(*args)[position] = ""
					}
				}
			default:
				// For all other types it must be separate argument
				if argument[1:] == o.sname {
					for i := position; i < position+o.size; i++ {
						(*args)[i] = ""
					}
				}
			}
		}
	}
}

func (o *arg) parse(args []string) error {
	// If unique do not allow more than one time
	if o.unique && o.parsed {
		return fmt.Errorf("[%s] can only be present once", o.name())
	}

	// If validation function provided -- execute, on error return it immediately
	if o.opts != nil && o.opts.Validate != nil {
		err := o.opts.Validate(args)
		if err != nil {
			return err
		}
	}

	switch o.result.(type) {
	case *help:
		helpText := o.parent.Help(nil)
		fmt.Print(helpText)
		os.Exit(0)
	case *bool:
		*o.result.(*bool) = true
		o.parsed = true
	case *int:
		if len(args) < 1 {
			return fmt.Errorf("[%s] must be followed by an integer", o.name())
		}
		if len(args) > 1 {
			return fmt.Errorf("[%s] followed by too many arguments", o.name())
		}
		val, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("[%s] bad interger value [%s]", o.name(), args[0])
		}
		*o.result.(*int) = val
		o.parsed = true
	case *float64:
		if len(args) < 1 {
			return fmt.Errorf("[%s] must be followed by a floating point number", o.name())
		}
		if len(args) > 1 {
			return fmt.Errorf("[%s] followed by too many arguments", o.name())
		}
		val, err := strconv.ParseFloat(args[0], 64)
		if err != nil {
			return fmt.Errorf("[%s] bad floating point value [%s]", o.name(), args[0])
		}
		*o.result.(*float64) = val
		o.parsed = true
	case *string:
		if len(args) < 1 {
			return fmt.Errorf("[%s] must be followed by a string", o.name())
		}
		if len(args) > 1 {
			return fmt.Errorf("[%s] followed by too many arguments", o.name())
		}
		// Selector case
		if o.selector != nil {
			match := false
			for _, v := range *o.selector {
				if args[0] == v {
					match = true
				}
			}
			if !match {
				return fmt.Errorf("bad value for [%s]. Allowed values are %v", o.name(), *o.selector)
			}
		}
		*o.result.(*string) = args[0]
		o.parsed = true
	case *os.File:
		if len(args) < 1 {
			return fmt.Errorf("[%s] must be followed by a path to file", o.name())
		}
		if len(args) > 1 {
			return fmt.Errorf("[%s] followed by too many arguments", o.name())
		}
		f, err := os.OpenFile(args[0], o.fileFlag, o.filePerm)
		if err != nil {
			return err
		}
		*o.result.(*os.File) = *f
		o.parsed = true
	case *[]string:
		if len(args) < 1 {
			return fmt.Errorf("[%s] must be followed by a string", o.name())
		}
		if len(args) > 1 {
			return fmt.Errorf("[%s] followed by too many arguments", o.name())
		}
		*o.result.(*[]string) = append(*o.result.(*[]string), args[0])
		o.parsed = true
	default:
		return fmt.Errorf("unsupported type [%t]", o.result)
	}
	return nil
}

func (o *arg) name() string {
	var name string
	if o.lname == "" {
		name = "-" + o.sname
	} else if o.sname == "" {
		name = "--" + o.lname
	} else {
		name = "-" + o.sname + "|" + "--" + o.lname
	}
	return name
}

func (o *arg) usage() string {
	var result string
	result = o.name()
	switch o.result.(type) {
	case *bool:
		break
	case *int:
		result = result + " <integer>"
	case *float64:
		result = result + " <float>"
	case *string:
		if o.selector != nil {
			result = result + " (" + strings.Join(*o.selector, "|") + ")"
		} else {
			result = result + " \"<value>\""
		}
	case *os.File:
		result = result + " <file>"
	case *[]string:
		result = result + " \"<value>\"" + " [" + result + " \"<value>\" ...]"
	default:
		break
	}
	if o.opts == nil || o.opts.Required == false {
		result = "[" + result + "]"
	}
	return result
}

func (o *arg) getHelpMessage() string {
	message := ""
	if len(o.opts.Help) > 0 {
		message += o.opts.Help
		if !o.opts.Required && o.opts.Default != nil {
			message += fmt.Sprintf(". Default: %v", o.opts.Default)
		}
	}
	return message
}

func (o *arg) setDefault() error {
	// Only set default if it was not parsed, and default value was defined
	if !o.parsed && o.opts != nil && o.opts.Default != nil {
		switch o.result.(type) {
		case *bool:
			if _, ok := o.opts.Default.(bool); !ok {
				return fmt.Errorf("cannot use default type [%T] as type [bool]", o.opts.Default)
			}
			*o.result.(*bool) = o.opts.Default.(bool)
		case *int:
			if _, ok := o.opts.Default.(int); !ok {
				return fmt.Errorf("cannot use default type [%T] as type [int]", o.opts.Default)
			}
			*o.result.(*int) = o.opts.Default.(int)
		case *float64:
			if _, ok := o.opts.Default.(float64); !ok {
				return fmt.Errorf("cannot use default type [%T] as type [float64]", o.opts.Default)
			}
			*o.result.(*float64) = o.opts.Default.(float64)
		case *string:
			if _, ok := o.opts.Default.(string); !ok {
				return fmt.Errorf("cannot use default type [%T] as type [string]", o.opts.Default)
			}
			*o.result.(*string) = o.opts.Default.(string)
		case *os.File:
			// In case of File we should get string as default value
			if v, ok := o.opts.Default.(string); ok {
				f, err := os.OpenFile(v, o.fileFlag, o.filePerm)
				if err != nil {
					return err
				}
				*o.result.(*os.File) = *f
			} else {
				return fmt.Errorf("cannot use default type [%T] as type [string]", o.opts.Default)
			}
		case *[]string:
			if _, ok := o.opts.Default.([]string); !ok {
				return fmt.Errorf("cannot use default type [%T] as type [[]string]", o.opts.Default)
			}
			*o.result.(*[]string) = o.opts.Default.([]string)
		}
	}

	return nil
}
