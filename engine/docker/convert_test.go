package docker

import (
	"reflect"
	"testing"
)

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
