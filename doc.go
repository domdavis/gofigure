// Package gofigure allows for configuration of an application using command
// line flags, environment variables and configuration files.
//
// # Configuration Sources
//
// Configuration values are loaded from multiple sources with the following
// priority (highest to lowest):
//   - Command-line flags (Flag, ShortFlag)
//   - Environment variables (EnvVar)
//   - Configuration files (Key)
//   - Default values (Default)
//
// # Thread Safety
//
// Configuration instances are NOT safe for concurrent use. All configuration
// setup (calls to Group, Add, etc.) and parsing (Parse, ParseUsing) should be
// performed from a single goroutine, typically during application
// initialisation. Once parsed, the configuration values themselves can be
// safely read from multiple goroutines as they are stored in user-provided
// variables.
//
// # Panics
//
// Several functions in this package will panic on programmer errors rather than
// returning errors. This includes:
//   - NewValue: Panics if name or description is empty, or if ptr is nil/invalid
//   - NewParameters: Panics if name is empty
//
// These panics indicate programming errors (invalid arguments) rather than
// runtime errors, and should be caught during development and testing.
package gofigure
