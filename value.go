package gofigure

import (
	"errors"
	"fmt"
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

// Type accepted by gofigure.
type Type interface {
	~bool | ~float32 | ~float64 | ~string |
		~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// Value validation errors.
var (
	ErrMissingName        = errors.New("value for Value.Name is empty")
	ErrMissingDescription = errors.New("value for Value.Description is empty")
	ErrNilPointer         = errors.New("value for Value.Ptr is nil")
	ErrInvalidType        = errors.New("invalid type")
)

// NewValue returns a new, valid value. An empty name, description, or an
// invalid ptr will result in a panic.
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

// Assign the value to the target, returning an error if assignment fails.
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
		v, e := strconv.ParseInt(s, 10, 64)
		r, err = int(v), e
	case t.Kind() == reflect.Int8:
		v, e := strconv.ParseInt(s, 10, 64)
		r, err = int8(v), e
	case t.Kind() == reflect.Int16:
		v, e := strconv.ParseInt(s, 10, 64)
		r, err = int16(v), e
	case t.Kind() == reflect.Int32:
		v, e := strconv.ParseInt(s, 10, 64)
		r, err = int32(v), e
	case t.Kind() == reflect.Int64:
		r, err = strconv.ParseInt(s, 10, 64)
	case t.Kind() == reflect.Uint:
		v, e := strconv.ParseUint(s, 10, 64)
		r, err = uint(v), e
	case t.Kind() == reflect.Uint8:
		v, e := strconv.ParseInt(s, 10, 64)
		r, err = uint8(v), e
	case t.Kind() == reflect.Uint16:
		v, e := strconv.ParseInt(s, 10, 64)
		r, err = uint16(v), e
	case t.Kind() == reflect.Uint32:
		v, e := strconv.ParseInt(s, 10, 64)
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

// Cast a float64 value to an integer value. If the value isn't a float64, or
// the type isn't an integer then Cast will simply return the value.
//
//nolint:cyclop // Case switch for all supported types.
func Cast(value, typeOf any) any {
	f, ok := value.(float64)

	if !ok {
		return value
	}

	t := reflect.TypeOf(typeOf)

	switch {
	case t.Kind() == reflect.Int:
		return int(f)
	case t.Kind() == reflect.Int8:
		return int8(f)
	case t.Kind() == reflect.Int16:
		return int16(f)
	case t.Kind() == reflect.Int32:
		return int32(f)
	case t.Kind() == reflect.Int64:
		return int64(f)
	case t.Kind() == reflect.Uint:
		return uint(f)
	case t.Kind() == reflect.Uint8:
		return uint8(f)
	case t.Kind() == reflect.Uint16:
		return uint16(f)
	case t.Kind() == reflect.Uint32:
		return uint32(f)
	case t.Kind() == reflect.Uint64:
		return uint64(f)
	default:
		return value
	}
}

// Dereference a value. If the value isn't a pointer then it is returned as is.
func Dereference(in any) any {
	if in == nil || reflect.TypeOf(in).Kind() != reflect.Ptr {
		return in
	}

	return reflect.ValueOf(in).Elem().Interface()
}
