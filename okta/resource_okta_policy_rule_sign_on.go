package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
)

func resourcePolicySignonRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicySignOnRuleCreate,
		ReadContext:   resourcePolicySignOnRuleRead,
		UpdateContext: resourcePolicySignOnRuleUpdate,
		DeleteContext: resourcePolicySignOnRuleDelete,
		Importer:      createPolicyRuleImporter(),

		Schema: buildRuleSchema(map[string]*schema.Schema{
			"authtype": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: stringInSlice([]string{"ANY", "RADIUS"}),
				Description:      "Authentication entrypoint: ANY or RADIUS.",
				Default:          "ANY",
			},
			"access": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: stringInSlice([]string{"ALLOW", "DENY"}),
				Description:      "Allow or deny access based on the rule conditions: ALLOW or DENY.",
				Default:          "ALLOW",
			},
			"mfa_required": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Require MFA.",
				Default:     false,
			},
			"mfa_prompt": { // mfa_require must be true
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: stringInSlice([]string{"DEVICE", "SESSION", "ALWAYS"}),
				Description:      "Prompt for MFA based on the device used, a factor session lifetime, or every sign on attempt: DEVICE, SESSION or ALWAYS",
			},
			"mfa_remember_device": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Remember MFA device.",
				Default:     false,
			},
			"mfa_lifetime": { // mfa_require must be true, mfaprompt must be SESSION
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Elapsed time before the next MFA challenge",
			},
			"session_idle": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Max minutes a session can be idle.",
				Default:     120,
			},
			"session_lifetime": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Max minutes a session is active: Disable = 0.",
				Default:     120,
			},
			"session_persistent": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether session cookies will last across browser sessions. Okta Administrators can never have persistent session cookies.",
				Default:     false,
			},
		}),
	}
}

func resourcePolicySignOnRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	template := buildSignOnPolicyRule(d)
	err := createRule(ctx, d, m, template, policyRuleSignOn)
	if err != nil {
		return diag.Errorf("failed to create sign-on policy rule: %v", err)
	}
	return resourcePolicySignOnRuleRead(ctx, d, m)
}

func resourcePolicySignOnRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	rule, err := getPolicyRule(ctx, d, m)
	if err != nil {
		return diag.Errorf("failed to get sign-on policy rule: %v", err)
	}
	if rule == nil {
		return nil
	}

	// Update with upstream state to prevent stale state
	_ = d.Set("authtype", rule.Conditions.AuthContext.AuthType)
	_ = d.Set("access", rule.Actions.Signon.Access)
	_ = d.Set("mfa_required", rule.Actions.Signon.RequireFactor)
	_ = d.Set("mfa_remember_device", rule.Actions.Signon.RememberDeviceByDefault)
	_ = d.Set("mfa_lifetime", rule.Actions.Signon.FactorLifetime)
	_ = d.Set("session_idle", rule.Actions.Signon.Session.MaxSessionIdleMinutes)
	_ = d.Set("session_lifetime", rule.Actions.Signon.Session.MaxSessionLifetimeMinutes)
	_ = d.Set("session_persistent", rule.Actions.Signon.Session.UsePersistentCookie)

	if rule.Actions.Signon.FactorPromptMode != "" {
		_ = d.Set("mfa_prompt", rule.Actions.Signon.FactorPromptMode)
	}
	err = syncRuleFromUpstream(d, rule)
	if err != nil {
		return diag.Errorf("failed to sync sign-on policy rule: %v", err)
	}
	return nil
}

func resourcePolicySignOnRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	template := buildSignOnPolicyRule(d)
	err := updateRule(ctx, d, m, template)
	if err != nil {
		return diag.Errorf("failed to update sign-on policy rule: %v", err)
	}
	return resourcePolicySignOnRuleRead(ctx, d, m)
}

func resourcePolicySignOnRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := deleteRule(ctx, d, m, true)
	if err != nil {
		return diag.Errorf("failed to delete MFA policy rule: %v", err)
	}
	return nil
}

// Build Policy Sign On Rule from resource data
func buildSignOnPolicyRule(d *schema.ResourceData) sdk.PolicyRule {
	template := sdk.SignOnPolicyRule()
	template.Name = d.Get("name").(string)
	template.Status = d.Get("status").(string)
	if priority, ok := d.GetOk("priority"); ok {
		template.Priority = int64(priority.(int))
	}
	template.Conditions = &okta.PolicyRuleConditions{
		Network: getNetwork(d),
		AuthContext: &okta.PolicyRuleAuthContextCondition{
			AuthType: d.Get("authtype").(string),
		},
		People: getUsers(d),
	}
	template.Actions = sdk.PolicyRuleActions{
		OktaSignOnPolicyRuleActions: &okta.OktaSignOnPolicyRuleActions{
			Signon: &okta.OktaSignOnPolicyRuleSignonActions{
				Access:                  d.Get("access").(string),
				FactorLifetime:          int64(d.Get("mfa_lifetime").(int)),
				FactorPromptMode:        d.Get("mfa_prompt").(string),
				RememberDeviceByDefault: boolPtr(d.Get("mfa_remember_device").(bool)),
				RequireFactor:           boolPtr(d.Get("mfa_required").(bool)),
				Session: &okta.OktaSignOnPolicyRuleSignonSessionActions{
					MaxSessionIdleMinutes:     int64(d.Get("session_idle").(int)),
					MaxSessionLifetimeMinutes: int64(d.Get("session_lifetime").(int)),
					UsePersistentCookie:       boolPtr(d.Get("session_persistent").(bool)),
				},
			},
		},
	}
	return template
}
