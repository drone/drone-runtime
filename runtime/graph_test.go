// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Drone Community
// License that can be found in the LICENSE file.

package runtime

import (
	"testing"

	"github.com/drone/drone-runtime/engine"
)

func TestIsSerial(t *testing.T) {
	spec := &engine.Spec{
		Steps: []*engine.Step{
			{Metadata: engine.Metadata{Name: "build"}},
			{Metadata: engine.Metadata{Name: "test"}},
		},
	}
	if isSerial(spec) == false {
		t.Errorf("Expect is serial true")
	}

	spec.Steps[1].DependsOn = []string{"build"}
	if isSerial(spec) == true {
		t.Errorf("Expect is serial false")
	}
}
