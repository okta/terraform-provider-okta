# Complete example showing collection management with resources and assignments
resource "okta_collection" "slack_permissions" {
  name        = "Slack Bot Permissions"
  description = "Entitlements for in-house Slack bot"
}

# Add app with entitlements to the collection
resource "okta_collection_resource" "slack_app" {
  collection_id = okta_collection.slack_permissions.id
  resource_orn  = "orn:okta:idp:${var.org_id}:apps:slack:${var.app_id}"
  
  entitlements {
    id = okta_entitlement.admin_features.id
    values {
      id = okta_entitlement.admin_features.values[0].id
    }
  }
  
  entitlements {
    id = okta_entitlement.user_features.id
    values {
      id = okta_entitlement.user_features.values[0].id
    }
    values {
      id = okta_entitlement.user_features.values[1].id
    }
  }
}

# Assign the collection to a group
resource "okta_collection_assignment" "engineering_team" {
  collection_id   = okta_collection.slack_permissions.id
  principal_id    = okta_group.engineers.id
  principal_type  = "OKTA_GROUP"
  actor          = "ADMIN"
}

# Temporary assignment with expiration
resource "okta_collection_assignment" "contractor_access" {
  collection_id   = okta_collection.slack_permissions.id
  principal_id    = okta_user.contractor.id
  principal_type  = "OKTA_USER"
  actor          = "API"
  expiration_time = "2024-12-31T23:59:59Z"
  time_zone       = "America/Los_Angeles"
}
