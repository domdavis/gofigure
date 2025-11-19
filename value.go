package gofigure

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"time"
)

// A Value is used to hold a configured value. The Value must be a pointer to
// the variable being set, and must satisfy Type. Once set the Value will
// contain the Source that provided the value.
type Value struct {
	Name        string
	Description string
	Ptr         any

	Source Source

	base any
}

// External types hold a path to an external configuration file.
type External string

// Int types accepted by gofigure.
type Int interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Uint types accepted by gofigure.
type Uint interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// Float types accepted by gofigure.
type Float interface{ ~float32 | ~float64 }

// Numeric types accepted by gofigure.
type Numeric interface{ Int | Uint | Float }

// Type accepted by gofigure.
type Type interface{ ~bool | ~string | Numeric }

// Value validation errors.
var (
	ErrMissingName        = errors.New("value for Value.Name is empty")
	ErrMissingDescription = errors.New("value for Value.Description is empty")
	ErrNilPointer         = errors.New("value for Value.Ptr is nil")
	ErrInvalidType        = errors.New("invalid type")
)

// NewValue returns a new, valid value.
//
// NewValue panics if:
//   - name is empty
//   - description is empty
//   - ptr is nil
//   - ptr points to an invalid type
//
// These panics indicate programming errors (invalid arguments) and should be
// caught during development and testing.
func NewValue[T Type](name string, ptr *T, value T, source Source, description string) *Value {
	s := &Value{Name: name, Description: description, Ptr: ptr, Source: source}

	if err := s.Validate(); err != nil {
		panic(err)
	}

	if source.Contains(Default) {
		*ptr = value
		s.base = value
	}

	return s
}

// Validate the setting returning an error if the Value lacks a name or
// description, if the Ptr is nil, or if the Ptr is of the incorrect type.
func (v *Value) Validate() error {
	switch {
	case v.Name == "":
		return ErrMissingName
	case v.Description == "":
		return fmt.Errorf("%w: (Value %s)", ErrMissingDescription, v.Name)
	case Dereference(v.Ptr) == nil:
		return fmt.Errorf("%w: (Value %s)", ErrNilPointer, v.Name)
	}

	switch v.Ptr.(type) {
	case *bool, *float32, *float64, *string, *time.Duration, *External,
		*int, *int8, *int16, *int32, *int64,
		*uint, *uint8, *uint16, *uint32, *uint64:
	default:
		return fmt.Errorf("%w for Value %s: %T", ErrInvalidType, v.Name, v.Ptr)
	}

	return nil
}

// Assign a value to the Value.Ptr, returning an error if the assignment
// cannot be made.
//
//nolint:cyclop,funlen // Case switch for all available types.
func (v *Value) Assign(value any, source Source) error {
	err := v.Validate()

	if err != nil {
		return fmt.Errorf("cannot assign %v to invalid setting: %w", value, err)
	}

	switch target := v.Ptr.(type) {
	case *bool:
		err = Assign(target, value)
	case *int:
		err = Assign(target, value)
	case *int8:
		err = Assign(target, value)
	case *int16:
		err = Assign(target, value)
	case *int32:
		err = Assign(target, value)
	case *int64:
		err = Assign(target, value)
	case *uint:
		err = Assign(target, value)
	case *uint8:
		err = Assign(target, value)
	case *uint16:
		err = Assign(target, value)
	case *uint32:
		err = Assign(target, value)
	case *uint64:
		err = Assign(target, value)
	case *float32:
		err = Assign(target, value)
	case *float64:
		err = Assign(target, value)
	case *string:
		err = Assign(target, value)
	case *time.Duration:
		err = Assign(target, value)
	case *External:
		if e, ok := value.(External); ok {
			err = Assign(target, e)
		} else if s, ok := value.(string); !ok {
			err = ErrInvalidType
		} else {
			err = Assign(target, External(s))
		}
	}

	if err != nil {
		return fmt.Errorf("%w: cannot assign %T to %T", err, value, v.Ptr)
	}

	v.Source = source

	return nil
}

// Assign the value to the target, returning an error if the assignment fails.
// Assign will attempt to coerce string values to the correct type.
func Assign[T Type](target *T, value any) (err error) {
	var (
		ok     bool
		typeOf T
	)

	value, err = Coerce(value, typeOf)

	if err != nil {
		return fmt.Errorf("assignment error: %w", err)
	}

	*target, ok = Cast(value, typeOf).(T)

	if !ok {
		return ErrInvalidType
	}

	return nil
}

// Coerce will coerce string values to the correct type. All other value types
// are simply returned as is. If the string value cannot be coerced to the type
// an error is returned.
//
//nolint:cyclop,funlen // Case switch for all available types.
func Coerce(value, typeOf any) (r any, err error) {
	s, ok := value.(string)

	if !ok {
		return value, nil
	}

	t := reflect.TypeOf(typeOf)

	switch {
	case t == reflect.TypeOf(time.Duration(1)):
		r, err = time.ParseDuration(s)
	case t.Kind() == reflect.Bool:
		r, err = strconv.ParseBool(s)
	case t.Kind() == reflect.Int:
		v, e := strconv.ParseInt(s, 10, 0)
		r, err = int(v), e
	case t.Kind() == reflect.Int8:
		v, e := strconv.ParseInt(s, 10, 8)
		r, err = int8(v), e
	case t.Kind() == reflect.Int16:
		v, e := strconv.ParseInt(s, 10, 16)
		r, err = int16(v), e
	case t.Kind() == reflect.Int32:
		v, e := strconv.ParseInt(s, 10, 32)
		r, err = int32(v), e
	case t.Kind() == reflect.Int64:
		r, err = strconv.ParseInt(s, 10, 64)
	case t.Kind() == reflect.Uint:
		v, e := strconv.ParseUint(s, 10, 0)
		r, err = uint(v), e
	case t.Kind() == reflect.Uint8:
		v, e := strconv.ParseUint(s, 10, 8)
		r, err = uint8(v), e
	case t.Kind() == reflect.Uint16:
		v, e := strconv.ParseUint(s, 10, 16)
		r, err = uint16(v), e
	case t.Kind() == reflect.Uint32:
		v, e := strconv.ParseUint(s, 10, 32)
		r, err = uint32(v), e
	case t.Kind() == reflect.Uint64:
		r, err = strconv.ParseUint(s, 10, 64)
	case t.Kind() == reflect.Float32:
		v, e := strconv.ParseFloat(s, 32)
		r, err = float32(v), e
	case t.Kind() == reflect.Float64:
		r, err = strconv.ParseFloat(s, 64)
	case t.Kind() == reflect.String:
		r = s
	default:
		r = s
		err = fmt.Errorf("%w: %v", ErrInvalidType, t)
	}

	if err != nil {
		err = fmt.Errorf("cannot coerce %q: %w", s, err)
	}

	return r, err
}

// Cast a numeric value to the target type. Handles float64 and int conversions
// to appropriately sized integer types. If the value isn't numeric, Cast will
// simply return the value.
//
//nolint:gosec // Bounds checking is handled before the cast.
func Cast(value, typeOf any) any {
	type clamp struct {
		max uint
		min int
		f   func(i int) any
	}

	var ranges = map[reflect.Kind]clamp{
		reflect.Int:    {max: math.MaxInt, min: math.MinInt, f: func(i int) any { return i }},
		reflect.Int8:   {max: math.MaxInt8, min: math.MinInt8, f: func(i int) any { return int8(i) }},
		reflect.Int16:  {max: math.MaxInt16, min: math.MinInt16, f: func(i int) any { return int16(i) }},
		reflect.Int32:  {max: math.MaxInt32, min: math.MinInt32, f: func(i int) any { return int32(i) }},
		reflect.Int64:  {max: math.MaxInt64, min: math.MinInt64, f: func(i int) any { return int64(i) }},
		reflect.Uint:   {max: math.MaxUint, min: 0, f: func(i int) any { return uint(i) }},
		reflect.Uint8:  {max: math.MaxUint8, min: 0, f: func(i int) any { return uint8(i) }},
		reflect.Uint16: {max: math.MaxUint16, min: 0, f: func(i int) any { return uint16(i) }},
		reflect.Uint32: {max: math.MaxUint32, min: 0, f: func(i int) any { return uint32(i) }},
		reflect.Uint64: {max: math.MaxUint64, min: 0, f: func(i int) any { return uint64(i) }},
	}

	var v int

	if f, float := value.(float64); float {
		v = int(f)
	} else if i, integer := value.(int); integer {
		v = i
	} else {
		return value
	}

	c, clampable := ranges[reflect.TypeOf(typeOf).Kind()]

	if !clampable {
		return value
	}

	if v < c.min || uint(v) > c.max {
		return value
	}

	return c.f(v)
}

// Dereference a value. If the value isn't a pointer then it is returned as is.
func Dereference(in any) any {
	if in == nil || reflect.TypeOf(in).Kind() != reflect.Ptr {
		return in
	}

	return reflect.ValueOf(in).Elem().Interface()
}
