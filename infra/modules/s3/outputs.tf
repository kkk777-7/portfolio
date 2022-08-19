output "arn" {
  description = "ARN of the bucket"
  value       = aws_s3_bucket.aws_sam_template.arn
}

output "name" {
  description = "Name (id) of the bucket"
  value       = aws_s3_bucket.aws_sam_template.id
}
