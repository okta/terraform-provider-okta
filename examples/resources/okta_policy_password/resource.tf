resource "okta_policy_password" "example" {
  name        = "example"
  type        = "PASSWORD"
  status      = "ACTIVE"
  description = "Example"

  conditions {
    people {
      groups {
        include = [data.okta_group.everyone.id]
      }
    }
  }

  settings {
    password {
      age {
        history_count = 4
      }
    }
  }
}
