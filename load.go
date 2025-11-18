package gofigure

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

// ErrLoadingJSON is returned if an external file cannot be read from the path
// or URL given.
var ErrLoadingJSON = errors.New("error loading JSON")

// ErrParsingJSON is returned if the external file is not valid JSON.
var ErrParsingJSON = errors.New("error parsing JSON")

// ErrLoadingConfig is given to a ConfigError if Load fails.
var ErrLoadingConfig = errors.New("error loading config")

const (
	// Maximum size for remote configuration files (10MB).
	maxConfigSize = 10 * 1024 * 1024

	// HTTP client timeout for the entire request (30 seconds).
	httpTimeout = 30 * time.Second

	// Connection timeout (10 seconds).
	dialTimeout = 10 * time.Second

	// TLS handshake timeout (10 seconds).
	tlsTimeout = 10 * time.Second

	// Timeout for reading response headers (10 seconds).
	responseHeaderTimeout = 10 * time.Second
)

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
		return data, fmt.Errorf("%w from %q: %w", ErrLoadingJSON, uri, err)
	} else if err = json.Unmarshal(b, &data); err != nil {
		return data, fmt.Errorf("%w %q: %w", ErrParsingJSON, uri, err)
	}

	return data, nil
}

func get(uri string) ([]byte, error) {
	client := &http.Client{
		Timeout: httpTimeout,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   dialTimeout,
				KeepAlive: dialTimeout,
			}).DialContext,
			TLSHandshakeTimeout:   tlsTimeout,
			ResponseHeaderTimeout: responseHeaderTimeout,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)

	if err != nil {
		return []byte{}, fmt.Errorf(
			"failed to request external config from %s: %w", uri, err)
	}

	r, err := client.Do(req)

	if err != nil {
		return []byte{}, fmt.Errorf(
			"failed to get external config from %s: %w", uri, err)
	}

	defer func() { _ = r.Body.Close() }()

	if r.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf(
			"%w: config fetch from %s failed with status %d: %s",
			ErrLoadingConfig, uri, r.StatusCode, r.Status)
	}

	b, err := io.ReadAll(io.LimitReader(r.Body, maxConfigSize))

	if err != nil {
		return []byte{}, fmt.Errorf(
			"failed to read external config from %s: %w", uri, err)
	}

	if len(b) >= maxConfigSize {
		return []byte{}, fmt.Errorf(
			"%w: config file from %s exceeds maximum size of %d bytes",
			ErrLoadingConfig, uri, maxConfigSize)
	}

	return b, nil
}
