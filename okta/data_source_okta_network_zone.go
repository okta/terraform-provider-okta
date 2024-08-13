package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	v5okta "github.com/okta/okta-sdk-golang/v5/okta"
)

func dataSourceNetworkZone() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkZoneRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"name"},
				Description:   "ID of the network zone to retrieve, conflicts with `name`.",
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id"},
				Description:   "Name of the network zone to retrieve, conflicts with `id`.",
			},
			"dynamic_locations": {
				Type:        schema.TypeSet,
				Computed:    true,
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
				Computed:    true,
				Description: "Type of proxy being controlled by this dynamic network zone - can be one of `Any`, `TorAnonymizer` or `NotTorAnonymizer`. Use with type `DYNAMIC`",
			},
			"gateways": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Array of values in CIDR/range form depending on the way it's been declared (i.e. CIDR will contain /suffix). Please check API docs for examples. Use with type `IP`",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"proxies": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Array of values in CIDR/range form depending on the way it's been declared (i.e. CIDR will contain /suffix). Please check API docs for examples. Can not be set if `usage` is set to `BLOCKLIST`. Use with type `IP`",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the Network Zone - can be `IP`, `DYNAMIC` or `DYNAMIC_V2` only",
			},
			"usage": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Zone's purpose: POLICY or BLOCKLIST",
			},
			"asns": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "List of asns included. Format of each array value: a string representation of an ASN numeric value. Use with type `DYNAMIC` or `DYNAMIC_V2`",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Network Status - can either be ACTIVE or INACTIVE only",
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
		Description: "Gets Okta Network Zone.",
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
		zone *v5okta.ListNetworkZones200ResponseInner
	)
	if id != "" {
		zone, _, err = getOktaV5ClientFromMetadata(m).NetworkZoneAPI.GetNetworkZone(ctx, id).Execute()
	} else {
		zone, err = findNetworkZoneByName(ctx, m, name)
	}
	if err != nil {
		return diag.Errorf("failed to find network zone: %v", err)
	}
	nzID, err := concreteNetworkzoneID(zone)
	if err != nil {
		return diag.Errorf("failed to create network zone: %v", err)
	}
	d.SetId(nzID)

	err = mapNetworkZoneToState(d, zone)
	if err != nil {
		return diag.Errorf("failed to set network zone properties: %v", err)
	}
	return nil
}

func findNetworkZoneByName(ctx context.Context, m interface{}, name string) (*v5okta.ListNetworkZones200ResponseInner, error) {
	client := getOktaV5ClientFromMetadata(m)
	zones, resp, err := client.NetworkZoneAPI.ListNetworkZones(ctx).Execute()
	if err != nil {
		return nil, err
	}
	for i := range zones {
		if getNetworkZoneName(zones[i]) == name {
			return &zones[i], nil
		}
	}
	for {
		var moreZones []v5okta.ListNetworkZones200ResponseInner
		if resp.HasNextPage() {
			resp, err = resp.Next(&moreZones)
			if err != nil {
				return nil, err
			}
			for i := range moreZones {
				if getNetworkZoneName(moreZones[i]) == name {
					return &zones[i], nil
				}
			}
		} else {
			break
		}
	}
	return nil, fmt.Errorf("network zone with name '%s' does not exist", name)
}

func getNetworkZoneName(data v5okta.ListNetworkZones200ResponseInner) string {
	var name string
	nz := data.GetActualInstance()
	switch v := nz.(type) {
	case *v5okta.DynamicNetworkZone:
		name = v.GetName()
	case *v5okta.EnhancedDynamicNetworkZone:
		name = v.GetName()
	case *v5okta.IPNetworkZone:
		name = v.GetName()
	}
	return name
}
