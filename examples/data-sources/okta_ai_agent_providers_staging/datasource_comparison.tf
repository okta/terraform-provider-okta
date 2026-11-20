data "okta_ai_agent_providers_staging" "test" {
  id          = "replace_with_staging_item_id"
  provider_id = "replace_with_provider_id"
}

output "agent_comparison" {
  description = "Comparison between existing and incoming agents"
  value = {
    existing_agent = {
      id       = try(data.okta_ai_agent_providers_staging.test.existing_agent[0].id, null)
      name     = try(data.okta_ai_agent_providers_staging.test.existing_agent[0].profile[0].name, null)
      platform = try(data.okta_ai_agent_providers_staging.test.existing_agent[0].profile[0].platform, null)
    }
    incoming_agent = {
      id       = try(data.okta_ai_agent_providers_staging.test.incoming_agent[0].id, null)
      name     = try(data.okta_ai_agent_providers_staging.test.incoming_agent[0].profile[0].name, null)
      platform = try(data.okta_ai_agent_providers_staging.test.incoming_agent[0].profile[0].platform, null)
    }
    requires_reconciliation = (
      try(data.okta_ai_agent_providers_staging.test.existing_agent[0].profile[0].name, null) !=
      try(data.okta_ai_agent_providers_staging.test.incoming_agent[0].profile[0].name, null)
    )
  }
}