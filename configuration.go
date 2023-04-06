package gofigure

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// Configuration for a program.
type Configuration struct {
	Help   bool
	Prefix string
	Groups []*Group

	groups   map[string]*Group
	external External
}

// Line in a Report. The values in the Line will respect Mask settings.
type Line struct {
	Name   string
	Values map[string]any
}

// Report holds the setting configuration in a way that can be reported to the
// user. Values in the Report will respect Mask settings.
type Report []Line

// ErrMissingRequiredOption is returned if a required option has not been set
// after parsing.
var ErrMissingRequiredOption = errors.New("missing required option")

const internalGroup = "Base Configuration"

// NewConfiguration returns a new Configuration set to use the given prefix for
// environment variables.
func NewConfiguration(envPrefix string) *Configuration {
	return &Configuration{Prefix: envPrefix}
}

// Group of Definitions for this Configuration.
func (c *Configuration) Group(name string) *Group {
	if c.groups == nil {
		c.groups = map[string]*Group{}
	}

	group, ok := c.groups[name]

	if !ok {
		group = &Group{Name: name}
		c.groups[name] = group
		c.Groups = append(c.Groups, group)
	}

	return group
}

// AddHelp will add a "help" flag to the set of options. If ShortFlag is set on
// the sources then a short flag of 'h' is also added. All other sources are
// ignored.
func (c *Configuration) AddHelp(sources Source) {
	use := Flag

	if sources.Contains(ShortFlag) {
		use |= ShortFlag
	}

	g := c.Group(internalGroup)
	g.Add(Optional("Help", "help", &c.Help, false, use, HideValue,
		"Display usage information"))
}

// AddConfigFile will add a "config" option to the set of options. If ShortFlag
// is set on the sources then a short flag of 'h' is also added. All other
// sources are ignored. The provided value will be used as a path or URI to load
// and external configuration file from.
func (c *Configuration) AddConfigFile(sources Source) {
	use := Flag

	if sources.Contains(ShortFlag) {
		use |= ShortFlag
	}

	g := c.Group(internalGroup)
	g.Add(Optional("Config File", "config", &c.external, "", use, MaskUnset,
		"Provide configuration from an external JSON file"))
}

// Report on the configuration, returning the values in a format that can be
// displayed to the user. Empty groups will be stripped from the Report. Values
// in the Report will respect the Mask setting.
func (c *Configuration) Report() Report {
	var report Report

	for _, group := range c.Groups {
		values := group.Values()

		if len(values) > 0 {
			report = append(report, Line{Name: group.Name, Values: values})
		}
	}

	return report
}

// Usage string for this set of Options.
func (c *Configuration) Usage() string {
	var options int

	b := strings.Builder{}
	b.WriteString("usage:\n")

	for _, group := range c.Groups {
		for _, setting := range group.Settings {
			var base string

			if len(setting.Parameters) == 0 {
				continue
			}

			b.WriteString("  ")
			b.WriteString(setting.Value.Name)
			b.WriteString(" ")
			b.WriteString(setting.Parameters.Format(c.Prefix))
			b.WriteString("\n    ")
			b.WriteString(setting.Value.Description)

			if setting.Value.base != nil && !setting.Mask.Contains(HideUnset) {
				base = fmt.Sprint(setting.Value.base)
			}

			if base != "" {
				b.WriteString(fmt.Sprintf(" (default: %v)", base))
			} else if setting.Value.base == nil {
				b.WriteString(" (required)")
			}

			b.WriteString("\n\n")

			options++
		}
	}

	if options == 0 {
		return "[no options]"
	}

	return b.String()
}

// Parse the Options.
func (c *Configuration) Parse() error {
	return c.ParseUsing(os.Args[1:])
}

// ParseUsing uses the given arguments as the set of command line arguments.
//
//nolint:cyclop // We're parsing a lot of sources here.
func (c *Configuration) ParseUsing(args []string) error {
	settings := Settings{}

	for _, group := range c.Groups {
		settings = append(settings, group.Settings...)
	}

	for i, setting := range settings {
		for j := range setting.Parameters {
			settings[i].Parameters[j].Stub = c.Prefix
		}
	}

	if flags, err := Flags(args); err != nil {
		return fmt.Errorf("could not parse command line arguments: %w", err)
	} else if err = settings.Map(flags); err != nil {
		return fmt.Errorf("invalid command line argument: %w", err)
	} else if err = settings.Map(Environment(c.Prefix, settings)); err != nil {
		return fmt.Errorf("invalid environment variable: %w", err)
	}

	for _, path := range settings.External() {
		if data, err := Load(path); err != nil {
			return fmt.Errorf("failed to load external configuration: %w", err)
		} else if err = settings.Map(data, ErrUnexpectedArgument); err != nil {
			return fmt.Errorf("invalid config value: %w", err)
		}
	}

	for _, setting := range settings {
		if setting.Value.Source == None {
			return NewConfigError(ErrMissingRequiredOption, fmt.Errorf("%w: %s",
				ErrMissingRequiredOption, setting.Parameters.Format(c.Prefix)),
				setting.Parameters...)
		}
	}

	return nil
}

// Format an error for user consumption. This will remove most of the technical
// details and leave a simple message as to why the configuration failed.
// Format should be used to report any errors to the user.
func (c *Configuration) Format(err error) string {
	var target ConfigError

	if errors.As(err, &target) {
		return target.Format(c.Prefix)
	}

	return err.Error()
}
