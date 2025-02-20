package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// data source to retrieve information on the Everyone Group
func DataSourceEveryoneGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEveryoneGroupRead,
		Schema: map[string]*schema.Schema{
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of group.",
			},
			"include_users": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Fetch group users, having default off cuts down on API calls.",
			},
		},
		Description: "Get the `Everyone` group from Okta.",
	}
}

func dataSourceEveryoneGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return findGroup(ctx, GroupProfileEveryone, d, meta, true)
}
