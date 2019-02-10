// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Drone Community
// License that can be found in the LICENSE file.

// +build !linux

package plugin

import (
	"errors"

	"github.com/drone/drone-runtime/engine"
)

// Symbol the symbol name used to lookup the plugin provider value.
const Symbol = "Engine"

// Open returns a Engine dynamically loaded from a plugin.
func Open(path string) (engine.Engine, error) {
	panic(
		errors.New("unsupported operating system"),
	)
}
