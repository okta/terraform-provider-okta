---
page_title: "Resource: okta_app_signon_policy_rules"
description: |-
  Manages multiple app sign-on policy rules for a single policy. This resource allows you to define all rules for a policy in a single configuration block, ensuring consistent priority ordering and avoiding drift issues.
---
# Resource: okta_app_signon_policy_rules

Manages multiple app sign-on policy rules for a single policy. This resource allows you to define all rules for a policy in a single configuration block, ensuring consistent priority ordering and avoiding drift issues.

~> **IMPORTANT:** This resource uses name-first matching to identify and update rules. When migrating from individual `okta_app_signon_policy_rule` resources, ensure rule names remain consistent to enable safe adoption without data loss.

~> **NOTE ON RENAMING RULES:** If you rename a rule without explicitly preserving its `id`, the provider will treat it as a deletion of the old rule and creation of a new rule. To rename a rule while preserving its configuration and ID, you must explicitly set the `id` attribute in your configuration before changing the `name`. For example:
```terraform
rule {
  id   = "rulAbc123" # Explicitly reference the existing rule ID
  name = "New Rule Name" # New name
  # ... other attributes
}
```
After applying with the explicit `id`, you can remove it in subsequent applies and the rule will be tracked by its new name.

## Example Usage

```terraform
resource "okta_app_signon_policy_rules" "example" {
  policy_id = okta_app_signon_policy.example.id

  rule {
    name                        = "High Priority Rule"
    priority                    = 1
    factor_mode                 = "2FA"
    re_authentication_frequency = "PT2H"
    status                      = "ACTIVE"
  }

  rule {
    name       = "Low Priority Rule"
    priority   = 2
    factor_mode = "1FA"
    access     = "ALLOW"
    status     = "ACTIVE"
  }

  rule {
    name   = "Deny Rule"
    priority = 3
    access = "DENY"
    status = "ACTIVE"
  }
}
```

## Argument Reference

### Required Arguments

- `policy_id` (String) - ID of the policy to manage rules for. Changing this forces a new resource.

### Optional Arguments

- `rule` (Block Set) - List of policy rules to manage. Each rule block supports the following:
  - `id` (String) - ID of the rule (computed, but can be specified to adopt an existing rule during migration).
  - `name` (String, Required) - Policy rule name. Must be unique within the policy.
  - `priority` (Number) - Priority of the rule. Lower numbers are evaluated first. When omitted, the provider will use the actual priority assigned by Okta.
  - `status` (String) - Status of the rule: `ACTIVE` or `INACTIVE`. Defaults to `ACTIVE`.
  - `access` (String) - Access decision: `ALLOW` or `DENY`. Defaults to `ALLOW`.
  - `factor_mode` (String) - Number of factors required: `1FA` or `2FA`. Defaults to `2FA`.
  - `type` (String) - Verification method type. Defaults to `ASSURANCE`.
  - `re_authentication_frequency` (String) - Re-authentication frequency in ISO 8601 duration format (e.g., `PT2H`). Defaults to `PT2H`.
  - `inactivity_period` (String) - Inactivity period before re-authentication in ISO 8601 duration format.
  - `network_connection` (String) - Network selection mode: `ANYWHERE`, `ZONE`, `ON_NETWORK`, or `OFF_NETWORK`. Defaults to `ANYWHERE`.
  - `network_includes` (List of Strings) - List of network zone IDs to include.
  - `network_excludes` (List of Strings) - List of network zone IDs to exclude.
  - `groups_included` (Set of Strings) - Set of group IDs to include in this rule.
  - `groups_excluded` (Set of Strings) - Set of group IDs to exclude from this rule.
  - `users_included` (Set of Strings) - Set of user IDs to include in this rule.
  - `users_excluded` (Set of Strings) - Set of user IDs to exclude from this rule.
  - `user_types_included` (Set of Strings) - Set of user type IDs to include.
  - `user_types_excluded` (Set of Strings) - Set of user type IDs to exclude.
  - `device_is_registered` (Boolean) - Require device to be registered with Okta Verify.
  - `device_is_managed` (Boolean) - Require device to be managed by a device management system.
  - `device_assurances_included` (Set of Strings) - Set of device assurance policy IDs to include.
  - `custom_expression` (String) - Custom Okta Expression Language condition for advanced matching.
  - `risk_score` (String) - Risk score level to match: `ANY`, `LOW`, `MEDIUM`, or `HIGH`. Defaults to `ANY`.
  - `constraints` (List of Strings) - List of authenticator constraints as JSON-encoded strings.
  - `platform_include` (Block Set) - Platform conditions to include. Each platform block supports:
    - `type` (String) - Platform type: `ANY`, `MOBILE`, or `DESKTOP`.
    - `os_type` (String) - OS type: `ANY`, `IOS`, `ANDROID`, `WINDOWS`, `OSX`, `MACOS`, `CHROMEOS`, or `OTHER`.
    - `os_expression` (String) - Custom OS expression for advanced matching.
  - `system` (Boolean, Computed) - Whether this is a system rule (e.g., Catch-all Rule). System rules cannot be modified.
  - `chains` (Block List) Authentication method chains. Only supports 5 items in the array. Each chain can support maximum 3 steps. To be used only with verification method type `AUTH_METHOD_CHAIN`.(see [below for nested schema](#nestedblock--chains))


<a id="nestedblock--chains"></a>
### Nested Schema for `chains`

Optional:

- `authenticationMethods` (List of Authentication Methods) (see [below for nested schema](#nestedblock--authenticationMethods))
- `next` (List) The next steps of the authentication method chain. This is an array of type `chains`. Only supports one item in the array.
- `reauthenticateIn` (String) Specifies how often the user is prompted for authentication using duration format for the time period. This parameter can't be set at the same time as the `re_authentication_frequency` field.

<a id="nestedblock--authenticationMethods"></a>
### Nested Schema for `authenticationMethods`
Required:
- `key` (String) A label that identifies the authenticator.
- `method` (String) Specifies the method used for the authenticator.

## Attributes Reference

- `id` (String) - The ID of this resource (same as `policy_id`).

## Priority Management

~> **IMPORTANT:** When managing priorities, follow Okta's top-down synchronization strategy to avoid automatic priority shifting. Always update rules in priority order (P1 → PN) to prevent cascading shifts.

### Priority Behavior

- Priorities are **0-indexed** but can be any positive integer.
- When you omit `priority` in config, Terraform will display the actual priority assigned by Okta in the state.
- Okta allows gaps in priorities (e.g., priorities 1, 2, 5, 99 are valid).
- **Maximum of 99 rules per policy:** You can define up to 99 custom rules. Priority 99 is reserved for the system Catch-all Rule.
- If you update multiple rules simultaneously and priorities conflict, Okta may automatically shift priorities to resolve conflicts. To prevent this:
  1. Use explicit, sequential priority values (1, 2, 3, 4, 5).
  2. Avoid updating rules with the same priority.
  3. Consider using temporary high priorities (100+) if reordering is needed.

### System Rules and Catch-all Rule

- **Catch-all Rule:** Every policy includes a system-level Catch-all Rule with priority 99. This is the default rule that applies if no other rules match.
- **Immutable:** System rules (like "Catch-all Rule") cannot be modified, deleted, or have their conditions/priorities changed.
- **Import behavior:** If you import a policy with a Catch-all Rule, it will appear in your state but cannot be managed by Terraform. Omit it from your configuration.

## Migration from okta_app_signon_policy_rule

If you currently use multiple individual `okta_app_signon_policy_rule` resources, follow these steps to safely migrate to the grouped resource:

### Step 1: Backup Your State

```bash
terraform state pull > terraform-state-backup.json
```

### Step 2: Identify Rule IDs and Names

List all rules in your policy to identify their IDs and names:

```bash
# Using Terraform state
terraform state show 'okta_app_signon_policy_rule.rule1'
```


### Step 3: Comment the existing configuration and remove the current config from terraform state without deleting the rules on the Okta org

```terraform
removed {
  from = okta_app_signon_policy_rule.rule1
  lifecycle {
    destroy = false
  }
}
removed {
  from = okta_app_signon_policy_rule.rule2
  lifecycle {
    destroy = false
  }
}
```

### Step 4: Create the Grouped Resource Configuration

Create a new `okta_app_signon_policy_rules` resource with the same rules. Use **rule names as the matching key** (do NOT include `id` unless adopting an existing rule):

```terraform
resource "okta_app_signon_policy_rules" "migrated" {
  policy_id = okta_app_signon_policy.my_policy.id

  rule {
    name                        = "High Priority Rule"
    priority                    = 1  # Explicitly set to avoid shifting
    factor_mode                 = "2FA"
    re_authentication_frequency = "PT2H"
    status                      = "ACTIVE"
  }

  rule {
    name       = "Second Rule"
    priority   = 2
    factor_mode = "1FA"
    status     = "ACTIVE"
  }

  rule {
    name   = "Deny Rule"
    priority = 3
    access = "DENY"
    status = "ACTIVE"
  }
}
```

### Step 5: Plan and Verify

Run `terraform plan` to see if any changes will be made:

```bash
terraform plan -out migration.tfplan
```

**Expected outcome:** Either "No changes" or only legitimate config updates. If the plan shows unexpected deletes or recreates:
- Verify rule names match exactly (case-sensitive).
- Check that attributes (access, factor_mode, etc.) match Okta's current state.
- Ensure you're not declaring system rules (e.g., Catch-all Rule).

### Step 6: Apply the Migration

```bash
terraform apply migration.tfplan
```

### Step 6: Verify Final State

```bash
terraform plan
# Output should show: No changes. Your infrastructure matches the configuration.
```

## Important Migration Notes

### Priority Differences

When migrating from `okta_app_signon_policy_rule` (single-rule resource) to `okta_app_signon_policy_rules` (grouped resource):

- **Old behavior:** Single-rule resource defaults to `priority = 0` in Terraform state, even if Okta assigned a different priority.
- **New behavior:** Grouped resource reflects the actual priority assigned by Okta. When you omit `priority`, the state will show the real Okta-assigned priority.

**Action:** Explicitly set `priority` values in your configuration to avoid unexpected changes during migration.

### Name-First Matching

The grouped resource matches rules by **name first**, not by ID:

- If you rename a rule in config, the provider will match it to the existing rule by the new name.
- If a name doesn't exist in Okta, a new rule is created.
- This design prevents accidental re-matching when state order differs from config order.

### System Rules

System rules (e.g., "Catch-all Rule") cannot be managed:

- If your policy has a system rule, do NOT declare it in your configuration.
- The provider will detect it during import/read but will not attempt to modify or delete it.
- System rules are read-only from Terraform's perspective.

## Troubleshooting

### Plan shows "1 to change" but config matches Okta

**Cause:** State and config rule order differ, triggering re-matching.

**Solution:** Verify rule names match exactly (case-sensitive) and attributes are correct. Run `terraform plan -refresh=true` to refresh state from Okta.

### Priorities shift unexpectedly after apply

**Cause:** Okta automatically shifts priorities when updates conflict. This is expected behavior per Okta's prioritization rules.

**Solution:** Use explicit, sequential priorities (1, 2, 3, ...) instead of 0-indexed or sparse values. Avoid updating multiple rules simultaneously with overlapping priorities.

### Error: "System rule cannot be modified"

**Cause:** You attempted to set conditions on a system rule (Catch-all Rule).

**Solution:** Remove the system rule from your configuration. System rules are immutable.

## Import

Import an existing policy's rules:

```bash
terraform import okta_app_signon_policy_rules.example <policy_id>
```

This will populate your state with all rules for the policy (including system rules). Update your configuration to match.

## See Also

- [okta_app_signon_policy](https://registry.terraform.io/providers/okta/okta/latest/docs/resources/app_signon_policy) - Create and manage app sign-on policies.
- [Okta Policy Rule Prioritization Guide](https://developer.okta.com/docs/guides/policy-rule-prioritization/main/) - Learn about priority management and top-down synchronization.
