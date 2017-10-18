package runtime

import (
	"testing"

	"github.com/drone/drone-runtime/engine"
	"github.com/drone/drone-runtime/runtime/mocks"
	"github.com/golang/mock/gomock"
)

func TestWithHooks(t *testing.T) {
	h := &Hook{}
	r := New(WithHooks(h))
	if r.hook != h {
		t.Errorf("Option does not set runtime hooks")
	}
}

func TestWithConfig(t *testing.T) {
	c := &engine.Config{}
	r := New(WithConfig(c))
	if r.config != c {
		t.Errorf("Option does not set runtime configuration")
	}
}

func TestWithFileSystem(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()

	fs := mocks.NewMockFileSystem(c)

	r := New(WithFileSystem(fs))
	if r.fs != fs {
		t.Errorf("Option does not set runtime virtual filesystem")
	}
}
