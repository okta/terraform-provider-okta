package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

func dataSourceUserProfileMappingSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserProfileMappingSourceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the source",
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceUserProfileMappingSourceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	mapping, err := getSupplementFromMetadata(m).FindProfileMappingSource(ctx, "user", "user", &query.Params{Limit: defaultPaginationLimit})
	if err != nil {
		return diag.Errorf("failed to find profile mapping source: %v", err)
	}
	d.SetId(mapping.ID)
	_ = d.Set("type", mapping.Type)
	_ = d.Set("name", mapping.Name)
	return nil
}
