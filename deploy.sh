#!/bin/sh

# Ensure the script returns an exit code on failure.
set -e

# Check the version.
sam --version

# Check the environment variables.
if [ -z "$MOBILEPHONE" ]; then
  echo "Error: you must set the environment variable: MOBILEPHONE (format: +12225557777)"
  exit
fi

if [ -z "$ACCESS_KEY" ]; then
  echo "Error: you must set the environment variable: ACCESS_KEY"
  exit
fi

if [ -z "$SECRET_KEY" ]; then
  echo "Error: you must set the environment variable: SECRET_KEY"
  exit
fi

# Set the application name.
APPNAME=auth

# Build the application and include the API in the binary.
GOOS=linux go build -ldflags "-X main.AWSAccessKeyID=${ACCESS_KEY} -X main.AWSSecretAccessKey=${SECRET_KEY}" -o ${APPNAME} main.go
zip handler.zip ./${APPNAME}

# Make an S3 bucket.
aws s3 mb s3://$(aws sts get-caller-identity \
--query 'Account' --output text)-${APPNAME}

# Copy the code to S3.
aws cloudformation package \
    --template-file template.json \
    --s3-bucket $(aws sts get-caller-identity \
    --query 'Account' --output text)-${APPNAME} \
    --output-template-file packaged-template.yml

# Deploy the API gateway and lambda.
aws cloudformation deploy \
    --template-file packaged-template.yml \
    --stack-name ${APPNAME} \
    --capabilities CAPABILITY_IAM \
    --parameter-overrides MobilePhone=${MOBILEPHONE}

# Output the outputs.
aws cloudformation describe-stacks --stack-name ${APPNAME} --query 'Stacks[0].Outputs'

# Clean up.
rm ./${APPNAME}
rm ./handler.zip
rm ./packaged-template.yml