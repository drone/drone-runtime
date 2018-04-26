package runtime

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"github.com/drone/drone-runtime/engine"
	"github.com/vincent-petithory/dataurl"
	"golang.org/x/sync/errgroup"
)

// Runtime executes a pipeline configuration.
type Runtime struct {
	engine engine.Engine
	config *engine.Config
	hook   *Hook
	start  int64
	error  error
	fs     FileSystem
}

// New returns a new runtime using the specified runtime configuration
// and runtime engine.
func New(opts ...Option) *Runtime {
	r := &Runtime{}
	r.hook = &Hook{}
	for _, opts := range opts {
		opts(r)
	}
	return r
}

// Run starts the pipeline and waits for it to complete.
func (r *Runtime) Run(ctx context.Context) error {
	return r.Resume(ctx, 0)
}

// Resume starts the pipeline at the specified stage and waits
// for it to complete.
func (r *Runtime) Resume(ctx context.Context, start int) error {
	defer func() {
		// note that we use a new context to destroy the
		// environment to ensure it is not in a canceled
		// state.
		r.engine.Destroy(
			context.Background(), r.config)
	}()

	r.error = nil
	r.start = time.Now().Unix()

	if r.hook.Before != nil {
		state := snapshot(r, nil, nil)
		if err := r.hook.Before(state); err != nil {
			return err
		}
	}

	if err := r.engine.Setup(ctx, r.config); err != nil {
		return err
	}

	for i, stage := range r.config.Stages {
		if i < start {
			continue
		}
		select {
		case <-ctx.Done():
			return ErrCancel
		case err := <-r.execAll(stage.Steps):
			if err != nil {
				r.error = err
			}
		}
	}

	if r.hook.After != nil {
		state := snapshot(r, nil, nil)
		if err := r.hook.After(state); err != nil {
			return err
		}
	}
	return r.error
}

func (r *Runtime) execAll(group []*engine.Step) <-chan error {
	var g errgroup.Group
	done := make(chan error)

	for _, step := range group {
		step := step
		g.Go(func() error {
			return r.exec(step)
		})
	}

	go func() {
		done <- g.Wait()
		close(done)
	}()
	return done
}

func (r *Runtime) exec(step *engine.Step) error {
	ctx := context.TODO()

	switch {
	case r.error != nil && step.OnFailure == false:
		return nil
	case r.error == nil && step.OnSuccess == false:
		return nil
	}

	if r.hook.BeforeEach != nil {
		state := snapshot(r, step, nil)
		if err := r.hook.BeforeEach(state); err == ErrSkip {
			return nil
		} else if err != nil {
			return err
		}
	}

	if err := r.engine.Create(ctx, step); err != nil {
		return err
	}

	if r.fs != nil {
		state := snapshot(r, step, nil)
		if err := restoreAll(state); err != nil {
			return err
		}
	}

	if err := r.engine.Start(ctx, step); err != nil {
		return err
	}

	rc, err := r.engine.Tail(ctx, step)
	if err != nil {
		return err
	}

	var g errgroup.Group
	state := snapshot(r, step, nil)
	g.Go(func() error {
		return stream(state, rc)
	})

	if step.Detached {
		return nil // do not wait for service containers to complete.
	}

	defer func() {
		g.Wait() // wait for background tasks to complete.
		rc.Close()
	}()

	wait, err := r.engine.Wait(ctx, step)
	if err != nil {
		return err
	}

	if r.hook.GotFile != nil {
		state := snapshot(r, step, wait)
		g.Go(func() error {
			return exportAll(state)
		})
	}

	if r.fs != nil {
		state := snapshot(r, step, wait)
		g.Go(func() error {
			return backupAll(state)
		})
	}

	err = g.Wait() // wait for background tasks to complete.

	if wait.OOMKilled {
		err = &OomError{
			Name: step.Name,
			Code: wait.ExitCode,
		}
	} else if wait.ExitCode != 0 {
		err = &ExitError{
			Name: step.Name,
			Code: wait.ExitCode,
		}
	}

	if r.hook.AfterEach != nil {
		state := snapshot(r, step, wait)
		return r.hook.AfterEach(state)
	}

	if step.ErrIgnore {
		return nil
	}
	return err
}

// helper function exports a single file or folder.
func stream(state *State, rc io.ReadCloser) error {
	defer rc.Close()

	w := newWriter(state)
	io.Copy(w, rc)

	if state.hook.GotLogs != nil {
		return state.hook.GotLogs(state, w.lines)
	}
	return nil
}

// helper function exports files and folders in parallel.
func exportAll(state *State) error {
	var g errgroup.Group
	for _, file := range state.Step.Exports {
		file := file
		g.Go(func() error {
			return export(state, file)
		})
	}
	return g.Wait()
}

// helper function exports a single file or folder.
func export(state *State, file *engine.File) error {
	ctx := context.TODO()

	path := file.Path
	mime := file.Mime

	rc, info, err := state.engine.Download(ctx, state.Step, path)
	if err != nil {
		return err
	}
	defer rc.Close()
	info.Mime = mime
	return state.hook.GotFile(state, info, rc)
}

// helper function to backup files and folders in parallel.
func backupAll(state *State) error {
	var g errgroup.Group
	for _, b := range state.Step.Backup {
		b := b
		g.Go(func() error {
			return backup(state, b)
		})
	}
	return g.Wait()
}

// helper function to backup a single file or folder.
func backup(s *State, b *engine.Snapshot) error {
	ctx := context.TODO()

	src, _, err := s.engine.Download(ctx, s.Step, b.Source)
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := s.fs.Create(b.Target)
	if err != nil {
		return err
	}
	_, err = io.Copy(dst, src)
	src.Close()
	dst.Close()
	return err
}

// helper function to restore files and folders serially.
func restoreAll(state *State) error {
	for _, b := range state.Step.Restore {
		if err := restore(state, b); err != nil {
			return err
		}
	}
	return nil
}

// helper function to restore a single file or folder.
func restore(s *State, b *engine.Snapshot) error {
	ctx := context.TODO()

	var rc io.ReadCloser
	if strings.HasPrefix(b.Source, "data:") {
		u, err := dataurl.DecodeString(b.Source)
		if err != nil {
			return err
		}
		r := bytes.NewBuffer(u.Data)
		rc = ioutil.NopCloser(r)
	} else {
		var err error
		rc, err = s.fs.Open(b.Source)
		if err != nil {
			return err
		}
	}
	defer rc.Close()

	return s.engine.Upload(ctx, s.Step, b.Target, rc)
}
