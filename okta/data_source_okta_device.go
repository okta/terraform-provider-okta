package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

func dataSourceDevice() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDeviceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"profile": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_display_name": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
	}
}

func dataSourceDeviceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var device *sdk.Device
	deviceID, ok := d.GetOk("id")
	if !ok {
        return diag.Errorf("user_id not specified") 
    }

    respDevice, _, err := getSupplementFromMetadata(m).GetDevice(ctx, deviceID.(string))
    if err != nil {
        return diag.Errorf("failed to get device by ID: %v", err)
    }
    device = respDevice

	if device == nil {
		return nil
	}

	d.SetId(device.ID)

	_ = d.Set("status", device.Status)

	profile := make(map[string]string)
	for k, v := range device.Profile {
		profile[k] = fmt.Sprint(v)
	}
	_ = d.Set("profile", profile)

	_ = d.Set("resource_type", device.ResourceType)

	resourceDisplayName := make(map[string]string)
	for k, v := range device.ResourceDisplayName {
		resourceDisplayName[k] = fmt.Sprint(v)
	}
	_ = d.Set("resource_display_name", device.ResourceDisplayName)

	_ = d.Set("resource_alternate_id", device.ResourceAlternateID)
	_ = d.Set("resource_id", device.ResourceID)

	return nil
}


