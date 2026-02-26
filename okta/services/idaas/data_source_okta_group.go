package idaas

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/okta/terraform-provider-okta/sdk/query"
)

func dataSourceGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGroupRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeString,
				Computed:      true,
				Optional:      true,
				ConflictsWith: []string{"name", "type"},
				Description:   "ID of group.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Name of group.",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Type of the group. When specified in the terraform resource, will act as a filter when searching for the group",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of group.",
			},
			"include_users": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Fetch group users, having default off cuts down on API calls.",
			},
			"users": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Users associated with the group. This can also be done per user.",
			},
			"delay_read_seconds": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Force delay of the group read by N seconds. Useful when eventual consistency of group information needs to be allowed for; for instance, when group rules are known to have been applied.",
			},
		},
		Description: "Get a group from Okta.",
	}
}

func dataSourceGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if n, ok := d.GetOk("delay_read_seconds"); ok {
		delay, err := strconv.Atoi(n.(string))
		if err == nil {
			logger(meta).Info("delaying group read by ", delay, " seconds")
			meta.(*config.Config).TimeOperations.Sleep(time.Duration(delay) * time.Second)
		} else {
			logger(meta).Warn("group read delay value ", n, " is not an integer")
		}
	}

	return findGroup(ctx, d.Get("name").(string), d, meta, false)
}

func findGroup(ctx context.Context, name string, d *schema.ResourceData, meta interface{}, isEveryone bool) diag.Diagnostics {
	var group *sdk.Group
	groupID, ok := d.GetOk("id")
	if ok {
		respGroup, _, err := getOktaClientFromMetadata(meta).Group.GetGroup(ctx, groupID.(string))
		if err != nil {
			return diag.Errorf("failed get group by ID: %v", err)
		}
		group = respGroup
	} else {
		// NOTE: Okta API query on name is effectively a starts with query, not
		// an exact match query.
		searchParams := &query.Params{Q: name}

		// NOTE: Uncertain of the OKTA /api/v1/groups API drifted during Classic
		// (when this data source was originally created) to OIE migration.
		// Currently, Okta API enforces unique names on all groups regardless of
		// type so type is essentially a meaningless parameter in OIE.  There
		// may be a case where imported groups allow duplicate names.
		//
		// As of Dec 2025 - Okta only enfoces unique names for groups with type `OKTA_GROUP`
		// groups imported from external systems are given type `APP_GROUP`
		// `APP_GROUP` can have the same name as an `OKTA_GROUP` or `APP_GROUP` from other import sources
		// e.g. a google group and a salesforce group can have the same name and will still be imported to Okta successfully
		t, okType := d.GetOk("type")
		if okType {
			searchParams.Filter = fmt.Sprintf("type eq \"%s\"", t.(string))
		}

		logger(meta).Info("looking for data source group", "query", searchParams.String())
		groups, _, err := getOktaClientFromMetadata(meta).Group.ListGroups(ctx, searchParams)
		if err != nil {
			return diag.Errorf("failed to query for groups: %v", err)
		}
		if len(groups) > 1 {
			logger(meta).Warn("data source group query matches", len(groups), "groups")
			for _, g := range groups {
				// exact match on name
				if g.Profile.Name == name {
					if okType && t.(string) == g.Type {
						// data source has type argument so take that into consideration also
						group = g
						break
					}
					if !okType {
						// otherwise consider name only
						group = g
						break
					}
				}
			}
			if group == nil {
				if okType {
					return diag.Errorf("group starting with name %q and type %q matches %d groups, select a more precise name parameter", name, d.Get("type").(string), len(groups))
				}
				return diag.Errorf("group starting with name %q matches %d groups, select a more precise name parameter", name, len(groups))
			}
		}
		if len(groups) < 1 {
			if okType {
				return diag.Errorf("group with name %q and type %q does not exist", name, d.Get("type").(string))
			}
			return diag.Errorf("group with name %q does not exist", name)
		}
		if len(groups) == 1 {
			group = groups[0]
			if group.Profile.Name != name {
				// keep old behavior that a fuzzy match is acceptable if query only returns one group
				logger(meta).Warn("group with exact name match was not found: using partial match which contains name as a substring", "name", group.Profile.Name)
			}
		}
	}
	d.SetId(group.Id)
	_ = d.Set("description", group.Profile.Description)
	if !isEveryone {
		_ = d.Set("type", group.Type)
		_ = d.Set("name", group.Profile.Name)
	}
	if !d.Get("include_users").(bool) {
		return nil
	}
	userIDList, err := listGroupUserIDs(ctx, meta, d.Id())
	if err != nil {
		return diag.Errorf("failed to list group user IDs: %v", err)
	}
	_ = d.Set("users", utils.ConvertStringSliceToSet(userIDList))
	return nil
}
