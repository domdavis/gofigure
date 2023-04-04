package gofigure_test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"bitbucket.org/idomdavis/gofigure"
	"github.com/stretchr/testify/assert"
)

func ExampleCoerce() {
	fmt.Println(gofigure.Coerce(1, 0))
	fmt.Println(gofigure.Coerce("1m0s", time.Nanosecond))
	fmt.Println(gofigure.Coerce("true", false))
	fmt.Println(gofigure.Coerce("1", 0))
	fmt.Println(gofigure.Coerce("1", int64(0)))
	fmt.Println(gofigure.Coerce("1", uint(0)))
	fmt.Println(gofigure.Coerce("1", uint64(0)))
	fmt.Println(gofigure.Coerce("1.1", float64(0)))
	fmt.Println(gofigure.Coerce("1", int8(0)))

	// Output:
	// 1 <nil>
	// 1m0s <nil>
	// true <nil>
	// 1 <nil>
	// 1 <nil>
	// 1 <nil>
	// 1 <nil>
	// 1.1 <nil>
	// 1 cannot coerce "1": invalid type: int8
}

func ExampleDereference() {
	i := 1

	fmt.Println(gofigure.Dereference(i))
	fmt.Println(gofigure.Dereference(&i))

	// Output:
	// 1
	// 1
}

func TestNewValue(t *testing.T) {
	t.Run("Name must be set", func(t *testing.T) {
		t.Parallel()

		assert.Panics(t, func() {
			gofigure.NewValue("", new(bool), false, 0, "description")
		})
	})

	t.Run("Description must be set", func(t *testing.T) {
		t.Parallel()

		assert.Panics(t, func() {
			gofigure.NewValue("name", new(bool), false, 0, "")
		})
	})

	t.Run("Ptr can't be nil", func(t *testing.T) {
		t.Parallel()

		assert.Panics(t, func() {
			gofigure.NewValue("name", nil, false, 0, "description")
		})
	})
}

func TestValue_Validate(t *testing.T) {
	t.Run("Ptr must be a valid type", func(t *testing.T) {
		t.Parallel()

		v := gofigure.Value{
			Name:        "name",
			Description: "description",
			Ptr:         0,
		}

		assert.ErrorIs(t, v.Validate(), gofigure.ErrInvalidType)
	})
}

func TestValue_Assign(t *testing.T) {
	for _, target := range []any{
		new(bool), new(int), new(int64), new(uint), new(uint64), new(float64),
		new(time.Duration), new(gofigure.External),
	} {
		func(target any) {
			name := fmt.Sprintf("string(%T) assigns to %[1]T", target)
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				in := fmt.Sprint(gofigure.Dereference(target))

				value := gofigure.Value{
					Name:        fmt.Sprintf("%T", target),
					Description: fmt.Sprintf("%T type", target),
					Ptr:         target,
					Source:      gofigure.AllSources,
				}

				err := value.Assign(in, gofigure.Flag)

				assert.NoError(t, err)
			})
		}(target)
	}

	t.Run("External accepts External types", func(t *testing.T) {
		t.Parallel()

		var external gofigure.External

		value := gofigure.Value{
			Name:        "external",
			Description: "external",
			Ptr:         &external,
			Source:      gofigure.AllSources,
		}

		err := value.Assign(gofigure.External("ext"), gofigure.Flag)

		assert.NoError(t, err)
	})

	t.Run("External can't be an invalid type", func(t *testing.T) {
		t.Parallel()

		var external gofigure.External

		value := gofigure.Value{
			Name:        "external",
			Description: "external",
			Ptr:         &external,
			Source:      gofigure.AllSources,
		}

		err := value.Assign(1, gofigure.Flag)

		assert.ErrorIs(t, err, gofigure.ErrInvalidType)
	})

	t.Run("Values must be valid", func(t *testing.T) {
		t.Parallel()

		value := gofigure.Value{
			Name:        "external",
			Description: "external",
			Ptr:         nil,
			Source:      gofigure.AllSources,
		}

		err := value.Assign(nil, gofigure.Flag)

		assert.ErrorIs(t, err, gofigure.ErrNilPointer)
	})
}

func TestAssign(t *testing.T) {
	t.Run("types must match", func(t *testing.T) {
		t.Parallel()

		var setting string

		err := gofigure.Assign(&setting, 1)

		assert.ErrorIs(t, err, gofigure.ErrInvalidType)
	})

	t.Run("string must be coercible", func(t *testing.T) {
		t.Parallel()

		var setting int

		err := gofigure.Assign(&setting, "s")

		assert.ErrorIs(t, err, strconv.ErrSyntax)
	})
}
