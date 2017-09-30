package chroot

import (
	"io"
	"os"
	"path/filepath"

	"github.com/drone/drone-runtime/runtime"
)

type chroot struct {
	base string
}

// New returns a new chroot file system.
func New(base string) (runtime.FileSystem, error) {
	return &chroot{base}, nil
}

func (fs *chroot) Open(path string) (io.ReadCloser, error) {
	path = filepath.Clean(path)
	return os.Open(
		filepath.Join(fs.base, path),
	)
}

func (fs *chroot) Stat(path string) (os.FileInfo, error) {
	path = filepath.Clean(path)
	return os.Stat(
		filepath.Join(fs.base, path),
	)
}

func (fs *chroot) Create(path string) (io.WriteCloser, error) {
	path = filepath.Clean(path)
	path = filepath.Join(fs.base, path)
	base := filepath.Dir(path)
	os.MkdirAll(base, 0700)
	return os.Create(path)
}

func (fs *chroot) Remove(path string) error {
	path = filepath.Clean(path)
	return os.Remove(
		filepath.Join(fs.base, path),
	)
}
