variable "lambda_bucket" {
  type = string
}

variable "repo_collector_role_arn" {
  type = string
}

variable "subnet_ids" {
  type = list(string)
}

variable "security_groups_ids" {
  type = list(string)
}