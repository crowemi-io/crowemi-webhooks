resource "google_service_account" "this" {
  account_id   = "srv-${local.service}-${local.env}"
  display_name = "srv_${replace(local.service, "-", "_")}_${lower(local.env)}"
  description  = "A service account for ${local.service} ${local.env}"
}

# data "google_iam_policy" "this" {
#   binding {
#     role = "roles/run.invoker"
#     members = [
#       "serviceAccount:${google_service_account.this.email}",
#     ]
#   }
# }
# resource "google_cloud_run_service_iam_policy" "this" {
#   location = google_cloud_run_v2_service.this.location
#   project  = google_cloud_run_v2_service.this.project
#   service  = google_cloud_run_v2_service.this.name

#   policy_data = data.google_iam_policy.this.policy_data
# }

resource "google_project_iam_member" "crowemi-log" {
  project = local.project
  role    = "roles/pubsub.publisher"
  member  = "serviceAccount:${google_service_account.this.email}"
}
