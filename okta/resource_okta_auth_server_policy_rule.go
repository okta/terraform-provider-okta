package okta

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

func resourceAuthServerPolicyRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAuthServerPolicyRuleCreate,
		ReadContext:   resourceAuthServerPolicyRuleRead,
		UpdateContext: resourceAuthServerPolicyRuleUpdate,
		DeleteContext: resourceAuthServerPolicyRuleDelete,
		Importer:      createNestedResourceImporter([]string{"auth_server_id", "policy_id", "id"}),
		Description: `Creates an Authorization Server Policy Rule.
This resource allows you to create and configure an Authorization Server Policy Rule.
-> This resource is concurrency safe. However, when creating/updating/deleting
multiple rules belonging to a policy, the Terraform meta argument
['depends_on'](https://www.terraform.io/language/meta-arguments/depends_on)
should be added to each rule chaining them all in sequence. Base the sequence on
the 'priority' property in ascending value.`,
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
			"system": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The rule is the system (default) rule for its associated policy",
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
					Type: schema.TypeString,
				},
				Description: "Accepted grant type values, `authorization_code`, `implicit`, `password`, `client_credentials`, `urn:ietf:params:oauth:grant-type:saml2-bearer` (*Early Access Property*), `urn:ietf:params:oauth:grant-type:token-exchange` (*Early Access Property*),`urn:ietf:params:oauth:grant-type:device_code` (*Early Access Property*), `interaction_code` (*OIE only*). For `implicit` value either `user_whitelist` or `group_whitelist` should be set.",
			},
			"scope_whitelist": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Scopes allowed for this policy rule. They can be whitelisted by name or all can be whitelisted with ` * `",
			},
			"access_token_lifetime_minutes": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     60,
				Description: "Lifetime of access token. Can be set to a value between 5 and 1440 minutes. Default is `60`.",
			},
			"refresh_token_lifetime_minutes": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Lifetime of refresh token.",
			},
			"refresh_token_window_minutes": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10080,
				Description: "Window in which a refresh token can be used. It can be a value between 5 and 2628000 (5 years) minutes. Default is `10080` (7 days).`refresh_token_window_minutes` must be between `access_token_lifetime_minutes` and `refresh_token_lifetime_minutes`.",
			},
			"inline_hook_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the inline token to trigger.",
			},
			"user_whitelist": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Specifies a set of Users to be included.",
			},
			"user_blacklist": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Specifies a set of Users to be excluded.",
			},
			"group_whitelist": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Specifies a set of Groups whose Users are to be included. Can be set to Group ID or to the following: `EVERYONE`.",
			},
			"group_blacklist": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Specifies a set of Groups whose Users are to be excluded.",
			},
		},
	}
}

func resourceAuthServerPolicyRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	oktaMutexKV.Lock(authServerPolicyRule)
	defer oktaMutexKV.Unlock(authServerPolicyRule)

	err := validateAuthServerPolicyRule(d)
	if err != nil {
		return diag.FromErr(err)
	}
	resp, _, err := getOktaClientFromMetadata(m).AuthorizationServer.CreateAuthorizationServerPolicyRule(
		ctx, d.Get("auth_server_id").(string), d.Get("policy_id").(string), buildAuthServerPolicyRule(d))
	if err != nil {
		return diag.Errorf("failed to create auth server policy rule: %v", err)
	}
	d.SetId(resp.Id)
	return resourceAuthServerPolicyRuleRead(ctx, d, m)
}

func resourceAuthServerPolicyRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	authServerPolicyRule, resp, err := getOktaClientFromMetadata(m).AuthorizationServer.GetAuthorizationServerPolicyRule(
		ctx, d.Get("auth_server_id").(string), d.Get("policy_id").(string), d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get auth server policy rule: %v", err)
	}
	if authServerPolicyRule == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("system", boolFromBoolPtr(authServerPolicyRule.System))
	_ = d.Set("name", authServerPolicyRule.Name)
	_ = d.Set("status", authServerPolicyRule.Status)
	if authServerPolicyRule.PriorityPtr != nil {
		_ = d.Set("priority", *authServerPolicyRule.PriorityPtr)
	}
	_ = d.Set("type", authServerPolicyRule.Type)
	if authServerPolicyRule.Actions.Token.AccessTokenLifetimeMinutesPtr != nil {
		_ = d.Set("access_token_lifetime_minutes", *authServerPolicyRule.Actions.Token.AccessTokenLifetimeMinutesPtr)
	}
	if authServerPolicyRule.Actions.Token.RefreshTokenLifetimeMinutesPtr != nil {
		_ = d.Set("refresh_token_lifetime_minutes", *authServerPolicyRule.Actions.Token.RefreshTokenLifetimeMinutesPtr)
	}
	if authServerPolicyRule.Actions.Token.RefreshTokenWindowMinutesPtr != nil {
		_ = d.Set("refresh_token_window_minutes", *authServerPolicyRule.Actions.Token.RefreshTokenWindowMinutesPtr)
	}

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
	oktaMutexKV.Lock(authServerPolicyRule)
	defer oktaMutexKV.Unlock(authServerPolicyRule)

	err := validateAuthServerPolicyRule(d)
	if err != nil {
		return diag.FromErr(err)
	}
	_, _, err = getOktaClientFromMetadata(m).AuthorizationServer.UpdateAuthorizationServerPolicyRule(
		ctx,
		d.Get("auth_server_id").(string),
		d.Get("policy_id").(string), d.Id(),
		buildAuthServerPolicyRule(d),
	)
	if err != nil {
		return diag.Errorf("failed to update auth server policy rule: %v", err)
	}
	if d.HasChange("status") {
		err := handleAuthServerPolicyRuleLifecycle(ctx, d, m)
		if err != nil {
			return err
		}
	}
	return resourceAuthServerPolicyRuleRead(ctx, d, m)
}

func handleAuthServerPolicyRuleLifecycle(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var err error
	oldStatus, newStatus := d.GetChange("status")
	if oldStatus == newStatus {
		return nil
	}
	if newStatus == statusActive {
		_, err = getOktaClientFromMetadata(m).AuthorizationServer.ActivateAuthorizationServerPolicyRule(ctx, d.Get("auth_server_id").(string),
			d.Get("policy_id").(string), d.Id())
	} else {
		_, err = getOktaClientFromMetadata(m).AuthorizationServer.DeactivateAuthorizationServerPolicyRule(ctx, d.Get("auth_server_id").(string),
			d.Get("policy_id").(string), d.Id())
	}
	if err != nil {
		return diag.Errorf("failed to change authorization server policy status: %v", err)
	}
	return nil
}

func resourceAuthServerPolicyRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	oktaMutexKV.Lock(authServerPolicyRule)
	defer oktaMutexKV.Unlock(authServerPolicyRule)

	_, err := getOktaClientFromMetadata(m).AuthorizationServer.DeleteAuthorizationServerPolicyRule(
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

func buildAuthServerPolicyRule(d *schema.ResourceData) sdk.AuthorizationServerPolicyRule {
	var hook *sdk.TokenAuthorizationServerPolicyRuleActionInlineHook
	inlineHook := d.Get("inline_hook_id").(string)
	if inlineHook != "" {
		hook = &sdk.TokenAuthorizationServerPolicyRuleActionInlineHook{
			Id: inlineHook,
		}
	}
	return sdk.AuthorizationServerPolicyRule{
		Name:        d.Get("name").(string),
		Status:      d.Get("status").(string),
		PriorityPtr: int64Ptr(d.Get("priority").(int)),
		Type:        d.Get("type").(string),
		Actions: &sdk.AuthorizationServerPolicyRuleActions{
			Token: &sdk.TokenAuthorizationServerPolicyRuleAction{
				AccessTokenLifetimeMinutesPtr:  int64Ptr(d.Get("access_token_lifetime_minutes").(int)),
				RefreshTokenLifetimeMinutesPtr: int64Ptr(d.Get("refresh_token_lifetime_minutes").(int)),
				RefreshTokenWindowMinutesPtr:   int64Ptr(d.Get("refresh_token_window_minutes").(int)),
				InlineHook:                     hook,
			},
		},
		Conditions: &sdk.AuthorizationServerPolicyRuleConditions{
			GrantTypes: &sdk.GrantTypePolicyRuleCondition{Include: convertInterfaceToStringSet(d.Get("grant_type_whitelist"))},
			Scopes:     &sdk.OAuth2ScopesMediationPolicyRuleCondition{Include: convertInterfaceToStringSet(d.Get("scope_whitelist"))},
			People: &sdk.PolicyPeopleCondition{
				Groups: &sdk.GroupCondition{
					Include: convertInterfaceToStringSet(d.Get("group_whitelist")),
					Exclude: convertInterfaceToStringSet(d.Get("group_blacklist")),
				},
				Users: &sdk.UserCondition{
					Include: convertInterfaceToStringSet(d.Get("user_whitelist")),
					Exclude: convertInterfaceToStringSet(d.Get("user_blacklist")),
				},
			},
		},
	}
}

func setPeopleAssignments(d *schema.ResourceData, c *sdk.PolicyPeopleCondition) error {
	if c.Groups != nil {
		err := setNonPrimitives(d, map[string]interface{}{
			"group_whitelist": convertStringSliceToSet(c.Groups.Include),
			"group_blacklist": convertStringSliceToSet(c.Groups.Exclude),
		})
		if err != nil {
			return err
		}
	} else {
		_ = setNonPrimitives(d, map[string]interface{}{
			"group_whitelist": convertStringSliceToSet([]string{}),
			"group_blacklist": convertStringSliceToSet([]string{}),
		})
	}
	return setNonPrimitives(d, map[string]interface{}{
		"user_whitelist": convertStringSliceToSet(c.Users.Include),
		"user_blacklist": convertStringSliceToSet(c.Users.Exclude),
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
