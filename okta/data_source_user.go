package okta

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/validation"

	"github.com/okta/okta-sdk-golang/okta/query"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceUserRead,

		Schema: map[string]*schema.Schema{
			"search": &schema.Schema{
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Filter to find a user, each filter will be concatenated with an AND clause. Please be aware profile properties must match what is in Okta, which is likely camel case",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Property name to search for. This requires the search feature be on. Please see Okta documentation on their filter API for users. https://developer.okta.com/docs/api/resources/users#list-users-with-search",
						},
						"value": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"comparison": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "eq",
							ValidateFunc: validation.StringInSlice([]string{"eq", "lt", "gt", "sw"}, true),
						},
					},
				},
			},
			"admin_roles": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"city": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"cost_center": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"country_code": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"custom_profile_attributes": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
			"department": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"division": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"employee_number": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"first_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_memberships": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"honorific_prefix": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"honorific_suffix": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"locale": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"login": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"manager": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"manager_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"middle_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"mobile_phone": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"nick_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"organization": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"postal_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"preferred_language": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"primary_phone": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"profile_url": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"second_email": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"street_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"timezone": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"title": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"zip_code": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceUserRead(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	users, _, err := client.User.ListUsers(&query.Params{Search: getSearchCriteria(d), Limit: 1})

	if err != nil {
		return fmt.Errorf("Error Getting User from Okta: %v", err)
	} else if len(users) < 1 {
		return errors.New("failed to locate user with provided parameters")
	}

	d.Set("status", mapStatus(users[0].Status))
	d.SetId(users[0].Id)

	if err = setUserProfileAttributes(d, users[0]); err != nil {
		return err
	}

	if err = setAdminRoles(d, client); err != nil {
		return err
	}

	return nil
}

func getSearchCriteria(d *schema.ResourceData) string {
	rawFilters := d.Get("search").(*schema.Set)
	filterList := make([]string, rawFilters.Len())
	for i, f := range rawFilters.List() {
		fmap := f.(map[string]interface{})
		filterList[i] = fmt.Sprintf(`%s %s "%s"`, fmap["name"], fmap["comparison"], fmap["value"])
	}

	return strings.Join(filterList, " and ")
}
