package gofigure

import "fmt"

// Mask definition values when they are logged.
type Mask uint8

const (
	// HideSet will not report the value if it's set.
	HideSet = Mask(1 << iota)

	// HideUnset will not report anything if the value is not set.
	HideUnset = Mask(1 << iota)

	// MaskSet will report Set rather than the definition value.
	MaskSet = Mask(1 << iota)

	// MaskUnset will report Unset rather than the definition value.
	MaskUnset = Mask(1 << iota)

	// DefaultIsSet sets the behaviour of Mask to treat a default value as
	// Set rather than NotSet.
	DefaultIsSet = Mask(1 << iota)

	// ReportValue will report the definition value.
	ReportValue = 0

	// HideValue will not report anything for the definition.
	HideValue = HideSet | HideUnset

	// MaskValue will report Set or Unset rather than the definition value.
	MaskValue = MaskSet | MaskUnset
)

const (
	// Set is used when MaskSet is specified.
	Set = "SET"

	// NotSet is used when MaskNotSet is specified.
	NotSet = "UNSET"

	// Invalid is used when a Value is invalid.
	Invalid = "INVALID"
)

// Contains returns true if the Mask contains the given Mask.
func (m Mask) Contains(mask Mask) bool {
	return m&mask != 0
}

func (m Mask) String() string {
	switch m {
	case ReportValue:
		return "Report value"
	case HideSet:
		return "Hide value if set"
	case HideUnset:
		return "Hide value if unset"
	case MaskSet:
		return "Mask value if set"
	case MaskUnset:
		return "Mask value if unset"
	case HideValue:
		return "Hide value"
	case MaskValue:
		return "Mask value"
	default:
		return fmt.Sprintf("Mask value: %d", m)
	}
}
