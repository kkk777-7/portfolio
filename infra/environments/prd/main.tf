module "food_search" {
  source = "../../modules/samprepare"

  bucket_name            = var.bucket_name
  iam_policy_name        = var.iam_policy_name
  iam_policy_description = var.iam_policy_description
  iam_role_name          = var.iam_role_name
  iam_role_description   = var.iam_role_description
  repo_name              = var.repo_name
}
