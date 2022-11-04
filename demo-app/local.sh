#!/bin/bash
export AWS_ACCESS_KEY_ID="test"
export AWS_SECRET_ACCESS_KEY="test"
export AWS_DEFAULT_REGION="us-east-1"

aws --endpoint-url=http://localstack:4566 s3 ls
aws --endpoint-url=http://localstack:4566 sqs list-queues
aws --endpoint-url=http://localstack:4566 kinesis list-streams
aws --endpoint-url=http://localstack:4566 lambda list-functions

sleep 360000