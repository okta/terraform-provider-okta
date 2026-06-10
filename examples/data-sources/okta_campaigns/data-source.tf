# List all campaigns (default sort by created)
data "okta_campaigns" "all" {}

# Example: filter by status
# data "okta_campaigns" "active" {
#   filter = "status eq \"ACTIVE\""
# }

# Example: limit and order
# data "okta_campaigns" "recent" {
#   limit    = 10
#   order_by = ["created"]
# }

output "campaign_ids" {
  value = [for c in data.okta_campaigns.all.campaigns : c.id]
}

output "campaign_names" {
  value = [for c in data.okta_campaigns.all.campaigns : c.name]
}
