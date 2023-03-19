// Package okta terraform configuration for an okta site
package okta

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/okta/terraform-provider-okta/okta/internal/mutexkv"
)

// Resource names, defined in place, used throughout the provider and tests
const (
	adminRoleCustom               = "okta_admin_role_custom"
	adminRoleCustomAssignments    = "okta_admin_role_custom_assignments"
	adminRoleTargets              = "okta_admin_role_targets"
	app                           = "okta_app"
	appAutoLogin                  = "okta_app_auto_login"
	appBasicAuth                  = "okta_app_basic_auth"
	appBookmark                   = "okta_app_bookmark"
	appGroupAssignment            = "okta_app_group_assignment"
	appGroupAssignments           = "okta_app_group_assignments"
	appMetadataSaml               = "okta_app_metadata_saml"
	appOAuth                      = "okta_app_oauth"
	appOAuthAPIScope              = "okta_app_oauth_api_scope"
	appOAuthPostLogoutRedirectURI = "okta_app_oauth_post_logout_redirect_uri"
	appOAuthRedirectURI           = "okta_app_oauth_redirect_uri"
	appSaml                       = "okta_app_saml"
	appSamlAppSettings            = "okta_app_saml_app_settings"
	appSecurePasswordStore        = "okta_app_secure_password_store"
	appSharedCredentials          = "okta_app_shared_credentials"
	appSignOnPolicy               = "okta_app_signon_policy"
	appSignOnPolicyRule           = "okta_app_signon_policy_rule"
	appSwa                        = "okta_app_swa"
	appThreeField                 = "okta_app_three_field"
	appUser                       = "okta_app_user"
	appUserAssignments            = "okta_app_user_assignments"
	appUserBaseSchemaProperty     = "okta_app_user_base_schema_property"
	appUserSchemaProperty         = "okta_app_user_schema_property"
	authenticator                 = "okta_authenticator"
	authServer                    = "okta_auth_server"
	authServerClaim               = "okta_auth_server_claim"
	authServerClaimDefault        = "okta_auth_server_claim_default"
	authServerClaims              = "okta_auth_server_claims"
	authServerDefault             = "okta_auth_server_default"
	authServerPolicy              = "okta_auth_server_policy"
	authServerPolicyRule          = "okta_auth_server_policy_rule"
	authServerScope               = "okta_auth_server_scope"
	authServerScopes              = "okta_auth_server_scopes"
	behavior                      = "okta_behavior"
	behaviors                     = "okta_behaviors"
	brand                         = "okta_brand"
	brands                        = "okta_brands"
	captcha                       = "okta_captcha"
	captchaOrgWideSettings        = "okta_captcha_org_wide_settings"
	defaultPolicies               = "okta_default_policies"
	defaultPolicy                 = "okta_default_policy"
	device                        = "okta_device"
	devices                       = "okta_devices"
	domain                        = "okta_domain"
	domainCertificate             = "okta_domain_certificate"
	domainVerification            = "okta_domain_verification"
	emailSender                   = "okta_email_sender"
	emailSenderVerification       = "okta_email_sender_verification"
	emailCustomization            = "okta_email_customization"
	emailCustomizations           = "okta_email_customizations"
	emailTemplate                 = "okta_email_template"
	emailTemplates                = "okta_email_templates"
	eventHook                     = "okta_event_hook"
	eventHookVerification         = "okta_event_hook_verification"
	factor                        = "okta_factor"
	factorTotp                    = "okta_factor_totp"
	group                         = "okta_group"
	groupEveryone                 = "okta_everyone_group"
	groupMembership               = "okta_group_membership"
	groupMemberships              = "okta_group_memberships"
	groupRole                     = "okta_group_role"
	groupRoles                    = "okta_group_roles"
	groupRule                     = "okta_group_rule"
	groups                        = "okta_groups"
	groupSchemaProperty           = "okta_group_schema_property"
	idpMetadataSaml               = "okta_idp_metadata_saml"
	idpOidc                       = "okta_idp_oidc"
	idpSaml                       = "okta_idp_saml"
	idpSamlKey                    = "okta_idp_saml_key"
	idpSocial                     = "okta_idp_social"
	inlineHook                    = "okta_inline_hook"
	linkDefinition                = "okta_link_definition"
	linkValue                     = "okta_link_value"
	networkZone                   = "okta_network_zone"
	orgConfiguration              = "okta_org_configuration"
	orgSupport                    = "okta_org_support"
	policy                        = "okta_policy"
	policyMfa                     = "okta_policy_mfa"
	policyMfaDefault              = "okta_policy_mfa_default"
	policyPassword                = "okta_policy_password"
	policyPasswordDefault         = "okta_policy_password_default"
	policyProfileEnrollment       = "okta_policy_profile_enrollment"
	policyProfileEnrollmentApps   = "okta_policy_profile_enrollment_apps"
	policyRuleIdpDiscovery        = "okta_policy_rule_idp_discovery"
	policyRuleMfa                 = "okta_policy_rule_mfa"
	policyRulePassword            = "okta_policy_rule_password"
	policyRuleProfileEnrollment   = "okta_policy_rule_profile_enrollment"
	policyRuleSignOn              = "okta_policy_rule_signon"
	policySignOn                  = "okta_policy_signon"
	profileMapping                = "okta_profile_mapping"
	rateLimiting                  = "okta_rate_limiting"
	resourceSet                   = "okta_resource_set"
	roleSubscription              = "okta_role_subscription"
	securityNotificationEmails    = "okta_security_notification_emails"
	templateEmail                 = "okta_template_email"
	templateSms                   = "okta_template_sms"
	theme                         = "okta_theme"
	themes                        = "okta_themes"
	threatInsightSettings         = "okta_threat_insight_settings"
	trustedOrigin                 = "okta_trusted_origin"
	trustedOrigins                = "okta_trusted_origins"
	user                          = "okta_user"
	userAdminRoles                = "okta_user_admin_roles"
	userBaseSchemaProperty        = "okta_user_base_schema_property"
	userFactorQuestion            = "okta_user_factor_question"
	userGroupMemberships          = "okta_user_group_memberships"
	userProfileMappingSource      = "okta_user_profile_mapping_source"
	users                         = "okta_users"
	userSchemaProperty            = "okta_user_schema_property"
	userSecurityQuestions         = "okta_user_security_questions"
	userType                      = "okta_user_type"
)

// Provider establishes a client connection to an okta site
// determined by its schema string values
func Provider() *schema.Provider {
	deprecatedPolicies := dataSourceDefaultPolicy()
	deprecatedPolicies.DeprecationMessage = "This data source will be deprecated in favor of okta_default_policy or okta_policy data sources."
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"org_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OKTA_ORG_NAME", nil),
				Description: "The organization to manage in Okta.",
			},
			"access_token": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("OKTA_ACCESS_TOKEN", nil),
				Description:   "Bearer token granting privileges to Okta API.",
				ConflictsWith: []string{"api_token", "client_id", "scopes", "private_key"},
			},
			"api_token": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("OKTA_API_TOKEN", nil),
				Description:   "API Token granting privileges to Okta API.",
				ConflictsWith: []string{"access_token", "client_id", "scopes", "private_key"},
			},
			"client_id": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("OKTA_API_CLIENT_ID", nil),
				Description:   "API Token granting privileges to Okta API.",
				ConflictsWith: []string{"access_token", "api_token"},
			},
			"scopes": {
				Type:          schema.TypeSet,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				DefaultFunc:   envDefaultSetFunc("OKTA_API_SCOPES", nil),
				Description:   "API Token granting privileges to Okta API.",
				ConflictsWith: []string{"access_token", "api_token"},
			},
			"private_key": {
				Optional:      true,
				Type:          schema.TypeString,
				DefaultFunc:   schema.EnvDefaultFunc("OKTA_API_PRIVATE_KEY", nil),
				Description:   "API Token granting privileges to Okta API.",
				ConflictsWith: []string{"access_token", "api_token"},
			},
			"private_key_id": {
				Optional:      true,
				Type:          schema.TypeString,
				DefaultFunc:   schema.EnvDefaultFunc("OKTA_API_PRIVATE_KEY_ID", nil),
				Description:   "API Token Id granting privileges to Okta API.",
				ConflictsWith: []string{"api_token"},
			},
			"base_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OKTA_BASE_URL", "okta.com"),
				Description: "The Okta url. (Use 'oktapreview.com' for Okta testing)",
			},
			"http_proxy": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OKTA_HTTP_PROXY", ""),
				Description: "Alternate HTTP proxy of scheme://hostname or scheme://hostname:port format",
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
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          5,
				ValidateDiagFunc: intAtMost(100), // Have to cut it off somewhere right?
				Description:      "maximum number of retries to attempt before erroring out.",
			},
			"parallelism": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Number of concurrent requests to make within a resource where bulk operations are not possible. Take note of https://developer.okta.com/docs/api/getting_started/rate-limits.",
			},
			"log_level": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          int(hclog.Error),
				ValidateDiagFunc: intBetween(1, 5),
				Description:      "providers log level. Minimum is 1 (TRACE), and maximum is 5 (ERROR)",
			},
			"max_api_capacity": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateDiagFunc: intBetween(1, 100),
				DefaultFunc:      schema.EnvDefaultFunc("MAX_API_CAPACITY", 100),
				Description: "(Experimental) sets what percentage of capacity the provider can use of the total rate limit " +
					"capacity while making calls to the Okta management API endpoints. Okta API operates in one minute buckets. " +
					"See Okta Management API Rate Limits: https://developer.okta.com/docs/reference/rl-global-mgmt/",
			},
			"request_timeout": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          0,
				ValidateDiagFunc: intBetween(0, 300),
				Description:      "Timeout for single request (in seconds) which is made to Okta, the default is `0` (means no limit is set). The maximum value can be `300`.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			adminRoleCustom:               resourceAdminRoleCustom(),
			adminRoleCustomAssignments:    resourceAdminRoleCustomAssignments(),
			adminRoleTargets:              resourceAdminRoleTargets(),
			appAutoLogin:                  resourceAppAutoLogin(),
			appBasicAuth:                  resourceAppBasicAuth(),
			appBookmark:                   resourceAppBookmark(),
			appGroupAssignment:            resourceAppGroupAssignment(),
			appGroupAssignments:           resourceAppGroupAssignments(),
			appOAuth:                      resourceAppOAuth(),
			appOAuthAPIScope:              resourceAppOAuthAPIScope(),
			appOAuthPostLogoutRedirectURI: resourceAppOAuthPostLogoutRedirectURI(),
			appOAuthRedirectURI:           resourceAppOAuthRedirectURI(),
			appSaml:                       resourceAppSaml(),
			appSamlAppSettings:            resourceAppSamlAppSettings(),
			appSecurePasswordStore:        resourceAppSecurePasswordStore(),
			appSharedCredentials:          resourceAppSharedCredentials(),
			appSignOnPolicy:               resourceAppSignOnPolicy(),
			appSignOnPolicyRule:           resourceAppSignOnPolicyRule(),
			appSwa:                        resourceAppSwa(),
			appThreeField:                 resourceAppThreeField(),
			appUser:                       resourceAppUser(),
			appUserBaseSchemaProperty:     resourceAppUserBaseSchemaProperty(),
			appUserSchemaProperty:         resourceAppUserSchemaProperty(),
			authenticator:                 resourceAuthenticator(),
			authServer:                    resourceAuthServer(),
			authServerClaim:               resourceAuthServerClaim(),
			authServerClaimDefault:        resourceAuthServerClaimDefault(),
			authServerDefault:             resourceAuthServerDefault(),
			authServerPolicy:              resourceAuthServerPolicy(),
			authServerPolicyRule:          resourceAuthServerPolicyRule(),
			authServerScope:               resourceAuthServerScope(),
			behavior:                      resourceBehavior(),
			brand:                         resourceBrand(),
			captcha:                       resourceCaptcha(),
			captchaOrgWideSettings:        resourceCaptchaOrgWideSettings(),
			domain:                        resourceDomain(),
			domainCertificate:             resourceDomainCertificate(),
			domainVerification:            resourceDomainVerification(),
			emailCustomization:            resourceEmailCustomization(),
			emailSender:                   resourceEmailSender(),
			emailSenderVerification:       resourceEmailSenderVerification(),
			eventHook:                     resourceEventHook(),
			eventHookVerification:         resourceEventHookVerification(),
			factor:                        resourceFactor(),
			factorTotp:                    resourceFactorTOTP(),
			group:                         resourceGroup(),
			groupMembership:               resourceGroupMembership(),
			groupMemberships:              resourceGroupMemberships(),
			groupRole:                     resourceGroupRole(),
			groupRoles:                    resourceGroupRoles(),
			groupRule:                     resourceGroupRule(),
			groupSchemaProperty:           resourceGroupCustomSchemaProperty(),
			idpOidc:                       resourceIdpOidc(),
			idpSaml:                       resourceIdpSaml(),
			idpSamlKey:                    resourceIdpSigningKey(),
			idpSocial:                     resourceIdpSocial(),
			inlineHook:                    resourceInlineHook(),
			linkDefinition:                resourceLinkDefinition(),
			linkValue:                     resourceLinkValue(),
			networkZone:                   resourceNetworkZone(),
			orgConfiguration:              resourceOrgConfiguration(),
			orgSupport:                    resourceOrgSupport(),
			policyMfa:                     resourcePolicyMfa(),
			policyMfaDefault:              resourcePolicyMfaDefault(),
			policyPassword:                resourcePolicyPassword(),
			policyPasswordDefault:         resourcePolicyPasswordDefault(),
			policyProfileEnrollment:       resourcePolicyProfileEnrollment(),
			policyProfileEnrollmentApps:   resourcePolicyProfileEnrollmentApps(),
			policyRuleIdpDiscovery:        resourcePolicyRuleIdpDiscovery(),
			policyRuleMfa:                 resourcePolicyMfaRule(),
			policyRulePassword:            resourcePolicyPasswordRule(),
			policyRuleProfileEnrollment:   resourcePolicyProfileEnrollmentRule(),
			policyRuleSignOn:              resourcePolicySignOnRule(),
			policySignOn:                  resourcePolicySignOn(),
			profileMapping:                resourceProfileMapping(),
			rateLimiting:                  resourceRateLimiting(),
			resourceSet:                   resourceResourceSet(),
			roleSubscription:              resourceRoleSubscription(),
			securityNotificationEmails:    resourceSecurityNotificationEmails(),
			templateEmail:                 resourceTemplateEmail(),
			templateSms:                   resourceTemplateSms(),
			theme:                         resourceTheme(),
			threatInsightSettings:         resourceThreatInsightSettings(),
			trustedOrigin:                 resourceTrustedOrigin(),
			user:                          resourceUser(),
			userAdminRoles:                resourceUserAdminRoles(),
			userBaseSchemaProperty:        resourceUserBaseSchemaProperty(),
			userFactorQuestion:            resourceUserFactorQuestion(),
			userGroupMemberships:          resourceUserGroupMemberships(),
			userSchemaProperty:            resourceUserCustomSchemaProperty(),
			userType:                      resourceUserType(),

			// The day I realized I was naming stuff wrong :'-(
			"okta_idp":                       deprecateIncorrectNaming(resourceIdpOidc(), idpOidc),
			"okta_saml_idp":                  deprecateIncorrectNaming(resourceIdpSaml(), idpSaml),
			"okta_saml_idp_signing_key":      deprecateIncorrectNaming(resourceIdpSigningKey(), idpSamlKey),
			"okta_social_idp":                deprecateIncorrectNaming(resourceIdpSocial(), idpSocial),
			"okta_bookmark_app":              deprecateIncorrectNaming(resourceAppBookmark(), appBookmark),
			"okta_saml_app":                  deprecateIncorrectNaming(resourceAppSaml(), appSaml),
			"okta_oauth_app":                 deprecateIncorrectNaming(resourceAppOAuth(), appOAuth),
			"okta_oauth_app_redirect_uri":    deprecateIncorrectNaming(resourceAppOAuthRedirectURI(), appOAuthRedirectURI),
			"okta_auto_login_app":            deprecateIncorrectNaming(resourceAppAutoLogin(), appAutoLogin),
			"okta_secure_password_store_app": deprecateIncorrectNaming(resourceAppSecurePasswordStore(), appSecurePasswordStore),
			"okta_three_field_app":           deprecateIncorrectNaming(resourceAppThreeField(), appThreeField),
			"okta_swa_app":                   deprecateIncorrectNaming(resourceAppSwa(), appSwa),
			"okta_password_policy":           deprecateIncorrectNaming(resourcePolicyPassword(), policyPassword),
			"okta_signon_policy":             deprecateIncorrectNaming(resourcePolicySignOn(), policySignOn),
			"okta_signon_policy_rule":        deprecateIncorrectNaming(resourcePolicySignOnRule(), policyRuleSignOn),
			"okta_password_policy_rule":      deprecateIncorrectNaming(resourcePolicyPasswordRule(), policyRulePassword),
			"okta_mfa_policy":                deprecateIncorrectNaming(resourcePolicyMfa(), policyMfa),
			"okta_mfa_policy_rule":           deprecateIncorrectNaming(resourcePolicyMfaRule(), policyRuleMfa),
			"okta_app_user_schema":           deprecateIncorrectNaming(resourceAppUserSchemaProperty(), appUserSchemaProperty),
			"okta_app_user_base_schema":      deprecateIncorrectNaming(resourceAppUserBaseSchemaProperty(), appUserBaseSchemaProperty),
			"okta_user_schema":               deprecateIncorrectNaming(resourceUserCustomSchemaProperty(), userSchemaProperty),
			"okta_user_base_schema":          deprecateIncorrectNaming(resourceUserBaseSchemaProperty(), userBaseSchemaProperty),
		},
		DataSourcesMap: map[string]*schema.Resource{
			app:                      dataSourceApp(),
			appGroupAssignments:      dataSourceAppGroupAssignments(),
			appMetadataSaml:          dataSourceAppMetadataSaml(),
			appOAuth:                 dataSourceAppOauth(),
			appSaml:                  dataSourceAppSaml(),
			appSignOnPolicy:          dataSourceAppSignOnPolicy(),
			appUserAssignments:       dataSourceAppUserAssignments(),
			authenticator:            dataSourceAuthenticator(),
			authServer:               dataSourceAuthServer(),
			authServerClaim:          dataSourceAuthServerClaim(),
			authServerClaims:         dataSourceAuthServerClaims(),
			authServerPolicy:         dataSourceAuthServerPolicy(),
			authServerScopes:         dataSourceAuthServerScopes(),
			behavior:                 dataSourceBehavior(),
			behaviors:                dataSourceBehaviors(),
			brand:                    dataSourceBrand(),
			brands:                   dataSourceBrands(),
			domain:                   dataSourceDomain(),
			device:                   dataSourceDevice(),
			devices:                  dataSourceDevices(),
			emailCustomization:       dataSourceEmailCustomization(),
			emailCustomizations:      dataSourceEmailCustomizations(),
			emailTemplate:            dataSourceEmailTemplate(),
			emailTemplates:           dataSourceEmailTemplates(),
			defaultPolicies:          deprecatedPolicies,
			defaultPolicy:            dataSourceDefaultPolicy(),
			group:                    dataSourceGroup(),
			groupEveryone:            dataSourceEveryoneGroup(),
			groups:                   dataSourceGroups(),
			idpMetadataSaml:          dataSourceIdpMetadataSaml(),
			idpOidc:                  dataSourceIdpOidc(),
			idpSaml:                  dataSourceIdpSaml(),
			idpSocial:                dataSourceIdpSocial(),
			networkZone:              dataSourceNetworkZone(),
			policy:                   dataSourcePolicy(),
			roleSubscription:         dataSourceRoleSubscription(),
			theme:                    dataSourceTheme(),
			themes:                   dataSourceThemes(),
			trustedOrigins:           dataSourceTrustedOrigins(),
			user:                     dataSourceUser(),
			userProfileMappingSource: dataSourceUserProfileMappingSource(),
			users:                    dataSourceUsers(),
			userSecurityQuestions:    dataSourceUserSecurityQuestions(),
			userType:                 dataSourceUserType(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func deprecateIncorrectNaming(d *schema.Resource, newResource string) *schema.Resource {
	d.DeprecationMessage = fmt.Sprintf("Resource is deprecated due to a correction in naming conventions, please use '%s' instead.", newResource)
	return d
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	log.Printf("[INFO] Initializing Okta client")
	config := Config{
		orgName:        d.Get("org_name").(string),
		domain:         d.Get("base_url").(string),
		apiToken:       d.Get("api_token").(string),
		accessToken:    d.Get("access_token").(string),
		clientID:       d.Get("client_id").(string),
		privateKey:     d.Get("private_key").(string),
		privateKeyId:   d.Get("private_key_id").(string),
		scopes:         convertInterfaceToStringSet(d.Get("scopes")),
		retryCount:     d.Get("max_retries").(int),
		parallelism:    d.Get("parallelism").(int),
		backoff:        d.Get("backoff").(bool),
		minWait:        d.Get("min_wait_seconds").(int),
		maxWait:        d.Get("max_wait_seconds").(int),
		logLevel:       d.Get("log_level").(int),
		requestTimeout: d.Get("request_timeout").(int),
		maxAPICapacity: d.Get("max_api_capacity").(int),
	}

	if httpProxy, ok := d.Get("http_proxy").(string); ok {
		config.httpProxy = httpProxy
	}

	if v := os.Getenv("OKTA_API_SCOPES"); v != "" && len(config.scopes) == 0 {
		config.scopes = strings.Split(v, ",")
	}
	if err := config.loadAndValidate(ctx); err != nil {
		return nil, diag.Errorf("[ERROR] invalid configuration: %v", err)
	}

	// Discover if the Okta Org is Classic or OIE
	if org, _, err := config.supplementClient.GetWellKnownOktaOrganization(ctx); err == nil {
		config.classicOrg = (org.Pipeline == "v1") // v1 == Classic, idx == OIE
	}

	return &config, nil
}

// This is a global MutexKV for use within this plugin.
var oktaMutexKV = mutexkv.NewMutexKV()

func envDefaultSetFunc(k string, dv interface{}) schema.SchemaDefaultFunc {
	return func() (interface{}, error) {
		if v := os.Getenv(k); v != "" {
			stringList := strings.Split(v, ",")
			arr := make([]interface{}, len(stringList))
			for i := range stringList {
				arr[i] = stringList[i]
			}
			return arr, nil
		}
		return dv, nil
	}
}

func isClassicOrg(m interface{}) bool {
	if config, ok := m.(*Config); ok && config.classicOrg {
		return true
	}
	return false
}

func oieOnlyFeatureError(kind, name string) diag.Diagnostics {
	url := fmt.Sprintf("https://registry.terraform.io/providers/okta/okta/latest/docs/%s/%s", kind, string(name[5:]))
	if kind == "resources" {
		kind = "resource"
	}
	if kind == "data-sources" {
		kind = "datasource"
	}
	return diag.Errorf("%q is a %s for OIE Orgs only, see %s", name, kind, url)
}

func resourceOIEOnlyFeatureError(name string) diag.Diagnostics {
	return oieOnlyFeatureError("resources", name)
}

func datasourceOIEOnlyFeatureError(name string) diag.Diagnostics {
	return oieOnlyFeatureError("data-sources", name)
}

func resourceFuncNoOp(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return nil
}

func int64Ptr(what int) *int64 {
	result := int64(what)
	return &result
}
