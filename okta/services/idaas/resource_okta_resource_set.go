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

// TODO: V6 BREAKING CHANGE - In V6, this resource will be split into:
// 1. okta_resource_set - Manages basic resource set metadata (label, description) only
// 2. okta_resource_set_resource - Manages individual resources with native IDs
// This change will align with Terraform best practices of one resource per API endpoint

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
			// TODO: V6 BREAKING CHANGE - These fields will be moved to okta_resource_set_resource resource
			// and will use native resource IDs instead of extracted links
			"resources": {
				Type:          schema.TypeSet,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"resources_orn"},
				ExactlyOneOf:  []string{"resources", "resources_orn"},
				Description:   "The endpoints that reference the resources to be included in the new Resource Set. At least one endpoint must be specified when creating resource set. Only one of 'resources' or 'resources_orn' can be specified.",
			},
			"resources_orn": {
				Type:          schema.TypeSet,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"resources"},
				ExactlyOneOf:  []string{"resources", "resources_orn"},
				Description:   "The orn(Okta Resource Name) that reference the resources to be included in the new Resource Set. At least one orn must be specified when creating resource set. Only one of 'resources' or 'resources_orn' can be specified.",
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
	if err := utils.SuppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get resource set: %v", err)
	}
	if rs == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("label", rs.Label)
	_ = d.Set("description", rs.Description)

	// TODO: V6 BREAKING CHANGE - This API call and resource processing will be moved to
	// okta_resource_set_resource resource in V6
	resources, err := listResourceSetResources(ctx, getAPISupplementFromMetadata(meta), d.Id())
	if err != nil {
		return diag.Errorf("failed to get list of resource set resources: %v", err)
	}
	if _, ok := d.GetOk("resources"); ok {
		_ = d.Set("resources", flattenResourceSetResourcesLinks(resources))
	} else if _, ok := d.GetOk("resources_orn"); ok {
		_ = d.Set("resources_orn", flattenResourceSetResourcesORN(resources))
	}
	return nil
}

func resourceResourceSetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := getAPISupplementFromMetadata(meta)
	if d.HasChanges("label", "description") {
		set, _ := buildResourceSet(d, false)
		_, _, err := client.UpdateResourceSet(ctx, d.Id(), *set)
		if err != nil {
			return diag.Errorf("failed to update resource set: %v", err)
		}
	}
	if !d.HasChange("resources") && !d.HasChange("resources_orn") {
		return nil
	}

	// TODO: V6 BREAKING CHANGE - This resource management logic will be moved to
	// okta_resource_set_resource resource in V6
	if d.HasChange("resources") {
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
	}

	if d.HasChange("resources_orn") {
		oldResourcesOrn, newResourcesOrn := d.GetChange("resources_orn")
		oldSetOrn := oldResourcesOrn.(*schema.Set)
		newSetOrn := newResourcesOrn.(*schema.Set)
		ornResourcesToAdd := utils.ConvertInterfaceArrToStringArr(newSetOrn.Difference(oldSetOrn).List())
		ornResourcesToRemove := utils.ConvertInterfaceArrToStringArr(oldSetOrn.Difference(newSetOrn).List())

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
	if err := utils.SuppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to delete resource set: %v", err)
	}
	return nil
}

// TODO: V6 BREAKING CHANGE - This function will be simplified to only handle basic metadata
// Resource management will be moved to okta_resource_set_resource resource
func buildResourceSet(d *schema.ResourceData, isNew bool) (*sdk.ResourceSet, error) {
	rs := &sdk.ResourceSet{
		Label:       d.Get("label").(string),
		Description: d.Get("description").(string),
	}
	if isNew {
		resourceLinks := utils.ConvertInterfaceToStringSetNullable(d.Get("resources"))
		resourceOrns := utils.ConvertInterfaceToStringSetNullable(d.Get("resources_orn"))

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

// TODO: V6 BREAKING CHANGE - This function will be moved to okta_resource_set_resource resource
// and will use native IDs instead of extracted links
func flattenResourceSetResourcesLinks(resources []*sdk.ResourceSetResource) *schema.Set {
	var arr []interface{}
	for _, res := range resources {
		if res.Links != nil {
			arr = append(arr, encodeResourceSetResourceLink(utils.LinksValue(res.Links, "self", "href")))
		}
	}
	return schema.NewSet(schema.HashString, arr)
}

// TODO: V6 BREAKING CHANGE - This function will be moved to okta_resource_set_resource resource
// and will use native IDs instead of extracted links
func flattenResourceSetResourcesORN(resources []*sdk.ResourceSetResource) *schema.Set {
	var arr []interface{}
	for _, res := range resources {
		if res.Orn != "" {
			urlStr := res.Orn
			arr = append(arr, urlStr)
		}
	}
	return schema.NewSet(schema.HashString, arr)
}

// TODO: V6 BREAKING CHANGE - This function will be moved to okta_resource_set_resource resource
// and will use native IDs instead of extracted links
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

// TODO: V6 BREAKING CHANGE - This function will be moved to okta_resource_set_resource resource
// and will use native IDs instead of extracted links
func addResourcesToResourceSet(ctx context.Context, client *sdk.APISupplement, resourceSetID string, links []string) error {
	if len(links) == 0 {
		return nil
	}
	var encodedLinks []string
	for _, link := range links {
		encodedLinks = append(encodedLinks, encodeResourceSetResourceLink(link))
	}
	_, err := client.AddResourceSetResources(ctx, resourceSetID, sdk.AddResourceSetResourcesRequest{Additions: encodedLinks})
	if err != nil {
		return fmt.Errorf("failed to add resources to the resource set: %v", err)
	}
	return nil
}

// TODO: V6 BREAKING CHANGE - This function will be moved to okta_resource_set_resource resource
// and will use native IDs instead of extracted links
func removeResourcesFromResourceSet(ctx context.Context, client *sdk.APISupplement, resourceSetID string, urls []string) error {
	resources, err := listResourceSetResources(ctx, client, resourceSetID)
	if err != nil {
		return fmt.Errorf("failed to get list of resource set resources: %v", err)
	}
	var escapedUrls []string
	for _, u := range urls {
		u1, err := escapeResourceSetResourceLink(u)
		if err != nil {
			return fmt.Errorf("failed to escape resource set resource link: %v", err)
		}
		escapedUrls = append(escapedUrls, u1)
	}

	for _, res := range resources {
		orn := res.Orn
		toDelete := false

		if res.Links != nil {
			url := utils.LinksValue(res.Links, "self", "href")
			toDelete = utils.Contains(escapedUrls, url)
		}

		toDelete = toDelete || utils.Contains(escapedUrls, orn)

		if toDelete {
			_, err := client.DeleteResourceSetResource(ctx, resourceSetID, res.Id)
			if err != nil {
				return fmt.Errorf("failed to remove %s resource from the resource set: %v", res.Id, err)
			}
		}
	}
	return nil
}

// TODO: V6 BREAKING CHANGE - This function will be moved to okta_resource_set_resource resource
// and will use native IDs instead of extracted links
func encodeResourceSetResourceLink(link string) string {
	parsedBaseUrl, err := url.Parse(link)
	if err != nil {
		return ""
	}

	filter := parsedBaseUrl.Query().Get("filter")
	q := parsedBaseUrl.Query()
	if filter != "" {
		q.Set("filter", filter)
	}
	parsedBaseUrl.RawQuery = q.Encode()
	return parsedBaseUrl.String()
}

// TODO: V6 BREAKING CHANGE - This function will be moved to okta_resource_set_resource resource
// and will use native IDs instead of extracted links
func escapeResourceSetResourceLink(link string) (string, error) {
	// Parse the URL first
	u, err := url.Parse(link)
	if err != nil {
		return link, fmt.Errorf("error parsing URL %s: %v", link, err)
	}

	// Escape query parameters using url.QueryEscape for each query parameter
	q := u.Query()

	// Escape each value in the query parameters
	for key, values := range q {
		for i, v := range values {
			// URL escape the value
			q[key][i] = url.QueryEscape(v)
		}
	}

	// Rebuild the URL with the modified query parameters
	u.RawQuery = q.Encode()
	return u.String(), nil
}
