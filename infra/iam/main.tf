module "policies" {
  source = "./policies"
}

module "roles" {
  source         = "./roles"
  rds_policy_arn = module.policies.can_access_rds_arn
  log_policy_arn = module.policies.can_log_arn

  can_manage_network_policy_arn = module.policies.can_manage_network_interfaces_arn
  can_get_db_password_arn = module.policies.can_get_db_password_arn

  can_access_github_secret_arn = module.policies.can_access_github_secret_arn
}

output "repo_collector_role_arn" {
  value = module.roles.repo_collector_role_arn
}