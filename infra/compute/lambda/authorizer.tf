data "aws_s3_object" "authorizer" {
  bucket = var.lambda_bucket
  key    = "authorizer.zip"
}

resource "aws_lambda_function" "authorizer" {
  function_name    = "authorizer"
  handler          = "bootstrap"
  runtime          = "provided.al2"
  s3_bucket        = var.lambda_bucket
  timeout          = 300
  s3_key           = "authorizer.zip"
  role             = var.repo_collector_role_arn
  source_code_hash = data.aws_s3_object.authorizer.version_id

  vpc_config {
    security_group_ids = var.security_groups_ids
    subnet_ids          = var.subnet_ids
  }

  environment {
    variables = {
      DEMO = "DEMO"
    }
  }


}

output "authorizer_invoke_arn" {
  value = aws_lambda_function.authorizer.invoke_arn
}

output "authorizer_lambda_name" {
  value = aws_lambda_function.authorizer.function_name
}