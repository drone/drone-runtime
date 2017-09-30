package runtime

//go:generate mockgen -source=fs.go -destination=mocks/fs.go -package=mocks

import (
	"io"
	"os"
)

// FileSystem interface defines an abstract file system used by the runtime
// to snapshot and restore a container file system.
type FileSystem interface {
	Open(string) (io.ReadCloser, error)
	Stat(string) (os.FileInfo, error)
	Create(string) (io.WriteCloser, error)
	Remove(string) error
}
