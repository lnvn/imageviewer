#!/bin/bash

# compile code
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main

# test with SAM
# replace --profile with your actual AWS profile
sam local invoke GreetingFunction --template sam-template/template.yaml -e sam-template/event.json --profile s3access
