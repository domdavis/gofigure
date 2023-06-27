package gofigure_test

import (
	"fmt"
	"time"

	"github.com/domdavis/gofigure"
)

func Example() {
	var settings struct {
		Mode    int
		Name    string
		Address string
		Timeout time.Duration
		TLS     bool

		Duration time.Duration
		Int      int
		Float    float64
	}

	config := gofigure.NewConfiguration("EXAMPLE")
	config.AddHelp(gofigure.CommandLine)
	config.AddConfigFile(gofigure.CommandLine)

	group := config.Group("settings")

	group.Add(gofigure.Required("App Name", "name", &settings.Name,
		gofigure.AllSources, gofigure.ReportValue, "Application name"))
	group.Add(gofigure.Required("Mode", "mode", &settings.Mode, gofigure.Flag,
		gofigure.ReportValue, "Mode indicator"))
	group.Add(gofigure.Optional("IP Address", "address", &settings.Address,
		"", gofigure.AllSources, gofigure.MaskUnset, "Remote server address"))
	group.Add(gofigure.Optional("Timeout", "timeout", &settings.Timeout,
		time.Minute, gofigure.AllSources, gofigure.ReportValue,
		"Remote server address"))
	group.Add(gofigure.Optional("TLS", "tls", &settings.TLS, false,
		gofigure.NamedSources, gofigure.ReportValue, "Use TLS"))

	types := config.Group("JSON types")

	types.Add(gofigure.Required("Duration", "duration", &settings.Duration,
		gofigure.Key, gofigure.ReportValue, "Duration type"))
	types.Add(gofigure.Required("Int", "int", &settings.Int,
		gofigure.Key, gofigure.ReportValue, "int type"))
	types.Add(gofigure.Required("Float", "float", &settings.Float,
		gofigure.Key, gofigure.ReportValue, "float type"))

	// Ordinarily this would be config.Parse().
	err := config.ParseUsing([]string{
		"-c", "testdata/config.json",
		"--mode", "3",
		"--name", "example",
		"-h",
	})

	if err != nil {
		fmt.Println(config.Format(err))
		fmt.Println(config.Usage())
	}

	fmt.Printf("%v\n\n", settings)

	if config.Help {
		fmt.Println(config.Usage())
	}

	// Output:
	// {3 example localhost:8000 1m0s false 1h0m0s 0 0}
	//
	// usage:
	//   Help [-h, --help]
	//     Display usage information
	//
	//   Config File [-c, --config]
	//     Provide configuration from an external JSON file
	//
	//   App Name [JSON key: "name", env EXAMPLE_NAME, -n, --name]
	//     Application name (required)
	//
	//   Mode [--mode]
	//     Mode indicator (required)
	//
	//   IP Address [JSON key: "address", env EXAMPLE_ADDRESS, -a, --address]
	//     Remote server address
	//
	//   Timeout [JSON key: "timeout", env EXAMPLE_TIMEOUT, -t, --timeout]
	//     Remote server address (default: 1m0s)
	//
	//   TLS [JSON key: "tls", env EXAMPLE_TLS, --tls]
	//     Use TLS (default: false)
	//
	//   Duration [JSON key: "duration"]
	//     Duration type (required)
	//
	//   Int [JSON key: "int"]
	//     int type (required)
	//
	//   Float [JSON key: "float"]
	//     float type (required)
}
