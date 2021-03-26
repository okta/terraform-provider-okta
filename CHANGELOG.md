# Changelog

## 3.11.0 (March 26, 2021)

ENHANCEMENTS:

* Add new `okta_app_oauth_api_scope` resource [#356](https://github.com/okta/terraform-provider-okta/pull/356). Thanks, [@mariussturm](https://github.com/mariussturm)!
* Remove `ForceNew` in case policy name changes to avoid policy resources recreation [#362](https://github.com/okta/terraform-provider-okta/pull/362). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Add hotp factor to the `okta_policy_mfa` resource [#363](https://github.com/okta/terraform-provider-okta/pull/363). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Remove unnecessary validations from the `okta_app_oauth` resource [#372](https://github.com/okta/terraform-provider-okta/pull/372). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Add `links` field to `okta_app`, `okta_app_saml` and `okta_app_oauth` data sources [#374](https://github.com/okta/terraform-provider-okta/pull/374). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Add new `okta_auth_server_default` resource [#375](https://github.com/okta/terraform-provider-okta/pull/375). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Add new `okta_policy_mfa_default` and `okta_policy_password_default` resources [#378](https://github.com/okta/terraform-provider-okta/pull/378). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Add `remove_assigned_users` field to the `okta_group_rule` resource [#388](https://github.com/okta/terraform-provider-okta/pull/388). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Add new `auth_server_claim_default` resource [#392](https://github.com/okta/terraform-provider-okta/pull/392). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Add `groups` and `users` fields to the `okta_app`, `okta_app_oauth` and `okta_app_saml` data sources [#395](https://github.com/okta/terraform-provider-okta/pull/395). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Add `id` field to the `okta_group` data source [#395](https://github.com/okta/terraform-provider-okta/pull/395). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Add new `auth_server_claim_default` resource [#392](https://github.com/okta/terraform-provider-okta/pull/392). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Add new `okta_groups` data source [#103](https://github.com/okta/terraform-provider-okta/pull/103). Thanks, [@bendrucker](https://github.com/bendrucker) and [@me](https://github.com/bogdanprodan-okta)!
* Several minor bug fixes and enhancements.

BUGS:

* Add group existence check to `okta_group_membership` resource [#380](https://github.com/okta/terraform-provider-okta/pull/380). Thanks, [@ymylei](https://github.com/ymylei)!
* Fix group assignment priority in the `okta_app_group_assignment` resource [#381](https://github.com/okta/terraform-provider-okta/pull/381). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Fixed status change in the `okta_auth_server_policy_rule` resource [](https://github.com/okta/terraform-provider-okta/pull/386). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Add operation retry to the `okta_group_role` resource [#390](https://github.com/okta/terraform-provider-okta/pull/390). Thanks, [@me](https://github.com/bogdanprodan-okta)!

## 3.10.1 (February 26, 2021)

ENHANCEMENTS:

* Add `retain_assignment` field to `okta_app_user` and `okta_app_group_assignment` resource [#330](https://github.com/oktadeveloper/terraform-provider-okta/pull/330). Thanks, [@Omicron7](https://github.com/Omicron7)!
* Add `target_app_list` field to the `okta_group_role` resource [#349](https://github.com/oktadeveloper/terraform-provider-okta/pull/349). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Add support for `OVERRIDE` value in `master` field and new `master_override_priority` field to the `okta_user_schema` resource [#351](https://github.com/oktadeveloper/terraform-provider-okta/pull/351). Thanks, [@me](https://github.com/bogdanprodan-okta)!

BUGS:

* Added wait to `okta_group_membership` resource [#335](https://github.com/oktadeveloper/terraform-provider-okta/pull/335). Thanks, [@ymylei](https://github.com/ymylei)!
* Fix set of `subject_match_attribute` value for `okta_idp_oidc` resource [#344](https://github.com/oktadeveloper/terraform-provider-okta/pull/344). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Fix resource validation [#348](https://github.com/oktadeveloper/terraform-provider-okta/pull/348). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Fix setup of empty `login_scopes` for `okta_app_oauth` resource [#352](https://github.com/oktadeveloper/terraform-provider-okta/pull/352). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Fix `okta_group_role` when removing all the items from `target_group_list` [#341](https://github.com/oktadeveloper/terraform-provider-okta/pull/341). Thanks, [@me](https://github.com/bogdanprodan-okta)!

## 3.10.0 (February 19, 2021)

ENHANCEMENTS:

* Add new `okta_auth_server_scopes` datasource [#336](https://github.com/oktadeveloper/terraform-provider-okta/pull/336). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Add new `okta_idp_social` datasource [#337](https://github.com/oktadeveloper/terraform-provider-okta/pull/337). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Several minor bug fixes and enhancements.

BUGS:

* Fix preconfigured `okta_app_swa` creation in case it has more that one sign-on modes [#328](https://github.com/oktadeveloper/terraform-provider-okta/pull/328). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Add force recreate in case `okta_app_user_schema` changes the `scope` value since it's a read-only attribute [#331](https://github.com/oktadeveloper/terraform-provider-okta/pull/331). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Fix false positive output when runnning `terraform plan`for the `okta_profile_mapping` resource in case `delete_when_absent` is set to `false` [#332](https://github.com/oktadeveloper/terraform-provider-okta/pull/332). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Fix `okta_app_oauth` validation [#333](https://github.com/oktadeveloper/terraform-provider-okta/pull/333) and [#340](https://github.com/oktadeveloper/terraform-provider-okta/pull/340). Thanks, [@me](https://github.com/bogdanprodan-okta)!

## 3.9.0 (February 12, 2021)

ENHANCEMENTS:

* Add new `okta_admin_role_targets` resource [#325](https://github.com/oktadeveloper/terraform-provider-okta/pull/325). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Add `target_group_list` field to the `okta_group_role` resource [#256](https://github.com/oktadeveloper/terraform-provider-okta/pull/256). Thanks, [@ymylei](https://github.com/ymylei)!

BUGS:

* Fixed `subject_match_attribute` setup in the `okta_idp_saml` resource [#320](https://github.com/oktadeveloper/terraform-provider-okta/pull/320). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Fixed `users` setup when importing `okta_group` resource [#323](https://github.com/oktadeveloper/terraform-provider-okta/pull/323). Thanks, [@me](https://github.com/bogdanprodan-okta)!

## 3.8.0 (February 1, 2021)

ENHANCEMENTS:

* Add support for OAuth Authorization for Okta API [#290](https://github.com/oktadeveloper/terraform-provider-okta/pull/290). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Make `key_id` optional for `okta_app_saml_metadata` [#128](https://github.com/oktadeveloper/terraform-provider-okta/pull/128). Thanks, [@cludden](https://github.com/cludden)!
* Add new `okta_group_membership` resource [#252](https://github.com/oktadeveloper/terraform-provider-okta/pull/252). Thanks, [@ymylei](https://github.com/ymylei)!
* Add new `okta_group_role` resource [#255](https://github.com/oktadeveloper/terraform-provider-okta/pull/255). Thanks, [@ymylei](https://github.com/ymylei)!
* Add new `okta_idp_oidc` data source [#286](https://github.com/oktadeveloper/terraform-provider-okta/pull/286). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Add new `okta_app_oauth` data source [#293](https://github.com/oktadeveloper/terraform-provider-okta/pull/293). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Add new `okta_auth_server_policy` data source [#298](https://github.com/oktadeveloper/terraform-provider-okta/pull/298). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Add `usage` field to the `okta_network_zone` resource [#271](https://github.com/oktadeveloper/terraform-provider-okta/pull/271). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Add `okta_email` factor to the `okta_policy_mfa` resource [#269](https://github.com/oktadeveloper/terraform-provider-okta/pull/269). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Add `id` field to the `okta_users` data source [#288](https://github.com/oktadeveloper/terraform-provider-okta/pull/288). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Add `union` field to the `app_user_schema` resource [#291](https://github.com/oktadeveloper/terraform-provider-okta/pull/291). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Add `implicit_assignment` field to the `okta_app_oauth` resource [120](https://github.com/oktadeveloper/terraform-provider-okta/pull/120). Thanks, [Justin Lewis](https://github.com/jlew)!
* Add `issuer` and `issuer_mode` fields to the `okta_auth_server` data resource [#301](https://github.com/oktadeveloper/terraform-provider-okta/pull/301). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Add `login_mode` and `login_scopes` to the `okta_app_oauth` resource [#311](https://github.com/oktadeveloper/terraform-provider-okta/pull/311). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Add `single_logout_issuer`, `single_logout_url` and `single_logout_certificate` fields to the `okta_app_saml` resource [#307](https://github.com/oktadeveloper/terraform-provider-okta/pull/307). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Add `metadata_url` field to the `okta_app_saml` resource [#316](https://github.com/oktadeveloper/terraform-provider-okta/pull/316). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Remove `acs_binding` and `acs_type` from `okta_idp_oidc` as (they are not supported)[(https://developer.okta.com/docs/reference/api/idps/#oauth-2-0-and-openid-connect-endpoints-object)] by this resource [#286](https://github.com/oktadeveloper/terraform-provider-okta/pull/286). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Deprecate `acs_binding` argument for `okta_idp_saml` resource, as it [can only be set to `HTTP-POST`](https://developer.okta.com/docs/reference/api/idps/#assertion-consumer-service-acs-endpoint-object) [#286](https://github.com/oktadeveloper/terraform-provider-okta/pull/286). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Add a retry on `404` error in case Okta lagging during resource creation. Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Add validation for all URL-type fields.
* Various code improvements and documentation updates. Thanks, [@me](https://github.com/bogdanprodan-okta)!

BUGS:

* Ignore special groups (`BUILT_IN` and `APP_GROUP`) in the `group_memberships` field [#118](https://github.com/oktadeveloper/terraform-provider-okta/pull/118). Thanks, [@rasta-rocket](https://github.com/rasta-rocket)!
* Fix `inline_hooks` delete operation if the hooks were removed outside the provider [#288](https://github.com/oktadeveloper/terraform-provider-okta/pull/288). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Fix `group_memberships` populating in the `okta_user` data source [#284](https://github.com/oktadeveloper/terraform-provider-okta/pull/284). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Fix terraform import for the `app_user_schema` resource [#291](https://github.com/oktadeveloper/terraform-provider-okta/pull/291). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Fix delete operation for `auth_server_claim` resource in case claim has type `SYSTEM` [#283](https://github.com/oktadeveloper/terraform-provider-okta/pull/283). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Remove redundant `description` field from the `okta_app_saml` resource [#278](https://github.com/oktadeveloper/terraform-provider-okta/pull/278). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Add suppress function for the `features` field in the `okta_app_saml` resource since it's not currently possible to create/update provisioning features via the API [296](https://github.com/oktadeveloper/terraform-provider-okta/pull/296). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Remove `OAUTH_AUTHORIZATION_POLICY` from `okta_default_policy` and `okta_policy` since it's not supported by Okta API [#298](https://github.com/oktadeveloper/terraform-provider-okta/pull/298). Use `okta_auth_server_policy` instead. Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Fix status change in the `okta_auth_server_policy` resource [#299](https://github.com/oktadeveloper/terraform-provider-okta/pull/299). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Fix `user_name_template_*` fields setup for the apps resource [#309](https://github.com/oktadeveloper/terraform-provider-okta/pull/309/files). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Fix `refresh_token_window_minutes` minimum value in the `okta_auth_server_policy_rule` resource [#314](https://github.com/oktadeveloper/terraform-provider-okta/pull/314). Thanks, [@me](https://github.com/bogdanprodan-okta)!
* Fix `attribute_statements` field validation in the `okta_app_saml` resource [#313](https://github.com/oktadeveloper/terraform-provider-okta/pull/313). Thanks, [@me](https://github.com/bogdanprodan-okta)!

## 3.7.4 (December 28, 2020)

ENHANCEMENTS:

* Add `dependabot` to automate dependency updates [#259](https://github.com/oktadeveloper/terraform-provider-okta/pull/259). Thanks [@jlosito](https://github.com/jlosito)!
* Add `max_clock_skew` property to IdP SAML resource [#263](https://github.com/oktadeveloper/terraform-provider-okta/pull/263). Thanks [@me](https://github.com/bogdanprodan-okta)!

BUGS:

* Fix panic caused by a null pointer in `okta_policy_password` resource. [#262](https://github.com/oktadeveloper/terraform-provider-okta/pull/262). Thanks [@me](https://github.com/bogdanprodan-okta)!
* Add retries for creating/updating `okta_user_schema` resource. [#262](https://github.com/oktadeveloper/terraform-provider-okta/pull/262). Thanks [@me](https://github.com/bogdanprodan-okta)!

## 3.7.3 (December 24, 2020)

ENHANCEMENTS:

* Add call recovery for Okta password policy [#248](https://github.com/oktadeveloper/terraform-provider-okta/pull/248). Thanks [@me](https://github.com/bogdanprodan-okta)!
* Update data okta_group docs [#251](https://github.com/oktadeveloper/terraform-provider-okta/pull/251). Thanks [@ymylei](https://github.com/ymylei)!
* Adds `pattern` property for `okta_*_schema` resources [#159](https://github.com/oktadeveloper/terraform-provider-okta/pull/159). Thanks [@fitzoh](https://github.com/fitzoh) and [@me](https://github.com/bogdanprodan-okta)!
* Add retries on connection timeouts errors [#246](https://github.com/oktadeveloper/terraform-provider-okta/issues/246). Thanks [@me](https://github.com/bogdanprodan-okta)!

BUGS:

* Fixed handling rule with `INVALID` status [#250](https://github.com/oktadeveloper/terraform-provider-okta/pull/250). Thanks [@ymylei](https://github.com/ymylei)!

## 3.7.2 (December 18, 2020)

ENHANCEMENTS:

* Add logs to group data source for different cases [#150](https://github.com/oktadeveloper/terraform-provider-okta/pull/150). Thanks [@nathanbartlett](https://github.com/nathanbartlett)!
* Added missing documentation [#245](https://github.com/oktadeveloper/terraform-provider-okta/pull/245). Thanks [@me](https://github.com/bogdanprodan-okta)!

BUGS:

* Fix default name for idp_discovery [#244](https://github.com/oktadeveloper/terraform-provider-okta/pull/244). Thanks [@nickerzb](https://github.com/nickerzb)!
* Fix okta auth server policy rule resource causing panic [#245](https://github.com/oktadeveloper/terraform-provider-okta/pull/245). Thanks [@SBerda](https://github.com/SBerda) for submitting the [issue](https://github.com/oktadeveloper/terraform-provider-okta/issues/202) and [@me](https://github.com/bogdanprodan-okta) for fixing it!
* Fix `key_years_valid` defaulting to `2` during resource import [#245](https://github.com/oktadeveloper/terraform-provider-okta/pull/245). Thanks [@btsteve](https://github.com/btsteve) for submitting the [issue](https://github.com/oktadeveloper/terraform-provider-okta/issues/201) and [@me](https://github.com/bogdanprodan-okta) for fixing it!

## 3.7.1 (December 16, 2020)

ENHANCEMENTS:

* Add validation for user type [#242](https://github.com/oktadeveloper/terraform-provider-okta/pull/242).

BUGS:

* Fix state refresh for `okta_user_base_schema` and `okta_user_schema` [#242](https://github.com/oktadeveloper/terraform-provider-okta/pull/242).

## 3.7.0 (December 15, 2020)

ENHANCEMENTS:

* Add user types support [#183](https://github.com/oktadeveloper/terraform-provider-okta/pull/183). Thanks, [@rajnadimpalli](https://github.com/rajnadimpalli) and [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Add type to data okta group [#217](https://github.com/oktadeveloper/terraform-provider-okta/pull/217). Thanks, [@dangoslen](https://github.com/dangoslen)!
* Add `acs_endpoints` to SAML app (okta_app_saml) definition [#226](https://github.com/oktadeveloper/terraform-provider-okta/pull/226). Thanks, [@pranjalranjan](https://github.com/pranjalranjan)!
* Update terraform-plugin-sdk libraries, added possibility to set provider's log level [#220](https://github.com/oktadeveloper/terraform-provider-okta/pull/220). Thanks, [@bryantbiggs](https://github.com/bryantbiggs) and [@bogdanprodan-okta!](https://github.com/bogdanprodan-okta)
* Overhaul idp_discovery_rule documentation [#228](https://github.com/oktadeveloper/terraform-provider-okta/pull/228). Thanks [@eatplaysleep](https://github.com/eatplaysleep)!
* General documentation updates [#224](https://github.com/oktadeveloper/terraform-provider-okta/pull/224). Thanks, [@bryantbiggs](https://github.com/bryantbiggs)!

BUGS:

* Changed `okta_app_basic_auth` optional fields to required [issue 223](https://github.com/oktadeveloper/terraform-provider-okta/issues/223). Thanks, [@bryantbiggs](https://github.com/bryantbiggs)!
* Add idp discovery to allowed list of default policies [#233](https://github.com/oktadeveloper/terraform-provider-okta/pull/233). Thanks, [@nickerzb](https://github.com/nickerzb)!

## 3.6.1 (November 14, 2020)

ENHANCEMENTS:

* Remove 3rd party Okta SDK [#215](https://github.com/oktadeveloper/terraform-provider-okta/pull/215). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)
* Enhance `okta_app_auto_login` resource [#164](https://github.com/oktadeveloper/terraform-provider-okta/pull/164). Thanks, [@isometry](https://github.com/isometry)!
* Add group name to the error for group data call [#156](https://github.com/oktadeveloper/terraform-provider-okta/pull/156). Thanks, [@ymylei](https://github.com/ymylei)!

BUGS:

* Fix population of the user 'status' attribute [#206](https://github.com/oktadeveloper/terraform-provider-okta/pull/206). Thanks, [@isometry](https://github.com/isometry)!

## 3.6.0 (October 12, 2020)

ENHANCEMENTS:

* Upgrade to Okta SDK 2.0.0 [#203](https://github.com/oktadeveloper/terraform-provider-okta/pull/203). Thanks a ton! [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)
* Fix validation false positive when api_token is set via environment variable. [#147](https://github.com/oktadeveloper/terraform-provider-okta/pull/147). Thanks, [@jgeurts](https://github.com/jgeurts)
* Update required to optional and more [#208](https://github.com/oktadeveloper/terraform-provider-okta/pull/208), Thanks, me! :smile:

BUGS:

* Update config.go [#207](https://github.com/oktadeveloper/terraform-provider-okta/pull/207), Thanks, me! :smile:

## 3.5.1 (October 9, 2020)

ENHANCEMENTS:

* Update config.go [#192](https://github.com/oktadeveloper/terraform-provider-okta/pull/192), Thanks, [@bretterer](https://github.com/bretterer)!

BUGS:

* Documentation: Update okta_idp_metadata_saml correct example [#173](https://github.com/oktadeveloper/terraform-provider-okta/pull/173), Thanks, [@gaurdro](https://github.com/gaurdro) and [@netflash](https://github.com/netflash)!
* Documentation: Update warning in app_group_assignment.html.markdown [#172](https://github.com/oktadeveloper/terraform-provider-okta/pull/172), Thanks, [@ssttgg](https://github.com/ssttgg)!
* Renaming Go module as per the organization move [#195](https://github.com/oktadeveloper/terraform-provider-okta/pull/195), Thanks, [@stack72](https://github.com/stack72)!

## 3.5.0 (August 31, 2020)

ENHANCEMENTS:

* Add password import inline hook type. [#168](https://github.com/oktadeveloper/terraform-provider-okta/pull/168), Thanks, [@noinarisak](https://github.com/noinarisak) aka me! :tada:
* Add external_namespace property for app_user_schema and user_schema. [#102](https://github.com/oktadeveloper/terraform-provider-okta/pull/102), Thanks, [@thehunt33r](https://github.com/thehunt33r)!

BUGS:

* Fix inline hook example code to match version that is supported. [#175](https://github.com/oktadeveloper/terraform-provider-okta/pull/175), Thanks, [@noinarisak](https://github.com/noinarisak) me again! :smiley:
* Update app_group_assignment.html.markdown. [#165](https://github.com/oktadeveloper/terraform-provider-okta/pull/165), Thanks, [snolan-amount](https://github.com/snolan-amount)!

## 3.4.1 (July 31, 2020)

RELEASE:

* First release under oktadeveloper organization with binary published to [registry.hashicorp.com](https://registry.terraform.io/).

## 3.4.0 (July 30, 2020)

ENHANCEMENTS:

* Add resource definition for Okta Event Hooks. [#14](https://github.com/terraform-providers/terraform-provider-okta/pull/14), Thanks, [@mbudnek](https://github.com/mbudnek)!
* Adding support for GROUP_MEMBERSHIP_ADMIN & REPORT_ADMIN. [#138](https://github.com/terraform-providers/terraform-provider-okta/pull/138) Thanks, [ymylei](https://github.com/ymylei)!

BUG FIXES:

* Documentation corrections. Thanks, to these fine individuals!
  * [#126](https://github.com/terraform-providers/terraform-provider-okta/pull/126) [@ChristophShyper](https://github.com/ChristophShyper)
  * [#127](https://github.com/terraform-providers/terraform-provider-okta/pull/127) [@thekbb](https://github.com/thekbb)
  * [#151](https://github.com/terraform-providers/terraform-provider-okta/pull/151) [@varrunramani-okta](https://github.com/varrunramani-okta)

## 3.3.0 (May 29, 2020)

ENHANCEMENTS:

* Add user lockout notification channels. [#15](https://github.com/terraform-providers/terraform-provider-okta/pull/15), Thanks, [@thehunt33r](https://github.com/thehunt33r)!
* Adding support for SMS template changes. [#18](https://github.com/terraform-providers/terraform-provider-okta/pull/18) Thanks, [@gusChan](https://github.com/gusChan)!

## 3.2.0 (April 03, 2020)

BUG FIXES:

* Documentation, `id` is an output of `app_oauth`. [#98]() Thanks, [beyondbill](https://github.com/beyondbill)!

ENHANCEMENTS:

* Improve app filtering and update Terraform SDK. [#97](https://github.com/terraform-providers/terraform-provider-okta/pull/97) Thanks, [quantumew](https://github.com/quantumew)! :tada:

## 3.1.1 (March 18, 2020)

ENHANCEMENTS:

* Add unique property to UserSchema. [#12](https://github.com/terraform-providers/terraform-provider-okta/pull/12) Thanks, [@gusChan](https://github.com/gusChan)!

## 3.1.0 (February 19, 2020)

RELEASE:

* First release under terraform-providers organization with binary published to releases.hashicorp.com

## 3.0.0 (October 16, 2019)

FEATURES:

* Updated provider to support Terraform v0.12.0

## 3.0.1

FEATURES:

* **New Resource:** `okta_inline_hook`

ENHANCEMENTS:

* Add missing okta_idp_saml settings

## 3.0.2

ENHANCEMENTS:

* Use backoff/retries functionality for XML API calls

## 3.0.3

FEATURES:

* **New Data Source:** okta_idp_saml

ENHANCEMENTS:

* Support import user by email

## 3.0.4

FEATURES:

* **New Data Source:** Add okta_app_saml data source
* **New Data Source:** Add okta_app_metadata_saml data source
* **New Data Source:** Add okta_idp_metadata_saml data source

ENHANCEMENTS:

* Change type of custom_profile_attributes from map to JSON string to support all types

BUG FIXES:

* Fix group filter bug, filter_type and filter_value were not being sync'd

## 3.0.5

BUG FIXES:

* Fix bug introduced in v3.0.4. User data source was not updated to the new caustom_profile_attribute type
* Added test to cover this scenario, tests were passiing

## 3.0.6

ENHANCEMENTS:

* Allow client_id to be set on OIDC application, while also maintaining the computed version. With some auth methods, such as basic auth, this is possible.

## 3.0.7

ENHANCEMENTS:

* Add group_assignments for SAML and social IdPs

## 3.0.8

ENHANCEMENTS:

* Add issuer_mode to social IdP. Our test org does not have a custom domain setup, thus it was working there but not in other orgs. Hard to test both scenarios in one org.

## 3.0.9

FEATURES:

* **New Resource:** `okta_template_email`
* **New Resource:** `okta_group_roles`

## 3.0.10

FEATURES:

* **New Resource:** `okta_network_zone`

## 3.0.11

BUG FIXES:

* Fix occasional panic when creating a user schema see [issue 144](https://github.com/terraform-providers/terraform-provider-okta/issues/144)
* Users in LOCKED_OUT state are unlocked when config is ACTIVE [issue 225](https://github.com/terraform-providers/terraform-provider-okta/issues/225)

## 3.0.12

BUG FIXES:

* Ensure schema does not panic after retry

## 3.0.13

FEATURES:

* **New Resource:** `okta_user_base_schema`

ENHANCEMENTS:

* Add missing attribute, match_type and match_attribute, on social idp resource

## 3.0.14

BUG FIXES:

* Fix logic around including/excluding networks on policy rules

## 3.0.15

ENHANCEMENTS:

* Update Okta SDK
* Filter out GROUP based admin roles when processing user `admin_roles` attribute

## 3.0.16

* Fix issues around `okta_policy_rule_idp_discovery`
  * `app_include` and `app_exlcude` were missing required properties
  * `user_identifier_type` was being added even when not defined, causing API errors
* Fix integer array type

## 3.0.17

FEATURES:

* **New Resource:** `okta_app_user_schema`
* **New Resource:** `okta_app_user_base_schema`
* **New Resource:** `okta_app_user` resource
* **New Resource:** `okta_app_group` resource

ENHANCEMENTS:

* Add `required` field to base schema

## 3.0.18

ENHANCEMENTS:

* Support SHA-1 signing algorithm on IdPs

BUG FIXES:

* Fix bug where audience is reset on IdP update because it is omitted from the payload

## 3.0.19

BUG FIXES:

* Fix diff issues around `okta_policy_rule_idp_discovery`
* Allow `provisioning_action` for IdPs to be set to `DISABLED`

## 3.0.20

BUG FIXES:

* Fix `okta_auth_server_claim`, `group_filter_type` could not be set to `STARTS_WITH` due to a typo

## 3.0.21

ENHANCEMENTS:

* Expose scope property on `okta_user_schema`
* Allow setting of OAuth application visibility settings

## 3.0.22

BUG FIXES:

* Send `profileMaster` along with IdP, so the config is recognized by Okta API
* Fix bug in SDK related to retries and the request body being empty on subsequent requests.

## 3.0.23

ENHANCEMENTS:

* Add `external_name` property to the `okta_app_user_schema` and `okta_user_schema`

## 3.0.24

ENHANCEMENTS:

* Support `profile` on `okta_oauth_app` resource

## 3.0.25

ENHANCEMENTS:

* Support setting an auth server scope as the default
* Support `profile` and `priority` on `okta_app_group_assignment`
* Support `profile` on `okta_app_user`

BUG FIXES:

* Fix bug with supporting `profile` on `okta_oauth_app` resource

## 3.0.26

ENHANCEMENTS:

* Support array enums in `okta_user_schema` and `okta_app_user_schema` as `array_enum` and `array_one_of`

## 3.0.27

ENHANCEMENTS:

* Update refresh token window validation to account for new upper limit of 5 years

## 3.0.28

BUG FIXES:

* Remove resource from state on 404. ([#269](https://github.com/terraform-providers/terraform-provider-cherryservers/issues/269))

## 3.0.29

BUG FIXES:

* Ensure we safely sync auth server properties. ([#299](https://github.com/terraform-providers/terraform-provider-cherryservers/issues/299))
* MANUAL rotation mode can only be set on an auth server on update. Ensure we run update after create for that scenario. ([#287](https://github.com/terraform-providers/terraform-provider-cherryservers/issues/287))

## 3.0.30

ENHANCEMENT:

* Update to new separate Terraform SDK ([#307](https://github.com/terraform-providers/terraform-provider-cherryservers/issues/307))

## 3.0.31

BUG FIXES:

* Ensure `okta_app_group_assignment` resource syncs using the right read function. ([#307](https://github.com/terraform-providers/terraform-provider-cherryservers/issues/307))

## 3.0.32

BUG FIXES:

* Ensure `okta_app_group_assignment` and `okta_app_user` resources properly take multiple ids on the import functions. ([#307](https://github.com/terraform-providers/terraform-provider-cherryservers/issues/307))
* Ensure `okta_user` does not error on 404 ([#313](https://github.com/terraform-providers/terraform-provider-cherryservers/issues/313))

## 3.0.33

FEATURES:

* **New Resource:** `okta_profile_mapping` ([#246](https://github.com/terraform-providers/terraform-provider-cherryservers/issues/246))
* **New Resource:** `okta_app_basic_auth` ([#329](https://github.com/terraform-providers/terraform-provider-cherryservers/issues/329))

## 3.0.34

BUG FIXES:

* Policy values could not be set to 0. Doing so resulted in the SDK omitting them, resulting in Okta resetting the values to default.

## 3.0.35

ENHANCEMENT:

* Require target_id on `okta_profile_mapping` to avoid ambiguity

FEATURES:

* **New Data Source:** `okta_user_profile_mapping_source` ([#340](https://github.com/terraform-providers/terraform-provider-cherryservers/issues/340))

## 3.0.36

BUG FIXES

* Schema merging helper function was mutating input schema causing side effects when used in a particular way. Used shallow copying to avoid this side effect. ([#338](https://github.com/terraform-providers/terraform-provider-cherryservers/issues/338))
* Ensure response is not nil when checking status code ([#307](https://github.com/terraform-providers/terraform-provider-cherryservers/issues/307))

## 3.0.37

BUG FIXES

* Ensure `index` is sync'd on import to avoid recreation.

## 3.0.38

ENHANCEMENT:

* Support `password`, `recovery_answer`, and `recovery_question` as attributes on the `okta_user` resource.
