package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/resources"
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
			// password_reset_access_control controls whether SSPR uses legacy or auth-policy access control.
			// The logic has been updated by AI to support new fields, review carefully
			"password_reset_access_control": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Determines whether the Self-Service Password Reset (SSPR) access is governed by an authentication policy or legacy behavior. Options: `LEGACY`, `AUTH_POLICY`. The logic has been updated by AI to support new fields, review carefully",
			},
			// password_reset_requirement defines the SSPR requirement settings (primary methods, step-up).
			// The logic has been updated by AI to support new fields, review carefully
			"password_reset_requirement": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Self-service password reset (SSPR) requirement settings. Use only when `password_reset_access_control = \"LEGACY\"`. The logic has been updated by AI to support new fields, review carefully",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"primary_methods": {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "List of primary authentication methods for SSPR. Options: `otp`, `push`, `sms`, `email`, `voice`.",
						},
						"step_up_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Determines whether step-up authentication is required for SSPR.",
						},
					},
				},
			},
		}),
	}
}

func resourcePolicyPasswordRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	template := buildPolicyRulePassword(d)
	err := createRule(ctx, d, meta, template, resources.OktaIDaaSPolicyRulePassword)
	if err != nil {
		return diag.Errorf("failed to create password policy rule: %v", err)
	}
	return resourcePolicyPasswordRuleRead(ctx, d, meta)
}

func resourcePolicyPasswordRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rule, err := getPolicyRule(ctx, d, meta)
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

	// Read SSPR requirement fields from upstream and set them in state.
	// The logic has been updated by AI to support new fields, review carefully
	if sspr := rule.Actions.SelfServicePasswordReset; sspr != nil && sspr.Requirement != nil {
		if sspr.Requirement.AccessControl != "" {
			_ = d.Set("password_reset_access_control", sspr.Requirement.AccessControl)
		}
		reqMap := map[string]interface{}{
			"primary_methods": []string{},
			"step_up_enabled": false,
		}
		if sspr.Requirement.Primary != nil {
			reqMap["primary_methods"] = sspr.Requirement.Primary.Methods
		}
		if sspr.Requirement.StepUp != nil {
			reqMap["step_up_enabled"] = sspr.Requirement.StepUp.Required
		}
		_ = d.Set("password_reset_requirement", []interface{}{reqMap})
	}

	err = syncRuleFromUpstream(d, rule)
	if err != nil {
		return diag.Errorf("failed to sync password policy rule: %v", err)
	}
	return nil
}

func resourcePolicyPasswordRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	template := buildPolicyRulePassword(d)
	err := updateRule(ctx, d, meta, template)
	if err != nil {
		return diag.Errorf("failed to update password policy rule: %v", err)
	}
	return resourcePolicyPasswordRuleRead(ctx, d, meta)
}

func resourcePolicyPasswordRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := deleteRule(ctx, d, meta, false)
	if err != nil {
		return diag.Errorf("failed to delete password policy rule: %v", err)
	}
	return nil
}

// buildPolicyRulePassword constructs the SDK policy rule from schema data.
// The SSPR action is built with an optional requirement block when password_reset_access_control is set.
// The logic has been updated by AI to support new fields, review carefully
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

	// Build the SSPR action, conditionally including the requirement block
	// when password_reset_access_control is explicitly configured.
	// The logic has been updated by AI to support new fields, review carefully
	ssprAction := &sdk.PasswordPolicyRuleAction{
		Access: d.Get("password_reset").(string),
	}
	if accessControl, ok := d.GetOk("password_reset_access_control"); ok {
		req := &sdk.PasswordPolicyRuleRequirement{
			AccessControl: accessControl.(string),
		}
		// Only build primary methods and step-up if the requirement block is provided.
		// The logic has been updated by AI to support new fields, review carefully
		if v, ok := d.GetOk("password_reset_requirement"); ok {
			reqList := v.([]interface{})
			if len(reqList) > 0 {
				reqBlock := reqList[0].(map[string]interface{})
				if methods, ok := reqBlock["primary_methods"]; ok {
					methodSet := methods.(*schema.Set)
					if methodSet.Len() > 0 {
						methodList := make([]string, 0, methodSet.Len())
						for _, m := range methodSet.List() {
							methodList = append(methodList, m.(string))
						}
						req.Primary = &sdk.PasswordPolicyRuleRequirementPrimary{
							Methods: methodList,
						}
					}
				}
				if stepUp, ok := reqBlock["step_up_enabled"]; ok {
					req.StepUp = &sdk.PasswordPolicyRuleRequirementStepUp{
						Required: stepUp.(bool),
					}
				}
			}
		}
		ssprAction.Requirement = req
	}

	template.Actions = sdk.SdkPolicyRuleActions{
		PasswordPolicyRuleActions: &sdk.PasswordPolicyRuleActions{
			PasswordChange: &sdk.PasswordPolicyRuleAction{
				Access: d.Get("password_change").(string),
			},
			SelfServicePasswordReset: ssprAction,
			SelfServiceUnlock: &sdk.PasswordPolicyRuleAction{
				Access: d.Get("password_unlock").(string),
			},
		},
	}
	return template
}
