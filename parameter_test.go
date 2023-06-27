package gofigure_test

import (
	"fmt"
	"testing"

	"github.com/domdavis/gofigure"
	"github.com/stretchr/testify/assert"
)

func ExampleNewParameters() {
	fmt.Println(gofigure.NewParameters("param", gofigure.AllSources))

	// Output:
	// [JSON key: "param" env PARAM -p --param]
}

func ExampleParameter_Matches() {
	a := gofigure.Parameter{Name: "param", Source: gofigure.AllSources}
	b := gofigure.Parameter{Name: "param", Source: gofigure.Flag}
	c := gofigure.Parameter{Name: "param", Source: gofigure.ShortFlag}

	fmt.Println(a.Matches(b))
	fmt.Println(b.Matches(c))
	fmt.Println(b.Matches(a))
	fmt.Println(c.Matches(b))

	// Output:
	// true
	// false
	// true
	// false
}

func ExampleParameters_Format() {
	p := gofigure.NewParameters("param", gofigure.AllSources)

	fmt.Println(p.Format("STUB"))

	// Output:
	// [JSON key: "param", env STUB_PARAM, -p, --param]
}

func TestNewParameters(t *testing.T) {
	t.Run("Parameter must have a name", func(t *testing.T) {
		t.Parallel()

		assert.Panics(t, func() { gofigure.NewParameters("", 0) })
	})
}
