variable "lambda_bucket" {
  type = string
}

variable "repo_collector_role_arn" {
  type = string
}

# To connect to RDS
variable "subnet_ids" {
  type = list(string)
}

variable "security_group_ids" {
  type = list(string)
}