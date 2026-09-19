# This test targets an email domain that has already been created AND verified
# out of band (DNS records published + propagated). Verifying an already-verified
# email domain must succeed (the provider detects the VERIFIED status and skips
# the verify call). A fresh email domain cannot be verified inside an acceptance
# test because there is no opportunity to publish DNS records mid-run.
#
# When recording: export TF_VAR_email_domain_id=<your verified email domain id>.
# After recording: set the default below to the recorded id so replay reproduces
# the same request URLs.
variable "email_domain_id" {
  type        = string
  description = "ID of a pre-existing, already-verified Okta email domain"
  default     = "OeD14kbe6kzsL858H698"
}

resource "okta_email_domain_verification" "test" {
  email_domain_id = var.email_domain_id
}
