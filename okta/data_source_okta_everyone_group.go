package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// data source to retrieve information on the Everyone Group
func dataSourceEveryoneGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEveryoneGroupRead,
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"include_users": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Fetch group users, having default off cuts down on API calls.",
			},
		},
	}
}

func dataSourceEveryoneGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return findGroup(ctx, groupProfileEveryone, d, m, true)
}
