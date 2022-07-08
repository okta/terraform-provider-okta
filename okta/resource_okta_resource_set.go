package okta

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"github.com/okta/terraform-provider-okta/sdk"
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
		Description: "Resource to manage administrative Role assignments for a User",
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
				Description: "The endpoints that reference the resources to be included in the new Resource Set",
			},
		},
	}
}

func resourceResourceSetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	set, err := buildResourceSet(d, true)
	if err != nil {
		return diag.Errorf("failed to create resource set: %v", err)
	}
	rs, _, err := getSupplementFromMetadata(m).CreateResourceSet(ctx, *set)
	if err != nil {
		return diag.Errorf("failed to create resource set: %v", err)
	}
	d.SetId(rs.Id)
	return resourceResourceSetRead(ctx, d, m)
}

func resourceResourceSetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	rs, resp, err := getSupplementFromMetadata(m).GetResourceSet(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get resource set: %v", err)
	}
	if rs == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("label", rs.Label)
	_ = d.Set("description", rs.Description)
	resources, err := listResourceSetResources(ctx, getSupplementFromMetadata(m), d.Id())
	if err != nil {
		return diag.Errorf("failed to get list of resource set resources: %v", err)
	}
	_ = d.Set("resources", flattenResourceSetResources(resources))
	return nil
}

func resourceResourceSetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getSupplementFromMetadata(m)
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
	resourcesToAdd := convertInterfaceArrToStringArr(newSet.Difference(oldSet).List())
	resourcesToRemove := convertInterfaceArrToStringArr(oldSet.Difference(newSet).List())
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

func resourceResourceSetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resp, err := getSupplementFromMetadata(m).DeleteResourceSet(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
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
		rs.Resources = convertInterfaceToStringSetNullable(d.Get("resources"))
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
		links := res.Links.(map[string]interface{})
		var url string
		for _, v := range links {
			for _, link := range v.(map[string]interface{}) {
				url = link.(string)
				break
			}
		}
		arr = append(arr, url)
	}
	return schema.NewSet(schema.HashString, arr)
}

func listResourceSetResources(ctx context.Context, client *sdk.APISupplement, id string) ([]*sdk.ResourceSetResource, error) {
	var resResources []*sdk.ResourceSetResource
	resources, _, err := client.ListResourceSetResources(ctx, id, &query.Params{Limit: defaultPaginationLimit})
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
		if nextURL := linksValue(resources.Links, "next", "href"); nextURL != "" {
			u, err := url.Parse(nextURL)
			if err != nil {
				break
			}
			// "links": { "next": { "href": "https://host/api/v1/iam/resource-sets/{id}/resources?after={afterId}&limit=100" } }
			after := u.Query().Get("after")
			resources, _, err = client.ListResourceSetResources(ctx, id, &query.Params{After: after, Limit: defaultPaginationLimit})
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
		links := res.Links.(map[string]interface{})
		var url string
		for _, v := range links {
			for _, link := range v.(map[string]interface{}) {
				url = link.(string)
				break
			}
		}
		if contains(urls, url) {
			_, err := client.DeleteResourceSetResource(ctx, resourceSetID, res.Id)
			if err != nil {
				return fmt.Errorf("failed to remove %s resource from the resource set: %v", res.Id, err)
			}
		}
	}
	return nil
}
