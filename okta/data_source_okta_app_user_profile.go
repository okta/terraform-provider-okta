package okta

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

func dataSourceAppUserProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppUserProfileRead,
		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the Okta App being queried for",
				ForceNew:    true,
			},
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the User associated with the application",
				ForceNew:    true,
			},
			"profile": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The application profile for the user",
				ForceNew:    true,
			},
		},
	}
}

func dataSourceAppUserProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := getOktaClientFromMetadata(m)
	appId := d.Get("app_id").(string)
	userId := d.Get("user_id").(string)

	var appUser *okta.AppUser

	appUser, _, err := client.Application.GetApplicationUser(ctx, appId, userId, &query.Params{})
	if err != nil {
		return diag.Errorf("unable to get profile for user (%s) assisgned to app (%s): %s", userId, appId, err)
	}

	jsonProfile, err := json.Marshal(appUser.Profile)
	if err != nil {
		return diag.Errorf("failed to marshal app user profile to JSON: %v", err)
	}
	_ = d.Set("profile", string(jsonProfile))

	_ = d.Set("user_id", userId)
	d.SetId(appId)
	return nil
}
