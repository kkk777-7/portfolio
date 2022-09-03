module "food_search" {
  source = "../../modules/samprepare"

  bucket_name            = var.bucket_name
  iam_policy_name        = var.iam_policy_name
  iam_policy_description = var.iam_policy_description
  iam_role_name          = var.iam_role_name
  iam_role_description   = var.iam_role_description
  repo_name              = var.repo_name
}

resource "aws_ssm_parameter" "channel_secret" {
  name        = var.channel_secret_name
  description = var.channel_secret_description
  type        = "SecureString"
  value       = var.channel_secret_value
}

resource "aws_ssm_parameter" "channel_token" {
  name        = var.channel_token_name
  description = var.channel_token_description
  type        = "SecureString"
  value       = var.channel_token_value
}

resource "aws_ssm_parameter" "hotpepper_key" {
  name        = var.hotpepper_key_name
  description = var.hotpepper_key_description
  type        = "SecureString"
  value       = var.hotpepper_key_value
}

resource "aws_ssm_parameter" "google_key" {
  name        = var.google_key_name
  description = var.google_key_description
  type        = "SecureString"
  value       = var.google_key_value
}

resource "aws_dynamodb_table" "users" {
  name           = var.dynamodb_name
  billing_mode   = "PROVISIONED"
  read_capacity  = 1
  write_capacity = 1
  hash_key       = "user_line_id"
  attribute {
    name = "user_line_id"
    type = "S"
  }
}
