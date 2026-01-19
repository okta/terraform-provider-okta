package idaas

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
)

func resourcePolicySignOnRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicySignOnRuleCreate,
		ReadContext:   resourcePolicySignOnRuleRead,
		UpdateContext: resourcePolicySignOnRuleUpdate,
		DeleteContext: resourcePolicySignOnRuleDelete,
		Importer:      createPolicyRuleImporter(),
		Description:   "Creates a Sign On Policy Rule. In case `Invalid condition type specified: riskScore.` error is thrown, set `risc_level` to an empty string, since this feature is not enabled.",
		Schema: buildRuleSchema(map[string]*schema.Schema{
			"authtype": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Authentication entrypoint: `ANY`, `RADIUS` or `LDAP_INTERFACE`. Default: `ANY`",
				Default:     "ANY",
			},
			"access": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Allow or deny access based on the rule conditions: `ALLOW`, `DENY` or `CHALLENGE`. Default: `ALLOW`",
				Default:     "ALLOW",
			},
			"mfa_required": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Require MFA. Default: `false`",
				Default:     false,
			},
			"mfa_prompt": { // mfa_require must be true
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Prompt for MFA based on the device used, a factor session lifetime, or every sign-on attempt: `DEVICE`, `SESSION` or`ALWAYS`.",
			},
			"mfa_remember_device": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Remember MFA device. Default: `false`",
				Default:     false,
			},
			"mfa_lifetime": { // mfa_require must be true, mfa_prompt must be SESSION
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Elapsed time before the next MFA challenge",
			},
			"session_idle": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Max minutes a session can be idle. Default: `120`",
				Default:     120,
			},
			"session_lifetime": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Max minutes a session is active: Disable = 0. Default: `120`",
				Default:     120,
			},
			"session_persistent": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether session cookies will last across browser sessions. Okta Administrators can never have persistent session cookies. Default: `false`",
				Default:     false,
			},
			"risc_level": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Risc level: ANY, LOW, MEDIUM or HIGH. Default: `ANY`",
				Deprecated:  "Attribute typo, switch to risk_level instead. Default: `ANY`",
			},
			"risk_level": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Risk level: ANY, LOW, MEDIUM or HIGH. Default: `ANY`",
			},
			"behaviors": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of behavior IDs",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"primary_factor": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Rule's primary factor. **WARNING** Ony works as a part of the Identity Engine. Valid values: `PASSWORD_IDP_ANY_FACTOR`, `PASSWORD_IDP`.",
			},
			"factor_sequence": {
				Type:     schema.TypeList,
				Optional: true,
				Description: `Auth factor sequences. Should be set if 'access = "CHALLENGE"'.
	- 'primary_criteria_provider' - (Required) Primary provider of the auth section.
	- 'primary_criteria_factor_type' - (Required) Primary factor type of the auth section.
	- 'secondary_criteria' - (Optional) Additional authentication steps.
	- 'provider' - (Required) Provider of the additional authentication step.
	- 'factor_type' - (Required) Factor type of the additional authentication step.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"primary_criteria_provider": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Factor provider",
						},
						"primary_criteria_factor_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Type of a Factor",
						},
						"secondary_criteria": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"provider": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Factor provider",
									},
									"factor_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Type of a Factor",
									},
								},
							},
						},
					},
				},
			},
			"identity_provider": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Apply rule based on the IdP used: `ANY`, `OKTA` or `SPECIFIC_IDP`. Default: `ANY`. ~> **WARNING**: Use of `identity_provider` requires a feature flag to be enabled.",
				Default:     "ANY",
			},
			"identity_provider_ids": { // identity_provider must be SPECIFIC_IDP
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "When identity_provider is `SPECIFIC_IDP` then this is the list of IdP IDs to apply the rule on",
			},
		}),
	}
}

func resourcePolicySignOnRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := validateSignOnPolicyRule(d)
	if err != nil {
		return diag.FromErr(err)
	}
	template := buildSignOnPolicyRule(d)
	err = createRule(ctx, d, meta, template, resources.OktaIDaaSPolicyRuleSignOn)
	if err != nil {
		return diag.Errorf("failed to create sign-on policy rule: %v", err)
	}
	return resourcePolicySignOnRuleRead(ctx, d, meta)
}

func resourcePolicySignOnRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rule, err := getPolicyRule(ctx, d, meta)
	if err != nil {
		return diag.Errorf("failed to get sign-on policy rule: %v", err)
	}
	if rule == nil {
		return nil
	}
	// Update with upstream state to prevent stale state
	_ = d.Set("authtype", rule.Conditions.AuthContext.AuthType)
	_ = d.Set("access", rule.Actions.SignOn.Access)
	_ = d.Set("mfa_required", rule.Actions.SignOn.RequireFactor)
	_ = d.Set("mfa_remember_device", rule.Actions.SignOn.RememberDeviceByDefault)
	_ = d.Set("mfa_lifetime", rule.Actions.SignOn.FactorLifetime)
	if rule.Actions.SignOn.Session.MaxSessionIdleMinutesPtr != nil {
		_ = d.Set("session_idle", rule.Actions.SignOn.Session.MaxSessionIdleMinutesPtr)
	}
	if rule.Actions.SignOn.Session.MaxSessionLifetimeMinutesPtr != nil {
		_ = d.Set("session_lifetime", rule.Actions.SignOn.Session.MaxSessionLifetimeMinutesPtr)
	}
	_ = d.Set("session_persistent", rule.Actions.SignOn.Session.UsePersistentCookie)
	if rule.Actions.SignOn.FactorPromptMode != "" {
		_ = d.Set("mfa_prompt", rule.Actions.SignOn.FactorPromptMode)
	}
	if rule.Actions.SignOn.PrimaryFactor != "" {
		_ = d.Set("primary_factor", rule.Actions.SignOn.PrimaryFactor)
	}
	if rule.Conditions != nil {
		if rule.Conditions.RiskScore != nil {
			curRiscLevel, riscLevelSet := d.GetOk("risc_level")
			curRiskLevel, riskLevelSet := d.GetOk("risk_level")
			if riskLevelSet {
				_ = d.Set("risk_level", rule.Conditions.RiskScore.Level)
				_ = d.Set("risc_level", curRiscLevel) // retain current value to avoid diff during plan.

			} else if riscLevelSet {
				_ = d.Set("risc_level", rule.Conditions.RiskScore.Level)
				_ = d.Set("risk_level", curRiskLevel) // retain current value to avoid diff during plan.
			}
		}
		if rule.Conditions.Risk != nil {
			err = utils.SetNonPrimitives(d, map[string]interface{}{
				"behaviors": utils.ConvertStringSliceToSet(rule.Conditions.Risk.Behaviors),
			})
			if err != nil {
				return diag.Errorf("failed to set sign-on policy rule behaviors: %v", err)
			}
		}
	}
	if rule.Conditions.IdentityProvider != nil {
		_ = d.Set("identity_provider", rule.Conditions.IdentityProvider.Provider)
		if rule.Conditions.IdentityProvider.Provider == "SPECIFIC_IDP" {
			_ = d.Set("identity_provider_ids", utils.ConvertStringSliceToInterfaceSlice(rule.Conditions.IdentityProvider.IdpIds))
		}
	}

	if rule.Actions.SignOn.Access == "CHALLENGE" {
		chain := rule.Actions.SignOn.Challenge.Chain
		arr := make([]map[string]interface{}, len(chain))
		for i, c := range chain {
			arr[i] = map[string]interface{}{
				"primary_criteria_provider":    c.Criteria[0].Provider,
				"primary_criteria_factor_type": c.Criteria[0].FactorType,
			}
			if len(c.Next) > 0 {
				scs := make([]map[string]interface{}, len(c.Next[0].Criteria))
				for j, sc := range c.Next[0].Criteria {
					scs[j] = map[string]interface{}{
						"provider":    sc.Provider,
						"factor_type": sc.FactorType,
					}
				}
				arr[i]["secondary_criteria"] = scs
			}
		}
		err = utils.SetNonPrimitives(d, map[string]interface{}{"factor_sequence": arr})
		if err != nil {
			return diag.Errorf("failed to set OAuth application properties: %v", err)
		}
	}
	err = syncRuleFromUpstream(d, rule)
	if err != nil {
		return diag.Errorf("failed to sync sign-on policy rule: %v", err)
	}
	return nil
}

func resourcePolicySignOnRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := validateSignOnPolicyRule(d)
	if err != nil {
		return diag.FromErr(err)
	}
	template := buildSignOnPolicyRule(d)
	err = updateRule(ctx, d, meta, template)
	if err != nil {
		return diag.Errorf("failed to update sign-on policy rule: %v", err)
	}
	return resourcePolicySignOnRuleRead(ctx, d, meta)
}

func resourcePolicySignOnRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := deleteRule(ctx, d, meta, true)
	if err != nil {
		return diag.Errorf("failed to delete MFA policy rule: %v", err)
	}
	return nil
}

// Build Policy Sign On Rule from resource data
func buildSignOnPolicyRule(d *schema.ResourceData) sdk.SdkPolicyRule {
	template := sdk.SignOnPolicyRule()
	template.Name = d.Get("name").(string)
	template.Status = d.Get("status").(string)
	if priority, ok := d.GetOk("priority"); ok {
		template.Priority = int64(priority.(int))
	}
	template.Conditions = &sdk.PolicyRuleConditions{
		AuthContext: &sdk.PolicyRuleAuthContextCondition{
			AuthType: d.Get("authtype").(string),
		},
		Network: buildPolicyNetworkCondition(d),
		People:  getUsers(d),
	}

	provider, ok := d.GetOk("identity_provider")
	if ok {
		template.Conditions.IdentityProvider = &sdk.IdentityProviderPolicyRuleCondition{
			Provider: provider.(string),
			IdpIds:   utils.ConvertInterfaceToStringArr(d.Get("identity_provider_ids")),
		}
	}

	bi, ok := d.GetOk("behaviors")
	if ok {
		template.Conditions.Risk = &sdk.RiskPolicyRuleCondition{
			Behaviors: utils.ConvertInterfaceToStringSetNullable(bi),
		}
	}
	riskLevel, riskLevelExists := d.GetOk("risk_level")
	riscLevel, riscLevelExists := d.GetOk("risc_level")
	if riskLevelExists {
		template.Conditions.RiskScore = &sdk.RiskScorePolicyRuleCondition{Level: riskLevel.(string)}
	} else if riscLevelExists {
		template.Conditions.RiskScore = &sdk.RiskScorePolicyRuleCondition{Level: riscLevel.(string)}
	}

	template.Actions = sdk.SdkPolicyRuleActions{
		SignOn: &sdk.SdkSignOnPolicyRuleSignOnActions{
			Access:                  d.Get("access").(string),
			FactorLifetime:          int64(d.Get("mfa_lifetime").(int)),
			FactorPromptMode:        d.Get("mfa_prompt").(string),
			RememberDeviceByDefault: utils.BoolPtr(d.Get("mfa_remember_device").(bool)),
			RequireFactor:           utils.BoolPtr(d.Get("mfa_required").(bool)),
			Session: &sdk.OktaSignOnPolicyRuleSignonSessionActions{
				MaxSessionIdleMinutesPtr:     utils.Int64Ptr(d.Get("session_idle").(int)),
				MaxSessionLifetimeMinutesPtr: utils.Int64Ptr(d.Get("session_lifetime").(int)),
				UsePersistentCookie:          utils.BoolPtr(d.Get("session_persistent").(bool)),
			},
		},
	}
	pf, ok := d.GetOk("primary_factor")
	if ok {
		template.Actions.SignOn.PrimaryFactor = pf.(string)
	}
	factorSeq := d.Get("factor_sequence").([]interface{})
	if len(factorSeq) == 0 {
		return template
	}
	template.Actions.SignOn.Challenge = &sdk.SdkSignOnPolicyRuleSignOnActionsChallenge{}
	chain := make([]sdk.SdkSignOnPolicyRuleSignOnActionsChallengeChain, len(factorSeq))
	for i := range factorSeq {
		chain[i] = sdk.SdkSignOnPolicyRuleSignOnActionsChallengeChain{
			Criteria: []sdk.SdkSignOnPolicyRuleSignOnActionsChallengeChainCriteria{
				{
					Provider:   d.Get(fmt.Sprintf("factor_sequence.%d.primary_criteria_provider", i)).(string),
					FactorType: d.Get(fmt.Sprintf("factor_sequence.%d.primary_criteria_factor_type", i)).(string),
				},
			},
		}
		secondaryCriteria := d.Get(fmt.Sprintf("factor_sequence.%d.secondary_criteria", i)).([]interface{})
		chain[i].Next = make([]sdk.SdkSignOnPolicyRuleSignOnActionsChallengeChainNext, 1)
		for j := range secondaryCriteria {
			chain[i].Next[0].Criteria = append(chain[i].Next[0].Criteria, sdk.SdkSignOnPolicyRuleSignOnActionsChallengeChainCriteria{
				Provider:   d.Get(fmt.Sprintf("factor_sequence.%d.secondary_criteria.%d.provider", i, j)).(string),
				FactorType: d.Get(fmt.Sprintf("factor_sequence.%d.secondary_criteria.%d.factor_type", i, j)).(string),
			})
		}
		if len(chain[i].Next[0].Criteria) == 0 {
			chain[i].Next = nil
		}
	}
	template.Actions.SignOn.Challenge.Chain = chain
	return template
}

func validateSignOnPolicyRule(d *schema.ResourceData) error {
	_, ok := d.GetOk("factor_sequence")
	isChallenge := d.Get("access").(string) == "CHALLENGE"
	if (!ok && isChallenge) || (ok && !isChallenge) {
		return errors.New("'factor_sequence' can only be set when access is 'CHALLENGE' and vice versa")
	}
	prompt, ok := d.GetOk("mfa_prompt")
	if !ok {
		return nil
	}
	if prompt.(string) != "DEVICE" {
		d, ok := d.GetOk("mfa_remember_device")
		if ok && d.(bool) {
			return errors.New("'mfa_remember_device' can only be set when mfa_prompt='DEVICE'")
		}
	}
	ip, ok := d.GetOk("identity_provider")
	if ok && ip == "SPECIFIC_IDP" && len(utils.ConvertInterfaceToStringArrNullable(d.Get("identity_provider_ids"))) < 1 {
		return errors.New("'identity_provider_ids' should have at least one element when 'identity_provider' is 'SPECIFIC_IDP'")
	}
	return nil
}
