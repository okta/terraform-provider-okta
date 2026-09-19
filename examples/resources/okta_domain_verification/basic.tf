# This test targets a custom domain that has already been created AND verified
# out of band (DNS records published + propagated). Verifying an already-verified
# domain must succeed (the provider detects the VERIFIED status and skips the
# verify call). A fresh domain cannot be verified inside an acceptance test
# because there is no opportunity to publish DNS records mid-run.
#
# When recording: export TF_VAR_domain_id=<your verified domain id>.
# After recording: set the default below to the recorded id so replay reproduces
# the same request URLs.
variable "domain_id" {
  type        = string
  description = "ID of a pre-existing, already-verified Okta custom domain"
  default     = "OcD14k9oy3uuzBt6Y698"
}

resource "okta_domain_verification" "test" {
  domain_id = var.domain_id
}
