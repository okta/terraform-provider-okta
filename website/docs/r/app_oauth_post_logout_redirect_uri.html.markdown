---
layout: 'okta'
page_title: 'Okta: okta_app_oauth_post_logout_redirect_uri'
sidebar_current: 'docs-okta-resource-app-oauth-post-logout-redirect-uri'
description: |-
  Manager app OAuth post logout redirect URI
---

# okta_app_oauth_post_logout_redirect_uri

This resource allows you to manage post logout redirection URI for use in redirect-based flows.

## Example Usage

```hcl
resource "okta_app_oauth" "test" {
  label          = "testAcc_replace_with_uuid"
  type           = "web"
  grant_types    = ["authorization_code"]
  response_types = ["code"]

  // Okta requires at least one redirect URI to create an app
  redirect_uris = ["myapp://callback"]
  
  post_logout_redirect_uris = ["https://www.example.com"]

  // Since Okta forces us to create it with a redirect URI we have to ignore future changes, they will be detected as config drift.
  lifecycle {
    ignore_changes = [redirect_uris]
  }
}

resource "okta_app_oauth_post_logout_redirect_uri" "test" {
  app_id = okta_app_oauth.test.id
  uri    = "https://www.example.com"
}
```

## Argument Reference

- `app_id` - (Required) OAuth application ID.

- `uri` - (Required) Post Logout Redirect URI to append to Okta OIDC application.

## Attributes Reference

- `id` - ID of the resource, equals to `uri`.

## Import

A post logout redirect URI can be imported via the Okta ID.

```
$ terraform import okta_app_oauth_post_logout_redirect_uri.example &#60;app id&#62;/&#60;uri&#62;
```
