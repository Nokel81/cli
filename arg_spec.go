package cli

import (
	"flag"
	"fmt"
)

// ArgSpec is an interface that allows us to mark flags as ones that can be parsed
// as arguments
type ArgSpec interface {
	fmt.Stringer

	// Apply Flag settings to the given flag set
	Parse(set *flag.FlagSet, values []string) error

	// Name returns the name that this arg will be accessed by
	AccessName() string

	// Required returns true if this arg is required
	IsRequired() bool

	// Slice returns true if this spec represents multiple values to be parsed
	IsSlice() bool

	// Represents the (min, max) permissable for this arg
	//
	// * max == 0 implies unlimited
	MaxLength() int
}
