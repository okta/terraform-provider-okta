variable "partition" {
  description = "The okta TLD without the .com"
  default     = "oktapreview"
}

variable "org_name" {
  description = "The Okta sub-domain"
}

variable "org_id" {
  description = "The organization ID, can be obtained from /api/v1/org"
}

variable "api_token" {
  description = "The administrator API token used to run terraform"
}
