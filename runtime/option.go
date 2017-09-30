package runtime

import (
	"github.com/drone/drone-runtime/engine"
)

// Option configures a Runtime option.
type Option func(*Runtime)

// WithEngine sets the Runtime engine.
func WithEngine(engine engine.Engine) Option {
	return func(r *Runtime) {
		r.engine = engine
	}
}

// WithConfig sets the Runtime configuration.
func WithConfig(c *engine.Config) Option {
	return func(r *Runtime) {
		r.config = c
	}
}

// WithHooks sets the Runtime tracer.
func WithHooks(h *Hook) Option {
	return func(r *Runtime) {
		if h != nil {
			r.hook = h
		}
	}
}

// WithFileSystem sets the Runtime virtual filesystem.
func WithFileSystem(fs FileSystem) Option {
	return func(r *Runtime) {
		r.fs = fs
	}
}
