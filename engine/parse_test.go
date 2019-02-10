// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Drone Community
// License that can be found in the LICENSE file.

package engine

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParse(t *testing.T) {
	spec, err := ParseString(mockSpecJSON)
	if err != nil {
		t.Error(err)
		return
	}

	if diff := cmp.Diff(mockSpec, spec); diff != "" {
		t.Errorf("Unxpected Parse results")
		t.Log(diff)
	}

	_, err = ParseString("[]")
	if err == nil {
		t.Errorf("Want parse error, got nil")
	}

}

func TestParseFile(t *testing.T) {
	f, err := ioutil.TempFile(os.TempDir(), "drone")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(f.Name())
	f.WriteString(mockSpecJSON)
	f.Close()

	_, err = ParseFile(f.Name())
	if err != nil {
		t.Error(err)
		return
	}

	_, err = ParseFile("/tmp/this/path/does/not/exist")
	if err == nil {
		t.Errorf("Want parse error, got nil")
	}
}

func init() {
	// when the test package initializes, encode
	// the spec and snapshot the value.
	data, _ := json.Marshal(mockSpec)
	mockSpecJSON = string(data)
}

var mockSpecJSON string

// this is a sample runtime specification file.
var mockSpec = &Spec{
	Metadata: Metadata{
		UID:       "metadata.uid",
		Namespace: "metadata.namespace",
		Name:      "metadata.name",
		Labels: map[string]string{
			"metadata.labels.key": "metadata.labels.value",
		},
	},
	Platform: Platform{
		OS:      "platform.os",
		Arch:    "platform.arch",
		Version: "platform.version",
		Variant: "platform.variant",
	},
	Secrets: []*Secret{
		{
			Metadata: Metadata{Name: "secrets.1.name"},
			Data:     "secrets.1.data",
		},
	},
	Files: []*File{
		{
			Metadata: Metadata{Name: "files.1.name"},
			Data:     []byte("files.1.data"),
		},
	},
	Docker: &DockerConfig{
		Volumes: []*Volume{
			{
				Metadata: Metadata{
					UID:       "volumes.1.metadata.uid",
					Namespace: "volumes.1.metadata.namespace",
					Name:      "volumes.1.metadata.name",
					Labels: map[string]string{
						"volumes.1.metadata.labels.key": "volumes.1.metadata.labels.value",
					},
				},
				EmptyDir: &VolumeEmptyDir{},
			},
			{
				Metadata: Metadata{
					UID:       "volumes.2.metadata.uid",
					Namespace: "volumes.2.metadata.namespace",
					Name:      "volumes.2.metadata.name",
					Labels: map[string]string{
						"volumes.2.metadata.labels.key": "volumes.2.metadata.labels.value",
					},
				},
				HostPath: &VolumeHostPath{
					Path: "volumes.2.host.path",
				},
			},
		},
		Auths: []*DockerAuth{
			{
				Address:  "auths.1.address",
				Username: "auths.1.username",
				Password: "auths.1.password",
			},
		},
	},
	Steps: []*Step{
		{
			Metadata: Metadata{
				UID:       "steps.1.metadata.uid",
				Namespace: "steps.1.metadata.namespace",
				Name:      "steps.1.metadata.name",
				Labels: map[string]string{
					"steps.1.metadata.labesl.key": "steps.1.metadata.labels.value",
				},
			},
			Detach:    true,
			DependsOn: []string{"steps.1.depends_on.1"},
			Docker: &DockerStep{
				Args:     []string{"steps.1.args.1"},
				Command:  []string{"steps.1.command.1"},
				Image:    "steps.1.image",
				Networks: []string{"steps.1.network"},
				Ports: []*Port{
					{
						Port:     3306,
						Host:     3307,
						Protocol: "TPC",
					},
				},
				Privileged: true,
				PullPolicy: PullIfNotExists,
			},
			Envs: map[string]string{
				"steps.1.envs.key": "steps.1.envs.value",
			},
			Files: []*FileMount{
				{
					Name: "steps.1.files.1.name",
					Path: "steps.1.files.1.path",
				},
			},
			IgnoreErr:    true,
			IgnoreStdout: true,
			IgnoreStderr: true,
			Resources:    &Resources{},
			RunPolicy:    RunAlways,
			Secrets: []*SecretVar{
				{
					Name: "steps.1.secrets.1.name",
					Env:  "steps.1.secrets.1.env",
				},
			},
			Volumes: []*VolumeMount{
				{
					Name: "steps.1.volumes.1.name",
					Path: "steps.1.volumes.1.path",
				},
			},
			WorkingDir: "steps.1.working_dir",
		},
	},
}
