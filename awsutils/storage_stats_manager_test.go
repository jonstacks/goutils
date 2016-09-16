package awsutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorageStatsManager(t *testing.T) {
	sm := NewStorageStatsManager()

  // All Storage Stats should be initialized
  for _, s := range storageClasses {
    assert.NotNil(t, sm.Stats[s])
  }
}

// func ExampleStorageStatsManager() {
// 	sm := NewStorageStatsManager()
//
// 	// Output:
// 	// (count: 3, size: 6046)
// }
