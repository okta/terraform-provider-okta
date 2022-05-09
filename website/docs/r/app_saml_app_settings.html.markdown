---
layout: 'okta'
page_title: 'Okta: okta_app_saml_app_settings'
sidebar_current: 'docs-okta-resource-app-saml-app-settings'
description: |-
  Manages app settings of the SAML application.
---

# okta_app_saml_app_settings

This resource allows you to manage app settings of the SAML Application . It's basically the same as
`app_settings_json` field in `okta_app_saml` resource and can be used in cases where settings require to be managed separately.

## Example Usage

```hcl
resource "okta_app_saml" "test" {
  preconfigured_app = "amazon_aws"
  label             = "Amazon AWS"
  status            = "ACTIVE"
}

resource "okta_app_saml_app_settings" "test" {
  app_id   = okta_app_saml.test.id
  settings = jsonencode(
  {
    "appFilter" : "okta",
    "awsEnvironmentType" : "aws.amazon",
    "groupFilter" : "aws_(?{{accountid}}\\\\d+)_(?{{role}}[a-zA-Z0-9+=,.@\\\\-_]+)",
    "joinAllRoles" : false,
    "loginURL" : "https://console.aws.amazon.com/ec2/home",
    "roleValuePattern" : "arn:aws:iam::$${accountid}:saml-provider/OKTA,arn:aws:iam::$${accountid}:role/$${role}",
    "sessionDuration" : 3200,
    "useGroupMapping" : false
  }
  )
}
```

## Argument Reference

The following arguments are supported:

- `app_id` - (Required) ID of the application.

- `settings` - (Required) Application settings in JSON format.

## Import

A settings for the SAML App can be imported via the Okta ID.

```
$ terraform import okta_app_saml_app_settings.example &#60;app id&#62;
```
