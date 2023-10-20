package okta

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/crc32"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk/query"
)

func dataSourceGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGroupsRead,
		Schema: map[string]*schema.Schema{
			"q": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Searches the name property of groups for matching value",
			},
			"search": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Searches for groups with a supported filtering expression for all attributes except for '_embedded', '_links', and 'objectClass'",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Type of the group. When specified in the terraform resource, will act as a filter when searching for the groups",
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Group ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Group name.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Group type.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Group description.",
						},
						"custom_profile_attributes": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Raw JSON containing all custom profile attributes. Likely only useful on groups of type",
						},
					},
				},
			},
		},
		Description: "Get a list of groups from Okta.",
	}
}

func dataSourceGroupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	qp := &query.Params{Limit: 10000}
	groupType, ok := d.GetOk("type")
	if ok {
		qp.Filter = fmt.Sprintf("type eq \"%s\"", groupType.(string))
	}
	q, ok := d.GetOk("q")
	if ok {
		qp.Q = q.(string)
	}
	search, ok := d.GetOk("search")
	if ok {
		qp.Search = search.(string)
	}
	groups, err := listGroups(ctx, getOktaClientFromMetadata(m), qp)
	if err != nil {
		return diag.Errorf("failed to list groups: %v", err)
	}
	d.SetId(fmt.Sprintf("%d", crc32.ChecksumIEEE([]byte(qp.String()))))
	arr := make([]map[string]interface{}, len(groups))
	for i := range groups {
		customProfile, err := json.Marshal(groups[i].Profile.GroupProfileMap)
		if err != nil {
			return diag.Errorf("failed to read custom profile attributes from group: %s", groups[i].Profile.Name)
		}
		arr[i] = map[string]interface{}{
			"id":                        groups[i].Id,
			"name":                      groups[i].Profile.Name,
			"type":                      groups[i].Type,
			"description":               groups[i].Profile.Description,
			"custom_profile_attributes": string(customProfile),
		}
	}
	_ = d.Set("groups", arr)
	return nil
}
