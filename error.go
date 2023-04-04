package gofigure

import (
	"fmt"
)

// ConfigError holds data about what specifically caused configuration to fail.
// This can e used to generate a user-friendly report.
type ConfigError struct {
	Cause    error
	Internal error

	Parameters Parameters
	Value      any
}

// NewConfigError will return a new ConfigError for the given errors and
// Parameters.
func NewConfigError(cause, internal error, parameters ...Parameter) ConfigError {
	return ConfigError{
		Cause:      cause,
		Internal:   internal,
		Parameters: parameters,
	}
}

// Format the error in a user-centric way.
func (c ConfigError) Format(prefix string) string {
	return fmt.Sprintf("%s: %s", c.Cause.Error(), c.Parameters.Format(prefix))
}

func (c ConfigError) Error() string {
	return c.Internal.Error()
}

func (c ConfigError) Unwrap() error {
	return c.Internal
}
