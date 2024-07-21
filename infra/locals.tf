locals {
  image_id = "${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.default.repository_id}/bee-ci:latest"
}
