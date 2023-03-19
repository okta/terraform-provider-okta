package okta

import (
	"context"
	"fmt"
	"hash/crc32"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"github.com/okta/terraform-provider-okta/sdk"
)

func dataSourceDevices() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDevicesRead,
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Searches for devices owned by the specified user_id",
			},
			"devices": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"profile": {
							Type:     schema.TypeMap,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_display_name": {
							Type:     schema.TypeMap,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"resource_alternate_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDevicesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	qp := &query.Params{Limit: defaultPaginationLimit}
	userID, ok := d.GetOk("user_id")
	var devices []*sdk.Device
	var err error
	if ok {
		qp.Expand = "users"

		respDevices, _, err := getSupplementFromMetadata(m).ListDevices(ctx, qp)
		if err != nil {
			return diag.Errorf("failed to list devices: %v", err)
		}

		devices = searchUserDevices(ctx, respDevices, userID.(string), m)
		if err != nil {
			return diag.Errorf("failed to find devices for specified user: %v", err)
		}
	} else {
		devices, _, err = getSupplementFromMetadata(m).ListDevices(ctx, qp)
		if err != nil {
			return diag.Errorf("failed to list devices: %v", err)
		}
	}
	d.SetId(fmt.Sprintf("%d", crc32.ChecksumIEEE([]byte(qp.String()))))
	arr := make([]map[string]interface{}, len(devices))
	for i := range devices {
		arr[i] = map[string]interface{}{
			"id":                    devices[i].ID,
			"status":                devices[i].Status,
			"resource_type":         devices[i].ResourceType,
			"resource_alternate_id": devices[i].ResourceAlternateID,
			"resource_id":           devices[i].ResourceID,
		}
		profile := make(map[string]string)
		for k, v := range devices[i].Profile {
			profile[k] = fmt.Sprint(v)
		}
		arr[i]["profile"] = profile
		resourceDisplayName := make(map[string]string)
		for k, v := range devices[i].ResourceDisplayName {
			resourceDisplayName[k] = fmt.Sprint(v)
		}
		arr[i]["resource_display_name"] = resourceDisplayName
	}
	err = d.Set("devices", arr)
	return diag.FromErr(err)
}

func searchUserDevices(ctx context.Context, devices []*sdk.Device, userID string, m interface{}) []*sdk.Device {
	var userDevices []*sdk.Device
	for _, respDevice := range devices {
		for _, respUser := range respDevice.Embedded["users"] {
			u := respUser["user"]
			for key, value := range u.(map[string]interface{}) {
				if key == "id" {
					if userID == value {
						userDevices = append(userDevices, respDevice)
					}
				}
			}
		}
	}

	return userDevices
}
