package awsutils

import (
	"fmt"
)

// StorageStat provides a data structure for storing storage information. It
// can be used to store stats for separate AWS storage classes.
type StorageStat struct {
	Name  string
	Count int64
	Size  int64
}

// NewStorageStat initializes a new StorageStat and returns a pointer to it
func NewStorageStat(name string) *StorageStat {
	return &StorageStat{name, 0, 0}
}

// AddObject adds an object to the Storage Stat by its size.
func (s *StorageStat) AddObject(size int64) {
	s.Count++
	s.Size += size
}

func (s *StorageStat) String() string {
	return fmt.Sprintf("%s(count: %d, size: %d)", s.Name, s.Count, s.Size)
}
