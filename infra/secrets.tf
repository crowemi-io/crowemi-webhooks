data "google_secret_manager_secret" "this" {
  secret_id = "${local.service}-${local.env}"

}
resource "google_secret_manager_secret_iam_member" "this" {
  secret_id = data.google_secret_manager_secret.this.secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.this.email}"
}
