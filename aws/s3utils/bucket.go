package s3utils

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// DefaultRegion is the default region for a bucket. Also known as US Standard.
const DefaultRegion = "us-east-1"

// Bucket is used to store a reference to the s3 client and information about a
// bucket. It also makes it easier to perform common opperations.
type Bucket struct {
	client       *s3.S3
	Name         string
	Region       string
	StatsManager *StorageStatsManager
}

// FindBucket returns a New Bucket object by looking for an Existing Bucket. If
// the bucket does not exist, it returns nil
func FindBucket(name string) *Bucket {
	if !bucketExists(name) {
		return nil
	}
	region, _ := getBucketRegion(name)
	client := getClientForRegion(region)
	return &Bucket{client, name, region, NewStorageStatsManager()}
}

// CreateBucket creates a bucket in the specified region
func CreateBucket(name, region string) (*Bucket, error) {
	client := getClientForRegion(region)
	input := &s3.CreateBucketInput{
		Bucket: aws.String(name),
	}

	_, err := client.CreateBucket(input)
	if err != nil {
		return nil, err
	}
	return &Bucket{client, name, region, NewStorageStatsManager()}, nil
}

// IsVersioningEnabled determines whether or not versioning is enabled on the
// bucket.
func (b *Bucket) IsVersioningEnabled() (bool, error) {
	input := &s3.GetBucketVersioningInput{
		Bucket: aws.String(b.Name),
	}
	resp, err := b.client.GetBucketVersioning(input)
	if err != nil {
		return false, err
	}
	status := derefStringOrEmpty(resp.Status)
	return status == s3.BucketVersioningStatusEnabled, nil
}

// EnableVersioning enables versioning for the bucket. It returns an error or
// nil if there was an issue
func (b *Bucket) EnableVersioning() error {
	input := &s3.PutBucketVersioningInput{
		Bucket: aws.String(b.Name),
		VersioningConfiguration: &s3.VersioningConfiguration{
			Status: aws.String(s3.BucketVersioningStatusEnabled),
		},
	}

	_, err := b.client.PutBucketVersioning(input)
	return err
}

func (b *Bucket) String() string {
	return fmt.Sprintf("%s (region: %s)", b.Name, b.Region)
}

// LifecycleRules returns a slice of pointers to s3.LifecycleRules.
// Unfortunately this throws a nasty error if ther are no lifecycle rules
func (b *Bucket) LifecycleRules() ([]*s3.LifecycleRule, error) {
	input := &s3.GetBucketLifecycleConfigurationInput{Bucket: aws.String(b.Name)}
	resp, err := b.client.GetBucketLifecycleConfiguration(input)
	if err != nil {
		return []*s3.LifecycleRule{}, err
	}
	return resp.Rules, nil
}

// AddLifecycleRule adds the given rule to the existing lifecycle configuration
func (b *Bucket) AddLifecycleRule(rule *s3.LifecycleRule) error {
	rules, err := b.LifecycleRules()
	found := false

	if err != nil {
		// There is probably no lifecycle configuration?? It would be nice if the
		// API let us handle that error better. So lets just assume no rules
	} else {
		// Loop over the rules looking to update a desired rule
		ruleName := ""
		for i := range rules {
			// If this rule is already in the Lifecycle Rules, update it
			ruleName = *rules[i].ID
			if ruleName == *rule.ID {
				found = true
				rules[i] = rule
			}
		}
	}

	if !found {
		// This rule does not yet exist in our rules, so lets go ahead and add it
		// to the rules
		rules = append(rules, rule)
	}

	// Now lets update the rules
	putInput := &s3.PutBucketLifecycleConfigurationInput{
		Bucket: aws.String(b.Name),
		LifecycleConfiguration: &s3.BucketLifecycleConfiguration{
			Rules: rules,
		},
	}
	validationErr := putInput.Validate()
	if validationErr != nil {
		return validationErr
	}

	_, putErr := b.client.PutBucketLifecycleConfiguration(putInput)
	if putErr != nil {
		return putErr
	}

	return nil
}

//
// Package Internal Methods
//

// BucketExists returns whether or not the bucket with the given name exists
func bucketExists(name string) bool {
	_, err := getBucketRegion(name)
	return err == nil
}

// Gets the buckets region
func getBucketRegion(name string) (string, error) {
	region := DefaultRegion
	client := getClientForRegion(region)
	input := &s3.GetBucketLocationInput{Bucket: aws.String(name)}
	resp, err := client.GetBucketLocation(input)
	if err != nil {
		return region, err
	}

	// If this is nil, fall back to using the DefaultRegion that is already set.
	if resp.LocationConstraint != nil {
		region = *resp.LocationConstraint
	}

	return region, nil
}
