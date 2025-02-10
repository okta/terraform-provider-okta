package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

func resourceAuthServerPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAuthServerPolicyCreate,
		ReadContext:   resourceAuthServerPolicyRead,
		UpdateContext: resourceAuthServerPolicyUpdate,
		DeleteContext: resourceAuthServerPolicyDelete,
		Importer:      createNestedResourceImporter([]string{"auth_server_id", "id"}),
		Description:   "Creates an Authorization Server Policy. This resource allows you to create and configure an Authorization Server Policy.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Auth Server Policy.",
			},
			"auth_server_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the Auth Server.",
			},
			"status": statusSchema,
			"priority": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Priority of the auth server policy",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The description of the Auth Server Policy.",
			},
			"client_whitelist": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "The clients to whitelist the policy for. `[ALL_CLIENTS]` is a special value that can be used to whitelist all clients, otherwise it is a list of client ids.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAuthServerPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.Get("status").(string) == statusInactive {
		return diag.Errorf("can not create an inactive auth server policy, only existing ones can be deactivated")
	}
	policy := buildAuthServerPolicy(d)
	respPolicy, _, err := getOktaClientFromMetadata(meta).AuthorizationServer.CreateAuthorizationServerPolicy(ctx, d.Get("auth_server_id").(string), policy)
	if err != nil {
		return diag.Errorf("failed to create authorization server policy: %v", err)
	}
	d.SetId(respPolicy.Id)
	return resourceAuthServerPolicyRead(ctx, d, meta)
}

func resourceAuthServerPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	policy, resp, err := getOktaClientFromMetadata(meta).AuthorizationServer.GetAuthorizationServerPolicy(ctx, d.Get("auth_server_id").(string), d.Id())
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
	if policy.PriorityPtr != nil {
		_ = d.Set("priority", policy.PriorityPtr)
	}
	_ = d.Set("client_whitelist", convertStringSliceToSet(policy.Conditions.Clients.Include))
	return nil
}

func resourceAuthServerPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	policy := buildAuthServerPolicy(d)
	_, _, err := getOktaClientFromMetadata(meta).AuthorizationServer.UpdateAuthorizationServerPolicy(ctx, d.Get("auth_server_id").(string), d.Id(), policy)
	if err != nil {
		return diag.Errorf("failed to update auth server policy: %v", err)
	}
	oldStatus, newStatus := d.GetChange("status")
	if oldStatus != newStatus {
		if newStatus == statusActive {
			_, err = getOktaClientFromMetadata(meta).AuthorizationServer.ActivateAuthorizationServerPolicy(ctx, d.Get("auth_server_id").(string), d.Id())
		} else {
			_, err = getOktaClientFromMetadata(meta).AuthorizationServer.DeactivateAuthorizationServerPolicy(ctx, d.Get("auth_server_id").(string), d.Id())
		}
		if err != nil {
			return diag.Errorf("failed to change authorization server policy status: %v", err)
		}
	}
	return resourceAuthServerPolicyRead(ctx, d, meta)
}

func resourceAuthServerPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_, err := getOktaClientFromMetadata(meta).AuthorizationServer.DeleteAuthorizationServerPolicy(ctx, d.Get("auth_server_id").(string), d.Id())
	if err != nil {
		return diag.Errorf("failed to delete auth server policy: %v", err)
	}
	return nil
}

func buildAuthServerPolicy(d *schema.ResourceData) sdk.AuthorizationServerPolicy {
	return sdk.AuthorizationServerPolicy{
		Name:        d.Get("name").(string),
		Type:        sdk.OauthAuthorizationPolicyType,
		Status:      d.Get("status").(string),
		PriorityPtr: int64Ptr(d.Get("priority").(int)),
		Description: d.Get("description").(string),
		Conditions: &sdk.PolicyRuleConditions{
			Clients: &sdk.ClientPolicyCondition{
				Include: convertInterfaceToStringSet(d.Get("client_whitelist")),
			},
		},
	}
}
