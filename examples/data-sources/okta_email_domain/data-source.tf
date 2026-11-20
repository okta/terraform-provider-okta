data "okta_email_domain" "example" {
  id = "<email_domain_id>"
}

output "email_domain_domain" {
  value = data.okta_email_domain.example.domain
}
