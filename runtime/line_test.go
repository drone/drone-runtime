// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Drone Non-Commercial License
// that can be found in the LICENSE file.

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
	state.config = &engine.Spec{}
	state.config.Secrets = []*engine.Secret{
		{Metadata: engine.Metadata{Name: "foo"}, Data: "bar"},
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

func TestLineWriterSingle(t *testing.T) {
	line := &Line{}
	hook := &Hook{}
	state := &State{}

	hook.GotLine = func(_ *State, l *Line) error {
		line = l
		return nil
	}
	state.hook = hook
	state.Step = &engine.Step{}
	state.config = &engine.Spec{}

	lw := newWriter(state)
	lw.num = 5
	lw.Write([]byte("foo\n"))

	if line == nil {
		t.Error("Expect LineFunc invoked")
	}
	if got, want := line.Message, "foo\n"; got != want {
		t.Errorf("Got line %q, want %q", got, want)
	}
	if got, want := line.Number, 5; got != want {
		t.Errorf("Got line %d, want %d", got, want)
	}
}

func TestLineWriterMulti(t *testing.T) {
	var lines []*Line
	hook := &Hook{}
	state := &State{}

	hook.GotLine = func(_ *State, l *Line) error {
		lines = append(lines, l)
		return nil
	}
	state.hook = hook
	state.Step = &engine.Step{}
	state.config = &engine.Spec{}

	newWriter(state).Write([]byte("foo\nbar\nbaz"))

	if len(lines) != 3 {
		t.Error("Expect LineFunc invoked")
	}
	if got, want := lines[1].Message, "bar\n"; got != want {
		t.Errorf("Got line %q, want %q", got, want)
	}
	if got, want := lines[1].Number, 1; got != want {
		t.Errorf("Got line %d, want %d", got, want)
	}
}

func TestLineReplacer(t *testing.T) {
	secrets := []*engine.Secret{
		{Metadata: engine.Metadata{Name: "foo"}, Data: "bar"},
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

	secrets = []*engine.Secret{}
	replacer = newReplacer(secrets)
	if replacer != nil {
		t.Errorf("Expect nil replacer when no masked secrets")
	}
}

func TestLineCircling(t *testing.T) {
	hook := &Hook{}
	state := &State{}

	state.hook = hook
	state.Step = &engine.Step{}
	state.config = &engine.Spec{}
	state.config.Secrets = []*engine.Secret{
		{Metadata: engine.Metadata{Name: "foo"}, Data: "bar"},
	}

	w := newWriter(state)
	w.limit = 25
	w.Write([]byte("foobar1"))
	w.Write([]byte("foobar2"))
	w.Write([]byte("foobar3"))

	if len(w.lines) != 2 {
		t.Errorf("Got %d lines, want %d lines", len(w.lines), 2)
	}
	if got, want := w.lines[0].Message, "foo********2"; got != want {
		t.Errorf("Got line %q, want %q", got, want)
	}
	if got, want := w.lines[1].Message, "foo********3"; got != want {
		t.Errorf("Got line %q, want %q", got, want)
	}
}
