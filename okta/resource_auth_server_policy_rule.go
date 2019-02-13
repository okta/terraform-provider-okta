package okta

import (
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAuthServerPolicyRule() *schema.Resource {
	return &schema.Resource{
		Create:   resourceAuthServerPolicyRuleCreate,
		Exists:   resourceAuthServerPolicyRuleExists,
		Read:     resourceAuthServerPolicyRuleRead,
		Update:   resourceAuthServerPolicyRuleUpdate,
		Delete:   resourceAuthServerPolicyRuleDelete,
		Importer: createNestedResourceImporter([]string{"auth_server_id", "policy_id", "id"}),

		Schema: map[string]*schema.Schema{
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "RESOURCE_ACCESS",
				Description: "Auth server policy rule type, unlikely this will be anything other then the default",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Auth server policy rule name",
			},
			"auth_server_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Auth server ID",
			},
			"policy_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Auth server policy ID",
			},
			"status": statusSchema,
			"priority": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Priority of the auth server policy rule",
			},
			"grant_type_whitelist": &schema.Schema{
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Accepted grant type values: authorization_code, implicit, password.",
			},
			"assignments": peopleSchema,
			"scope_whitelist": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func buildAuthServerPolicyRule(d *schema.ResourceData) *AuthorizationServerPolicyRule {
	return &AuthorizationServerPolicyRule{
		Name:     d.Get("name").(string),
		Status:   d.Get("status").(string),
		Priority: d.Get("priority").(int),
		Type:     d.Get("type").(string),
	}
}

func resourceAuthServerPolicyRuleCreate(d *schema.ResourceData, m interface{}) error {
	authServerPolicyRule := buildAuthServerPolicyRule(d)
	c := getSupplementFromMetadata(m)
	authServerId := d.Get("auth_server_id").(string)
	policyId := d.Get("policy_id").(string)
	responseAuthServerPolicyRule, _, err := c.CreateAuthorizationServerPolicyRule(authServerId, policyId, *authServerPolicyRule, nil)
	if err != nil {
		return err
	}

	d.SetId(responseAuthServerPolicyRule.Id)

	return resourceAuthServerPolicyRuleRead(d, m)
}

func resourceAuthServerPolicyRuleExists(d *schema.ResourceData, m interface{}) (bool, error) {
	g, err := fetchAuthServerPolicyRule(d, m)

	return err == nil && g != nil, err
}

func resourceAuthServerPolicyRuleRead(d *schema.ResourceData, m interface{}) error {
	authServerPolicyRule, err := fetchAuthServerPolicyRule(d, m)
	if err != nil {
		return err
	}

	d.Set("name", authServerPolicyRule.Name)
	d.Set("status", authServerPolicyRule.Status)
	d.Set("priority", authServerPolicyRule.Priority)
	d.Set("type", authServerPolicyRule.Type)
	d.Set("people", flattenPeopleConditions(authServerPolicyRule.Conditions.People))

	return nil
}

func resourceAuthServerPolicyRuleUpdate(d *schema.ResourceData, m interface{}) error {
	authServerPolicyRule := buildAuthServerPolicyRule(d)
	c := getSupplementFromMetadata(m)
	authServerId := d.Get("auth_server_id").(string)
	policyId := d.Get("policy_id").(string)
	_, _, err := c.UpdateAuthorizationServerPolicyRule(authServerId, policyId, d.Id(), *authServerPolicyRule, nil)
	if err != nil {
		return err
	}

	return resourceAuthServerPolicyRuleRead(d, m)
}

func resourceAuthServerPolicyRuleDelete(d *schema.ResourceData, m interface{}) error {
	_, err := getSupplementFromMetadata(m).DeleteAuthorizationServerPolicyRule(d.Get("auth_server_id").(string), d.Get("policy_id").(string), d.Id())

	return err
}

func fetchAuthServerPolicyRule(d *schema.ResourceData, m interface{}) (*AuthorizationServerPolicyRule, error) {
	c := getSupplementFromMetadata(m)
	authServerId := d.Get("auth_server_id").(string)
	policyId := d.Get("policy_id").(string)
	auth, resp, err := c.GetAuthorizationServerPolicyRule(authServerId, policyId, d.Id(), AuthorizationServerPolicyRule{})

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	return auth, err
}
