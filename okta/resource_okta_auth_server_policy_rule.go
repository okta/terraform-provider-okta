package okta

import (
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/oktadeveloper/terraform-provider-okta/sdk"
)

func resourceAuthServerPolicyRule() *schema.Resource {
	return &schema.Resource{
		Create:   resourceAuthServerPolicyRuleCreate,
		Exists:   resourceAuthServerPolicyRuleExists,
		Read:     resourceAuthServerPolicyRuleRead,
		Update:   resourceAuthServerPolicyRuleUpdate,
		Delete:   resourceAuthServerPolicyRuleDelete,
		Importer: createNestedResourceImporter([]string{"auth_server_id", "policy_id", "id"}),

		Schema: addPeopleAssignments(map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "RESOURCE_ACCESS",
				Description: "Auth server policy rule type, unlikely this will be anything other then the default",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Auth server policy rule name",
			},
			"auth_server_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Auth server ID",
			},
			"policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Auth server policy ID",
			},
			"status": statusSchema,
			"priority": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Priority of the auth server policy rule",
			},
			"grant_type_whitelist": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Accepted grant type values: authorization_code, implicit, password.",
			},
			"scope_whitelist": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"access_token_lifetime_minutes": {
				Type:     schema.TypeInt,
				Optional: true,
				// 5 mins - 1 day
				ValidateFunc: validation.IntBetween(5, 1440),
				Default:      60,
			},
			"refresh_token_lifetime_minutes": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"refresh_token_window_minutes": {
				Type:     schema.TypeInt,
				Optional: true,
				// 10 mins - 5 years
				ValidateFunc: validation.IntBetween(10, 2628000),
				Default:      10080,
			},
			"inline_hook_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		}),
	}
}

func buildAuthServerPolicyRule(d *schema.ResourceData) *sdk.AuthorizationServerPolicyRule {
	var hook *sdk.AuthServerInlineHook

	inlineHook := d.Get("inline_hook_id").(string)

	if inlineHook != "" {
		hook = &sdk.AuthServerInlineHook{
			Id: inlineHook,
		}
	}

	return &sdk.AuthorizationServerPolicyRule{
		Name:     d.Get("name").(string),
		Status:   d.Get("status").(string),
		Priority: d.Get("priority").(int),
		Type:     d.Get("type").(string),
		Actions: &sdk.AuthorizationServerPolicyRuleActions{
			Token: &sdk.TokenActions{
				AccessTokenLifetimeMinutes:  d.Get("access_token_lifetime_minutes").(int),
				RefreshTokenLifetimeMinutes: d.Get("refresh_token_lifetime_minutes").(int),
				RefreshTokenWindowMinutes:   d.Get("refresh_token_window_minutes").(int),
				InlineHook:                  hook,
			},
		},
		Conditions: &sdk.AuthorizationServerPolicyRuleConditions{
			GrantTypes: &sdk.Whitelist{Include: convertInterfaceToStringSet(d.Get("grant_type_whitelist"))},
			Scopes:     &sdk.Whitelist{Include: convertInterfaceToStringSet(d.Get("scope_whitelist"))},
			People:     getPeopleConditions(d),
		},
	}
}

func resourceAuthServerPolicyRuleCreate(d *schema.ResourceData, m interface{}) error {
	authServerPolicyRule := buildAuthServerPolicyRule(d)
	c := getSupplementFromMetadata(m)
	authServerID := d.Get("auth_server_id").(string)
	policyID := d.Get("policy_id").(string)
	responseAuthServerPolicyRule, _, err := c.CreateAuthorizationServerPolicyRule(authServerID, policyID, *authServerPolicyRule, nil)
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

	if authServerPolicyRule == nil {
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	_ = d.Set("name", authServerPolicyRule.Name)
	_ = d.Set("status", authServerPolicyRule.Status)
	_ = d.Set("priority", authServerPolicyRule.Priority)
	_ = d.Set("type", authServerPolicyRule.Type)

	if authServerPolicyRule.Actions.Token.InlineHook != nil {
		_ = d.Set("inline_hook_id", authServerPolicyRule.Actions.Token.InlineHook.Id)
	}

	err = setNonPrimitives(d, map[string]interface{}{
		"grant_type_whitelist": authServerPolicyRule.Conditions.GrantTypes.Include,
		"scope_whitelist":      authServerPolicyRule.Conditions.Scopes.Include,
	})
	if err != nil {
		return err
	}

	return setPeopleAssignments(d, authServerPolicyRule.Conditions.People)
}

func resourceAuthServerPolicyRuleUpdate(d *schema.ResourceData, m interface{}) error {
	authServerPolicyRule := buildAuthServerPolicyRule(d)
	c := getSupplementFromMetadata(m)
	authServerID := d.Get("auth_server_id").(string)
	policyID := d.Get("policy_id").(string)
	_, _, err := c.UpdateAuthorizationServerPolicyRule(authServerID, policyID, d.Id(), *authServerPolicyRule, nil)
	if err != nil {
		return err
	}

	return resourceAuthServerPolicyRuleRead(d, m)
}

func resourceAuthServerPolicyRuleDelete(d *schema.ResourceData, m interface{}) error {
	authServerID := d.Get("auth_server_id").(string)
	policyID := d.Get("policy_id").(string)
	_, err := getSupplementFromMetadata(m).DeleteAuthorizationServerPolicyRule(authServerID, policyID, d.Id())

	return err
}

func fetchAuthServerPolicyRule(d *schema.ResourceData, m interface{}) (*sdk.AuthorizationServerPolicyRule, error) {
	c := getSupplementFromMetadata(m)
	authServerID := d.Get("auth_server_id").(string)
	policyID := d.Get("policy_id").(string)
	auth, resp, err := c.GetAuthorizationServerPolicyRule(authServerID, policyID, d.Id(), sdk.AuthorizationServerPolicyRule{})

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	return auth, err
}
