data "okta_org_metadata" "test" {}
locals {
  okta_org_url = try(
    data.okta_org_metadata.test.domains.alternate,
    data.okta_org_metadata.test.domains.organization
  )
}
