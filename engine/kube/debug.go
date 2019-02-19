// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Drone Non-Commercial License
// that can be found in the LICENSE file.

package kube

import (
	"bytes"

	"github.com/drone/drone-runtime/engine"
	"github.com/ghodss/yaml"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	documentBegin = "---\n"
	documentEnd   = "...\n"
)

// Print encodes returns specification as a Kubernetes
// multi-document yaml configuration file, in string format.
func Print(spec *engine.Spec) string {
	buf := new(bytes.Buffer)

	//
	// Secret Encoding.
	//

	for _, secret := range spec.Secrets {
		buf.WriteString(documentBegin)
		res := toSecret(spec, secret)
		res.Namespace = spec.Metadata.Namespace
		res.Kind = "Secret"
		res.Type = "Opaque"
		raw, _ := yaml.Marshal(res)
		buf.Write(raw)
	}

	//
	// Config Map Encoding.
	//

	for _, file := range spec.Files {
		res := &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name: file.Metadata.UID,
			},
			Data: map[string]string{
				file.Metadata.UID: string(file.Data),
			},
		}
		res.Namespace = spec.Metadata.Namespace
		res.Kind = "ConfigMap"
		buf.WriteString(documentBegin)
		raw, _ := yaml.Marshal(res)
		buf.Write(raw)
	}

	//
	// Step Encoding.
	//

	for _, step := range spec.Steps {
		buf.WriteString(documentBegin)
		res := toPod(spec, step)
		res.Namespace = spec.Metadata.Namespace
		res.Kind = "Pod"
		raw, _ := yaml.Marshal(res)
		buf.Write(raw)

		if len(step.Docker.Ports) != 0 {
			buf.WriteString(documentBegin)
			res := toService(spec, step)
			res.Namespace = spec.Metadata.Namespace
			res.Kind = "Service"
			raw, _ := yaml.Marshal(res)
			buf.Write(raw)
		}
	}

	buf.WriteString(documentEnd)
	return buf.String()
}
