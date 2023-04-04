package gofigure

import (
	"errors"
	"fmt"
	"strings"
)

// ErrUnexpectedArgument is returned if an unexpected argument is passed.
var ErrUnexpectedArgument = errors.New("unexpected argument")

// Flags will build a set of Options from the given argument list.
func Flags(args []string) (Options, error) {
	const (
		short = 1
		long  = 2
		flag  = "-"
	)

	options := Options{}

	for len(args) > 0 {
		parameter := Parameter{}

		switch strings.Count(args[0], flag) {
		case long:
			parameter.Source = Flag
		case short:
			parameter.Source = ShortFlag
		default:
			return options, NewConfigError(ErrUnexpectedArgument,
				fmt.Errorf("%w: %s", ErrUnexpectedArgument, args[0]),
				Parameter{Name: args[0], Source: CommandLine})
		}

		name := strings.TrimLeft(args[0], flag)

		args = args[1:]
		parameter.Name = name

		if len(args) == 0 || strings.HasPrefix(args[0], flag) {
			options[parameter] = "true"
		} else {
			options[parameter] = args[0]
			args = args[1:]
		}
	}

	return options, nil
}
