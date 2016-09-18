package s3utils

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/s3"
)

// Storage Classes
const (
	StorageClassStandard          = "STANDARD"
	StorageClassStandardIA        = "STANDARD_IA"
	StorageClassReducedRedundancy = "REDUCED_REDUNDANCY"
	StorageClassGlacier           = "GLACIER"
)

var storageClasses = []string{
	StorageClassStandard,
	StorageClassStandardIA,
	StorageClassReducedRedundancy,
	StorageClassGlacier,
}

// StorageStatsManager is used to manage the stats for different AWS Storage
// Classes in an easy to use interface.
type StorageStatsManager struct {
	Stats   map[string]*StorageStat
	classes *[]string // This is just here so we can keep track of the ordering
}

// NewStorageStatsManager returns a new & initialized StorageStatsManger
func NewStorageStatsManager() *StorageStatsManager {
	sm := &StorageStatsManager{
		Stats:   make(map[string]*StorageStat),
		classes: &storageClasses,
	}
	for _, s := range *sm.classes {
		sm.Stats[s] = NewStorageStat()
	}
	return sm
}

// AddObject adds an object to the storage stat it belongs to
func (sm *StorageStatsManager) AddObject(obj *s3.Object) {
	sm.Stats[*obj.StorageClass].AddObject(*obj.Size)
}

// Size returns the total size across all storage stats
func (sm *StorageStatsManager) Size() int64 {
	size := int64(0)
	for _, v := range sm.Stats {
		size += v.Size
	}
	return size
}

// Count returns the total count across all storage stats
func (sm *StorageStatsManager) Count() int64 {
	count := int64(0)
	for _, v := range sm.Stats {
		count += v.Count
	}
	return count
}

func (sm *StorageStatsManager) String() string {
	stats := make([]string, 0, len(sm.Stats))

	for _, k := range *sm.classes {
		stats = append(stats, fmt.Sprintf("%s%s", k, sm.Stats[k]))
	}
	x := strings.Join(stats, ", ")
	return fmt.Sprintf("[%s]", x)
}
