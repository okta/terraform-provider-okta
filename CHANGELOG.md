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
