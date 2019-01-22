package runtime

import (
	"testing"

	"github.com/drone/drone-runtime/engine"
)

func TestWithHooks(t *testing.T) {
	h := &Hook{}
	r := New(WithHooks(h))
	if r.hook != h {
		t.Errorf("Option does not set runtime hooks")
	}
}

func TestWithConfig(t *testing.T) {
	c := &engine.Spec{}
	r := New(WithConfig(c))
	if r.config != c {
		t.Errorf("Option does not set runtime configuration")
	}
}
