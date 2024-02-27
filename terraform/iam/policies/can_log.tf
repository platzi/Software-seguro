resource "aws_iam_policy" "can_log" {
  name        = "can-log-platzi"
  path        = "/"
  description = "Allow log to cloudwatch"

  policy = jsonencode(
    {
      Version : "2012-10-17",
      Statement : [
        {
          Effect = "Allow",
          Action = [
            "logs:CreateLogGroup",
            "logs:CreateLogStream",
            "logs:PutLogEvents",
          ],
          Resource = "*"
        }
      ]
  })
}

output "can_log_arn" {
  value = aws_iam_policy.can_log.arn
}