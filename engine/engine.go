package engine

//go:generate mockgen -source=engine.go -destination=mocks/engine.go
//-package=mocks -imports=.=github.com/drone/drone-runtime/engine
import (
	"context"
	"io"
)

// Factory defines a runtime engine factory.
type Factory interface {
	Create(*Spec) Engine
}

// Engine defines a runtime engine for pipeline execution.
type Engine interface {
	// Setup the pipeline environment.
	Setup(context.Context) error

	// Create creates the pipeline state.
	Create(context.Context, *Step) error

	// Start the pipeline step.
	Start(context.Context, *Step) error

	// Wait for the pipeline step to complete and returns the completion results.
	Wait(context.Context, *Step) (*State, error)

	// Tail the pipeline step logs.
	Tail(context.Context, *Step) (io.ReadCloser, error)

	// Destroy the pipeline environment.
	Destroy(context.Context) error
}
