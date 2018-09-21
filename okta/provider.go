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
		orgName:  d.Get("org_name").(string),
		domain:   d.Get("base_url").(string),
		apiToken: d.Get("api_token").(string),
	}
	if err := config.loadAndValidate(); err != nil {
		return nil, fmt.Errorf("[ERROR] Error initializing the Okta SDK clients: %v", err)
	}
	return &config, nil
}
