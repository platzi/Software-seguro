resource "aws_iam_policy" "can_get_db_password" {
  name        = "can-get-db-password"
  path        = "/"
  description = "Allow access to retrieve secrets from Secrets Manager"

  policy = jsonencode({
    Version : "2012-10-17",
    Statement : [
      {
        Effect : "Allow",
        Action : [
          "secretsmanager:GetSecretValue"
        ],
        Resource : [
            # Copy ARN from AWS
          "arn:aws:secretsmanager:us-east-2:533267258008:secret:rds!db-30adf398-fea8-4d75-9920-bab3601b979e-3ZcqbR"
        ]
      }
    ]
  })
}

output "can_get_db_password_arn" {
  value = aws_iam_policy.can_get_db_password.arn
}
