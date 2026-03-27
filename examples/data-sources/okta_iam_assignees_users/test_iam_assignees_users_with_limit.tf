data "okta_iam_assignees_users" "test" {
  limit = 50
}

output "iam_assignees_users" {
  value = data.okta_iam_assignees_users.test.users
} 
