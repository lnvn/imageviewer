#!/bin/bash

# compile code
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main

# test with SAM
# replace --profile with your actual AWS profile
sam local invoke GreetingFunction --template aws-sam/template.yaml -e aws-sam/event.json --profile s3access
