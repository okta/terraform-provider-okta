resource "okta_domain" "test" {
  name   = "testAcc-replace_with_uuid.example.com"
  verify = false
}

resource "okta_domain_certificate" "test" {
  domain_id = okta_domain.test.id
  type      = "PEM"

  #certificate = file("cert.pem")
  certificate = <<CERT
-----BEGIN CERTIFICATE-----
MIIFNzCCBB+gAwIBAgISBAXomJWRama3ypu8TIxdA9wzMA0GCSqGSIb3DQEBCwUA
...
NSgRtSXq11j8O4JONi8EXe7cEtvzUiLR5PL3itsK2svtrZ9jIwQ95wOPaA==
-----END CERTIFICATE-----
CERT

  #certificate_chain = file("chain.pem")
  certificate_chain = <<CHAIN
-----BEGIN CERTIFICATE-----
MIIFFjCCAv6gAwIBAgIRAJErCErPDBinU/bWLiWnX1owDQYJKoZIhvcNAQELBQAw
...
Dfvp7OOGAN6dEOM4+qR9sdjoSYKEBpsr6GtPAQw4dy753ec5
-----END CERTIFICATE-----
  CHAIN

  #private_key = file("privkey.pem")
  private_key = <<PK
-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQC5cyk6x63iBJSW
...
nUFLNE8pXSnsqb0eOL74f3uQ
-----END PRIVATE KEY-----
PK
}
