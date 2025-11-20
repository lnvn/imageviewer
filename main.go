package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Global variable for the S3 client, initialized during the cold start
var s3Client *s3.Client

func init() {
	// The init function now loads the AWS configuration and creates the S3 client.
	// Load the AWS configuration from the environment (credentials, region)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Printf("FATAL: unable to load SDK config, %v\n", err)
		os.Exit(1)
	}

	// Create the S3 service client using the loaded configuration
	s3Client = s3.NewFromConfig(cfg)
}

// HandleRequest lists the first 5 S3 buckets and returns a greeting plus the list.
func HandleRequest(ctx context.Context) (string, error) {
	name := "AWS User" 
	
	result, err := s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return "", fmt.Errorf("failed to list buckets: %w", err)
	}

	var bucketNames []string
	for i, bucket := range result.Buckets {
		if i < 5 { // Limit to the first 5 for brevity
			bucketNames = append(bucketNames, *bucket.Name)
		}
	}
	
	var bucketList string
	if len(bucketNames) > 0 {
		bucketList = fmt.Sprintf("Found %d buckets (first 5: %s).", 
			len(result.Buckets), strings.Join(bucketNames, ", "))
	} else {
		bucketList = "No S3 buckets found."
	}

	// 4. Return the combined message
	greeting := fmt.Sprintf("Hello, %s! S3 Status: %s", name, bucketList)
	return greeting, nil
}

func main() {
	lambda.Start(HandleRequest)
}