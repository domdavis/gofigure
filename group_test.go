package gofigure_test

import (
	"fmt"

	"bitbucket.org/idomdavis/gofigure"
)

func ExampleGroup_Values() {
	var a, b string

	g := gofigure.Group{}

	g.Add(gofigure.Required("a", "a", &a, gofigure.Flag, gofigure.HideUnset, "a"))
	g.Add(gofigure.Required("b", "b", &b, gofigure.Flag, gofigure.MaskValue, "b"))

	fmt.Println(g.Values())

	// Output:
	// map[b:UNSET]
}
