module "policies" {
  source = "./policies"
}

module "roles" {
  source         = "./roles"
  rds_policy_arn = module.policies.can_access_rds_arn
  log_policy_arn = module.policies.can_log_arn
}

output "repo_collector_role_arn" {
  value = module.roles.repo_collector_role_arn
}