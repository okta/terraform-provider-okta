---
page_title: "Resource: okta_app_features"
description: |-
  Manages Okta application features. This resource allows you to configure provisioning capabilities for applications, including user provisioning (outbound) and inbound provisioning settings.
---
# Resource: okta_app_features

Manages Okta application features. This resource allows you to configure provisioning capabilities for applications, including user provisioning (outbound) and inbound provisioning settings.

~> **NOTE:** This resource cannot be deleted via Terraform. Application features are managed by Okta and can only be updated or read.

~> **NOTE:** This resource is only supported with a limited subset of OIN applications, see the [api docs](https://developer.okta.com/docs/api/openapi/okta-management/management/tag/ApplicationFeatures/) for more details.

## Example Usage

### User Provisioning Feature

```terraform
resource "okta_app_features" "user_provisioning" {
  app_id = okta_app_saml.example.id
  name   = "USER_PROVISIONING"
  status = "ENABLED"

  capabilities {
    create {
      lifecycle_create {
        status = "ENABLED"
      }
    }

    update {
      lifecycle_deactivate {
        status = "ENABLED"
      }

      password {
        change = "CHANGE"
        seed   = "OKTA"
        status = "ENABLED"
      }

      profile {
        status = "ENABLED"
      }
    }
  }
}
```

### Inbound Provisioning Feature

```terraform
resource "okta_app_features" "inbound_provisioning" {
  app_id = okta_app_oauth.example.id
  name   = "INBOUND_PROVISIONING"
  status = "ENABLED"

  capabilities {
    import_rules {
      user_create_and_match {
        exact_match_criteria        = "USERNAME"
        allow_partial_match         = true
        auto_activate_new_users     = false
        autoconfirm_exact_match     = false
        autoconfirm_new_users       = false
        autoconfirm_partial_match   = false
      }
    }

    import_settings {
      username {
        username_format     = "EMAIL"
        username_expression = ""
      }

      schedule {
        status = "DISABLED"
        
        full_import {
          expression = "0 0 * * *"
          timezone   = "America/New_York"
        }
        
        incremental_import {
          expression = "0 */6 * * *"
          timezone   = "America/New_York"
        }
      }
    }
  }
}
```

### Complete User Provisioning Configuration

```terraform
resource "okta_app_features" "complete_provisioning" {
  app_id = okta_app_saml.example.id
  name   = "USER_PROVISIONING"
  status = "ENABLED"

  capabilities {
    create {
      lifecycle_create {
        status = "ENABLED"
      }
    }

    update {
      lifecycle_deactivate {
        status = "ENABLED"
      }

      password {
        change = "CHANGE"
        seed   = "RANDOM"
        status = "ENABLED"
      }

      profile {
        status = "ENABLED"
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `app_id` - (Required) The ID of the application to configure features for.
* `name` - (Required) The name of the feature to configure. Valid values:
    * `USER_PROVISIONING` - User profiles are pushed from Okta to the third-party app.
    * `INBOUND_PROVISIONING` - User profiles are imported from the third-party app into Okta.
* `status` - (Optional) The status of the feature. Valid values are `ENABLED` or `DISABLED`.
* `description` - (Optional) Description of the feature.
* `capabilities` - (Optional) Configuration block for feature capabilities. See [Capabilities](#capabilities) below.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the app feature in the format `{app_id}/{feature_name}`.

### Capabilities

The `capabilities` block supports different configuration blocks depending on the feature type:

#### Create Capabilities (User Provisioning)

* `create` - (Optional) Block for create lifecycle settings:
    * `lifecycle_create` - (Optional) Block for create lifecycle configuration:
        * `status` - (Optional) Status of the create lifecycle setting. Valid values are `ENABLED` or `DISABLED`.

#### Update Capabilities (User Provisioning)

* `update` - (Optional) Block for update settings:
    * `lifecycle_deactivate` - (Optional) Block for deactivation lifecycle configuration:
        * `status` - (Optional) Status of the deactivate lifecycle setting. Valid values are `ENABLED` or `DISABLED`.
    * `password` - (Optional) Block for password synchronization settings:
        * `change` - (Optional) Determines password change behavior. Valid values are `CHANGE` or `KEEP_EXISTING`.
        * `seed` - (Optional) Determines password source. Valid values are `OKTA` or `RANDOM`.
        * `status` - (Optional) Status of password sync. Valid values are `ENABLED` or `DISABLED`.
    * `profile` - (Optional) Block for profile update settings:
        * `status` - (Optional) Status of profile updates. Valid values are `ENABLED` or `DISABLED`.

#### Import Rules (Inbound Provisioning)

* `import_rules` - (Optional) Block for import rules configuration:
    * `user_create_and_match` - (Optional) Block for user matching and creation rules:
        * `exact_match_criteria` - (Optional) Attribute used for exact matching (e.g., `USERNAME`, `EMAIL`).
        * `allow_partial_match` - (Optional) Whether to allow partial matching based on first and last names.
        * `auto_activate_new_users` - (Optional) Whether imported new users are automatically activated.
        * `autoconfirm_exact_match` - (Optional) Whether exact-matched users are automatically confirmed.
        * `autoconfirm_new_users` - (Optional) Whether imported new users are automatically confirmed.
        * `autoconfirm_partial_match` - (Optional) Whether partially matched users are automatically confirmed.

#### Import Settings (Inbound Provisioning)

* `import_settings` - (Optional) Block for import settings configuration:
    * `username` - (Optional) Block for username configuration:
        * `username_format` - (Optional) Format for usernames (e.g., `EMAIL`, `CUSTOM`).
        * `username_expression` - (Optional) Okta Expression Language statement for custom username format.
    * `schedule` - (Optional) Block for import schedule configuration:
        * `status` - (Optional) Status of the import schedule. Valid values are `ENABLED` or `DISABLED`.
        * `full_import` - (Optional) Block for full import schedule:
            * `expression` - (Optional) UNIX cron expression for full import schedule.
            * `timezone` - (Optional) IANA timezone name for the schedule.
        * `incremental_import` - (Optional) Block for incremental import schedule:
            * `expression` - (Optional) UNIX cron expression for incremental import schedule.
            * `timezone` - (Optional) IANA timezone name for the schedule.

## Import

App features can be imported using the format `{app_id}/{feature_name}`:

```bash
terraform import okta_app_features.example 0oarblaf7hWdLawNg1d7/USER_PROVISIONING
```

```bash
terraform import okta_app_features.inbound 0oarblaf7hWdLawNg1d7/INBOUND_PROVISIONING
```

## Behavior Notes

### Deletion Behavior

This resource cannot be deleted via Terraform. Running `terraform destroy` will show a warning but will not actually delete the feature configuration. Application features are managed by Okta and persist with the application.

### Cron Expression Examples

For import schedules, use standard UNIX cron format:

- `0 0 * * *` - Daily at midnight
- `0 */6 * * *` - Every 6 hours
- `0 0 * * 0` - Weekly on Sunday at midnight
- `0 2 1 * *` - Monthly on the 1st at 2 AM
