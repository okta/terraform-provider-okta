package okta

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

func dataSourceUserProfileMappingSource() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceUserProfileMappingSourceRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the source",
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceUserProfileMappingSourceRead(d *schema.ResourceData, m interface{}) error {
	client := getSupplementFromMetadata(m)

	mapping, _, err := client.FindProfileMappingSource("user", "user", &query.Params{})
	if err != nil {
		return fmt.Errorf("error Listing User Profile Mapping Source in Okta: %v", err)
	}

	d.SetId(mapping.ID)
	_ = d.Set("type", mapping.Type)
	_ = d.Set("name", mapping.Name)

	return nil
}
