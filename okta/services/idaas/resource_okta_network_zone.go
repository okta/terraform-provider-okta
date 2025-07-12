package idaas

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	v5okta "github.com/okta/okta-sdk-golang/v5/okta"
	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

const defaultEnhancedDynamicZone = "DefaultEnhancedDynamicZone"

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
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Network Status - can either be `ACTIVE` or `INACTIVE` only",
				ValidateFunc: validation.StringInSlice([]string{"ACTIVE", "INACTIVE"}, false),
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
			"use_as_exempt_list": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Indicates that this network zone is used as an exempt list. Only applicable to IP zones. This parameter is required when updating the DefaultExemptIpZone to allow IPs through the blocklist.",
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
	if err := utils.SuppressErrorOn404_V5(resp, err); err != nil {
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

	err = updateDefaultEnhancedDynamicZone(d, meta)
	if err != nil {
		return diag.Errorf("failed to update network zone: %v", err)
	}

	payload, err := buildNetworkZone(d)
	if err != nil {
		return diag.FromErr(err)
	}
	// Check if this is an exempt zone update
	if useAsExemptList, ok := d.GetOk("use_as_exempt_list"); ok && useAsExemptList.(bool) {
		// For exempt zones, we need to add the useAsExemptList=true query parameter
		// Since the SDK doesn't directly support this, we'll construct the request manually
		err = updateNetworkZoneWithExemptList(ctx, meta, d.Id(), payload)
		if err != nil {
			return diag.Errorf("failed to update exempt network zone: %v", err)
		}
	} else {
		_, _, err = getOktaV5ClientFromMetadata(meta).NetworkZoneAPI.ReplaceNetworkZone(ctx, d.Id()).Zone(payload).Execute()
		if err != nil {
			return diag.Errorf("failed to update network zone: %v", err)
		}
	}
	return resourceNetworkZoneRead(ctx, d, meta)
}

func resourceNetworkZoneDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_, resp, err := getOktaV5ClientFromMetadata(meta).NetworkZoneAPI.DeactivateNetworkZone(ctx, d.Id()).Execute()
	if err := utils.SuppressErrorOn404_V5(resp, err); err != nil {
		return diag.Errorf("failed to deactivate network zone: %v", err)
	}
	resp, err = getOktaV5ClientFromMetadata(meta).NetworkZoneAPI.DeleteNetworkZone(ctx, d.Id()).Execute()
	if err := utils.SuppressErrorOn404_V5(resp, err); err != nil {
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
		dynz.SetAsns(utils.ConvertInterfaceToStringSetNullable(d.Get("asns")))
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
		asns := v5okta.EnhancedDynamicNetworkZoneAllOfAsns{Include: utils.ConvertInterfaceToStringSetNullable(d.Get("asns"))}
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

		ipService := v5okta.EnhancedDynamicNetworkZoneAllOfIpServiceCategories{Include: utils.ConvertInterfaceToStringSetNullable(d.Get("ip_service_categories_include")), Exclude: utils.ConvertInterfaceToStringSetNullable(d.Get("ip_service_categories_exclude"))}
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

	// Validate that use_as_exempt_list is only used with IP zones
	if useAsExemptList, ok := d.GetOk("use_as_exempt_list"); ok && useAsExemptList.(bool) {
		if d.Get("type").(string) != "IP" {
			return fmt.Errorf(`use_as_exempt_list can only be set to true for IP zones`)
		}
	}

	return nil
}

func updateDefaultEnhancedDynamicZone(d *schema.ResourceData, meta interface{}) error {
	status, ok := d.GetOk("status")
	if d.Get("name").(string) == defaultEnhancedDynamicZone && ok {
		switch status.(string) {
		case "ACTIVE":
			_, _, err := getOktaV5ClientFromMetadata(meta).NetworkZoneAPI.ActivateNetworkZone(context.Background(), d.Id()).Execute()
			return err
		case "INACTIVE":
			_, _, err := getOktaV5ClientFromMetadata(meta).NetworkZoneAPI.DeactivateNetworkZone(context.Background(), d.Id()).Execute()
			return err
		}
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
		_ = d.Set("asns", utils.ConvertStringSliceToSetNullable(v.GetAsns()))
		err = utils.SetNonPrimitives(d, map[string]interface{}{
			"dynamic_locations": flattenDynamicLocations(v.GetLocations()),
		})
	case *v5okta.EnhancedDynamicNetworkZone:
		_ = d.Set("name", v.GetName())
		_ = d.Set("type", v.GetType())
		_ = d.Set("status", v.GetStatus())
		_ = d.Set("usage", v.GetUsage())
		if asns, ok := v.GetAsnsOk(); ok {
			err = utils.SetNonPrimitives(d, map[string]interface{}{
				"asns": utils.ConvertStringSliceToSetNullable(asns.GetInclude()),
			})
		}
		if location, ok := v.GetLocationsOk(); ok {
			err = utils.SetNonPrimitives(d, map[string]interface{}{
				"dynamic_locations":         flattenDynamicLocations(location.GetInclude()),
				"dynamic_locations_exclude": flattenDynamicLocations(location.GetExclude()),
			})
		}
		if ips, ok := v.GetIpServiceCategoriesOk(); ok {
			err = utils.SetNonPrimitives(d, map[string]interface{}{
				"ip_service_categories_include": utils.ConvertStringSliceToSetNullable(ips.GetInclude()),
				"ip_service_categories_exclude": utils.ConvertStringSliceToSetNullable(ips.GetExclude()),
			})
		}
	case *v5okta.IPNetworkZone:
		_ = d.Set("name", v.GetName())
		_ = d.Set("type", v.GetType())
		_ = d.Set("status", v.GetStatus())
		_ = d.Set("usage", v.GetUsage())
		err = utils.SetNonPrimitives(d, map[string]interface{}{
			"gateways": flattenAddresses(v.GetGateways()),
			"proxies":  flattenAddresses(v.GetProxies()),
		})
	}
	return err
}

// updateNetworkZoneWithExemptList makes a custom HTTP request to update a network zone
// with the useAsExemptList field included in the JSON body
func updateNetworkZoneWithExemptList(ctx context.Context, meta interface{}, zoneID string, payload v5okta.ListNetworkZones200ResponseInner) error {
	// Get the configuration from meta
	cfg := meta.(*config.Config)

	// Build the URL (no query parameter needed)
	// The cfg.Domain is just the base domain like "okta.com",
	// but we need the full org URL like "https://trial-7001215.okta.com"
	baseURL := fmt.Sprintf("https://%s.%s", cfg.OrgName, strings.TrimSuffix(cfg.Domain, "/"))
	endpoint := fmt.Sprintf("/api/v1/zones/%s", zoneID)
	fullURL := baseURL + endpoint

	// Convert SDK payload to map for manipulation
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	// Parse JSON into a map so we can add the useAsExemptList field
	var payloadMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &payloadMap); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %v", err)
	}

	// Add the useAsExemptList field to the JSON payload
	payloadMap["useAsExemptList"] = true

	// Re-marshal with the added field
	finalJsonData, err := json.Marshal(payloadMap)
	if err != nil {
		return fmt.Errorf("failed to marshal final payload: %v", err)
	}

	// Create the HTTP request
	req, err := http.NewRequestWithContext(ctx, "PUT", fullURL, bytes.NewBuffer(finalJsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Set required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "SSWS "+cfg.ApiToken)
	req.Header.Set("User-Agent", "terraform-provider-okta")

	// Make the request
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
