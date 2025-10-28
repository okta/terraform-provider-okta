resource "okta_end_user_my_requests" "example" {
  entry_id = "cen123456789abcdefgh"

  requester_field_values {
    id    = "abcdefgh-0123-4567-8910-hgfedcba123"
    value = "I need access to complete my certification."
  }

  requester_field_values {
    id    = "ijklmnop-a12b2-c3d4-e5f6-abcdefghi"
    value = "For 5 days"
  }

  requester_field_values {
    id    = "tuvwxyz-0123-456-8910-zyxwvut0123"
    value = "Yes"
  }
}
