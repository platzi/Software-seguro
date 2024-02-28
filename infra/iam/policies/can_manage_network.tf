// role for solve InvalidParameterValueException: The provided execution role does not have permissions to call CreateNetworkInterface on EC2
resource "aws_iam_policy" "can_manage_network_interfaces" {
  name        = "can-manage-network-interfaces"
  description = "Allows the ability to manage network interfaces"

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Action = [
          "ec2:CreateNetworkInterface",
          "ec2:DeleteNetworkInterface",
          "ec2:DescribeNetworkInterfaces",
          "ec2:DetachNetworkInterface",
          "ec2:AttachNetworkInterface",
        ]
        Effect   = "Allow",
        Resource = "*",
      },
    ],
  })
}

output "can_manage_network_interfaces_arn" {
  value = aws_iam_policy.can_manage_network_interfaces.arn
}