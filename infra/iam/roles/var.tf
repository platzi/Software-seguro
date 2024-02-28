variable "rds_policy_arn" {
  type = string
}

variable "log_policy_arn" {
  type = string
}

# For connecting to RDS
variable "can_manage_network_policy_arn" {
  type = string
}