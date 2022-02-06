terraform {
  required_version = "~> 1.foo.bar"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.34.0"
    }
  }
}