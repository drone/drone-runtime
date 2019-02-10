// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Drone Community
// License that can be found in the LICENSE file.

package runtime

import "testing"

func TestExitError(t *testing.T) {
	err := ExitError{
		Name: "build",
		Code: 255,
	}
	got, want := err.Error(), "build : exit code 255"
	if got != want {
		t.Errorf("Want error message %q, got %q", want, got)
	}
}

func TestOomError(t *testing.T) {
	err := OomError{
		Name: "build",
	}
	got, want := err.Error(), "build : received oom kill"
	if got != want {
		t.Errorf("Want error message %q, got %q", want, got)
	}
}
