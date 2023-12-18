package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

func resourcePolicyPasswordRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyPasswordRuleCreate,
		ReadContext:   resourcePolicyPasswordRuleRead,
		UpdateContext: resourcePolicyPasswordRuleUpdate,
		DeleteContext: resourcePolicyPasswordRuleDelete,
		Importer:      createPolicyRuleImporter(),
		Description:   "Creates a Password Policy Rule. This resource allows you to create and configure a Password Policy Rule.",
		Schema: buildRuleSchema(map[string]*schema.Schema{
			"password_change": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Allow or deny a user to change their password: `ALLOW` or `DENY`. Default: `ALLOW`",
				Default:     "ALLOW",
			},
			"password_reset": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Allow or deny a user to reset their password: `ALLOW` or `DENY`. Default: `ALLOW`",
				Default:     "ALLOW",
			},
			"password_unlock": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Allow or deny a user to unlock. Default: `DENY`",
				Default:     "DENY",
			},
		}),
	}
}

func resourcePolicyPasswordRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	template := buildPolicyRulePassword(d)
	err := createRule(ctx, d, m, template, policyRulePassword)
	if err != nil {
		return diag.Errorf("failed to create password policy rule: %v", err)
	}
	return resourcePolicyPasswordRuleRead(ctx, d, m)
}

func resourcePolicyPasswordRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	rule, err := getPolicyRule(ctx, d, m)
	if err != nil {
		return diag.Errorf("failed to get password policy rule: %v", err)
	}
	if rule == nil {
		return nil
	}
	// Update with upstream state to prevent stale state
	_ = d.Set("password_change", rule.Actions.PasswordChange.Access)
	_ = d.Set("password_unlock", rule.Actions.SelfServiceUnlock.Access)
	_ = d.Set("password_reset", rule.Actions.SelfServicePasswordReset.Access)
	err = syncRuleFromUpstream(d, rule)
	if err != nil {
		return diag.Errorf("failed to sync password policy rule: %v", err)
	}
	return nil
}

func resourcePolicyPasswordRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	template := buildPolicyRulePassword(d)
	err := updateRule(ctx, d, m, template)
	if err != nil {
		return diag.Errorf("failed to update password policy rule: %v", err)
	}
	return resourcePolicyPasswordRuleRead(ctx, d, m)
}

func resourcePolicyPasswordRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := deleteRule(ctx, d, m, false)
	if err != nil {
		return diag.Errorf("failed to delete password policy rule: %v", err)
	}
	return nil
}

// build password policy rule from schema data
func buildPolicyRulePassword(d *schema.ResourceData) sdk.SdkPolicyRule {
	template := sdk.PasswordPolicyRule()
	template.Name = d.Get("name").(string)
	template.Status = d.Get("status").(string)
	if priority, ok := d.GetOk("priority"); ok {
		template.Priority = int64(priority.(int))
	}
	template.Conditions = &sdk.PolicyRuleConditions{
		Network: buildPolicyNetworkCondition(d),
		People:  getUsers(d),
	}
	template.Actions = sdk.SdkPolicyRuleActions{
		PasswordPolicyRuleActions: &sdk.PasswordPolicyRuleActions{
			PasswordChange: &sdk.PasswordPolicyRuleAction{
				Access: d.Get("password_change").(string),
			},
			SelfServicePasswordReset: &sdk.PasswordPolicyRuleAction{
				Access: d.Get("password_reset").(string),
			},
			SelfServiceUnlock: &sdk.PasswordPolicyRuleAction{
				Access: d.Get("password_unlock").(string),
			},
		},
	}
	return template
}
