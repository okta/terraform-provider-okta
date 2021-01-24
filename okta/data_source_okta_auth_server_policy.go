package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAuthServerPolicy() *schema.Resource {
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
				Type:     schema.TypeString,
				Computed: true,
			},
			"assigned_clients": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceAuthServerPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	policies, _, err := getOktaClientFromMetadata(m).AuthorizationServer.ListAuthorizationServerPolicies(ctx, d.Get("auth_server_id").(string))
	if err != nil {
		return diag.Errorf("failed to list auth server policies: %v", err)
	}
	name := d.Get("name").(string)
	for _, policy := range policies {
		if policy.Name == name {
			d.SetId(policy.Id)
			_ = d.Set("description", policy.Description)
			_ = d.Set("assigned_clients", convertStringSetToInterface(policy.Conditions.Clients.Include))
			return nil
		}
	}
	return diag.Errorf("auth server policy with name '%s' does not exist", name)
}
