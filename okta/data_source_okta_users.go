package okta

import (
	"context"
	"fmt"
	"hash/crc32"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

func dataSourceUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUsersRead,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Find users based on group membership using the id of the group.",
			},
			"include_groups": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Fetch group memberships for each user",
			},
			"search": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Filter to find a user, each filter will be concatenated with an AND clause. Please be aware profile properties must match what is in Okta, which is likely camel case",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Property name to search for. This requires the search feature be on. Please see Okta documentation on their filter API for users. https://developer.okta.com/docs/api/resources/users#list-users-with-search",
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"comparison": {
							Type:             schema.TypeString,
							Optional:         true,
							Default:          "eq",
							ValidateDiagFunc: elemInSlice([]string{"eq", "lt", "gt", "sw"}),
						},
					},
				},
			},
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: buildSchema(userProfileDataSchema,
						map[string]*schema.Schema{
							"id": {
								Type:     schema.TypeString,
								Computed: true,
							},
						}),
				},
			},
		},
	}
}

func dataSourceUsersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var (
		users []*okta.User
		id    string
		err   error
	)

	client := getOktaClientFromMetadata(m)

	if groupId, ok := d.GetOk("group_id"); ok {
		id = groupId.(string)
		users, err = listGroupUsers(ctx, m, id)
	} else if _, ok := d.GetOk("search"); ok {
		params := &query.Params{Search: getSearchCriteria(d), Limit: defaultPaginationLimit, SortOrder: "0"}
		id = fmt.Sprintf("%d", crc32.ChecksumIEEE([]byte(params.String())))
		users, err = collectUsers(ctx, client, params)
	} else {
		return diag.Errorf("must specify either group_id or search attributes")
	}

	if err != nil {
		return diag.Errorf("failed to list users: %v", err)
	}
	d.SetId(id)
	shouldGetGroups := d.Get("include_groups").(bool)
	arr := make([]map[string]interface{}, len(users))
	for i, user := range users {
		rawMap := flattenUser(user)
		rawMap["id"] = user.Id
		if shouldGetGroups {
			groups, err := getGroupsForUser(ctx, user.Id, client)
			if err != nil {
				return diag.Errorf("failed to list users: %v", err)
			}
			rawMap["group_memberships"] = groups
		}
		arr[i] = rawMap
	}
	_ = d.Set("users", arr)
	return nil
}

func collectUsers(ctx context.Context, client *okta.Client, qp *query.Params) ([]*okta.User, error) {
	users, resp, err := client.User.ListUsers(ctx, qp)
	if err != nil {
		return nil, err
	}
	for resp.HasNextPage() {
		var nextUsers []*okta.User
		resp, err = resp.Next(ctx, &nextUsers)
		if err != nil {
			return nil, err
		}
		for i := range nextUsers {
			users = append(users, nextUsers[i])
		}
	}
	return users, nil
}
