package cli

import (
	"flag"
	"fmt"
	"strconv"
)

// BoolFlag is a flag with type bool
type BoolFlag struct {
	Name        string
	Aliases     []string
	Usage       string
	EnvVars     []string
	FilePath    string
	Required    bool
	Hidden      bool
	Value       bool
	DefaultText string
	Destination *bool
	HasBeenSet  bool
}

// IsSet returns whether or not the flag has been set through env or file
func (f *BoolFlag) IsSet() bool {
	return f.HasBeenSet
}

// String returns a readable representation of this value
// (for usage defaults)
func (f *BoolFlag) String() string {
	return FlagStringer(f)
}

// Names returns the names of the flag
func (f *BoolFlag) Names() []string {
	return flagNames(f.Name, f.Aliases)
}

// IsRequired returns whether or not the flag is required
func (f *BoolFlag) IsRequired() bool {
	return f.Required
}

// TakesValue returns true of the flag takes a value, otherwise false
func (f *BoolFlag) TakesValue() bool {
	return false
}

// GetUsage returns the usage string for the flag
func (f *BoolFlag) GetUsage() string {
	return f.Usage
}

// GetValue returns the flags value as string representation and an empty
// string if the flag takes no value at all.
func (f *BoolFlag) GetValue() string {
	return ""
}

// Apply populates the flag given the flag set and environment
func (f *BoolFlag) Apply(set *flag.FlagSet) error {
	if val, ok := flagFromEnvOrFile(f.EnvVars, f.FilePath); ok {
		if val != "" {
			valBool, err := strconv.ParseBool(val)

			if err != nil {
				return fmt.Errorf("could not parse %q as bool value for flag %s: %s", val, f.Name, err)
			}

			f.Value = valBool
			f.HasBeenSet = true
		}
	}

	for _, name := range f.Names() {
		if f.Destination != nil {
			set.BoolVar(f.Destination, name, f.Value, f.Usage)
			continue
		}
		set.Bool(name, f.Value, f.Usage)
	}

	return nil
}

// Bool looks up the value of a local BoolFlag, returns
// false if not found
func (c *Context) Bool(name string) bool {
	if fs := lookupFlagSet(name, c); fs != nil {
		return lookupBool(name, fs)
	}
	return false
}

func lookupBool(name string, set *flag.FlagSet) bool {
	f := set.Lookup(name)
	if f != nil {
		parsed, err := strconv.ParseBool(f.Value.String())
		if err != nil {
			return false
		}
		return parsed
	}
	return false
}

type BoolArg struct {
	Name        string
	Usage       string
	Required    bool
	Value       bool
	Destination *bool
}

func (a *BoolArg) String() string {
	if a.Required {
		return fmt.Sprintf("[%s]", a.Name)
	} else {
		return fmt.Sprintf("<%s>", a.Name)
	}
}

func (a *BoolArg) Parse(set *flag.FlagSet, values []string) error {
	if len(values) == 0 {
		if a.Required {
			return fmt.Errorf("no value provided for required bool arg %q", a.Name)
		}

		return nil
	}

	if len(values) > 1 {
		return fmt.Errorf("too many values provided for required bool arg %q", a.Name)
	}

	val := values[0]
	valBool, err := strconv.ParseBool(val)
	if err != nil {
		return fmt.Errorf("could not parse %q as bool value for flag %s: %s", val, a.Name, err)
	}

	a.Value = valBool
	if a.Destination != nil {
		set.BoolVar(a.Destination, a.Name, a.Value, a.Usage)
	} else {
		set.Bool(a.Name, a.Value, a.Usage)
	}

	return nil
}

func (a *BoolArg) AccessName() string {
	return a.Name
}

func (a *BoolArg) IsRequired() bool {
	return a.Required
}

func (a *BoolArg) IsSlice() bool {
	return false
}

func (a *BoolArg) MaxLength() int {
	return 0
}
