package gofigure_test

import (
	"fmt"
	"testing"

	"github.com/domdavis/gofigure"
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

	t.Run("Flags can include -'s", func(t *testing.T) {
		t.Parallel()

		flags, err := gofigure.Flags([]string{"-a-b"})

		assert.NoError(t, err)
		assert.Len(t, flags, 1)
		assert.Equal(t, "true", flags[flag("a-b", gofigure.ShortFlag)])
	})
}

func flag(name string, source gofigure.Source) gofigure.Parameter {
	a := gofigure.Parameter{Name: name, Source: source}

	return a
}

func FuzzParseFlags(f *testing.F) {
	f.Add([]byte("-v"))
	f.Add([]byte("--name value"))
	f.Add([]byte("-a A -b B"))
	f.Add([]byte("--flag=value"))
	f.Add([]byte("-x"))
	f.Add([]byte("--long-flag"))
	f.Add([]byte("-abc"))
	f.Add([]byte("--key value --other"))
	f.Add([]byte("-"))
	f.Add([]byte("--"))

	f.Fuzz(func(t *testing.T, input []byte) {
		args := []string{string(input)}

		_, _ = gofigure.Flags(args)
	})
}

func FuzzParseFlagsMultiple(f *testing.F) {
	f.Add("-v", "true")
	f.Add("--name", "value")
	f.Add("-a", "A")
	f.Add("--key", "")
	f.Add("-x", "--y")
	f.Add("--flag", "123")
	f.Add("-a-b", "test")

	f.Fuzz(func(t *testing.T, flag1, flag2 string) {
		args := []string{flag1, flag2}

		_, _ = gofigure.Flags(args)
	})
}
