package gofigure_test

import (
	"fmt"
	"testing"

	"bitbucket.org/idomdavis/gofigure"
	"github.com/stretchr/testify/assert"
)

func TestSetting_Display(t *testing.T) {
	const assignedValue = "value"

	for _, data := range []struct {
		Mask    gofigure.Mask
		Set     bool
		Display bool
		Expect  string
	}{
		{
			Mask:    gofigure.ReportValue,
			Set:     false,
			Display: true,
			Expect:  "",
		},
		{
			Mask:    gofigure.ReportValue,
			Set:     true,
			Display: true,
			Expect:  assignedValue,
		},
		{
			Mask:    gofigure.HideSet,
			Set:     false,
			Display: true,
			Expect:  "",
		},
		{
			Mask:    gofigure.HideSet,
			Set:     true,
			Display: false,
			Expect:  "",
		},
		{
			Mask:    gofigure.HideUnset,
			Set:     false,
			Display: false,
			Expect:  "",
		},
		{
			Mask:    gofigure.HideUnset,
			Set:     true,
			Display: true,
			Expect:  assignedValue,
		},
		{
			Mask:    gofigure.MaskSet,
			Set:     false,
			Display: true,
			Expect:  "",
		},
		{
			Mask:    gofigure.MaskSet,
			Set:     true,
			Display: true,
			Expect:  gofigure.Set,
		},
		{
			Mask:    gofigure.MaskUnset,
			Set:     false,
			Display: true,
			Expect:  gofigure.NotSet,
		},
		{
			Mask:    gofigure.MaskUnset,
			Set:     true,
			Display: true,
			Expect:  assignedValue,
		},
	} {
		func(mask gofigure.Mask, set, display bool, expect string) {
			t.Run(fmt.Sprintf("%s (Set: %t", mask, set), func(t *testing.T) {
				t.Parallel()

				var value string

				setting := gofigure.Required("setting", "setting", &value,
					gofigure.AllSources, mask, "Example Setting")

				if set {
					err := setting.Value.Assign(assignedValue, gofigure.Flag)

					assert.NoError(t, err)
				}

				v, ok := setting.Display()

				assert.Equal(t, display, ok)
				assert.Equal(t, expect, v)
			})
		}(data.Mask, data.Set, data.Display, data.Expect)
	}

	t.Run("An empty setting is invalid and wont display", func(t *testing.T) {
		t.Parallel()

		setting := gofigure.Setting{}
		display, ok := setting.Display()

		assert.False(t, ok)
		assert.Equal(t, gofigure.Invalid, display)
	})

	t.Run("An invalid setting displays as such", func(t *testing.T) {
		t.Parallel()

		setting := gofigure.Setting{Value: new(gofigure.Value)}
		display, ok := setting.Display()

		assert.True(t, ok)
		assert.Equal(t, gofigure.Invalid, display)
	})

	t.Run("Invalid settings respect masks", func(t *testing.T) {
		t.Parallel()

		setting := gofigure.Setting{Value: new(gofigure.Value), Mask: gofigure.HideValue}
		display, ok := setting.Display()

		assert.False(t, ok)
		assert.Equal(t, "", display)
	})
}

func TestSettings_Map(t *testing.T) {
	t.Run("Map will error on extra parameters", func(t *testing.T) {
		t.Parallel()

		var value int

		settings := gofigure.Settings{gofigure.Required("arg", "arg", &value,
			gofigure.AllSources, gofigure.ReportValue, "test setting")}
		parameter := gofigure.Parameter{Name: "extra", Source: gofigure.Key}
		options := map[gofigure.Parameter]any{parameter: "string"}

		assert.ErrorIs(t, settings.Map(options), gofigure.ErrUnexpectedArgument)
	})

	t.Run("Map will error if am option can't be mapped", func(t *testing.T) {
		t.Parallel()

		var value int

		settings := gofigure.Settings{gofigure.Required("arg", "arg", &value,
			gofigure.AllSources, gofigure.ReportValue, "test setting")}
		parameter := gofigure.Parameter{Name: "arg", Source: gofigure.Flag}
		options := map[gofigure.Parameter]any{parameter: "string"}

		assert.Error(t, settings.Map(options))
	})

	t.Run("Map will ignore given errors", func(t *testing.T) {
		t.Parallel()

		var value int

		settings := gofigure.Settings{gofigure.Required("arg", "arg", &value,
			gofigure.AllSources, gofigure.ReportValue, "test setting")}
		parameter := gofigure.Parameter{Name: "extra", Source: gofigure.Key}
		options := map[gofigure.Parameter]any{parameter: "string"}

		assert.NoError(t, settings.Map(options, gofigure.ErrUnexpectedArgument))
	})
}
