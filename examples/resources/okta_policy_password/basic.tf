data "okta_group" "all" {
  name = "Everyone"
}

resource "okta_policy_password" "test" {
  name        = "testAcc_replace_with_uuid"
  type        = "PASSWORD"
  status      = "ACTIVE"
  description = "Terraform Acceptance Test Password Policy"

  conditions {
    people {
      groups {
        include = [data.okta_group.all.id]
      }
    }
  }

  settings {
    password {
      age {
        history_count = 5
      }
    }
  }
}
