terraform {
  required_version = "= 1.4.2"

  required_providers {
    okta = {
      source  = "terraform.local/local/okta"
      version = "3.44.0"
    }
  }
}

provider "okta" {
  org_name  = var.org_name
  base_url  = "${var.partition}.com"
  api_token = var.api_token
}

resource "okta_resource_set" "my_resources" {
  label       = "my_resources"
  description = "All my resources"

  resources = [
    "https://${var.org_name}.${var.partition}.com/api/v1/groups",
    "orn:oktapreview:idp:${var.org_id}:customizations",
    # "orn:oktapreview:idp:${var.org_id}:apps:workday",
    "https://${var.org_name}.${var.partition}.com/api/v1/apps?filter=name+eq+%22workday%22",
    # "https://${var.org_name}.${var.partition}.com/api/v1/apps?filter=name+eq+\"workday\"",
  ]
}
