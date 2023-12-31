package gofigure_test

import (
	"fmt"
	"math"

	"github.com/domdavis/gofigure"
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
	fmt.Println(gofigure.None)
	fmt.Println(gofigure.Default)
	fmt.Println(gofigure.Key)
	fmt.Println(gofigure.EnvVar)
	fmt.Println(gofigure.ShortFlag)
	fmt.Println(gofigure.Flag)
	fmt.Println(gofigure.Reference)
	fmt.Println(gofigure.Source(math.MaxUint8))
	fmt.Println(gofigure.Source(0))

	// Output:
	// none
	// default value
	// config file key
	// environment variable
	// short flag
	// flag
	// reference value
	// config file
	// source
}
