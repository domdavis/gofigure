package gofigure_test

import (
	"fmt"

	"bitbucket.org/idomdavis/gofigure"
)

func ExampleMask_String() {
	fmt.Println(gofigure.HideValue)
	fmt.Println(gofigure.MaskValue)
	fmt.Println(gofigure.DefaultIsSet)

	// Output:
	// Hide value
	// Mask value
	// Mask value: 16
}
