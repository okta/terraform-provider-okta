package idaas

import (
	"context"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceResourceSets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceResourceSetsRead,
		Schema: map[string]*schema.Schema{
			"resource_sets": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of resource sets",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique ID for the resource set",
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
				},
			},
		},
		Description: "List all resource sets with pagination support. This data source retrieves all resource sets in the organization.",
	}
}

func dataSourceResourceSetsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := getOktaV5ClientFromMetadata(meta).ResourceSetAPI

	// Get all resource sets with pagination
	var allResourceSets []map[string]interface{}
	var after string

	for {
		request := client.ListResourceSets(ctx)
		if after != "" {
			request = request.After(after)
		}

		resourceSets, resp, err := request.Execute()
		if err != nil {
			return diag.Errorf("failed to list resource sets: %v", err)
		}

		// Process the resource sets
		for _, rs := range resourceSets.GetResourceSets() {
			resourceSet := map[string]interface{}{
				"id":          rs.GetId(),
				"label":       rs.GetLabel(),
				"description": rs.GetDescription(),
			}

			if rs.HasCreated() {
				resourceSet["created"] = rs.GetCreated().Format("2006-01-02T15:04:05.000Z")
			}
			if rs.HasLastUpdated() {
				resourceSet["last_updated"] = rs.GetLastUpdated().Format("2006-01-02T15:04:05.000Z")
			}

			allResourceSets = append(allResourceSets, resourceSet)
		}

		// Check for pagination
		if resp.HasNextPage() {
			nextPageURL := resp.NextPage()
			if nextPageURL == "" {
				break
			}
			// Extract the 'after' parameter from the next page URL
			if u, err := url.Parse(nextPageURL); err == nil {
				after = u.Query().Get("after")
				if after == "" {
					break
				}
			} else {
				break
			}
		} else {
			break
		}
	}

	_ = d.Set("resource_sets", allResourceSets)
	d.SetId("resource_sets_list")

	return nil
}
