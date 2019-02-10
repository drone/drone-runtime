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

import (
	"encoding/json"
	"io"
	"os"
	"strings"
)

// Parse parses the pipeline config from an io.Reader.
func Parse(r io.Reader) (*Spec, error) {
	cfg := Spec{}
	err := json.NewDecoder(r).Decode(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// ParseFile parses the pipeline config from a file.
func ParseFile(path string) (*Spec, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Parse(f)
}

// ParseString parses the pipeline config from a string.
func ParseString(s string) (*Spec, error) {
	return Parse(
		strings.NewReader(s),
	)
}
