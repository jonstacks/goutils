package s3utils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Get client returns a generic S3 Client
func getClient() *s3.S3 {
	return s3.New(getSession())
}

// getClientForRegion returns a s3 client for the given region
func getClientForRegion(region string) *s3.S3 {
	return s3.New(getSession(), aws.NewConfig().WithRegion(region))
}
