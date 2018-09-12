package okta

import (
	"fmt"
	"log"

	articulateOkta "github.com/articulate/oktasdk-go/okta"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

// global var to determine if our policy rule is a system policy rule
var systemPolicyRule bool = false

func resourcePolicyRules() *schema.Resource {
	return &schema.Resource{
		Create: resourcePolicyRuleCreate,
		Read:   resourcePolicyRuleRead,
		Update: resourcePolicyRuleUpdate,
		Delete: resourcePolicyRuleDelete,

		CustomizeDiff: func(d *schema.ResourceDiff, v interface{}) error {
			// user cannot edit a default policy rule
			if d.Get("name").(string) == "Default Rule" {
				return fmt.Errorf("You cannot edit a default Policy Rule")
			}

			// user cannot change policyid, name, or type for an existing policy
			prev, _ := d.GetChange("name")
			if prev.(string) != "" {
				if d.HasChange("policyid") || d.HasChange("type") || d.HasChange("name") {
					return fmt.Errorf("You cannot change the policyid field, name field, or type field of an existing Policy Rule")
				}
			}

			// add custom error messages if user supplies options not supported by the policy rule
			switch d.Get("type").(string) {
			case "PASSWORD":
				if len(d.Get("conditions.0.authtype").(string)) > 0 {
					return fmt.Errorf("authtype condition options not supported in the Okta Password Policy Rule")
				}
				if len(d.Get("actions.0.signon").([]interface{})) > 0 {
					return fmt.Errorf("signon action options not supported in the Okta Password Policy Rule")
				}

			case "SIGN_ON":
				if len(d.Get("actions.0.password").([]interface{})) > 0 {
					return fmt.Errorf("password action options not supported in the Okta SignOn Policy Rule")
				}

			case "MFA_ENROLL":

			}

			return nil
		},

		Schema: map[string]*schema.Schema{
			"policyid": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Policy ID of the Rule",
			},
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"SIGN_ON", "PASSWORD", "MFA_ENROLL"}, false),
				Description:  "Policy Rule Type: SIGN_ON, PASSWORD, or MFA_ENROLL",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Policy Rule Name",
			},
			"priority": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Policy Rule Priority",
			},
			"status": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ACTIVE",
				ValidateFunc: validation.StringInSlice([]string{"ACTIVE", "INACTIVE"}, false),
				Description:  "Policy Rule Status: ACTIVE or INACTIVE. Default = ACTIVE",
			},
			"conditions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Conditions that must be met during Policy Evaluation",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"users": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of User IDs to Exclude",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"network": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Network selection mode & a set of network zones to be included or excluded",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"connection": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"ANYWHERE", "ZONE", "ON_NETWORK", "OFF_NETWORK"}, false),
										Description:  "Network selection mode: ANYWHERE, ZONE, ON_NETWORK, or OFF_NETWORK. Default = ANYWHERE",
									},
									"include": {
										Type:          schema.TypeList,
										Optional:      true,
										Description:   "The zones to include",
										ConflictsWith: []string{"conditions.0.network.0.exclude"},
										Elem:          &schema.Schema{Type: schema.TypeString},
									},
									"exclude": {
										Type:          schema.TypeList,
										Optional:      true,
										Description:   "The zones to exclude",
										ConflictsWith: []string{"conditions.0.network.0.include"},
										Elem:          &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"authtype": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"ANY", "RADIUS"}, false),
							Description:  "Authentication entrypoint: ANY or RADIUS. Default = ANY",
						},
					},
				},
			},
			"actions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Actions for a rule",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"password": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Password Policy Rule actions",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"passwordchange": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"ALLOW", "DENY"}, false),
										Description:  "Allow or deny a user to change their password: ALLOW or DENY. Default = ALLOW",
									},
									"passwordreset": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"ALLOW", "DENY"}, false),
										Description:  "Allow or deny a user to reset their password: ALLOW or DENY. Default = ALLOW",
									},
									"passwordunlock": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"ALLOW", "DENY"}, false),
										Description:  "Allow or deny a user to unock their password: ALLOW or DENY. Default = DENY",
									},
								},
							},
						},
						"signon": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "SignOn Policy Rule actions",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"access": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"ALLOW", "DENY"}, false),
										Description:  "Allow or deny access based on the rule conditions: ALLOW or DENY. Default = ALLOW",
									},
									"requiremfa": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Require MFA. Default = false",
									},
									"mfaprompt": { // requiremfa must be true
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"DEVICE", "SESSION", "ALWAYS"}, false),
										Description:  "Prompt for MFA based on the device used, a factor session lifetime, or every sign on attempt: DEVICE, SESSION or ALWAYS",
									},
									"remembermfadevice": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Remember MFA device. Default = false",
									},
									"mfalifetime": { // requiremfa must be true, mfaprompt must be SESSION
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Elapsed time before the next MFA challenge",
									},
									"sessionidle": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Max minutes a session can be idle. Default = 120",
									},
									"sessionlifetime": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Max minutes a session is active: Disable = 0. Default = 120",
									},
									"persistentcookie": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether session cookies will last across broswer sessions. Okta Administrators can never have persistent session cookies. Default = false",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourcePolicyRuleCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Creating Policy Rule %v", d.Get("name").(string))
	client := m.(*Config).articulateOktaClient

	var ruleID string
	exists := false
	_, _, err := client.Policies.GetPolicy(d.Get("policyid").(string))
	if err != nil {
		return fmt.Errorf("[ERROR] Error Listing Policy in Okta: %v", err)
	}

	currentPolicyRules, _, err := client.Policies.GetPolicyRules(d.Get("policyid").(string))
	if err != nil {
		return fmt.Errorf("[ERROR] Error Listing Policy Rules in Okta: %v", err)
	}
	if currentPolicyRules != nil {
		for _, rule := range currentPolicyRules.Rules {
			if rule.Name == d.Get("name").(string) {
				ruleID = rule.ID
				systemPolicyRule = rule.System
				exists = true
				break
			}
		}
	}

	if exists == true {
		log.Printf("[INFO] Policy Rule %v already exists in Okta. Adding to Terraform.", d.Get("name").(string))
		d.SetId(ruleID)
	} else {
		switch d.Get("type").(string) {
		case "PASSWORD":
			err = policyRulePassword("create", d, m)
			if err != nil {
				return err
			}

		case "SIGN_ON":
			err = policyRuleSignOn("create", d, m)
			if err != nil {
				return err
			}

		case "MFA_ENROLL":
			return fmt.Errorf("[ERROR] MFA Policy Rule not supported in this terraform provider at this time")

		}
	}

	return nil
}

func resourcePolicyRuleRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] List Policy Rule %v", d.Get("name").(string))

	exists, err := policyRuleExists(d, m)
	if err != nil {
		return err
	}
	if exists == false {
		// if the policy rule does not exist in okta, delete from terraform state
		d.SetId("")
		return nil
	}

	return nil
}

func resourcePolicyRuleUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Update Policy Rule %v", d.Get("name").(string))
	d.Partial(true)

	exists, err := policyRuleExists(d, m)
	if err != nil {
		return err
	}
	if exists == true {
		switch d.Get("type").(string) {
		case "PASSWORD":
			err = policyRulePassword("update", d, m)
			if err != nil {
				return err
			}

		case "SIGN_ON":
			err = policyRuleSignOn("update", d, m)
			if err != nil {
				return err
			}

		case "MFA_ENROLL":
			return fmt.Errorf("[ERROR] MFA Policy Rule not supported in this terraform provider at this time")

		}
	} else {
		return fmt.Errorf("[ERROR] Error Policy not found in Okta: %v", err)
	}
	d.Partial(false)

	return nil
}

func resourcePolicyRuleDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Delete Policy Rule %v", d.Get("name").(string))
	client := m.(*Config).articulateOktaClient

	exists, err := policyRuleExists(d, m)
	if err != nil {
		return err
	}
	if exists == true {
		if systemPolicyRule == true {
			log.Printf("[INFO] Policy Rule %v is a System Policy, cannot delete from Okta", d.Get("name").(string))
		} else {
			_, err = client.Policies.DeletePolicyRule(d.Get("policyid").(string), d.Id())
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

// check if policy rule exists in Okta
func policyRuleExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*Config).articulateOktaClient

	_, _, err := client.Policies.GetPolicy(d.Get("policyid").(string))
	if client.OktaErrorCode == "E0000007" {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("[ERROR] Error Listing Policy in Okta: %v", err)
	}

	rule, _, err := client.Policies.GetPolicyRule(d.Get("policyid").(string), d.Id())
	if client.OktaErrorCode == "E0000007" {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("[ERROR] Error Listing Policy Rule in Okta: %v", err)
	}
	systemPolicyRule = rule.System
	return true, nil

}

// populate policy rule conditions with the terraform schema conditions fields
func policyRuleConditions(d *schema.ResourceData) ([]string, error) {
	users := make([]string, 0)
	if len(d.Get("conditions.0.users").([]interface{})) > 0 {
		for _, vals := range d.Get("conditions.0.users").([]interface{}) {
			users = append(users, vals.(string))
		}
	}
	if len(d.Get("conditions.0.network").([]interface{})) > 0 {
		return users, fmt.Errorf("[ERROR] network condition not supported in this terraform provider at this time")
	}
	if len(d.Get("conditions.0.authtype").(string)) > 0 {
		return users, fmt.Errorf("[ERROR] authtype condition not supported in this terraform provider at this time")
	}
	return users, nil
}

// activate or deactivate a policy rule according to the terraform schema status field
func policyRuleActivate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Config).articulateOktaClient

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

func getDefaultNetwork() *articulateOkta.Network {
	return &articulateOkta.Network{
		Connection: "ANYWHERE", // Okta required default
	}
}

func getExcludedUsers(excludedUsers []string) *articulateOkta.People {
	return &articulateOkta.People{
		Users: &articulateOkta.Users{
			Exclude: excludedUsers,
		},
	}
}

func getIncludedGroups(includedGroups []string) *articulateOkta.People {
	return &articulateOkta.People{
		Groups: &articulateOkta.Groups{
			Include: includedGroups,
		},
	}
}

func getDefaultAuthContext() *articulateOkta.AuthContext {
	return &articulateOkta.AuthContext{
		AuthType: "ANY", // Okta required default
	}
}

// create or update a password policy rule
func policyRulePassword(action string, d *schema.ResourceData, m interface{}) error {
	client := m.(*Config).articulateOktaClient

	template := client.Policies.PasswordRule()
	template.Name = d.Get("name").(string)
	template.Type = d.Get("type").(string)
	template.Status = d.Get("status").(string)
	if priority, ok := d.GetOk("priority"); ok {
		template.Priority = priority.(int)
	}

	users, err := policyRuleConditions(d)
	if err != nil {
		return err
	}
	template.Conditions = &articulateOkta.PolicyConditions{
		Network: getDefaultNetwork(),
		People:  getExcludedUsers(users),
	}

	// Okta defaults
	// we add the defaults here & not in the schema map to avoid defaults appearing in the terraform plan diff
	template.Actions.PasswordChange.Access = "ALLOW"
	template.Actions.SelfServicePasswordReset.Access = "ALLOW"
	template.Actions.SelfServiceUnlock.Access = "DENY"

	if len(d.Get("actions.0.password").([]interface{})) > 0 {
		if passwordchange, ok := d.GetOk("actions.0.password.0.passwordchange"); ok {
			template.Actions.PasswordChange.Access = passwordchange.(string)
		}
		if passwordreset, ok := d.GetOk("actions.0.password.0.passwordreset"); ok {
			template.Actions.SelfServicePasswordReset.Access = passwordreset.(string)
		}
		if passwordunlock, ok := d.GetOk("actions.0.password.0.passwordunlock"); ok {
			template.Actions.SelfServiceUnlock.Access = passwordunlock.(string)
		}
	}

	switch action {
	case "create":
		rule, _, err := client.Policies.CreatePolicyRule(d.Get("policyid").(string), template)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Creating Policy Rule: %v", err)
		}
		log.Printf("[INFO] Okta Policy Rule Created: %+v", rule)
		log.Printf("[INFO] Adding Policy Rule to Terraform")
		d.SetId(rule.ID)

		err = policyRuleActivate(d, m)
		if err != nil {
			return err
		}

	case "update":
		rule, _, err := client.Policies.UpdatePolicyRule(d.Get("policyid").(string), d.Id(), template)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Updating Policy Rule: %v", err)
		}
		log.Printf("[INFO] Okta Policy Rule Updated: %+v", rule)

		err = policyRuleActivate(d, m)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("[ERROR] policyPasswordRule action only supports \"create\" and \"update\"")
	}
	return nil
}

// create or update a signon policy rule
func policyRuleSignOn(action string, d *schema.ResourceData, m interface{}) error {
	client := m.(*Config).articulateOktaClient

	template := client.Policies.SignOnRule()
	template.Name = d.Get("name").(string)
	template.Type = d.Get("type").(string)
	template.Status = d.Get("status").(string)
	if priority, ok := d.GetOk("priority"); ok {
		template.Priority = priority.(int)
	}

	users, err := policyRuleConditions(d)
	if err != nil {
		return err
	}

	template.Conditions = &articulateOkta.PolicyConditions{
		Network:     getDefaultNetwork(),
		AuthContext: getDefaultAuthContext(),
		People:      getExcludedUsers(users),
	}

	// Okta defaults
	// we add the defaults here & not in the schema map to avoid defaults appearing in the terraform plan diff
	template.Actions.SignOn.Access = "ALLOW"
	template.Actions.SignOn.RequireFactor = false
	//template.Actions.SignOn.FactorPromptMode = "SESSION"
	//template.Actions.SignOn.RememberDeviceByDefault = false
	//template.Actions.SignOn.FactorLifetime = ?
	template.Actions.SignOn.Session.MaxSessionIdleMinutes = 120
	template.Actions.SignOn.Session.MaxSessionLifetimeMinutes = 120
	template.Actions.SignOn.Session.UsePersistentCookie = false

	if len(d.Get("actions.0.signon").([]interface{})) > 0 {
		if access, ok := d.GetOk("actions.0.signon.0.access"); ok {
			template.Actions.SignOn.Access = access.(string)
		}
		//if requiremfa, ok := d.GetOk("actions.0.signon.0.requiremfa"); ok {
		if _, ok := d.GetOk("actions.0.signon.0.requiremfa"); ok {
			return fmt.Errorf("[ERROR] mfa signon actions not supported in this terraform provider at this time")
		}
		//if mfaprompt, ok := d.GetOk("actions.0.signon.0.mfaprompt"); ok {
		if _, ok := d.GetOk("actions.0.signon.0.mfaprompt"); ok {
			return fmt.Errorf("[ERROR] mfa signon actions not supported in this terraform provider at this time")
		}
		//if remembermfadevice, ok := d.GetOk("actions.0.signon.0.remembermfadevice"); ok {
		if _, ok := d.GetOk("actions.0.signon.0.remembermfadevice"); ok {
			return fmt.Errorf("[ERROR] mfa signon actions not supported in this terraform provider at this time")
		}
		//if mfalifetime, ok := d.GetOk("actions.0.signon.0.mfalifetime"); ok {
		if _, ok := d.GetOk("actions.0.signon.0.mfalifetime"); ok {
			return fmt.Errorf("[ERROR] mfa signon actions not supported in this terraform provider at this time")
		}
		if sessionidle, ok := d.GetOk("actions.0.signon.0.sessionidle"); ok {
			template.Actions.SignOn.Session.MaxSessionIdleMinutes = sessionidle.(int)
		}
		if sessionlifetime, ok := d.GetOk("actions.0.signon.0.sessionlifetime"); ok {
			template.Actions.SignOn.Session.MaxSessionLifetimeMinutes = sessionlifetime.(int)
		}
		if persistentcookie, ok := d.GetOk("actions.0.signon.0.persistentcookie"); ok {
			template.Actions.SignOn.Session.UsePersistentCookie = persistentcookie.(bool)
		}
	}

	switch action {
	case "create":
		rule, _, err := client.Policies.CreatePolicyRule(d.Get("policyid").(string), template)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Creating Policy Rule: %v", err)
		}
		log.Printf("[INFO] Okta Policy Rule Created: %+v", rule)
		log.Printf("[INFO] Adding Policy Rule to Terraform")
		d.SetId(rule.ID)

		err = policyRuleActivate(d, m)
		if err != nil {
			return err
		}

	case "update":
		rule, _, err := client.Policies.UpdatePolicyRule(d.Get("policyid").(string), d.Id(), template)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Updating Policy Rule: %v", err)
		}
		log.Printf("[INFO] Okta Policy Updated: %+v", rule)

		err = policyRuleActivate(d, m)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("[ERROR] policySignOn rule action only supports \"create\" and \"update\"")
	}
	return nil
}
