# Configure AWS provider
provider "aws" {
  region = "us-west-2"  # Change this to your desired region
}

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.0"
    }
  }
}

resource "aws_dynamodb_table" "todo_table" {
  name           = "TodoTable"            # The name of the DynamoDB table
  billing_mode   = "PROVISIONED"          # Provisioned or on-demand mode
  read_capacity  = 5                      # Number of read capacity units
  write_capacity = 5                      # Number of write capacity units

  # Define the key schema
  hash_key       = "UserID"               # Partition key (Primary Key)
  range_key      = "ID"                   # Sort key (Secondary Key)

  # Define the table attributes
  attribute {
    name = "UserID"
    type = "S"  # S stands for String type
  }

  attribute {
    name = "ID"
    type = "S"  # S stands for String type
  }

  tags = {
    Environment = "dev"  # Optional, for tagging the resource
    Project     = "DynamoDBExample"
  }
}

resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0"

  tags = {
    Name = "main"
  }
}

resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.main.id

}