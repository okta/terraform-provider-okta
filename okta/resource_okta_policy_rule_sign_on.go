package okta

import (
	"fmt"
	"log"

	articulateOkta "github.com/articulate/oktasdk-go/okta"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourcePolicySignonRule() *schema.Resource {
	return &schema.Resource{
		Exists:   resourcePolicyRuleExists,
		Create:   resourcePolicySignonRuleCreate,
		Read:     resourcePolicySignonRuleRead,
		Update:   resourcePolicySignonRuleUpdate,
		Delete:   resourcePolicySignonRuleDelete,
		Importer: createPolicyRuleImporter(),

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

func resourcePolicySignonRuleCreate(d *schema.ResourceData, m interface{}) error {
	if err := ensureNotDefaultRule(d); err != nil {
		return err
	}

	log.Printf("[INFO] Creating Policy Rule %v", d.Get("name").(string))
	template, err := buildSignOnPolicyRule(d, m)
	if err != nil {
		return err
	}

	rule, err := createRule(d, m, template, policyRuleSignOn)
	if err != nil {
		return err
	}

	// We want to put this under Terraform's control even if priority is invalid.
	d.SetId(rule.ID)
	err = validatePriority(template.Priority, rule.Priority)
	if err != nil {
		return err
	}

	return resourcePolicySignonRuleRead(d, m)
}

func resourcePolicySignonRuleRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] List Policy Rule %v", d.Get("name").(string))

	rule, err := getPolicyRule(d, m)

	if rule == nil {
		// if the policy rule does not exist in okta, delete from terraform state
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
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

	if rule.Actions.FactorPromptMode != "" {
		_ = d.Set("mfa_prompt", rule.Actions.FactorPromptMode)
	}

	return syncRuleFromUpstream(d, rule)
}

func resourcePolicySignonRuleUpdate(d *schema.ResourceData, m interface{}) error {
	if err := ensureNotDefaultRule(d); err != nil {
		return err
	}

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

	return resourcePolicySignonRuleRead(d, m)
}

func resourcePolicySignonRuleDelete(d *schema.ResourceData, m interface{}) error {
	if err := ensureNotDefaultRule(d); err != nil {
		return err
	}

	log.Printf("[INFO] Delete Policy Rule %v", d.Get("name").(string))
	client := m.(*Config).articulateOktaClient

	rule, err := getPolicyRule(d, m)

	if err != nil {
		return err
	}

	if rule != nil && rule.ID != "" {
		if rule.System {
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

	template.Actions.SignOn.RequireFactor = d.Get("mfa_required").(bool)
	template.Actions.SignOn.FactorPromptMode = d.Get("mfa_prompt").(string)
	template.Actions.SignOn.RememberDeviceByDefault = d.Get("mfa_remember_device").(bool)
	template.Actions.SignOn.FactorLifetime = d.Get("mfa_lifetime").(int)
	template.Actions.SignOn.Session.MaxSessionIdleMinutes = d.Get("session_idle").(int)
	template.Actions.SignOn.Session.MaxSessionLifetimeMinutes = d.Get("session_lifetime").(int)
	template.Actions.SignOn.Session.UsePersistentCookie = d.Get("session_persistent").(bool)
	template.Actions.SignOn.Access = d.Get("access").(string)

	return template, nil
}
