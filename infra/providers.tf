provider "google" {
  credentials = var.google_credentials
  project     = var.google_project_id
  region      = "us-west1"
}
