package okta

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/terraform-provider-okta/sdk"
)

func resourceAuthServerPolicyRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAuthServerPolicyRuleCreate,
		ReadContext:   resourceAuthServerPolicyRuleRead,
		UpdateContext: resourceAuthServerPolicyRuleUpdate,
		DeleteContext: resourceAuthServerPolicyRuleDelete,
		Importer:      createNestedResourceImporter([]string{"auth_server_id", "policy_id", "id"}),
		Schema: map[string]*schema.Schema{
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
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: stringInSlice([]string{authorizationCode, implicit, password, clientCredentials}),
				},
				Description: "Accepted grant type values: authorization_code, implicit, password, client_credentials",
			},
			"scope_whitelist": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"access_token_lifetime_minutes": {
				Type:     schema.TypeInt,
				Optional: true,
				// 5 minutes - 1 day
				ValidateDiagFunc: intBetween(5, 1440),
				Default:          60,
			},
			"refresh_token_lifetime_minutes": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"refresh_token_window_minutes": {
				Type:     schema.TypeInt,
				Optional: true,
				// 5 minutes - 5 years
				ValidateDiagFunc: intBetween(5, 2628000),
				Default:          10080,
			},
			"inline_hook_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_whitelist": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"user_blacklist": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"group_whitelist": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"group_blacklist": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
		},
	}
}

func resourceAuthServerPolicyRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := validateAuthServerPolicyRule(d)
	if err != nil {
		return diag.FromErr(err)
	}
	authServerPolicyRule := buildAuthServerPolicyRule(d)
	responseAuthServerPolicyRule, _, err := getSupplementFromMetadata(m).CreateAuthorizationServerPolicyRule(
		ctx,
		d.Get("auth_server_id").(string),
		d.Get("policy_id").(string),
		*authServerPolicyRule,
		nil,
	)
	if err != nil {
		return diag.Errorf("failed to create auth server policy rule: %v", err)
	}
	d.SetId(responseAuthServerPolicyRule.Id)
	return resourceAuthServerPolicyRuleRead(ctx, d, m)
}

func resourceAuthServerPolicyRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	authServerPolicyRule, resp, err := getSupplementFromMetadata(m).GetAuthorizationServerPolicyRule(
		ctx, d.Get("auth_server_id").(string),
		d.Get("policy_id").(string),
		d.Id(),
		sdk.AuthorizationServerPolicyRule{},
	)
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get auth server policy rule: %v", err)
	}
	if authServerPolicyRule == nil {
		d.SetId("")
		return nil
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
		return diag.Errorf("failed to read auth server rule: %v", err)
	}
	err = setPeopleAssignments(d, authServerPolicyRule.Conditions.People)
	if err != nil {
		return diag.Errorf("failed to read auth server rule: %v", err)
	}
	return nil
}

func resourceAuthServerPolicyRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := validateAuthServerPolicyRule(d)
	if err != nil {
		return diag.FromErr(err)
	}
	authServerPolicyRule := buildAuthServerPolicyRule(d)
	_, _, err = getSupplementFromMetadata(m).UpdateAuthorizationServerPolicyRule(
		ctx,
		d.Get("auth_server_id").(string),
		d.Get("policy_id").(string), d.Id(),
		*authServerPolicyRule,
		nil,
	)
	if err != nil {
		return diag.Errorf("failed to update auth server policy rule: %v", err)
	}
	return resourceAuthServerPolicyRuleRead(ctx, d, m)
}

func resourceAuthServerPolicyRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	_, err := getSupplementFromMetadata(m).DeleteAuthorizationServerPolicyRule(
		ctx,
		d.Get("auth_server_id").(string),
		d.Get("policy_id").(string),
		d.Id(),
	)
	if err != nil {
		return diag.Errorf("failed to delete auth server policy rule: %v", err)
	}
	return nil
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
			GrantTypes: &sdk.Conditions{Include: convertInterfaceToStringSet(d.Get("grant_type_whitelist"))},
			Scopes:     &sdk.Conditions{Include: convertInterfaceToStringSet(d.Get("scope_whitelist"))},
			People:     getPeopleConditions(d),
		},
	}
}

func setPeopleAssignments(d *schema.ResourceData, c *okta.GroupRulePeopleCondition) error {
	if c.Groups != nil {
		err := setNonPrimitives(d, map[string]interface{}{
			"group_whitelist": convertStringSetToInterface(c.Groups.Include),
			"group_blacklist": convertStringSetToInterface(c.Groups.Exclude),
		})
		if err != nil {
			return err
		}
	} else {
		_ = setNonPrimitives(d, map[string]interface{}{
			"group_whitelist": convertStringSetToInterface([]string{}),
			"group_blacklist": convertStringSetToInterface([]string{}),
		})
	}
	return setNonPrimitives(d, map[string]interface{}{
		"user_whitelist": convertStringSetToInterface(c.Users.Include),
		"user_blacklist": convertStringSetToInterface(c.Users.Exclude),
	})
}

func validateAuthServerPolicyRule(d *schema.ResourceData) error {
	if w, ok := d.GetOk("grant_type_whitelist"); ok {
		for _, v := range convertInterfaceToStringSet(w) {
			if v != implicit {
				continue
			}
			_, okUsers := d.GetOk("user_whitelist")
			_, okGroups := d.GetOk("group_whitelist")
			if !okUsers && !okGroups {
				return fmt.Errorf(`at least "user_whitelist" or "group_whitelist" should be provided when using '%s' in "grant_type_whitelist"`, implicit)
			}
		}
	}
	rtlm := d.Get("refresh_token_lifetime_minutes").(int)
	atlm := d.Get("access_token_lifetime_minutes").(int)
	rtwm := d.Get("refresh_token_window_minutes").(int)
	if rtlm > 0 && rtlm < atlm {
		return errors.New("'refresh_token_lifetime_minutes' must be greater than or equal to 'access_token_lifetime_minutes'")
	}
	if rtlm > 0 && (atlm > rtwm || rtlm < rtwm) {
		return errors.New("'refresh_token_window_minutes' must be between 'access_token_lifetime_minutes' and 'refresh_token_lifetime_minutes'")
	}
	return nil
}
