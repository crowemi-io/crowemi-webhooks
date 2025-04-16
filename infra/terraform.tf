terraform {
  cloud {

    organization = "crowemi-io"

    workspaces {
      name = "crowemi-bot-dev"
    }
  }
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "5.20.0"
    }
  }
}
