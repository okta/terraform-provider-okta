---
layout: 'okta'
page_title: 'Okta: okta_domain_certificate'
sidebar_current: 'docs-okta-resource-domain-certificate'
description: |-
  Manages certificate for the domain.
---

# okta_domain_certificate

Manages certificate for the domain.

This resource's `certificate`, `private_key`, and `certificate_chain` attributes
hold actual PEM values and can be referred to by other configs requiring
certificate and private key inputs. This is inline with TF's [best
practices](https://developer.hashicorp.com/terraform/plugin/sdkv2/best-practices/sensitive-state#don-t-encrypt-state)
of not encrypting state.

See [Let's Encrypt Certbot notes](#lets-encrypt-certbot) at the end of this
documentation for notes on how to generate a domain certificate with Let's
Encrypt Certbot

## Example Usage

```hcl
resource "okta_domain" "example" {
  name = "www.example.com"
}

resource "okta_domain_certificate" "test" {
  domain_id   = okta_domain.test.id
  type        = "PEM"

  certificate_chain = file("www.example.com/chain.pem")

  #certificate = file("www.example.com/cert.pem")
  certificate = <<CERT
-----BEGIN CERTIFICATE-----
MIIFNzCCBB+gAwIBAgISBAXomJWRama3ypu8TIxdA9wzMA0GCSqGSIb3DQEBCwUA
MDIxCzAJBgNVBAYTAlVTMRYwFAYDVQQKEw1MZXQncyBFbmNyeXB0MQswCQYDVQQD
EwJSMzAeFw0yMTAyMTAwNTEzMDVaFw0yMTA1MTEwNTEzMDVaMCQxIjAgBgNVBAMT
GWFuaXRhdGVzdC5zaWdtYW5ldGNvcnAudXMwggEiMA0GCSqGSIb3DQEBAQUAA4IB
DwAwggEKAoIBAQC5cyk6x63iBJSWvtgsOBqIxfO8euPHcRnyWsL9dsvnbNyOnyvc
qFWxdiW3sh2cItzYtoN1Zfgj5lWGOVXbHxP0VaNG9fHVX3+NHP6LFHQz92BzAYQm
pqi9zaP/aKJklk6LdPFbVLGhuZfm34+ijW9YsgLTKR2WTaZJK5QtamVVmP+VsSCl
a2ifFzjz2FCkMMEc/Y0zUyP+en/mbL71K+VnpZdlEC1s38EvjRTFKFZTKVw5wpWg
CZQq/AZYj9RxR23IIuRcUJ8TQ2pyoc3kIXPWjiIarSgBlA8G9kCsxgzXP2RyLwKr
IBIo+qyHweifpPYW28ipdSbPjiypAMdpbGLDAgMBAAGjggJTMIICTzAOBgNVHQ8B
Af8EBAMCBaAwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMAwGA1UdEwEB
/wQCMAAwHQYDVR0OBBYEFPVZKiovtIK4Av/IBUQeLUs29pT6MB8GA1UdIwQYMBaA
FBQusxe3WFbLrlAJQOYfr52LFMLGMFUGCCsGAQUFBwEBBEkwRzAhBggrBgEFBQcw
AYYVaHR0cDovL3IzLm8ubGVuY3Iub3JnMCIGCCsGAQUFBzAChhZodHRwOi8vcjMu
aS5sZW5jci5vcmcvMCQGA1UdEQQdMBuCGWFuaXRhdGVzdC5zaWdtYW5ldGNvcnAu
dXMwTAYDVR0gBEUwQzAIBgZngQwBAgEwNwYLKwYBBAGC3xMBAQEwKDAmBggrBgEF
BQcCARYaaHR0cDovL2Nwcy5sZXRzZW5jcnlwdC5vcmcwggEDBgorBgEEAdZ5AgQC
BIH0BIHxAO8AdgBc3EOS/uarRUSxXprUVuYQN/vV+kfcoXOUsl7m9scOygAAAXeK
kmOsAAAEAwBHMEUCIQDSudPEWXk969BT8yz3ag6BJWCMRU5tefEw9nXEQMsh5gIg
UmfGIuUlcNNI5PydVIHj+zns+SR8P7zfd3FIxW4gK0QAdQD2XJQv0XcwIhRUGAgw
lFaO400TGTO/3wwvIAvMTvFk4wAAAXeKkmOlAAAEAwBGMEQCIHQkr2qOGuInvonv
W4vvdI61nraax5V6SC3E0D2JSO91AiBVhpX4BBafRAh36r7l8LrxAfxBM3CjBmAC
q8fUrWfIWDANBgkqhkiG9w0BAQsFAAOCAQEAgGDMKXofKpDdv5kkID3s5GrKdzaj
jFmb/6kyqd1E6eGXZAewCP1EF5BVvR6lBP2aRXiZ6sJVZktoIfztZnbxBGgbPHfv
R3iXIG6fxkklzR9Y8puPMBFadANE/QV78tIRAlyaqeSNsoxHi7ssQjHTP111B2lf
3KmuTpsruut1UesEJcPReLk/1xTkRx262wAncach5Wp+6GWWduTZYJbsNFyrK1RP
YQ0qYpP9wt2qR+DGaRUBG8i1XLnZS8pkyxtKhVw/a5Fowt+NqCpEBjjJiWJRSGnG
NSgRtSXq11j8O4JONi8EXe7cEtvzUiLR5PL3itsK2svtrZ9jIwQ95wOPaA==
-----END CERTIFICATE-----
CERT

  #private_key = file("www.example.com/privkey.pem")
  private_key = <<PK
-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQC5cyk6x63iBJSW
vtgsOBqIxfO8euPHcRnyWsL9dsvnbNyOnyvcqFWxdiW3sh2cItzYtoN1Zfgj5lWG
OVXbHxP0VaNG9fHVX3+NHP6LFHQz92BzAYQmpqi9zaP/aKJklk6LdPFbVLGhuZfm
34+ijW9YsgLTKR2WTaZJK5QtamVVmP+VsSCla2ifFzjz2FCkMMEc/Y0zUyP+en/m
bL71K+VnpZdlEC1s38EvjRTFKFZTKVw5wpWgCZQq/AZYj9RxR23IIuRcUJ8TQ2py
oc3kIXPWjiIarSgBlA8G9kCsxgzXP2RyLwKrIBIo+qyHweifpPYW28ipdSbPjiyp
AMdpbGLDAgMBAAECggEAUXVfT91z6IqghhKwO8QtC5T/+fN06B8rCYSKj/FFoZL0
0oTiLFuYwImoCadoUDQUE/Efj0rKE2LSgFHg/44IItQXE01m+5WmHmL1ADxsyoLH
z9yDosKj7jNM7RyV8F8Bg0pL1hU+rU4rhhL/MaS0mx4eFYjC4UmcWBmXTdelSVJa
kvXvQLT5y86bqh7tqMjM/kALTWRz5CgNJFk/ONA1yo5RTX9S7SIXimBgAvuGqP8i
MPEhJou7U3DfzXVfvP8byqNdsZs6ZNhG3wXspl61mRyrY+51SOaNLA7Bkji7x4bH
Nw6mJI0IJTAP9oc1Z8fYeMuxT1bfuD7VOupSP0mAMQKBgQDk+KuyQkmPymeP/Wwu
II4DUpleVzxTK9obMQQoCEEElbQ6+jTb+8ixP0bWLvBXg/rX734j7OWfn/bljWLH
XLrSoqQZF1+XMVeY4g4wx9UuTK/D2n791zdOgQivxbIPdWL3a4ap86ar8uyMgJu8
BLXfFBAOc+9myqUkbeO7wt0e6QKBgQDPV04jPtIJoMrggpQDNreGrANKOmsXWxj4
OHW13QNdJ2KGQpoTdoqQ8ZmlxuA8Bf2RjHsnB2kgGVTVQR74zRib4MByhvsdhvVm
F2LNsJoIDfqtv3c+oj13VonRUGuzUeJpwT/snyaL+jQ/ZZcYz0jDgDhIODTcFYj8
DMSD5SHgywKBgHH6MwWuJ44TNBAiF2qyu959jGjAxf+k0ZI9iRMgYLUWjDvbdtqW
cCWDGRDfFraJtSEuTz003GzkJPPJuIUC7OCTI1p2HxhU8ITi6itwHfdJJyk4J4TW
T+qdIqTUpTk6tsPw23zYE3x+lS+viVZDhgEArKl1HpOthh0nMnixnH6ZAoGBAKGn
V+xy1h9bldFk/TFkP8Jn6ki9MzGKfPVKT7vzDORcCJzU4Hu8OFy5gSmW3Mzvfrsz
4/CR/oxgM5vwoc0pWr5thJ3GT5K93iYypX3o6q7M91zvonDa3UFl3x2qrc2pUfVS
DhzWGJ+Z+5JSCnP1aK3EEh18dPoCcELTUYPj6X3xAoGBALAllTb3RCIaqIqk+s3Y
6KDzikgwGM6j9lmOI2MH4XmCVym4Z40YGK5nxulDh2Ihn/n9zm13Z7ul2DJwgQSO
0zBc7/CMOsMEBaNXuKL8Qj4enJXMtub4waQ/ywqHIdc50YaPI5Ax8dD/10h9M6Qc
nUFLNE8pXSnsqb0eOL74f3uQ
-----END PRIVATE KEY-----
PK
}
```

## Argument Reference

The following arguments are supported:

- `domain_id` - (Required) Domain ID.

- `type` - (Required) Certificate type. Valid value is `"PEM"`.

- `certificate` - (Required) Certificate content.

- `private_key` - (Required) Certificate private key.

- `certificate_chain` - (Required) Certificate certificate chain.

## Import

This resource does not support importing.

## Let's Encrypt Certbot

This example demonstrates generatoring a domain certificate with letsencrypt
certbot https://letsencrypt.org/getting-started/

```
$ certbot certonly --manual --preferred-challenges dns --key-type rsa -d [DOMAIN]
```

Use letsencrypt's certbot to generate domain certificates in RSA output mode.
The generator's output corresponds to `okta_domain_certificate` fields in the
following manner.

Okta Field          | Certbot file
--------------------|--------------
`certificate`       | `cert.pem`
`certificate_chain` | `chain.pem`
`private_key`       | `privkey.pem`