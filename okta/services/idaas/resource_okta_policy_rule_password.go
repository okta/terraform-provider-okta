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
			"users_included": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Set of User IDs to include in this rule.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"groups_excluded": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Set of Group IDs to exclude from this rule.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"groups_included": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Set of Group IDs to include in this rule.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
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
			"password_reset_access_control": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Determines whether the Self-Service Password Reset (SSPR) access is governed by an authentication policy or legacy behavior. Options: `LEGACY`, `AUTH_POLICY`.",
			},
			// password_reset_requirement defines the SSPR requirement settings (primary methods, step-up).
			"password_reset_requirement": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Self-service password reset (SSPR) requirement settings. Use only when `password_reset_access_control = \"LEGACY\"`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"method_constraints": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Constraints on the values specified in the methods array. Specifying a constraint limits methods to specific authenticator(s). Currently, Google OTP is the only accepted constraint.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"method": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The method to constrain (e.g. `otp`).",
									},
									"allowed_authenticators": {
										Type:        schema.TypeSet,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Keys of the authenticators allowed for this method (e.g. `google_otp`).",
									},
								},
							},
						},
						"primary_methods": {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Authenticator methods allowed for the initial authentication step of password recovery. Method otp requires a constraint limiting it to a Google authenticator. Options: `otp`, `push`, `sms`, `email`, `voice`.",
						},
						"step_up_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Whether a secondary authenticator is required for password reset (`stepUp.required`). The following are three valid configurations: `required=false`, `required=true` with no methods to use any SSO authenticator, and `required=true` with `security_question` as the method.",
						},
						"step_up_methods": {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Authenticator methods required for the secondary authentication step of password recovery. Specify only when `step_up_enabled = true` and `security_question` is permitted. Items value: `security_question`.",
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
		inner, resp, createErr := getOktaV6ClientFromMetadata(meta).PolicyAPI.CreatePolicyRule(ctx, policyID).PolicyRule(wrapper).Execute()
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
		_, deactErr := getOktaV6ClientFromMetadata(meta).PolicyAPI.DeactivatePolicyRule(ctx, policyID, pr.GetId()).Execute()
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

	inner, resp, err := getOktaV6ClientFromMetadata(meta).PolicyAPI.GetPolicyRule(ctx, policyID, d.Id()).Execute()
	if err != nil {
		// Rule (or parent policy) no longer exists — remove from state.
		if resp != nil && resp.Response != nil && resp.StatusCode == http.StatusNotFound {
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
		// Sync people conditions.
		if people := conds.People; people != nil {
			if users := people.Users; users != nil {
				_ = d.Set("users_excluded", utils.ConvertStringSliceToSetNullable(users.Exclude))
				_ = d.Set("users_included", utils.ConvertStringSliceToSetNullable(users.Include))
			}
			if groups := people.Groups; groups != nil {
				_ = d.Set("groups_excluded", utils.ConvertStringSliceToSetNullable(groups.Exclude))
				_ = d.Set("groups_included", utils.ConvertStringSliceToSetNullable(groups.Include))
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
			if req := sspr.Requirement; req != nil {
				if ac := req.AccessControl; ac != nil {
					_ = d.Set("password_reset_access_control", ac)
				}
				reqMap := map[string]interface{}{
					"primary_methods":    []string{},
					"method_constraints": []interface{}{},
					"step_up_enabled":    false,
					"step_up_methods":    []string{},
				}
				if primary := req.Primary; primary != nil {
					reqMap["primary_methods"] = primary.Methods
					methodConstraints := make([]interface{}, 0, len(primary.MethodConstraints))
					for _, mc := range primary.MethodConstraints {
						allowedAuths := make([]interface{}, 0, len(mc.GetAllowedAuthenticators()))
						for _, identity := range mc.GetAllowedAuthenticators() {
							allowedAuths = append(allowedAuths, identity.GetKey())
						}
						methodConstraints = append(methodConstraints, map[string]interface{}{
							"method":                 mc.GetMethod(),
							"allowed_authenticators": allowedAuths,
						})
					}
					reqMap["method_constraints"] = methodConstraints
				}
				if stepUp := req.StepUp; stepUp != nil {
					reqMap["step_up_enabled"] = stepUp.GetRequired()
					reqMap["step_up_methods"] = stepUp.Methods
				}
				// Only persist the requirement block when it contains meaningful content.
				// If the API returns an empty requirement (e.g. with AUTH_POLICY), omitting
				// the block prevents a permanent diff for configs that don't set it.
				primaryMethods, _ := reqMap["primary_methods"].([]string)
				methodConstraints, _ := reqMap["method_constraints"].([]interface{})
				stepUpEnabled, _ := reqMap["step_up_enabled"].(bool)
				stepUpMethods, _ := reqMap["step_up_methods"].([]string)
				hasContent := len(primaryMethods) > 0 || len(methodConstraints) > 0 || stepUpEnabled || len(stepUpMethods) > 0
				if hasContent {
					_ = d.Set("password_reset_requirement", []interface{}{reqMap})
				} else {
					_ = d.Set("password_reset_requirement", []interface{}{})
				}
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
	updatedInner, _, err := getOktaV6ClientFromMetadata(meta).PolicyAPI.ReplacePolicyRule(ctx, policyID, d.Id()).PolicyRule(wrapper).Execute()
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

	resp, err := getOktaV6ClientFromMetadata(meta).PolicyAPI.DeletePolicyRule(ctx, policyID, d.Id()).Execute()
	if err != nil {
		// Already gone — not an error.
		if resp != nil && resp.Response != nil && resp.StatusCode == http.StatusNotFound {
			return nil
		}
		return diag.Errorf("failed to delete password policy rule: %v", err)
	}
	return nil
}

// buildPolicyRulePassword constructs a V6 PasswordPolicyRule from schema data.
// The SSPR action is built with an optional requirement block when password_reset_access_control is set.
// Uses the V6 okta-sdk-golang client types directly, removing the dependency on the internal SDK.
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

	// Build the people conditions (users and groups).
	peopleCond := v6okta.PolicyPeopleCondition{}
	hasPeople := false

	userCond := v6okta.UserCondition{}
	if exclude, ok := d.GetOk("users_excluded"); ok {
		userCond.Exclude = utils.ConvertInterfaceToStringSet(exclude)
		hasPeople = true
	}
	if include, ok := d.GetOk("users_included"); ok {
		userCond.Include = utils.ConvertInterfaceToStringSet(include)
		hasPeople = true
	}
	if hasPeople {
		peopleCond.SetUsers(userCond)
	}

	groupCond := v6okta.GroupCondition{}
	hasGroups := false
	if exclude, ok := d.GetOk("groups_excluded"); ok {
		groupCond.Exclude = utils.ConvertInterfaceToStringSet(exclude)
		hasGroups = true
	}
	if include, ok := d.GetOk("groups_included"); ok {
		groupCond.Include = utils.ConvertInterfaceToStringSet(include)
		hasGroups = true
	}
	if hasGroups {
		peopleCond.SetGroups(groupCond)
	}

	if hasPeople || hasGroups {
		conds.SetPeople(peopleCond)
	}
	rule.SetConditions(conds)

	// Build the SSPR action, conditionally including the requirement block
	// when password_reset_access_control is explicitly configured.
	ssprAction := v6okta.SelfServicePasswordResetAction{}
	access := d.Get("password_reset").(string)
	ssprAction.SetAccess(access)

	if accessControl, ok := d.GetOk("password_reset_access_control"); ok {
		req := v6okta.SsprRequirement{}
		req.SetAccessControl(accessControl.(string))

		// Only build primary methods and step-up if the requirement block is provided.
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
						if methodConstraintsRaw, ok := reqBlock["method_constraints"]; ok {
							methodConstraints := methodConstraintsRaw.([]interface{})
							if len(methodConstraints) > 0 {
								authenticatorMethodConstraints := make([]v6okta.AuthenticatorMethodConstraint, 0, len(methodConstraints))
								for _, methodConstraint := range methodConstraints {
									mcMap := methodConstraint.(map[string]interface{})
									constraint := v6okta.NewAuthenticatorMethodConstraint()
									if method, ok := mcMap["method"].(string); ok && method != "" {
										constraint.SetMethod(method)
									}
									if allowedAuthenticatorsRaw, ok := mcMap["allowed_authenticators"]; ok {
										identities := []v6okta.AuthenticatorIdentity{}
										if allowedAuthenticatorsSet, ok := allowedAuthenticatorsRaw.(*schema.Set); ok {
											allowedAuthenticatorsList := allowedAuthenticatorsSet.List()
											for _, key := range allowedAuthenticatorsList {
												if allowedAuthenticator, ok := key.(string); ok {
													identity := v6okta.NewAuthenticatorIdentity()
													identity.SetKey(allowedAuthenticator)
													identities = append(identities, *identity)
												}
											}
											if len(identities) > 0 {
												constraint.SetAllowedAuthenticators(identities)
											}
										}
									}
									authenticatorMethodConstraints = append(authenticatorMethodConstraints, *constraint)
								}
								primary.SetMethodConstraints(authenticatorMethodConstraints)
							}
						}
						primary.SetMethods(methodList)
						req.SetPrimary(primary)
					}
				}
				stepUpReq := v6okta.SsprStepUpRequirement{}
				if stepUp, ok := reqBlock["step_up_enabled"]; ok {
					stepUpReq.SetRequired(stepUp.(bool))
				}
				if stepUpMethodsRaw, ok := reqBlock["step_up_methods"]; ok {
					if methodSet, ok := stepUpMethodsRaw.(*schema.Set); ok && methodSet.Len() > 0 {
						methods := make([]string, 0, methodSet.Len())
						for _, m := range methodSet.List() {
							methods = append(methods, m.(string))
						}
						stepUpReq.SetMethods(methods)
					}
				}
				req.SetStepUp(stepUpReq)
			}
		}
		// Okta requires requirement.primary and requirement.stepUp to be non-null
		// whenever a requirement object is sent, even for AUTH_POLICY. Set defaults
		// if the requirement block was not explicitly configured.
		if req.Primary == nil {
			req.SetPrimary(v6okta.SsprPrimaryRequirement{Methods: []string{}})
		}
		if req.StepUp == nil {
			defaultStepUp := v6okta.SsprStepUpRequirement{}
			defaultStepUp.SetRequired(false)
			req.SetStepUp(defaultStepUp)
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
