package okta

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/sdk"
)

func dataSourceUserType() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserTypeRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceUserTypeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	userTypes, _, err := getOktaClientFromMetadata(m).UserType.ListUserTypes(ctx)
	if err != nil {
		return diag.Errorf("failed to list user types: %v", err)
	}
	name := d.Get("name").(string)
	var userType *sdk.UserType
	for _, ut := range userTypes {
		if strings.EqualFold(name, ut.Name) {
			userType = ut
		}
	}
	if userType == nil {
		return diag.Errorf("user type '%s' does not exist", name)
	}
	d.SetId(userType.Id)
	_ = d.Set("display_name", userType.DisplayName)
	_ = d.Set("description", userType.Description)

	return nil
}
