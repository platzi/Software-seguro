data "aws_s3_object" "get-metrics" {
  bucket = var.lambda_bucket
  key    = "get-metrics.zip"
}

resource "aws_lambda_function" "get-metrics" {
  function_name    = "get-metrics"
  handler          = "bootstrap"
  runtime          = "provided.al2"
  s3_bucket        = var.lambda_bucket
  timeout          = 300
  s3_key           = "get-metrics.zip"
  role             = var.repo_collector_role_arn
  source_code_hash = data.aws_s3_object.get-metrics.version_id

  environment {
    variables = {
      DEMO     = "DEMO"
    }
  }
}

output "get-metrics_invoke_arn" {
  value = aws_lambda_function.get-metrics.invoke_arn
}

output "get-metrics_lambda_name" {
  value = aws_lambda_function.get-metrics.function_name
}