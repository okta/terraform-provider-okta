package okta

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/okta/okta-sdk-golang/okta/query"
)

func dataSourceGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGroupsRead,

		Schema: map[string]*schema.Schema{
			"q": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Searches the name property of groups for matching value",
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"description": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceGroupsRead(d *schema.ResourceData, m interface{}) error {
	client := getOktaClientFromMetadata(m)
	q := d.Get("q").(string)

	groups, _, err := client.Group.ListGroups(&query.Params{Q: q})

	if err != nil {
		return fmt.Errorf("failed to query for groups: %v", err)
	}

	d.SetId(fmt.Sprintf("%d", hashcode.String(q)))
	arr := make([]map[string]interface{}, len(groups))

	for i, group := range groups {
		arr[i] = map[string]interface{}{
			"name":        group.Profile.Name,
			"description": group.Profile.Description,
		}
	}

	return d.Set("groups", arr)
}
