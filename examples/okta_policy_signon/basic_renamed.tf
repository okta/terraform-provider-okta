data okta_group all {
  name = "Everyone"
}

resource okta_policy_signon test {
  name            = "testAccUpdated_replace_with_uuid"
  status          = "ACTIVE"
  description     = "Terraform Acceptance Test SignOn Policy"
  groups_included = ["${data.okta_group.all.id}"]
}
