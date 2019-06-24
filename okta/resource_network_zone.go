package okta

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceNetworkZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkZoneCreate,
		Read:   resourceNetworkZoneRead,
		Update: resourceNetworkZoneUpdate,
		Delete: resourceNetworkZoneDelete,
		// Exists: resourceNetworkZoneExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"gateway_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "IP addresses (range or CIDR form) of this zone",
				ValidateFunc: validation.StringInSlice([]string{"CIDR", "RANGE"}, false),
			},
			"gateway_values": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Array of values in CIDR/range form depending on the type specified",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Network Zone Resource",
			},
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"IP", "DYNAMIC"}, false),
				Description:  "Type of the Network Zone - can either be IP or DYNAMIC only",
			},
		},
	}
}

func resourceNetworkZoneCreate(d *schema.ResourceData, m interface{}) error {
	client := getSupplementFromMetadata(m)
	networkZone := buildNetworkZone(d, m)
	networkZone, _, err := client.CreateNetworkZone(*networkZone, nil)
	if err != nil {
		return err
	}

	d.SetId(networkZone.ID)
	return resourceNetworkZoneRead(d, m)
}

func resourceNetworkZoneRead(d *schema.ResourceData, m interface{}) error {
	zone, _, err := getSupplementFromMetadata(m).GetNetworkZone(d.Id())
	if err != nil {
		return err
	}
	d.Set("name", zone.Name)
	d.Set("type", zone.Type)

	return setNonPrimitives(d, map[string]interface{}{
		// TODO
		// "channel": flattenHookChannel(hook.Channel),
		// "headers": flattenHeaders(hook.Channel),
		// "auth":    flattenAuth(d, hook.Channel),
	})
}

func resourceNetworkZoneUpdate(d *schema.ResourceData, m interface{}) error {
	client := getSupplementFromMetadata(m)
	networkZone := buildNetworkZone(d, m)
	networkZone, _, err := client.UpdateNetworkZone(d.Id(), *networkZone, nil)

	if err != nil {
		return err
	}

	return resourceNetworkZoneRead(d, m)
}

func resourceNetworkZoneDelete(d *schema.ResourceData, m interface{}) error {
	client := getSupplementFromMetadata(m)

	res, err := client.DeleteNetworkZone(d.Id())
	if err != nil {
		return responseErr(res, err)
	}

	return err
}

func buildNetworkZone(d *schema.ResourceData, m interface{}) *NetworkZone {
	gatewayList := []*Gateways{}
	if values, ok := d.GetOk("gateway_values"); ok {
		for _, value := range values.(*schema.Set).List() {
			gatewayList = append(gatewayList, &Gateways{Type: d.Get("gateway_type").(string), Value: value.(string)})
		}
	}

	return &NetworkZone{
		Name:     d.Get("name").(string),
		Type:     d.Get("type").(string),
		Gateways: gatewayList,
	}
}

// func buildGateways(d *schema.ResourceData, m interface{}) *Gateways {
// 	// if _, ok := d.GetOk("gateway_type"); !ok {
// 	// 	return nil
// 	// }

// 	return &Gateways{
// 		Type:  getStringValue(d, "gateway_type"),
// 		Value: getStringValue(d, "gateway_values"),
// 	}
// }
