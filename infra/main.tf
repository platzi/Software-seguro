terraform {
  backend "s3" {
    bucket = "platzi-tf-security"
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

module "compute" {
  source                  = "./compute"
  lambda_bucket           = module.s3.lambda_bucket
  repo_collector_role_arn = module.iam.repo_collector_role_arn

  # To connect to RDS
  # subnet_ids              = [module.vpc.main_subnet_id, module.vpc.availability_subnet_id]
  subnet_ids              = ["subnet-038e3c64ac224e62e", "subnet-0317446266e5ccaea", "subnet-043cb47352e975b2b"]
}