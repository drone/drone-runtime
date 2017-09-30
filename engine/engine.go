package engine

//go:generate mockgen -source=engine.go -destination=mocks/engine.go -package=mocks -imports=.=github.com/drone/drone-runtime/engine

import "io"

// Engine defines a runtime engine for pipeline execution.
type Engine interface {
	// Setup the pipeline environment.
	Setup(*Config) error

	// Create creates the pipeline state.
	Create(*Step) error

	// Start the pipeline step.
	Start(*Step) error

	// Wait for the pipeline step to complete and returns the completion results.
	Wait(*Step) (*State, error)

	// Tail the pipeline step logs.
	Tail(*Step) (io.ReadCloser, error)

	// Upload uploads the file or data to the environment.
	Upload(*Step, string, io.Reader) error

	// Download downloads the file or data from the environment.
	Download(*Step, string) (io.ReadCloser, *FileInfo, error)

	// Destroy the pipeline environment.
	Destroy(*Config) error
}
