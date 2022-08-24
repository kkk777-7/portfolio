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
