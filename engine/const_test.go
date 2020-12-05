// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Drone Non-Commercial License
// that can be found in the LICENSE file.

package engine

import (
	"bytes"
	"encoding/json"
	"testing"
)

//
// runtime policy unit tests.
//

func TestRunPolicy_Marshal(t *testing.T) {
	tests := []struct {
		policy RunPolicy
		data   string
	}{
		{
			policy: RunAlways,
			data:   `"always"`,
		},
		{
			policy: RunOnFailure,
			data:   `"on-failure"`,
		},
		{
			policy: RunOnSuccess,
			data:   `"on-success"`,
		},
		{
			policy: RunNever,
			data:   `"never"`,
		},
	}
	for _, test := range tests {
		data, err := json.Marshal(&test.policy)
		if err != nil {
			t.Error(err)
			return
		}
		if bytes.Equal([]byte(test.data), data) == false {
			t.Errorf("Failed to marshal policy %s", test.policy)
		}
	}
}

func TestRunPolicy_Unmarshal(t *testing.T) {
	tests := []struct {
		policy RunPolicy
		data   string
	}{
		{
			policy: RunAlways,
			data:   `"always"`,
		},
		{
			policy: RunOnFailure,
			data:   `"on-failure"`,
		},
		{
			policy: RunOnSuccess,
			data:   `"on-success"`,
		},
		{
			policy: RunNever,
			data:   `"never"`,
		},
		{
			// no policy should default to on-success
			policy: RunOnSuccess,
			data:   `""`,
		},
	}
	for _, test := range tests {
		var policy RunPolicy
		err := json.Unmarshal([]byte(test.data), &policy)
		if err != nil {
			t.Error(err)
			return
		}
		if got, want := policy, test.policy; got != want {
			t.Errorf("Want policy %q, got %q", want, got)
		}
	}
}

func TestRunPolicy_UnmarshalTypeError(t *testing.T) {
	var policy RunPolicy
	err := json.Unmarshal([]byte("[]"), &policy)
	if _, ok := err.(*json.UnmarshalTypeError); !ok {
		t.Errorf("Expect unmarshal error return when JSON invalid")
	}
}

func TestRunPolicy_String(t *testing.T) {
	tests := []struct {
		policy RunPolicy
		value  string
	}{
		{
			policy: RunAlways,
			value:  "always",
		},
		{
			policy: RunOnFailure,
			value:  "on-failure",
		},
		{
			policy: RunOnSuccess,
			value:  "on-success",
		},
	}
	for _, test := range tests {
		if got, want := test.policy.String(), test.value; got != want {
			t.Errorf("Want policy string %q, got %q", want, got)
		}
	}
}

//
// pull policy unit tests.
//

func TestPullPolicy_Marshal(t *testing.T) {
	tests := []struct {
		policy PullPolicy
		data   string
	}{
		{
			policy: PullAlways,
			data:   `"always"`,
		},
		{
			policy: PullDefault,
			data:   `"default"`,
		},
		{
			policy: PullIfNotExists,
			data:   `"if-not-exists"`,
		},
		{
			policy: PullNever,
			data:   `"never"`,
		},
	}
	for _, test := range tests {
		data, err := json.Marshal(&test.policy)
		if err != nil {
			t.Error(err)
			return
		}
		if bytes.Equal([]byte(test.data), data) == false {
			t.Errorf("Failed to marshal policy %s", test.policy)
		}
	}
}

func TestPullPolicy_Unmarshal(t *testing.T) {
	tests := []struct {
		policy PullPolicy
		data   string
	}{
		{
			policy: PullAlways,
			data:   `"always"`,
		},
		{
			policy: PullDefault,
			data:   `"default"`,
		},
		{
			policy: PullIfNotExists,
			data:   `"if-not-exists"`,
		},
		{
			policy: PullNever,
			data:   `"never"`,
		},
		{
			// no policy should default to on-success
			policy: PullDefault,
			data:   `""`,
		},
	}
	for _, test := range tests {
		var policy PullPolicy
		err := json.Unmarshal([]byte(test.data), &policy)
		if err != nil {
			t.Error(err)
			return
		}
		if got, want := policy, test.policy; got != want {
			t.Errorf("Want policy %q, got %q", want, got)
		}
	}
}

func TestPullPolicy_UnmarshalTypeError(t *testing.T) {
	var policy PullPolicy
	err := json.Unmarshal([]byte("[]"), &policy)
	if _, ok := err.(*json.UnmarshalTypeError); !ok {
		t.Errorf("Expect unmarshal error return when JSON invalid")
	}
}

func TestPullPolicy_String(t *testing.T) {
	tests := []struct {
		policy PullPolicy
		value  string
	}{
		{
			policy: PullAlways,
			value:  "always",
		},
		{
			policy: PullDefault,
			value:  "default",
		},
		{
			policy: PullIfNotExists,
			value:  "if-not-exists",
		},
		{
			policy: PullNever,
			value:  "never",
		},
	}
	for _, test := range tests {
		if got, want := test.policy.String(), test.value; got != want {
			t.Errorf("Want policy string %q, got %q", want, got)
		}
	}
}

//
// volume host path type unit tests.
//

func TestVolumeHostPathType_Marshal(t *testing.T) {
	tests := []struct {
		hostPathType VolumeHostPathType
		data         string
	}{
		{
			hostPathType: HostPathDirectoryOrCreate,
			data:         `"dir-or-create"`,
		},
		{
			hostPathType: HostPathDirectory,
			data:         `"path-dir"`,
		},
		{
			hostPathType: HostPathFileOrCreate,
			data:         `"file-or-create"`,
		},
		{
			hostPathType: HostPathFile,
			data:         `"file"`,
		},
		{
			hostPathType: HostPathSocket,
			data:         `"socket"`,
		},
		{
			hostPathType: HostPathCharDev,
			data:         `"char-dev"`,
		},
		{
			hostPathType: HostPathBlockDev,
			data:         `"block-dev"`,
		},
	}
	for _, test := range tests {
		data, err := json.Marshal(&test.hostPathType)
		if err != nil {
			t.Error(err)
			return
		}
		if bytes.Equal([]byte(test.data), data) == false {
			t.Errorf("Failed to marshal host path type %s", test.hostPathType)
		}
	}
}

func TestRunVolumeHostPathType_Unmarshal(t *testing.T) {
	tests := []struct {
		hostPathType VolumeHostPathType
		data         string
	}{
		{
			hostPathType: HostPathDirectoryOrCreate,
			data:         `"dir-or-create"`,
		},
		{
			hostPathType: HostPathDirectory,
			data:         `"path-dir"`,
		},
		{
			hostPathType: HostPathFileOrCreate,
			data:         `"file-or-create"`,
		},
		{
			hostPathType: HostPathFile,
			data:         `"file"`,
		},
		{
			hostPathType: HostPathSocket,
			data:         `"socket"`,
		},
		{
			hostPathType: HostPathCharDev,
			data:         `"char-dev"`,
		},
		{
			hostPathType: HostPathBlockDev,
			data:         `"block-dev"`,
		},
	}
	for _, test := range tests {
		var hostPathType VolumeHostPathType
		err := json.Unmarshal([]byte(test.data), &hostPathType)
		if err != nil {
			t.Error(err)
			return
		}
		if got, want := hostPathType, test.hostPathType; got != want {
			t.Errorf("Want host path type %q, got %q", want, got)
		}
	}
}

func TestVolumeHostPathType_UnmarshalTypeError(t *testing.T) {
	var hostPathType VolumeHostPathType
	err := json.Unmarshal([]byte("[]"), &hostPathType)
	if _, ok := err.(*json.UnmarshalTypeError); !ok {
		t.Errorf("Expect unmarshal error return when JSON invalid")
	}
}

func TestVolumeHostPathType_String(t *testing.T) {
	tests := []struct {
		hostPathType VolumeHostPathType
		value        string
	}{
		{
			hostPathType: HostPathDirectoryOrCreate,
			value:        "dir-or-create",
		},
		{
			hostPathType: HostPathDirectory,
			value:        "path-dir",
		},
		{
			hostPathType: HostPathFileOrCreate,
			value:        "file-or-create",
		},
		{
			hostPathType: HostPathFile,
			value:        "file",
		},
		{
			hostPathType: HostPathSocket,
			value:        "socket",
		},
		{
			hostPathType: HostPathCharDev,
			value:        "char-dev",
		},
		{
			hostPathType: HostPathBlockDev,
			value:        "block-dev",
		},
	}
	for _, test := range tests {
		if got, want := test.hostPathType.String(), test.value; got != want {
			t.Errorf("Want host path type string %q, got %q", want, got)
		}
	}
}
