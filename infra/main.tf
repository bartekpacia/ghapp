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
  name            = var.project_name
  project_id      = var.project_id
  billing_account = var.billing_account_id

  # labels = {
  #   "firebase" = "enabled"
  # }
}

# resource "google_project_service" "artifact_registry" {
#   project = google_project.my_project.project_id
#   service = "artifactregistry.googleapis.com"
# }

# resource "google_project_service" "cloud_run" {
#   project = google_project.my_project.project_id
#   service = "run.googleapis.com"
# }
