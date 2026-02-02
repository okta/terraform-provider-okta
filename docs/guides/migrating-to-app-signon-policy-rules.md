# Migrating to `okta_app_signon_policy_rules` Resource

This guide explains how to migrate from using multiple individual `okta_app_signon_policy_rule` resources to the new consolidated `okta_app_signon_policy_rules` resource.

## Why Migrate?

The individual `okta_app_signon_policy_rule` resource has a known issue with **priority drift**. When managing multiple rules for the same policy, Terraform may send API requests in an unpredictable order, causing Okta to auto-assign priorities differently than specified. This results in perpetual diffs in your Terraform plans.

The new `okta_app_signon_policy_rules` resource solves this by:
- Managing all rules for a policy as a single atomic unit
- Processing rules in priority order (lowest priority number first)
- Ensuring consistent state after each apply
- Handling name conflicts during rule updates automatically

## Prerequisites

Before migrating, ensure you have:
1. Terraform 1.0 or later
2. Okta Terraform Provider version X.X.X or later (with `okta_app_signon_policy_rules` support)
3. A backup of your current Terraform state
4. Knowledge of all existing rules and their configurations

## Migration Steps

### Step 1: Document Your Existing Rules

First, list all your existing `okta_app_signon_policy_rule` resources for the policy you want to migrate. Note down:
- Rule names
- Priorities
- All configuration settings (constraints, conditions, etc.)

Example of existing configuration:
```hcl
resource "okta_app_signon_policy_rule" "rule1" {
  name               = "MFA-Required"
  policy_id          = okta_app_signon_policy.my_policy.id
  priority           = 1
  status             = "ACTIVE"
  factor_mode        = "2FA"
  network_connection = "ANYWHERE"
  constraints = [jsonencode({
    "authenticationMethods" : [
      {
        "key" : "okta_verify",
        "method" : "signed_nonce"
      }
    ]
  })]
}

resource "okta_app_signon_policy_rule" "rule2" {
  name               = "Password-Only"
  policy_id          = okta_app_signon_policy.my_policy.id
  priority           = 2
  status             = "ACTIVE"
  factor_mode        = "1FA"
  network_connection = "ANYWHERE"
  constraints = [jsonencode({
    "authenticationMethods" : [
      {
        "key" : "okta_password",
        "method" : "password"
      }
    ]
  })]
}
```

### Step 2: Remove Old Resources from State (NOT from Okta)

**IMPORTANT**: This step removes resources from Terraform state only. It does NOT delete the actual rules in Okta.

```bash
# Remove each individual rule resource from state
terraform state rm okta_app_signon_policy_rule.rule1
terraform state rm okta_app_signon_policy_rule.rule2
# ... repeat for all rules
```

### Step 3: Comment Out or Remove Old Resource Definitions

In your Terraform configuration, comment out or remove the old `okta_app_signon_policy_rule` resources:

```hcl
# COMMENTED OUT - Migrated to okta_app_signon_policy_rules
# resource "okta_app_signon_policy_rule" "rule1" {
#   ...
# }
# resource "okta_app_signon_policy_rule" "rule2" {
#   ...
# }
```

### Step 4: Create the New Consolidated Resource

Create a new `okta_app_signon_policy_rules` resource with all your rules:

```hcl
resource "okta_app_signon_policy_rules" "my_policy_rules" {
  policy_id = okta_app_signon_policy.my_policy.id

  rule {
    name               = "MFA-Required"
    priority           = 1
    status             = "ACTIVE"
    factor_mode        = "2FA"
    network_connection = "ANYWHERE"
    constraints = [jsonencode({
      "authenticationMethods" : [
        {
          "key" : "okta_verify",
          "method" : "signed_nonce"
        }
      ]
    })]
  }

  rule {
    name               = "Password-Only"
    priority           = 2
    status             = "ACTIVE"
    factor_mode        = "1FA"
    network_connection = "ANYWHERE"
    constraints = [jsonencode({
      "authenticationMethods" : [
        {
          "key" : "okta_password",
          "method" : "password"
        }
      ]
    })]
  }
}
```

### Step 5: Import Existing Rules

Import the existing rules from Okta into the new resource:

```bash
terraform import okta_app_signon_policy_rules.my_policy_rules <policy_id>
```

Replace `<policy_id>` with your actual policy ID (e.g., `rstuibbxieCTOq8vq1d7`).

### Step 6: Run Terraform Plan

After import, run a plan to see what changes Terraform detects:

```bash
terraform plan
```

## Understanding the Plan Output After Import

### Expected Behavior

After import, you may see a plan that appears to show "name changes" on rules. **This is a display artifact**, not an actual rename operation.

Example plan output:
```
~ rule {
    ~ name = "Rule2" -> "Rule1"   # This is a POSITIONAL diff, not a real rename
    ...
  }
~ rule {
    ~ name = "Rule1" -> "Rule2"   # This is a POSITIONAL diff, not a real rename
    ...
  }
```

**Why this happens**: 
- Import stores rules sorted by priority (or API order)
- Your config may have rules in a different order
- Terraform compares lists by POSITION, not by name
- The plan shows what's different at each position

**What actually happens during apply**:
- The provider matches rules by NAME, not position
- Each rule is updated with the config that matches its name
- No rules are actually renamed
- After apply, state order matches your config order

### When to Be Concerned

You should review carefully if you see:
1. **Constraints changing** between rules (e.g., MFA constraints appearing on a password-only rule)
2. **Groups or users changing** between rules
3. **Actual rule deletions** (rules being removed entirely)

If you see configuration CONTENT moving between rules (not just positional name swaps), verify that your config rule names match exactly what exists in Okta.

## Important Considerations

### 1. Rule Names Must Be Unique

Within a policy, each rule must have a unique name. The provider uses names to match configuration to existing rules.

### 2. System Rules Cannot Be Managed

The "Catch-all Rule" (system rule) is automatically excluded during import and cannot be managed by this resource. It will always exist with priority 99.

### 3. Priority Gaps Are Allowed

Okta ACCESS_POLICY rules allow gaps in priorities. You can have rules with priorities 1, 3, 5 without needing 2 and 4.

### 4. Order Your Config by Priority (Recommended)

For clearest plan output, order your `rule` blocks by priority:

```hcl
resource "okta_app_signon_policy_rules" "my_policy_rules" {
  policy_id = okta_app_signon_policy.my_policy.id

  rule {
    name     = "Highest-Priority-Rule"
    priority = 1
    # ...
  }

  rule {
    name     = "Second-Priority-Rule"
    priority = 2
    # ...
  }

  rule {
    name     = "Third-Priority-Rule"
    priority = 3
    # ...
  }
}
```

### 5. Name Changes Require Temporary Names

If you need to swap rule names (e.g., rename "RuleA" to "RuleB" and "RuleB" to "RuleA"), the provider automatically handles this by:
1. Renaming the conflicting rule to a temporary name
2. Applying the first rename
3. Applying the second rename

This happens transparently during apply.

### 6. Deleting Rules

To delete a rule, simply remove its `rule` block from the configuration. The provider will delete rules that are in state but not in the plan.

### 7. Adding New Rules

To add a new rule, add a new `rule` block with a unique name. The provider will create it with the specified priority.

## Rollback Procedure

If you need to rollback to the old individual resources:

1. Remove the new resource from state:
   ```bash
   terraform state rm okta_app_signon_policy_rules.my_policy_rules
   ```

2. Uncomment your old resource definitions

3. Import each rule individually:
   ```bash
   terraform import okta_app_signon_policy_rule.rule1 <policy_id>/<rule_id>
   terraform import okta_app_signon_policy_rule.rule2 <policy_id>/<rule_id>
   ```

## Troubleshooting

### "Policy rule name already in use" Error

This occurs when trying to rename a rule to a name that's already taken by another rule. The provider should handle this automatically with temporary names, but if you see this error:

1. Ensure you're using the latest provider version
2. Check that your rule names in config match what's in Okta
3. Try running apply again (the retry logic should handle transient conflicts)

### "Provider produced inconsistent result after apply" Error

This typically occurs when:
1. Rule IDs in the plan don't match what's returned after apply
2. Usually caused by mismatched rule ordering between import and config

**Solution**: 
1. Remove the resource from state: `terraform state rm okta_app_signon_policy_rules.my_policy_rules`
2. Re-import: `terraform import okta_app_signon_policy_rules.my_policy_rules <policy_id>`
3. Run apply again

### Empty Sets Showing as Changes ([] -> null)

After import, you may see changes like:
```
- groups_excluded = [] -> null
- users_included  = [] -> null
```

This is expected and harmless. It occurs because the imported state has empty sets, but your config doesn't specify these optional fields. After the first apply, subsequent plans will show no changes.

## Example: Complete Migration

### Before (Old Way)
```hcl
resource "okta_app_signon_policy" "my_app_policy" {
  name        = "My App Policy"
  description = "Policy for my application"
}

resource "okta_app_signon_policy_rule" "mfa_rule" {
  name               = "Require-MFA"
  policy_id          = okta_app_signon_policy.my_app_policy.id
  priority           = 1
  factor_mode        = "2FA"
  constraints = [jsonencode({
    "authenticationMethods": [{"key": "okta_verify", "method": "signed_nonce"}]
  })]
}

resource "okta_app_signon_policy_rule" "password_rule" {
  name               = "Password-Fallback"
  policy_id          = okta_app_signon_policy.my_app_policy.id
  priority           = 2
  factor_mode        = "1FA"
  constraints = [jsonencode({
    "authenticationMethods": [{"key": "okta_password", "method": "password"}]
  })]
}
```

### After (New Way)
```hcl
resource "okta_app_signon_policy" "my_app_policy" {
  name        = "My App Policy"
  description = "Policy for my application"
}

resource "okta_app_signon_policy_rules" "my_app_policy_rules" {
  policy_id = okta_app_signon_policy.my_app_policy.id

  rule {
    name        = "Require-MFA"
    priority    = 1
    factor_mode = "2FA"
    constraints = [jsonencode({
      "authenticationMethods": [{"key": "okta_verify", "method": "signed_nonce"}]
    })]
  }

  rule {
    name        = "Password-Fallback"
    priority    = 2
    factor_mode = "1FA"
    constraints = [jsonencode({
      "authenticationMethods": [{"key": "okta_password", "method": "password"}]
    })]
  }
}
```

### Migration Commands
```bash
# 1. Backup state
cp terraform.tfstate terraform.tfstate.backup

# 2. Remove old resources from state
terraform state rm okta_app_signon_policy_rule.mfa_rule
terraform state rm okta_app_signon_policy_rule.password_rule

# 3. Update your .tf files (replace old resources with new consolidated resource)

# 4. Import existing rules into new resource
terraform import okta_app_signon_policy_rules.my_app_policy_rules <policy_id>

# 5. Plan and apply
terraform plan
terraform apply
```

## Support

If you encounter issues during migration, please:
1. Check the [Terraform Provider for Okta documentation](https://registry.terraform.io/providers/okta/okta/latest/docs)
2. Search existing [GitHub issues](https://github.com/okta/terraform-provider-okta/issues)
3. Open a new issue with:
   - Provider version
   - Terraform version
   - Relevant configuration (sanitized)
   - Full error message
   - Steps to reproduce
