package main

import (
	"context"
	"fmt"
	"os"
	"github.com/aws/aws-lambda-go/lambda"
)

var defaultName string

func init() {
	defaultName = os.Getenv("GREETING_NAME")
	if defaultName == "" {
		defaultName = "Stranger"
	}
}

func HandleRequest(ctx context.Context) (string, error) {
	return fmt.Sprintf("Hello, %s!", defaultName), nil
}

func main() {
	lambda.Start(HandleRequest)
}