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
	"net/url"
	"strings"
)

// LookupVolume is a helper function that will lookup the
// named volume.
func LookupVolume(spec *Spec, name string) (*Volume, bool) {
	if spec.Docker == nil {
		return nil, false
	}
	for _, vol := range spec.Docker.Volumes {
		if vol.Metadata.Name == name {
			return vol, true
		}
	}
	return nil, false
}

// LookupSecret is a helper function that will lookup the
// named secret.
func LookupSecret(spec *Spec, secret *SecretVar) (*Secret, bool) {
	for _, sec := range spec.Secrets {
		if sec.Metadata.Name == secret.Name {
			return sec, true
		}
	}
	return nil, false
}

// LookupFile is a helper function that will lookup the
// named file.
func LookupFile(spec *Spec, name string) (*File, bool) {
	for _, file := range spec.Files {
		if file.Metadata.Name == name {
			return file, true
		}
	}
	return nil, false
}

// LookupAuth is a helper function that will lookup the
// docker credentials by hostname.
func LookupAuth(spec *Spec, domain string) (*DockerAuth, bool) {
	if spec.Docker == nil {
		return nil, false
	}
	for _, auth := range spec.Docker.Auths {
		host := auth.Address

		// the auth address could be a fully qualified
		// url in which case, we should parse so we can
		// extract the domain name.
		if strings.HasPrefix(host, "http://") ||
			strings.HasPrefix(host, "https://") {
			uri, err := url.Parse(auth.Address)
			if err != nil {
				continue
			}
			host = uri.Host
		}

		// we need to account for the legacy docker
		// index domain name, which should match the
		// normalized domain name.
		if host == "index.docker.io" {
			host = "docker.io"
		}

		if host == domain {
			return auth, true
		}
	}
	return nil, false
}
