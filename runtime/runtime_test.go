package runtime

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"testing"

	"github.com/drone/drone-runtime/engine"
	"github.com/drone/drone-runtime/engine/mocks"
	"github.com/golang/mock/gomock"
)

func TestRun(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()

	conf := &engine.Config{
		Stages: []*engine.Stage{
			{
				Name: "stage_0",
				Steps: []*engine.Step{
					{
						Name: "step_0",
						Exports: []*engine.File{
							{Path: "/etc/hosts", Mime: "text/plain"},
						},
						OnSuccess: true,
					},
				},
			},
		},
	}

	buf := ioutil.NopCloser(bytes.NewBufferString(""))

	state := new(engine.State)

	ctx := context.TODO()

	mock := mocks.NewMockEngine(c)
	mock.EXPECT().Setup(ctx, conf)
	mock.EXPECT().Destroy(ctx, conf)
	mock.EXPECT().Tail(ctx, conf.Stages[0].Steps[0]).Return(buf, nil)
	mock.EXPECT().Wait(ctx, conf.Stages[0].Steps[0]).Return(state, nil)
	mock.EXPECT().Create(ctx, conf.Stages[0].Steps[0])
	mock.EXPECT().Start(ctx, conf.Stages[0].Steps[0])

	run := New(
		WithEngine(mock),
		WithConfig(conf),
	)
	err := run.Run(context.Background())
	if err != nil {
		t.Error(err)
	}

	// TODO test Before
	// TODO test BeforeEach
	// TODO test After
	// TODO test AfterEach
	// TODO test GotFile
	// TODO test GotLine
	// TODO test GotLogs
}

// TestResume verifies the runtime resumes execution at the specified stage
// and skips previous stages.
func TestResume(t *testing.T) {
	t.Skip()
}

// TestRunOnSuccessTrue verifies the runtime executes a container if the
// OnSuccess flag is True and the pipeline is in a passing state.
func TestRunOnSuccessTrue(t *testing.T) {
	t.Skip()
}

// TestRunOnSuccessFalse verifies the runtime skips a container if the
// OnSuccess flag is False and the pipeline is in a passing state.
func TestRunOnSuccessFalse(t *testing.T) {
	t.Skip()
}

// TestRunOnFailureTrue verifies the runtime executes a container if the
// OnFailure flag is True and the pipeline is in a failing state.
func TestRunOnFailureTrue(t *testing.T) {
	t.Skip()
}

// TestRunOnFailureFalse verifies the runtime skips a container if the
// OnFailure flag is False and the pipeline is in a failing state.
func TestRunOnFailureFalse(t *testing.T) {
	t.Skip()
}

// TestRunDetached verifies the runtime executes a container in the background
// and does not wait for it to execute when the detached flag is true.
func TestRunDetached(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()

	conf := &engine.Config{
		Stages: []*engine.Stage{
			{
				Name: "stage_0",
				Steps: []*engine.Step{
					{Name: "step_0", OnSuccess: true, Detached: true},
				},
			},
		},
	}

	buf := ioutil.NopCloser(bytes.NewBufferString(""))

	ctx := context.TODO()

	mock := mocks.NewMockEngine(c)
	mock.EXPECT().Setup(ctx, conf)
	mock.EXPECT().Destroy(ctx, conf)
	mock.EXPECT().Tail(ctx, conf.Stages[0].Steps[0]).Return(buf, nil)
	mock.EXPECT().Create(ctx, conf.Stages[0].Steps[0])
	mock.EXPECT().Start(ctx, conf.Stages[0].Steps[0])

	run := New(
		WithEngine(mock),
		WithConfig(conf),
	)

	err := run.Run(context.Background())
	if err != nil {
		t.Error(err)
	}
}

// TestRunError verifies the runtime exits when the docker engine returns an
// error doing a routine operation, like waiting for a container to exit.
func TestRunError(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()

	conf := &engine.Config{
		Stages: []*engine.Stage{
			{
				Name: "stage_0",
				Steps: []*engine.Step{
					{Name: "step_0", OnSuccess: true},
				},
			},
		},
	}

	err := errors.New("dummy error")

	ctx := context.TODO()

	mock := mocks.NewMockEngine(c)
	mock.EXPECT().Setup(ctx, conf)
	mock.EXPECT().Destroy(ctx, conf)
	mock.EXPECT().Tail(ctx, conf.Stages[0].Steps[0]).Return(nil, err)
	mock.EXPECT().Create(ctx, conf.Stages[0].Steps[0])
	mock.EXPECT().Start(ctx, conf.Stages[0].Steps[0])
	run := New(
		WithEngine(mock),
		WithConfig(conf),
	)

	if got, want := run.Run(context.Background()), err; got != want {
		t.Error("Want Engine error returned from runtime")
	}
}

// TestRunErrorExit verifies the runtime exits when a step returns a non-zero
// exit code. The runtime must return an ExitError with the container exit code.
func TestRunErrorExit(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()

	conf := &engine.Config{
		Stages: []*engine.Stage{
			{
				Name: "stage_0",
				Steps: []*engine.Step{
					{Name: "step_0", OnSuccess: true},
				},
			},
		},
	}

	buf := ioutil.NopCloser(bytes.NewBufferString(""))

	state := &engine.State{
		ExitCode: 255,
	}

	ctx := context.TODO()

	mock := mocks.NewMockEngine(c)
	mock.EXPECT().Setup(ctx, conf)
	mock.EXPECT().Destroy(ctx, conf)
	mock.EXPECT().Tail(ctx, conf.Stages[0].Steps[0]).Return(buf, nil)
	mock.EXPECT().Wait(ctx, conf.Stages[0].Steps[0]).Return(state, nil)
	mock.EXPECT().Create(ctx, conf.Stages[0].Steps[0])
	mock.EXPECT().Start(ctx, conf.Stages[0].Steps[0])
	run := New(
		WithEngine(mock),
		WithConfig(conf),
	)

	err := run.Run(context.Background())
	if err == nil {
		t.Errorf("Want error returned from runtime, got nil")
	}
	errExit, ok := err.(*ExitError)
	if !ok {
		t.Errorf("Want ExitError returned from runtime")
		return
	}
	if got, want := errExit.Code, state.ExitCode; got != want {
		t.Errorf("Want exit code %d, got %d", want, got)
	}
	if got, want := errExit.Name, "step_0"; got != want {
		t.Errorf("Want step name %s, got %s", want, got)
	}
}

// TestRunErrorOom verifies the runtime exits when a step returns with
// out-of-memory killed is true. The runtime must return an OomError.
func TestRunErrorOom(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()

	conf := &engine.Config{
		Stages: []*engine.Stage{
			{
				Name: "stage_0",
				Steps: []*engine.Step{
					{Name: "step_0", OnSuccess: true},
				},
			},
		},
	}

	buf := ioutil.NopCloser(bytes.NewBufferString(""))

	state := &engine.State{
		OOMKilled: true,
		ExitCode:  255,
	}

	ctx := context.TODO()

	mock := mocks.NewMockEngine(c)
	mock.EXPECT().Setup(ctx, conf)
	mock.EXPECT().Destroy(ctx, conf)
	mock.EXPECT().Tail(ctx, conf.Stages[0].Steps[0]).Return(buf, nil)
	mock.EXPECT().Wait(ctx, conf.Stages[0].Steps[0]).Return(state, nil)
	mock.EXPECT().Create(ctx, conf.Stages[0].Steps[0])
	mock.EXPECT().Start(ctx, conf.Stages[0].Steps[0])
	run := New(
		WithEngine(mock),
		WithConfig(conf),
	)

	err := run.Run(context.Background())
	if err == nil {
		t.Errorf("Want error returned from runtime, got nil")
	}
	errOOM, ok := err.(*OomError)
	if !ok {
		t.Errorf("Want OomError returned from runtime")
		return
	}
	if got, want := errOOM.Name, "step_0"; got != want {
		t.Errorf("Want step name %s, got %s", want, got)
	}
}

// TestRunCancel verifies the runtime exits when context.Done and returns an
// ErrCancel. It also verifies the runtime exits immediately and does not
// execute additional steps.
func TestRunCancel(t *testing.T) {
	t.Skipf("this test panics when cancel() is invoked")
	t.SkipNow()

	c := gomock.NewController(t)
	defer c.Finish()

	conf := &engine.Config{
		Stages: []*engine.Stage{
			{
				Name: "stage_0",
				Steps: []*engine.Step{
					{Name: "step_0", OnSuccess: true},
				},
			},
		},
	}

	ctx := context.TODO()

	mock := mocks.NewMockEngine(c)
	mock.EXPECT().Setup(ctx, conf)
	mock.EXPECT().Destroy(ctx, conf)

	run := New(
		WithEngine(mock),
		WithConfig(conf),
	)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	err := run.Run(ctx)
	if err != ErrCancel {
		t.Errorf("Expect ErrCancel when context is cancelled, got %s", err)
	}
}
