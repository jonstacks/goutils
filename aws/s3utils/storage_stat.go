package s3utils

import (
	"fmt"
)

// StorageStat provides a data structure for storing storage information. It
// can be used to store stats for separate AWS storage classes.
type StorageStat struct {
	Count int64
	Size  int64
}

// NewStorageStat initializes a new StorageStat and returns a pointer to it
func NewStorageStat() *StorageStat {
	return &StorageStat{0, 0}
}

// AddObject adds an object to the Storage Stat by its size.
func (s *StorageStat) AddObject(size int64) {
	s.Count++
	s.Size += size
}

// AverageFileSize returns the average size, in bytes, of the files the
// StorageStat counts
func (s *StorageStat) AverageFileSize() float64 {
	return float64(s.Size) / float64(s.Count)
}

func (s *StorageStat) String() string {
	return fmt.Sprintf("(count: %d, size: %d)", s.Count, s.Size)
}
