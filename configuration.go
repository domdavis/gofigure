package gofigure

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// Configuration for a program.
type Configuration struct {
	Help   bool
	Prefix string
	Groups []*Group

	groups map[string]*Group
}

// ErrMissingRequiredOption is returned if a required option has not been set
// after parsing.
var ErrMissingRequiredOption = errors.New("missing required option")

const internalGroup = "internal"

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

// AddHelp will add a "help" flag to the set of Options. If ShortFlag is set on
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

// Log the configuration using the given logger. Masks will be respected.
func (c *Configuration) Log(logger *logrus.Logger) {
	for _, group := range c.Groups {
		values := group.Values()

		if len(values) > 0 {
			logger.WithFields(values).Infof("%s Configuration", group.Name)
		}
	}
}

// Usage string for this set of Options.
func (c *Configuration) Usage() string {
	var options int

	b := strings.Builder{}
	b.WriteString("usage:\n")

	for _, group := range c.Groups {
		for _, setting := range group.Settings {
			if setting.Value.Source.Contains(None) {
				continue
			}

			b.WriteString("  ")
			b.WriteString(setting.Value.Name)
			b.WriteString(" ")
			b.WriteString(setting.Parameters.Format(c.Prefix))
			b.WriteString("\n    ")
			b.WriteString(setting.Value.Description)
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
