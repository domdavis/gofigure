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

// ErrLoadingConfig is returned if a config file cannot be read from the
// path or URL given.
var ErrLoadingConfig = errors.New("error loading config")

// ErrParsingConfig is returned if the config file is not valid JSON.
var ErrParsingConfig = errors.New("error parsing config")

// Load external Options from a URI. The external file can be any JSON
// object.
func Load(uri string) (Options, error) {
	var data map[string]any

	options := Options{}
	f := os.ReadFile

	if strings.HasPrefix(uri, "http://") || strings.HasPrefix(uri, "https://") {
		f = get
	}

	if b, err := f(uri); err != nil {
		return options, NewConfigError(ErrLoadingConfig,
			fmt.Errorf("%w from %q: %s", ErrLoadingConfig, uri, err.Error()),
			Parameter{Name: uri, Source: configFile})
	} else if err = json.Unmarshal(b, &data); err != nil {
		return options, NewConfigError(ErrLoadingConfig,
			fmt.Errorf("%w %q: %s", ErrParsingConfig, uri, err.Error()),
			Parameter{Name: uri, Source: configFile})
	}

	for k, v := range data {
		options[Parameter{Name: k, Source: Key}] = v
	}

	return options, nil
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
