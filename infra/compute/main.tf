module "lambda" {
  source                  = "./lambda"
  lambda_bucket           = var.lambda_bucket
  repo_collector_role_arn = var.repo_collector_role_arn

  # Cuando conectemos a la VPC
  #security_group_ids      = [module.security_groups.main_security_group_id]
  security_group_ids      = ["sg-0ea2aa6ce985f84f6"] #default
  subnet_ids              = var.subnet_ids
}