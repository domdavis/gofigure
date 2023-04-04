package gofigure

import (
	"os"
)

// Environment Options defined by the Settings.
func Environment(prefix string, settings Settings) Options {
	vars := Options{}

	for _, setting := range settings {
		for _, parameter := range setting.Parameters {
			if !parameter.Source.Contains(EnvVar) || parameter.Name == "" {
				continue
			}

			parameter.Stub = prefix

			if value := os.Getenv(parameter.FullName()); value != "" {
				vars[parameter] = value
			}
		}
	}

	return vars
}
