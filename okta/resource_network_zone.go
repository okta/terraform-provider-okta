package okta

import (
	"fmt"
	"log"

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
				Type:        schema.TypeList,
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
	log.Printf("[INFO] Create Network Zone %v", d.Get("name").(string))
	return fmt.Errorf("[ERROR] Okta Network Zone not implemented")
}

func resourceNetworkZoneRead(d *schema.ResourceData, m interface{}) error {
	return fmt.Errorf("[ERROR] Okta Network Zone not implemented")
}

func resourceNetworkZoneUpdate(d *schema.ResourceData, m interface{}) error {
	return fmt.Errorf("[ERROR] Okta Network Zone not implemented")
}

func resourceNetworkZoneDelete(d *schema.ResourceData, m interface{}) error {
	return fmt.Errorf("[ERROR] Okta Network Zone not implemented")
}
