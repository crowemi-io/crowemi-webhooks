resource "google_storage_bucket" "this" {
  count                       = var.env == "prod" ? 1 : 0
  name                        = local.service
  location                    = local.region
  uniform_bucket_level_access = true
}
