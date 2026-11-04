data "okta_ai_agent_providers_staging" "test" {
  id          = "replace_with_staging_item_id"
  provider_id = "replace_with_provider_id"
}

output "staging_type" {
  description = "Type of the staging item"
  value       = data.okta_ai_agent_providers_staging.test.type
}