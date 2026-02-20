package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceResourceSet() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceResourceSetRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID or label of the resource set to retrieve",
			},
			"label": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique name given to the Resource Set",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A description of the Resource Set",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp when the resource set was created",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp when the resource set was last updated",
			},
		},
		Description: "Get a resource set by ID or label. This data source retrieves basic metadata about a resource set including its label, description, and timestamps.",
	}
}

func dataSourceResourceSetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := getOktaV5ClientFromMetadata(meta).ResourceSetAPI
	resourceSetID := d.Get("id").(string)

	// Get the resource set metadata
	rs, _, err := client.GetResourceSet(ctx, resourceSetID).Execute()
	if err != nil {
		return diag.Errorf("failed to get resource set with ID '%s', error: %v", resourceSetID, err)
	}
	if rs == nil {
		return diag.Errorf("resource set with ID '%s' not found", resourceSetID)
	}

	// Set the ID and basic fields
	d.SetId(rs.GetId())
	_ = d.Set("label", rs.GetLabel())
	_ = d.Set("description", rs.GetDescription())

	if rs.HasCreated() {
		_ = d.Set("created", rs.GetCreated().Format("2006-01-02T15:04:05.000Z"))
	}
	if rs.HasLastUpdated() {
		_ = d.Set("last_updated", rs.GetLastUpdated().Format("2006-01-02T15:04:05.000Z"))
	}

	return nil
}
