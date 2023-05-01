package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

func dataSourceNetworkZone() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkZoneRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
			"dynamic_locations": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Array of locations ISO-3166-1(2). Format code: countryCode OR countryCode-regionCode",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"dynamic_proxy_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of proxy being controlled by this network zone",
			},
			"gateways": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Array of values in CIDR/range form depending on the way it's been declared (i.e. CIDR will contain /suffix). Please check API docs for examples",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"proxies": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Array of values in CIDR/range form depending on the way it's been declared (i.e. CIDR will contain /suffix). Please check API docs for examples",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the Network Zone - can either be IP or DYNAMIC only",
			},
			"usage": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Zone's purpose: POLICY or BLOCKLIST",
			},
			"asns": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Format of each array value: a string representation of an ASN numeric value",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceNetworkZoneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Get("id").(string)
	name := d.Get("name").(string)
	if id == "" && name == "" {
		return diag.Errorf("config must provide either 'id' or 'name' to retrieve the network zone")
	}
	var (
		err  error
		zone *sdk.NetworkZone
	)
	if id != "" {
		zone, _, err = getOktaClientFromMetadata(m).NetworkZone.GetNetworkZone(ctx, id)
	} else {
		zone, err = findNetworkZoneByName(ctx, m, name)
	}
	if err != nil {
		return diag.Errorf("failed to find network zone: %v", err)
	}
	d.SetId(zone.Id)
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

func findNetworkZoneByName(ctx context.Context, m interface{}, name string) (*sdk.NetworkZone, error) {
	client := getOktaClientFromMetadata(m)
	zones, resp, err := client.NetworkZone.ListNetworkZones(ctx, nil)
	if err != nil {
		return nil, err
	}
	for i := range zones {
		if zones[i].Name == name {
			return zones[i], nil
		}
	}
	for {
		var moreZones []*sdk.NetworkZone
		if resp.HasNextPage() {
			resp, err = resp.Next(ctx, &moreZones)
			if err != nil {
				return nil, err
			}
			for i := range moreZones {
				if moreZones[i].Name == name {
					return moreZones[i], nil
				}
			}
		} else {
			break
		}
	}
	return nil, fmt.Errorf("network zone with name '%s' does not exist", name)
}
