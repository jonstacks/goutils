package buildutils

import (
	"io/ioutil"
	"os"
)

// TempDir is a temporary Directory
type TempDir struct {
	dir string
}

// NewTempDir creates a new temporary directory.
func NewTempDir(dir, prefix string) (*TempDir, error) {
	tmpDir, err := ioutil.TempDir(dir, prefix)
	if err != nil {
		return nil, err
	}
	td := &TempDir{dir: tmpDir}
	return td, nil
}

// Dir The temp directory created
func (td *TempDir) Dir() string {
	return td.dir
}

// Do is a wrapper that automatically cleans up a temp directory when the
// function that is invoked returns
func (td *TempDir) Do(f func(dir string) error) error {
	// Make it hard to forgot to clean up after yourself.
	defer func(d string) { os.RemoveAll(d) }(td.Dir())

	return f(td.Dir())
}
