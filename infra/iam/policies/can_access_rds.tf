resource "aws_iam_policy" "can_access_rds" {
  name        = "can-access-rds"
  path        = "/"
  description = "Allow manage RDS databases for queries"

  policy = jsonencode(
    {
      Version : "2012-10-17",
      Statement : [
        {
          Effect : "Allow",
          Action : [
            "rds-db:connect"
          ],
          Resource : [
            "arn:aws:rds-db:us-east-2:*:*"
          ]
        }
      ]
    }
  )
}

output "can_access_rds_arn" {
  value = aws_iam_policy.can_access_rds.arn
}