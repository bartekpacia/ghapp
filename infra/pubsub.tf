

# resource "google_pubsub_schema" "analytics" {
#   project    = google_project.default.project_id
#   name       = "analytics-schema"
#   type       = "PROTOCOL_BUFFER"
#   definition = file("${path.module}/analytics-schema.proto")
# }

# resource "google_pubsub_topic" "analytics" {
#   project = google_project.default.project_id
#   name    = "analytics-topic"

#   depends_on = [google_pubsub_schema.analytics]
#   schema_settings {
#     schema   = "projects/${google_project.default.project_id}/schemas/${google_pubsub_schema.analytics.name}"
#     encoding = "JSON"
#   }
# }

# resource "google_pubsub_subscription" "analytics" {
#   project = google_project.default.project_id
#   name    = "analytics-sub"
#   topic   = google_pubsub_topic.analytics.id

#   ack_deadline_seconds = 10

#   push_config {
#     push_endpoint = "https://webhook.site/d017e8ca-c877-4dcf-8d18-8c314d9282ab"
#   }
# }
