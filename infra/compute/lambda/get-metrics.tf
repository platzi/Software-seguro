data "aws_s3_object" "get_metrics" {
  bucket = var.lambda_bucket
  key    = "get-metrics.zip"
}

resource "aws_lambda_function" "get_metrics" {
  function_name    = "get-metrics"
  handler          = "bootstrap"
  runtime          = "provided.al2"
  s3_bucket        = var.lambda_bucket
  timeout          = 300
  s3_key           = "get-metrics.zip"
  role             = var.repo_collector_role_arn
  source_code_hash = data.aws_s3_object.get_metrics.version_id

  vpc_config {
    security_group_ids = var.security_group_ids
    subnet_ids         = var.subnet_ids
  }

  environment {
    variables = {
      DB_HOST     = "platzi-course.c9sqmwee2hyj.us-east-2.rds.amazonaws.com" # Update with current db host
      DB_PORT     = "5432"
      DB_USER     = "platzi"
      DB_PASSWORD = "rds!db-e4e77491-2a21-4121-9788-59fdbbe7f38c" # Update with secret name
      DB_NAME     = "postgres"
    }
  }
}

output "get_metrics_invoke_arn" {
  value = aws_lambda_function.get_metrics.invoke_arn
}

output "get_metrics_lambda_name" {
  value = aws_lambda_function.get_metrics.function_name
}