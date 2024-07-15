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

  # Problems with user_project_override:
  #  https://github.com/hashicorp/terraform-provider-google/issues/14174
  user_project_override = false
}

resource "google_project" "default" {
  provider        = google
  name            = "Bee CI"
  project_id      = var.project_id
  billing_account = var.billing_account_id
}

variable "required_services" {
  description = "List of APIs necessary for this project"
  type        = list(string)
  default = [
    "cloudresourcemanager.googleapis.com", # cannot be enabled through Terraform ?
    "serviceusage.googleapis.com",         # cannot be enabled through Terraform ?
    "cloudbuild.googleapis.com",
  ]
}

resource "google_project_service" "default" {
  project  = google_project.default.project_id
  for_each = toset(var.required_services)
  service  = each.key

  disable_on_destroy = true
}

resource "google_artifact_registry_repository" "default" {
  project       = google_project.default.project_id
  location      = var.region
  repository_id = "bee-ci"
  format        = "DOCKER"
  description   = "Default repo for our Docker images"
}

resource "google_cloudbuild_trigger" "default" {
  project = google_project.default.project_id
  trigger_template {
    branch_name = "master"
    repo_name   = "your-github-repo"

  }

  # TODO: Migrate to
  #  https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/cloudbuild_trigger#example-usage---cloudbuild-trigger-build
  filename = "cloudbuild.yaml"
}

resource "google_cloud_run_service" "default" {
  project  = google_project.default.project_id
  name     = "main-cloud-run-service"
  location = var.region

  template {
    spec {
      containers {
        image = "gcr.io/bee-ci/bee-ci:latest"
        ports {
          container_port = 8080
        }
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}

data "google_iam_policy" "noauth" {
  binding {
    role = "roles/run.invoker"
    members = [
      "allUsers",
    ]
  }
}

resource "google_cloud_run_service_iam_policy" "noauth" {
  location = google_cloud_run_service.default.location
  project  = google_cloud_run_service.default.project
  service  = google_cloud_run_service.default.name

  policy_data = data.google_iam_policy.noauth.policy_data
}

