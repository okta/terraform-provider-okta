package idaas

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

func dataSourceAppUserProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppUserProfileRead,
		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the Okta App.",
			},
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the user assigned to the app.",
			},
			"profile": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The JSON profile of the App User.",
			},
		},
		Description: "Get the profile for a user assigned to an Okta application.",
	}
}

func dataSourceAppUserProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appID := d.Get("app_id").(string)
	userID := d.Get("user_id").(string)

	u, resp, err := getOktaClientFromMetadata(meta).Application.GetApplicationUser(ctx, appID, userID, nil)
	if utils.Is404(resp) {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get application's user: %v", err)
	}

	var rawProfile string
	if u.Profile != nil {
		p, _ := json.Marshal(u.Profile)
		rawProfile = string(p)
	}
	_ = d.Set("profile", rawProfile)
	// Use a composite ID for uniqueness
	d.SetId(appID + ":" + userID)
	return nil
}
