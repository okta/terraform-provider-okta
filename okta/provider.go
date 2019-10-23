// Package okta terraform configuration for an okta site
package okta

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Resource names, defined in place, used throughout the provider and tests
const (
	appAutoLogin           = "okta_app_auto_login"
	appBookmark            = "okta_app_bookmark"
	appBasicAuth           = "okta_app_basic_auth"
	appGroupAssignment     = "okta_app_group_assignment"
	appUser                = "okta_app_user"
	appOAuth               = "okta_app_oauth"
	appOAuthRedirectUri    = "okta_app_oauth_redirect_uri"
	appSaml                = "okta_app_saml"
	appSecurePasswordStore = "okta_app_secure_password_store"
	appSwa                 = "okta_app_swa"
	appThreeField          = "okta_app_three_field"
	appUserSchema          = "okta_app_user_schema"
	appUserBaseSchema      = "okta_app_user_base_schema"
	authServer             = "okta_auth_server"
	authServerClaim        = "okta_auth_server_claim"
	authServerPolicy       = "okta_auth_server_policy"
	authServerPolicyRule   = "okta_auth_server_policy_rule"
	authServerScope        = "okta_auth_server_scope"
	factor                 = "okta_factor"
	groupRoles             = "okta_group_roles"
	groupRule              = "okta_group_rule"
	identityProvider       = "okta_identity_provider"
	idpResource            = "okta_idp_oidc"
	idpSaml                = "okta_idp_saml"
	idpSamlKey             = "okta_idp_saml_key"
	idpSocial              = "okta_idp_social"
	inlineHook             = "okta_inline_hook"
	networkZone            = "okta_network_zone"
	oktaGroup              = "okta_group"
	oktaProfileMapping     = "okta_profile_mapping"
	oktaUser               = "okta_user"
	policyMfa              = "okta_policy_mfa"
	policyPassword         = "okta_policy_password"
	policyRuleIdpDiscovery = "okta_policy_rule_idp_discovery"
	policyRuleMfa          = "okta_policy_rule_mfa"
	policyRulePassword     = "okta_policy_rule_password"
	policyRuleSignOn       = "okta_policy_rule_signon"
	policySignOn           = "okta_policy_signon"
	templateEmail          = "okta_template_email"
	trustedOrigin          = "okta_trusted_origin"
	userBaseSchema         = "okta_user_base_schema"
	userSchema             = "okta_user_schema"
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
			appAutoLogin:           resourceAppAutoLogin(),
			appBookmark:            resourceAppBookmark(),
			appBasicAuth:           resourceAppBasicAuth(),
			appGroupAssignment:     resourceAppGroupAssignment(),
			appUser:                resourceAppUser(),
			appOAuth:               resourceAppOAuth(),
			appOAuthRedirectUri:    resourceAppOAuthRedirectUri(),
			appSaml:                resourceAppSaml(),
			appSecurePasswordStore: resourceAppSecurePasswordStore(),
			appSwa:                 resourceAppSwa(),
			appThreeField:          resourceAppThreeField(),
			appUserSchema:          resourceAppUserSchema(),
			appUserBaseSchema:      resourceAppUserBaseSchema(),
			authServer:             resourceAuthServer(),
			authServerClaim:        resourceAuthServerClaim(),
			authServerPolicy:       resourceAuthServerPolicy(),
			authServerPolicyRule:   resourceAuthServerPolicyRule(),
			authServerScope:        resourceAuthServerScope(),
			factor:                 resourceFactor(),
			groupRoles:             resourceGroupRoles(),
			groupRule:              resourceGroupRule(),
			idpResource:            resourceIdpOidc(),
			idpSaml:                resourceIdpSaml(),
			idpSamlKey:             resourceIdpSigningKey(),
			idpSocial:              resourceIdpSocial(),
			inlineHook:             resourceInlineHook(),
			networkZone:            resourceNetworkZone(),
			oktaGroup:              resourceGroup(),
			oktaProfileMapping:     resourceOktaProfileMapping(),
			oktaUser:               resourceUser(),
			policyMfa:              resourcePolicyMfa(),
			policyPassword:         resourcePolicyPassword(),
			policyRuleIdpDiscovery: resourcePolicyRuleIdpDiscovery(),
			policyRuleMfa:          resourcePolicyMfaRule(),
			policyRulePassword:     resourcePolicyPasswordRule(),
			policyRuleSignOn:       resourcePolicySignonRule(),
			policySignOn:           resourcePolicySignon(),
			templateEmail:          resourceTemplateEmail(),
			trustedOrigin:          resourceTrustedOrigin(),
			userSchema:             resourceUserSchema(),
			userBaseSchema:         resourceUserBaseSchema(),

			// The day I realized I was naming stuff wrong :'-(
			"okta_idp":                       deprecateIncorrectNaming(resourceIdpOidc(), idpResource),
			"okta_saml_idp":                  deprecateIncorrectNaming(resourceIdpSaml(), idpSaml),
			"okta_saml_idp_signing_key":      deprecateIncorrectNaming(resourceIdpSigningKey(), idpSamlKey),
			"okta_social_idp":                deprecateIncorrectNaming(resourceIdpSocial(), idpSocial),
			"okta_bookmark_app":              deprecateIncorrectNaming(resourceAppBookmark(), appBookmark),
			"okta_saml_app":                  deprecateIncorrectNaming(resourceAppSaml(), appSaml),
			"okta_oauth_app":                 deprecateIncorrectNaming(resourceAppOAuth(), appOAuth),
			"okta_oauth_app_redirect_uri":    deprecateIncorrectNaming(resourceAppOAuthRedirectUri(), appOAuthRedirectUri),
			"okta_auto_login_app":            deprecateIncorrectNaming(resourceAppAutoLogin(), appAutoLogin),
			"okta_secure_password_store_app": deprecateIncorrectNaming(resourceAppSecurePasswordStore(), appSecurePasswordStore),
			"okta_three_field_app":           deprecateIncorrectNaming(resourceAppThreeField(), appThreeField),
			"okta_swa_app":                   deprecateIncorrectNaming(resourceAppSwa(), appSwa),
			"okta_password_policy":           deprecateIncorrectNaming(resourcePolicyPassword(), policyPassword),
			"okta_signon_policy":             deprecateIncorrectNaming(resourcePolicySignon(), policySignOn),
			"okta_signon_policy_rule":        deprecateIncorrectNaming(resourcePolicySignonRule(), policyRuleSignOn),
			"okta_password_policy_rule":      deprecateIncorrectNaming(resourcePolicyPasswordRule(), policyRulePassword),
			"okta_mfa_policy":                deprecateIncorrectNaming(resourcePolicyMfa(), policyMfa),
			"okta_mfa_policy_rule":           deprecateIncorrectNaming(resourcePolicyMfaRule(), policyRuleMfa),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"okta_app":               dataSourceApp(),
			"okta_app_saml":          dataSourceAppSaml(),
			"okta_app_metadata_saml": dataSourceAppMetadataSaml(),
			"okta_default_policies":  deprecatedPolicies,
			"okta_default_policy":    dataSourceDefaultPolicies(),
			"okta_everyone_group":    dataSourceEveryoneGroup(),
			"okta_group":             dataSourceGroup(),
			"okta_idp_metadata_saml": dataSourceIdpMetadataSaml(),
			"okta_idp_saml":          dataSourceIdpSaml(),
			"okta_policy":            dataSourcePolicy(),
			"okta_user":              dataSourceUser(),
			"okta_users":             dataSourceUsers(),
			authServer:               dataSourceAuthServer(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func deprecateIncorrectNaming(d *schema.Resource, newResource string) *schema.Resource {
	d.DeprecationMessage = fmt.Sprintf("Resource is deprecated due to a correction in naming conventions, please use %s instead.", newResource)
	return d
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
