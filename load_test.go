package gofigure_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"bitbucket.org/idomdavis/gofigure"
	"github.com/stretchr/testify/assert"
)

func ExampleLoad() {
	d, err := gofigure.Load("testdata/config.json")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(d)

	// Output:
	// [address:localhost:8000, duration:1h, float:0, int:0, name:overridden]
}

func TestGet(t *testing.T) {
	t.Run("An error calling the target will be reported", func(t *testing.T) {
		t.Parallel()

		_, err := gofigure.Get("https://notfound")

		assert.ErrorIs(t, err, gofigure.ErrLoadingJSON)
	})

	t.Run("An invalid url will fail", func(t *testing.T) {
		t.Parallel()

		_, err := gofigure.Get("https://notfound\t")

		assert.ErrorIs(t, err, gofigure.ErrLoadingJSON)
	})

	t.Run("An invalid body will fail", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Length", "1")
			}))

		_, err := gofigure.Get(server.URL)

		assert.ErrorIs(t, err, gofigure.ErrLoadingJSON)
	})

	t.Run("Invalid JSON will fail", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(`{"key": `))
			}))

		_, err := gofigure.Get(server.URL)

		assert.ErrorIs(t, err, gofigure.ErrParsingJSON)
	})
}
