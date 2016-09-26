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
	return putErr
}

// IsReplicated returns whether or not the replication is enabled for the
// bucket. Unfortunately, since this relies on the ReplicationConfig method
// and the underlying service client returns an error if the config doesn't
// exist, we have to rely on an error meaning that replication is not
// configured.
func (b *Bucket) IsReplicated() bool {
	if _, err := b.ReplicationConfig(); err != nil {
		return false
	}
	return true
}

// Permissions returns a list of permissions on the bucket.
func (b *Bucket) Permissions() (*s3.GetBucketAclOutput, error) {
	input := &s3.GetBucketAclInput{Bucket: aws.String(b.Name)}
	resp, err := b.client.GetBucketAcl(input)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// CopyPermissionsFrom copies all ACLs from the source bucket to this bucket,
// overwriting this buckets ACLs if present.
func (b *Bucket) CopyPermissionsFrom(src *Bucket) error {
	// Get the current Permissions from the current
	input := &s3.GetBucketAclInput{Bucket: aws.String(src.Name)}
	resp, err := src.client.GetBucketAcl(input)
	if err != nil {
		return err
	}

	// Set the default type for the grants grantee. Not sure why this doesn't
	// come back as part of the original request
	for i := range resp.Grants {
		switch {
		case resp.Grants[i].Grantee.ID != nil:
			// assume we are dealing with a cannonical User
			resp.Grants[i].Grantee.Type = aws.String(s3.TypeCanonicalUser)
		case resp.Grants[i].Grantee.URI != nil:
			resp.Grants[i].Grantee.Type = aws.String(s3.TypeGroup)
		}

		if err = resp.Grants[i].Grantee.Validate(); err != nil {
			return err
		}
	}

	// Now copy those permissions to our bucket
	putInput := &s3.PutBucketAclInput{
		Bucket: aws.String(b.Name),
		AccessControlPolicy: &s3.AccessControlPolicy{
			Grants: resp.Grants,
			Owner:  resp.Owner,
		},
	}
	_, err = b.client.PutBucketAcl(putInput)
	if err != nil {
		return err
	}
	return nil
}

// ProfileStorage iterates over the objects in the bucket and adds their stats
// to the StorageStatsManager.
func (b *Bucket) ProfileStorage() error {
	params := &s3.ListObjectsV2Input{
		Bucket:  aws.String(b.Name), // Required
		MaxKeys: aws.Int64(1000),
	}
	err := params.Validate()
	if err != nil {
		return err
	}

	b.client.ListObjectsV2Pages(params,
		func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			for _, c := range page.Contents {
				b.StatsManager.AddObject(c)
			}
			return !lastPage
		},
	)
	return nil
}

// // Replicate adds configuration to replicate this bucket to the destination
// // Bucket.
// func (b *Bucket) Replicate(dest *Bucket, storageClass string) error {
// 	// The desired Rule we will be creating:
// 	rule := &s3.ReplicationRule{
// 		ID:     aws.String(b.Name),
// 		Prefix: aws.String("/"),
// 		Status: aws.String(s3.ReplicationRuleStatusEnabled),
// 		Destination: &s3.Destination{
// 			Bucket:       aws.String(dest.Name),
// 			StorageClass: aws.String(storageClass),
// 		},
// 	}
//
// 	// First see if this bucket is replicated.
// 	config, err := b.ReplicationConfig()
// 	if err != nil {
// 		// No replication config exists, we need to create a config
// 		config = &s3.ReplicationConfiguration{
// 			Role: aws.String(createReplicationRole(b.Name, dest.Name)),
// 			Rules: []*s3.ReplicationRule{
// 				rule,
// 			},
// 		}
// 	} else {
// 		// There is an existing configuration. Great lets use that and make sure our
// 		// rule is in there.
// 		found := false // Whether or not we found the rule in the existing rules
// 		ruleName := ""
// 		for i := range config.Rules {
// 			ruleName = *config.Rules[i].ID
// 			if ruleName == *rule.ID {
// 				found = true
// 				config.Rules[i] = rule
// 			}
// 		}
//
// 		if !found {
// 			// This rule does not yet exist in our rules, so lets go ahead and add it
// 			// to the rules
// 			config.Rules = append(config.Rules, rule)
// 		}
// 	}
//
// 	// We have a config, go ahead and put the configuration
// 	input := &s3.PutBucketReplicationInput{
// 		Bucket: aws.String(b.Name),
// 		ReplicationConfiguration: config,
// 	}
//
// 	_, err = b.client.PutBucketReplication(input)
// 	return err
// }

// ReplicationConfig returns the s3 Replication Configuration. If the bucket
// does not have a ReplicationConfiguration, it will return an error.
func (b *Bucket) ReplicationConfig() (*s3.ReplicationConfiguration, error) {
	input := &s3.GetBucketReplicationInput{Bucket: aws.String(b.Name)}
	resp, err := b.client.GetBucketReplication(input)
	if err != nil {
		return resp.ReplicationConfiguration, err
	}
	return resp.ReplicationConfiguration, nil
}

//
// Package Internal Methods
//

// bucketExists returns whether or not the bucket with the given name exists
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
