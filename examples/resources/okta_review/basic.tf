resource "okta_user" "test" {
  first_name = "TestAcc"
  last_name  = "Smith"
  login      = "testAcc-replace_with_uuid@example.com"
  email      = "testAcc-replace_with_uuid@example.com"
}

resource "okta_review" "test" {
  campaign_id = "ici1mo8jm7rv8vTXv1d7"
  reviewer_id = okta_user.test.id
  review_ids = [
    "icr1mobpl8z8xIPLv1d7"
  ]
  reviewer_level = "FIRST"
  note           = "John Smith is on leave for this month. His manager Tim will be the reviewer instead."
}