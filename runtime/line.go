package runtime

import (
	"strings"
	"time"

	"github.com/drone/drone-runtime/engine"
)

// Line represents a line in the container logs.
type Line struct {
	Number    int    `json:"pos,omitempty"`
	Message   string `json:"out,omitempty"`
	Timestamp int64  `json:"time,omitempty"`
}

type lineWriter struct {
	num   int
	now   time.Time
	rep   *strings.Replacer
	state *State
	lines []*Line
}

func newWriter(state *State) *lineWriter {
	w := &lineWriter{}
	w.num = 0
	w.now = time.Now().UTC()
	w.state = state
	w.rep = newReplacer(state.Step.Secrets)
	return w
}

func (w *lineWriter) Write(p []byte) (n int, err error) {
	out := string(p)
	if w.rep != nil {
		out = w.rep.Replace(out)
	}

	line := &Line{
		Number:    w.num,
		Message:   out,
		Timestamp: int64(time.Since(w.now).Seconds()),
	}

	if w.state.hook.GotLine != nil {
		w.state.hook.GotLine(w.state, line)
	}
	w.num++

	w.lines = append(w.lines, line)
	return len(p), nil
}

func newReplacer(secrets []*engine.Secret) *strings.Replacer {
	var oldnew []string
	for _, secret := range secrets {
		if secret.Mask {
			oldnew = append(oldnew, secret.Value)
			oldnew = append(oldnew, "********")
		}
	}
	if len(oldnew) == 0 {
		return nil
	}
	return strings.NewReplacer(oldnew...)
}
