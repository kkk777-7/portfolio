variable "bucket_name" {
  description = "Name of the s3 bucket. Must be unique."
  type        = string
}

variable "iam_policy_name" {
  description = "Name of the IAM Policy. Must be unique."
  type        = string
}

variable "iam_policy_description" {
  description = "Description of the IAM Policy."
  type        = string
}

variable "repo_name" {
  description = "Name of the Github Repository."
  type        = string
}

variable "iam_role_name" {
  description = "Name of the IAM role. Must be unique."
  type        = string
}

variable "iam_role_description" {
  description = "Description of the IAM role."
  type        = string
}
