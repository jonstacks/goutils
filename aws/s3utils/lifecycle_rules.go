package s3utils

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// LifecycleDeletePreviousVersions returns a pointer to a s3.LifecycleRule that
// cleans up the bucket after the specified number of days. It also includes a
// policy to clean up multipart uploads after 7 days
func LifecycleDeletePreviousVersions(days int) *s3.LifecycleRule {
	id := fmt.Sprintf("Delete old versions after %d days", days)
	return &s3.LifecycleRule{
		ID:     aws.String(id),
		Prefix: aws.String(""),
		Status: aws.String("Enabled"),
		NoncurrentVersionExpiration: &s3.NoncurrentVersionExpiration{
			NoncurrentDays: aws.Int64(int64(days)),
		},
		AbortIncompleteMultipartUpload: &s3.AbortIncompleteMultipartUpload{
			DaysAfterInitiation: aws.Int64(7),
		},
	}
}
