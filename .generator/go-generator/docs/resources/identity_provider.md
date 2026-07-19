---
page_title: "okta_identity_provider Resource - terraform-provider-okta"
description: |-
  Manages an Okta Identity Provider resource.
---

# okta_identity_provider (Resource)

Manages an Okta Identity Provider.

## Example Usage

```hcl
resource "okta_identity_provider" "example" {
  inquiry_template_id = "<inquiry-template-id>"

  # Optional fields
  # issuer_mode = "<issuer_mode>"
  # name = "Example Name"
  # action = "<action>"
  # include = "<include>"
  # exclude = "<exclude>"
}
```

## Schema

### Required

- `inquiry_template_id` (String) The ID of the inquiry template from your Persona dashboard. The inquiry template always starts with `itmpl`. Applies to the `IDV_PERSONA` IdP type.

### Optional

- `issuer_mode` (String) Indicates whether Okta uses the original Okta org domain URL or a custom domain URL in the request to the social IdP
- `name` (String) Unique name for the IdP
- `action` (String) Specifies the account linking action for an IdP user
- `include` (List) Specifies the allowlist of Group identifiers to match against. Group memberships are restricted to type `OKTA_GROUP`.
- `exclude` (List) Specifies the blocklist of user identifiers to exclude from account linking
- `exclude_admins` (Boolean) Specifies whether admin users should be excluded from account linking
- `max_clock_skew` (Number) Maximum allowable clock skew when processing messages from the IdP
- `action` (String) Specifies the user provisioning action during authentication when an IdP user isn't linked to an existing Okta user. * To successfully provision a new Okta user, you must enable just-in-time (JIT) ...
- `action` (String) Specifies the action during authentication when an IdP user is linked to a previously deprovisioned Okta user
- `action` (String) Specifies the action during authentication when an IdP user is linked to a previously suspended Okta user
- `action` (String) Provisioning action for the IdP user's group memberships  | Enum     | Description                                                                                                                   ...
- `assignments` (List) List of `OKTA_GROUP` group identifiers to add an IdP user as a member with the `ASSIGN` action
- `filter` (List) Allowlist of `OKTA_GROUP` group identifiers for the `APPEND` or `SYNC` provisioning action
- `source_attribute_name` (String) IdP user profile attribute name (case-insensitive) for an array value that contains group memberships
- `profile_master` (Boolean) Determines if the IdP should act as a source of truth for user profile attributes
- `filter` (String) Optional [regular expression pattern](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Regular_expressions) used to filter untrusted IdP usernames. * As a best security practice, y...
- `match_attribute` (String) Okta user profile attribute for matching a transformed IdP username. Only for matchType `CUSTOM_ATTRIBUTE`. The `matchAttribute` must be a valid Okta user profile attribute of one of the following ...
- `match_type` (String) Determines the Okta user profile attribute match conditions for account linking and authentication of the transformed IdP username
- `template` (String) Template
- `trust_claims` (Boolean) Indicates whether to trust authentication claims from the IdP
- `aal_value` (String) The [authentication assurance level](https://developers.login.gov/oidc/#aal-values) (AAL) value for the Login.gov IdP. See [Add a Login.gov IdP](https://developer.okta.com/docs/guides/add-logingov-...
- `additional_amr` (List) The additional Assurance Methods References (AMR) values for Smart Card IdPs. Applies to `X509` IdP type.
- `ial_value` (String) The [type of identity verification](https://developers.login.gov/oidc/#ial-values) (IAL) value for the Login.gov IdP. See [Add a Login.gov IdP](https://developer.okta.com/docs/guides/add-logingov-i...
- `privacy_policy` (String) A URL that links to the privacy policy for the IDV vendor
- `terms_of_use` (String) A URL that links to the terms of use for the IDV vendor
- `vendor_display_name` (String) The display name of the IDV vendor
- `protocol` (String) IdP-specific protocol settings for endpoints, bindings, and algorithms used to connect with the IdP and validate messages
- `status` (String) Status
- `type` (String) The IdP object's `type` property identifies the social or enterprise IdP used for authentication. Each IdP uses a specific protocol, therefore the `protocol` object must correspond with the IdP `ty...

### Read-Only

- `id` (String) The unique identifier for the resource.
- `created` (String) Timestamp when the object was created
- `last_updated` (String) Timestamp when the object was last updated
