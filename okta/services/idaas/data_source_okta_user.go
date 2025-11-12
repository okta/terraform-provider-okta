package idaas

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk"
	"github.com/okta/terraform-provider-okta/sdk/query"
)

var userSearchSchemaDescription = "Filter to find " +
	"user/users. Each filter will be concatenated with " +
	"the compound search operator. Please be aware profile properties " +
	"must match what is in Okta, which is likely camel case. " +
	"Expression is a free form expression filter " +
	"https://developer.okta.com/docs/reference/core-okta-api/#filter . " +
	"The set name/value/comparison properties will be ignored if expression is present"

var userSearchSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Property name to search for. This requires the search feature be on. Please see Okta documentation on their filter API for users. https://developer.okta.com/docs/api/resources/users#list-users-with-search",
	},
	"value": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"comparison": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "eq",
	},
	"expression": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "A raw search expression string. This requires the search feature be on. Please see Okta documentation on their filter API for users. https://developer.okta.com/docs/api/resources/users#list-users-with-search",
	},
}

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserRead,
		Schema: buildUserDataSourceSchema(map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Description: "Retrieve a single user based on their id",
				Optional:    true,
			},
			"search": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: userSearchSchemaDescription,
				Elem: &schema.Resource{
					Schema: userSearchSchema,
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
				Description: "Force delay of the user read by N seconds. Useful when eventual consistency of user information needs to be allowed for.",
			},
			"skip_groups": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Do not populate user groups information (prevents additional API call)",
			},
			"skip_roles": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Do not populate user roles information (prevents additional API call)",
			},
			"realm_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Realm ID associated with the user.",
			},
		}),
		Description: "Get a single users from Okta.",
	}
}

func dataSourceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if n, ok := d.GetOk("delay_read_seconds"); ok {
		delay, err := strconv.Atoi(n.(string))
		if err == nil {
			logger(meta).Info("delaying user read by ", delay, " seconds")
			meta.(*config.Config).TimeOperations.Sleep(time.Duration(delay) * time.Second)
		} else {
			logger(meta).Warn("user read delay value ", n, " is not an integer")
		}
	}

	client := getOktaClientFromMetadata(meta)
	var user *sdk.User
	var err error
	userID, ok := d.GetOk("user_id")
	_, searchCriteriaOk := d.GetOk("search")
	if ok {
		logger(meta).Info("reading user by ID", "id", userID.(string))
		user, _, err = client.User.GetUser(ctx, userID.(string))
		if err != nil {
			return diag.Errorf("failed to get user: %v", err)
		}
	} else if searchCriteriaOk {
		var users []*sdk.User
		sc := getSearchCriteria(d)
		logger(meta).Info("reading user using search", "search", sc)
		users, _, err = client.User.ListUsers(ctx, &query.Params{Search: sc, Limit: 1})
		if err != nil {
			return diag.Errorf("failed to list users: %v", err)
		} else if len(users) < 1 {
			return diag.Errorf("no users found using search criteria: %+v", sc)
		}
		user = users[0]
	}
	d.SetId(user.Id)
	rawMap := flattenUser(user, []string{})
	err = utils.SetNonPrimitives(d, rawMap)
	if err != nil {
		return diag.Errorf("failed to set user's properties: %v", err)
	}
	if val := d.Get("skip_roles"); val != nil {
		if skip, ok := val.(bool); ok && !skip {
			err = setAdminRoles(ctx, d, meta)
			if err != nil {
				return diag.Errorf("failed to set user's admin roles: %v", err)
			}
			err = setRoles(ctx, d, meta)
			if err != nil {
				return diag.Errorf("failed to set user's roles: %v", err)
			}
		}
	}

	if val := d.Get("skip_groups"); val != nil {
		if skip, ok := val.(bool); ok && !skip {
			err = setAllGroups(ctx, d, client)
			if err != nil {
				return diag.Errorf("failed to set user's groups: %v", err)
			}
		}
	}

	return nil
}

func getSearchCriteria(d *schema.ResourceData) string {
	rawFilters := d.Get("search").(*schema.Set)
	filterList := make([]string, rawFilters.Len())
	for i, f := range rawFilters.List() {
		fmap := f.(map[string]interface{})
		rawExpression := fmap["expression"]
		if rawExpression != "" && rawExpression != nil {
			filterList[i] = fmt.Sprintf(`%s`, fmap["expression"])
			continue
		}

		// Need to set up the filter clause to allow comparisons that do not
		// accept a right hand argument and those that do.
		// profile.email pr
		filterList[i] = fmt.Sprintf(`%s %s`, fmap["name"], fmap["comparison"])
		if fmap["value"] != "" {
			// profile.email eq "example@example.com"
			filterList[i] = fmt.Sprintf(`%s "%s"`, filterList[i], fmap["value"])
		}
	}

	operator := " and "
	cso := d.Get("compound_search_operator")
	if cso != nil && cso.(string) != "" {
		operator = fmt.Sprintf(" %s ", cso.(string))
	}
	return strings.Join(filterList, operator)
}
