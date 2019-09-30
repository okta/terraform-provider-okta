For Release v3.0.0:

* Updated provider protocol version to v5 to support Terraform v0.12.0

For Release v3.0.1

* Add some missing okta_idp_saml settings
* Add registration inline hook type

For Release v3.0.2

* Use backoff/retries functionality for XML API calls

For Release v3.0.3

* Add okta_idp_saml data source
* Support import user by email

For Release v3.0.4

* Change type of custom_profile_attributes from map to JSON string to support all types
* Add okta_app_saml data source
* Add okta_app_metadata_saml data source
* Add okta_idp_metadata_saml data source
* Fix group filter bug, filter_type and filter_value were not being sync'd

For Release v3.0.5

* Fix bug introduced in v3.0.4. User data source was not updated to the new caustom_profile_attribute type
* Added test to cover this scenario, tests were passiing

For Release v3.0.6

* Allow client_id to be set on OIDC application, while also maintaining the computed version. With some auth methods, such as basic auth, this is possible.

For Release v3.0.7

* Add group_assignments for SAML and social IdPs

For Release v3.0.8

* Add issuer_mode to social IdP. Our test org does not have a custom domain setup, thus it was working there but not in other orgs. Hard to test both scenarios in one org.

For Release v3.0.9

* Add okta_template_email resource for defining Custom Email Templates
* Add okta_group_roles resource for defining the admin roles tied to a group

For Release v3.0.10

* Add okta_network_zone resource

For Release v3.0.11

* Fix ocassional panic when creating a user schema see https://github.com/articulate/terraform-provider-okta/issues/144
* Users in LOCKED_OUT state are unlocked when config is ACTIVE https://github.com/articulate/terraform-provider-okta/issues/225

For Release v3.0.12

* Ensure schema does not panic after retry :smh:

For Release v3.0.13

* Add okta_user_base_schema resource for managing base schema properties
* Add missing attribute, match_type and match_attribute, on social idp resource

For Release v3.0.14

* Fix logic around including/excluding networks on policy rules

For Release v3.0.15

* Update Okta SDK
* Filter out GROUP based admin roles when processing user `admin_roles` attribute

For Release v3.0.16

* Fix issues around `okta_policy_rule_idp_discovery`
    * `app_include` and `app_exlcude` were missing required properties
    * `user_identifier_type` was being added even when not defined, causing API errors
* Fix integer array type

For Release v3.0.17

* Add okta_app_user_schema resource
* Add okta_app_user_base_schema resource
* Add `required` field to base schema
* Add `okta_app_user` resource
* Add `okta_app_group` resource

For Release v3.0.18

* Support SHA-1 signing algorithm on IdPs
* Fix bug where audience is reset on IdP update because it is omitted from the payload

For Release v3.0.19

* Fix diff issues around `okta_policy_rule_idp_discovery`
* Allow `provisioning_action` for IdPs to be set to `DISABLED`

For Release v3.0.20

* Fix `okta_auth_server_claim`, `group_filter_type` could not be set to `STARTS_WITH` due to a typo

For Release v3.0.21

* Expose scope property on `okta_user_schema`
* Allow setting of OAuth application visibility settings

For Release v3.0.22

* Send `profileMaster` along with IdP, so the config is recognized by Okta API
* Fix bug in SDK related to retries and the request body being empty on subsequent requests.
* Various updates related to Hashicorp's review process that aren't necessarily functionality related, see https://github.com/articulate/terraform-provider-okta/pull/271

For Release v3.0.23

* Add `external_name` property to the `okta_app_user_schema` and `okta_user_schema`

For Release v3.0.24

* Support `profile` on `okta_oauth_app` resource

For Release v3.0.25

* Support setting an auth server scope as the default
* Fix bug with supporting `profile` on `okta_oauth_app` resource
* Support `profile` and `priority` on `okta_app_group_assignment`
* Support `profile` on `okta_app_user`

For Release v3.0.26

* Support array enums in `okta_user_schema` and `okta_app_user_schema` as `array_enum` and `array_one_of`

For Release v3.0.27

* Update refresh token window validation to account for new upper limit of 5 years 

For Release v3.0.28

* Remove resource from state on 404.
