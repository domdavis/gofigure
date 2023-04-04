package gofigure

import "math"

// Source of a Value, either accepted sources (i.e. what sources set a Value),
// or defining Source that set the Value. Combine multiple sources with | when
// defining accepted sources (e.g. Flag | EnvVar).
type Source uint8

const (
	// None indicates no source has set the Value, or that no source can set
	// this value.
	None = Source(1 << iota)

	// A Default Value has been used. No other sources have overwritten this.
	Default = Source(1 << iota)

	// A Key in a JSON file. Can be used either to indicate the Value can be set
	// via this Key, or has been set via this Key.
	Key = Source(1 << iota)

	// An EnvVar or environment variable. Can be used either to indicate the
	// Value can be set via this environment variable, or that it has been set
	// via this environment variable.
	EnvVar = Source(1 << iota)

	// A ShortFlag on the command line. Can be used either to indicate the
	// Value can be set via this flag, or that it has been set via this flag.
	ShortFlag = Source(1 << iota)

	// A Flag on the command line. Can be used either to indicate the
	// Value can be set via this flag, or that it has been set via this flag.
	Flag = Source(1 << iota)
)

const (
	// CommandLine indicates an option can come from a short or long flag.
	CommandLine = Flag | ShortFlag

	// NamedSources indicates an option should be set from a flag, environment
	// variable, or JSON key.
	NamedSources = Flag | EnvVar | Key

	// AllSources indicates an option should be set from a short flag, flag,
	// environment variable, or JSON key.
	AllSources = NamedSources | ShortFlag
)

// ConfigFile is used for error reporting purposes.
const configFile = Source(math.MaxUint8)

// Contains returns true if the Source contains the given Source.
func (s Source) Contains(source Source) bool {
	return s&source != 0
}

func (s Source) String() string {
	switch s {
	case None:
		return "none"
	case Default:
		return "default value"
	case Key:
		return "config file key"
	case EnvVar:
		return "environment value"
	case ShortFlag:
		return "short flag"
	case Flag:
		return "flag"
	case configFile:
		return "config file"
	default:
		return "source"
	}
}
