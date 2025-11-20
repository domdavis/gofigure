package gofigure_test

import (
	"testing"

	"github.com/domdavis/gofigure"
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

func TestEnvironmentEdgeCases(t *testing.T) {
	t.Run("Very long environment variable value", func(t *testing.T) {
		var longValue string

		longString := make([]byte, 10000)
		for i := range longString {
			longString[i] = 'a'
		}
		expected := string(longString)

		t.Setenv("LONG_VAR", expected)

		settings := gofigure.Settings{
			gofigure.Required("long", "long-var", &longValue, gofigure.EnvVar, gofigure.ReportValue, "long var"),
		}

		m := gofigure.Environment("", settings)

		assert.NotEmpty(t, m)
	})

	t.Run("Unicode in environment variable value", func(t *testing.T) {
		var unicodeValue string

		expected := "Hello 世界 🌍 Здравствуй مرحبا"

		t.Setenv("UNICODE_VAR", expected)

		settings := gofigure.Settings{
			gofigure.Required("unicode", "unicode-var", &unicodeValue, gofigure.EnvVar, gofigure.ReportValue, "unicode var"),
		}

		m := gofigure.Environment("", settings)

		assert.NotEmpty(t, m)
	})

	t.Run("Empty environment variable value", func(t *testing.T) {
		var emptyValue string

		t.Setenv("EMPTY_VAR", "")

		settings := gofigure.Settings{
			gofigure.Required("empty", "empty-var", &emptyValue, gofigure.EnvVar, gofigure.ReportValue, "empty var"),
		}

		m := gofigure.Environment("", settings)

		// Empty environment variables may or may not be included in the map
		// depending on the implementation
		assert.NotNil(t, m)
	})

	t.Run("Environment variable with whitespace", func(t *testing.T) {
		var whitespaceValue string

		expected := "  spaces  and\ttabs\nand\rnewlines  "

		t.Setenv("WHITESPACE_VAR", expected)

		settings := gofigure.Settings{
			gofigure.Required("ws", "whitespace-var", &whitespaceValue, gofigure.EnvVar, gofigure.ReportValue, "whitespace var"),
		}

		m := gofigure.Environment("", settings)

		assert.NotEmpty(t, m)
	})

	t.Run("Multiple hyphens in environment variable name", func(t *testing.T) {
		var multiDashValue string

		t.Setenv("MULTI_DASH_VAR_NAME", "multi-dash")

		settings := gofigure.Settings{
			gofigure.Required("md", "multi-dash-var-name", &multiDashValue, gofigure.EnvVar, gofigure.ReportValue, "multi dash"),
		}

		m := gofigure.Environment("", settings)

		assert.NotEmpty(t, m)
	})
}

func FuzzEnvironmentValues(f *testing.F) {
	f.Add("test")
	f.Add("123")
	f.Add("true")
	f.Add("1h30m")
	f.Add("")
	f.Add("  spaces  ")
	f.Add("unicode-世界")
	f.Add("path/to/file")
	f.Add("key value")
	f.Add("!@#$%^&*()")

	f.Fuzz(func(t *testing.T, value string) {
		// Skip values that contain null bytes or other characters
		// that can't be set as environment variable values
		for i := 0; i < len(value); i++ {
			if value[i] == 0 || value[i] == '=' {
				t.Skip("skipping value with null byte or equals sign")
			}
		}

		var strValue string

		t.Setenv("FUZZ_VAR", value)

		settings := gofigure.Settings{
			gofigure.Optional("fuzz", "fuzz-var", &strValue, "", gofigure.EnvVar, gofigure.ReportValue, "fuzz test"),
		}

		_ = gofigure.Environment("", settings)
	})
}
