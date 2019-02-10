// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Drone Community
// License that can be found in the LICENSE file.

package engine

import "testing"

//
// File Lookup Tests
//

func TestLookupFile(t *testing.T) {
	want := &File{Metadata: Metadata{Name: "foo"}}
	spec := &Spec{
		Files: []*File{want},
	}
	got, ok := LookupFile(spec, "foo")
	if !ok {
		t.Errorf("Expect file found")
	}
	if got != want {
		t.Errorf("Expect file returned")
	}
}

func TestLookupFile_NotFound(t *testing.T) {
	want := &File{Metadata: Metadata{Name: "foo"}}
	spec := &Spec{
		Files: []*File{want},
	}
	got, ok := LookupFile(spec, "bar")
	if ok {
		t.Errorf("Expect file not found")
	}
	if got != nil {
		t.Errorf("Expect file not returned")
	}
}

//
// Secret Lookup Tests
//

func TestLookupSecret(t *testing.T) {
	want := &Secret{Metadata: Metadata{Name: "foo"}}
	spec := &Spec{
		Secrets: []*Secret{want},
	}
	got, ok := LookupSecret(spec, &SecretVar{Name: "foo"})
	if !ok {
		t.Errorf("Expect secret found")
	}
	if got != want {
		t.Errorf("Expect secret returned")
	}
}

func TestLookupSecret_NotFound(t *testing.T) {
	want := &Secret{Metadata: Metadata{Name: "foo"}}
	spec := &Spec{
		Secrets: []*Secret{want},
	}
	got, ok := LookupSecret(spec, &SecretVar{Name: "bar"})
	if ok {
		t.Errorf("Expect volume not found")
	}
	if got != nil {
		t.Errorf("Expect volume not returned")
	}
}

//
// Volume Lookup Tests
//

func TestLookupVolume(t *testing.T) {
	want := &Volume{Metadata: Metadata{Name: "foo"}}
	spec := &Spec{
		Docker: &DockerConfig{
			Volumes: []*Volume{want},
		},
	}
	got, ok := LookupVolume(spec, "foo")
	if !ok {
		t.Errorf("Expect volume found")
	}
	if got != want {
		t.Errorf("Expect volume returned")
	}
}

func TestLookupVolume_NotFound(t *testing.T) {
	volume := &Volume{Metadata: Metadata{Name: "foo"}}
	spec := &Spec{
		Docker: &DockerConfig{
			Volumes: []*Volume{volume},
		},
	}
	got, ok := LookupVolume(spec, "bar")
	if ok {
		t.Errorf("Expect volume not found")
	}
	if got != nil {
		t.Errorf("Expect volume not returned")
	}
}

func TestLookupVolume_NotDocker(t *testing.T) {
	_, ok := LookupVolume(&Spec{}, "foo")
	if ok {
		t.Fail()
	}
}

//
// Auth Lookup Tests
//

func TestLookupAuth(t *testing.T) {
	tests := []string{"docker.io", "index.docker.io", "https://index.docker.io/v1", "http://docker.io/v2"}
	for _, test := range tests {
		want := &DockerAuth{Address: test}
		spec := &Spec{
			Docker: &DockerConfig{
				Auths: []*DockerAuth{want},
			},
		}
		got, ok := LookupAuth(spec, "docker.io")
		if !ok {
			t.Errorf("Expect auth found for %s", test)
		}
		if got != want {
			t.Errorf("Expect auth returned for %s", test)
		}
	}
}

func TestLookupAuth_NotFound(t *testing.T) {
	want := &DockerAuth{Address: "gcr.io"}
	spec := &Spec{
		Docker: &DockerConfig{
			Auths: []*DockerAuth{want},
		},
	}
	got, ok := LookupAuth(spec, "docker.io")
	if ok {
		t.Errorf("Expect auth not found")
	}
	if got != nil {
		t.Errorf("Expect auth not returned")
	}
}

func TestLookupAuth_NotDocker(t *testing.T) {
	_, ok := LookupAuth(&Spec{}, "foo")
	if ok {
		t.Fail()
	}
}

func TestLookupAuth_InvalidRegistry(t *testing.T) {
	want := &DockerAuth{Address: "http://192.168.0.%31"}
	spec := &Spec{
		Docker: &DockerConfig{
			Auths: []*DockerAuth{want},
		},
	}
	_, ok := LookupAuth(spec, "192.168.0.%31")
	if ok {
		t.Fail()
	}
}
