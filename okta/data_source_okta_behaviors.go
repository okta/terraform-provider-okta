package okta

import (
	"context"
	"fmt"
	"hash/crc32"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

func dataSourceBehaviors() *schema.Resource {
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
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"settings": {
							Type:     schema.TypeMap,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBehaviorsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	qp := &query.Params{Limit: defaultPaginationLimit}
	q, ok := d.GetOk("q")
	if ok {
		qp.Q = q.(string)
	}
	behaviors, _, err := getSupplementFromMetadata(m).ListBehaviors(ctx, qp)
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
