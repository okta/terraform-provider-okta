package okta

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/okta/okta-sdk-golang/okta/query"
)

func dataSourceUserType() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceUserTypeRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"display_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceUserTypeRead(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	userType, err := getSupplementFromMetadata(m).FindUserType(name, &query.Params{})
	if err != nil {
		return err
	}

	if userType == nil {
		return fmt.Errorf("No user type found with provided name %s", name)
	}

	d.SetId(userType.Id)
	d.Set("name", userType.Name)
	d.Set("display_name", userType.DisplayName)
	d.Set("description", userType.Description)

	return nil
}
