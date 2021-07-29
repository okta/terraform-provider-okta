package okta

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/terraform-provider-okta/sdk"
)

func resourcePolicySignOnRule() *schema.Resource {
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
				ValidateDiagFunc: elemInSlice([]string{"ANY", "RADIUS", "LDAP_INTERFACE"}),
				Description:      "Authentication entrypoint: ANY, RADIUS or LDAP_INTERFACE",
				Default:          "ANY",
			},
			"access": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: elemInSlice([]string{"ALLOW", "DENY", "CHALLENGE"}),
				Description:      "Allow or deny access based on the rule conditions: ALLOW, DENY or CHALLENGE.",
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
				ValidateDiagFunc: elemInSlice([]string{"DEVICE", "SESSION", "ALWAYS"}),
				Description:      "Prompt for MFA based on the device used, a factor session lifetime, or every sign-on attempt: DEVICE, SESSION or ALWAYS",
			},
			"mfa_remember_device": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Remember MFA device.",
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
			"risc_level": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: elemInSlice([]string{"", "ANY", "LOW", "MEDIUM", "HIGH"}),
				Description:      "Risc level: ANY, LOW, MEDIUM or HIGH",
				Default:          "ANY",
			},
			"behaviors": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of behavior IDs",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"factor_sequence": {
				Type:     schema.TypeList,
				Optional: true,
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
		}),
	}
}

func resourcePolicySignOnRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := validateSignOnPolicyRule(d)
	if err != nil {
		return diag.FromErr(err)
	}
	template := buildSignOnPolicyRule(d)
	err = createRule(ctx, d, m, template, policyRuleSignOn)
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
	_ = d.Set("access", rule.Actions.SignOn.Access)
	_ = d.Set("mfa_required", rule.Actions.SignOn.RequireFactor)
	_ = d.Set("mfa_remember_device", rule.Actions.SignOn.RememberDeviceByDefault)
	_ = d.Set("mfa_lifetime", rule.Actions.SignOn.FactorLifetime)
	_ = d.Set("session_idle", rule.Actions.SignOn.Session.MaxSessionIdleMinutes)
	_ = d.Set("session_lifetime", rule.Actions.SignOn.Session.MaxSessionLifetimeMinutes)
	_ = d.Set("session_persistent", rule.Actions.SignOn.Session.UsePersistentCookie)
	if rule.Actions.SignOn.FactorPromptMode != "" {
		_ = d.Set("mfa_prompt", rule.Actions.SignOn.FactorPromptMode)
	}
	if rule.Conditions != nil {
		if rule.Conditions.RiskScore != nil {
			_ = d.Set("risc_level", rule.Conditions.RiskScore.Level)
		}
		if rule.Conditions.Risk != nil {
			err = setNonPrimitives(d, map[string]interface{}{
				"behaviors": convertStringSetToInterface(rule.Conditions.Risk.Behaviors),
			})
			if err != nil {
				return diag.Errorf("failed to set sign-on policy rule behaviors: %v", err)
			}
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
		err = setNonPrimitives(d, map[string]interface{}{"factor_sequence": arr})
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

func resourcePolicySignOnRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := validateSignOnPolicyRule(d)
	if err != nil {
		return diag.FromErr(err)
	}
	template := buildSignOnPolicyRule(d)
	err = updateRule(ctx, d, m, template)
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
		AuthContext: &okta.PolicyRuleAuthContextCondition{
			AuthType: d.Get("authtype").(string),
		},
		Network: getNetwork(d),
		People:  getUsers(d),
	}
	bi, ok := d.GetOk("behaviors")
	if ok {
		template.Conditions.Risk = &okta.RiskPolicyRuleCondition{
			Behaviors: convertInterfaceToStringSetNullable(bi),
		}
	}
	ri, ok := d.GetOk("risc_level")
	if ok {
		template.Conditions.RiskScore = &okta.RiskScorePolicyRuleCondition{
			Level: ri.(string),
		}
	}
	template.Actions = sdk.PolicyRuleActions{
		SignOn: &sdk.SignOnPolicyRuleSignOnActions{
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
	}
	factorSeq := d.Get("factor_sequence").([]interface{})
	if len(factorSeq) == 0 {
		return template
	}
	template.Actions.SignOn.Challenge = &sdk.SignOnPolicyRuleSignOnActionsChallenge{}
	chain := make([]sdk.SignOnPolicyRuleSignOnActionsChallengeChain, len(factorSeq))
	for i := range factorSeq {
		chain[i] = sdk.SignOnPolicyRuleSignOnActionsChallengeChain{
			Criteria: []sdk.SignOnPolicyRuleSignOnActionsChallengeChainCriteria{
				{
					Provider:   d.Get(fmt.Sprintf("factor_sequence.%d.primary_criteria_provider", i)).(string),
					FactorType: d.Get(fmt.Sprintf("factor_sequence.%d.primary_criteria_factor_type", i)).(string),
				},
			},
		}
		secondaryCriteria := d.Get(fmt.Sprintf("factor_sequence.%d.secondary_criteria", i)).([]interface{})
		chain[i].Next = make([]sdk.SignOnPolicyRuleSignOnActionsChallengeChainNext, 1)
		for j := range secondaryCriteria {
			chain[i].Next[0].Criteria = append(chain[i].Next[0].Criteria, sdk.SignOnPolicyRuleSignOnActionsChallengeChainCriteria{
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
	return nil
}
