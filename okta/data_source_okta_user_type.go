package okta

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func dataSourceUserType() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceUserTypeRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceUserTypeRead(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	userTypes, _, err := getOktaClientFromMetadata(m).UserType.ListUserTypes(context.Background())
	if err != nil {
		return err
	}
	var userType *okta.UserType
	for _, ut := range userTypes {
		if strings.EqualFold(name, ut.Name) {
			userType = ut
		}
	}

	if userType == nil {
		return fmt.Errorf("no user type found with provided name %s", name)
	}

	d.SetId(userType.Id)
	_ = d.Set("name", userType.Name)
	_ = d.Set("display_name", userType.DisplayName)
	_ = d.Set("description", userType.Description)

	return nil
}
