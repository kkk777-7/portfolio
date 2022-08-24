data "aws_caller_identity" "self" {}

########################
#  s3
########################
resource "aws_s3_bucket" "aws_sam_template" {
  bucket = var.bucket_name
  versioning {
    enabled = true
  }
  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm = "AES256"
      }
    }
  }
}

########################
#  iam
########################

data "aws_iam_policy_document" "sam_actions" {
  statement {
    actions = [
      "iam:*",
      "s3:*",
      "lambda:*",
      "ecr:*",
      "cloudformation:*",
      "apigateway:*"
    ]
    resources = [
      "*"
    ]
  }
}

resource "aws_iam_policy" "sam_actions" {
  name        = var.iam_policy_name
  description = var.iam_policy_description
  policy      = data.aws_iam_policy_document.sam_actions.json
}

data "aws_iam_policy_document" "github_actions_openid_connect" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRoleWithWebIdentity"]
    principals {
      type        = "Federated"
      identifiers = ["arn:aws:iam::${data.aws_caller_identity.self.account_id}:oidc-provider/token.actions.githubusercontent.com"]
    }
    condition {
      test     = "StringEquals"
      variable = "token.actions.githubusercontent.com:aud"
      values   = ["sts.amazonaws.com"]
    }
    condition {
      test     = "StringLike"
      variable = "token.actions.githubusercontent.com:sub"
      values   = [var.repo_name]
    }
  }
}

resource "aws_iam_role" "sam_actions" {
  name               = var.iam_role_name
  description        = var.iam_role_description
  assume_role_policy = data.aws_iam_policy_document.github_actions_openid_connect.json
}

resource "aws_iam_role_policy_attachment" "sam_actions_attach" {
  role       = aws_iam_role.sam_actions.name
  policy_arn = aws_iam_policy.sam_actions.arn
}
