package gofigure_test

import (
	"testing"

	"bitbucket.org/idomdavis/gofigure"
	"github.com/stretchr/testify/assert"
)

//nolint:paralleltest // Setting environment variables.
func TestEnvironment(t *testing.T) {
	t.Run("Stubs and underscores are correctly applied", func(t *testing.T) {
		m := gofigure.Environment("A", setupEnvironment(t))

		assert.Equal(t, "[A_DASH_TEST:test, A_TEST:stub]", m.String())
	})

	t.Run("No stub is correctly handled", func(t *testing.T) {
		m := gofigure.Environment("", setupEnvironment(t))

		assert.Equal(t, "[ATTEST:no stub]", m.String())
	})
}

func setupEnvironment(t *testing.T) gofigure.Settings {
	t.Helper()

	var a, b, c, d string

	t.Setenv("A_TEST", "stub")
	t.Setenv("A_DASH_TEST", "test")
	t.Setenv("ATTEST", "no stub")

	return gofigure.Settings{
		gofigure.Required("a", "test", &a, gofigure.EnvVar, gofigure.ReportValue, "a"),
		gofigure.Required("b", "dash-test", &b, gofigure.EnvVar, gofigure.ReportValue, "b"),
		gofigure.Required("c", "attest", &b, gofigure.EnvVar, gofigure.ReportValue, "c"),
		gofigure.Required("d", "unset", &c, gofigure.EnvVar, gofigure.ReportValue, "d"),
		gofigure.Required("e", "flag", &d, gofigure.Flag, gofigure.ReportValue, "e"),
	}
}
