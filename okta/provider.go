// Package okta terraform configuration for an okta site
package okta

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cenkalti/backoff"
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
	apps                          = "okta_apps"
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
	appOAuthRoleAssignment        = "okta_app_oauth_role_assignment"
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
	defaultPolicy                 = "okta_default_policy"
	domain                        = "okta_domain"
	domainCertificate             = "okta_domain_certificate"
	domainVerification            = "okta_domain_verification"
	emailDomain                   = "okta_email_domain"
	emailDomainVerification       = "okta_email_domain_verification"
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
	groupOwner                    = "okta_group_owner"
	groupEveryone                 = "okta_everyone_group"
	groupMemberships              = "okta_group_memberships"
	groupRole                     = "okta_group_role"
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
	logStream                     = "okta_log_stream"
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

// oktaMutexKV is a global MutexKV for use within this plugin
var oktaMutexKV = mutexkv.NewMutexKV()

// Provider establishes a client connection to an okta site
// determined by its schema string values
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"org_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The organization to manage in Okta.",
			},
			"access_token": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Bearer token granting privileges to Okta API.",
				ConflictsWith: []string{"api_token", "client_id", "scopes", "private_key"},
			},
			"api_token": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "API Token granting privileges to Okta API.",
				ConflictsWith: []string{"access_token", "client_id", "scopes", "private_key"},
			},
			"client_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "API Token granting privileges to Okta API.",
				ConflictsWith: []string{"access_token", "api_token"},
			},
			"scopes": {
				Type:          schema.TypeSet,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Description:   "API Token granting privileges to Okta API.",
				ConflictsWith: []string{"access_token", "api_token"},
			},
			"private_key": {
				Optional:      true,
				Type:          schema.TypeString,
				Description:   "API Token granting privileges to Okta API.",
				ConflictsWith: []string{"access_token", "api_token"},
			},
			"private_key_id": {
				Optional:      true,
				Type:          schema.TypeString,
				Description:   "API Token Id granting privileges to Okta API.",
				ConflictsWith: []string{"api_token"},
			},
			"base_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Okta url. (Use 'oktapreview.com' for Okta testing)",
			},
			"http_proxy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Alternate HTTP proxy of scheme://hostname or scheme://hostname:port format",
			},
			"backoff": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Use exponential back off strategy for rate limits.",
			},
			"min_wait_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "minimum seconds to wait when rate limit is hit. We use exponential backoffs when backoff is enabled.",
			},
			"max_wait_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "maximum seconds to wait when rate limit is hit. We use exponential backoffs when backoff is enabled.",
			},
			"max_retries": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateDiagFunc: intAtMost(100),
				Description:      "maximum number of retries to attempt before erroring out.",
			},
			"parallelism": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Number of concurrent requests to make within a resource where bulk operations are not possible. Take note of https://developer.okta.com/docs/api/getting_started/rate-limits.",
			},
			"log_level": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateDiagFunc: intBetween(1, 5),
				Description:      "providers log level. Minimum is 1 (TRACE), and maximum is 5 (ERROR)",
			},
			"max_api_capacity": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateDiagFunc: intBetween(1, 100),
				Description: "(Experimental) sets what percentage of capacity the provider can use of the total rate limit " +
					"capacity while making calls to the Okta management API endpoints. Okta API operates in one minute buckets. " +
					"See Okta Management API Rate Limits: https://developer.okta.com/docs/reference/rl-global-mgmt/",
			},
			"request_timeout": {
				Type:             schema.TypeInt,
				Optional:         true,
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
			captcha:                       resourceCaptcha(),
			captchaOrgWideSettings:        resourceCaptchaOrgWideSettings(),
			domain:                        resourceDomain(),
			domainCertificate:             resourceDomainCertificate(),
			domainVerification:            resourceDomainVerification(),
			emailCustomization:            resourceEmailCustomization(),
			emailDomain:                   resourceEmailDomain(),
			emailDomainVerification:       resourceEmailDomainVerification(),
			emailSender:                   resourceEmailSender(),
			emailSenderVerification:       resourceEmailSenderVerification(),
			eventHook:                     resourceEventHook(),
			eventHookVerification:         resourceEventHookVerification(),
			factor:                        resourceFactor(),
			factorTotp:                    resourceFactorTOTP(),
			group:                         resourceGroup(),
			groupMemberships:              resourceGroupMemberships(),
			groupRole:                     resourceGroupRole(),
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
		},
		DataSourcesMap: map[string]*schema.Resource{
			app:                      dataSourceApp(),
			appGroupAssignments:      dataSourceAppGroupAssignments(),
			appMetadataSaml:          dataSourceAppMetadataSaml(),
			appOAuth:                 dataSourceAppOauth(),
			appSaml:                  dataSourceAppSaml(),
			appSignOnPolicy:          dataSourceAppSignOnPolicy(),
			appUserAssignments:       dataSourceAppUserAssignments(),
			appUserProfile:           dataSourceAppUserProfile(),
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
			emailCustomization:       dataSourceEmailCustomization(),
			emailCustomizations:      dataSourceEmailCustomizations(),
			emailTemplate:            dataSourceEmailTemplate(),
			emailTemplates:           dataSourceEmailTemplates(),
			defaultPolicy:            dataSourceDefaultPolicy(),
			group:                    dataSourceGroup(),
			groupEveryone:            dataSourceEveryoneGroup(),
			groupRule:                dataSourceGroupRule(),
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
		},
		ConfigureContextFunc: providerConfigure,
	}
}

// providerConfigure is only called once when a terraform command is run but it
// will be called many times while running different ACC tests
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	log.Printf("[INFO] Initializing Okta client")
	config := NewConfig(d)
	if err := config.loadClients(ctx); err != nil {
		return nil, diag.Errorf("[ERROR] failed to load sdk clients: %v", err)
	}
	config.SetTimeOperations(NewProductionTimeOperations())

	// NOTE: production runtime needs to know about VCR test environment for
	// this one case where the validate function calls GET /api/v1/users/me to
	// quickly verify if the operator's auth settings are correct.
	if os.Getenv("OKTA_VCR_TF_ACC") == "" {
		if err := config.verifyCredentials(ctx); err != nil {
			return nil, diag.Errorf("[ERROR] failed validate configuration: %v", err)
		}
	}

	return config, nil
}

func isClassicOrg(ctx context.Context, m interface{}) bool {
	if config, ok := m.(*Config); ok && config.IsClassicOrg(ctx) {
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

// newExponentialBackOffWithContext helper to dry up creating a backoff object that is exponential and has context
func newExponentialBackOffWithContext(ctx context.Context, maxElapsedTime time.Duration) backoff.BackOffContext {
	bOff := backoff.NewExponentialBackOff()
	bOff.MaxElapsedTime = maxElapsedTime

	// NOTE: backoff.BackOffContext is an interface that embeds backoff.Backoff
	// so the greater context is considered on backoff.Retry
	return backoff.WithContext(bOff, ctx)
}

// doNotRetry helper function to flag if provider should be using backoff.Retry
func doNotRetry(m interface{}, err error) bool {
	return m.(*Config).timeOperations.DoNotRetry(err)
}
