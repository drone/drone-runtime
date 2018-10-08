package docker

import (
	"reflect"
	"testing"

	"github.com/drone/drone-runtime/engine"

	"docker.io/go-docker/api/types/mount"
	"github.com/google/go-cmp/cmp"
)

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
