package idaas

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/okta/terraform-provider-okta/sdk/query"
)

func ResourceResourceSet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceResourceSetCreate,
		ReadContext:   resourceResourceSetRead,
		UpdateContext: resourceResourceSetUpdate,
		DeleteContext: resourceResourceSetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: `Manages Resource Sets as custom collections of resources. This resource allows the creation and manipulation of Okta Resource Sets as custom collections of Okta resources. You can use Okta Resource Sets to assign Custom Roles to administrators who are scoped to the designated resources. 
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
	set, err := buildResourceSet(d, true)
	if err != nil {
		return diag.Errorf("failed to create resource set: %v", err)
	}
	rs, _, err := GetAPISupplementFromMetadata(meta).CreateResourceSet(ctx, *set)
	if err != nil {
		return diag.Errorf("failed to create resource set: %v", err)
	}
	d.SetId(rs.Id)
	return resourceResourceSetRead(ctx, d, meta)
}

func resourceResourceSetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rs, resp, err := GetAPISupplementFromMetadata(meta).GetResourceSet(ctx, d.Id())
	if err := utils.SuppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get resource set: %v", err)
	}
	if rs == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("label", rs.Label)
	_ = d.Set("description", rs.Description)
	resources, err := listResourceSetResources(ctx, GetAPISupplementFromMetadata(meta), d.Id())
	if err != nil {
		return diag.Errorf("failed to get list of resource set resources: %v", err)
	}
	_ = d.Set("resources", flattenResourceSetResources(resources))
	return nil
}

func resourceResourceSetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := GetAPISupplementFromMetadata(meta)
	if d.HasChanges("label", "description") {
		set, _ := buildResourceSet(d, false)
		_, _, err := client.UpdateResourceSet(ctx, d.Id(), *set)
		if err != nil {
			return diag.Errorf("failed to update resource set: %v", err)
		}
	}
	if !d.HasChange("resources") {
		return nil
	}
	oldResources, newResources := d.GetChange("resources")
	oldSet := oldResources.(*schema.Set)
	newSet := newResources.(*schema.Set)
	resourcesToAdd := utils.ConvertInterfaceArrToStringArr(newSet.Difference(oldSet).List())
	resourcesToRemove := utils.ConvertInterfaceArrToStringArr(oldSet.Difference(newSet).List())
	err := addResourcesToResourceSet(ctx, client, d.Id(), resourcesToAdd)
	if err != nil {
		return diag.FromErr(err)
	}
	err = removeResourcesFromResourceSet(ctx, client, d.Id(), resourcesToRemove)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceResourceSetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	resp, err := GetAPISupplementFromMetadata(meta).DeleteResourceSet(ctx, d.Id())
	if err := utils.SuppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to delete resource set: %v", err)
	}
	return nil
}

func buildResourceSet(d *schema.ResourceData, isNew bool) (*sdk.ResourceSet, error) {
	rs := &sdk.ResourceSet{
		Label:       d.Get("label").(string),
		Description: d.Get("description").(string),
	}
	if isNew {
		rs.Resources = utils.ConvertInterfaceToStringSetNullable(d.Get("resources"))
		if len(rs.Resources) == 0 {
			return nil, errors.New("at least one resource must be specified when creating resource set")
		}
	} else {
		rs.Id = d.Id()
	}
	return rs, nil
}

func flattenResourceSetResources(resources []*sdk.ResourceSetResource) *schema.Set {
	var arr []interface{}
	for _, res := range resources {
		if res.Links != nil {
			links := res.Links.(map[string]interface{})
			var urlStr string
			if selfURL, ok := links["self"]; ok {
				if _, ok := selfURL.(map[string]interface{}); ok {
					for _, link := range selfURL.(map[string]interface{}) {
						urlStr = link.(string)
						break
					}
				}
			}
			arr = append(arr, urlStr)
		}
	}
	return schema.NewSet(schema.HashString, arr)
}

func listResourceSetResources(ctx context.Context, client *sdk.APISupplement, id string) ([]*sdk.ResourceSetResource, error) {
	var resResources []*sdk.ResourceSetResource
	resources, _, err := client.ListResourceSetResources(ctx, id, &query.Params{Limit: utils.DefaultPaginationLimit})
	if err != nil {
		return nil, err
	}
	resResources = append(resResources, resources.Resources...)
	for {
		// NOTE: The resources endpoint /api/v1/iam/resource-sets/%s/resources
		// is not returning pagination in the headers. Make use of the _links
		// object in the response body. Convert implemenation style back to
		// resp.HasNextPage() if/when that endpoint starts to have pagination
		// information in its headers and/or when this code is supported by
		// okta-sdk-golang instead of the local SDK.
		if nextURL := utils.LinksValue(resources.Links, "next", "href"); nextURL != "" {
			u, err := url.Parse(nextURL)
			if err != nil {
				break
			}
			// "links": { "next": { "href": "https://host/api/v1/iam/resource-sets/{id}/resources?after={afterId}&limit=100" } }
			after := u.Query().Get("after")
			resources, _, err = client.ListResourceSetResources(ctx, id, &query.Params{After: after, Limit: utils.DefaultPaginationLimit})
			if err != nil {
				return nil, err
			}
			resResources = append(resResources, resources.Resources...)
		} else {
			break
		}
	}
	return resResources, nil
}

func addResourcesToResourceSet(ctx context.Context, client *sdk.APISupplement, resourceSetID string, links []string) error {
	if len(links) == 0 {
		return nil
	}
	_, err := client.AddResourceSetResources(ctx, resourceSetID, sdk.AddResourceSetResourcesRequest{Additions: links})
	if err != nil {
		return fmt.Errorf("failed to add resources to the resource set: %v", err)
	}
	return nil
}

func removeResourcesFromResourceSet(ctx context.Context, client *sdk.APISupplement, resourceSetID string, urls []string) error {
	resources, err := listResourceSetResources(ctx, client, resourceSetID)
	if err != nil {
		return fmt.Errorf("failed to get list of resource set resources: %v", err)
	}
	for _, res := range resources {
		if res.Links == nil {
			continue
		}
		links := res.Links.(map[string]interface{})
		var url string
		for _, v := range links {
			for _, link := range v.(map[string]interface{}) {
				url = link.(string)
				break
			}
		}
		if utils.Contains(urls, url) {
			_, err := client.DeleteResourceSetResource(ctx, resourceSetID, res.Id)
			if err != nil {
				return fmt.Errorf("failed to remove %s resource from the resource set: %v", res.Id, err)
			}
		}
	}
	return nil
}
