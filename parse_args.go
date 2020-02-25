package cli

import (
	"flag"
	"fmt"
)

type argInfo struct {
	req      []ArgSpec
	opt      []ArgSpec
	firstReq int
	firstOpt int

	// this helps makes sure that all the provided args mean something
	maxOptArgs   int
	unlimOptArgs bool

	// these help make sure that the optional args mean anything
	reqArgs int
}

// verifyArgSpec is used to make sure that the specs provided to the command are resonable
//
// The following rules are checked:
//
// * The slice of optional arguments are contiguous
//
// * The slice of required arguments are contiguous
//
// * All required arguments must of a fixed size (including slice types)
func verifyArgSpec(specs []ArgSpec) (*argInfo, error) {
	res := &argInfo{}

	if len(specs) == 0 {
		return res, nil
	}

	l, r := 0, 0
	req := specs[0].IsRequired()

	for i, spec := range specs[1:] {
		if spec.IsRequired() != req {
			break
		}

		r = i
	}

	if req {
		res.req = specs[l:r]
		res.opt = specs[r:]
		res.firstOpt = r
	} else {
		res.opt = specs[l:r]
		res.req = specs[r:]
		res.firstReq = r
	}

	req = !req
	for i, spec := range specs[l:] {
		if spec.IsRequired() != req {
			if req {
				// when from required -> optional
				return nil, fmt.Errorf(`Switching back to "optional" ArgSpec at index %d. `, i)
			} else {
				// when from optional -> required
				return nil, fmt.Errorf(`Switching back to "required" ArgSpec at index %d. `, i)
			}
		}
	}

	for i, spec := range res.req {
		if spec.IsSlice() {
			max := spec.MaxLength()

			if max == 0 {
				return nil, fmt.Errorf(`has unbounded size for "required"+"slice" ArgSpec at index %d`, res.firstReq+i)
			}

			res.reqArgs += int(max)
		} else {
			res.reqArgs++
		}
	}

	for _, spec := range res.opt {
		if spec.IsSlice() {
			max := spec.MaxLength()
			if max == 0 {
				res.unlimOptArgs = true
			}

			res.maxOptArgs += int(max)
		} else {
			res.maxOptArgs++
		}
	}

	return res, nil
}

func parseArgs(set *flag.FlagSet, info *argInfo, args []string) error {
	argsForReq := args[info.firstReq : info.firstReq+info.reqArgs]
	var argsForOpt []string
	if info.firstReq == 0 {
		// required args before optional
		argsForOpt = args[info.firstReq+info.reqArgs:]
	} else {
		// optional args before required
		argsForOpt = args[:info.firstReq]
	}

	err := parseRequiredArgs(set, info.req, argsForReq)
	if err != nil {
		return err
	}

	return parseOptionalArgs(set, info.opt, argsForOpt)
}

func parseRequiredArgs(set *flag.FlagSet, argSpec []ArgSpec, args []string) error {
	for front := argSpec[0]; len(argSpec) > 0; front, argSpec = argSpec[1], argSpec[1:] {
		count := 1
		if front.IsSlice() {
			count = int(front.MaxLength())
		}

		err := front.Parse(set, args[:count])
		if err != nil {
			return err
		}

		args = args[count:]
	}

	return nil
}

func parseOptionalArgs(set *flag.FlagSet, argSpec []ArgSpec, args []string) error {
	for front := argSpec[0]; len(argSpec) > 0 && len(args) > 0; front, argSpec = argSpec[1], argSpec[1:] {
		count := 1
		if front.IsSlice() {
			if max := front.MaxLength(); max > count {
				count = max
			}
			if len(args) < count {
				count = len(args)
			}
		}

		err := front.Parse(set, args[:count])
		if err != nil {
			return err
		}

		args = args[count:]
	}

	return nil
}
