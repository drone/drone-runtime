// Copyright 2019 Drone IO, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package engine

//go:generate mockgen -source=engine.go -destination=mocks/engine.go

import (
	"context"
	"io"
)

// Engine defines a runtime engine for pipeline execution.
type Engine interface {
	// Setup the pipeline environment.
	Setup(context.Context, *Spec) error

	// Create creates the pipeline state.
	Create(context.Context, *Spec, *Step) error

	// Start the pipeline step.
	Start(context.Context, *Spec, *Step) error

	// Wait for the pipeline step to complete and returns the completion results.
	Wait(context.Context, *Spec, *Step) (*State, error)

	// Tail the pipeline step logs.
	Tail(context.Context, *Spec, *Step) (io.ReadCloser, error)

	// Destroy the pipeline environment.
	Destroy(context.Context, *Spec) error
}
