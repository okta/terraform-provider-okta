package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

func resourcePolicyProfileEnrollmentRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyProfileEnrollmentRuleCreate,
		ReadContext:   resourcePolicyProfileEnrollmentRuleRead,
		UpdateContext: resourcePolicyProfileEnrollmentRuleUpdate,
		DeleteContext: resourcePolicyProfileEnrollmentRuleDelete,
		Importer:      createPolicyRuleImporter(),
		Description: `Creates a Profile Enrollment Policy Rule.
		
~> **WARNING:** This feature is only available as a part of the Identity Engine. [Contact support](mailto:dev-inquiries@okta.com) for further information.
A [profile enrollment
policy](https://developer.okta.com/docs/reference/api/policy/#profile-enrollment-policy)
is limited to one default rule. This resource does not create a rule for an
enrollment policy, it allows the default policy rule to be updated.`,
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the policy",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the rule",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the rule",
			},
			"inline_hook_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of a Registration Inline Hook",
			},
			"target_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of a Group that this User should be added to",
			},
			"unknown_user_action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Which action should be taken if this User is new. Valid values are: `DENY`, `REGISTER`",
			},
			"ui_schema_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Value created by the backend. If present all policy updates must include this attribute/value.",
			},
			"email_verification": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether email verification should occur before access is granted. Default: `true`.",
				Default:     true,
			},
			"access": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Allow or deny access based on the rule conditions. Valid values are: `ALLOW`, `DENY`. Default: `ALLOW`.",
				Default:     "ALLOW",
			},
			"profile_attributes": {
				Type:     schema.TypeList,
				Optional: true,
				Description: `A list of attributes to prompt the user during registration or progressive profiling. Where defined on the User schema, these attributes are persisted in the User profile. Non-schema attributes may also be added, which aren't persisted to the User's profile, but are included in requests to the registration inline hook. A maximum of 10 Profile properties is supported.
	- 'label' - (Required) A display-friendly label for this property
	- 'name' - (Required) The name of a User Profile property
	- 'required' - (Required) Indicates if this property is required for enrollment. Default is 'false'.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"label": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "A display-friendly label for this property",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of a User Profile property",
						},
						"required": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates if this property is required for enrollment",
							Default:     false,
						},
					},
				},
			},
			"progressive_profiling_action": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enabled or disabled progressive profiling action rule conditions: `ENABLED` or `DISABLED`. Default: `DISABLED`",
				Default:     "DISABLED",
			},
			"enroll_authenticator_types": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Enrolls authenticator types",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

// resourcePolicyProfileEnrollmentRuleCreate
// True create does not exist for the rule of a profile enrollment policy.
// "This type of policy can only have one policy rule, so it's not possible to
// create other rules. Instead, consider editing the default one to meet your
// needs."
// https://developer.okta.com/docs/reference/api/policy/#profile-enrollment-policy
func resourcePolicyProfileEnrollmentRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if isClassicOrg(ctx, m) {
		return resourceOIEOnlyFeatureError(policyRuleProfileEnrollment)
	}

	policy, _, err := getAPISupplementFromMetadata(m).GetPolicy(ctx, d.Get("policy_id").(string))
	if err != nil {
		return diag.Errorf("failed to get profile enrollment policy: %v", err)
	}
	if policy.Type != sdk.ProfileEnrollmentPolicyType {
		return diag.Errorf("provided policy is not of type %s", sdk.ProfileEnrollmentPolicyType)
	}
	rules, _, err := getAPISupplementFromMetadata(m).ListPolicyRules(ctx, d.Get("policy_id").(string))
	if err != nil {
		return diag.Errorf("failed to get list profile enrollment policy rules: %v", err)
	}
	if len(rules) == 0 {
		return diag.Errorf("this policy should contain one default Catch-All rule, but it doesn't")
	}
	updateRule, err := buildPolicyRuleProfileEnrollment(ctx, m, d, rules[0].Id)
	if err != nil {
		return diag.Errorf("failed to prepare update of existing profile enrollment policy rule: %v", err)
	}
	rule, _, err := getAPISupplementFromMetadata(m).UpdatePolicyRule(ctx, d.Get("policy_id").(string), rules[0].Id, *updateRule)
	if err != nil {
		return diag.Errorf("failed to update existing profile enrollment policy rule: %v", err)
	}
	d.SetId(rule.Id)
	return resourcePolicyProfileEnrollmentRuleRead(ctx, d, m)
}

func resourcePolicyProfileEnrollmentRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if isClassicOrg(ctx, m) {
		return resourceOIEOnlyFeatureError(policyRuleProfileEnrollment)
	}

	rule, resp, err := getAPISupplementFromMetadata(m).GetPolicyRule(ctx, d.Get("policy_id").(string), d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get profile enrollment policy rule: %v", err)
	}
	if rule == nil || rule.Type != sdk.ProfileEnrollmentPolicyType {
		d.SetId("")
		return nil
	}
	_ = d.Set("status", rule.Status)
	_ = d.Set("name", rule.Name)
	if len(rule.Actions.ProfileEnrollment.PreRegistrationInlineHooks) != 0 {
		_ = d.Set("inline_hook_id", rule.Actions.ProfileEnrollment.PreRegistrationInlineHooks[0].InlineHookId)
	}
	if len(rule.Actions.ProfileEnrollment.TargetGroupIds) != 0 {
		_ = d.Set("target_group_id", rule.Actions.ProfileEnrollment.TargetGroupIds[0])
	}
	_ = d.Set("unknown_user_action", rule.Actions.ProfileEnrollment.UnknownUserAction)
	_ = d.Set("ui_schema_id", rule.Actions.ProfileEnrollment.UiSchemaId)
	_ = d.Set("email_verification", rule.Actions.ProfileEnrollment.ActivationRequirements.EmailVerification)
	_ = d.Set("access", rule.Actions.ProfileEnrollment.Access)
	arr := make([]map[string]interface{}, len(rule.Actions.ProfileEnrollment.ProfileAttributes))
	for i := range rule.Actions.ProfileEnrollment.ProfileAttributes {
		arr[i] = map[string]interface{}{
			"label":    rule.Actions.ProfileEnrollment.ProfileAttributes[i].Label,
			"name":     rule.Actions.ProfileEnrollment.ProfileAttributes[i].Name,
			"required": *rule.Actions.ProfileEnrollment.ProfileAttributes[i].Required,
		}
	}
	_ = d.Set("profile_attributes", arr)
	_ = d.Set("enroll_authenticator_types", convertStringSliceToSetNullable(rule.Actions.ProfileEnrollment.EnrollAuthenticatorTypes))
	return nil
}

func resourcePolicyProfileEnrollmentRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if isClassicOrg(ctx, m) {
		return resourceOIEOnlyFeatureError(policyRuleProfileEnrollment)
	}

	updateRule, err := buildPolicyRuleProfileEnrollment(ctx, m, d, d.Id())
	if err != nil {
		return diag.Errorf("failed to prepare update profile enrollment policy rule: %v", err)
	}
	_, _, err = getAPISupplementFromMetadata(m).UpdatePolicyRule(ctx, d.Get("policy_id").(string), d.Id(), *updateRule)
	if err != nil {
		return diag.Errorf("failed to update profile enrollment policy rule: %v", err)
	}
	return resourcePolicyProfileEnrollmentRuleRead(ctx, d, m)
}

// You cannot delete a default rule in a policy
func resourcePolicyProfileEnrollmentRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if isClassicOrg(ctx, m) {
		return resourceOIEOnlyFeatureError(policyRuleProfileEnrollment)
	}

	return nil
}

// buildPolicyRuleProfileEnrollment build profile enrollment policy rule from
// copy of existing rule
func buildPolicyRuleProfileEnrollment(ctx context.Context, m interface{}, d *schema.ResourceData, id string) (*sdk.SdkPolicyRule, error) {
	rule, resp, err := getAPISupplementFromMetadata(m).GetPolicyRule(ctx, d.Get("policy_id").(string), id)
	if err = suppressErrorOn404(resp, err); err != nil {
		return nil, err
	}

	// Given that the Okta API is sensitive to attributes that are already set
	// on the one rule always present on a profile enrollment policy:
	// 1. get a current copy of the rule
	// 2. use the copy to prepopulate vaules for the update rule
	// 3. apply values from the resource config to the update rule

	// First, API requires these attributes always be set, prepopulate with
	// existing values
	updateRule := sdk.ProfileEnrollmentPolicyRule()
	updateRule.Id = rule.Id
	updateRule.Name = rule.Name
	updateRule.Priority = rule.Priority
	updateRule.System = rule.System
	updateRule.Status = rule.Status

	// Additionally, only prepopulate the attributes that are already set on the
	// current rule

	ruleAction := sdk.NewProfileEnrollmentPolicyRuleAction()

	// access
	if rule.Actions.ProfileEnrollment.Access != "" {
		ruleAction.Access = rule.Actions.ProfileEnrollment.Access
	}
	if access, ok := d.GetOk("access"); ok {
		ruleAction.Access = access.(string)
	}
	if progressiveProfilingAction, ok := d.GetOk("progressive_profiling_action"); ok {
		ruleAction.ProgressiveProfilingAction = progressiveProfilingAction.(string)
	}

	activationRequirements := sdk.NewProfileEnrollmentPolicyRuleActivationRequirement()
	// email_verification is set on the config use it's value, else fallback to the rule copy
	if ev, _ := d.GetOk("email_verification"); ev != nil {
		activationRequirements.EmailVerification = boolPtr(ev.(bool))
	} else {
		activationRequirements.EmailVerification = ruleAction.ActivationRequirements.EmailVerification
	}

	// unknown_user_action
	ruleAction.ActivationRequirements = activationRequirements
	if rule.Actions.ProfileEnrollment.UnknownUserAction != "" {
		ruleAction.UnknownUserAction = rule.Actions.ProfileEnrollment.UnknownUserAction
	}
	if uua, ok := d.GetOk("unknown_user_action"); ok {
		ruleAction.UnknownUserAction = uua.(string)
	}

	// ui_schema_id
	if rule.Actions.ProfileEnrollment.UiSchemaId != "" {
		ruleAction.UiSchemaId = rule.Actions.ProfileEnrollment.UiSchemaId
	}
	if usi, ok := d.GetOk("ui_schema_id"); ok {
		ruleAction.UiSchemaId = usi.(string)
	}

	if eat, ok := d.GetOk("enroll_authenticator_types"); ok {
		ruleAction.EnrollAuthenticatorTypes = convertInterfaceToStringSetNullable(eat)
	}

	updateRule.Actions = sdk.SdkPolicyRuleActions{
		ProfileEnrollment: ruleAction,
	}

	// inline_hook_id
	if len(rule.Actions.ProfileEnrollment.PreRegistrationInlineHooks) != 0 {
		updateRule.Actions.ProfileEnrollment.PreRegistrationInlineHooks = rule.Actions.ProfileEnrollment.PreRegistrationInlineHooks
	}
	if hook, ok := d.GetOk("inline_hook_id"); ok {
		updateRule.Actions.ProfileEnrollment.PreRegistrationInlineHooks = []*sdk.PreRegistrationInlineHook{{InlineHookId: hook.(string)}}
	}

	// target_group_id
	if len(rule.Actions.ProfileEnrollment.TargetGroupIds) != 0 {
		updateRule.Actions.ProfileEnrollment.TargetGroupIds = rule.Actions.ProfileEnrollment.TargetGroupIds
	}
	if targetGroup, ok := d.GetOk("target_group_id"); ok {
		updateRule.Actions.ProfileEnrollment.TargetGroupIds = []string{targetGroup.(string)}
	}

	pa, ok := d.GetOk("profile_attributes")
	if !ok {
		return &updateRule, nil
	}

	attributes := make([]*sdk.ProfileEnrollmentPolicyRuleProfileAttribute, len(pa.([]interface{})))
	for i := range pa.([]interface{}) {
		attributes[i] = &sdk.ProfileEnrollmentPolicyRuleProfileAttribute{
			Label:    d.Get(fmt.Sprintf("profile_attributes.%d.label", i)).(string),
			Name:     d.Get(fmt.Sprintf("profile_attributes.%d.name", i)).(string),
			Required: boolPtr(d.Get(fmt.Sprintf("profile_attributes.%d.required", i)).(bool)),
		}
	}
	updateRule.Actions.ProfileEnrollment.ProfileAttributes = attributes

	return &updateRule, nil
}
