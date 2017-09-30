package runtime

import (
	"testing"

	"github.com/drone/drone-runtime/engine"
)

func TestLineWriter(t *testing.T) {
	var (
		line  *Line
		hook  = new(Hook)
		state = new(State)
	)
	hook.GotLine = func(_ *State, l *Line) error {
		line = l
		return nil
	}
	state.hook = hook
	state.Step = new(engine.Step)
	state.Step.Secrets = []*engine.Secret{
		{Name: "foo", Value: "bar", Mask: true},
	}

	newWriter(state).Write([]byte("foobar"))

	if line == nil {
		t.Error("Expect LineFunc invoked")
	}
	if got, want := line.Message, "foo********"; got != want {
		t.Errorf("Got line %q, want %q", got, want)
	}
	if got, want := line.Number, 0; got != want {
		t.Errorf("Got line %d, want %d", got, want)
	}
}

func TestLineReplacer(t *testing.T) {
	secrets := []*engine.Secret{
		{Name: "foo", Value: "bar", Mask: true},
	}
	replacer := newReplacer(secrets)
	if replacer == nil {
		t.Errorf("Expect non-nil replacer when masked secrets")
	}
	if got, want := replacer.Replace("foobar"), "foo********"; got != want {
		t.Errorf("Expect %q replaced with value %q", got, want)
	}

	// ensure the replacer is nil when the secret list is empty
	// or contains no masked secrets.

	secrets = []*engine.Secret{
		{Name: "foo", Value: "bar", Mask: false},
	}
	replacer = newReplacer(secrets)
	if replacer != nil {
		t.Errorf("Expect nil replacer when no masked secrets")
	}
}
