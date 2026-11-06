# Enough groups and mappings for refreshes to be served from listings: ten groups
# start the org's group listing, and two mappings start the app's mapping listing.
# depends_on chains fix the order of reads, so the same IDs are read directly
# before each listing on every run and the recorded requests replay exactly.

resource "okta_app_swa" "listing" {
  app_links_json = jsonencode(
    {
      deelhr_link = true
    }
  )
  label                   = "testAcc_replace_with_uuid"
  preconfigured_app       = "deelhr"
  status                  = "ACTIVE"
  user_name_template_type = "BUILT_IN"
}

resource "okta_group" "g0" {
  name = "testAcc_replace_with_uuid_0"
}

resource "okta_group" "g1" {
  name       = "testAcc_replace_with_uuid_1"
  depends_on = [okta_group.g0]
}

resource "okta_group" "g2" {
  name       = "testAcc_replace_with_uuid_2"
  depends_on = [okta_group.g1]
}

resource "okta_group" "g3" {
  name       = "testAcc_replace_with_uuid_3"
  depends_on = [okta_group.g2]
}

resource "okta_group" "g4" {
  name       = "testAcc_replace_with_uuid_4"
  depends_on = [okta_group.g3]
}

resource "okta_group" "g5" {
  name       = "testAcc_replace_with_uuid_5"
  depends_on = [okta_group.g4]
}

resource "okta_group" "g6" {
  name       = "testAcc_replace_with_uuid_6"
  depends_on = [okta_group.g5]
}

resource "okta_group" "g7" {
  name       = "testAcc_replace_with_uuid_7"
  depends_on = [okta_group.g6]
}

resource "okta_group" "g8" {
  name       = "testAcc_replace_with_uuid_8"
  depends_on = [okta_group.g7]
}

resource "okta_group" "g9" {
  name       = "testAcc_replace_with_uuid_9"
  depends_on = [okta_group.g8]
}

resource "okta_push_group" "p0" {
  app_id                         = okta_app_swa.listing.id
  source_group_id                = okta_group.g0.id
  status                         = "ACTIVE"
  delete_target_group_on_destroy = true
  depends_on                     = [okta_group.g9]
}

resource "okta_push_group" "p1" {
  app_id                         = okta_app_swa.listing.id
  source_group_id                = okta_group.g1.id
  status                         = "ACTIVE"
  delete_target_group_on_destroy = true
  depends_on                     = [okta_push_group.p0]
}

resource "okta_group_memberships" "g0" {
  group_id        = okta_group.g0.id
  users           = []
  track_all_users = true
  depends_on      = [okta_push_group.p1]
}
