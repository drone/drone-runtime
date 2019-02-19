// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Drone Non-Commercial License
// that can be found in the LICENSE file.

package kube

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/drone/drone-runtime/engine"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// TODO(bradrydzewski) enable container resource limits.

// helper function converts environment variable
// string data to kubernetes variables.
func toEnv(spec *engine.Spec, step *engine.Step) []v1.EnvVar {
	var to []v1.EnvVar
	for k, v := range step.Envs {
		to = append(to, v1.EnvVar{
			Name:  k,
			Value: v,
		})
	}
	to = append(to, v1.EnvVar{
		Name: "KUBERNETES_NODE",
		ValueFrom: &v1.EnvVarSource{
			FieldRef: &v1.ObjectFieldSelector{
				FieldPath: "spec.nodeName",
			},
		},
	})
	for _, secret := range step.Secrets {
		sec, ok := engine.LookupSecret(spec, secret)
		if !ok {
			continue
		}
		optional := true
		to = append(to, v1.EnvVar{
			Name: secret.Env,
			ValueFrom: &v1.EnvVarSource{
				SecretKeyRef: &v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{
						Name: sec.Metadata.UID,
					},
					Key:      sec.Metadata.UID,
					Optional: &optional,
				},
			},
		})
	}
	return to
}

// helper function converts the engine pull policy
// to the kubernetes pull policy constant.
func toPullPolicy(from engine.PullPolicy) v1.PullPolicy {
	switch from {
	case engine.PullAlways:
		return v1.PullAlways
	case engine.PullNever:
		return v1.PullNever
	case engine.PullIfNotExists:
		return v1.PullIfNotPresent
	default:
		return v1.PullIfNotPresent
	}
}

// helper function converts the engine secret object
// to the kubernetes secret object.
func toSecret(spec *engine.Spec, from *engine.Secret) *v1.Secret {
	return &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: from.Metadata.UID,
		},
		Type: "Opaque",
		StringData: map[string]string{
			from.Metadata.UID: from.Data,
		},
	}
}

func toConfigVolumes(spec *engine.Spec, step *engine.Step) []v1.Volume {
	var to []v1.Volume
	for _, mount := range step.Files {
		file, ok := engine.LookupFile(spec, mount.Name)
		if !ok {
			continue
		}
		mode := int32(mount.Mode)
		volume := v1.Volume{Name: file.Metadata.UID}

		optional := false
		volume.ConfigMap = &v1.ConfigMapVolumeSource{
			LocalObjectReference: v1.LocalObjectReference{
				Name: file.Metadata.UID,
			},
			Optional: &optional,
			Items: []v1.KeyToPath{
				{
					Key:  file.Metadata.UID,
					Path: path.Base(mount.Path), // use the base path. document this.
					Mode: &mode,
				},
			},
		}
		to = append(to, volume)
	}
	return to
}

func toConfigMounts(spec *engine.Spec, step *engine.Step) []v1.VolumeMount {
	var to []v1.VolumeMount
	for _, mount := range step.Files {
		file, ok := engine.LookupFile(spec, mount.Name)
		if !ok {
			continue
		}
		volume := v1.VolumeMount{
			Name:      file.Metadata.UID,
			MountPath: path.Dir(mount.Path), // mount the config map here, using the base path
		}
		to = append(to, volume)
	}
	return to
}

func toVolumes(spec *engine.Spec, step *engine.Step) []v1.Volume {
	var to []v1.Volume
	for _, mount := range step.Volumes {
		vol, ok := engine.LookupVolume(spec, mount.Name)
		if !ok {
			continue
		}
		volume := v1.Volume{Name: vol.Metadata.UID}
		source := v1.HostPathDirectoryOrCreate
		if vol.HostPath != nil {
			volume.HostPath = &v1.HostPathVolumeSource{
				Path: vol.HostPath.Path,
				Type: &source,
			}
		}
		if vol.EmptyDir != nil {
			// volume.EmptyDir = &v1.EmptyDirVolumeSource{}

			// NOTE the empty_dir cannot be shared across multiple
			// pods so we emulate its behavior, and mount a temp
			// directory on the host machine that can be shared
			// between pods. This means we are responsible for deleting
			// these directories.
			volume.HostPath = &v1.HostPathVolumeSource{
				Path: filepath.Join("/tmp", "drone", spec.Metadata.Namespace, vol.Metadata.UID),
				Type: &source,
			}
		}
		to = append(to, volume)
	}
	return to
}

func toVolumeMounts(spec *engine.Spec, step *engine.Step) []v1.VolumeMount {
	var to []v1.VolumeMount
	for _, mount := range step.Volumes {
		vol, ok := engine.LookupVolume(spec, mount.Name)
		if !ok {
			continue
		}
		to = append(to, v1.VolumeMount{
			Name:      vol.Metadata.UID,
			MountPath: mount.Path,
		})
	}
	return to
}

func toPorts(step *engine.Step) []v1.ContainerPort {
	if len(step.Docker.Ports) == 0 {
		return nil
	}
	var ports []v1.ContainerPort
	for _, port := range step.Docker.Ports {
		ports = append(ports, v1.ContainerPort{
			ContainerPort: int32(port.Port),
		})
	}
	return ports
}

// helper function returns a kubernetes namespace
// for the given specification.
func toNamespace(spec *engine.Spec) *v1.Namespace {
	return &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   spec.Metadata.Namespace,
			Labels: spec.Metadata.Labels,
		},
	}
}

func toResources(step *engine.Step) v1.ResourceRequirements {
	var resources v1.ResourceRequirements
	if step.Resources != nil && step.Resources.Limits != nil {
		resources.Limits = v1.ResourceList{}
		if step.Resources.Limits.Memory > int64(0) {
			resources.Limits[v1.ResourceMemory] = *resource.NewQuantity(
				step.Resources.Limits.Memory, resource.BinarySI)
		}
		if step.Resources.Limits.CPU > int64(0) {
			resources.Limits[v1.ResourceCPU] = *resource.NewMilliQuantity(
				step.Resources.Limits.CPU, resource.DecimalSI)
		}
	}
	if step.Resources != nil && step.Resources.Requests != nil {
		resources.Requests = v1.ResourceList{}
		if step.Resources.Requests.Memory > int64(0) {
			resources.Requests[v1.ResourceMemory] = *resource.NewQuantity(
				step.Resources.Requests.Memory, resource.BinarySI)
		}
		if step.Resources.Requests.CPU > int64(0) {
			resources.Requests[v1.ResourceCPU] = *resource.NewMilliQuantity(
				step.Resources.Requests.CPU, resource.DecimalSI)
		}
	}
	return resources
}

// helper function returns a kubernetes pod for the
// given step and specification.
func toPod(spec *engine.Spec, step *engine.Step) *v1.Pod {
	var volumes []v1.Volume
	volumes = append(volumes, toVolumes(spec, step)...)
	volumes = append(volumes, toConfigVolumes(spec, step)...)

	var mounts []v1.VolumeMount
	mounts = append(mounts, toVolumeMounts(spec, step)...)
	mounts = append(mounts, toConfigMounts(spec, step)...)

	var pullSecrets []v1.LocalObjectReference
	if len(spec.Docker.Auths) > 0 {
		pullSecrets = []v1.LocalObjectReference{{
			Name: "docker-auth-config", // TODO move name to a const
		}}
	}

	return &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      step.Metadata.UID,
			Namespace: step.Metadata.Namespace,
			Labels:    step.Metadata.Labels,
		},
		Spec: v1.PodSpec{
			AutomountServiceAccountToken: boolptr(false),
			RestartPolicy:                v1.RestartPolicyNever,
			Containers: []v1.Container{{
				Name:            step.Metadata.UID,
				Image:           step.Docker.Image,
				ImagePullPolicy: toPullPolicy(step.Docker.PullPolicy),
				Command:         step.Docker.Command,
				Args:            step.Docker.Args,
				WorkingDir:      step.WorkingDir,
				SecurityContext: &v1.SecurityContext{
					Privileged: &step.Docker.Privileged,
				},
				Env:          toEnv(spec, step),
				VolumeMounts: mounts,
				Ports:        toPorts(step),
				Resources:    toResources(step),
			}},
			ImagePullSecrets: pullSecrets,
			Volumes:          volumes,
		},
	}
}

// helper function returns a kubernetes service for the
// given step and specification.
func toService(spec *engine.Spec, step *engine.Step) *v1.Service {
	var ports []v1.ServicePort
	for _, p := range step.Docker.Ports {
		source := p.Port
		target := p.Host
		if target == 0 {
			target = source
		}
		ports = append(ports, v1.ServicePort{
			Port: int32(source),
			TargetPort: intstr.IntOrString{
				IntVal: int32(target),
			},
		})
	}
	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      toDNS(step.Metadata.Name),
			Namespace: step.Metadata.Namespace,
		},
		Spec: v1.ServiceSpec{
			Type: v1.ServiceTypeClusterIP,
			Selector: map[string]string{
				"io.drone.step.name": step.Metadata.Name,
			},
			Ports: ports,
		},
	}
}

func toDNS(i string) string {
	return strings.Replace(i, "_", "-", -1)
}

func boolptr(v bool) *bool {
	return &v
}

func stringptr(v string) *string {
	return &v
}
