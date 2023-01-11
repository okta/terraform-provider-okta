---
layout: "okta"
page_title: "Okta: okta_app_ws_federation"
sidebar_current: "docs-okta-resource-app-ws-federation"
description: |-
  Manages a WS-Federation Application.
---

# okta_app_ws_federation

This resource allows you to create and configure a WS-Federation Application.

-> If you receive the error `You do not have permission to access the feature
you are requesting` [contact support](mailto:dev-inquiries@okta.com) and
request feature flag `ADVANCED_SSO` be applied to your org.

## Example Usage

```hcl
resource "okta_app_ws_federation" "exampleWsFedApp" {
  label    = "exampleWsFedApp"
  site_url = "https://signin.example.com/saml"
  reply_url = "https://example.com"
  reply_override = false
  name_id_format = "urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified"
  audience_restriction = "https://signin.example.com"
  authn_context_class_ref = "urn:oasis:names:tc:SAML:2.0:ac:classes:PasswordProtectedTransport"
  group_filter = "app1.*"
  group_name = "username"
  group_value_format = "dn"
  username_attribute = "username"
  attribute_statements = "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/givenname|bob|,http://schemas.xmlsoap.org/ws/2005/05/identity/claims/surname|hope|"
  visibility = false
  status = "ACTIVE"
}
```

## Argument Reference

The following arguments are supported:

- `attribute_statements` - (Optional) Defines custom SAML attribute statements`

- `audience_restriction` - (Optional) The assertion containing a bearer subject confirmation MUST contain an Audience Restriction including the service provider's unique identifier as an Audience.

- `authn_context_class_ref` - (Optional) Specifies the Authentication Context for the issued SAML Assertion.

- `auto_submit_toolbar` - (Optional) Display auto submit toolbar

- `group_filter` - (Optional) An expression that will be used to filter groups. If the Okta group name matches the expression, the group name will be included in the SAML Assertion Attribute Statement.

- `group_name` - (Optional) Specifies the SAML attribute name for a user's group memberships.

- `group_value_format` - (Optional) Specifies the SAML assertion attribute value for filtered groups.

- `hide_ios` -  (Optional) Do not display application icon on mobile app

- `hide_web` - (Optional) Do not display application icon to users  

- `name_id_format` - (Optional) Name ID Format.

- `realm` - (Optional) The trust realm for the Web Application.

- `reply_url` - (Optional) The ReplyTo URL to which responses are directed.

- `reply_override` - (Optional) Enable web application to override ReplyTo URL with reply param.

- `site_url` - (Optional) Launch URL for the Web Application.
 
- `status` - (Optional) Activation status of the application

- `username_attribute` - (Optional) Specifies additional username attribute statements to include in the SAML Assertion.

- `visibility` - (Optional) Application icon visibility to users.


## Attributes Reference

- `id` - id of application.

- `name` - Name assigned to the application by Okta.

- `site_url` - Launch URL for the Web Application.

- `realm` - The trust realm for the Web Application.

- `reply_url` - The ReplyTo URL to which responses are directed.

- `reply_override` - Enable web application to override ReplyTo URL with reply param.

- `name_id_format` - Name ID Format.

- `audience_restriction` - The assertion containing a bearer subject confirmation MUST contain an Audience Restriction including the service provider's unique identifier as an Audience.

- `authn_context_class_ref` - Specifies the Authentication Context for the issued SAML Assertion.

- `group_filter` - An expression that will be used to filter groups. If the Okta group name matches the expression, the group name will be included in the SAML Assertion Attribute Statement

- `group_name` - Specifies the SAML attribute name for a user's group memberships.

- `group_value_format` - Specifies the SAML assertion attribute value for filtered groups.

- `username_attribute` - Specifies additional username attribute statements to include in the SAML Assertion.

- `attribute_statements` - Defines custom SAML attribute statements.

- `visibility` - Application icon visibility to users.

## Timeouts

The `timeouts` block allows you to specify custom [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions: 

- `create` - Create timeout if syncing users/groups (default 1 hour).

- `update` - Update timeout if syncing users/groups (default 1 hour).

- `read` - Read timeout if syncing users/groups (default 1 hour).

## Import

A WS-Federation App can be imported via the Okta ID.

```
$ terraform import okta_app_ws_federation.example &#60;app id&#62;
```

