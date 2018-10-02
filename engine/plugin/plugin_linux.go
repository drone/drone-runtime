// +build go1.8,linux

package plugin

import (
	"plugin"

	"github.com/drone/drone-runtime/engine"
)

// Symbol the symbol name used to lookup the plugin provider value.
const Symbol = "Factory"

// Open returns a Factory dynamically loaded from a plugin.
func Open(path string) (engine.Factory, error) {
	lib, err := plugin.Open(path)
	if err != nil {
		return nil, err
	}
	provider, err := lib.Lookup(Symbol)
	if err != nil {
		return nil, err
	}
	return provider.(func() (engine.Factory, error))()
}
