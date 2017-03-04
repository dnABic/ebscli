provider "aws" {
  endpoints {
    s3  = "http://awsserver:5000"
    ec2 = "http://awsserver:5000"
    iam = "http://awsserver:5000"
    elb = "http://awsserver:5000"
  }

  access_key = "the_key"
  secret_key = "the_secret"
  region     = "us-east-1"
}

resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"

  tags {
    Name = "test"
  }
}
