package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePolicyRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the policy",
				Required:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "Policy type, see https://developer.okta.com/docs/reference/api/policy/#policy-object",
				Required:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the policy",
			},
		},
		Description: "Get a policy from Okta.",
	}
}

func dataSourcePolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	policy, err := findPolicyByNameAndType(ctx, meta, d.Get("name").(string), d.Get("type").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(policy.Id)
	_ = d.Set("status", policy.Status)
	return nil
}
