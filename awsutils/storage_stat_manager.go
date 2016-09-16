package awsutils

import "github.com/aws/aws-sdk-go/service/s3"

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
	Stats map[string]*StorageStat
}

// NewStorageStatsManager returns a new & initialized StorageStatsManger
func NewStorageStatsManager() *StorageStatsManager {
	sm := &StorageStatsManager{Stats: make(map[string]*StorageStat)}
	for _, s := range storageClasses {
		sm.Stats[s] = NewStorageStat()
	}
	return sm
}

// AddObject adds an object to the storage stat it belongs to
func (sm *StorageStatsManager) AddObject(obj *s3.Object) {
	sm.Stats[*obj.StorageClass].AddObject(*obj.Size)
}
