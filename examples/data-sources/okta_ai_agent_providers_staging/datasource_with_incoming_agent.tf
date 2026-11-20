data "okta_ai_agent_providers_staging" "test" {
  id          = "replace_with_staging_item_id"
  provider_id = "replace_with_provider_id"
}

output "incoming_agent_info" {
  description = "Information about the incoming agent in the staging item"
  value = {
    id          = try(data.okta_ai_agent_providers_staging.test.incoming_agent[0].id, null)
    name        = try(data.okta_ai_agent_providers_staging.test.incoming_agent[0].profile[0].name, null)
    description = try(data.okta_ai_agent_providers_staging.test.incoming_agent[0].profile[0].description, null)
    platform    = try(data.okta_ai_agent_providers_staging.test.incoming_agent[0].profile[0].platform, null)
    status      = try(data.okta_ai_agent_providers_staging.test.incoming_agent[0].status, null)
  }
}