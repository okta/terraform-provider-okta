package idaas

import (
	"context"
	"fmt"
	"hash/crc32"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/okta/terraform-provider-okta/sdk/query"
)

func dataSourceUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUsersRead,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Find users based on group membership using the id of the group.",
				ConflictsWith: []string{"search"},
			},
			"include_groups": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Fetch group memberships for each user",
			},
			"include_roles": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Fetch user roles for each user",
			},
			"search": {
				Type:          schema.TypeSet,
				Optional:      true,
				Description:   userSearchSchemaDescription,
				ConflictsWith: []string{"group_id"},
				Elem: &schema.Resource{
					Schema: userSearchSchema,
				},
			},
			"users": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "collection of users retrieved from Okta.",
				Elem: &schema.Resource{
					Schema: utils.BuildSchema(userProfileDataSchema,
						map[string]*schema.Schema{
							"id": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"realm_id": {
								Type:        schema.TypeString,
								Computed:    true,
								Description: "The Realm ID associated with the user.",
							},
						}),
				},
			},
			"compound_search_operator": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "and",
				Description: "Search operator used when joining multiple search clauses",
			},
			"delay_read_seconds": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Force delay of the users read by N seconds. Useful when eventual consistency of users information needs to be allowed for.",
			},
		},
		Description: "Get a list of users from Okta.",
	}
}

func dataSourceUsersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if n, ok := d.GetOk("delay_read_seconds"); ok {
		delay, err := strconv.Atoi(n.(string))
		if err == nil {
			logger(meta).Info("delaying users read by ", delay, " seconds")
			meta.(*config.Config).TimeOperations.Sleep(time.Duration(delay) * time.Second)
		} else {
			logger(meta).Warn("users read delay value ", n, " is not an integer")
		}
	}

	var (
		users []*sdk.User
		id    string
		err   error
	)

	client := getOktaClientFromMetadata(meta)

	if groupId, ok := d.GetOk("group_id"); ok {
		id = groupId.(string)
		users, err = listGroupUsers(ctx, meta, id)
	} else if _, ok := d.GetOk("search"); ok {
		params := &query.Params{Search: getSearchCriteria(d), Limit: utils.DefaultPaginationLimit, SortOrder: "0"}
		id = fmt.Sprintf("%d", crc32.ChecksumIEEE([]byte(params.String())))
		users, err = collectUsers(ctx, client, params)
	} else {
		return diag.Errorf("must specify either group_id or search attributes")
	}

	if err != nil {
		return diag.Errorf("failed to list users: %v", err)
	}
	d.SetId(id)
	includeGroups := d.Get("include_groups").(bool)
	includeRoles := d.Get("include_roles").(bool)
	arr := make([]map[string]interface{}, len(users))
	for i, user := range users {
		rawMap := flattenUser(user, []string{})
		rawMap["id"] = user.Id
		if includeGroups {
			groups, err := getGroupsForUser(ctx, user.Id, client)
			if err != nil {
				return diag.Errorf("failed to list users: %v", err)
			}
			rawMap["group_memberships"] = groups
		}
		if includeRoles {
			roles, _, err := getAdminRoles(ctx, user.Id, client)
			if err != nil {
				return diag.Errorf("failed to set user's admin roles: %v", err)
			}
			rawMap["admin_roles"] = roles
		}
		arr[i] = rawMap
	}

	_ = d.Set("users", arr)
	return nil
}

func collectUsers(ctx context.Context, client *sdk.Client, qp *query.Params) ([]*sdk.User, error) {
	users, resp, err := client.User.ListUsers(ctx, qp)
	if err != nil {
		return nil, err
	}
	for resp.HasNextPage() {
		var nextUsers []*sdk.User
		resp, err = resp.Next(ctx, &nextUsers)
		if err != nil {
			return nil, err
		}
		users = append(users, nextUsers...)
	}
	return users, nil
}
