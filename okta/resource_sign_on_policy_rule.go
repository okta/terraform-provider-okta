package okta

import (
	"fmt"
	"log"

	articulateOkta "github.com/articulate/oktasdk-go/okta"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceSignOnPolicyRule() *schema.Resource {
	return &schema.Resource{
		Exists:   resourcePolicyRuleExists,
		Create:   resourceSignOnPolicyRuleCreate,
		Read:     resourceSignOnPolicyRuleRead,
		Update:   resourceSignOnPolicyRuleUpdate,
		Delete:   resourceSignOnPolicyRuleDelete,
		Importer: createPolicyRuleImporter(),

		CustomizeDiff: func(d *schema.ResourceDiff, v interface{}) error {
			// user cannot edit a default policy rule
			if d.Get("name").(string) == "Default Rule" {
				return fmt.Errorf("You cannot edit a default Policy Rule")
			}

			return nil
		},

		Schema: buildRuleSchema(map[string]*schema.Schema{
			"authtype": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ANY", "RADIUS"}, false),
				Description:  "Authentication entrypoint: ANY or RADIUS.",
				Default:      "ANY",
			},
			"access": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ALLOW", "DENY"}, false),
				Description:  "Allow or deny access based on the rule conditions: ALLOW or DENY.",
				Default:      "ALLOW",
			},
			"mfa_required": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Require MFA.",
				Default:     false,
			},
			"mfa_prompt": { // mfa_require must be true
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"DEVICE", "SESSION", "ALWAYS"}, false),
				Description:  "Prompt for MFA based on the device used, a factor session lifetime, or every sign on attempt: DEVICE, SESSION or ALWAYS",
			},
			"mfa_remember_device": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Remember MFA device.",
				Default:     false,
			},
			"mfa_lifetime": { // mfa_require must be true, mfaprompt must be SESSION
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
		}),
	}
}

func resourceSignOnPolicyRuleCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Creating Policy Rule %v", d.Get("name").(string))
	template, err := buildSignOnPolicyRule(d, m)
	if err != nil {
		return err
	}

	rule, err := createRule(d, m, template, signOnPolicyRule)
	if err != nil {
		return err
	}

	// We want to put this under Terraform's control even if priority is invalid.
	d.SetId(rule.ID)
	err = validatePriority(template.Priority, rule.Priority)
	if err != nil {
		return err
	}

	return resourceSignOnPolicyRuleRead(d, m)
}

func resourceSignOnPolicyRuleRead(d *schema.ResourceData, m interface{}) error {
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
	d.Set("authtype", rule.Conditions.AuthContext.AuthType)
	d.Set("access", rule.Actions.SignOn.Access)
	d.Set("mfa_required", rule.Actions.SignOn.RequireFactor)
	d.Set("mfa_remember_device", rule.Actions.SignOn.RememberDeviceByDefault)
	d.Set("mfa_lifetime", rule.Actions.SignOn.FactorLifetime)
	d.Set("session_idle", rule.Actions.SignOn.Session.MaxSessionIdleMinutes)
	d.Set("session_lifetime", rule.Actions.SignOn.Session.MaxSessionLifetimeMinutes)
	d.Set("session_persistent", rule.Actions.SignOn.Session.UsePersistentCookie)

	if rule.Actions.FactorPromptMode != "" {
		d.Set("mfa_prompt", rule.Actions.FactorPromptMode)
	}

	return syncRuleFromUpstream(d, rule)
}

func resourceSignOnPolicyRuleUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Update Policy Rule %v", d.Get("name").(string))
	template, err := buildSignOnPolicyRule(d, m)
	if err != nil {
		return err
	}

	rule, err := updateRule(d, m, template)
	if err != nil {
		return err
	}

	err = validatePriority(template.Priority, rule.Priority)
	if err != nil {
		return err
	}

	return resourceSignOnPolicyRuleRead(d, m)
}

func resourceSignOnPolicyRuleDelete(d *schema.ResourceData, m interface{}) error {
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

// Build Policy Sign On Rule from resource data
func buildSignOnPolicyRule(d *schema.ResourceData, m interface{}) (articulateOkta.SignOnRule, error) {
	client := getClientFromMetadata(m)
	template := client.Policies.SignOnRule()
	template.Name = d.Get("name").(string)
	template.Status = d.Get("status").(string)
	template.Type = singOnPolicyRuleType
	if priority, ok := d.GetOk("priority"); ok {
		template.Priority = priority.(int)
	}

	template.Conditions = &articulateOkta.PolicyConditions{
		Network: getNetwork(d),
		AuthContext: &articulateOkta.AuthContext{
			AuthType: d.Get("authtype").(string),
		},
		People: getUsers(d),
	}

	// Hardcoded since MFA is not implemented.
	template.Actions.SignOn.RequireFactor = false
	template.Actions.SignOn.Session.MaxSessionIdleMinutes = d.Get("session_idle").(int)
	template.Actions.SignOn.Session.MaxSessionLifetimeMinutes = d.Get("session_lifetime").(int)
	template.Actions.SignOn.Session.UsePersistentCookie = d.Get("session_persistent").(bool)
	template.Actions.SignOn.Access = d.Get("access").(string)

	// Preserving existing errors here, looks like the MFA rule needs to be there in order for these to work.
	//if required, ok := d.GetOk("mfa_required"); ok {
	if _, ok := d.GetOk("mfa_required"); ok {
		return template, fmt.Errorf("[ERROR] mfa signon actions not supported in this terraform provider at this time")
	}
	//if prompt, ok := d.GetOk("mfa_prompt"); ok {
	if _, ok := d.GetOk("mfa_prompt"); ok {
		return template, fmt.Errorf("[ERROR] mfa signon actions not supported in this terraform provider at this time")
	}
	//if remember, ok := d.GetOk("mfa_remember_device"); ok {
	if _, ok := d.GetOk("mfa_remember_device"); ok {
		return template, fmt.Errorf("[ERROR] mfa signon actions not supported in this terraform provider at this time")
	}
	//if lifetime, ok := d.GetOk("mfa_lifetime"); ok {
	if _, ok := d.GetOk("mfa_lifetime"); ok {
		return template, fmt.Errorf("[ERROR] mfa signon actions not supported in this terraform provider at this time")
	}

	return template, nil
}
