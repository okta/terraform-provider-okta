data "okta_ai_agent_providers_staging" "test" {
  id          = "aicufjl5wp7W2cZl5ao5"
  provider_id = "aipued6hgSvnXYSVPao5"
}

output "staging_type" {
  description = "Type of the staging item"
  value       = data.okta_ai_agent_providers_staging.test.type
}