terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "5.38.0"
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

  disable_on_destroy = false
}
