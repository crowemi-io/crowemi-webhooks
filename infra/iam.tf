resource "google_service_account" "this" {
  account_id   = "srv-${local.service}-${local.env}"
  display_name = "srv_${replace(local.service, "-", "_")}_${lower(local.env)}"
  description  = "A service account for ${local.service} ${local.env}"
}


resource "google_cloud_run_service_iam_policy" "unauthenticated" {
  location    = local.region
  project     = local.project
  service     = google_cloud_run_v2_service.this.name

  policy_data = data.google_iam_policy.unauthenticated_access.policy_data
}

data "google_iam_policy" "unauthenticated_access" {
  binding {
    role = "roles/run.invoker"
    members = [
      "allUsers",
    ]
  }
}

resource "google_project_iam_member" "crowemi-log" {
  # TODO: move this to crowemi-log module
  project = local.project
  role    = "roles/pubsub.publisher"
  member  = "serviceAccount:${google_service_account.this.email}"
}
