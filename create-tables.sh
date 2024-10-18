#!/usr/bin/env bash

set -e

aws dynamodb create-table \
  --table-name TodoItems \
  --attribute-definitions \
      AttributeName=ID,AttributeType=S \
      AttributeName=UserID,AttributeType=S \
  --key-schema \
    AttributeName=UserID,KeyType=HASH \
    AttributeName=ID,KeyType=RANGE \
  --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
  --endpoint-url http://localhost:8000


