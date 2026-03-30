package idaas

import (
	"context"
	"net/http"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	v6okta "github.com/okta/okta-sdk-golang/v6/okta"
	"github.com/okta/terraform-provider-okta/okta/utils"
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
	// Prevent creating a rule named "Default Rule" which is managed by Okta.
	if err := ensureNotDefaultRule(d); err != nil {
		return diag.FromErr(err)
	}
	policyID := d.Get("policy_id").(string)
	if policyID == "" {
		return diag.Errorf("'policy_id' field should be set")
	}

	rule := buildPolicyRulePassword(d)
	// Wrap the typed rule in the union type required by the V6 PolicyAPI.
	wrapper := v6okta.PasswordPolicyRuleAsListPolicyRules200ResponseInner(&rule)

	// Use exponential backoff to handle transient 500 errors from the Okta API.
	var createdInner *v6okta.ListPolicyRules200ResponseInner
	boc := utils.NewExponentialBackOffWithContext(ctx, backoff.DefaultMaxElapsedTime)
	err := backoff.Retry(func() error {
		inner, resp, createErr := getOktaV6ClientFromMetadata(meta).PolicyAPI.
			CreatePolicyRule(ctx, policyID).PolicyRule(wrapper).Execute()
		if doNotRetry(meta, createErr) {
			return backoff.Permanent(createErr)
		}
		if createErr != nil {
			return backoff.Permanent(createErr)
		}
		if resp.StatusCode == http.StatusInternalServerError {
			return createErr
		}
		createdInner = inner
		return nil
	}, boc)
	if err != nil {
		return diag.Errorf("failed to create password policy rule: %v", err)
	}

	pr := createdInner.PasswordPolicyRule
	d.SetId(pr.GetId())

	// If the desired status is INACTIVE, deactivate immediately after creation
	// since the Okta API always creates rules in ACTIVE state.
	if d.Get("status").(string) == StatusInactive {
		_, deactErr := getOktaV6ClientFromMetadata(meta).PolicyAPI.
			DeactivatePolicyRule(ctx, policyID, pr.GetId()).Execute()
		if deactErr != nil {
			return diag.Errorf("failed to deactivate password policy rule on creation: %v", deactErr)
		}
	}

	// Validate that the API honoured the requested priority.
	if err := utils.ValidatePriority(int64(rule.GetPriority()), int64(pr.GetPriority())); err != nil {
		return diag.FromErr(err)
	}

	return resourcePolicyPasswordRuleRead(ctx, d, meta)
}

func resourcePolicyPasswordRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	policyID := d.Get("policy_id").(string)
	if policyID == "" {
		return diag.Errorf("'policy_id' field should be set")
	}

	inner, resp, err := getOktaV6ClientFromMetadata(meta).PolicyAPI.
		GetPolicyRule(ctx, policyID, d.Id()).Execute()
	if err != nil {
		// Rule (or parent policy) no longer exists — remove from state.
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to get password policy rule: %v", err)
	}
	if inner == nil || inner.PasswordPolicyRule == nil {
		d.SetId("")
		return nil
	}

	rule := inner.PasswordPolicyRule

	// Sync base rule fields.
	_ = d.Set("name", rule.GetName())
	_ = d.Set("status", rule.GetStatus())
	_ = d.Set("priority", int(rule.GetPriority()))

	// Sync network conditions.
	if conds := rule.Conditions; conds != nil {
		if net := conds.Network; net != nil {
			_ = d.Set("network_connection", net.GetConnection())
			if len(net.Include) > 0 {
				_ = d.Set("network_includes", utils.ConvertStringSliceToInterfaceSlice(net.Include))
			}
			if len(net.Exclude) > 0 {
				_ = d.Set("network_excludes", utils.ConvertStringSliceToInterfaceSlice(net.Exclude))
			}
		}
		// Sync users_excluded.
		if people := conds.People; people != nil {
			if users := people.Users; users != nil {
				_ = d.Set("users_excluded", utils.ConvertStringSliceToSetNullable(users.Exclude))
			}
		}
	}

	// Sync password action fields.
	if actions := rule.Actions; actions != nil {
		if pc := actions.PasswordChange; pc != nil {
			_ = d.Set("password_change", pc.GetAccess())
		}
		if su := actions.SelfServiceUnlock; su != nil {
			_ = d.Set("password_unlock", su.GetAccess())
		}
		// Sync SSPR action, including the requirement block added for GH-2559.
		if sspr := actions.SelfServicePasswordReset; sspr != nil {
			_ = d.Set("password_reset", sspr.GetAccess())
			// Read SSPR requirement fields from upstream and set them in state.
			// The logic has been updated by AI to support new fields, review carefully
			if req := sspr.Requirement; req != nil {
				if ac := req.AccessControl; ac != nil {
					_ = d.Set("password_reset_access_control", *ac)
				}
				reqMap := map[string]interface{}{
					"primary_methods": []string{},
					"step_up_enabled": false,
				}
				if primary := req.Primary; primary != nil {
					reqMap["primary_methods"] = primary.Methods
				}
				if stepUp := req.StepUp; stepUp != nil {
					reqMap["step_up_enabled"] = stepUp.GetRequired()
				}
				_ = d.Set("password_reset_requirement", []interface{}{reqMap})
			}
		}
	}

	return nil
}

func resourcePolicyPasswordRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if err := ensureNotDefaultRule(d); err != nil {
		return diag.FromErr(err)
	}
	policyID := d.Get("policy_id").(string)
	if policyID == "" {
		return diag.Errorf("'policy_id' field should be set")
	}

	rule := buildPolicyRulePassword(d)
	wrapper := v6okta.PasswordPolicyRuleAsListPolicyRules200ResponseInner(&rule)

	// ReplacePolicyRule is the V6 equivalent of UpdatePolicyRule (HTTP PUT).
	updatedInner, _, err := getOktaV6ClientFromMetadata(meta).PolicyAPI.
		ReplacePolicyRule(ctx, policyID, d.Id()).PolicyRule(wrapper).Execute()
	if err != nil {
		return diag.Errorf("failed to update password policy rule: %v", err)
	}

	if err := utils.ValidatePriority(int64(rule.GetPriority()), int64(updatedInner.PasswordPolicyRule.GetPriority())); err != nil {
		return diag.FromErr(err)
	}

	// Activate or deactivate the rule to match the desired status.
	client := getOktaV6ClientFromMetadata(meta).PolicyAPI
	if d.Get("status").(string) == StatusActive {
		if _, actErr := client.ActivatePolicyRule(ctx, policyID, d.Id()).Execute(); actErr != nil {
			return diag.Errorf("failed to activate password policy rule: %v", actErr)
		}
	} else if d.Get("status").(string) == StatusInactive {
		if _, deactErr := client.DeactivatePolicyRule(ctx, policyID, d.Id()).Execute(); deactErr != nil {
			return diag.Errorf("failed to deactivate password policy rule: %v", deactErr)
		}
	}

	return resourcePolicyPasswordRuleRead(ctx, d, meta)
}

func resourcePolicyPasswordRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if err := ensureNotDefaultRule(d); err != nil {
		return diag.FromErr(err)
	}
	policyID := d.Get("policy_id").(string)
	if policyID == "" {
		return diag.Errorf("'policy_id' field should be set")
	}

	resp, err := getOktaV6ClientFromMetadata(meta).PolicyAPI.
		DeletePolicyRule(ctx, policyID, d.Id()).Execute()
	if err != nil {
		// Already gone — not an error.
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil
		}
		return diag.Errorf("failed to delete password policy rule: %v", err)
	}
	return nil
}

// buildPolicyRulePassword constructs a V6 PasswordPolicyRule from schema data.
// The SSPR action is built with an optional requirement block when password_reset_access_control is set.
// Uses the V6 okta-sdk-golang client types directly, removing the dependency on the internal SDK.
// The logic has been updated by AI to support new fields, review carefully
func buildPolicyRulePassword(d *schema.ResourceData) v6okta.PasswordPolicyRule {
	rule := v6okta.NewPasswordPolicyRule()
	rule.SetType("PASSWORD")
	rule.SetName(d.Get("name").(string))
	rule.SetStatus(d.Get("status").(string))
	if priority, ok := d.GetOk("priority"); ok {
		rule.SetPriority(int32(priority.(int)))
	}

	// Build network and people conditions using V6 types.
	networkCond := v6okta.PolicyNetworkCondition{}
	networkCond.SetConnection(d.Get("network_connection").(string))
	networkCond.Include = utils.ConvertInterfaceToStringArrNullable(d.Get("network_includes"))
	networkCond.Exclude = utils.ConvertInterfaceToStringArrNullable(d.Get("network_excludes"))

	conds := v6okta.PasswordPolicyRuleConditions{}
	conds.SetNetwork(networkCond)

	// Build the people/users_excluded condition.
	if exclude, ok := d.GetOk("users_excluded"); ok {
		userCond := v6okta.UserCondition{}
		userCond.Exclude = utils.ConvertInterfaceToStringSet(exclude)
		peopleCond := v6okta.PolicyPeopleCondition{}
		peopleCond.SetUsers(userCond)
		conds.SetPeople(peopleCond)
	}
	rule.SetConditions(conds)

	// Build the SSPR action, conditionally including the requirement block
	// when password_reset_access_control is explicitly configured.
	// The logic has been updated by AI to support new fields, review carefully
	ssprAction := v6okta.SelfServicePasswordResetAction{}
	access := d.Get("password_reset").(string)
	ssprAction.SetAccess(access)

	if accessControl, ok := d.GetOk("password_reset_access_control"); ok {
		req := v6okta.SsprRequirement{}
		req.SetAccessControl(accessControl.(string))

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
						primary := v6okta.SsprPrimaryRequirement{}
						primary.SetMethods(methodList)
						req.SetPrimary(primary)
					}
				}
				if stepUp, ok := reqBlock["step_up_enabled"]; ok {
					stepUpReq := v6okta.SsprStepUpRequirement{}
					stepUpReq.SetRequired(stepUp.(bool))
					req.SetStepUp(stepUpReq)
				}
			}
		}
		ssprAction.SetRequirement(req)
	}

	actions := v6okta.PasswordPolicyRuleActions{}
	pcAction := v6okta.PasswordPolicyRuleAction{}
	pcAction.SetAccess(d.Get("password_change").(string))
	actions.SetPasswordChange(pcAction)

	actions.SetSelfServicePasswordReset(ssprAction)

	suAction := v6okta.PasswordPolicyRuleAction{}
	suAction.SetAccess(d.Get("password_unlock").(string))
	actions.SetSelfServiceUnlock(suAction)

	rule.SetActions(actions)
	return *rule
}
