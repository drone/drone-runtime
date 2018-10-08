package docker

import (
	"reflect"
	"testing"

	"github.com/drone/drone-runtime/engine"

	"docker.io/go-docker/api/types/container"
	"docker.io/go-docker/api/types/mount"
	"docker.io/go-docker/api/types/network"
	"github.com/google/go-cmp/cmp"
)

func TestToConfig(t *testing.T) {
	step := &engine.Step{
		Metadata: engine.Metadata{
			UID:    "123",
			Name:   "test",
			Labels: map[string]string{},
		},
		Envs: map[string]string{
			"GOOS": "linux",
		},
		Docker: &engine.DockerStep{
			Image:   "golang:latest",
			Command: []string{"/bin/sh"},
			Args:    []string{"-c", "go build; go test -v"},
		},
		WorkingDir: "/workspace",
		Secrets: []*engine.SecretVar{
			{
				Name: "password",
				Env:  "HTTP_PASSWORD",
			},
		},
	}
	spec := &engine.Spec{
		Metadata: engine.Metadata{
			UID: "abc123",
		},
		Steps: []*engine.Step{step},
		Secrets: []*engine.Secret{
			{
				Name: "password",
				Data: "correct-horse-battery-staple",
			},
		},
	}
	a := &container.Config{
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
		Entrypoint:   step.Docker.Command,
		Cmd:          step.Docker.Args,
		Env: []string{
			"GOOS=linux",
			"HTTP_PASSWORD=correct-horse-battery-staple",
		},
	}
	b := toConfig(spec, step)
	if diff := cmp.Diff(a, b); diff != "" {
		t.Errorf("Unexpected container.Config")
		t.Log(diff)
	}
}

func TestToHostConfig(t *testing.T) {
	step := &engine.Step{
		Metadata: engine.Metadata{
			UID:    "123",
			Name:   "test",
			Labels: map[string]string{},
		},
		Docker: &engine.DockerStep{
			Image:      "golang:latest",
			Command:    []string{"/bin/sh"},
			Args:       []string{"-c", "go build; go test -v"},
			Privileged: true,
			ExtraHosts: []string{"host.company.com"},
			DNS:        []string{"8.8.8.8"},
			DNSSearch:  []string{"dns.company.com"},
		},
		Resources: &engine.Resources{
			Limits: &engine.ResourceObject{
				Memory: 10000,
			},
		},
		Volumes: []*engine.VolumeMount{
			{Name: "foo", Path: "/foo"},
			{Name: "bar", Path: "/baz"},
		},
	}
	spec := &engine.Spec{
		Metadata: engine.Metadata{
			UID: "abc123",
		},
		Steps: []*engine.Step{step},
		Docker: &engine.DockerConfig{
			Volumes: []*engine.Volume{
				{
					Metadata: engine.Metadata{Name: "foo", UID: "1"},
					EmptyDir: &engine.VolumeEmptyDir{},
				},
				{
					Metadata: engine.Metadata{Name: "bar", UID: "2"},
					HostPath: &engine.VolumeHostPath{Path: "/bar"},
				},
			},
		},
	}
	a := &container.HostConfig{
		Privileged: true,
		LogConfig: container.LogConfig{
			Type: "json-file",
		},
		Binds:      []string{"1:/foo"},
		DNS:        []string{"8.8.8.8"},
		DNSSearch:  []string{"dns.company.com"},
		ExtraHosts: []string{"host.company.com"},
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: "/bar",
				Target: "/baz",
			},
		},
		Resources: container.Resources{
			Memory: 10000,
		},
	}
	b := toHostConfig(spec, step)
	if diff := cmp.Diff(a, b); diff != "" {
		t.Errorf("Unexpected container.HostConfig")
		t.Log(diff)
	}

	// we ensure that privileged mode is always mapped
	// correctly. better to be safe ...

	step.Docker.Privileged = false
	b = toHostConfig(spec, step)
	if b.Privileged {
		t.Errorf("Expect privileged mode disabled")
	}
}

func TestToNetConfig(t *testing.T) {
	step := &engine.Step{
		Metadata: engine.Metadata{
			Name: "redis",
		},
	}
	spec := &engine.Spec{
		Metadata: engine.Metadata{
			UID: "abc123",
		},
		Steps: []*engine.Step{step},
	}
	a := toNetConfig(spec, step)
	b := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			"abc123": &network.EndpointSettings{
				Aliases:   []string{"redis"},
				NetworkID: "abc123"},
		},
	}
	if diff := cmp.Diff(a, b); diff != "" {
		t.Errorf("Unexpected network configuration")
		t.Log(diff)
	}
}

func TestToVolumeSlice(t *testing.T) {
	step := &engine.Step{
		Volumes: []*engine.VolumeMount{
			{Name: "foo", Path: "/foo"},
			{Name: "bar", Path: "/bar"},
			{Name: "baz", Path: "/baz"},
		},
	}
	spec := &engine.Spec{
		Steps: []*engine.Step{step},
		Docker: &engine.DockerConfig{
			Volumes: []*engine.Volume{
				{
					Metadata: engine.Metadata{Name: "foo", UID: "1"},
					EmptyDir: &engine.VolumeEmptyDir{},
				},
				{
					Metadata: engine.Metadata{Name: "bar", UID: "2"},
					HostPath: &engine.VolumeHostPath{Path: "/bar"},
				},
			},
		},
	}

	a := toVolumeSlice(spec, step)
	b := []string{"1:/foo"}
	if diff := cmp.Diff(a, b); diff != "" {
		t.Errorf("Unexpected volume slice")
		t.Log(diff)
	}
}

func TestToVolumeMounts(t *testing.T) {
	step := &engine.Step{
		Volumes: []*engine.VolumeMount{
			{Name: "foo", Path: "/foo"},
			{Name: "bar", Path: "/bar"},
			{Name: "baz", Path: "/baz"},
		},
	}
	spec := &engine.Spec{
		Steps: []*engine.Step{step},
		Docker: &engine.DockerConfig{
			Volumes: []*engine.Volume{
				{
					Metadata: engine.Metadata{Name: "foo", UID: "1"},
					EmptyDir: &engine.VolumeEmptyDir{},
				},
				{
					Metadata: engine.Metadata{Name: "bar", UID: "2"},
					HostPath: &engine.VolumeHostPath{Path: "/tmp"},
				},
			},
		},
	}

	a := toVolumeMounts(spec, step)
	b := []mount.Mount{
		{Type: mount.TypeBind, Source: "/tmp", Target: "/bar"},
	}
	if diff := cmp.Diff(a, b); diff != "" {
		t.Errorf("Unexpected volume mounts")
		t.Log(diff)
	}

	step.Volumes = []*engine.VolumeMount{}
	if toVolumeMounts(spec, step) != nil {
		t.Errorf("Expect nil volume mount")
	}
}

func TestToEnv(t *testing.T) {
	kv := map[string]string{
		"foo": "bar",
	}
	want := []string{"foo=bar"}
	got := toEnv(kv)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Want environment variables %v, got %v", want, got)
	}
}

func TestToMount(t *testing.T) {
	tests := []struct {
		source *engine.Volume
		target *engine.VolumeMount
		result mount.Mount
	}{
		// volume mount
		{
			source: &engine.Volume{EmptyDir: &engine.VolumeEmptyDir{}},
			target: &engine.VolumeMount{Path: "/foo"},
			result: mount.Mount{Type: mount.TypeVolume, Target: "/foo"},
		},
		// tmpfs mount
		{
			source: &engine.Volume{EmptyDir: &engine.VolumeEmptyDir{Medium: "memory", SizeLimit: 10000}},
			target: &engine.VolumeMount{Path: "/foo"},
			result: mount.Mount{Type: mount.TypeTmpfs, Target: "/foo", TmpfsOptions: &mount.TmpfsOptions{SizeBytes: 10000, Mode: 0700}},
		},
		// bind mount
		{
			source: &engine.Volume{HostPath: &engine.VolumeHostPath{Path: "/foo"}},
			target: &engine.VolumeMount{Path: "/bar"},
			result: mount.Mount{Type: mount.TypeBind, Source: "/foo", Target: "/bar"},
		},
		// named pipe
		{
			source: &engine.Volume{HostPath: &engine.VolumeHostPath{Path: `\\.\pipe\docker_engine`}},
			target: &engine.VolumeMount{Path: `\\.\pipe\docker_engine`},
			result: mount.Mount{Type: mount.TypeNamedPipe, Source: `\\.\pipe\docker_engine`, Target: `\\.\pipe\docker_engine`},
		},
	}
	for _, test := range tests {
		result := toMount(test.source, test.target)
		if diff := cmp.Diff(result, test.result); diff != "" {
			t.Error("Unexpected mount value")
			t.Log(diff)
		}
	}
}

func TestToVolumeType(t *testing.T) {
	tests := []struct {
		volume *engine.Volume
		value  mount.Type
	}{
		{
			volume: &engine.Volume{},
			value:  mount.TypeBind,
		},
		{
			volume: &engine.Volume{HostPath: &engine.VolumeHostPath{Path: "/tmp"}},
			value:  mount.TypeBind,
		},
		{
			volume: &engine.Volume{HostPath: &engine.VolumeHostPath{Path: `\\.\pipe\docker_engine`}},
			value:  mount.TypeNamedPipe,
		},
		{
			volume: &engine.Volume{EmptyDir: &engine.VolumeEmptyDir{Medium: "memory"}},
			value:  mount.TypeTmpfs,
		},
		{
			volume: &engine.Volume{EmptyDir: &engine.VolumeEmptyDir{}},
			value:  mount.TypeVolume,
		},
	}
	for _, test := range tests {
		if got, want := toVolumeType(test.volume), test.value; got != want {
			t.Errorf("Want mount type %v, got %v", want, got)
		}
	}
}

func TestIsBindMount(t *testing.T) {
	tests := []struct {
		volume *engine.Volume
		value  bool
	}{
		{
			volume: &engine.Volume{},
			value:  false,
		},
		{
			volume: &engine.Volume{HostPath: &engine.VolumeHostPath{Path: "/tmp"}},
			value:  true,
		},
	}
	for _, test := range tests {
		if got, want := isBindMount(test.volume), test.value; got != want {
			t.Errorf("Want is bind mount %v, got %v", want, got)
		}
	}
}

func TestIsNamedPipe(t *testing.T) {
	tests := []struct {
		volume *engine.Volume
		value  bool
	}{
		{
			volume: &engine.Volume{},
			value:  false,
		},
		{
			volume: &engine.Volume{HostPath: &engine.VolumeHostPath{Path: "/tmp"}},
			value:  false,
		},
		{
			volume: &engine.Volume{HostPath: &engine.VolumeHostPath{Path: `\\.\pipe\docker_engine`}},
			value:  true,
		},
	}
	for _, test := range tests {
		if got, want := isNamedPipe(test.volume), test.value; got != want {
			t.Errorf("Want is named pipe %v, got %v", want, got)
		}
	}
}

func TestIsTempfs(t *testing.T) {
	tests := []struct {
		volume *engine.Volume
		value  bool
	}{
		{
			volume: &engine.Volume{},
			value:  false,
		},
		{
			volume: &engine.Volume{EmptyDir: &engine.VolumeEmptyDir{}},
			value:  false,
		},
		{
			volume: &engine.Volume{EmptyDir: &engine.VolumeEmptyDir{Medium: "memory"}},
			value:  true,
		},
	}
	for _, test := range tests {
		if got, want := isTempfs(test.volume), test.value; got != want {
			t.Errorf("Want is temp fs %v, got %v", want, got)
		}
	}
}

func TestIsDataVolume(t *testing.T) {
	tests := []struct {
		volume *engine.Volume
		value  bool
	}{
		{
			volume: &engine.Volume{},
			value:  false,
		},
		{
			volume: &engine.Volume{EmptyDir: &engine.VolumeEmptyDir{Medium: "memory"}},
			value:  false,
		},
		{
			volume: &engine.Volume{EmptyDir: &engine.VolumeEmptyDir{}},
			value:  true,
		},
	}
	for i, test := range tests {
		if got, want := isDataVolume(test.volume), test.value; got != want {
			t.Errorf("Want is data volume %v, got %v at index %d", want, got, i)
		}
	}
}
