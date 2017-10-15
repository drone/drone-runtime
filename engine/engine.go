package engine

//go:generate mockgen -source=engine.go -destination=mocks/engine.go -package=mocks -imports=.=github.com/drone/drone-runtime/engine

import (
	"context"
	"io"
)

// Engine defines a runtime engine for pipeline execution.
type Engine interface {
	// Setup the pipeline environment.
	Setup(context.Context, *Config) error

	// Create creates the pipeline state.
	Create(context.Context, *Step) error

	// Start the pipeline step.
	Start(context.Context, *Step) error

	// Wait for the pipeline step to complete and returns the completion results.
	Wait(context.Context, *Step) (*State, error)

	// Tail the pipeline step logs.
	Tail(context.Context, *Step) (io.ReadCloser, error)

	// Upload uploads the file or data to the environment.
	Upload(context.Context, *Step, string, io.Reader) error

	// Download downloads the file or data from the environment.
	Download(context.Context, *Step, string) (io.ReadCloser, *FileInfo, error)

	// Destroy the pipeline environment.
	Destroy(context.Context, *Config) error
}
