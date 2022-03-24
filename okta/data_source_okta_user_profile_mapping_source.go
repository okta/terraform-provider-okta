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
	mappings, resp, err := getOktaClientFromMetadata(m).ProfileMapping.ListProfileMappings(ctx, &query.Params{Limit: defaultPaginationLimit})
	if err != nil {
		return diag.Errorf("failed to list mappings: %v", err)
	}
	name := "user"
	typ := "user"
	for {
		for _, mapping := range mappings {
			target := mapping.Target
			source := mapping.Source
			if target.Name == name && target.Type == typ {
				d.SetId(target.Id)
				_ = d.Set("type", target.Type)
				_ = d.Set("name", target.Name)
				return nil
			} else if source.Name == name && source.Type == typ {
				d.SetId(source.Id)
				_ = d.Set("type", source.Type)
				_ = d.Set("name", source.Name)
				return nil
			}
		}
		if resp.HasNextPage() {
			resp, err = resp.Next(ctx, &mappings)
			if err != nil {
				return diag.Errorf("failed to find profile mapping source: %v", err)
			}
			continue
		} else {
			break
		}
	}

	return nil
}
