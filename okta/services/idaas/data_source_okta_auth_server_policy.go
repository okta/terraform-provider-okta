package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

func DataSourceAuthServerPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAuthServerPolicyRead,
		Schema: map[string]*schema.Schema{
			"auth_server_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Auth server ID",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the policy",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of authorization server policy.",
			},
			"assigned_clients": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of clients this policy is assigned to. `[ALL_CLIENTS]` is a special value when policy is assigned to all clients.",
			},
			"priority": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Priority of the auth server policy",
			},
		},
		Description: "Get an authorization server policy from Okta.",
	}
}

func dataSourceAuthServerPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	policies, _, err := GetOktaClientFromMetadata(meta).AuthorizationServer.ListAuthorizationServerPolicies(ctx, d.Get("auth_server_id").(string))
	if err != nil {
		return diag.Errorf("failed to list auth server policies: %v", err)
	}
	name := d.Get("name").(string)
	for _, policy := range policies {
		if policy.Name != name {
			continue
		}
		d.SetId(policy.Id)
		_ = d.Set("description", policy.Description)
		_ = d.Set("assigned_clients", utils.ConvertStringSliceToSet(policy.Conditions.Clients.Include))
		if policy.PriorityPtr != nil {
			_ = d.Set("priority", policy.PriorityPtr)
		}
		return nil
	}
	return diag.Errorf("auth server policy with name '%s' does not exist", name)
}
