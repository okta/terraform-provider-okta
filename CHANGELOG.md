# Changelog

## 3.42.0 (February 10, 2023)

### NEW - RESOURCES, DATA SOURCES, PROPERTIES, ATTRIBUTES, ENV VARS:

* New data source [`okta_domain`](https://registry.terraform.io/providers/okta/okta/latest/docs/data-sources/domain) see PR 1447 notes in BUG FIXES
* Actual PEM text values in `okta_domain_certificate` for attributes `certificate`, `certificate_chain`, and `private_key`, see PR 1447 notes in BUG FIXES
* New attribute `roles` in data source `okta_user` [#1437](https://github.com/okta/terraform-provider-okta/pull/1437). Thanks, [@duytiennguyen-okta](https://github.com/duytiennguyen-okta)!

### BUG FIXES:

* Don't md5sum to save space on `okta_domain_certificate` values for attributes `certificate`, `certificate_chain`, and `private_key`, per TF best practices [#1447](https://github.com/okta/terraform-provider-okta/pull/1447). Thanks, [@monde](https://github.com/monde)!
* Remove org type restrictions and artificial input check on `type` attribute for  data source `okta_policy`  [#1445](https://github.com/okta/terraform-provider-okta/pull/1445). Thanks, [@monde](https://github.com/monde)!
* Improve resource `okta_app_saml` documentation  [#1439](https://github.com/okta/terraform-provider-okta/pull/1439). Thanks, [@exitcode0](https://github.com/exitcode0)!

## 3.41.0 (January 27, 2023)

### PROJECT IMPROVEMENTS:

* Enable okta_password authenticator for okta_policy_mfa [#1210](https://github.com/okta/terraform-provider-okta/pull/1210). Tests [#1427](https://github.com/okta/terraform-provider-okta/pull/1427). Thanks, [@nickrmc83](https://github.com/nickrmc83)!
* Update resource documentation with link to role-type api doc references [#1430](https://github.com/okta/terraform-provider-okta/pull/1430). Thanks, [@noinarisak](https://github.com/noinarisak)!

## 3.40.0 (January 09, 2023)

### BUG FIXES:

* Fixes ThreatInsight Configuration Continuously Reordering [#1398](https://github.com/okta/terraform-provider-okta/pull/1398). Thanks, [@duytiennguyen-okta](https://github.com/duytiennguyen-okta)!
* Fixes rate limit accounting for `/api/v1/authorizationServers` endpoints [#1420](https://github.com/okta/terraform-provider-okta/pull/1420). Thanks, [@monde](https://github.com/monde)!

### PROJECT IMPROVEMENTS:

* Improve `app_user_base_schema_property` documentation [#1407](https://github.com/okta/terraform-provider-okta/pull/1407). Thanks, [@robgero](https://github.com/robgero)!
* Fix `TestAccOktaAppSignOnPolicy` ACC test [#1412](https://github.com/okta/terraform-provider-okta/pull/1412). Thanks, [@noinarisak](https://github.com/noinarisak)!

## 3.39.0 (November 18, 2022)

### NEW - RESOURCES, DATA SOURCES, PROPERTIES, ATTRIBUTES, ENV VARS:

* `okta_authenticator` resource and data source [#1379](https://github.com/okta/terraform-provider-okta/pull/1379). Thanks, [@monde](https://github.com/monde)!
  * Added argment `provider_json` allowing provider information to be set with JSON on the authenticator
  * Improved resource behavior in regards to Okta API's hard create, soft create, and soft delete of authenticators
  * Improved data source and resource documentation

* Added `authentication_policy` argument to resource `okta_app_bookmark` [#1376](https://github.com/okta/terraform-provider-okta/pull/1376). Thanks, [@jakezarobsky-8451](https://github.com/jakezarobsky-8451)!

* `okta_user` resrouce [#1372](https://github.com/okta/terraform-provider-okta/pull/1372). Thanks, [@monde](https://github.com/monde)!
  * Adds `skip_roles` flag to allow explicit gating on the attempt to set roles
  * Swallows and warns on 403 errors when roles API is called and API token is less than super admin scope
  * Improved data source and resource documentation

### ENHANCEMENTS:

* `okta_idp_saml` gracefully handles 401 errors when setting profile mapping [#1355](https://github.com/okta/terraform-provider-okta/pull/1355)/[#1369](https://github.com/okta/terraform-provider-okta/pull/1369). Thanks, [@deorus](https://github.com/deorus)!
* Rate limits handler rules are generated from Okta service's actual code [#1356](https://github.com/okta/terraform-provider-okta/pull/1356). Thanks, [@monde](https://github.com/monde)!

### BUG FIXES:

* Address parallel API calls in `okta_user_base_schema_property` resource [#1351](https://github.com/okta/terraform-provider-okta/pull/1351). Thanks, [@monde](https://github.com/monde)!

### PROJECT IMPROVEMENTS:

* Updated `okta_app_user_schema_property`, `okta_auth_server_policy`, and `okta_auth_server_policy_rule` resource documentation [#1348](https://github.com/okta/terraform-provider-okta/pull/1348). Thanks, [@zlitberg](https://github.com/zlitberg)!
* Document a PEM and JWKS example for the `okta_app_oauth` resource [#1350](https://github.com/okta/terraform-provider-okta/pull/1350). Thanks, [@monde](https://github.com/monde)!

## 3.38.0 (October 28, 2022)

BUG FIXES:

* Address potential panic in resource `okta_app_group_assignments`'s `profile` attribute [#1345](https://github.com/okta/terraform-provider-okta/pull/1345). Thanks, [@monde](https://github.com/monde)!
* Address potential panic in resource `okta_inline_hook`s `auth` attribute [#1337](https://github.com/okta/terraform-provider-okta/pull/1337). Thanks, [@monde](https://github.com/monde)!
* Fully document and refine `okta_app_oauth`'s `pkce_required` attribute required if `token_endpoint_auth_method` is "none" [#1327](https://github.com/okta/terraform-provider-okta/pull/1327). Thanks, [@monde](https://github.com/monde)!

## 3.37.0 (October 04, 2022)

NEW - RESOURCES, DATA SOURCES, PROPERTIES, ATTRIBUTES, ENV VARS:

* Add `ui_schema_id` property to resource `okta_policy_rule_profile_enrollment` [#1324](https://github.com/okta/terraform-provider-okta/pull/1324). Thanks, [@monde](https://github.com/monde)!
* Add `CUSTOM` to list of group role types in datasource `okta_role_subscription` [#1320](https://github.com/okta/terraform-provider-okta/pull/1320). Thanks, [@monde](https://github.com/monde)!
* From PR [#1322](https://github.com/okta/terraform-provider-okta/pull/1322). Thanks, [@monde](https://github.com/monde)!
  * Improved resource `okta_email_customization` behavior with new property `force_is_default` with regards to the `is_default` property
  * Added explicit errors for Classic orgs trying to make use of OIE features. Error messages refer to corresponding online documentation
    * datasource `okta_app_signon_policy`
    * datasource `okta_authenticator`
    * resource `okta_app_signon_policy`
    * resource `okta_authenticator`
    * resource `okta_captcha`
    * resource `okta_captcha_org_wide_settings`
    * resource `okta_policy_profile_enrollment`
    * resource `okta_policy_profile_enrollment_apps`
    * resource `okta_policy_rule_profile_enrollment`

BUG FIXES:

* Fixed `okta_app_user_schema_property` for non string enum types [#1316](https://github.com/okta/terraform-provider-okta/pull/1316). Thanks, [@duytiennguyen-okta](https://github.com/duytiennguyen-okta)!
* From PR [#1322](https://github.com/okta/terraform-provider-okta/pull/1322). Thanks, [@monde](https://github.com/monde)!
  * Fixed (unreported) bug where resource `okta_org_configuration` would null out org settings 
  * Fixed an ACC test with resource `okta_user_schema_property` that would cause a incorrect login flow blocking out the admin
  * Fixed/improved sms template tests
  * Marked the schema enum boolean tests skip as there is an issue with the public API / monolith
  * Cleaned up code paths for default/system policy getting/setting for apps and policies
  * Fixed and/or cleaned up a number of other ACC tests

PROJECT IMPROVEMENTS:

* Correct `okta_email_customization` docs [#1310](https://github.com/okta/terraform-provider-okta/pull/1310). Thanks, [@lucascantor](https://github.com/lucascantor)!


## 3.36.0 (September 14, 2022)

NEW - RESOURCES, DATA SOURCES, PROPERTIES, ATTRIBUTES, ENV VARS:

* Add `client_secret` attribute on data source `okta_app_oauth` [#1307](https://github.com/okta/terraform-provider-okta/pull/1307) Thanks, [@dkulchinsky](https://github.com/dkulchinsky), [@monde](https://github.com/monde), [@rickardp](https://github.com/rickardp)!
  * oauth app data source: allow to retrieve client_secret [#1285](https://github.com/okta/terraform-provider-okta/pull/1285)
  * client_secret is missing from okta_app_oauth data source [#1279](https://github.com/okta/terraform-provider-okta/issues/1279)
  * Added support for retrieving client secret from okta_app_oauth data source [#1280](https://github.com/okta/terraform-provider-okta/pull/1280)
* Adds `pkce_required` property to resource `okta_app_oauth` [#1305](https://github.com/okta/terraform-provider-okta/pull/1305) Thanks, [@monde](https://github.com/monde)!
  * Add support to pkce_required property for OIDC app integrations [#1241](https://github.com/okta/terraform-provider-okta/issues/1241)
* Schema updates for `okta_idp_oidc` and `okta_idp_social` [#1297](https://github.com/okta/terraform-provider-okta/pull/1297) Thanks, [@monde](https://github.com/monde)!
  * okta_idp_oidc does not support DYNAMIC issuer_mode [#1288](https://github.com/okta/terraform-provider-okta/issues/1288)
  * Okta Social IDP with Type Github [#1293](https://github.com/okta/terraform-provider-okta/issues/1293)

BUG FIXES:

* Policy Rule Retry On InternalServerError [#1273](https://github.com/okta/terraform-provider-okta/pull/1273) Thanks, [@ymylei](https://github.com/ymylei)!
* Set SAML Features To Computed [#1272](https://github.com/okta/terraform-provider-okta/pull/1272) Thanks, [@ymylei](https://github.com/ymylei)!
* Errors when adding user to group are incorrectly ignored. [#1301](https://github.com/okta/terraform-provider-okta/pull/1301) Thanks, [@monde](https://github.com/monde)!
  * prevent error overwrite in addGroupMember [#1269](https://github.com/okta/terraform-provider-okta/issues/1269)
* Okta Group Schema Null Handling [#1271](https://github.com/okta/terraform-provider-okta/pull/1271) Thanks, [@ymylei](https://github.com/ymylei)!
* Diff Suppression on SLO Certs [#1270](https://github.com/okta/terraform-provider-okta/pull/1270) Thanks, [@ymylei](https://github.com/ymylei)!
* Nil guard on app.Settings.OauthClient [#1300](https://github.com/okta/terraform-provider-okta/pull/1300) Thanks, [@monde](https://github.com/monde)!
  * Provider crashes when doing a data source lookup of an app with different type than the label it is using for the lookup. * [#1082](https://github.com/okta/terraform-provider-okta/issues/1082)
* Nil guard on resource `set _links` value [#1299](https://github.com/okta/terraform-provider-okta/pull/1299) Thanks, [@monde](https://github.com/monde)!
  * Error when creating okta_resource_set [#1278](https://github.com/okta/terraform-provider-okta/issues/1278)
* Guard from nil pointer dereference [#1298](https://github.com/okta/terraform-provider-okta/pull/1298) Thanks, [@monde](https://github.com/monde)!
  * Plugin crash when importing okta_policy_signon [#1294](https://github.com/okta/terraform-provider-okta/issues/1294)

PROJECT IMPROVEMENTS:

* Variable Types Update - Documentation [#1276](https://github.com/okta/terraform-provider-okta/pull/1276)
Thanks, [@pro4tlzz](https://github.com/pro4tlzz)!
* Update brand.html.markdown [#1281](https://github.com/okta/terraform-provider-okta/pull/1281) Thanks, [@monde](https://github.com/monde)!
* Update theme.html.markdown [#1282](https://github.com/okta/terraform-provider-okta/pull/1282) Thanks, [@monde](https://github.com/monde)!

## 3.35.0 (August 25, 2022)

NEW - RESOURCES, DATA SOURCES, PROPERTIES, ATTRIBUTES, ENV VARS:

* Adds customizable Timeouts to resources/data that rely on syncing users and groups to avoid context.DeadlineExceeded
 [#1207](https://github.com/okta/terraform-provider-okta/pull/1207). Thanks, [@emanor-okta](https://github.com/emanor-okta)!
  * [Terraform documentation: Resources - Retries and Customizable Timeouts](https://www.terraform.io/plugin/sdkv2/resources/retries-and-customizable-timeouts)
  * Resources: `okta_app_auto_login`, `okta_app_basic_auth`, `okta_app_bookmark`, `okta_app_group_assignment`, `okta_app_oauth`, `okta_app_saml`, `okta_app_secure_password_store`, `okta_app_shared_credentials`, `okta_app_swa`

BUG FIXES:

* Correctly collect network zones in datasource `okta_network_zone` [#1239](https://github.com/okta/terraform-provider-okta/pull/1239). Thanks, [@natmariam](https://github.com/natmariam)!
* Adding `CHROMEOS` to `os_type` in `platform_include` [#1261](https://github.com/okta/terraform-provider-okta/pull/1261). Thanks, [@monde](https://github.com/monde)!
* Update okta-sdk-golang that correctly caches OAuth2 access tokens [#1262](https://github.com/okta/terraform-provider-okta/pull/1262). Thanks, [@monde](https://github.com/monde)!
* Update role types validation on resource `okta_role_subscription` [#1265](https://github.com/okta/terraform-provider-okta/pull/1265). Thanks, [@monde](https://github.com/monde)!
* Correct pagination to list all email templates on data source `okta_email_templates` [#1266](https://github.com/okta/terraform-provider-okta/pull/1266). Thanks, [@monde](https://github.com/monde)!

PROJECT IMPROVEMENTS:

* Show current version for provider config in documentation [#1256](https://github.com/okta/terraform-provider-okta/pull/1256). Thanks, [@ErelAdoni](https://github.com/ErelAdoni)!
* Code clean up from go vet and format [#1264](https://github.com/okta/terraform-provider-okta/pull/1264). Thanks, [@monde](https://github.com/monde)!

## 3.34.0 (August 12, 2022)

BUG FIXES:

* Fix concurrency issue in resource `okta_auth_server_policy_rule` that could cause 500s in the Okta API as well as not preserve priority ordering even when `depends_on`is present [#1248](https://github.com/okta/terraform-provider-okta/pull/1248). Thanks, [@monde](https://github.com/monde)!

PROJECT IMPROVEMENTS:

* Fix typo provider test [#1229](https://github.com/okta/terraform-provider-okta/pull/1229). Thanks, [@lukas-hetzenecker](https://github.com/lukas-hetzenecker)!

## 3.33.0 (August 02, 2022)

BUG FIXES:

* Fix "error invalid configuration" error introduced in v3.32.0 release; includes unit tests to verify fix. [#1234](https://github.com/okta/terraform-provider-okta/pull/1234). Thanks, [@ericnorris](https://github.com/ericnorris)!

## 3.32.0 (July 29, 2022)

NEW - RESOURCES, DATA SOURCES, PROPERTIES, ATTRIBUTES, ENV VARS:

* Add keys attribute to okta_app_saml resource [#1206](https://github.com/okta/terraform-provider-okta/pull/1206). Thanks, [@ericnorrisl](https://github.com/ericnorris) and [@slichtenthal](https://github.com/slichtenthal)!
* Export the app embed url for saml apps [#1215](https://github.com/okta/terraform-provider-okta/pull/1215). Thanks, [@felixcolaci](https://github.com/felixcolaci)!
* Ability to configure the provider with an access (Bearer) token [#1222](https://github.com/okta/terraform-provider-okta/pull/1222). Thanks, [@ericnorrisl](https://github.com/ericnorris)!
* Add `privateKeyId` private key signing support available in okta-sdk-golang client [#1223](https://github.com/okta/terraform-provider-okta/pull/1223). Thanks, [@powellchristoph](https://github.com/powellchristoph)!

BUG FIXES:

* Fix "no default policy found" bug, includes ability for provider to discover if it is running against an OIE or Classic org [#1224](https://github.com/okta/terraform-provider-okta/pull/1224). Thanks, [@monde](https://github.com/monde)!

## 3.31.0 (July 08, 2022)

NEW - RESOURCES, DATA SOURCES, PROPERTIES, ATTRIBUTES, ENV VARS:

* New resource `okta_app_signon_policy` [#1193](https://github.com/okta/terraform-provider-okta/pull/1193). Thanks, [@felixcolaci](https://github.com/felixcolaci)!
* Added property `inactivity_period` to resource `okta_app_signon_policy_rule` [#1184](https://github.com/okta/terraform-provider-okta/pull/1181). Thanks, [@monde](https://github.com/monde)!
* Property `issuer_mode` can be `"CUSTOM_URL"`, `"ORG_URL"`, or `"DYNAMIC"` on resource `okta_auth_server_default` [#1197](https://github.com/okta/terraform-provider-okta/pull/1197). Thanks, [@monde](https://github.com/monde)!

BUG FIXES:

* Correct API endpoint and call for resource `okta_policy_profile_enrollment_apps` [#1191](https://github.com/okta/terraform-provider-okta/pull/1191). Thanks, [@felixcolaci](https://github.com/felixcolaci)!
* Fix resources pagination in resource `okta_resource_set` for resource items greater than 100 [#1196](https://github.com/okta/terraform-provider-okta/pull/1196). Thanks, [@monde](https://github.com/monde)!

ENHANCEMENTS:

* Update documentation on resource `okta_policy_mfa` and `okta_policy_mfa_default` for required FF `OKTA_MFA_POLICY` and when FF `ENG_ENABLE_OPTIONAL_PASSWORD_ENROLLMENT` is enabled [#1176](https://github.com/okta/terraform-provider-okta/pull/1176). Thanks, [@monde](https://github.com/monde)!

## 3.30.0 (June 22, 2022)

BUG FIXES:

* Correct issuer mode value in embedded `groups_claim` of an `okta_app_oauth` resource [#1167](https://github.com/okta/terraform-provider-okta/pull/1167). Thanks, [@monde](https://github.com/monde)!
* Resource `okta_app_oauth` property`redirect_uris` is a list, not a set, and needs to maintain order. [#1171](https://github.com/okta/terraform-provider-okta/pull/1171). Thanks, [@monde](https://github.com/monde)!
* Fix JSON serialization errors that group and user schemas experience when `enum` and `one_of` properties are utilized with a `type` value other than `string` [#1178](https://github.com/okta/terraform-provider-okta/pull/1178). Thanks, [@monde](https://github.com/monde)!

ENHANCEMENTS:

* Add `no-stalebot` label exemption for GH stalebot action [#1180](https://github.com/okta/terraform-provider-okta/pull/1180). Thanks, [@monde](https://github.com/monde)!


## 3.29.0 (June 09, 2022)

ENHANCEMENTS:

* HTTP proxy feature with `OKTA_HTTP_PROXY` alternative to `OKTA_ORG_NAME`+`OKTA_BASE_URL` [#1142](https://github.com/okta/terraform-provider-okta/pull/1142). Thanks, [@ido50](https://github.com/ido50)!
* Full support for Duo authenticator [#1146](https://github.com/okta/terraform-provider-okta/pull/1146). Thanks, [@monde](https://github.com/monde)!
* Improve data source `okta_user` and `okta_users` and a bug fix [#1159](https://github.com/okta/terraform-provider-okta/pull/1159). Thanks, [@exitcode0](https://github.com/exitcode0), [@monde](https://github.com/monde)!
* Update latest list of Custom Role Permission properties on resource `okta_admin_role_custom` [#1160](https://github.com/okta/terraform-provider-okta/pull/1160). Thanks, [@tim-fitzgerald](https://github.com/tim-fitzgerald)!

BUG FIXES:

* Remove incorrect attributes `response_signature_algorithm`, and `response_signature_scope` from resource `okta_idp_oidc` [#1156](https://github.com/okta/terraform-provider-okta/pull/1156). Thanks, [@monde](https://github.com/monde)!
* Reestablish old behavior of `okta_group_memberships` resource, add toggle to track all users [#1161](https://github.com/okta/terraform-provider-okta/pull/1161). Thanks, [@monde](https://github.com/monde)!

PROJECT IMPROVEMENTS:

* Fix typo in data source `okta_email_template` documentation [#1157](https://github.com/okta/terraform-provider-okta/pull/1157). Thanks, [@monde](https://github.com/monde)!
* ACC tests maintenance [#1158](https://github.com/okta/terraform-provider-okta/pull/1158). Thanks, [@monde](https://github.com/monde)!

NEW - RESOURCES, DATA SOURCES, PROPERTIES, ATTRIBUTES, ENV VARS:

* ENV VAR
  * `OKTA_HTTP_PROXY` alternative to `OKTA_ORG_NAME`+`OKTA_BASE_URL`
* Data Sources
  * `okta_user`
    * `delay_read_seconds` property to assist dealing with data eventual consistency
  * `okta_users`
    * `include_roles` property to signal admin roles for each user should also be gathered
    * `delay_read_seconds` property to assist dealing with data eventual consistency
* Resources
  * `okta_group_memberships`
    * `track_all_users` track all users of group, not just those when resource was initialized

## 3.28.0 (May 24, 2022)

ENHANCEMENTS:

* Add `system` attribute to `okta_auth_server_scope` resource [#1112](https://github.com/okta/terraform-provider-okta/pull/1112). Thanks, [@monde](https://github.com/monde)!
* Refine search criteria precision in `okta_app` data source [#1115](https://github.com/okta/terraform-provider-okta/pull/1115). Thanks, [@monde](https://github.com/monde)!
* `okta_group` adds delay argument; Refine `okta_group_memberships` resource and add tests. Update documentation [#1120](https://github.com/okta/terraform-provider-okta/pull/1120). Thanks, [@monde](https://github.com/monde)!
* Add `com.okta.telephony.provider` hook type to `okta_inline_hooks` resource [#1132](https://github.com/okta/terraform-provider-okta/pull/1132). Thanks, [@monde](https://github.com/monde)!

BUG FIXES:

* Fix type in custom role permissions for `okta_admin_role_custom` resource [#1116](https://github.com/okta/terraform-provider-okta/pull/1116). Thanks, [@faurel](https://github.com/faurel)!
* Fix pagination bug in `okta_group_memberships` [#1125](https://github.com/okta/terraform-provider-okta/pull/1125). Thanks, [@monde](https://github.com/monde)!
* Reverted commit on `okta_policy_rule_sign_on` resource that adversely affected `SPECIFIC_IDP` [#1133](https://github.com/okta/terraform-provider-okta/pull/1133). Thanks, [@monde](https://github.com/monde)!
* Corrected signature defaults on `okta_idp_oidc`, `okta_idp_saml`, and `okta_idp_social` resources [#1134](https://github.com/okta/terraform-provider-okta/pull/1134). Thanks, [@monde](https://github.com/monde)!
* Fixed regression on `okta_group_memberships` resource with 0 users [#1138](https://github.com/okta/terraform-provider-okta/pull/1138). Thanks, [@exitcode0](https://github.com/exitcode0)!

PROJECT IMPROVEMENTS:

* Update `okta_template_email` documentation [#1113](https://github.com/okta/terraform-provider-okta/pull/1113). Thanks, [@monde](https://github.com/monde)!
* ACC Test for `okta_rate_limiting` resource and update documentation [#1121](https://github.com/okta/terraform-provider-okta/pull/1121). Thanks, [@monde](https://github.com/monde)!
* Note that `okta_group_membership` is deprecated in the documentation [#1122](https://github.com/okta/terraform-provider-okta/pull/1122). Thanks, [@monde](https://github.com/monde)!
* Update documentation on `okta_app_oauth` explaining how reset a client secret [#1127](https://github.com/okta/terraform-provider-okta/pull/1127). Thanks, [@monde](https://github.com/monde)!
* Update deprecation notice on `okta_template_email` resource documentation [#1136](https://github.com/okta/terraform-provider-okta/pull/1136). Thanks, [@monde](https://github.com/monde)!
* ACC Test on `okta_group_memberships` resource with 0 users [#1139](https://github.com/okta/terraform-provider-okta/pull/1139). Thanks, [@monde](https://github.com/monde)!

## 3.27.0 (May 13, 2022)

ENHANCEMENTS:

* Data sources and resources for branded themes [#1104](https://github.com/okta/terraform-provider-okta/pull/1104). Thanks, [@monde](https://github.com/monde)!
  * Data Sources
    * `okta_themes`
    * `okta_theme`
  * Resources
    * `okta_theme`

BUG FIXES:
* Soft revert of diff suppress on `okta_policy_password` and `okta_policy_password_default` resources [#1108](https://github.com/okta/terraform-provider-okta/pull/1108). Thanks, [@monde](https://github.com/monde)!

PROJECT IMPROVEMENTS:

* Removed confusing and inaccurate information about Duo and Yubikey support in resource `okta_authenticator` [#1103](https://github.com/okta/terraform-provider-okta/pull/1103). Thanks, [@monde](https://github.com/monde)!
* Fixed formatting in docs for a markdown rendering quirk of the Terraform Registry [#1096](https://github.com/okta/terraform-provider-okta/pull/1096). Thanks, [@monde](https://github.com/monde)!

## 3.26.0 (May 06, 2022)

ENHANCEMENTS:

* Data sources and resources for branded email customization [#1089](https://github.com/okta/terraform-provider-okta/pull/1089). Thanks, [@monde](https://github.com/monde)!
  * Data Sources
    * `okta_brands`
    * `okta_brand`
    * `okta_email_customizations`
    * `okta_email_customization`
    * `okta_email_templates`
    * `okta_email_template`
  * Resources
    * `okta_brand`
    * `okta_email_customization`
* Allow user lookup by group membership; data source `okta_users` gets `group_id` property. [#998](https://github.com/okta/terraform-provider-okta/pull/998). Thanks, [@BrentSouza](https://github.com/BrentSouza)!

PROJECT IMPROVEMENTS:

* Note `browser` type for SPA apps in app_oauth.html.markdown documentation [#580](https://github.com/okta/terraform-provider-okta/issues/580). Thanks, [@monde](https://github.com/monde)!
* Add docs to represent USER_ADMIN in group_role.html.markdown documentation [#1075](https://github.com/okta/terraform-provider-okta/pull/1075). Thanks, [@naveen-vijay](https://github.com/naveen-vijay)!

## 3.25.1 (April 26, 2022)

BUGS:

 * Fix incomplete `compound_search_operator` on data source `okta_users`. [#1077](https://github.com/okta/terraform-provider-okta/issues/1077). Thanks, [@monde](https://github.com/monde)!
 * Fix default value regression on `okta_policy_rule_sign_on` for `identity_provider` attribute. [#1079](https://github.com/okta/terraform-provider-okta/issues/1079). Thanks, [@monde](https://github.com/monde)!

## 3.25.0 (April 21, 2022)

ENHANCEMENTS:
* Upgrade okta-sdk-golang to v2.12.1. [#1001](https://github.com/okta/terraform-provider-okta/pull/1001). Thanks, [@monde](https://github.com/monde)!
  * Removing/Updating local sdk code
    * Application.UploadApplicationLogo
    * Authenticator
    * EnrollFactor
    * LinkedObjects
    * PasswordPolicy
    * ProfileMapping
    * Subscription
    * UserFactor
  * Fixed ACC tests
    * TestAccOktaAppSignOnPolicyRule
    * TestAccOktaDataSourceIdpSocial_read
    * TestAccOktaDefaultPasswordPolicy
    * TestAccOktaIdpSocial_crud
    * TestAccOktaPolicyPassword_crud
    * TestAccOktaPolicySignOn_crud
    * TestAccAppOAuthApplication_postLogoutRedirectCrud
  * Backoff/retry on application delete
* Update okta_app_saml resource documentation. [#1076](https://github.com/okta/terraform-provider-okta/pull/1076). Thanks, [@jphuynh](https://github.com/jphuynh)!

## 3.24.0 (April 15, 2022)

ENHANCEMENTS:
* Document group rule name max and min length [#1068](https://github.com/okta/terraform-provider-okta/pull/1068). Thanks, [@monde](https://github.com/monde)!

BUGS:

* Correctly change password on Okta user resource [#1060](https://github.com/okta/terraform-provider-okta/pull/1060). Thanks, [@BalaGanaparthi](https://github.com/BalaGanaparthi)!
  * Uses change password flow if old password is present
  * Uses set password flow if only password is present

## 3.23.0 (April 08, 2022)

ENHANCEMENTS:

* Okta User and Okta Users search can use free form filter [#1027](https://github.com/okta/terraform-provider-okta/pull/1027). Thanks, [@cbrgm](https://github.com/cbrgm)!
* Uniqueness of logo file is by SHA only, not SHA and local file path [#1039](https://github.com/okta/terraform-provider-okta/pull/1039). Thanks, [@bobtfish](https://github.com/bobtfish)!
* Improve Okta Groups custom profile attributes for use in Terraform expressions [#1041](https://github.com/okta/terraform-provider-okta/pull/1041). Thanks, [@exitcode0](https://github.com/exitcode0)!

PROJECT IMPROVEMENTS:

* Add valid options for status field in user.html.markdown documentation [#1040](https://github.com/okta/terraform-provider-okta/pull/1040). Thanks, [@exitcode0](https://github.com/exitcode0)!
* Fix markdown typo in role_subscription.html.markdown documentation [#1049](https://github.com/okta/terraform-provider-okta/pull/1049). Thanks, [@lucascantor](https://github.com/lucascantor)!
* Fix markdown typo in role_subscription.html.markdown documentation [#1050](https://github.com/okta/terraform-provider-okta/pull/1050). Thanks, [@lucascantor](https://github.com/lucascantor)!

BUGS:
* Add missing valid custom role permissions [#1023](https://github.com/okta/terraform-provider-okta/pull/1023). Thanks, [@lucascantor](https://github.com/lucascantor)!
* Fix default auth server id when activate/deactivate it [#1045](https://github.com/okta/terraform-provider-okta/pull/1045). Thanks, [@peijiinsg](https://github.com/peijiinsg)!
* Panic bumper on buildEnum helper used with schemas [#1048](https://github.com/okta/terraform-provider-okta/pull/1048). Thanks, [@monde](https://github.com/monde)!

## 3.22.1 (March 11, 2022)

ENHANCEMENTS:

* Added `skip_groups` and `skip_roles` parameters to data source `okta_user` to suppress additional API calls when that data is not required. [#1011](https://github.com/okta/terraform-provider-okta/pull/1011). Thanks, [@monde](https://github.com/monde)!
* Update email temaplate names list on resource `okta_template_email`. [#1012](https://github.com/okta/terraform-provider-okta/pull/1012). Thanks, [@monde](https://github.com/monde)!

## 3.22.0 (March 03, 2022)

ENHANCEMENTS:

* Added new `okta_policy_profile_enrollment_apps` resource [#973](https://github.com/okta/terraform-provider-okta/pull/973). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Added "DYNAMIC" option to the `issuer_mode` in the `okta_auth_server` resource [#977](https://github.com/okta/terraform-provider-okta/pull/977). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Clean up provider argument conflicts documentation [#987](https://github.com/okta/terraform-provider-okta/pull/987). Thanks, [@monde](https://github.com/monde)!
* Update all App docs to match provider schema [#995](https://github.com/okta/terraform-provider-okta/pull/995). Thanks, [@virgofx](https://github.com/virgofx)!

BUGS:

* Correct ipd related error messages [#985](https://github.com/okta/terraform-provider-okta/pull/985). Thanks, [@monde](https://github.com/monde)!

## 3.21.0 (February 10, 2022)

ENHANCEMENTS:

* Added `okta_app_oauth_post_logout_redirect_uri` resource and improved request concurrency handling [#931](https://github.com/okta/terraform-provider-okta/pull/931). Thanks, [@jmaness](https://github.com/jmaness), and [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Added `LDAP` option to the `auth_provider` field in the `okta_policy_password` resource [#961](https://github.com/okta/terraform-provider-okta/pull/961). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Added new `priority` field to the `okta_auth_server_policy` data source [#965](https://github.com/okta/terraform-provider-okta/pull/965). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Added new option to the `issuer_mode` field in the `okta_app_oauth` resource [#966](https://github.com/okta/terraform-provider-okta/pull/966). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!

PROJECT IMPROVEMENTS:

* Updated docs regarding `okta_policy_rule_idp_discovery` [#964](https://github.com/okta/terraform-provider-okta/pull/964). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!

BUGS:

* Fixed import for the `okta_factor` resource [#960](https://github.com/okta/terraform-provider-okta/pull/960). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Fixed import for the `okta_policy_rule_mfa` resource [#962](https://github.com/okta/terraform-provider-okta/pull/962). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Fixed import for the `okta_group_schema_property` resource [#963](https://github.com/okta/terraform-provider-okta/pull/963). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!

## 3.20.8 (February 9, 2022)

ENHANCEMENTS:

* Removed default value for `identity_provider` field on the `okta_policy_rule_sign_on`[#955](https://github.com/okta/terraform-provider-okta/pull/955). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Added new `expire_password_on_create` field to the `okta_user` resource [#956](https://github.com/okta/terraform-provider-okta/pull/956). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Added new `user_type_id` field to the `okta_idp_oidc` and `okta_idp_saml` resources [#957](https://github.com/okta/terraform-provider-okta/pull/957). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!

## 3.20.7 (February 7, 2022)

PROJECT IMPROVEMENTS:

* Added a GH CI workflow to protect master branch [#948](https://github.com/okta/terraform-provider-okta/pull/948). Thanks, [@ymylei](https://github.com/ymylei)!

BUGS:

* Set a high limit on `client.Group.ListGroups` query data source Okta Groups [#946](https://github.com/okta/terraform-provider-okta/pull/929). Thanks, [@monde](https://github.com/monde)!

## 3.20.6 (February 3, 2022)

ENHANCEMENTS:

* Added new `identity_provider` and `identity_provider_ids` fields to the `okta_policy_rule_signon` resource [#942](https://github.com/okta/terraform-provider-okta/pull/942). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!

## 3.20.5 (February 2, 2022)

BUGS:

* Whiffed setting the user agent correctly, fixed for release.


## 3.20.4 (February 2, 2022)

ENHANCEMENTS:

* Add OIE support for MFA policies [#919](https://github.com/okta/terraform-provider-okta/pull/919). Thanks, [@virgofx](https://github.com/virgofx)!

BUGS:

* SAML SLO Cert Fix [#923](https://github.com/okta/terraform-provider-okta/pull/923). Thanks, [@ymylei](https://github.com/ymylei)!
* Nil bumper on `*sdk.ClientRateLimitMode` returned from rate limiting [#929](https://github.com/okta/terraform-provider-okta/pull/929). Thanks, [@monde](https://github.com/monde)!
* API Mutex Fix For `apps/{id}` endpoint [#933](https://github.com/okta/terraform-provider-okta/pull/933). Thanks, [@ymylei](https://github.com/ymylei)!
* Ensure okta_authenticator settings are ordered to prevent whitespace [#936](https://github.com/okta/terraform-provider-okta/pull/936). Thanks, [@virgofx](https://github.com/virgofx)!
* Ensure VERIFIED domains return true [#937](https://github.com/okta/terraform-provider-okta/pull/937). Thanks, [@virgofx](https://github.com/virgofx)!
* Fixed group search in the `okta_groups` data source [#938](https://github.com/okta/terraform-provider-okta/pull/938). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!

PROJECT IMPROVEMENTS:

* Updated dev and build tools [#912](https://github.com/okta/terraform-provider-okta/pull/912). Thanks, [@ymylei](https://github.com/ymylei)!
* Fixed TF logo [#918](https://github.com/okta/terraform-provider-okta/pull/918). Thanks, [@exitcode0](https://github.com/exitcode0)!
* Update profile mapping docs with OAuth2 scopes [#928](https://github.com/okta/terraform-provider-okta/pull/928). Thanks, [@virgofx](https://github.com/virgofx)!

## 3.20.3 (January 14, 2022)

ENHANCEMENTS:

* Added new `custom_profile_attributes` field to the `okta_group` resource [#851](https://github.com/okta/terraform-provider-okta/pull/851). Thanks, [@ymylei](https://github.com/ymylei)!
* Updated list of valid Okta OAuth scopes [#897](https://github.com/okta/terraform-provider-okta/pull/897). Thanks, [@virgofx](https://github.com/virgofx)!
* Added missing role type to the `okta_role_subscription` resource [#863](https://github.com/okta/terraform-provider-okta/pull/863). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Added new `certificate_source_type` field to the `okta_domain` resource [#899](https://github.com/okta/terraform-provider-okta/pull/899). Thanks, [@virgofx](https://github.com/virgofx)!
* Made `okta_authenticator` importable [#907](https://github.com/okta/terraform-provider-okta/pull/907). Thanks, [@virgofx](https://github.com/virgofx)!

BUGS:

* Fixed `okta_domain_verification` resource [#899](https://github.com/okta/terraform-provider-okta/pull/899). Thanks, [@virgofx](https://github.com/virgofx)!

## 3.20.2 (December 8, 2021)

ENHANCEMENTS:

* Added new `password_inline_hook` field to the `okta_user` resource [#849](https://github.com/okta/terraform-provider-okta/pull/849). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!

BUGS:

* Fixed `okta_domain` import [#845](https://github.com/okta/terraform-provider-okta/pull/845). Thanks, [quantumew](https://github.com/quantumew)!
* Fixed documentation [#848](https://github.com/okta/terraform-provider-okta/pull/848). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!

## 3.20.1 (December 3, 2021)

ENHANCEMENTS:

* Added new `apple_kid`, `apple_private_key` and `apple_team_id` fields to the `okta_idp_social` resource [#842](https://github.com/okta/terraform-provider-okta/pull/842). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Fixed docs for `okta_rate_limiting` resource [#827](https://github.com/okta/terraform-provider-okta/pull/827). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Fixed example in docs for `okta_idp_saml_key` resource [#824](https://github.com/okta/terraform-provider-okta/pull/824). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!

## 3.20.0 (November 23, 2021)

ENHANCEMENTS:

* Added new `okta_rate_limiting` resource [#803](https://github.com/okta/terraform-provider-okta/pull/803). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Added new `okta_captcha` and `okta_captcha_org_wide_settings` resources [#821](https://github.com/okta/terraform-provider-okta/pull/821). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Fixed example in docs for `okta_group` resource [#814](https://github.com/okta/terraform-provider-okta/pull/814). Thanks, [@tim-fitzgerald](https://github.com/tim-fitzgerald)!

BUGS:

* Fixed pagination bug in `okta_group_memberships` resource [#810](https://github.com/okta/terraform-provider-okta/pull/810). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Added missing fields to `okta_app_oauth` resource [#817](https://github.com/okta/terraform-provider-okta/pull/817). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!

## 3.19.0 (November 12, 2021)

ENHANCEMENTS:

* Added new `okta_admin_role_custom`, `okta_admin_role_custom_assignments` and `okta_resource_set` resources [#789](https://github.com/okta/terraform-provider-okta/pull/789). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Field `always_include_in_token` is now editable for all the default claims except `sub` [#790](https://github.com/okta/terraform-provider-okta/pull/790). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Added new `okta_link_definition` and `okta_link_value` resources [#794](https://github.com/okta/terraform-provider-okta/pull/794). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Added new `primary_factor` field to the `okta_policy_rule_signon` resource [#796](https://github.com/okta/terraform-provider-okta/pull/796). **IMPORTANT NOTE:** Available only for the organizations with Identity Engine. Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!

BUGS:

* Change authenticator status in case it's different from the state's one during resource creation [#782](https://github.com/okta/terraform-provider-okta/pull/782). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Numerus documentation fixes [#783](https://github.com/okta/terraform-provider-okta/pull/783), [#785](https://github.com/okta/terraform-provider-okta/pull/785)
and [#792](https://github.com/okta/terraform-provider-okta/pull/792). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta) and [@deepu105](https://github.com/deepu105)!
* Fixed provider crash [#795](https://github.com/okta/terraform-provider-okta/pull/795). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!

## 3.18.0 (November 2, 2021)

ENHANCEMENTS:

* Added new `okta_event_hook_verification` resource [#752](https://github.com/okta/terraform-provider-okta/pull/752). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Added new `app_include` and `app_exclude` fields to the `okta_policy_rule_mfa` resource [#762](https://github.com/okta/terraform-provider-okta/pull/762), [#771](https://github.com/okta/terraform-provider-okta/pull/771). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Added new `okta_trusted_origins` data source [#766](https://github.com/okta/terraform-provider-okta/pull/766). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Added `redirect_url` and `checkbox` fields to the `okta_app_swa` resource [#767](https://github.com/okta/terraform-provider-okta/pull/767). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Added new `user_name_template_push_status` field to some of the `okta_app_*` related resources [#769](https://github.com/okta/terraform-provider-okta/pull/769). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Added new `old_password` field to the `okta_user` resource [#765](https://github.com/okta/terraform-provider-okta/pull/765) and check for ability to change or set a password. Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!

BUGS:

* Fixed name matching for `okta_auth_server` data source [#764](https://github.com/okta/terraform-provider-okta/pull/764). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!

## 3.17.0 (October 26, 2021)

**IMPORTANT NOTE:** This release contains resources that are only available as a part of the Identity Engine. [Contact support](mailto:dev-inquiries@okta.com) for further information.

ENHANCEMENTS:

* Added new `okta_authenticator` resource and datasource [#708](https://github.com/okta/terraform-provider-okta/pull/708). Thanks, [@monde](https://github.com/monde) and [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Added new `okta_role_subscription` resource and datasource [#746](https://github.com/okta/terraform-provider-okta/pull/746). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Added new `okta_org_support` and `okta_org_configuration` resources [#749](https://github.com/okta/terraform-provider-okta/pull/749). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Added new `always_apply` field to the `okta_profile_mapping` resource [#750](https://github.com/okta/terraform-provider-okta/pull/750). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!

## 3.16.0 (October 22, 2021)

**IMPORTANT NOTE:** This release contains resources that are only available as a part of the Identity Engine. [Contact support](mailto:dev-inquiries@okta.com) for further information.

ENHANCEMENTS:

* Updated the list of supported scopes [#712](https://github.com/okta/terraform-provider-okta/pull/712). Thanks, [@boekkooi-lengoo](https://github.com/boekkooi-lengoo)!
* Added new `okta_app_signon_policy` and `okta_app_sign_on_policy_rule` resources [#714](https://github.com/okta/terraform-provider-okta/pull/714). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Added `preconfigured_app` field to the `okta_app_shared_credentials` resource [#723](https://github.com/okta/terraform-provider-okta/pull/723). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Added new `okta_network_zone` datasource [#726](https://github.com/okta/terraform-provider-okta/pull/726). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Added new `okta_security_notification_emails` and `okta_threat_insight_settings` resources [#728](https://github.com/okta/terraform-provider-okta/pull/728). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Added new `okta_policy_rule_profile_enrollment` and `okta_policy_profile_enrollment` resources [#731](https://github.com/okta/terraform-provider-okta/pull/731). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Added new `okta_auth_server_claims` and `okta_auth_server_claim` data sources [#734](https://github.com/okta/terraform-provider-okta/pull/734). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Added `disable_notifications` field to the `okta_user_admin_roles` and `okta_group_role` resources [#735](https://github.com/okta/terraform-provider-okta/pull/735). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!

BUGS:

* Fixed concurrent app logo upload [#716](https://github.com/okta/terraform-provider-okta/pull/716). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Fixed scopes diff bug [#737](https://github.com/okta/terraform-provider-okta/pull/737). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Minor tweaks to the provider's rate limiter [#719](https://github.com/okta/terraform-provider-okta/pull/719). Thanks, [@monde](https://github.com/monde) and [@phi1ipp](https://github.com/phi1ipp)!
* Made `priority` an optional parameter of `okta_app_group_assignment` [#741](https://github.com/okta/terraform-provider-okta/pull/741). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!

## 3.15.0 (October 11, 2021)

ENHANCEMENTS:

* Added new `okta_app_saml_app_settings` resource [#692](https://github.com/okta/terraform-provider-okta/pull/692). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Added new `okta_email_sender` and `okta_email_sender_verification` resources [#697](https://github.com/okta/terraform-provider-okta/pull/697). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Resource `okta_idp_saml_key` is now updatable [#698](https://github.com/okta/terraform-provider-okta/pull/698). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Added `implicit_assignment` field to the `okta_app_saml` resource [#703](https://github.com/okta/terraform-provider-okta/pull/703). Thanks, [@ashwini-desai](https://github.com/ashwini-desai)!

BUGS:

* Fixed delete operation for `okta_profile_mapping` resource [#693](https://github.com/okta/terraform-provider-okta/pull/693). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Included `404` check for `okta_app_user` resource in case app no longer exists [#695](https://github.com/okta/terraform-provider-okta/pull/695). Thanks, [@ymylei](https://github.com/ymylei)!
* Minor fix for API rate limiting [#700](https://github.com/okta/terraform-provider-okta/pull/700). Thanks, [@monde](https://github.com/monde) and [@phi1ipp](https://github.com/phi1ipp)!
* Fixed schema-related resources to handle numeric arrays properly [#702](https://github.com/okta/terraform-provider-okta/pull/702). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!

## 3.14.0 (October 7, 2021)

ENHANCEMENTS:

* Added new `okta_domain_verification` and `okta_domain_certificate` resources [#687](https://github.com/okta/terraform-provider-okta/pull/687). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Added new `okta_group_schema_property` resource [#688](https://github.com/okta/terraform-provider-okta/pull/688). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Added `skip_users` and `skip_groups` fields to the app-related data sources [#677](https://github.com/okta/terraform-provider-okta/pull/677). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta) and [@Philipp](https://github.com/phi1ipp)!
* Added new grant type values to the `okta_app_oauth` and `okta_auth_server_policy_rule` resources [#691](https://github.com/okta/terraform-provider-okta/pull/691). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!

BUGS:
* `okta_app_oauth.groups_claim` field won't be requested if it's not set in the config [#668](https://github.com/okta/terraform-provider-okta/pull/668). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Fixed panic in `okta_auth_server` data source [#679](https://github.com/okta/terraform-provider-okta/pull/679). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Fixed false positive `profile` field set in `okta_app_group_assignments` resource [#689](https://github.com/okta/terraform-provider-okta/pull/689). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!

## 3.13.13 (September 23, 2021)

BUGS:
* Another attempt to fix constant change-loops in the `okta_app_group_assignments` resource [#664](https://github.com/okta/terraform-provider-okta/pull/664). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!

## 3.13.12 (September 22, 2021)

BUGS:
* Fixed false users sync for `okta_group` resource [#661](https://github.com/okta/terraform-provider-okta/pull/661). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!

## 3.13.11 (September 21, 2021)

ENHANCEMENTS:
* Added `skip_users` to the `okta_group` resource (check latest documentation for the usage of these fields) [#646](https://github.com/okta/terraform-provider-okta/pull/646). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Added new `users_excluded` field to the `okta_group_rule` resource [#651](https://github.com/okta/terraform-provider-okta/pull/651). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!

BUGS:
* Fixed constant change-loops in the `okta_app_group_assignments` resource [#644](https://github.com/okta/terraform-provider-okta/pull/644). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Fixed typo and deprecation warning in the documentation for `okta_app_user` resource [#645](https://github.com/okta/terraform-provider-okta/pull/645). Thanks, [@SaffatHasan](https://github.com/SaffatHasan)!
* Fixed `okta_group_role` resource update in case of several roles are being updated [#646](https://github.com/okta/terraform-provider-okta/pull/646). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Terraform will attempt to remove `okta_user_schema_property` resource several times in case the resource still exists in the organization [#656](https://github.com/okta/terraform-provider-okta/pull/656). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!

## 3.13.10 (September 13, 2021)

BUGS:
* Fixed the way `okta_policy_mfa` resource store its factors in the state [#641](https://github.com/okta/terraform-provider-okta/pull/641). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Fixed provider crash when using policy rules resources [#641](https://github.com/okta/terraform-provider-okta/pull/641). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!

## 3.13.9 (September 10, 2021)

ENHANCEMENTS:
* Added `app_settings_json` to the `okta_app_oauth` resource [#627](https://github.com/okta/terraform-provider-okta/pull/627). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Added `skip_users` and `skip_groups` to the `okta_app_*` resources (check latest documentation for the usage of these fields) [#633](https://github.com/okta/terraform-provider-okta/pull/633). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!

BUGS:
* Fixed resource import of the `okta_app_group_assignments` [#630](https://github.com/okta/terraform-provider-okta/pull/630). Thanks, [@Philipp](https://github.com/phi1ipp)!
* Fixed creation of multiple app user schema properties for new (recently created) apps. [#634](https://github.com/okta/terraform-provider-okta/pull/634). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Fixed description for the app logo field [#639](https://github.com/okta/terraform-provider-okta/pull/639). Thanks, [@sklarsa](https://github.com/sklarsa)!

## 3.13.8 (September 1, 2021)

ENHANCEMENTS:
* Add `credentials_scheme`, `reveal_password`, `shared_username` and `shared_password` to the `okta_app_three_field` resource [#619](https://github.com/okta/terraform-provider-okta/pull/619). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Add `password_hash` to the `okta_user` resource [#622](https://github.com/okta/terraform-provider-okta/pull/622). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!

BUGS:
* Fix import of `accessibility_login_redirect_url` field in the `okta_app_saml` resource [#613](https://github.com/okta/terraform-provider-okta/pull/613). Thanks, [@Philipp](https://github.com/phi1ipp)!
* Fix create/update operations for the `okta_app_user_custom_schema_property` resource [#606](https://github.com/okta/terraform-provider-okta/pull/606). Thanks, [@Philipp](https://github.com/phi1ipp)!
* Fix provider crash when importing `okta_app_oauth` resource [#616](https://github.com/okta/terraform-provider-okta/pull/616). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Fix `group_memberships` field setup for `okta_user` data source [#615](https://github.com/okta/terraform-provider-okta/pull/615). Thanks, [@BrentSouza](https://github.com/BrentSouza)!
* Fix provider crash when `okta_policy_rule_idp_discovery` does not exist [#622](https://github.com/okta/terraform-provider-okta/pull/622). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!

## 3.13.7 (Aug 23, 2021)

ENHANCEMENTS:
* Add `asns` field to the `okta_network_zone` resource [#599](https://github.com/okta/terraform-provider-okta/pull/599). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Add `app_links_json` to the `okta_app_saml` resource [#601](https://github.com/okta/terraform-provider-okta/pull/601). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Add `app_settings_json` to the `okta_app_auto_login` resource [#602](https://github.com/okta/terraform-provider-okta/pull/602). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!

BUGS:
* Fix `*_token_*` fields setup when importing `okta_auth_server_policy_rule` resource [#600](https://github.com/okta/terraform-provider-okta/pull/600). Thanks, [@Philipp](https://github.com/phi1ipp)!
* Governed Transport is now handling nil response in `postRequestHook` func [#603](https://github.com/okta/terraform-provider-okta/pull/603). Thanks, [@Mike](https://github.com/monde)!

## 3.13.6 (Aug 18, 2021)

ENHANCEMENTS:
* Add `saml_version` field to the `okta_app_saml` resource [#593](https://github.com/okta/terraform-provider-okta/pull/593). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!

BUGS:
* Fixed provider crash when using `okta_template_sms` without `translations` [#592](https://github.com/okta/terraform-provider-okta/pull/592). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!

## 3.13.5 (Aug 17, 2021)

ENHANCEMENTS:
* Add `admin_note` and `enduser_note` to all `okta_app_*` resources [#589](https://github.com/okta/terraform-provider-okta/pull/589). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!

BUGS:
* Fixed bug in config validator [#589](https://github.com/okta/terraform-provider-okta/pull/589). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!

## 3.13.4 (Aug 16, 2021)

ENHANCEMENTS:
* Add auth config validator [#567](https://github.com/okta/terraform-provider-okta/pull/579). Thanks, [@bendrucker](https://github.com/bendrucker)!

BUGS:
* Fix unmarshalling error for `okta_network_zone` resource [#586](https://github.com/okta/terraform-provider-okta/pull/586). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Fix `pattern` property setup in `okta_user_schema_property` [#583](https://github.com/okta/terraform-provider-okta/pull/583). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!

## 3.13.3 (Aug 12, 2021)

BUGS:
* Fix `OKTA_API_SCOPES` not being set via env variable [#574](https://github.com/okta/terraform-provider-okta/pull/574). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!

## 3.13.2 (Aug 12, 2021)

ENHANCEMENTS:
* Minor tweaks for the API governor [#569](https://github.com/okta/terraform-provider-okta/pull/569). Thanks, [@monde](https://github.com/monde)!
* Use more methods from official Okta Golang SDK [#567](https://github.com/okta/terraform-provider-okta/pull/567). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Provider will now terminate in case of invalid credentials [#571](https://github.com/okta/terraform-provider-okta/pull/571). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!

BUGS:
* Fix `OKTA_API_SCOPES` env var parsing [#570](https://github.com/okta/terraform-provider-okta/pull/570). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Fix `target_app_list` and `target_group_list` fields behavior in `okta_group_role` resource [#570](https://github.com/okta/terraform-provider-okta/pull/570). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!

## 3.13.1 (Aug 6, 2021)

ENHANCEMENTS:

* Add `inline_hook_id` field to the `okta_app_saml` resource [#561](https://github.com/okta/terraform-provider-okta/pull/561). Thanks, [@noinarisak](https://github.com/noinarisak)!
* Add experimental `max_api_capacity` configuration field to the provider. Thanks, [@monde](https://github.com/monde)!

BUGS:

* Fixed users and groups assignment for the application resources [#565](https://github.com/okta/terraform-provider-okta/pull/565). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!

## 3.13.0 (Jul 29, 2021)

ENHANCEMENTS:

* Add new `user_factor_question` resource [#551](https://github.com/okta/terraform-provider-okta/pull/551). Thanks, [@pengyuwang-okta](https://github.com/pengyuwang-okta)!
* Add new `okta_behavior` resource [#552](https://github.com/okta/terraform-provider-okta/pull/552). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Add new `okta_user_security_questions` data source [#552](https://github.com/okta/terraform-provider-okta/pull/552). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!

## 3.12.1 (Jul 24, 2021)

BUGS:

* Fix provider crash caused by the `okta_policy_rule_signon` resource [#543](https://github.com/okta/terraform-provider-okta/pull/543). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Fix permissions field set behaviour in o`kta_app_user_schema_property` resource [#543](https://github.com/okta/terraform-provider-okta/pull/543). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Reverted the changes regarding the users field in the `okta_group` resource that was introducing breaking change [#543](https://github.com/okta/terraform-provider-okta/pull/543). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!

## 3.12.0 (Jul 20, 2021)

ENHANCEMENTS:

* Add new `okta_app_group_assignments` resource [#401](https://github.com/okta/terraform-provider-okta/pull/401). Thanks, [@edulop91](https://github.com/edulop91)!
* Add new `okta_user_group_memberships` resource [#416](https://github.com/okta/terraform-provider-okta/pull/416). Thanks, [@ymylei](https://github.com/ymylei)!
* Add `logo` and `logo_url` fields to all the `okta_app_*` related resources [#423](https://github.com/okta/terraform-provider-okta/pull/423) and [#514](https://github.com/okta/terraform-provider-okta/pull/514). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta) and [@gavinbunney](https://github.com/gavinbunney) for the fix!
* Add new `okta_group_memberships` resource [#427](https://github.com/okta/terraform-provider-okta/pull/427). Thanks, [@ymylei](https://github.com/ymylei)!
* Add `display_name` field to the `okta_auth_server_scope` resource [#433](https://github.com/okta/terraform-provider-okta/pull/433). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Add new `okta_app_shared_credentials` resource [#446](https://github.com/okta/terraform-provider-okta/pull/446). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Add `groups_claim` field to the `okta_app_oauth` resource [#468](https://github.com/okta/terraform-provider-okta/pull/468). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Add `wildcard_redirect` field to the `okta_app_oauth` resource [#474](https://github.com/okta/terraform-provider-okta/pull/474). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Add new `okta_app_group_assignments` data source [#498](https://github.com/okta/terraform-provider-okta/pull/498). Thanks, [@ymylei](https://github.com/ymylei)!
* Add new `okta_app_user_assignments` data source [#501](https://github.com/okta/terraform-provider-okta/pull/501). Thanks, [@ymylei](https://github.com/ymylei)!
* Add new `okta_user_admin_roles` resource [#518](https://github.com/okta/terraform-provider-okta/pull/518). Thanks, [@gavinbunney](https://github.com/gavinbunney)!
* Add new `okta_factor_totp` resource [#519](https://github.com/okta/terraform-provider-okta/pull/519). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Add `dynamic_proxy_type` field to the `okta_network_zone` resource [#522](https://github.com/okta/terraform-provider-okta/pull/522). [@gavinbunney](https://github.com/gavinbunney)!
* Add `issuer_mode` field to the `okta_auth_server_default` resource [#524](https://github.com/okta/terraform-provider-okta/pull/524). [@gavinbunney](https://github.com/gavinbunney)!
* Add `risc_level`, `behaviors` and `factor_sequence` fields to the `okta_policy_rule_signon` resource [#526](https://github.com/okta/terraform-provider-okta/pull/526). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Add new `okta_behavior` data source [#526](https://github.com/okta/terraform-provider-okta/pull/526). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Add new `okta_domain` resource [#530](https://github.com/okta/terraform-provider-okta/pull/530). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!

BUGS:

* Suppress 404 in case group role was removed outside of the terraform [#417](https://github.com/okta/terraform-provider-okta/pull/417). Thanks, [@ymylei](https://github.com/ymylei)!
* Don't recreate `okta_user` resource in case `login` field is changed [#435](https://github.com/okta/terraform-provider-okta/pull/435/files). Thanks, [@ymylei](https://github.com/ymylei)!
* Fixed attribute statements setup for preconfigured apps [#439](https://github.com/okta/terraform-provider-okta/pull/439). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!
* Don't recreate schema related resources in case `array_enum`, `array_one_of`, `enum` or `one_of` have changed [@531](https://github.com/okta/terraform-provider-okta/pull/531/files). Thanks, [@bogdanprodan-okta](https://github.com/bogdanprodan-okta)!

#### Special thanks to [@JeffAshton](https://github.com/JeffAshton), [@jeffg-hpe](https://github.com/jeffg-hpe), [@jtdoepke](https://github.com/jtdoepke), [@thatguysimon](https://github.com/thatguysimon), [@ymylei](https://github.com/ymylei), [@joshowen](https://github.com/joshowen), [@AlexanderProschek](https://github.com/AlexanderProschek), [@gavinbunney](https://github.com/gavinbunney) for a lot of various documentation fixes and code improvements!!!


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
