---
page_title: "Data Source: okta_app_features"
description: |-
  Get an application of any kind from Okta.
---
# Data Source: okta_app_features

Retrieves a Feature object for an app.

## Example Usage

### Basic Usage

```terraform
data "okta_app_features" "example" {
  app_id = "0oarblaf7hWdLawNg1d7"
  name   = "INBOUND_PROVISIONING"
}
```

### User Provisioning Feature

```terraform
data "okta_app_features" "user_provisioning" {
  app_id = okta_app_saml.example.id
  name   = "USER_PROVISIONING"
}

output "provisioning_status" {
  value = data.okta_app_features.user_provisioning.status
}

output "password_sync_enabled" {
  value = data.okta_app_features.user_provisioning.capabilities.update.password.status
}
```

### Inbound Provisioning Feature

```hcl
data "okta_app_features" "inbound_provisioning" {
  app_id = okta_app_oauth.example.id
  name   = "INBOUND_PROVISIONING"
}

output "import_schedule_status" {
  value = data.okta_app_features.inbound_provisioning.capabilities.import_settings.schedule.status
}

output "auto_activate_users" {
  value = data.okta_app_features.inbound_provisioning.capabilities.import_rules.user_create_and_match.auto_activate_new_users
}
```

## Argument Reference

The following arguments are supported:

* `app_id` - (Required) The ID of the application to retrieve features for.
* `name` - (Required) The name of the feature to retrieve. Valid values include:
    * `USER_PROVISIONING` - User profiles are pushed from Okta to the third-party app.
    * `INBOUND_PROVISIONING` - User profiles are imported from the third-party app into Okta.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the app feature.
* `description` - Description of the feature.
* `status` - The current status of the feature. Valid values are `ENABLED` or `DISABLED`.
* `capabilities` - A block containing the feature capabilities configuration. See [Capabilities](#capabilities) below.

### Capabilities

The `capabilities` block contains different configuration blocks depending on the feature type:

#### Create Capabilities

* `create` - Block for create lifecycle settings:
    * `lifecycle_create` - Block for create lifecycle configuration:
        * `status` - (String) Status of the create lifecycle setting. Valid values are `ENABLED` or `DISABLED`.

#### Update Capabilities

* `update` - Block for update settings:
    * `lifecycle_deactivate` - Block for deactivation lifecycle configuration:
        * `status` - (String) Status of the deactivate lifecycle setting. Valid values are `ENABLED` or `DISABLED`.
    * `password` - Block for password synchronization settings:
        * `change` - (String) Determines password change behavior. Valid values are `CHANGE` or `KEEP_EXISTING`.
        * `seed` - (String) Determines password source. Valid values are `OKTA` or `RANDOM`.
        * `status` - (String) Status of password sync. Valid values are `ENABLED` or `DISABLED`.
    * `profile` - Block for profile update settings:
        * `status` - (String) Status of profile updates. Valid values are `ENABLED` or `DISABLED`.

#### Import Rules (Inbound Provisioning)

* `import_rules` - Block for import rules configuration:
    * `user_create_and_match` - Block for user matching and creation rules:
        * `exact_match_criteria` - (String) Attribute used for exact matching (e.g., `USERNAME`, `EMAIL`).
        * `allow_partial_match` - (Boolean) Whether to allow partial matching based on first and last names.
        * `auto_activate_new_users` - (Boolean) Whether imported new users are automatically activated.
        * `autoconfirm_exact_match` - (Boolean) Whether exact-matched users are automatically confirmed.
        * `autoconfirm_new_users` - (Boolean) Whether imported new users are automatically confirmed.
        * `autoconfirm_partial_match` - (Boolean) Whether partially matched users are automatically confirmed.

#### Import Settings (Inbound Provisioning)

* `import_settings` - Block for import settings configuration:
    * `username` - Block for username configuration:
        * `username_format` - (String) Format for usernames (e.g., `EMAIL`, `CUSTOM`).
        * `username_expression` - (String) Okta Expression Language statement for custom username format.
    * `schedule` - Block for import schedule configuration:
        * `status` - (String) Status of the import schedule. Valid values are `ENABLED` or `DISABLED`.
        * `full_import` - Block for full import schedule:
            * `expression` - (String) UNIX cron expression for full import schedule.
            * `timezone` - (String) IANA timezone name for the schedule.
        * `incremental_import` - Block for incremental import schedule:
            * `expression` - (String) UNIX cron expression for incremental import schedule.
            * `timezone` - (String) IANA timezone name for the schedule.
