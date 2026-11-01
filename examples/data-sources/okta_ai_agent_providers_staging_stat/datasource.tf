data "okta_ai_agent_providers_staging_stat" "test" {
  provider_id = "aipty2bdDsW4txI6Nao5"
}

output "staging_counts" {
  description = "Counts of staging items by type for the AI agent provider"
  value = {
    duplicate = data.okta_ai_agent_providers_staging_stat.test.items[0].duplicate_count
    ignored   = data.okta_ai_agent_providers_staging_stat.test.items[0].ignored_count
    new       = data.okta_ai_agent_providers_staging_stat.test.items[0].new_count
  }
}
