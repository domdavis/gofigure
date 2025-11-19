package gofigure_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/domdavis/gofigure"
	"github.com/stretchr/testify/assert"
)

func ExampleLoad() {
	d, err := gofigure.Load("testdata/config.json")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(d)

	// Output:
	// [address:localhost:8000, duration:1h, float:12.3, int:12, name:overridden]
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

	t.Run("HTTP 404 status will fail", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			}))

		_, err := gofigure.Get(server.URL)

		assert.ErrorIs(t, err, gofigure.ErrLoadingJSON)
		assert.ErrorIs(t, err, gofigure.ErrLoadingConfig)
	})

	t.Run("HTTP 500 status will fail", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			}))

		_, err := gofigure.Get(server.URL)

		assert.ErrorIs(t, err, gofigure.ErrLoadingJSON)
		assert.ErrorIs(t, err, gofigure.ErrLoadingConfig)
	})

	t.Run("Config exceeding maximum size will fail", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				data := make([]byte, 11*1024*1024)
				for i := range data {
					data[i] = 'a'
				}
				_, _ = w.Write(data)
			}))

		_, err := gofigure.Get(server.URL)

		assert.ErrorIs(t, err, gofigure.ErrLoadingJSON)
		assert.ErrorIs(t, err, gofigure.ErrLoadingConfig)
	})

	t.Run("Cancelled context will fail", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"key": "value"}`))
			}))

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := gofigure.GetWithContext(ctx, server.URL)

		assert.ErrorIs(t, err, gofigure.ErrLoadingJSON)
		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("Context timeout will fail", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(100 * time.Millisecond)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"key": "value"}`))
			}))

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		_, err := gofigure.GetWithContext(ctx, server.URL)

		assert.ErrorIs(t, err, gofigure.ErrLoadingJSON)
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})

	t.Run("Deeply nested JSON is parsed correctly", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{
					"level1": {
						"level2": {
							"level3": {
								"level4": {
									"level5": "deep value"
								}
							}
						}
					},
					"key": "value"
				}`))
			}))

		data, err := gofigure.Get(server.URL)

		assert.NoError(t, err)
		assert.NotNil(t, data)
		assert.Contains(t, data, "key")
		assert.Contains(t, data, "level1")
	})

	t.Run("JSON with null values is handled", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{
					"nullValue": null,
					"stringValue": "test",
					"numberValue": 42
				}`))
			}))

		data, err := gofigure.Get(server.URL)

		assert.NoError(t, err)
		assert.NotNil(t, data)
		assert.Nil(t, data["nullValue"])
		assert.Equal(t, "test", data["stringValue"])
		assert.Equal(t, float64(42), data["numberValue"])
	})

	t.Run("JSON with arrays is parsed", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{
					"items": [1, 2, 3],
					"names": ["alice", "bob"],
					"key": "value"
				}`))
			}))

		data, err := gofigure.Get(server.URL)

		assert.NoError(t, err)
		assert.NotNil(t, data)
		assert.Contains(t, data, "items")
		assert.Contains(t, data, "names")
		assert.Contains(t, data, "key")
	})

	t.Run("JSON with mixed types is handled", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{
					"string": "text",
					"integer": 42,
					"float": 3.14,
					"boolean": true,
					"null": null,
					"object": {"nested": "value"},
					"array": [1, 2, 3]
				}`))
			}))

		data, err := gofigure.Get(server.URL)

		assert.NoError(t, err)
		assert.NotNil(t, data)
		assert.Equal(t, "text", data["string"])
		assert.Equal(t, float64(42), data["integer"])
		assert.Equal(t, 3.14, data["float"])
		assert.Equal(t, true, data["boolean"])
		assert.Nil(t, data["null"])
	})

	t.Run("Empty JSON object is handled", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{}`))
			}))

		data, err := gofigure.Get(server.URL)

		assert.NoError(t, err)
		assert.NotNil(t, data)
		assert.Empty(t, data)
	})

	t.Run("HTTP redirects are followed", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/redirect" {
					http.Redirect(w, r, "/target", http.StatusFound)

					return
				}
				if r.URL.Path == "/target" {
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write([]byte(`{"redirected": true}`))

					return
				}
				w.WriteHeader(http.StatusNotFound)
			}))

		data, err := gofigure.Get(server.URL + "/redirect")

		assert.NoError(t, err)
		assert.NotNil(t, data)
		assert.Equal(t, true, data["redirected"])
	})

	t.Run("Multiple HTTP redirects are followed", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/redirect1":
					http.Redirect(w, r, "/redirect2", http.StatusMovedPermanently)
				case "/redirect2":
					http.Redirect(w, r, "/target", http.StatusFound)
				case "/target":
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write([]byte(`{"hops": 2}`))
				default:
					w.WriteHeader(http.StatusNotFound)
				}
			}))

		data, err := gofigure.Get(server.URL + "/redirect1")

		assert.NoError(t, err)
		assert.NotNil(t, data)
		assert.Equal(t, float64(2), data["hops"])
	})
}
