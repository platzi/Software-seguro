data "aws_s3_object" "handle_github_webhook" {
  bucket = var.lambda_bucket
  key    = "handle-github-webhook.zip"
}

resource "aws_lambda_function" "handle_github_webhook" {
  function_name    = "handle-github-webhook"
  handler          = "bootstrap"
  runtime          = "provided.al2"
  s3_bucket        = var.lambda_bucket
  timeout          = 300
  s3_key           = "handle-github-webhook.zip"
  role             = var.repo_collector_role_arn
  source_code_hash = data.aws_s3_object.handle_github_webhook.version_id

  vpc_config {
    security_group_ids = var.security_group_ids
    subnet_ids         = var.subnet_ids
  }

  environment {
    variables = {
      DEMO = "DEMO"
    }
  }
}

output "handle_github_webhook_invoke_arn" {
  value = aws_lambda_function.handle_github_webhook.invoke_arn
}

output "handle_github_webhook_lambda_name" {
  value = aws_lambda_function.handle_github_webhook.function_name
}