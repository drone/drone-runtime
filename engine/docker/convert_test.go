package docker

import (
	"reflect"
	"testing"

	"github.com/drone/drone-runtime/engine"
)

func TestDockerConvertConfig(t *testing.T) {
	t.SkipNow()
}

func TestDockerConvertHostConfig(t *testing.T) {
	t.SkipNow()
}

func TestDockerConvertNetwork(t *testing.T) {
	t.SkipNow()
}

func TestDockerConvertVolume(t *testing.T) {
	t.SkipNow()
}

func TestDockerConvertDevice(t *testing.T) {
	from := []engine.DeviceMapping{
		{Source: "/dev/ttyUSB0", Target: "/dev/ttyUSB1"},
	}
	to := toDev(from)
	if len(to) != 1 {
		t.Errorf("Expect device converted to docker.DeviceMapping")
		return
	}
	if got, want := to[0].CgroupPermissions, "rwm"; got != want {
		t.Errorf("Got device cgroup permission %s, want %s", got, want)
	}
	if got, want := to[0].PathOnHost, from[0].Source; got != want {
		t.Errorf("Got device host path %s, want %s", got, want)
	}
	if got, want := to[0].PathInContainer, from[0].Target; got != want {
		t.Errorf("Got device container path %s, want %s", got, want)
	}
}

func TestDockerConvertEnviron(t *testing.T) {
	kv := map[string]string{
		"foo": "bar",
	}
	want := []string{"foo=bar"}
	got := toEnv(kv)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Want environment variables %v, got %v", want, got)
	}
}

func TestDockerEncodeAuthToBase64(t *testing.T) {
	auth := engine.Auth{
		Username: "spaceghost",
		Password: "dianarossfan",
	}
	want := "eyJ1c2VybmFtZSI6InNwYWNlZ2hvc3QiLCJwYXNzd29yZCI6ImRpYW5hcm9zc2ZhbiJ9"
	got := encodeAuthToBase64(auth)
	if got != want {
		t.Errorf("Got base64 encoded string %q, want %q", got, want)
	}
}
