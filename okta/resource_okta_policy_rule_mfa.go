package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
)

func resourcePolicyMfaRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyMfaRuleCreate,
		ReadContext:   resourcePolicyMfaRuleRead,
		UpdateContext: resourcePolicyMfaRuleUpdate,
		DeleteContext: resourcePolicyMfaRuleDelete,
		Importer:      createPolicyRuleImporter(),
		Schema: buildRuleSchema(map[string]*schema.Schema{
			"enroll": {
				Type:             schema.TypeString,
				ValidateDiagFunc: stringInSlice([]string{"CHALLENGE", "LOGIN", "NEVER"}),
				Default:          "CHALLENGE",
				Optional:         true,
				Description:      "Should the user be enrolled the first time they LOGIN, the next time they are CHALLENGEd, or NEVER?",
			},
		}),
	}
}

func resourcePolicyMfaRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	template := buildMfaPolicyRule(d)
	err := createRule(ctx, d, m, template, policyRulePassword)
	if err != nil {
		return diag.Errorf("failed to create MFA policy rule: %v", err)
	}
	return resourcePolicyMfaRuleRead(ctx, d, m)
}

func resourcePolicyMfaRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	rule, err := getPolicyRule(ctx, d, m)
	if err != nil {
		return diag.Errorf("failed to get MFA policy rule: %v", err)
	}
	if rule == nil {
		return nil
	}
	err = syncRuleFromUpstream(d, rule)
	if err != nil {
		return diag.Errorf("failed to sync MFA policy rule: %v", err)
	}
	return nil
}

func resourcePolicyMfaRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	template := buildMfaPolicyRule(d)
	err := updateRule(ctx, d, m, template)
	if err != nil {
		return diag.Errorf("failed to update MFA policy rule: %v", err)
	}
	return resourcePolicyMfaRuleRead(ctx, d, m)
}

func resourcePolicyMfaRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := deleteRule(ctx, d, m, false)
	if err != nil {
		return diag.Errorf("failed to delete MFA policy rule: %v", err)
	}
	return nil
}

// build password policy rule from schema data
func buildMfaPolicyRule(d *schema.ResourceData) sdk.PolicyRule {
	rule := sdk.MfaPolicyRule()
	rule.Name = d.Get("name").(string)
	rule.Status = d.Get("status").(string)
	if priority, ok := d.GetOk("priority"); ok {
		rule.Priority = int64(priority.(int))
	}
	rule.Conditions = &okta.PolicyRuleConditions{
		Network: getNetwork(d),
		People:  getUsers(d),
	}
	if enroll, ok := d.GetOk("enroll"); ok {
		rule.Actions = sdk.PolicyRuleActions{
			Enroll: &sdk.Enroll{
				Self: enroll.(string),
			},
		}
	}
	return rule
}
