module "food_search" {
  source = "../../modules/s3"

  bucket_name = var.bucket_name
}
