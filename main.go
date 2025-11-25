package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Global variable for the S3 client
var s3Client *s3.Client

func init() {
	// The init function now loads the AWS configuration and creates the S3 client.
	// Load the AWS configuration from the environment (credentials, region)
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "ap-southeast-1" // Default fallback
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		fmt.Printf("FATAL: unable to load SDK config, %v\n", err)
		os.Exit(1)
	}

	// Create the S3 service client using the loaded configuration
	s3Client = s3.NewFromConfig(cfg)
}

// HandleRequest downloads the file from S3 and returns it as a Base64 encoded string.
func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	bucketName := os.Getenv("S3_BUCKET_NAME")
	filename := os.Getenv("FILE_NAME")

	if bucketName == "" || filename == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "S3_BUCKET_NAME and FILE_NAME environment variables are required",
		}, nil
	}

	// Get the object from S3
	result, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &filename,
	})
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Failed to get file from S3: %v", err),
		}, nil
	}
	defer result.Body.Close()

	// Read the content
	bodyBytes, err := io.ReadAll(result.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Failed to read file content: %v", err),
		}, nil
	}

	// Encode to Base64
	encodedBody := base64.StdEncoding.EncodeToString(bodyBytes)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "image/jpeg", // Assuming JPEG for now, could be dynamic
		},
		Body:            encodedBody,
		IsBase64Encoded: true,
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
