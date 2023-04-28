package gofigure

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// ErrLoadingJSON is returned if an external file cannot be read from the path
// or URL given.
var ErrLoadingJSON = errors.New("error loading JSON")

// ErrParsingJSON is returned if the external file is not valid JSON.
var ErrParsingJSON = errors.New("error parsing JSON")

// ErrLoadingConfig is given to a ConfigError if Load fails.
var ErrLoadingConfig = errors.New("error loading config")

// Load external Options from a URI. The external file can be any JSON
// object.
func Load(uri string) (Options, error) {
	options := Options{}

	data, err := Get(uri)

	if err != nil {
		return options, NewConfigError(ErrLoadingConfig,
			fmt.Errorf("failed to get external config: %w", err),
			Parameter{Name: uri, Source: configFile})
	}

	for k, v := range data {
		options[Parameter{Name: k, Source: Key}] = v
	}

	return options, nil
}

// Get a JSON object from an external source.
func Get(uri string) (map[string]any, error) {
	var data map[string]any

	f := os.ReadFile

	if strings.HasPrefix(uri, "http://") || strings.HasPrefix(uri, "https://") {
		f = get
	}

	if b, err := f(uri); err != nil {
		return data, fmt.Errorf("%w from %q: %s", ErrLoadingJSON, uri, err.Error())
	} else if err = json.Unmarshal(b, &data); err != nil {
		return data, fmt.Errorf("%w %q: %s", ErrParsingJSON, uri, err.Error())
	}

	return data, nil
}

func get(uri string) ([]byte, error) {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)

	if err != nil {
		return []byte{}, fmt.Errorf(
			"failed to request external config from %s: %w", uri, err)
	}

	r, err := http.DefaultClient.Do(req)

	if err != nil {
		return []byte{}, fmt.Errorf(
			"failed to get external config from %s: %w", uri, err)
	}

	//nolint:errcheck // Not a huge amount we can do here.
	defer func() { _ = r.Body.Close() }()

	b, err := io.ReadAll(r.Body)

	if err != nil {
		return []byte{}, fmt.Errorf(
			"failed to read external config from %s: %w", uri, err)
	}

	return b, nil
}
