package gofigure_test

import (
	"testing"

	"bitbucket.org/idomdavis/gofigure"
	"github.com/stretchr/testify/assert"
)

func TestConfiguration_Format(t *testing.T) {
	t.Run("Unexpected command line arguments are reported", func(t *testing.T) {
		t.Parallel()

		config := gofigure.NewConfiguration("")
		err := config.ParseUsing([]string{"/invalid"})
		msg := config.Format(err)

		assert.Equal(t, "unexpected argument: [/invalid]", msg)
	})

	t.Run("Extra flags are reported", func(t *testing.T) {
		t.Parallel()

		config := gofigure.NewConfiguration("")
		err := config.ParseUsing([]string{"--extra"})
		msg := config.Format(err)

		assert.Equal(t, "unexpected argument: [--extra]", msg)
	})

	//nolint:paralleltest // Testing environment variables.
	t.Run("Invalid environment variables are reported", func(t *testing.T) {
		var value int

		t.Setenv("GOFIGURE_TEST_VALUE", "string")
		config := gofigure.NewConfiguration("GOFIGURE_TEST")
		config.Group("test").Add(gofigure.Required("value", "value", &value,
			gofigure.EnvVar, gofigure.ReportValue, "test value"))

		err := config.ParseUsing([]string{})
		msg := config.Format(err)

		assert.Equal(t, "invalid value 'string': [env GOFIGURE_TEST_VALUE]", msg)
	})

	t.Run("Load errors are reported", func(t *testing.T) {
		t.Parallel()

		var file gofigure.External
		config := gofigure.NewConfiguration("")
		config.Group("test").Add(gofigure.Required("config", "c", &file,
			gofigure.ShortFlag, gofigure.ReportValue, "config file"))

		err := config.ParseUsing([]string{"-c", "/does/not/exist"})
		msg := config.Format(err)

		assert.Equal(t, "error loading config: [file: /does/not/exist]", msg)
	})

	t.Run("JSON parse errors are reported", func(t *testing.T) {
		t.Parallel()

		var (
			value int
			file  gofigure.External
		)

		config := gofigure.NewConfiguration("")
		group := config.Group("test")

		group.Add(gofigure.Required("config", "c", &file,
			gofigure.ShortFlag, gofigure.ReportValue, "config file"))
		group.Add(gofigure.Required("name", "name", &value,
			gofigure.Key, gofigure.ReportValue, "name from json"))

		err := config.ParseUsing([]string{"-c", "testdata/config.json"})
		msg := config.Format(err)

		assert.Equal(t, "invalid value 'overridden': [JSON key: \"name\"]", msg)
	})

	t.Run("Missing parameters are reported", func(t *testing.T) {
		t.Parallel()

		var value int

		config := gofigure.NewConfiguration("TEST")
		config.Group("test").Add(gofigure.Required("required", "required", &value,
			gofigure.AllSources, gofigure.ReportValue, "required value"))

		err := config.ParseUsing([]string{})
		msg := config.Format(err)

		assert.Equal(t, "missing required option: [JSON key: \"required\", "+
			"env TEST_REQUIRED, -r, --required]", msg)
	})

	t.Run("Standard errors are just reported as is", func(t *testing.T) {
		t.Parallel()

		config := gofigure.NewConfiguration("")
		msg := config.Format(assert.AnError)

		assert.Equal(t, assert.AnError.Error(), msg)
	})
}

func TestConfiguration_Parse(t *testing.T) {
	t.Run("Parse will use the OS arguments", func(t *testing.T) {
		t.Parallel()

		var value string

		config := gofigure.NewConfiguration("TEST")
		group := config.Group("settings")

		group.Add(gofigure.Required("Value", "unset", &value, 0, 0, "Test Parse"))

		assert.Error(t, config.Parse())
	})
}

func TestConfiguration_Usage(t *testing.T) {
	t.Run("Empty options are displayed correctly", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "[no options]", gofigure.NewConfiguration("").Usage())
	})

	t.Run("Empty groups are ignored", func(t *testing.T) {
		t.Parallel()

		var value string

		c := gofigure.NewConfiguration("")
		g := c.Group("empty")
		g.Add(gofigure.Required("opt", "opt", &value, gofigure.None,
			gofigure.HideValue, "hidden option"))

		assert.Equal(t, "[no options]", c.Usage())
	})
}

func TestConfiguration_Report(t *testing.T) {
	t.Run("Settings will be correctly logged", func(t *testing.T) {
		t.Parallel()

		var a, b, c int

		config := gofigure.NewConfiguration("")
		config.AddHelp(gofigure.CommandLine)

		group := config.Group("settings")
		group.Add(gofigure.Required("A", "a", &a, gofigure.ShortFlag,
			gofigure.ReportValue, "A setting"))
		group.Add(gofigure.Required("B", "b", &b, gofigure.ShortFlag,
			gofigure.MaskSet, "B setting"))
		group.Add(gofigure.Required("C", "c", &c, gofigure.ShortFlag,
			gofigure.HideSet, "C setting"))

		err := config.ParseUsing([]string{"-a", "1", "-b", "2", "-c", "3"})

		assert.NoError(t, err)

		report := config.Report()

		assert.Len(t, report, 1)
		assert.Equal(t, report[0].Values["A"], "1")
		assert.Equal(t, report[0].Values["B"], "SET")
		assert.Nil(t, report[0].Values["C"])
	})
}
