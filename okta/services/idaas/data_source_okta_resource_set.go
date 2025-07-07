package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

// TODO: BREAKING CHANGE - In a future major release, this datasource will be split into:
// 1. okta_resource_set - Returns basic metadata (id, label, description, timestamps)
// 2. okta_resource_set_resources - Returns resources with native IDs instead of extracted links
// This change will align with Terraform best practices of one resource per API endpoint
// and use native IDs instead of brittle link extraction.

func dataSourceResourceSet() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceResourceSetRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the resource set to retrieve",
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
			// TODO: BREAKING CHANGE - These fields will be moved to okta_resource_set_resources datasource
			// and will use native resource IDs instead of extracted links
			"resources": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The endpoints that reference the resources included in the Resource Set",
			},
			"resources_orn": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The orn(Okta Resource Name) that reference the resources included in the Resource Set",
			},
		},
		Description: "Get a resource set by ID. This data source allows you to retrieve the details of a resource set including its resources, which can be used in lifecycle preconditions to prevent users from being granted admin over themselves.",
	}
}

func dataSourceResourceSetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := getAPISupplementFromMetadata(meta)
	resourceSetID := d.Get("id").(string)

	// Get the resource set metadata
	rs, resp, err := client.GetResourceSet(ctx, resourceSetID)
	if err := utils.SuppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get resource set with ID '%s', error: %v", resourceSetID, err)
	}
	if rs == nil {
		return diag.Errorf("resource set with ID '%s' not found", resourceSetID)
	}

	// Set the ID and basic fields
	d.SetId(rs.Id)
	_ = d.Set("label", rs.Label)
	_ = d.Set("description", rs.Description)

	// TODO: BREAKING CHANGE - This API call and resource processing will be moved to
	// okta_resource_set_resources datasource in a future major release
	// Get the resources in the resource set using the shared utility function
	resources, err := listResourceSetResources(ctx, client, resourceSetID)
	if err != nil {
		return diag.Errorf("failed to get list of resource set resources for resource set ID '%s', error: %v", resourceSetID, err)
	}

	// TODO: BREAKING CHANGE - This brittle link extraction will be replaced with native IDs
	// Use the shared utility functions to process resources consistently with the resource
	linksSet := flattenResourceSetResourcesLinks(resources)
	ornsSet := flattenResourceSetResourcesORN(resources)

	// Set the appropriate fields based on what we found
	if linksSet.Len() > 0 {
		_ = d.Set("resources", linksSet)
	}
	if ornsSet.Len() > 0 {
		_ = d.Set("resources_orn", ornsSet)
	}

	return nil
}
