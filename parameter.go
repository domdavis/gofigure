package gofigure

import (
	"fmt"
	"strings"
)

// Parameter used to set a Value.
type Parameter struct {
	Name   string
	Source Source

	Stub string
}

// Parameters that can be used to set a Value.
type Parameters []Parameter

// NewParameters returns a set of named parameters for the given sources.
// Combine multiple sources with | (e.g. Flag | EnvVar). The given name is used
// for each source with Flag and Key using the name as is, EnvSuffix set to the
// uppercase version of the name, and ShortFlag set to the first character of
// name.
func NewParameters(name string, sources Source) Parameters {
	var p []Parameter

	if name == "" {
		panic("parameter must have a name")
	}

	if sources.Contains(Key) {
		p = append(p, Parameter{Name: name, Source: Key})
	}

	if sources.Contains(EnvVar) {
		p = append(p, Parameter{Name: strings.ToUpper(name), Source: EnvVar})
	}

	if sources.Contains(ShortFlag) {
		p = append(p, Parameter{Name: string(name[0]), Source: ShortFlag})
	}

	if sources.Contains(Flag) {
		p = append(p, Parameter{Name: name, Source: Flag})
	}

	return p
}

// FullName returns the fully formatted name for the parameter.
func (p Parameter) FullName() string {
	if p.Source != EnvVar || p.Stub == "" {
		return p.Name
	}

	return fmt.Sprintf("%s_%s", p.Stub, p.Name)
}

func (p Parameter) String() string {
	switch p.Source {
	case Flag:
		return "--" + p.Name
	case ShortFlag:
		return "-" + p.Name
	case EnvVar:
		return fmt.Sprintf("env %s", p.FullName())
	case Key:
		return fmt.Sprintf("JSON key: %q", p.Name)
	case configFile:
		return "file: " + p.Name
	default:
		return p.Name
	}
}

// Matches returns true if the Argument matches the Parameter.
func (p Parameter) Matches(parameter Parameter) bool {
	return p.Name == parameter.Name && p.Source.Contains(parameter.Source)
}

// Format the given Parameters into a human-readable string.
func (p Parameters) Format(prefix string) string {
	parameters := make([]string, len(p))

	for i, parameter := range p {
		parameter.Stub = prefix
		parameters[i] = parameter.String()
	}

	return fmt.Sprintf("[%s]", strings.Join(parameters, ", "))
}
