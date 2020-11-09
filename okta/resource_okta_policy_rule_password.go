package okta

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
)

func resourcePolicyPasswordRule() *schema.Resource {
	return &schema.Resource{
		Exists:   resourcePolicyRuleExists,
		Create:   resourcePolicyPasswordRuleCreate,
		Read:     resourcePolicyPasswordRuleRead,
		Update:   resourcePolicyPasswordRuleUpdate,
		Delete:   resourcePolicyPasswordRuleDelete,
		Importer: createPolicyRuleImporter(),

		Schema: buildRuleSchema(map[string]*schema.Schema{
			"password_change": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ALLOW", "DENY"}, false),
				Description:  "Allow or deny a user to change their password: ALLOW or DENY. Default = ALLOW",
				Default:      "ALLOW",
			},
			"password_reset": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ALLOW", "DENY"}, false),
				Description:  "Allow or deny a user to reset their password: ALLOW or DENY. Default = ALLOW",
				Default:      "ALLOW",
			},
			"password_unlock": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ALLOW", "DENY"}, false),
				Description:  "Allow or deny a user to unlock. Default = DENY",
				Default:      "DENY",
			},
		}),
	}
}

func resourcePolicyPasswordRuleCreate(d *schema.ResourceData, m interface{}) error {
	if err := ensureNotDefaultRule(d); err != nil {
		return err
	}

	log.Printf("[INFO] Creating Policy Rule %v", d.Get("name").(string))
	template := buildPolicyRulePassword(d)
	rule, err := createRule(d, m, template, policyRulePassword)
	if err != nil {
		return err
	}

	// We want to put this under Terraform's control even if priority is invalid.
	d.SetId(rule.Id)
	err = validatePriority(template.Priority, rule.Priority)
	if err != nil {
		return err
	}

	return resourcePolicyPasswordRuleRead(d, m)
}

func resourcePolicyPasswordRuleRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] List Policy Rule %v", d.Get("name").(string))

	rule, err := getPolicyRule(d, m)

	if rule == nil {
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	// Update with upstream state to prevent stale state
	_ = d.Set("password_change", rule.Actions.PasswordChange.Access)
	_ = d.Set("password_unlock", rule.Actions.SelfServiceUnlock.Access)
	_ = d.Set("password_reset", rule.Actions.SelfServicePasswordReset.Access)

	return syncRuleFromUpstream(d, rule)
}

func resourcePolicyPasswordRuleUpdate(d *schema.ResourceData, m interface{}) error {
	if err := ensureNotDefaultRule(d); err != nil {
		return err
	}

	log.Printf("[INFO] Update Policy Rule %v", d.Get("name").(string))
	template := buildPolicyRulePassword(d)

	rule, err := updateRule(d, m, template)
	if err != nil {
		return err
	}

	err = validatePriority(template.Priority, rule.Priority)
	if err != nil {
		return err
	}

	return resourcePolicyPasswordRuleRead(d, m)
}

func resourcePolicyPasswordRuleDelete(d *schema.ResourceData, m interface{}) error {
	return deleteRule(d, m)
}

// activate or deactivate a policy rule according to the terraform schema status field
func policyRuleActivate(d *schema.ResourceData, m interface{}) error {
	client := getSupplementFromMetadata(m)

	if d.Get("status").(string) == "ACTIVE" {
		_, err := client.ActivatePolicyRule(context.Background(), d.Get("policyid").(string), d.Id())
		if err != nil {
			return fmt.Errorf("failed to activate policy rule: %v", err)
		}
	}
	if d.Get("status").(string) == "INACTIVE" {
		_, err := client.DeactivatePolicyRule(context.Background(), d.Get("policyid").(string), d.Id())
		if err != nil {
			return fmt.Errorf("failed to deactivate policy rule: %v", err)
		}
	}
	return nil
}

// build password policy rule from schema data
func buildPolicyRulePassword(d *schema.ResourceData) sdk.PolicyRule {
	template := sdk.PasswordPolicyRule()
	template.Name = d.Get("name").(string)
	template.Status = d.Get("status").(string)
	if priority, ok := d.GetOk("priority"); ok {
		template.Priority = int64(priority.(int))
	}
	template.Conditions = &okta.PolicyRuleConditions{
		Network: getNetwork(d),
		People:  getUsers(d),
	}
	template.Actions = sdk.PolicyRuleActions{
		PasswordPolicyRuleActions: &okta.PasswordPolicyRuleActions{
			PasswordChange: &okta.PasswordPolicyRuleAction{
				Access: d.Get("password_change").(string),
			},
			SelfServicePasswordReset: &okta.PasswordPolicyRuleAction{
				Access: d.Get("password_reset").(string),
			},
			SelfServiceUnlock: &okta.PasswordPolicyRuleAction{
				Access: d.Get("password_unlock").(string),
			},
		},
	}
	return template
}
