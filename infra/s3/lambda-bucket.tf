resource "aws_s3_bucket" "lambda_bucket" {
  bucket = "camilaleniss-lambda-platzi"
}

resource "aws_s3_bucket_versioning" "lambda_bucket" {
  bucket = aws_s3_bucket.lambda_bucket.id
  versioning_configuration {
    status = "Enabled"
  }
}

output "lambda_bucket" {
  value = aws_s3_bucket.lambda_bucket.bucket
}