terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "5.37.0"
    }
  }
}

provider "google" {
  project     = var.project_id
  region      = var.region
  zone        = var.zone
  credentials = file(var.credentials_file)

  user_project_override = true
}

resource "google_project" "default" {
  provider        = google
  name            = "Bee CI"
  project_id      = var.project_id
  billing_account = var.billing_account_id

  # labels = {
  #   "firebase" = "enabled"
  # }
}

#resource "google_project_service" "default" {
#  provider = google
#  project  = google_project.default.project_id
#  for_each = toset([
#    "serviceusage.googleapis.com",
#    #"artifactregistry.googleapis.com",
#    #"run.googleapis.com",
#  ])
#  service = each.key
#
#  disable_on_destroy = true
#}


resource "google_cloud_run_service" "default" {
  name     = "main-cloud-run-service"
  location = var.region

  template {
    spec {
      containers {
        image = "gcr.io/bee-ci/bee-ci:latest"
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}

