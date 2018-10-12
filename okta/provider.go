// Package okta terraform configuration for an okta site
package okta

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

const signOnPolicyRule = "okta_signon_policy_rule"
const passwordPolicyRule = "okta_password_policy_rule"
const passwordPolicy = "okta_password_policy"
const signOnPolicy = "okta_signon_policy"
const oAuthApp = "okta_oauth_app"
const samlApp = "okta_saml_app"
const swaApp = "okta_swa_app"
const autoLoginApp = "okta_auto_login_app"
const securePasswordStoreApp = "okta_secure_password_store_app"
const threeFieldApp = "okta_three_field_app"

// Provider establishes a client connection to an okta site
// determined by its schema string values
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"org_name": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OKTA_ORG_NAME", nil),
				Description: "The organization to manage in Okta.",
			},
			"api_token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OKTA_API_TOKEN", nil),
				Description: "API Token granting privileges to Okta API.",
			},
			"base_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OKTA_BASE_URL", "okta.com"),
				Description: "The Okta url. (Use 'oktapreview.com' for Okta testing)",
			},
			"parallelism": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Number of concurrent requests to make within a resource where bulk operations are not possible. Take note of https://developer.okta.com/docs/api/getting_started/rate-limits.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"okta_group":             resourceGroup(),
			"okta_identity_provider": resourceIdentityProvider(),
			passwordPolicy:           resourcePasswordPolicy(),
			signOnPolicy:             resourceSignOnPolicy(),
			signOnPolicyRule:         resourceSignOnPolicyRule(),
			passwordPolicyRule:       resourcePasswordPolicyRule(),
			"okta_trusted_origin":    resourceTrustedOrigin(),
			"okta_user_schemas":      resourceUserSchemas(),
			"okta_user":              resourceUser(),
			oAuthApp:                 resourceOAuthApp(),
			samlApp:                  resourceSamlApp(),
			autoLoginApp:             resourceAutoLoginApp(),
			securePasswordStoreApp:   resourceSecurePasswordStoreApp(),
			// Bug in the SDK preventing the use of this resource https://github.com/okta/okta-sdk-golang/pull/40
			//threeFieldApp:            resourceThreeFieldApp(),
			swaApp: resourceSwaApp(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"okta_everyone_group":   dataSourceEveryoneGroup(),
			"okta_default_policies": dataSourceDefaultPolicies(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	log.Printf("[INFO] Initializing Okta client")

	config := Config{
		orgName:     d.Get("org_name").(string),
		domain:      d.Get("base_url").(string),
		apiToken:    d.Get("api_token").(string),
		parallelism: d.Get("parallelism").(int),
	}
	if err := config.loadAndValidate(); err != nil {
		return nil, fmt.Errorf("[ERROR] Error initializing the Okta SDK clients: %v", err)
	}
	return &config, nil
}
