data "okta_ai_agent_providers_profile_mapping" "test" {
  provider_id = "aipt6xdwmimVy0xBcao5"
}

output "profile_mappings" {
  description = "Profile mappings for the AI agent provider"
  value = [
    for item in data.okta_ai_agent_providers_profile_mapping.test.items : {
      target_field      = item.target_field
      source_expression = item.source_expression
      per_platform      = item.per_platform
    }
  ]
}