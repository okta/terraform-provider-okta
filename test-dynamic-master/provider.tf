terraform {
  required_providers {
    okta = {
      source = "okta/okta"
    }
  }
}

provider "okta" {
  org_name    = "dcp-tf-testing-2026-01-06"
  base_url    = "oktapreview.com"
  api_token   = "00tH3Vj3gaSzfokbFECgKqYTltXy6eaAnmxbGojH-M"
  parallelism = 1
}
