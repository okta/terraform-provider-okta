package idaas

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

func ResourceAppOAuthPostLogoutRedirectURI() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "managing the post logout redirect URI should be done directly on an oauth app resource",
		CreateContext:      resourceAppOAuthRedirectURICreate("okta_app_oauth_post_logout_redirect_uri"),
		ReadContext:        resourceAppOAuthRedirectURIRead("okta_app_oauth_post_logout_redirect_uri"),
		UpdateContext:      resourceAppOAuthRedirectURIUpdate("okta_app_oauth_post_logout_redirect_uri"),
		DeleteContext:      resourceAppOAuthRedirectURIDelete("okta_app_oauth_post_logout_redirect_uri"),
		Importer:           utils.CreateCustomNestedResourceImporter([]string{"app_id", "id"}, "Expecting the following format: <app_id>/<uri>"),
		Description:        "This resource allows you to manage post logout redirection URI for use in redirect-based flows.",
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
