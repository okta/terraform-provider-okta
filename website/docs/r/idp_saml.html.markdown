---
layout: 'okta'
page_title: 'Okta: okta_idp_saml'
sidebar_current: 'docs-okta-resource-idp-saml'
description: |-
  Creates a SAML Identity Provider.
---

# okta_idp_saml

Creates a SAML Identity Provider.

This resource allows you to create and configure a SAML Identity Provider.

## Example Usage

```hcl
resource "okta_idp_saml" "example" {
  name                     = "testAcc_replace_with_uuid"
  acs_type                 = "INSTANCE"
  sso_url                  = "https://idp.example.com"
  sso_destination          = "https://idp.example.com"
  sso_binding              = "HTTP-POST"
  username_template        = "idpuser.email"
  kid                      = "${okta_idp_saml_key.test.id}"
  issuer                   = "https://idp.example.com"
  request_signature_scope  = "REQUEST"
  response_signature_scope = "ANY"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The Application's display name.

- `kid` - (Required) The ID of the signing key.

- `sso_url` - (Required) URL of binding-specific endpoint to send an AuthnRequest message to IdP.

- `issuer` - (Required) URI that identifies the issuer.

- `acs_type` - (Optional) The type of ACS. It can be `"INSTANCE"` or `"ORG"`.

- `sso_binding` - (Optional) The method of making an SSO request. It can be set to `"HTTP-POST"` or `"HTTP-REDIRECT"`.

- `sso_destination` - (Optional) URI reference indicating the address to which the AuthnRequest message is sent.

- `name_format` - (Optional) The name identifier format to use. By default `"urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified"`.

- `subject_format` - (Optional) The name format. By default `"urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified"`.

- `subject_filter` - (Optional) Optional regular expression pattern used to filter untrusted IdP usernames.

- `issuer_mode` - (Optional) Indicates whether Okta uses the original Okta org domain URL, or a custom domain URL. It can be `"ORG_URL"` or `"CUSTOM_URL"`.

- `max_clock_skew` - (Optional) Maximum allowable clock-skew when processing messages from the IdP.

- `status` - (Optional) Status of the IdP.

- `account_link_action` - (Optional) Specifies the account linking action for an IdP user.

- `account_link_group_include` - (Optional) Group memberships to determine link candidates.

- `provisioning_action` - (Optional) Provisioning action for an IdP user during authentication.

- `deprovisioned_action` - (Optional) Action for a previously deprovisioned IdP user during authentication. Can be `"NONE"` or `"REACTIVATE"`.

- `suspended_action` - (Optional) Action for a previously suspended IdP user during authentication. Can be set to `"NONE"` or `"UNSUSPEND"`

- `groups_action` - (Optional) Provisioning action for IdP user's group memberships. It can be `"NONE"`, `"SYNC"`, `"APPEND"`, or `"ASSIGN"`.

- `groups_attribute` - (Optional) IdP user profile attribute name (case-insensitive) for an array value that contains group memberships.

- `groups_assignment` - (Optional) List of Okta Group IDs to add an IdP user as a member with the `"ASSIGN"` `groups_action`.

- `groups_filter` - (Optional) Whitelist of Okta Group identifiers that are allowed for the `"APPEND"` or `"SYNC"` `groups_action`.

- `username_template` - (Optional) Okta EL Expression to generate or transform a unique username for the IdP user.

- `subject_match_type` - (Optional) Determines the Okta user profile attribute match conditions for account linking and authentication of the transformed IdP username. By default, it is set to `"USERNAME"`. It can be set to `"USERNAME"`, `"EMAIL"`, `"USERNAME_OR_EMAIL"` or `"CUSTOM_ATTRIBUTE"`.

- `subject_match_attribute` - (Optional) Okta user profile attribute for matching transformed IdP username. Only for matchType `"CUSTOM_ATTRIBUTE"`.

- `profile_master` - (Optional) Determines if the IdP should act as a source of truth for user profile attributes.

- `request_signature_algorithm` - (Optional) The XML digital signature algorithm used when signing an AuthnRequest message.

- `request_signature_scope` - (Optional) Specifies whether to digitally sign an AuthnRequest messages to the IdP. It can be `"REQUEST"` or `"NONE"`.

- `response_signature_algorithm` - (Optional) The minimum XML digital signature algorithm allowed when verifying a SAMLResponse message or Assertion element.

- `response_signature_scope` - (Optional) Specifies whether to verify a SAMLResponse message or Assertion element XML digital signature. It can be `"RESPONSE"`, `"ASSERTION"`, or `"ANY"`.

## Attributes Reference

- `id` - ID of the IdP.

- `type` - Type of the IdP.

- `audience` - The audience restriction for the IdP.

## Import

An SAML IdP can be imported via the Okta ID.

```
$ terraform import okta_idp_saml.example <idp id>
```
