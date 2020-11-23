package okta

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
)

func resourcePolicyMfaRule() *schema.Resource {
	return &schema.Resource{
		Exists:   resourcePolicyRuleExists,
		Create:   resourcePolicyMfaRuleCreate,
		Read:     resourcePolicyMfaRuleRead,
		Update:   resourcePolicyMfaRuleUpdate,
		Delete:   resourcePolicyMfaRuleDelete,
		Importer: createPolicyRuleImporter(),
		Schema: buildRuleSchema(map[string]*schema.Schema{
			"enroll": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"CHALLENGE", "LOGIN", "NEVER"}, false),
				Default:      "CHALLENGE",
				Optional:     true,
				Description:  "Should the user be enrolled the first time they LOGIN, the next time they are CHALLENGEd, or NEVER?",
			},
		}),
	}
}

func resourcePolicyMfaRuleCreate(d *schema.ResourceData, m interface{}) error {
	if err := ensureNotDefaultRule(d); err != nil {
		return err
	}
	template := buildMfaPolicyRule(d)
	rule, err := createRule(d, m, template, policyRulePassword)
	if err != nil {
		return err
	}
	d.SetId(rule.Id)
	err = validatePriority(template.Priority, rule.Priority)
	if err != nil {
		return err
	}
	return resourcePolicyMfaRuleRead(d, m)
}

func resourcePolicyMfaRuleRead(d *schema.ResourceData, m interface{}) error {
	rule, err := getPolicyRule(d, m)
	if rule == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}
	return syncRuleFromUpstream(d, rule)
}

func resourcePolicyMfaRuleUpdate(d *schema.ResourceData, m interface{}) error {
	if err := ensureNotDefaultRule(d); err != nil {
		return err
	}
	template := buildMfaPolicyRule(d)
	rule, err := updateRule(d, m, template)
	if err != nil {
		return err
	}
	err = validatePriority(template.Priority, rule.Priority)
	if err != nil {
		return err
	}
	return resourcePolicyMfaRuleRead(d, m)
}

func resourcePolicyMfaRuleDelete(d *schema.ResourceData, m interface{}) error {
	return deleteRule(d, m)
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
