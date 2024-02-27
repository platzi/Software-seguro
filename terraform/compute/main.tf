module "lambda" {
  source = "./lambda"
  lambda_bucket           = var.lambda_bucket
  repo_collector_role_arn = var.repo_collector_role_arn
}