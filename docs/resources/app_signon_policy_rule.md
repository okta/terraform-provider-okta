---
page_title: "Resource: okta_app_signon_policy_rule"
description: |-
  
---

# Resource: okta_app_signon_policy_rule





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Policy Rule Name
- `policy_id` (String) ID of the policy

### Optional

- `access` (String) Allow or deny access based on the rule conditions: ALLOW or DENY
- `constraints` (List of String) An array that contains nested Authenticator Constraint objects that are organized by the Authenticator class
- `custom_expression` (String) This is an optional advanced setting. If the expression is formatted incorrectly or conflicts with conditions set above, the rule may not match any users.
- `device_assurances_included` (Set of String) List of device assurance IDs to include
- `device_is_managed` (Boolean) If the device is managed. A device is managed if it's managed by a device management system. When managed is passed, registered must also be included and must be set to true.
- `device_is_registered` (Boolean) If the device is registered. A device is registered if the User enrolls with Okta Verify that is installed on the device.
- `factor_mode` (String) The number of factors required to satisfy this assurance level
- `groups_excluded` (Set of String) List of group IDs to exclude
- `groups_included` (Set of String) List of group IDs to include
- `inactivity_period` (String) The inactivity duration after which the end user must re-authenticate. Use the ISO 8601 Period format for recurring time intervals.
- `network_connection` (String) Network selection mode: ANYWHERE, ZONE, ON_NETWORK, or OFF_NETWORK.
- `network_excludes` (List of String) The zones to exclude
- `network_includes` (List of String) The zones to include
- `platform_include` (Block Set) (see [below for nested schema](#nestedblock--platform_include))
- `priority` (Number) Priority of the rule.
- `re_authentication_frequency` (String) The duration after which the end user must re-authenticate, regardless of user activity. Use the ISO 8601 Period format for recurring time intervals. PT0S - Every sign-in attempt, PT43800H - Once per session
- `risk_score` (String) The risk score specifies a particular level of risk to match on: ANY, LOW, MEDIUM, HIGH
- `status` (String) Status of the rule
- `type` (String) The Verification Method type
- `user_types_excluded` (Set of String) Set of User Type IDs to exclude
- `user_types_included` (Set of String) Set of User Type IDs to include
- `users_excluded` (Set of String) Set of User IDs to exclude
- `users_included` (Set of String) Set of User IDs to include

### Read-Only

- `id` (String) The ID of this resource.
- `system` (Boolean) Often the "Catch-all Rule" this rule is the system (default) rule for its associated policy

<a id="nestedblock--platform_include"></a>
### Nested Schema for `platform_include`

Optional:

- `os_expression` (String) Only available with OTHER OS type
- `os_type` (String)
- `type` (String)


