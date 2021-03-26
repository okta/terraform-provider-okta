package okta

import (
	"context"
	"fmt"
	"hash/crc32"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
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
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Type of the group. When specified in the terraform resource, will act as a filter when searching for the groups",
				ValidateDiagFunc: stringInSlice([]string{"OKTA_GROUP", "APP_GROUP", "BUILT_IN"}),
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceGroupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	qp := &query.Params{Limit: defaultPaginationLimit}
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
		arr[i] = map[string]interface{}{
			"id":          groups[i].Id,
			"name":        groups[i].Profile.Name,
			"type":        groups[i].Type,
			"description": groups[i].Profile.Description,
		}
	}
	_ = d.Set("groups", arr)
	return nil
}
