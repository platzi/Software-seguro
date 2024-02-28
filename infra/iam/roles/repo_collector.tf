# Repo collector role manage information about the repositories
# It is used by the lambda function to collect information about the repositories
# and provide it to the other services
data "aws_iam_policy_document" "instance_assume_role_policy" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "repo_collector" {
  name               = "repo-collector-platzi"
  path               = "/"
  assume_role_policy = data.aws_iam_policy_document.instance_assume_role_policy.json
}

resource "aws_iam_role_policy_attachment" "can_access_rds" {
  role       = aws_iam_role.repo_collector.name
  policy_arn = var.rds_policy_arn
}

resource "aws_iam_role_policy_attachment" "can_log" {
  role       = aws_iam_role.repo_collector.name
  policy_arn = var.log_policy_arn
}

output "repo_collector_role_arn" {
  value = aws_iam_role.repo_collector.arn
}