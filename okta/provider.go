// Package okta terraform configuration for an okta site
package okta

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/hashicorp/terraform/terraform"
)

// Resource names, defined in place, used throughout the provider and tests
const (
	authServer             = "okta_auth_server"
	authServerPolicy       = "okta_auth_server_policy"
	authServerPolicyRule   = "okta_auth_server_policy_rule"
	authServerClaim        = "okta_auth_server_claim"
	authServerScope        = "okta_auth_server_scope"
	autoLoginApp           = "okta_auto_login_app"
	bookmarkApp            = "okta_bookmark_app"
	factor                 = "okta_factor"
	identityProvider       = "okta_identity_provider"
	idpResource            = "okta_idp"
	samlIdp                = "okta_saml_idp"
	inlineHook             = "okta_inline_hook"
	mfaPolicy              = "okta_mfa_policy"
	mfaPolicyRule          = "okta_mfa_policy_rule"
	oAuthApp               = "okta_oauth_app"
	oAuthAppRedirectUri    = "okta_oauth_app_redirect_uri"
	oktaGroup              = "okta_group"
	groupRule              = "okta_group_rule"
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
	userBaseSchema         = "okta_user_base_schema"
)

// Provider establishes a client connection to an okta site
// determined by its schema string values
func Provider() terraform.ResourceProvider {
	deprecatedPolicies := dataSourceDefaultPolicies()
	deprecatedPolicies.DeprecationMessage = "This data source will be deprecated in favor of okta_default_policy or okta_policy data sources."

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
			"backoff": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Use exponential back off strategy for rate limits.",
			},
			"min_wait_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     30,
				Description: "minimum seconds to wait when rate limit is hit. We use exponential backoffs when backoff is enabled.",
			},
			"max_wait_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     300,
				Description: "maximum seconds to wait when rate limit is hit. We use exponential backoffs when backoff is enabled.",
			},
			"max_retries": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      5,
				ValidateFunc: validation.IntAtMost(100), // Have to cut it off somewhere right?
				Description:  "maximum number of retries to attempt before erroring out.",
			},
			"parallelism": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Number of concurrent requests to make within a resource where bulk operations are not possible. Take note of https://developer.okta.com/docs/api/getting_started/rate-limits.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			oktaGroup:              resourceGroup(),
			passwordPolicy:         resourcePasswordPolicy(),
			signOnPolicy:           resourceSignOnPolicy(),
			signOnPolicyRule:       resourceSignOnPolicyRule(),
			passwordPolicyRule:     resourcePasswordPolicyRule(),
			mfaPolicy:              resourceMfaPolicy(),
			mfaPolicyRule:          resourceMfaPolicyRule(),
			trustedOrigin:          resourceTrustedOrigin(),
			userSchema:             resourceUserSchema(),
			oktaUser:               resourceUser(),
			oAuthApp:               resourceOAuthApp(),
			oAuthAppRedirectUri:    resourceOAuthAppRedirectUri(),
			samlApp:                resourceSamlApp(),
			autoLoginApp:           resourceAutoLoginApp(),
			securePasswordStoreApp: resourceSecurePasswordStoreApp(),
			threeFieldApp:          resourceThreeFieldApp(),
			swaApp:                 resourceSwaApp(),
			factor:                 resourceFactor(),
			groupRule:              resourceGroupRule(),
			authServer:             resourceAuthServer(),
			authServerClaim:        resourceAuthServerClaim(),
			authServerPolicy:       resourceAuthServerPolicy(),
			authServerPolicyRule:   resourceAuthServerPolicyRule(),
			authServerScope:        resourceAuthServerScope(),
			bookmarkApp:            resourceBookmarkApp(),
			inlineHook:             resourceInlineHook(),
			idpResource:            resourceIdp(),
			samlIdp:                resourceSamlIdp(),

			// Below resources will be deprecated
			"okta_user_schemas": resourceUserSchemas(),
			identityProvider:    resourceIdentityProvider(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			authServer:              dataSourceAuthServer(),
			"okta_everyone_group":   dataSourceEveryoneGroup(),
			"okta_default_policies": deprecatedPolicies,
			"okta_default_policy":   dataSourceDefaultPolicies(),
			"okta_policy":           dataSourcePolicy(),
			"okta_group":            dataSourceGroup(),
			"okta_app":              dataSourceApp(),
			"okta_user":             dataSourceUser(),
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
		maxWait:     d.Get("max_wait_seconds").(int),
		minWait:     d.Get("min_wait_seconds").(int),
		backoff:     d.Get("backoff").(bool),
	}
	if err := config.loadAndValidate(); err != nil {
		return nil, fmt.Errorf("[ERROR] Error initializing the Okta SDK clients: %v", err)
	}
	return &config, nil
}
