package okta

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/okta/okta-sdk-golang/okta/query"
)

func dataSourceUserProfileMappingSource() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceUserProfileMappingSourceRead,

		Schema: map[string]*schema.Schema{
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
		return fmt.Errorf("Error Listing User Profile Mapping Source in Okta: %v", err)
	}

	d.SetId(mapping.ID)
	d.Set("type", mapping.Type)
	d.Set("name", mapping.Name)

	return nil
}
