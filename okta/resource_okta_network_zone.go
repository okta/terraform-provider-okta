package okta

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	v5okta "github.com/okta/okta-sdk-golang/v5/okta"
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
		Description: "Creates an Okta Network Zone. This resource allows you to create and configure an Okta Network Zone.",
		Schema: map[string]*schema.Schema{
			"dynamic_locations": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Array of locations ISO-3166-1(2) included. Format code: countryCode OR countryCode-regionCode. Use with type `DYNAMIC` or `DYNAMIC_V2`",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"dynamic_locations_exclude": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Array of locations ISO-3166-1(2) excluded. Format code: countryCode OR countryCode-regionCode. Use with type `DYNAMIC_V2`",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"dynamic_proxy_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Type of proxy being controlled by this dynamic network zone - can be one of `Any`, `TorAnonymizer` or `NotTorAnonymizer`. Use with type `DYNAMIC`",
			},
			"gateways": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Array of values in CIDR/range form depending on the way it's been declared (i.e. CIDR will contain /suffix). Please check API docs for examples. Use with type `IP`",
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
				Description: "Array of values in CIDR/range form depending on the way it's been declared (i.e. CIDR will contain /suffix). Please check API docs for examples. Can not be set if `usage` is set to `BLOCKLIST`. Use with type `IP`",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of the Network Zone - can be `IP`, `DYNAMIC` or `DYNAMIC_V2` only",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Network Status - can either be `ACTIVE` or `INACTIVE` only",
			},
			"usage": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Usage of the Network Zone - can be either `POLICY` or `BLOCKLIST`. By default, it is `POLICY`",
				Default:     "POLICY",
			},
			"asns": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of asns included. Format of each array value: a string representation of an ASN numeric value. Use with type `DYNAMIC` or `DYNAMIC_V2`",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"ip_service_categories_include": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of ip service included. Use with type `DYNAMIC_V2`",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"ip_service_categories_exclude": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of ip service excluded. Use with type `DYNAMIC_V2`",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceNetworkZoneCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := validateNetworkZone(d)
	if err != nil {
		return diag.FromErr(err)
	}
	payload, err := buildNetworkZone(d)
	if err != nil {
		return diag.FromErr(err)
	}
	zone, _, err := getOktaV5ClientFromMetadata(meta).NetworkZoneAPI.CreateNetworkZone(ctx).Zone(payload).Execute()
	if err != nil {
		return diag.Errorf("failed to create network zone: %v", err)
	}
	nzID, err := concreteNetworkzoneID(zone)
	if err != nil {
		return diag.Errorf("failed to create network zone: %v", err)
	}
	d.SetId(nzID)
	return resourceNetworkZoneRead(ctx, d, meta)
}

func resourceNetworkZoneRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zone, resp, err := getOktaV5ClientFromMetadata(meta).NetworkZoneAPI.GetNetworkZone(ctx, d.Id()).Execute()
	if err := v5suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get network zone: %v", err)
	}
	if zone == nil {
		d.SetId("")
		return nil
	}
	err = mapNetworkZoneToState(d, zone)
	if err != nil {
		return diag.Errorf("failed to set network zone properties: %v", err)
	}
	return nil
}

func resourceNetworkZoneUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := validateNetworkZone(d)
	if err != nil {
		return diag.FromErr(err)
	}
	payload, err := buildNetworkZone(d)
	if err != nil {
		return diag.FromErr(err)
	}
	_, _, err = getOktaV5ClientFromMetadata(meta).NetworkZoneAPI.ReplaceNetworkZone(ctx, d.Id()).Zone(payload).Execute()
	if err != nil {
		return diag.Errorf("failed to update network zone: %v", err)
	}
	return resourceNetworkZoneRead(ctx, d, meta)
}

func resourceNetworkZoneDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_, resp, err := getOktaV5ClientFromMetadata(meta).NetworkZoneAPI.DeactivateNetworkZone(ctx, d.Id()).Execute()
	if err := v5suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to deactivate network zone: %v", err)
	}
	resp, err = getOktaV5ClientFromMetadata(meta).NetworkZoneAPI.DeleteNetworkZone(ctx, d.Id()).Execute()
	if err := v5suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to delete network zone: %v", err)
	}
	return nil
}

func buildNetworkZone(d *schema.ResourceData) (v5okta.ListNetworkZones200ResponseInner, error) {
	var resp v5okta.ListNetworkZones200ResponseInner
	zoneType := d.Get("type").(string)
	switch zoneType {
	case "IP":
		ipnz := v5okta.IPNetworkZone{}
		ipnz.SetName(d.Get("name").(string))
		ipnz.SetType(zoneType)
		ipnz.SetUsage(d.Get("usage").(string))
		if values, ok := d.GetOk("gateways"); ok {
			ipnz.SetGateways(buildAddressObjList(values.(*schema.Set)))
		}
		if values, ok := d.GetOk("proxies"); ok {
			ipnz.SetProxies(buildAddressObjList(values.(*schema.Set)))
		}
		if status, ok := d.GetOk("status"); ok {
			ipnz.SetStatus(status.(string))
		}
		resp.IPNetworkZone = &ipnz
		return resp, nil
	case "DYNAMIC":
		dynz := v5okta.DynamicNetworkZone{}
		dynz.SetName(d.Get("name").(string))
		dynz.SetType(zoneType)
		dynz.SetUsage(d.Get("usage").(string))
		dynz.SetProxyType(d.Get("dynamic_proxy_type").(string))
		dynz.SetAsns(convertInterfaceToStringSetNullable(d.Get("asns")))
		var locationsList []v5okta.NetworkZoneLocation
		if values, ok := d.GetOk("dynamic_locations"); ok {
			locationsList = buildLocationList(values.(*schema.Set))
		}
		dynz.SetLocations(locationsList)
		if status, ok := d.GetOk("status"); ok {
			dynz.SetStatus(status.(string))
		}
		resp.DynamicNetworkZone = &dynz
		return resp, nil
	case "DYNAMIC_V2":
		dyv2nz := v5okta.EnhancedDynamicNetworkZone{}
		dyv2nz.SetName(d.Get("name").(string))
		dyv2nz.SetType(zoneType)
		dyv2nz.SetUsage(d.Get("usage").(string))
		asns := v5okta.EnhancedDynamicNetworkZoneAllOfAsns{Include: convertInterfaceToStringSetNullable(d.Get("asns"))}
		dyv2nz.SetAsns(asns)
		var locationsListInclude []v5okta.NetworkZoneLocation
		if values, ok := d.GetOk("dynamic_locations"); ok {
			locationsListInclude = buildLocationList(values.(*schema.Set))
		}
		var locationsListExclude []v5okta.NetworkZoneLocation
		if values, ok := d.GetOk("dynamic_locations_exclude"); ok {
			locationsListExclude = buildLocationList(values.(*schema.Set))
		}
		locations := v5okta.EnhancedDynamicNetworkZoneAllOfLocations{Include: locationsListInclude, Exclude: locationsListExclude}
		dyv2nz.SetLocations(locations)

		ipService := v5okta.EnhancedDynamicNetworkZoneAllOfIpServiceCategories{Include: convertInterfaceToStringSetNullable(d.Get("ip_service_categories_include")), Exclude: convertInterfaceToStringSetNullable(d.Get("ip_service_categories_exclude"))}
		dyv2nz.SetIpServiceCategories(ipService)
		if status, ok := d.GetOk("status"); ok {
			dyv2nz.SetStatus(status.(string))
		}
		resp.EnhancedDynamicNetworkZone = &dyv2nz
		return resp, nil
	default:
		return resp, fmt.Errorf("unknown network zone type %v", zoneType)
	}
}

func buildAddressObjList(values *schema.Set) []v5okta.NetworkZoneAddress {
	var addressType string
	var addressObjList []v5okta.NetworkZoneAddress
	for _, value := range values.List() {
		if strings.Contains(value.(string), "/") {
			addressType = "CIDR"
		} else {
			addressType = "RANGE"
		}
		obj := v5okta.NetworkZoneAddress{}
		obj.SetType(addressType)
		obj.SetValue(value.(string))
		addressObjList = append(addressObjList, obj)
	}
	return addressObjList
}

func buildLocationList(values *schema.Set) []v5okta.NetworkZoneLocation {
	var locationsList []v5okta.NetworkZoneLocation
	for _, value := range values.List() {
		if strings.Contains(value.(string), "-") {
			obj := v5okta.NetworkZoneLocation{}
			obj.SetCountry(strings.Split(value.(string), "-")[0])
			obj.SetRegion(value.(string))
			locationsList = append(locationsList, obj)
		} else {
			obj := v5okta.NetworkZoneLocation{}
			obj.SetCountry(value.(string))
			locationsList = append(locationsList, obj)
		}
	}
	return locationsList
}

func flattenAddresses(gateways []v5okta.NetworkZoneAddress) interface{} {
	if len(gateways) == 0 {
		return nil
	}
	arr := make([]interface{}, len(gateways))
	for i := range gateways {
		arr[i] = gateways[i].GetValue()
	}
	return schema.NewSet(schema.HashString, arr)
}

func flattenDynamicLocations(locations []v5okta.NetworkZoneLocation) interface{} {
	if len(locations) == 0 {
		return nil
	}
	arr := make([]interface{}, len(locations))
	for i := range locations {
		if strings.Contains(locations[i].GetRegion(), "-") {
			arr[i] = locations[i].GetRegion()
		} else {
			arr[i] = locations[i].GetCountry()
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

func concreteNetworkzoneID(src *v5okta.ListNetworkZones200ResponseInner) (id string, err error) {
	if src == nil {
		return "", errors.New("list network zone response is nil")
	}
	nz := src.GetActualInstance()
	if nz == nil {
		return "", errors.New("okta list network zone response does not contain a concrete type")
	}
	switch v := nz.(type) {
	case *v5okta.DynamicNetworkZone:
		id = v.GetId()
	case *v5okta.EnhancedDynamicNetworkZone:
		id = v.GetId()
	case *v5okta.IPNetworkZone:
		id = v.GetId()
	}
	if id == "" {
		err = fmt.Errorf("list network zone response does not contain a concrete type %T", src)
	}
	return
}

func mapNetworkZoneToState(d *schema.ResourceData, data *v5okta.ListNetworkZones200ResponseInner) error {
	if data == nil {
		return errors.New("list network zone response is nil")
	}
	nz := data.GetActualInstance()
	if nz == nil {
		return errors.New("okta list network zone response does not contain a concrete type")
	}
	var err error
	switch v := nz.(type) {
	case *v5okta.DynamicNetworkZone:
		_ = d.Set("name", v.GetName())
		_ = d.Set("type", v.GetType())
		_ = d.Set("status", v.GetStatus())
		_ = d.Set("usage", v.GetUsage())
		_ = d.Set("dynamic_proxy_type", v.GetProxyType())
		_ = d.Set("asns", convertStringSliceToSetNullable(v.GetAsns()))
		err = setNonPrimitives(d, map[string]interface{}{
			"dynamic_locations": flattenDynamicLocations(v.GetLocations()),
		})
	case *v5okta.EnhancedDynamicNetworkZone:
		_ = d.Set("name", v.GetName())
		_ = d.Set("type", v.GetType())
		_ = d.Set("status", v.GetStatus())
		_ = d.Set("usage", v.GetUsage())
		if asns, ok := v.GetAsnsOk(); ok {
			err = setNonPrimitives(d, map[string]interface{}{
				"asns": convertStringSliceToSetNullable(asns.GetInclude()),
			})
		}
		if location, ok := v.GetLocationsOk(); ok {
			err = setNonPrimitives(d, map[string]interface{}{
				"dynamic_locations":         flattenDynamicLocations(location.GetInclude()),
				"dynamic_locations_exclude": flattenDynamicLocations(location.GetExclude()),
			})
		}
		if ips, ok := v.GetIpServiceCategoriesOk(); ok {
			err = setNonPrimitives(d, map[string]interface{}{
				"ip_service_categories_include": convertStringSliceToSetNullable(ips.GetInclude()),
				"ip_service_categories_exclude": convertStringSliceToSetNullable(ips.GetExclude()),
			})
		}
	case *v5okta.IPNetworkZone:
		_ = d.Set("name", v.GetName())
		_ = d.Set("type", v.GetType())
		_ = d.Set("status", v.GetStatus())
		_ = d.Set("usage", v.GetUsage())
		err = setNonPrimitives(d, map[string]interface{}{
			"gateways": flattenAddresses(v.GetGateways()),
			"proxies":  flattenAddresses(v.GetProxies()),
		})
	}
	return err
}
