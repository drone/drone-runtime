// +build go1.8,linux

package plugin

import (
	"io"
	"plugin"

	"github.com/drone/drone-runtime/engine"
)

// Symbol the symbol name used to lookup the plugin provider value.
const Symbol = "Engine"

// Open returns an Engine dynamically loaded from a plugin.
func Open(path string) (engine.Engine, error) {
	lib, err := plugin.Open(path)
	if err != nil {
		return nil, err
	}
	provider, err := lib.Lookup(Symbol)
	if err != nil {
		return nil, err
	}
	return provider.(func() (engine.Engine, error))()
}

type pluginEngine struct {
	plugin *plugin.Plugin
}

func (p *pluginEngine) Setup(config *engine.Config) error {
	return nil
}

func (p *pluginEngine) Exec(*engine.Step) error {
	return nil
}

func (p *pluginEngine) Wait(*engine.Step) (*engine.State, error) {
	return nil, nil
}

func (p *pluginEngine) Tail(*engine.Step) (io.ReadCloser, error) {
	return nil, nil
}

func (p *pluginEngine) Copy(*engine.Step, string) (io.ReadCloser, *engine.FileInfo, error) {
	return nil, nil, nil
}

func (p *pluginEngine) Destroy(*engine.Config) error {
	return nil
}
