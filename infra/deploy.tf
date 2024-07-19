# This file contains the deployment pipeline:
# GitHub -> Cloud Build -> Artifact Registry -> Cloud Run

#resource "google_cloudbuildv2_connection" "default" {
#  project  = google_project.default.project_id
#  location = var.region
#  name     = "my-github-connection-tf"
#
#  github_config {
#   app_installation_id = var.github_app_installation_id 
#  }
#}

resource "google_cloudbuild_trigger" "default" {
  project  = google_project.default.project_id
  location = var.region
  name     = "my-cloudbuild-trigger-tf"

  github {
    owner = "bartekpacia"
    name  = "ghapp"
    push {
      branch = ".*"
    }
  }

  build {
    step {
      name             = "ubuntu"
      args             = ["echo", "hello world!"]
      allow_exit_codes = [0]
    }

    step {
      name = "gcr.io/cloud-builders/docker"
      args = ["build", "-t", local.image_id, "."]
    }

    images = [local.image_id]

    options {
      logging = "CLOUD_LOGGING_ONLY"
    }
  }
}

# resource "google_cloudbuildv2_repository" "name" {

# }


resource "google_artifact_registry_repository" "default" {
  project       = google_project.default.project_id
  location      = var.region
  repository_id = "my-artifact-registry-repository-tf"
  format        = "DOCKER"
  description   = "Default repo for our Docker images"
}


resource "google_cloud_run_service" "default" {
  project  = google_project.default.project_id
  location = var.region
  name     = "my-cloud-run-service-tf"

  template {
    spec {
      containers {
        image = local.image_id
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
    role    = "roles/run.invoker"
    members = ["allUsers"]
  }
}

resource "google_cloud_run_service_iam_policy" "noauth" {
  project  = google_cloud_run_service.default.project
  location = google_cloud_run_service.default.location
  service  = google_cloud_run_service.default.name

  policy_data = data.google_iam_policy.noauth.policy_data
}
