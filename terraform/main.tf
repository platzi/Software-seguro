terraform {
  backend "s3" {
    bucket = "camileniss-terraform"
    key    = "terraform.tfstate"
    region = "us-east-2"
  }
}

provider "aws" {
  region = "us-east-2"
}

module "iam" {
  source = "./iam"
}

module "s3" {
  source = "./s3"
}