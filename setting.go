package gofigure

import (
	"errors"
	"fmt"
)

// Setting in a Configuration. A Setting takes values from a set of Parameters
// and applies them to a Value. The Mask is used when generating a Display
// value.
type Setting struct {
	Value      *Value
	Parameters Parameters
	Mask       Mask
}

type Settings []*Setting

// ErrInvalidValue is used when an option can't be mapped.
var ErrInvalidValue = errors.New("invalid value")

// Optional Setting uses the given default value if no value is provided via its
// parameters. The parameters are constructed using the param value and the
// defined sources. Combine multiple sources with | (e.g. Flag | EnvVar). The
// given name is used for each source with Flag and Key using the name as is,
// EnvSuffix set to the uppercase version of the name, and ShortFlag set to the
// first character of name.
func Optional[T Type](name, param string, ptr *T, value T, sources Source,
	mask Mask, description string) *Setting {
	return &Setting{
		Value:      NewValue(name, ptr, value, Default, description),
		Parameters: NewParameters(param, sources),
		Mask:       mask,
	}
}

// Required Setting must be set via one of its Parameters. The parameters are
// constructed using the param value and the defined sources. Combine multiple
// sources with | (e.g. Flag | EnvVar). The given name is used for each source
// with Flag and Key using the name as is, EnvSuffix set to the uppercase
// version of the name, and ShortFlag set to the first character of name.
func Required[T Type](name, param string, ptr *T, sources Source, mask Mask,
	description string) *Setting {
	var value T

	return &Setting{
		Value:      NewValue(name, ptr, value, None, description),
		Parameters: NewParameters(param, sources),
		Mask:       mask,
	}
}

// Accepts returns true if this Setting accepts the given Parameter.
func (s Setting) Accepts(parameter Parameter) bool {
	for _, p := range s.Parameters {
		if p.Matches(parameter) {
			return s.Value.Source < parameter.Source
		}
	}

	return false
}

// Display string for the Setting. The string should only be displayed if
// Display returns true, otherwise it should be hidden.
func (s Setting) Display() (string, bool) {
	var value string

	unset := None

	if s.Value == nil {
		return Invalid, false
	} else if err := s.Value.Validate(); err != nil {
		value = Invalid
	} else {
		value = fmt.Sprintf("%v", Dereference(s.Value.Ptr))
	}

	if !s.Mask.Contains(DefaultIsSet) {
		unset |= Default
	}

	return s.display(value, unset)
}

func (s Setting) display(value string, unset Source) (string, bool) {
	display := true

	switch {
	case s.Value.Source.Contains(unset) && s.Mask.Contains(HideUnset):
		value = ""
		display = false
	case s.Value.Source.Contains(unset) && s.Mask.Contains(MaskUnset):
		value = NotSet
	case !s.Value.Source.Contains(unset) && s.Mask.Contains(HideSet):
		value = ""
		display = false
	case !s.Value.Source.Contains(unset) && s.Mask.Contains(MaskSet):
		value = Set
	}

	return value, display
}

// External configuration file paths defined by these settings.
func (s Settings) External() []string {
	var externals []string

	for _, setting := range s {
		path, ok := setting.Value.Ptr.(*External)

		if ok && *path != "" {
			externals = append(externals, string(*path))
		}
	}

	return externals
}

// Map the options to the settings, ignoring any errors provided.
func (s Settings) Map(options map[Parameter]any, ignore ...error) error {
	for parameter, value := range options {
		err := s.Apply(parameter, value)

		for _, e := range ignore {
			if errors.Is(err, e) {
				err = nil

				break
			}
		}

		switch {
		case errors.Is(err, ErrUnexpectedArgument):
			return NewConfigError(ErrUnexpectedArgument, err, parameter)
		case err != nil:
			return NewConfigError(fmt.Errorf("%w '%v'", ErrInvalidValue, value),
				err, parameter)
		}
	}

	return nil
}

// Apply the Parameter to the correct Setting in the set. Apply will return an
// error if the relevant Setting cannot be set, or if no Settings have been set.
func (s Settings) Apply(parameter Parameter, value any) error {
	for _, setting := range s {
		if !setting.Accepts(parameter) {
			continue
		} else if err := setting.Value.Assign(value, parameter.Source); err != nil {
			return fmt.Errorf("failed to apply %s: %w", parameter, err)
		} else {
			return nil
		}
	}

	return fmt.Errorf("%w: %s", ErrUnexpectedArgument, parameter)
}
