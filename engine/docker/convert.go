package docker

import (
	"strings"

	"github.com/drone/drone-runtime/engine"

	"docker.io/go-docker/api/types/container"
	"docker.io/go-docker/api/types/mount"
	"docker.io/go-docker/api/types/network"
)

// returns a container configuration.
func toConfig(spec *engine.Spec, step *engine.Step) *container.Config {
	config := &container.Config{
		Image:        step.Docker.Image,
		Labels:       step.Metadata.Labels,
		WorkingDir:   step.WorkingDir,
		AttachStdin:  false,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          false,
		OpenStdin:    false,
		StdinOnce:    false,
		ArgsEscaped:  false,
	}

	if len(step.Envs) != 0 {
		config.Env = toEnv(step.Envs)
	}
	for _, name := range step.Secrets {
		secret, ok := engine.LookupSecret(spec, name)
		if ok {
			config.Env = append(config.Env, secret.Name+"="+secret.Data)
		}
	}
	if len(step.Docker.Args) != 0 {
		config.Cmd = step.Docker.Args
	}
	if len(step.Docker.Command) != 0 {
		config.Entrypoint = step.Docker.Command
	}
	if len(step.Volumes) != 0 {
		config.Volumes = toVolumeSet(spec, step)
	}
	return config
}

// TODO: add dns
// TODO: add dns_search
// TODO: add extra_hosts
// TODO: add shmsize
// TODO: add tmpfs
// TODO: add devices
// TODO: set resource limits
// TODO: set port bindings
// returns a container host configuration.
func toHostConfig(spec *engine.Spec, proc *engine.Step) *container.HostConfig {
	config := &container.HostConfig{
		// TODO: map resources to proc.Resources.Limits
		// Resources: container.Resources{},
		LogConfig: container.LogConfig{
			Type: "json-file",
		},
		Privileged: proc.Docker.Privileged,
		// TODO: set shmsize for docker-based (e.g. non-kubernetes) installs
		// ShmSize: 0,
	}

	// TODO: set DNS and host settings for docker-based (e.g. non-kubernetes) installs
	// config.DNS
	// config.DNSSearch
	// config.ExtraHosts

	if len(proc.Devices) != 0 {
		// 	config.Devices = toDevices(proc.Devices)
	}
	if len(proc.Volumes) != 0 {
		config.Binds = toVolumeSlice(spec, proc)
		config.Mounts = toVolumeMounts(spec, proc)
	}
	// config.Tmpfs = map[string]string{}
	// for _, path := range proc.Tmpfs {
	// 	if strings.Index(path, ":") == -1 {
	// 		config.Tmpfs[path] = ""
	// 		continue
	// 	}
	// 	parts := strings.Split(path, ":")
	// 	config.Tmpfs[parts[0]] = parts[1]
	// }

	return config
}

func toNetConfig(spec *engine.Spec, proc *engine.Step) *network.NetworkingConfig {
	endpoints := map[string]*network.EndpointSettings{}
	endpoints[spec.Metadata.UID] = &network.EndpointSettings{
		NetworkID: spec.Metadata.UID,
		Aliases:   []string{proc.Metadata.Name},
	}
	return &network.NetworkingConfig{
		EndpointsConfig: endpoints,
	}
}

// helper function that converts a slice of volume paths to a set of
// unique volume names.
func toVolumeSet(spec *engine.Spec, step *engine.Step) map[string]struct{} {
	to := map[string]struct{}{}
	for _, mount := range step.Volumes {
		volume, ok := engine.LookupVolume(spec, mount.Name)
		if !ok {
			continue
		}
		if volume.HostPath != nil {
			if strings.HasPrefix(volume.HostPath.Path, `\\.\pipe\`) {
				continue
			}
		}
		to[mount.Path] = struct{}{}
	}
	return to
}

func toVolumeSlice(spec *engine.Spec, step *engine.Step) []string {
	var to []string
	for _, mount := range step.Volumes {
		volume, ok := engine.LookupVolume(spec, mount.Name)
		if !ok {
			continue
		}
		if volume.HostPath != nil {
			if strings.HasPrefix(volume.HostPath.Path, `\\.\pipe\`) {
				continue
			}
		}
		var path string
		if volume.HostPath != nil {
			path = volume.HostPath.Path + ":" + mount.Path
		} else {
			path = volume.Metadata.UID + ":" + mount.Path
		}
		to = append(to, path)
	}
	return to
}

func toVolumeMounts(spec *engine.Spec, step *engine.Step) []mount.Mount {
	var mounts []mount.Mount
	for _, target := range step.Volumes {
		source, ok := engine.LookupVolume(spec, target.Name)
		if !ok {
			continue
		}
		if source.HostPath == nil {
			continue
		}
		if strings.HasPrefix(source.HostPath.Path, `\\.\pipe\`) == false {
			continue
		}
		mounts = append(mounts, mount.Mount{
			Source: source.HostPath.Path,
			Target: target.Path,
			Type:   "npipe",
		})
	}
	if len(mounts) == 0 {
		return nil
	}
	return mounts
}

// helper function that converts a key value map of environment variables to a
// string slice in key=value format.
func toEnv(env map[string]string) []string {
	var envs []string
	for k, v := range env {
		envs = append(envs, k+"="+v)
	}
	return envs
}

// // helper function that converts a slice of device paths to a slice of
// // container.DeviceMapping.
// func toDevices(from []*engine.DeviceMapping) []container.DeviceMapping {
// 	var to []container.DeviceMapping
// 	for _, device := range from {
// 		to = append(to, container.DeviceMapping{
// 			PathOnHost:        device.Source,
// 			PathInContainer:   device.Target,
// 			CgroupPermissions: "rwm",
// 		})
// 	}
// 	return to
// }
