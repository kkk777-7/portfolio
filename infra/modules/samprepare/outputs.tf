output "s3_arn" {
  description = "ARN of the bucket"
  value       = aws_s3_bucket.aws_sam_template.arn
}

output "s3_name" {
  description = "Name (id) of the bucket"
  value       = aws_s3_bucket.aws_sam_template.id
}
