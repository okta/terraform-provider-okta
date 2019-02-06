// Package okta terraform configuration for an okta site
package okta

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/validation"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Resource names, defined in place, used throughout the provider and tests
const (
	autoLoginApp           = "okta_auto_login_app"
	factor                 = "okta_factor"
	identityProvider       = "okta_identity_provider"
	mfaPolicy              = "okta_mfa_policy"
	mfaPolicyRule          = "okta_mfa_policy_rule"
	oAuthApp               = "okta_oauth_app"
	oktaGroup              = "okta_group"
	oktaUser               = "okta_user"
	passwordPolicy         = "okta_password_policy"
	passwordPolicyRule     = "okta_password_policy_rule"
	samlApp                = "okta_saml_app"
	securePasswordStoreApp = "okta_secure_password_store_app"
	signOnPolicy           = "okta_signon_policy"
	signOnPolicyRule       = "okta_signon_policy_rule"
	swaApp                 = "okta_swa_app"
	threeFieldApp          = "okta_three_field_app"
	trustedOrigin          = "okta_trusted_origin"
	userSchema             = "okta_user_schema"
)

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
			"max_retries": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      5,
				ValidateFunc: validation.IntAtMost(100), // Have to cut it off somewhere right?
				Description:  "maximum number of retries to attempt before erroring out. This is also related to back offs when a 429 HTTP status code is received.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			oktaGroup:          resourceGroup(),
			identityProvider:   resourceIdentityProvider(),
			passwordPolicy:     resourcePasswordPolicy(),
			signOnPolicy:       resourceSignOnPolicy(),
			signOnPolicyRule:   resourceSignOnPolicyRule(),
			passwordPolicyRule: resourcePasswordPolicyRule(),
			mfaPolicy:          resourceMfaPolicy(),
			mfaPolicyRule:      resourceMfaPolicyRule(),
			trustedOrigin:      resourceTrustedOrigin(),
			// Will be deprecated
			"okta_user_schemas":    resourceUserSchemas(),
			userSchema:             resourceUserSchema(),
			oktaUser:               resourceUser(),
			oAuthApp:               resourceOAuthApp(),
			samlApp:                resourceSamlApp(),
			autoLoginApp:           resourceAutoLoginApp(),
			securePasswordStoreApp: resourceSecurePasswordStoreApp(),
			threeFieldApp:          resourceThreeFieldApp(),
			swaApp:                 resourceSwaApp(),
			factor:                 resourceFactor(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"okta_everyone_group":   dataSourceEveryoneGroup(),
			"okta_default_policies": dataSourceDefaultPolicies(),
			"okta_group":            dataSourceGroup(),
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
		retryCount:  d.Get("max_retries").(int),
	}
	if err := config.loadAndValidate(); err != nil {
		return nil, fmt.Errorf("[ERROR] Error initializing the Okta SDK clients: %v", err)
	}
	return &config, nil
}
