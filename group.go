package gofigure

// Group of Settings.
type Group struct {
	Name     string
	Settings []*Setting
}

// Add a Setting to the Group.
func (g *Group) Add(setting *Setting) {
	g.Settings = append(g.Settings, setting)
}

// Values set on this group. Values will strip any Setting with a Mask that
// indicates it should be hidden.
func (g *Group) Values() map[string]any {
	values := map[string]any{}

	for _, setting := range g.Settings {
		if value, ok := setting.Display(); ok {
			values[setting.Value.Name] = value
		}
	}

	return values
}
