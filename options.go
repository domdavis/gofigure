package gofigure

import (
	"fmt"
	"sort"
	"strings"
)

// Options use to configure an application.
type Options map[Parameter]any

func (o Options) String() string {
	var i int

	keys := make([]Parameter, len(o))
	values := make([]string, len(o))

	for k := range o {
		keys[i] = k
		i++
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i].FullName() < keys[j].FullName()
	})

	for i, key := range keys {
		values[i] = fmt.Sprintf("%s:%v", key.FullName(), o[key])
	}

	return fmt.Sprintf("[%s]", strings.Join(values, ", "))
}
