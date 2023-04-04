package gofigure_test

import (
	"fmt"
	"math"

	"bitbucket.org/idomdavis/gofigure"
)

func ExampleSource_Contains() {
	sources := gofigure.Flag | gofigure.ShortFlag

	fmt.Println(sources.Contains(gofigure.Flag))
	fmt.Println(sources.Contains(gofigure.Key))

	sources = gofigure.AllSources

	fmt.Println(sources.Contains(gofigure.Key))

	// Output:
	// true
	// false
	// true
}

func ExampleSource_String() {
	fmt.Println(gofigure.Default)
	fmt.Println(gofigure.Key)
	fmt.Println(gofigure.EnvVar)
	fmt.Println(gofigure.ShortFlag)
	fmt.Println(gofigure.Flag)
	fmt.Println(gofigure.Source(math.MaxUint8))
	fmt.Println(gofigure.None)

	// Output:
	// default value
	// config file key
	// environment value
	// short flag
	// flag
	// config file
	// source
}
