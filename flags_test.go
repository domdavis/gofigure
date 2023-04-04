package gofigure_test

import (
	"fmt"
	"testing"

	"bitbucket.org/idomdavis/gofigure"
	"github.com/stretchr/testify/assert"
)

func ExampleFlags() {
	flags, err := gofigure.Flags([]string{"-v", "--name", "example"})

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(flags)

	// Output:
	// [name:example, v:true]
}

func TestFlags(t *testing.T) {
	t.Run("A single boolean flag is allowed", func(t *testing.T) {
		t.Parallel()

		flags, err := gofigure.Flags([]string{"-v"})

		assert.NoError(t, err)
		assert.Len(t, flags, 1)
		assert.Equal(t, "true", flags[flag("v", gofigure.ShortFlag)])
	})

	t.Run("Multiple boolean flags are allows", func(t *testing.T) {
		t.Parallel()

		flags, err := gofigure.Flags([]string{"-a", "-b", "-c"})

		assert.NoError(t, err)
		assert.Len(t, flags, 3)
		assert.Equal(t, "true", flags[flag("a", gofigure.ShortFlag)])
		assert.Equal(t, "true", flags[flag("b", gofigure.ShortFlag)])
		assert.Equal(t, "true", flags[flag("c", gofigure.ShortFlag)])
	})

	t.Run("Single flags are allowed", func(t *testing.T) {
		t.Parallel()

		flags, err := gofigure.Flags([]string{"-a", "A"})

		assert.NoError(t, err)
		assert.Len(t, flags, 1)
		assert.Equal(t, "A", flags[flag("a", gofigure.ShortFlag)])
	})

	t.Run("Multiple flags are allows", func(t *testing.T) {
		t.Parallel()

		flags, err := gofigure.Flags([]string{"-a", "A", "-b", "B", "-c", "C"})

		assert.NoError(t, err)
		assert.Len(t, flags, 3)
		assert.Equal(t, "A", flags[flag("a", gofigure.ShortFlag)])
		assert.Equal(t, "B", flags[flag("b", gofigure.ShortFlag)])
		assert.Equal(t, "C", flags[flag("c", gofigure.ShortFlag)])
	})

	t.Run("Mixed flags are allowed", func(t *testing.T) {
		t.Parallel()

		flags, err := gofigure.Flags([]string{"-a", "A", "-b", "-c", "C"})

		assert.NoError(t, err)
		assert.Len(t, flags, 3)
		assert.Equal(t, "A", flags[flag("a", gofigure.ShortFlag)])
		assert.Equal(t, "true", flags[flag("b", gofigure.ShortFlag)])
		assert.Equal(t, "C", flags[flag("c", gofigure.ShortFlag)])
	})

	t.Run("A badly formatted arg set errors", func(t *testing.T) {
		t.Parallel()

		_, err := gofigure.Flags([]string{"-a", "A", "b", "-c", "C"})

		assert.ErrorIs(t, err, gofigure.ErrUnexpectedArgument)
	})
}

func flag(name string, source gofigure.Source) gofigure.Parameter {
	a := gofigure.Parameter{Name: name, Source: source}

	return a
}
