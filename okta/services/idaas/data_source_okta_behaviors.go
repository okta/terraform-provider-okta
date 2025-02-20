package idaas

import (
	"context"
	"fmt"
	"hash/crc32"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
	"github.com/okta/terraform-provider-okta/sdk/query"
)

func DataSourceBehaviors() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBehaviorsRead,
		Schema: map[string]*schema.Schema{
			"q": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Searches the name property of behaviors for matching value",
			},
			"behaviors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Behavior ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Behavior name.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Behavior status.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Behavior type.",
						},
						"settings": {
							Type:        schema.TypeMap,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "Map of behavior settings.",
						},
					},
				},
			},
		},
		Description: "Get a behaviors by search criteria.",
	}
}

func dataSourceBehaviorsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	qp := &query.Params{Limit: utils.DefaultPaginationLimit}
	q, ok := d.GetOk("q")
	if ok {
		qp.Q = q.(string)
	}
	behaviors, _, err := GetAPISupplementFromMetadata(meta).ListBehaviors(ctx, qp)
	if err != nil {
		return diag.Errorf("failed to list behaviors: %v", err)
	}
	d.SetId(fmt.Sprintf("%d", crc32.ChecksumIEEE([]byte(qp.String()))))
	arr := make([]map[string]interface{}, len(behaviors))
	for i := range behaviors {
		arr[i] = map[string]interface{}{
			"id":     behaviors[i].ID,
			"name":   behaviors[i].Name,
			"type":   behaviors[i].Type,
			"status": behaviors[i].Status,
		}
		settings := make(map[string]string)
		for k, v := range behaviors[i].Settings {
			settings[k] = fmt.Sprint(v)
		}
		arr[i]["settings"] = settings
	}
	err = d.Set("behaviors", arr)
	return diag.FromErr(err)
}
