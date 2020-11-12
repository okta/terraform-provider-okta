package okta

import (
	"fmt"
	"log"

	articulateOkta "github.com/articulate/oktasdk-go/okta"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
	client := getClientFromMetadata(m)
	template := buildpolicyRulePassword(d, client)
	rule, err := createRule(d, m, template, policyRulePassword)
	if err != nil {
		return err
	}

	// We want to put this under Terraform's control even if priority is invalid.
	d.SetId(rule.ID)
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
	client := getClientFromMetadata(m)
	template := buildpolicyRulePassword(d, client)

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
	if err := ensureNotDefaultRule(d); err != nil {
		return err
	}

	log.Printf("[INFO] Delete Policy Rule %v", d.Get("name").(string))
	client := getClientFromMetadata(m)

	_, err := client.Policies.DeletePolicyRule(d.Get("policyid").(string), d.Id())
	if err != nil {
		return fmt.Errorf("[ERROR] Error Deleting Policy Rule from Okta: %v", err)
	}

	// remove the policy rule resource from terraform
	d.SetId("")

	return nil
}

// activate or deactivate a policy rule according to the terraform schema status field
func policyRuleActivate(d *schema.ResourceData, m interface{}) error {
	client := getClientFromMetadata(m)

	if d.Get("status").(string) == "ACTIVE" {
		_, err := client.Policies.ActivatePolicyRule(d.Get("policyid").(string), d.Id())
		if err != nil {
			return fmt.Errorf("[ERROR] Error Activating Policy Rule: %v", err)
		}
	}
	if d.Get("status").(string) == "INACTIVE" {
		_, err := client.Policies.DeactivatePolicyRule(d.Get("policyid").(string), d.Id())
		if err != nil {
			return fmt.Errorf("[ERROR] Error Deactivating Policy Rule: %v", err)
		}
	}
	return nil
}

// build password policy rule from schema data
func buildpolicyRulePassword(d *schema.ResourceData, client *articulateOkta.Client) articulateOkta.PasswordRule {
	template := client.Policies.PasswordRule()
	template.Name = d.Get("name").(string)
	template.Status = d.Get("status").(string)
	if priority, ok := d.GetOk("priority"); ok {
		template.Priority = priority.(int)
	}

	template.Conditions = &articulateOkta.PolicyConditions{
		Network: getNetwork(d),
		People:  getUsers(d),
	}

	template.Actions.PasswordChange.Access = d.Get("password_change").(string)
	template.Actions.SelfServicePasswordReset.Access = d.Get("password_reset").(string)
	template.Actions.SelfServiceUnlock.Access = d.Get("password_unlock").(string)

	return template
}
