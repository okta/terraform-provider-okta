---
layout: 'okta'
page_title: 'Okta: okta_app_oauth_redirect_uri'
sidebar_current: 'docs-okta-resource-app-oauth-redirect-uri'
description: |-
  Manager app OAuth redirect URI
---

# okta_app_oauth_redirect_uri

This resource allows you to manage redirection URI for use in redirect-based flows.

~> `okta_app_oauth_redirect_uri` has been marked deprecated and will be removed
in the v5 release of the provider. Operators should manage the redirect URIs for
an oauth app directly on that resource.

## Example Usage

```hcl
resource "okta_app_oauth" "test" {
  label          = "testAcc_replace_with_uuid"
  type           = "web"
  grant_types    = ["authorization_code"]
  response_types = ["code"]

  // Okta requires at least one redirect URI to create an app
  redirect_uris = ["myapp://callback"]

  // Since Okta forces us to create it with a redirect URI we have to ignore future changes, they will be detected as config drift.
  lifecycle {
    ignore_changes = [redirect_uris]
  }
}

resource "okta_app_oauth_redirect_uri" "test" {
  app_id = okta_app_oauth.test.id
  uri    = "http://google.com"
}
```

## Argument Reference

- `app_id` - (Required) OAuth application ID. Note: `app_id` can not be changed once set.

- `uri` - (Required) Redirect URI to append to Okta OIDC application.

## Attributes Reference

- `id` - ID of the resource, equals to `uri`.

## Import

A redirect URI can be imported via the Okta ID.

```
$ terraform import okta_app_oauth_redirect_uri.example &#60;app id&#62;/&#60;uri&#62;
```
