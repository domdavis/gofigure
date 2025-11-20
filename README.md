# Go figure!

![Build](https://img.shields.io/github/actions/workflow/status/domdavis/gofigure/build.yml)
![Issues](https://img.shields.io/github/issues/domdavis/gofigure?style=plastic)
![Pull Requests](https://img.shields.io/github/issues-pr/domdavis/gofigure?style=plastic)
[![Go Doc](http://img.shields.io/badge/godoc-reference-5272B4.svg?style=plastic)](http://godoc.org/github.com/domdavis/gofigure)
[![License](https://img.shields.io/badge/license-MIT-green?style=plastic)](https://opensource.org/licenses/MIT)

A type-safe Go configuration library supporting multiple sources: command-line flags, environment variables, local 
config files, and remote HTTP config files. You define your configuration structure and say "Go figure!"

## Features

- **Multiple Configuration Sources**: Command-line flags, environment variables, JSON files (local and remote)
- **Type-Safe**: Uses Go generics for compile-time type safety
- **Source Priority**: Configurable precedence (flags > env vars > config files > defaults)
- **Flexible Display**: Hide or mask sensitive configuration values
- **Context Support**: Cancellation and timeout control for remote config loading
- **Comprehensive Security**: HTTP timeouts, size limits, proper error handling
- **100% Test Coverage**: Including fuzz testing for robustness

## Installation

```sh
go get -u github.com/domdavis/gofigure
```

Requires Go 1.18 or later (uses generics).

## Quick Start

```
package main

import (
    "fmt"
    "time"

    "github.com/domdavis/gofigure"
)

type Config struct {
    Host     string
    Port     int
    Timeout  time.Duration
    Debug    bool
    ApiKey   string
}

func main() {
    var cfg Config

    // Create configuration with defaults
    config := gofigure.New("myapp", "1.0.0")

    config.Add(
        gofigure.Optional("host", "host", &cfg.Host, "localhost",
            gofigure.Flag|gofigure.EnvVar, gofigure.ReportValue,
            "Server hostname"),
        gofigure.Optional("port", "port", &cfg.Port, 8080,
            gofigure.Flag|gofigure.EnvVar, gofigure.ReportValue,
            "Server port"),
        gofigure.Optional("timeout", "timeout", &cfg.Timeout, 30*time.Second,
            gofigure.Flag|gofigure.EnvVar, gofigure.ReportValue,
            "Request timeout"),
        gofigure.Optional("debug", "debug", &cfg.Debug, false,
            gofigure.Flag|gofigure.EnvVar, gofigure.ReportValue,
            "Enable debug mode"),
        gofigure.Required("api-key", "api-key", &cfg.ApiKey,
            gofigure.EnvVar, gofigure.MaskValue,
            "API authentication key"),
    )

    // Parse configuration from all sources
    if err := config.Parse(); err != nil {
        fmt.Println(config.Usage())
        panic(err)
    }

    // Configuration is now ready to use
    fmt.Printf("Server: %s:%d\n", cfg.Host, cfg.Port)
}
```

## Configuration Source Priority

Sources are checked in the following order (the highest to lowest priority):

1. **Command-line Flags** (`--flag`, `-f`)
2. **Environment Variables** (`MYAPP_VARIABLE`)
3. **Configuration Files** (JSON)
4. **Default Values**

The first source that provides a value wins. This allows you to set defaults in code, override with config files, and 
override those with environment variables or command-line flags.

## Usage Examples

### Command-Line Flags

```
// Support both long and short flags
gofigure.Optional("verbose", "v", &verbose, false,
    gofigure.Flag|gofigure.ShortFlag, gofigure.ReportValue,
    "Enable verbose output")

// Run with: ./myapp --verbose
// or:       ./myapp -v
```

### Environment Variables

```
// Automatically converts "api-key" to "MYAPP_API_KEY"
gofigure.Required("api-key", "api-key", &apiKey,
    gofigure.EnvVar, gofigure.MaskValue,
    "API authentication key")

// Run with: MYAPP_API_KEY=secret123 ./myapp
```

### JSON Configuration Files

Loading a local file:
```
options, err := gofigure.Load("config.json")
if err != nil {
    panic(err)
}
```

Loading a remote file with context:
```
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

options, err := gofigure.LoadWithContext(ctx, "https://config.example.com/app.json")
if err != nil {
    panic(err)
}
```

Example `config.json`:
```json
{
    "host": "api.example.com",
    "port": 443,
    "timeout": "1m30s",
    "debug": false
}
```

### Masking Sensitive Values

```
// Use MaskValue to hide sensitive data in logs/output
gofigure.Required("password", "password", &password,
    gofigure.EnvVar, gofigure.MaskValue,
    "Database password")

// Use HideIfSet to show placeholder when set
gofigure.Required("token", "token", &token,
    gofigure.EnvVar, gofigure.HideIfSet,
    "Authentication token")
```

When printed, masked values show as `***` and hidden values show as `<set>`.

## Security Considerations

### Remote Configuration Loading

When loading configuration from remote URLs:

- **Timeouts**: 30-second total timeout, 10-second connection timeout
- **Size Limits**: Maximum 10MB config file size
- **HTTPS Recommended**: Always use HTTPS for remote configs in production
- **Status Validation**: Only HTTP 200 responses are accepted
- **Context Support**: Use `LoadWithContext()` for cancellation control

```
// Production example with HTTPS and timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

options, err := gofigure.LoadWithContext(ctx, "https://secure-config.example.com/prod.json")
if err != nil {
    log.Fatalf("Failed to load config: %v", err)
}
```

### Sensitive Data Handling

- Always use `gofigure.MaskValue` for passwords, tokens, and API keys
- Never log configuration objects directly in production
- Use environment variables for secrets rather than config files
- Consider using dedicated secret management for production deployments

## Thread Safety

**Important**: Configuration instances are NOT safe for concurrent use. All configuration setup (`Group()`, `Add()`, 
etc.) and parsing (`Parse()`, `ParseUsing()`) must be performed from a single goroutine, typically during application 
initialisation.

```
// Good: Single-threaded initialization
func main() {
    config := gofigure.New("myapp", "1.0.0")
    config.Add(/* settings */)
    if err := config.Parse(); err != nil {
        panic(err)
    }

    // Config values can now be read concurrently
    go worker1(cfg)
    go worker2(cfg)
}

// Bad: Concurrent configuration setup
func main() {
    config := gofigure.New("myapp", "1.0.0")
    go config.Add(setting1)  // Race condition!
    go config.Add(setting2)  // Race condition!
}
```

## Error Handling

### Panics

The following functions panic on programmer errors (invalid arguments):

- `NewValue()` - panics if name or description are empty, or ptr is nil.
- `NewParameters()` - panics if name is empty.

These panics indicate bugs in your code and should be caught during development/testing.

### Runtime Errors

Configuration parsing returns errors for runtime issues:

```
if err := config.Parse(); err != nil {
    // Check for specific error types
    var configErr *gofigure.ConfigError
    if errors.As(err, &configErr) {
        fmt.Printf("Config error: %v\n", configErr)
        fmt.Println(config.Usage())
    }
    return err
}
```

## Documentation

- [GoDoc](https://godoc.org/github.com/domdavis/gofigure): Complete API reference
- [Examples](https://godoc.org/github.com/domdavis/gofigure#pkg-examples): Executable examples

## Contributing

Contributions are welcome! Please:

1. Write tests for new functionality
2. Run `go test ./...` and ensure all tests pass
3. Run `go test -fuzz=FuzzCoerceInt -fuzztime=10s` for fuzz testing
4. Follow existing code style and conventions
5. Update documentation as needed

## License

MIT Licence; see [LICENSE](LICENSE) file for details.
