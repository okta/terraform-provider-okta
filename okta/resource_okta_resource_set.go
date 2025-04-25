package okta

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/okta/terraform-provider-okta/sdk/query"
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
				Type:          schema.TypeSet,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"resources_orn"},
				ExactlyOneOf:  []string{"resources", "resources_orn"},
				Description:   "The endpoints that reference the resources to be included in the new Resource Set. At least one endpoint must be specified when creating resource set.",
			},
			"resources_orn": {
				Type:          schema.TypeSet,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"resources"},
				ExactlyOneOf:  []string{"resources", "resources_orn"},
				Description:   "The endpoints that reference the resources to be included in the new Resource Set. At least one endpoint must be specified when creating resource set.",
			},
		},
	}
}

func resourceResourceSetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	set, err := buildResourceSet(d, true)
	if err != nil {
		return diag.Errorf("failed to create resource set: %v", err)
	}
	rs, _, err := getAPISupplementFromMetadata(meta).CreateResourceSet(ctx, *set)
	if err != nil {
		return diag.Errorf("failed to create resource set: %v", err)
	}
	d.SetId(rs.Id)
	return resourceResourceSetRead(ctx, d, meta)
}

func resourceResourceSetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rs, resp, err := getAPISupplementFromMetadata(meta).GetResourceSet(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get resource set: %v", err)
	}
	if rs == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("label", rs.Label)
	_ = d.Set("description", rs.Description)
	resources, err := listResourceSetResources(ctx, getAPISupplementFromMetadata(meta), d.Id())
	if err != nil {
		return diag.Errorf("failed to get list of resource set resources: %v", err)
	}
	logger(meta).Info("Got resource set", "links", flattenResourceSetResourcesLinks(resources))
	logger(meta).Info("Got resource set", "orn", flattenResourceSetResourcesORN(resources))
	if _, ok := d.GetOk("resources"); ok {
		_ = d.Set("resources", flattenResourceSetResourcesLinks(resources))
	} else if _, ok := d.GetOk("resources_orn"); ok {
		_ = d.Set("resources_orn", flattenResourceSetResourcesORN(resources))
	}
	return nil
}

//how do i know what value to compare against, since i will get both orn and rest api url, but how do i know whether to compare url or the orns with the .tf config file

func resourceResourceSetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := getAPISupplementFromMetadata(meta)
	if d.HasChanges("label", "description") {
		set, _ := buildResourceSet(d, false)
		logger(meta).Info("Updating resource set", "resources", set)
		_, _, err := client.UpdateResourceSet(ctx, d.Id(), *set)
		if err != nil {
			return diag.Errorf("failed to update resource set: %v", err)
		}
	}
	if !d.HasChange("resources") && !d.HasChange("resources_orn") {
		return nil
	}

	logger(meta).Info("Updating resource set", "Has change", d.HasChange("resources"))
	if d.HasChange("resources") {
		oldResources, newResources := d.GetChange("resources")
		oldSet := oldResources.(*schema.Set)
		newSet := newResources.(*schema.Set)
		resourcesToAdd := convertInterfaceArrToStringArr(newSet.Difference(oldSet).List())
		resourcesToRemove := convertInterfaceArrToStringArr(oldSet.Difference(newSet).List())

		logger(meta).Info("Updating resource set", "resourcesToAdd", resourcesToAdd, "resourcesToRemove", resourcesToRemove)
		err := addResourcesToResourceSet(ctx, client, d.Id(), resourcesToAdd)
		if err != nil {
			return diag.FromErr(err)
		}
		err = removeResourcesFromResourceSet(ctx, client, d.Id(), resourcesToRemove)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	//return nil

	if d.HasChange("resources_orn") {
		oldResourcesOrn, newResourcesOrn := d.GetChange("resources_orn")
		oldSetOrn := oldResourcesOrn.(*schema.Set)
		newSetOrn := newResourcesOrn.(*schema.Set)
		ornResourcesToAdd := convertInterfaceArrToStringArr(newSetOrn.Difference(oldSetOrn).List())
		ornResourcesToRemove := convertInterfaceArrToStringArr(oldSetOrn.Difference(newSetOrn).List())

		logger(meta).Info("Updating resource set", "resourcesToAdd", ornResourcesToAdd, "resourcesToRemove", ornResourcesToRemove)
		err := addResourcesToResourceSet(ctx, client, d.Id(), ornResourcesToAdd)
		if err != nil {
			return diag.FromErr(err)
		}
		err = removeResourcesFromResourceSet(ctx, client, d.Id(), ornResourcesToRemove)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}

func resourceResourceSetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	resp, err := getAPISupplementFromMetadata(meta).DeleteResourceSet(ctx, d.Id())
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
		resourceLinks := convertInterfaceToStringSetNullable(d.Get("resources"))
		resourceOrns := convertInterfaceToStringSetNullable(d.Get("resources_orn"))

		var resource []string
		resource = append(resource, resourceLinks...)
		resource = append(resource, resourceOrns...)

		rs.Resources = resource
		if len(rs.Resources) == 0 {
			return nil, errors.New("at least one resource must be specified when creating resource set")
		}
	} else {
		rs.Id = d.Id()
	}
	return rs, nil
}

func flattenResourceSetResourcesLinks(resources []*sdk.ResourceSetResource) *schema.Set {
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
			fmt.Println("flattenResourceSetResourcesLinks", urlStr)
			arr = append(arr, urlStr)
		}
	}
	return schema.NewSet(schema.HashString, arr)
}

func flattenResourceSetResourcesORN(resources []*sdk.ResourceSetResource) *schema.Set {
	var arr []interface{}
	for _, res := range resources {
		fmt.Println("Mapping resource set resources to ORN", "additional properties", res.Orn)
		if res.Orn != "" {
			orns := res.Orn
			var urlStr string
			urlStr = orns
			fmt.Println("flattenResourceSetResourcesLinks orn", urlStr)
			arr = append(arr, urlStr)
		}
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
	fmt.Println("AddResourceSetResources", resourceSetID, links)
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
		if contains(urls, url) {
			_, err := client.DeleteResourceSetResource(ctx, resourceSetID, res.Id)
			if err != nil {
				return fmt.Errorf("failed to remove %s resource from the resource set: %v", res.Id, err)
			}
		}
	}
	return nil
}
