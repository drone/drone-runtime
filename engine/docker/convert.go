package docker

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/drone/drone-runtime/engine"
)

// returns a container configuration.
func toConfig(proc *engine.Step) *container.Config {
	config := &container.Config{
		Image:        proc.Image,
		Labels:       proc.Labels,
		WorkingDir:   proc.WorkingDir,
		AttachStdout: true,
		AttachStderr: true,
	}
	if len(proc.Environment) != 0 {
		config.Env = toEnv(proc.Environment)
	}
	if len(proc.Secrets) != 0 {
		for _, secret := range proc.Secrets {
			config.Env = append(config.Env, secret.Name+"="+secret.Value)
		}
	}
	if len(proc.Command) != 0 {
		config.Cmd = proc.Command
	}
	if len(proc.Entrypoint) != 0 {
		config.Entrypoint = proc.Entrypoint
	}
	if len(proc.Volumes) != 0 {
		config.Volumes = toVolumeSet(proc.Volumes)
	}
	return config
}

// returns a container host configuration.
func toHostConfig(proc *engine.Step) *container.HostConfig {
	config := &container.HostConfig{
		Resources: container.Resources{
			CPUQuota:   proc.CPUQuota,
			CPUShares:  proc.CPUShares,
			CpusetCpus: proc.CPUSet,
			Memory:     proc.MemLimit,
			MemorySwap: proc.MemSwapLimit,
		},
		LogConfig: container.LogConfig{
			Type: "json-file",
		},
		Privileged: proc.Privileged,
		ShmSize:    proc.ShmSize,
		Sysctls:    proc.Sysctls,
	}

	// if len(proc.VolumesFrom) != 0 {
	// 	config.VolumesFrom = proc.VolumesFrom
	// }
	if len(proc.NetworkMode) != 0 {
		config.NetworkMode = container.NetworkMode(proc.NetworkMode)
	}
	if len(proc.IpcMode) != 0 {
		config.IpcMode = container.IpcMode(proc.IpcMode)
	}
	if len(proc.DNS) != 0 {
		config.DNS = proc.DNS
	}
	if len(proc.DNSSearch) != 0 {
		config.DNSSearch = proc.DNSSearch
	}
	if len(proc.ExtraHosts) != 0 {
		config.ExtraHosts = proc.ExtraHosts
	}
	if len(proc.Devices) != 0 {
		config.Devices = toDevices(proc.Devices)
	}
	if len(proc.Volumes) != 0 {
		config.Binds = toVolumeSlice(proc.Volumes)
	}
	config.Tmpfs = map[string]string{}
	for _, path := range proc.Tmpfs {
		if strings.Index(path, ":") == -1 {
			config.Tmpfs[path] = ""
			continue
		}
		parts := strings.Split(path, ":")
		config.Tmpfs[parts[0]] = parts[1]
	}
	// if proc.OomKillDisable {
	// 	config.OomKillDisable = &proc.OomKillDisable
	// }

	return config
}

// helper function that converts a slice of volume paths to a set of
// unique volume names.
func toVolumeSet(from []*engine.VolumeMapping) map[string]struct{} {
	to := map[string]struct{}{}
	for _, v := range from {
		to[v.Target] = struct{}{}
	}
	return to
}

func toVolumeSlice(from []*engine.VolumeMapping) []string {
	var to []string
	for _, v := range from {
		var path string
		if v.Name != "" {
			path = v.Name + ":" + v.Target
		} else {
			path = v.Source + ":" + v.Target
		}
		to = append(to, path)
	}
	return to
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

// helper function that converts a slice of device paths to a slice of
// container.DeviceMapping.
func toDevices(from []*engine.DeviceMapping) []container.DeviceMapping {
	var to []container.DeviceMapping
	for _, device := range from {
		to = append(to, container.DeviceMapping{
			PathOnHost:        device.Source,
			PathInContainer:   device.Target,
			CgroupPermissions: "rwm",
		})
	}
	return to
}

// helper function that serializes the auth configuration as JSON
// base64 payload.
func encodeAuthToBase64(authConfig engine.Auth) string {
	buf, _ := json.Marshal(authConfig)
	return base64.URLEncoding.EncodeToString(buf)
}
