<!--- Please keep this note for the community --->

### Community Note

- Please vote on this issue by adding a 👍 [reaction](https://blog.github.com/2016-03-10-add-reactions-to-pull-requests-issues-and-comments/) to the original issue to help the community and maintainers prioritize this request
- Please do not leave "+1" or other comments that do not add relevant new information or questions, they generate extra noise for issue followers and do not help prioritize the request
- If you are interested in working on this issue or have submitted a pull request, please leave a comment

<!--- Thank you for keeping this note for the community --->

### Description

This request is to add 2 new resources to support WS Federation Apps in Okta. We confirmed that the Okta SDK for Golang supports this. 

### New or Affected Resource(s)

- okta_app_ws_federation
- data_source_okta_app_ws_federation

### Potential Terraform Configuration

<!--- Information about code formatting: https://help.github.com/articles/basic-writing-and-formatting-syntax/#quoting-code --->

```hcl
# Resource - Okta WS Federated App
resource "okta_app_ws_federation" "example" {
    label    = "example"
	site_url = "https://signin.example.com/saml"
	realm = "example"
	reply_to_url = "https://example.com"
	allow_reply_to_override = false
    name_id_format = "uid"
    audience_restriction = "https://signin.example.com"
    assert_authentication_context = "Kerberos"
    group_filter = "app1.*"
    group_attribute_name = "username"
    group_attribute_value = "dn"
    username_attribute = "username"
    custom_attribute_statements = "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/givenname|${user.firstName}|,http://schemas.xmlsoap.org/ws/2005/05/identity/claims/surname|${user.lastName}|"
    visibility = true
    signature_algorithm = "SHA256"     
    digest_algorithm = "SHA256"        
}

# Data Source - Okta WS Federated App
data "okta_app_ws_federation" "example" {
  id = "0ob3otzg0CHSgPcjZ0z9"
}
```

### References

Okta Golang SDK link for WS Federated Apps below:

- [wsFederationApplication.go](https://github.com/okta/okta-sdk-golang/blob/master/okta/wsFederationApplication.go)
- [wsFederationApplicationSettings.go](https://github.com/okta/okta-sdk-golang/blob/master/okta/wsFederationApplicationSettings.go)
- [wsFederationApplicationSettingsApplication.go](https://github.com/okta/okta-sdk-golang/blob/master/okta/wsFederationApplicationSettingsApplication.go)

- #0000
