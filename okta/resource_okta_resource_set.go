package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	v5okta "github.com/okta/okta-sdk-golang/v5/okta"
)

func resourceResourceSet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceResourceSetCreate,
		ReadContext:   resourceResourceSetRead,
		UpdateContext: resourceResourceSetUpdate,
		DeleteContext: resourceResourceSetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: `Manages Resource Sets as custom collections of resources. This resource allows 
			the creation and manipulation of Okta Resource Sets as custom collections of Okta resources. 
			You can use Okta Resource Sets to assign Custom Roles to administrators who are scoped to 
			the designated resources.
			The 'resources' field supports the following:
				- Apps
				- Groups 
				- All Users within a Group
				- All Users within the org
				- All Groups within the org
				- All Apps within the org
				- All Apps of the same type`,
		Schema: map[string]*schema.Schema{
			"label": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique name given to the Resource Set",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A description of the Resource Set",
			},
			"resources": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The endpoints that reference the resources to be included in the new Resource Set. At least one endpoint must be specified when creating resource set.",
			},
		},
	}
}

func resourceResourceSetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := getOktaV5ClientFromMetadata(meta)

	label := d.Get("label").(string)
	description := d.Get("description").(string)
	resources := d.Get("resources").(*schema.Set).List()
	if len(resources) == 0 {
		d.SetId("")
		return diag.Errorf("at least one resource must be specified when creating resource set")
	}

	rs := v5okta.CreateResourceSetRequest{
		Label:       &label,
		Description: &description,
		Resources:   convertInterfaceArrToStringArr(resources),
	}
	apiRequest := client.ResourceSetAPI.CreateResourceSet(ctx).Instance(rs)
	resourceSet, _, err := apiRequest.Execute()
	if err != nil {
		return diag.Errorf("failed to create resource set: %v", err)
	}
	d.SetId(*resourceSet.Id)
	return resourceResourceSetRead(ctx, d, meta)
}

func resourceResourceSetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// rs, resp, err := getAPISupplementFromMetadata(meta).GetResourceSet(ctx, d.Id())
	client := getOktaV5ClientFromMetadata(meta)
	rs, resp, err := client.ResourceSetAPI.GetResourceSet(ctx, d.Id()).Execute()
	if err != nil {
		if v5suppressErrorOn404(resp, err) == nil {
			d.SetId("") // resource set does not exist
			return nil
		}
		return diag.Errorf("failed to get resource set: %v", err)
	}

	if rs == nil { // resource set does not exist
		d.SetId("")
		return nil
	}

	_ = d.Set("label", rs.Label)
	_ = d.Set("description", rs.Description)

	resources, err := listResourceSetResources(ctx, client, d.Id())
	if err != nil {
		return diag.Errorf("failed to get list of resource set resources: %v", err)
	}

	_ = d.Set("resources", flattenResourceSetResources(resources))

	d.SetId(*rs.Id)
	return nil
}

func resourceResourceSetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := getOktaV5ClientFromMetadata(meta)

	if d.HasChanges("label", "description") {
		id := d.Id()
		label := d.Get("label").(string)
		description := d.Get("description").(string)

		rs := v5okta.ResourceSet{
			Id:          &id,
			Label:       &label,
			Description: &description,
		}

		apiRequest := client.ResourceSetAPI.ReplaceResourceSet(ctx, id).Instance(rs)
		if _, _, err := apiRequest.Execute(); err != nil {
			return diag.Errorf("failed to update resource set (ID: %s, Label: %s): %v", id, label, err)
		}
	}

	if !d.HasChange("resources") {
		return resourceResourceSetRead(ctx, d, meta)
	}

	oldResources, newResources := d.GetChange("resources")
	oldSet := oldResources.(*schema.Set)
	newSet := newResources.(*schema.Set)
	resourcesToAdd := convertInterfaceArrToStringArr(newSet.Difference(oldSet).List())
	resourcesToRemove := convertInterfaceArrToStringArr(oldSet.Difference(newSet).List())

	if len(resourcesToAdd) > 0 {
		addResourcesRequest := v5okta.ResourceSetResourcePatchRequest{
			Additions: resourcesToAdd,
		}
		_, _, err := client.ResourceSetAPI.AddResourceSetResource(ctx, d.Id()).Instance(addResourcesRequest).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if len(resourcesToRemove) > 0 {
		err := removeResourcesFromResourceSet(ctx, client, d.Id(), resourcesToRemove)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceResourceSetRead(ctx, d, meta)
}

func resourceResourceSetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	resp, err := getAPISupplementFromMetadata(meta).DeleteResourceSet(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to delete resource set: %v", err)
	}
	return nil
}

func flattenResourceSetResources(resources []v5okta.ResourceSetResource) *schema.Set {
	var arr []interface{}
	for _, res := range resources {
		if res.Links == nil {
			continue
		}
		links := res.Links.Self
		var urlStr string
		if links != nil {
			if href := links.Href; href != "" {
				urlStr = href
			}
		}
		arr = append(arr, urlStr)
	}
	return schema.NewSet(schema.HashString, arr)
}

func listResourceSetResources(ctx context.Context, client *v5okta.APIClient, id string) ([]v5okta.ResourceSetResource, error) {
	var allResources []v5okta.ResourceSetResource
	resources, resp, err := client.ResourceSetAPI.ListResourceSetResources(ctx, id).Execute()
	if err != nil {
		return nil, err
	}
	allResources = append(allResources, resources.Resources...)

	// handle pagination
	for resp.HasNextPage() {
		resp, err = resp.Next(&resources)
		if err != nil {
			return nil, err
		}
		allResources = append(allResources, resources.Resources...)
	}
	return allResources, nil
}

func removeResourcesFromResourceSet(ctx context.Context, client *v5okta.APIClient, resourceSetID string, urls []string) error {
	resources, err := listResourceSetResources(ctx, client, resourceSetID)
	if err != nil {
		return fmt.Errorf("failed to get list of resource set resources: %v", err)
	}

	for _, res := range resources {
		if res.Links == nil {
			continue
		}
		if res.Links.Self == nil {
			continue
		}
		if res.Links.Self.Href == "" {
			continue
		}
		if contains(urls, res.Links.Self.Href) {
			resp, err := client.ResourceSetAPI.DeleteResourceSetResource(ctx, resourceSetID, *res.Id).Execute()
			if err != nil {
				if v5suppressErrorOn404(resp, err) == nil { // couldn't delete resource as we got a 404
					continue
				}
				return fmt.Errorf("failed to remove %s resource from the resource set: %v", *res.Id, err)
			}
		}
	}
	return nil
}
