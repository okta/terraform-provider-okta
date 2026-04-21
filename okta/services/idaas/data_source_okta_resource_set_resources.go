package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceResourceSetResources() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceResourceSetResourcesRead,
		Schema: map[string]*schema.Schema{
			"resource_set_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the resource set to retrieve resources for",
			},
			"resources": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of resources in the resource set",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique ID of the resource set resource object",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the resource",
						},
						"created": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Timestamp when the resource set resource object was created",
						},
						"last_updated": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Timestamp when this object was last updated",
						},
					},
				},
			},
		},
		Description: "Get the resources that make up a specific resource set. This data source retrieves all resources associated with a resource set including their metadata.",
	}
}

func dataSourceResourceSetResourcesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := getOktaV5ClientFromMetadata(meta).ResourceSetAPI
	resourceSetID := d.Get("resource_set_id").(string)

	// Get resources in the resource set
	resourceSetResources, _, err := client.ListResourceSetResources(ctx, resourceSetID).Execute()
	if err != nil {
		return diag.Errorf("failed to list resource set resources for resource set ID '%s': %v", resourceSetID, err)
	}

	// Process the resources
	var allResources []map[string]interface{}
	for _, resource := range resourceSetResources.GetResources() {
		resourceData := map[string]interface{}{
			"id": resource.GetId(),
		}

		if resource.HasDescription() {
			resourceData["description"] = resource.GetDescription()
		}
		if resource.HasCreated() {
			resourceData["created"] = resource.GetCreated().Format("2006-01-02T15:04:05.000Z")
		}
		if resource.HasLastUpdated() {
			resourceData["last_updated"] = resource.GetLastUpdated().Format("2006-01-02T15:04:05.000Z")
		}

		allResources = append(allResources, resourceData)
	}

	_ = d.Set("resources", allResources)
	d.SetId(resourceSetID)

	return nil
}
