package docker

import (
	"archive/tar"
	"bytes"
	"io"
	"testing"

	"github.com/drone/drone-runtime/engine"
)

func TestCreateTarfile(t *testing.T) {
	file := &engine.File{
		Data: []byte("hello world"),
	}
	mount := &engine.FileMount{
		Path: "/tmp/greeting.txt",
		Mode: 0644,
	}
	d := createTarfile(file, mount)

	r := bytes.NewReader(d)
	tr := tar.NewReader(r)

	hdr, err := tr.Next()
	if err != nil {
		t.Error(err)
		return
	}
	if got, want := hdr.Mode, mount.Mode; got != want {
		t.Errorf("Unexpected file mode. Want %d got %d", want, got)
	}
	if got, want := hdr.Size, len(file.Data); got != int64(want) {
		t.Errorf("Unexpected file size. Want %d got %d", want, got)
	}
	if got, want := hdr.Name, mount.Path; got != want {
		t.Errorf("Unexpected file name. Want %s got %s", want, got)
	}
	buf := new(bytes.Buffer)
	io.Copy(buf, tr)
	if got, want := buf.String(), string(file.Data); got != want {
		t.Errorf("Unexpected file contents. Want %q got %q", want, got)
	}
}
