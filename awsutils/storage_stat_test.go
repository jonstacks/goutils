package awsutils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testName = "glacier"
var zero = int64(0)

func TestNewStorageStat(t *testing.T) {
	s := NewStorageStat(testName)

	assert.Equal(t, testName, s.Name, "Storage Stat should initialize name")
	assert.Equal(t, zero, s.Count, "Storage Stat's initial Count should be 0")
	assert.Equal(t, zero, s.Size, "Storage Stat's initial Size should be 0")
}

func TestAddObjectStorageStat(t *testing.T) {
	s := NewStorageStat(testName)

	for i := 1; i < 5; i++ {
		s.AddObject(int64(i) * int64(i))
	}

	assert.Equal(t, int64(30), s.Size, "Storage Stat should add objects size correctly")
}

func ExampleStorageStat() {
	s := NewStorageStat(testName)
	s.AddObject(345)
	s.AddObject(23)
	s.AddObject(5678)
	fmt.Println(s)
	// Output:
	// glacier(count: 3, size: 6046)
}
