package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/terraform-provider-okta/sdk"
)

func resourceAuthServerPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAuthServerPolicyCreate,
		ReadContext:   resourceAuthServerPolicyRead,
		UpdateContext: resourceAuthServerPolicyUpdate,
		DeleteContext: resourceAuthServerPolicyDelete,
		Importer:      createNestedResourceImporter([]string{"auth_server_id", "id"}),
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     sdk.OauthAuthorizationPolicyType,
				Description: "Auth server policy type, unlikely this will be anything other then the default",
				Deprecated:  "Policy type can only be of value 'OAUTH_AUTHORIZATION_POLICY', so this will be removed in the future, or set as 'Computed' value",
				DiffSuppressFunc: func(string, string, string, *schema.ResourceData) bool {
					return true
				},
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"auth_server_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": statusSchema,
			"priority": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Priority of the auth server policy",
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"client_whitelist": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Use [\"ALL_CLIENTS\"] when unsure.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAuthServerPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if d.Get("status").(string) == statusInactive {
		return diag.Errorf("can not create an inactive auth server policy, only existing ones can be deactivated")
	}
	policy := buildAuthServerPolicy(d)
	respPolicy, _, err := getOktaClientFromMetadata(m).AuthorizationServer.CreateAuthorizationServerPolicy(ctx, d.Get("auth_server_id").(string), policy)
	if err != nil {
		return diag.Errorf("failed to create authorization server policy: %v", err)
	}
	d.SetId(respPolicy.Id)
	return resourceAuthServerPolicyRead(ctx, d, m)
}

func resourceAuthServerPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	policy, resp, err := getOktaClientFromMetadata(m).AuthorizationServer.GetAuthorizationServerPolicy(ctx, d.Get("auth_server_id").(string), d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get auth server policy: %v", err)
	}
	if policy == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("name", policy.Name)
	_ = d.Set("description", policy.Description)
	_ = d.Set("status", policy.Status)
	_ = d.Set("priority", policy.Priority)
	_ = d.Set("client_whitelist", convertStringSetToInterface(policy.Conditions.Clients.Include))
	return nil
}

func resourceAuthServerPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	policy := buildAuthServerPolicy(d)
	_, _, err := getOktaClientFromMetadata(m).AuthorizationServer.UpdateAuthorizationServerPolicy(ctx, d.Get("auth_server_id").(string), d.Id(), policy)
	if err != nil {
		return diag.Errorf("failed to update auth server policy: %v", err)
	}
	oldStatus, newStatus := d.GetChange("status")
	if oldStatus != newStatus {
		if newStatus == statusActive {
			_, err = getSupplementFromMetadata(m).ActivateAuthorizationServerPolicy(ctx, d.Get("auth_server_id").(string), d.Id())
		} else {
			_, err = getSupplementFromMetadata(m).DeactivateAuthorizationServerPolicy(ctx, d.Get("auth_server_id").(string), d.Id())
		}
		if err != nil {
			return diag.Errorf("failed to change authorization server policy status: %v", err)
		}
	}
	return resourceAuthServerPolicyRead(ctx, d, m)
}

func resourceAuthServerPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	_, err := getOktaClientFromMetadata(m).AuthorizationServer.DeleteAuthorizationServerPolicy(ctx, d.Get("auth_server_id").(string), d.Id())
	if err != nil {
		return diag.Errorf("failed to delete auth server policy: %v", err)
	}
	return nil
}

func buildAuthServerPolicy(d *schema.ResourceData) okta.Policy {
	return okta.Policy{
		Name:        d.Get("name").(string),
		Type:        sdk.OauthAuthorizationPolicyType,
		Status:      d.Get("status").(string),
		Priority:    int64(d.Get("priority").(int)),
		Description: d.Get("description").(string),
		Conditions: &okta.PolicyRuleConditions{
			Clients: &okta.ClientPolicyCondition{
				Include: convertInterfaceToStringSet(d.Get("client_whitelist")),
			},
		},
	}
}
