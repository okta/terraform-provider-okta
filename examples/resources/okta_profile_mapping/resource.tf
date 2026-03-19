data "okta_user_profile_mapping_source" "user" {}

resource "okta_profile_mapping" "example" {
  source_id          = "<source id>"
  target_id          = data.okta_user_profile_mapping_source.user.id
  delete_when_absent = true

  mappings {
    id         = "firstName"
    expression = "appuser.firstName"
  }

  mappings {
    id         = "lastName"
    expression = "appuser.lastName"
  }

  mappings {
    id         = "email"
    expression = "appuser.email"
  }

  mappings {
    id         = "login"
    expression = "appuser.email"
  }
}
