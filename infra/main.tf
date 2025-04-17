locals {
  region       = var.google_region
  service      = var.service_name
  project      = var.google_project_id
  env          = lower(var.env)
  name         = "${local.service}-${local.env}"
}

resource "google_cloud_run_v2_service" "this" {
  name         = local.name
  project      = local.project
  location     = local.region
  launch_stage = "BETA"
  ingress      = "INGRESS_TRAFFIC_ALL"
  template {
    containers {
      image = "us-west1-docker.pkg.dev/${local.project}/crowemi-io/${local.name}:${var.docker_image_tag}"
      ports {
        container_port = 8003
      }
      env {
        name = "CONFIG"
        value_source {
          secret_key_ref {
            secret  = data.google_secret_manager_secret.this.secret_id
            version = "latest"
          }
        }
      }
    }
    scaling {
      max_instance_count = 1
    }
    vpc_access {
      network_interfaces {
        network    = "crowemi-io-network"
        subnetwork = "crowemi-io-subnet-01"
      }
      egress = "ALL_TRAFFIC"
    }
    service_account = google_service_account.this.email
  }
}
