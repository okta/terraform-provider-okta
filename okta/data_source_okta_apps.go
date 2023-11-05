package okta

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/crc32"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceApps() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppsRead,
		Schema: map[string]*schema.Schema{
			"active_only": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Search only active applications.",
			},
			"label": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"label_prefix"},
				Description:   "Searches for applications whose label or name property matches this value exactly",
			},
			"label_prefix": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"label"},
				Description:   "Searches for applications whose label or name property begins with this value",
			},
			"apps": {
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
						"label": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"links": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAppsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	filters, err := getAppsFilters(d)
	if err != nil {
		return diag.Errorf("invalid apps filters: %v", err)
	}

	appsResponse, err := listApps(ctx, getOktaClientFromMetadata(m), filters, 200)
	if err != nil {
		return diag.Errorf("failed to list apps: %v", err)
	}
	d.SetId(fmt.Sprintf("%d", crc32.ChecksumIEEE([]byte(filters.String()))))

	if len(appsResponse) > 0 {
		appsArr := []map[string]interface{}{}
		for _, app := range appsResponse {
			// Okta API for list apps uses a starts with query on label and name.
			// This can yield unexpected results if an exact match is desired.
			// If requested, drop apps that don't match the exact value.
			if filters.Label != "" && app.Label != filters.Label && app.Name != filters.Label {
				continue
			}
			links, _ := json.Marshal(app.Links)
			app := map[string]interface{}{
				"id":     app.Id,
				"name":   app.Name,
				"label":  app.Label,
				"status": app.Status,
				"links":  string(links),
			}
			appsArr = append(appsArr, app)
		}
		_ = d.Set("apps", appsArr)
	} else {
		_ = d.Set("apps", make([]map[string]interface{}, 0))
	}

	return nil
}
