package okta

import (
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-okta/sdk"
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
			"dynamic_locations": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Array of locations ISO-3166-1(2). Format code: countryCode OR countryCode-regionCode",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"gateways": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Array of values in CIDR/range form depending on the way it's been declared (i.e. CIDR will contain /suffix). Please check API docs for examples",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Network Zone Resource",
			},
			"proxies": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Array of values in CIDR/range form depending on the way it's been declared (i.e. CIDR will contain /suffix). Please check API docs for examples",
				Elem:        &schema.Schema{Type: schema.TypeString},
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
	zone, resp, err := getSupplementFromMetadata(m).GetNetworkZone(d.Id())

	if is404(resp.StatusCode) {
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	d.Set("name", zone.Name)
	d.Set("type", zone.Type)

	return setNonPrimitives(d, map[string]interface{}{
		// TODO
		// "gateways" 		: flattenHookGateways(),
		// "proxies" 		: flattenProxies(hook.Channel),
		// "dynamic_locations" 	: flattenDynamicLocations(d, hook.Channel),
	})
}

func resourceNetworkZoneUpdate(d *schema.ResourceData, m interface{}) error {
	client := getSupplementFromMetadata(m)
	networkZone := buildNetworkZone(d, m)
	_, _, err := client.UpdateNetworkZone(d.Id(), *networkZone, nil)

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

func buildNetworkZone(d *schema.ResourceData, m interface{}) *sdk.NetworkZone {
	gatewaysList := []*sdk.AddressObj{}
	proxiesList := []*sdk.AddressObj{}
	locationsList := []*sdk.Location{}
	zoneType := d.Get("type").(string)

	if strings.TrimRight(zoneType, "\n") == "IP" {
		if values, ok := d.GetOk("gateways"); ok {
			gatewaysList = buildAddressObjList(values.(*schema.Set))
		}
		if values, ok := d.GetOk("proxies"); ok {
			proxiesList = buildAddressObjList(values.(*schema.Set))
		}
	} else {
		if values, ok := d.GetOk("dynamic_locations"); ok {
			for _, value := range values.(*schema.Set).List() {
				if strings.Contains(value.(string), "-") {
					locationsList = append(locationsList, &sdk.Location{Country: strings.Split(value.(string), "-")[0], Region: value.(string)})
				} else {
					locationsList = append(locationsList, &sdk.Location{Country: value.(string)})
				}
			}
		}
	}

	return &sdk.NetworkZone{
		Name:      d.Get("name").(string),
		Type:      zoneType,
		Gateways:  gatewaysList,
		Locations: locationsList,
		Proxies:   proxiesList,
	}
}

func buildAddressObjList(values *schema.Set) []*sdk.AddressObj {
	var addressType string
	addressObjList := []*sdk.AddressObj{}

	for _, value := range values.List() {
		if strings.Contains(value.(string), "/") {
			addressType = "CIDR"
		} else {
			addressType = "RANGE"
		}
		addressObjList = append(addressObjList, &sdk.AddressObj{Type: addressType, Value: value.(string)})
	}
	return addressObjList
}
