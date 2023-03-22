package okta

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

func resourceNetworkZone() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkZoneCreate,
		ReadContext:   resourceNetworkZoneRead,
		UpdateContext: resourceNetworkZoneUpdate,
		DeleteContext: resourceNetworkZoneDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"dynamic_locations": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Array of locations ISO-3166-1(2). Format code: countryCode OR countryCode-regionCode",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"dynamic_proxy_type": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: elemInSlice([]string{"Any", "TorAnonymizer", "NotTorAnonymizer"}),
				Description:      "Type of proxy being controlled by this network zone",
			},
			"gateways": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Array of values in CIDR/range form depending on the way it's been declared (i.e. CIDR will contain /suffix). Please check API docs for examples",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"name": {
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
			"type": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: elemInSlice([]string{"IP", "DYNAMIC"}),
				Description:      "Type of the Network Zone - can either be IP or DYNAMIC only",
			},
			"usage": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Zone's purpose: POLICY or BLOCKLIST",
				ValidateDiagFunc: elemInSlice([]string{"POLICY", "BLOCKLIST"}),
				Default:          "POLICY",
			},
			"asns": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Format of each array value: a string representation of an ASN numeric value",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceNetworkZoneCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := validateNetworkZone(d)
	if err != nil {
		return diag.FromErr(err)
	}
	zone, _, err := getOktaClientFromMetadata(m).NetworkZone.CreateNetworkZone(ctx, buildNetworkZone(d))
	if err != nil {
		return diag.Errorf("failed to create network zone: %v", err)
	}
	d.SetId(zone.Id)
	return resourceNetworkZoneRead(ctx, d, m)
}

func resourceNetworkZoneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	zone, resp, err := getOktaClientFromMetadata(m).NetworkZone.GetNetworkZone(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get network zone: %v", err)
	}
	if zone == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("name", zone.Name)
	_ = d.Set("type", zone.Type)
	_ = d.Set("usage", zone.Usage)
	_ = d.Set("dynamic_proxy_type", zone.ProxyType)
	_ = d.Set("asns", convertStringSliceToSetNullable(zone.Asns))
	err = setNonPrimitives(d, map[string]interface{}{
		"gateways":          flattenAddresses(zone.Gateways),
		"proxies":           flattenAddresses(zone.Proxies),
		"dynamic_locations": flattenDynamicLocations(zone.Locations),
	})
	if err != nil {
		return diag.Errorf("failed to set network zone properties: %v", err)
	}
	return nil
}

func resourceNetworkZoneUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := validateNetworkZone(d)
	if err != nil {
		return diag.FromErr(err)
	}
	_, _, err = getOktaClientFromMetadata(m).NetworkZone.UpdateNetworkZone(ctx, d.Id(), buildNetworkZone(d))
	if err != nil {
		return diag.Errorf("failed to update network zone: %v", err)
	}
	return resourceNetworkZoneRead(ctx, d, m)
}

func resourceNetworkZoneDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resp, err := getOktaClientFromMetadata(m).NetworkZone.DeleteNetworkZone(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to delete network zone: %v", err)
	}
	return nil
}

func buildNetworkZone(d *schema.ResourceData) sdk.NetworkZone {
	var gatewaysList, proxiesList []*sdk.NetworkZoneAddress
	var locationsList []*sdk.NetworkZoneLocation
	zoneType := d.Get("type").(string)
	proxyType := d.Get("dynamic_proxy_type").(string)

	if zoneType == "IP" {
		if values, ok := d.GetOk("gateways"); ok {
			gatewaysList = buildAddressObjList(values.(*schema.Set))
		}
		if values, ok := d.GetOk("proxies"); ok {
			proxiesList = buildAddressObjList(values.(*schema.Set))
		}
	} else if values, ok := d.GetOk("dynamic_locations"); ok {
		for _, value := range values.(*schema.Set).List() {
			if strings.Contains(value.(string), "-") {
				locationsList = append(locationsList, &sdk.NetworkZoneLocation{Country: strings.Split(value.(string), "-")[0], Region: value.(string)})
			} else {
				locationsList = append(locationsList, &sdk.NetworkZoneLocation{Country: value.(string)})
			}
		}
	}

	return sdk.NetworkZone{
		Asns:      convertInterfaceToStringSetNullable(d.Get("asns")),
		Name:      d.Get("name").(string),
		Type:      zoneType,
		Gateways:  gatewaysList,
		Locations: locationsList,
		Proxies:   proxiesList,
		ProxyType: proxyType,
		Usage:     d.Get("usage").(string),
	}
}

func buildAddressObjList(values *schema.Set) []*sdk.NetworkZoneAddress {
	var addressType string
	var addressObjList []*sdk.NetworkZoneAddress
	for _, value := range values.List() {
		if strings.Contains(value.(string), "/") {
			addressType = "CIDR"
		} else {
			addressType = "RANGE"
		}
		addressObjList = append(addressObjList, &sdk.NetworkZoneAddress{Type: addressType, Value: value.(string)})
	}
	return addressObjList
}

func flattenAddresses(gateways []*sdk.NetworkZoneAddress) interface{} {
	if len(gateways) == 0 {
		return nil
	}
	arr := make([]interface{}, len(gateways))
	for i := range gateways {
		arr[i] = gateways[i].Value
	}
	return schema.NewSet(schema.HashString, arr)
}

func flattenDynamicLocations(locations []*sdk.NetworkZoneLocation) interface{} {
	if len(locations) == 0 {
		return nil
	}
	arr := make([]interface{}, len(locations))
	for i := range locations {
		if strings.Contains(locations[i].Region, "-") {
			arr[i] = locations[i].Region
		} else {
			arr[i] = locations[i].Country
		}
	}
	return schema.NewSet(schema.HashString, arr)
}

func validateNetworkZone(d *schema.ResourceData) error {
	proxies, ok := d.GetOk("proxies")
	if d.Get("usage").(string) != "POLICY" && ok && proxies.(*schema.Set).Len() != 0 {
		return fmt.Errorf(`zones with usage = "BLOCKLIST" cannot have trusted proxies`)
	}
	return nil
}
