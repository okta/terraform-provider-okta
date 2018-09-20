package okta

import (
	"fmt"
	"log"

	articulateOkta "github.com/articulate/oktasdk-go/okta"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourcePasswordPolicyRule() *schema.Resource {
	return &schema.Resource{
		Create: resourcePasswordPolicyRuleCreate,
		Read:   resourcePasswordPolicyRuleRead,
		Update: resourcePasswordPolicyRuleUpdate,
		Delete: resourcePasswordPolicyRuleDelete,

		CustomizeDiff: func(d *schema.ResourceDiff, v interface{}) error {
			// user cannot edit a default policy rule
			if d.Get("name").(string) == "Default Rule" {
				return fmt.Errorf("You cannot edit a default Policy Rule")
			}

			return nil
		},

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

func resourcePasswordPolicyRuleCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Creating Policy Rule %v", d.Get("name").(string))
	client := getClientFromMetadata(m)
	template := buildPasswordPolicyRule(d, client)
	rule, err := createRule(d, m, template)
	if err != nil {
		return err
	}

	err = validatePriority(template.Priority, rule.Priority)
	if err != nil {
		return err
	}
	d.SetId(rule.ID)

	return resourcePasswordPolicyRuleRead(d, m)
}

func resourcePasswordPolicyRuleRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] List Policy Rule %v", d.Get("name").(string))

	rule, err := getPolicyRule(d, m)
	if err != nil {
		return err
	}
	if rule == nil {
		// if the policy rule does not exist in okta, delete from terraform state
		d.SetId("")
		return nil
	}

	// Update with upstream state to prevent stale state
	d.Set("password_change", rule.Actions.PasswordChange.Access)
	d.Set("password_unlock", rule.Actions.SelfServiceUnlock.Access)
	d.Set("password_reset", rule.Actions.SelfServicePasswordReset.Access)

	return syncRuleFromUpstream(d, rule)
}

func resourcePasswordPolicyRuleUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Update Policy Rule %v", d.Get("name").(string))
	client := getClientFromMetadata(m)
	template := buildPasswordPolicyRule(d, client)

	rule, err := updateRule(d, m, template)
	if err != nil {
		return err
	}

	err = validatePriority(template.Priority, rule.Priority)
	if err != nil {
		return err
	}

	return resourcePasswordPolicyRuleRead(d, m)
}

func resourcePasswordPolicyRuleDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Delete Policy Rule %v", d.Get("name").(string))
	client := m.(*Config).articulateOktaClient

	rule, err := getPolicyRule(d, m)
	if err != nil {
		return err
	}
	if rule != nil && rule.ID != "" {
		if rule.System == true {
			log.Printf("[INFO] Policy Rule %v is a System Policy, cannot delete from Okta", d.Get("name").(string))
		} else {
			_, err = client.Policies.DeletePolicyRule(d.Get("policyid").(string), rule.ID)
			if err != nil {
				return fmt.Errorf("[ERROR] Error Deleting Policy Rule from Okta: %v", err)
			}
		}
	} else {
		log.Printf("[INFO] Policy Rule not found in Okta, removing from terraform")
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
func buildPasswordPolicyRule(d *schema.ResourceData, client *articulateOkta.Client) articulateOkta.PasswordRule {
	template := client.Policies.PasswordRule()
	template.Name = d.Get("name").(string)
	template.Type = passwordPolicyType
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
