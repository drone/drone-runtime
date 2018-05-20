package runtime

import (
	"testing"

	"github.com/drone/drone-runtime/engine"
)

func TestLineWriter(t *testing.T) {
	line := &Line{}
	hook := &Hook{}
	state := &State{}

	hook.GotLine = func(_ *State, l *Line) error {
		line = l
		return nil
	}
	state.hook = hook
	state.Step = &engine.Step{}
	state.Step.Secrets = []*engine.Secret{
		{Name: "foo", Value: "bar", Mask: true},
	}

	newWriter(state).Write([]byte("foobar"))

	if line == nil {
		t.Error("Expect LineFunc invoked")
	}
	if got, want := line.Message, "foo********\n"; got != want {
		t.Errorf("Got line %q, want %q", got, want)
	}
	if got, want := line.Number, 0; got != want {
		t.Errorf("Got line %d, want %d", got, want)
	}
}

func TestMultiLineWriter(t *testing.T) {
	state := &State{}
	state.hook = &Hook{}
	state.Step = &engine.Step{}

	w := newWriter(state)

	written, err := w.Write([]byte("foo\nbar\n"))
	if err != nil {
		t.Errorf("Expect no error but got: %v", err)
	}
	if written != 8 {
		t.Errorf("Expect to write 8 chars but written: %d", written)
	}

	if len(w.lines) != 2 {
		t.Errorf("Expect 2 lines to be crated, got %d", len(w.lines))
	}

	expected := []Line{
		{Number: 0, Message: "foo\n", Timestamp: 0},
		{Number: 1, Message: "bar\n", Timestamp: 0},
	}
	for i, exp := range expected {
		if w.lines[i].Number != exp.Number {
			t.Errorf("Got line number %d, want %d", w.lines[i].Number, exp.Number)
		}
		if w.lines[i].Message != exp.Message {
			t.Errorf("Got line %s, want %s", w.lines[i].Message, exp.Message)
		}
		if w.lines[i].Timestamp != exp.Timestamp {
			t.Errorf("Got line timestamp %d, want %d", w.lines[i].Timestamp, exp.Timestamp)
		}
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
