package gofigure_test

import (
	"testing"

	"bitbucket.org/idomdavis/gofigure"
	"github.com/stretchr/testify/assert"
)

//nolint:paralleltest // Setting environment variables.
func TestEnvironment(t *testing.T) {
	t.Run("Stubs are correctly applied", func(t *testing.T) {
		m := gofigure.Environment("A", setupEnvironment(t))

		assert.Equal(t, "[A_TEST:stub]", m.String())
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
	t.Setenv("ATTEST", "no stub")

	return gofigure.Settings{
		gofigure.Required("a", "TEST", &a, gofigure.EnvVar, gofigure.ReportValue, "a"),
		gofigure.Required("b", "ATTEST", &b, gofigure.EnvVar, gofigure.ReportValue, "b"),
		gofigure.Required("c", "UNSET", &c, gofigure.EnvVar, gofigure.ReportValue, "c"),
		gofigure.Required("d", "FLAG", &d, gofigure.Flag, gofigure.ReportValue, "d"),
	}
}
