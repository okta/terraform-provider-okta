package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/terraform-provider-okta/sdk"
)

func resourcePolicyProfileEnrollmentRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyProfileEnrollmentRuleCreate,
		ReadContext:   resourcePolicyProfileEnrollmentRuleRead,
		UpdateContext: resourcePolicyProfileEnrollmentRuleUpdate,
		DeleteContext: resourcePolicyProfileEnrollmentRuleDelete,
		Importer:      createPolicyRuleImporter(),
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
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: elemInSlice([]string{"DENY", "REGISTER"}),
				Description:      "Which action should be taken if this User is new",
			},
			"email_verification": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether email verification should occur before access is granted",
				Default:     true,
			},
			"access": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: elemInSlice([]string{"ALLOW", "DENY"}),
				Description:      "Allow or deny access based on the rule conditions: ALLOW or DENY",
				Default:          "ALLOW",
			},
			"profile_attributes": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of attributes to prompt the user during registration or progressive profiling",
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
		},
	}
}

func resourcePolicyProfileEnrollmentRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	policy, _, err := getSupplementFromMetadata(m).GetPolicy(ctx, d.Get("policy_id").(string))
	if err != nil {
		return diag.Errorf("failed to get profile enrollment policy: %v", err)
	}
	if policy.Type != sdk.ProfileEnrollmentPolicyType {
		return diag.Errorf("provided policy is not of type %s", sdk.ProfileEnrollmentPolicyType)
	}
	rules, _, err := getSupplementFromMetadata(m).ListPolicyRules(ctx, d.Get("policy_id").(string))
	if err != nil {
		return diag.Errorf("failed to get list profile enrollment policy rules: %v", err)
	}
	if len(rules) == 0 {
		return diag.Errorf("this policy should contain one default Catch-All rule, but it doesn't")
	}
	rule, _, err := getSupplementFromMetadata(m).UpdatePolicyRule(ctx, d.Get("policy_id").(string), rules[0].Id, buildPolicyRuleProfileEnrollment(d, rules[0].Id))
	if err != nil {
		return diag.Errorf("failed to create profile enrollment policy rule: %v", err)
	}
	d.SetId(rule.Id)
	return resourcePolicyProfileEnrollmentRuleRead(ctx, d, m)
}

func resourcePolicyProfileEnrollmentRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	rule, resp, err := getSupplementFromMetadata(m).GetPolicyRule(ctx, d.Get("policy_id").(string), d.Id())
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
	_ = d.Set("email_verification", *rule.Actions.ProfileEnrollment.ActivationRequirements.EmailVerification)
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
	return nil
}

func resourcePolicyProfileEnrollmentRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	_, _, err := getSupplementFromMetadata(m).UpdatePolicyRule(ctx, d.Get("policy_id").(string), d.Id(), buildPolicyRuleProfileEnrollment(d, d.Id()))
	if err != nil {
		return diag.Errorf("failed to update profile enrollment policy rule: %v", err)
	}
	return resourcePolicyProfileEnrollmentRuleRead(ctx, d, m)
}

//  You cannot delete a default rule in a policy
func resourcePolicyProfileEnrollmentRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

// build profile enrollment policy rule from schema data
func buildPolicyRuleProfileEnrollment(d *schema.ResourceData, id string) sdk.PolicyRule {
	rule := sdk.ProfileEnrollmentPolicyRule()
	rule.Id = id
	rule.Name = "Catch-all Rule" // read-only
	rule.Priority = 99           // read-only
	rule.System = boolPtr(true)  // read-only
	rule.Status = statusActive
	rule.Actions = sdk.PolicyRuleActions{
		ProfileEnrollment: &okta.ProfileEnrollmentPolicyRuleAction{
			Access: d.Get("access").(string),
			ActivationRequirements: &okta.ProfileEnrollmentPolicyRuleActivationRequirement{
				EmailVerification: boolPtr(d.Get("email_verification").(bool)),
			},
			UnknownUserAction: d.Get("unknown_user_action").(string),
		},
	}
	hook, ok := d.GetOk("inline_hook_id")
	if ok {
		rule.Actions.ProfileEnrollment.PreRegistrationInlineHooks = []*okta.PreRegistrationInlineHook{{InlineHookId: hook.(string)}}
	}
	targetGroup, ok := d.GetOk("target_group_id")
	if ok {
		rule.Actions.ProfileEnrollment.TargetGroupIds = []string{targetGroup.(string)}
	}
	pa, ok := d.GetOk("profile_attributes")
	if !ok {
		return rule
	}
	attributes := make([]*okta.ProfileEnrollmentPolicyRuleProfileAttribute, len(pa.([]interface{})))
	for i := range pa.([]interface{}) {
		attributes[i] = &okta.ProfileEnrollmentPolicyRuleProfileAttribute{
			Label:    d.Get(fmt.Sprintf("profile_attributes.%d.label", i)).(string),
			Name:     d.Get(fmt.Sprintf("profile_attributes.%d.name", i)).(string),
			Required: boolPtr(d.Get(fmt.Sprintf("profile_attributes.%d.required", i)).(bool)),
		}
	}
	rule.Actions.ProfileEnrollment.ProfileAttributes = attributes
	return rule
}
