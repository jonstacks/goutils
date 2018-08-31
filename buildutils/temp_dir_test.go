package buildutils

import (
	"os"
	"testing"
)

func TestTempDir(t *testing.T) {
	td, err := NewTempDir("", "buildutils-")
	if err != nil {
		t.Errorf("Error creating temp directory: %s", err)
	}

	// Check that temp directory exists
	_, err = os.Stat(td.Dir())
	if err != nil {
		if os.IsNotExist(err) {
			t.Errorf("Directory should exist: %s", err)
		}
		t.Errorf("Unknown error: %s", err)
	}

	// With the temp directory
	err = td.Do(func(dir string) error {
		return nil
	})
	if err != nil {
		t.Errorf("TempDir.Do returned an error: %s", err)
	}

	// Check that temp directory no longer exists
	_, err = os.Stat(td.Dir())
	if !os.IsNotExist(err) {
		t.Errorf("Directory should not exist, but it does: %s", td.Dir())
	}
}
