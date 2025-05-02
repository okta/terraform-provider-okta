resource "okta_realm_assignment" "test" {
  name                 = "TestAcc Example Realm Assignment Updated"
  priority             = 111
  status               = "ACTIVE"
  profile_source_id    = okta_idp_saml.test.id
  condition_expression = "user.profile.login.contains(\"@acctestupdated.com\")"
  realm_id             = okta_realm.test.id
}

resource "okta_realm" "test" {
  name       = "TestAcc Example Assignment Realm"
  realm_type = "DEFAULT"
}

resource "okta_idp_saml" "test" {
  name                     = "testAcc_replace_with_uuid"
  acs_type                 = "INSTANCE"
  sso_url                  = "https://idp.example.com"
  sso_destination          = "https://idp.example.com"
  sso_binding              = "HTTP-POST"
  username_template        = "idpuser.email"
  kid                      = okta_idp_saml_key.test.id
  issuer                   = "https://idp.example.com"
  request_signature_scope  = "REQUEST"
  response_signature_scope = "ANY"
  profile_master           = true
}

resource "okta_idp_saml_key" "test" {
  x5c = ["MIIDnjCCAoagAwIBAgIGAVG3MN+PMA0GCSqGSIb3DQEBBQUAMIGPMQswCQYDVQQGEwJVUzETMBEGA1UECAwKQ2FsaWZvcm5p\nYTEWMBQGA1UEBwwNU2FuIEZyYW5jaXNjbzENMAsGA1UECgwET2t0YTEUMBIGA1UECwwLU1NPUHJvdmlkZXIxEDAOBgNVBAMM\nB2V4YW1wbGUxHDAaBgkqhkiG9w0BCQEWDWluZm9Ab2t0YS5jb20wHhcNMTUxMjE4MjIyMjMyWhcNMjUxMjE4MjIyMzMyWjCB\njzELMAkGA1UEBhMCVVMxEzARBgNVBAgMCkNhbGlmb3JuaWExFjAUBgNVBAcMDVNhbiBGcmFuY2lzY28xDTALBgNVBAoMBE9r\ndGExFDASBgNVBAsMC1NTT1Byb3ZpZGVyMRAwDgYDVQQDDAdleGFtcGxlMRwwGgYJKoZIhvcNAQkBFg1pbmZvQG9rdGEuY29t\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAtcnyvuVCrsFEKCwHDenS3Ocjed8eWDv3zLtD2K/iZfE8BMj2wpTf\nn6Ry8zCYey3mWlKdxIybnV9amrujGRnE0ab6Q16v9D6RlFQLOG6dwqoRKuZy33Uyg8PGdEudZjGbWuKCqqXEp+UKALJHV+k4\nwWeVH8g5d1n3KyR2TVajVJpCrPhLFmq1Il4G/IUnPe4MvjXqB6CpKkog1+ThWsItPRJPAM+RweFHXq7KfChXsYE7Mmfuly8s\nDQlvBmQyxZnFHVuiPfCvGHJjpvHy11YlHdOjfgqHRvZbmo30+y0X/oY/yV4YEJ00LL6eJWU4wi7ViY3HP6/VCdRjHoRdr5L/\nDwIDAQABMA0GCSqGSIb3DQEBBQUAA4IBAQCzzhOFkvyYLNFj2WDcq1YqD4sBy1iCia9QpRH3rjQvMKDwQDYWbi6EdOX0TQ/I\nYR7UWGj+2pXd6v0t33lYtoKocp/4lUvT3tfBnWZ5KnObi+J2uY2teUqoYkASN7F+GRPVOuMVoVgm05ss8tuMb2dLc9vsx93s\nDt+XlMTv/2qi5VPwaDtqduKkzwW9lUfn4xIMkTiVvCpe0X2HneD2Bpuao3/U8Rk0uiPfq6TooWaoW3kjsmErhEAs9bA7xuqo\n1KKY9CdHcFhkSsMhoeaZylZHtzbnoipUlQKSLMdJQiiYZQ0bYL83/Ta9fulr1EERICMFt3GUmtYaZZKHpWSfdJp9"]
}

