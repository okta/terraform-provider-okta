resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name  = "Smith"
  login      = "testAcc-replace_with_uuid@example.com"
  email      = "testAcc-replace_with_uuid@example.com"
}

resource "okta_review" "test" {
  campaign_id = "icizigd86iM9sOcbN1d6"
  reviewer_id = okta_user.test.id
  review_ids = [
    "icrztblxbBFiVKepb1d6"
  ]
  reviewer_level = "FIRST"
  note           = "John Smith is on leave for this month. His manager Tim will be the reviewer instead."
}