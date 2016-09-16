package awsutils

import (
	"fmt"
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

func ExampleStorageStatsManager() {
	sm := NewStorageStatsManager()
	fmt.Println(sm)
	// Output:
	// [STANDARD(count: 0, size: 0), STANDARD_IA(count: 0, size: 0), REDUCED_REDUNDANCY(count: 0, size: 0), GLACIER(count: 0, size: 0)]
}
