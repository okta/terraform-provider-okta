package okta

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppOAuthPostLogoutRedirectURI() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppOAuthRedirectURICreate("okta_app_oauth_post_logout_redirect_uri"),
		ReadContext:   resourceAppOAuthRedirectURIRead("okta_app_oauth_post_logout_redirect_uri"),
		UpdateContext: resourceAppOAuthRedirectURIUpdate("okta_app_oauth_post_logout_redirect_uri"),
		DeleteContext: resourceAppOAuthRedirectURIDelete("okta_app_oauth_post_logout_redirect_uri"),
		Importer:      createCustomNestedResourceImporter([]string{"app_id", "id"}, "Expecting the following format: <app_id>/<uri>"),
		Schema: map[string]*schema.Schema{
			"app_id": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "OAuth application ID.",
			},
			"uri": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Post Logout Redirect URI to append to Okta OIDC application.",
			},
		},
	}
}
