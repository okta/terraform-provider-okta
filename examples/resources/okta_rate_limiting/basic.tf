resource "okta_rate_limiting" "example" {
  default_mode = "ENFORCE"
  use_case_mode_overrides{
    login_page= "ENFORCE"
  }
}
