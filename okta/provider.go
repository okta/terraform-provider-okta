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
)

// Resource names, defined in place, used throughout the provider and tests
const (
	adminRoleTargets           = "okta_admin_role_targets"
	appAutoLogin               = "okta_app_auto_login"
	appBookmark                = "okta_app_bookmark"
	appBasicAuth               = "okta_app_basic_auth"
	appGroupAssignment         = "okta_app_group_assignment"
	appGroupAssignments        = "okta_app_group_assignments"
	appUser                    = "okta_app_user"
	appOAuth                   = "okta_app_oauth"
	appOAuthAPIScope           = "okta_app_oauth_api_scope"
	appOAuthRedirectURI        = "okta_app_oauth_redirect_uri"
	appSaml                    = "okta_app_saml"
	appSignOnPolicy            = "okta_app_signon_policy"
	appSignOnPolicyRule        = "okta_app_signon_policy_rule"
	appSamlAppSettings         = "okta_app_saml_app_settings"
	appSecurePasswordStore     = "okta_app_secure_password_store"
	appSwa                     = "okta_app_swa"
	appSharedCredentials       = "okta_app_shared_credentials"
	appThreeField              = "okta_app_three_field"
	appUserSchemaProperty      = "okta_app_user_schema_property"
	appUserBaseSchemaProperty  = "okta_app_user_base_schema_property"
	authServer                 = "okta_auth_server"
	authServerDefault          = "okta_auth_server_default"
	authServerClaim            = "okta_auth_server_claim"
	authServerClaimDefault     = "okta_auth_server_claim_default"
	authServerPolicy           = "okta_auth_server_policy"
	authServerPolicyRule       = "okta_auth_server_policy_rule"
	authServerScope            = "okta_auth_server_scope"
	behavior                   = "okta_behavior"
	behaviors                  = "okta_behaviors"
	domain                     = "okta_domain"
	domainVerification         = "okta_domain_verification"
	domainCertificate          = "okta_domain_certificate"
	eventHook                  = "okta_event_hook"
	emailSender                = "okta_email_sender"
	emailSenderVerification    = "okta_email_sender_verification"
	factor                     = "okta_factor"
	factorTotp                 = "okta_factor_totp"
	groupRole                  = "okta_group_role"
	groupRoles                 = "okta_group_roles"
	groupRule                  = "okta_group_rule"
	groupSchemaProperty        = "okta_group_schema_property"
	idpOidc                    = "okta_idp_oidc"
	idpSaml                    = "okta_idp_saml"
	idpSamlKey                 = "okta_idp_saml_key"
	idpSocial                  = "okta_idp_social"
	inlineHook                 = "okta_inline_hook"
	networkZone                = "okta_network_zone"
	oktaGroup                  = "okta_group"
	oktaGroups                 = "okta_groups"
	oktaGroupMembership        = "okta_group_membership"
	oktaGroupMemberships       = "okta_group_memberships"
	oktaProfileMapping         = "okta_profile_mapping"
	oktaUser                   = "okta_user"
	policyMfa                  = "okta_policy_mfa"
	policyMfaDefault           = "okta_policy_mfa_default"
	policyPassword             = "okta_policy_password"
	policyPasswordDefault      = "okta_policy_password_default"
	policyRuleIdpDiscovery     = "okta_policy_rule_idp_discovery"
	policyRuleMfa              = "okta_policy_rule_mfa"
	policyRulePassword         = "okta_policy_rule_password"
	policyRuleSignOn           = "okta_policy_rule_signon"
	policySignOn               = "okta_policy_signon"
	securityNotificationEmails = "okta_security_notification_emails"
	templateEmail              = "okta_template_email"
	templateSms                = "okta_template_sms"
	threatInsightSettings      = "okta_threat_insight_settings"
	trustedOrigin              = "okta_trusted_origin"
	userAdminRoles             = "okta_user_admin_roles"
	userBaseSchemaProperty     = "okta_user_base_schema_property"
	userFactorQuestion         = "okta_user_factor_question"
	userGroupMemberships       = "okta_user_group_memberships"
	userSecurityQuestions      = "okta_user_security_questions"
	userSchemaProperty         = "okta_user_schema_property"
	userType                   = "okta_user_type"
)

// Provider establishes a client connection to an okta site
// determined by its schema string values
func Provider() *schema.Provider {
	deprecatedPolicies := dataSourceDefaultPolicies()
	deprecatedPolicies.DeprecationMessage = "This data source will be deprecated in favor of okta_default_policy or okta_policy data sources."
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"org_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OKTA_ORG_NAME", nil),
				Description: "The organization to manage in Okta.",
			},
			"api_token": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("OKTA_API_TOKEN", nil),
				Description:   "API Token granting privileges to Okta API.",
				ConflictsWith: []string{"client_id", "scopes", "private_key"},
			},
			"client_id": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc("OKTA_API_CLIENT_ID", nil),
				Description:   "API Token granting privileges to Okta API.",
				ConflictsWith: []string{"api_token"},
			},
			"scopes": {
				Type:          schema.TypeSet,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				DefaultFunc:   envDefaultSetFunc("OKTA_API_SCOPES", nil),
				Description:   "API Token granting privileges to Okta API.",
				ConflictsWith: []string{"api_token"},
			},
			"private_key": {
				Optional:      true,
				Type:          schema.TypeString,
				DefaultFunc:   schema.EnvDefaultFunc("OKTA_API_PRIVATE_KEY", nil),
				Description:   "API Token granting privileges to Okta API.",
				ConflictsWith: []string{"api_token"},
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
			adminRoleTargets:           resourceAdminRoleTargets(),
			appAutoLogin:               resourceAppAutoLogin(),
			appBookmark:                resourceAppBookmark(),
			appBasicAuth:               resourceAppBasicAuth(),
			appGroupAssignment:         resourceAppGroupAssignment(),
			appGroupAssignments:        resourceAppGroupAssignments(),
			appUser:                    resourceAppUser(),
			appOAuth:                   resourceAppOAuth(),
			appOAuthAPIScope:           resourceAppOAuthAPIScope(),
			appOAuthRedirectURI:        resourceAppOAuthRedirectURI(),
			appSaml:                    resourceAppSaml(),
			appSamlAppSettings:         resourceAppSamlAppSettings(),
			appSecurePasswordStore:     resourceAppSecurePasswordStore(),
			appSignOnPolicyRule:        resourceAppSignOnPolicyRule(),
			appSwa:                     resourceAppSwa(),
			appSharedCredentials:       resourceAppSharedCredentials(),
			appThreeField:              resourceAppThreeField(),
			appUserSchemaProperty:      resourceAppUserSchemaProperty(),
			appUserBaseSchemaProperty:  resourceAppUserBaseSchemaProperty(),
			authServer:                 resourceAuthServer(),
			authServerDefault:          resourceAuthServerDefault(),
			authServerClaim:            resourceAuthServerClaim(),
			authServerClaimDefault:     resourceAuthServerClaimDefault(),
			authServerPolicy:           resourceAuthServerPolicy(),
			authServerPolicyRule:       resourceAuthServerPolicyRule(),
			authServerScope:            resourceAuthServerScope(),
			behavior:                   resourceBehavior(),
			domain:                     resourceDomain(),
			domainCertificate:          resourceDomainCertificate(),
			domainVerification:         resourceDomainVerification(),
			eventHook:                  resourceEventHook(),
			emailSender:                resourceEmailSender(),
			emailSenderVerification:    resourceEmailSenderVerification(),
			factor:                     resourceFactor(),
			factorTotp:                 resourceFactorTOTP(),
			groupRole:                  resourceGroupRole(),
			groupRoles:                 resourceGroupRoles(),
			groupRule:                  resourceGroupRule(),
			groupSchemaProperty:        resourceGroupCustomSchemaProperty(),
			idpOidc:                    resourceIdpOidc(),
			idpSaml:                    resourceIdpSaml(),
			idpSamlKey:                 resourceIdpSigningKey(),
			idpSocial:                  resourceIdpSocial(),
			inlineHook:                 resourceInlineHook(),
			networkZone:                resourceNetworkZone(),
			oktaGroup:                  resourceGroup(),
			oktaGroupMembership:        resourceGroupMembership(),
			oktaGroupMemberships:       resourceGroupMemberships(),
			oktaProfileMapping:         resourceOktaProfileMapping(),
			oktaUser:                   resourceUser(),
			policyMfa:                  resourcePolicyMfa(),
			policyMfaDefault:           resourcePolicyMfaDefault(),
			policyPassword:             resourcePolicyPassword(),
			policyPasswordDefault:      resourcePolicyPasswordDefault(),
			policySignOn:               resourcePolicySignOn(),
			policyRuleIdpDiscovery:     resourcePolicyRuleIdpDiscovery(),
			policyRuleMfa:              resourcePolicyMfaRule(),
			policyRulePassword:         resourcePolicyPasswordRule(),
			policyRuleSignOn:           resourcePolicySignOnRule(),
			securityNotificationEmails: resourceSecurityNotificationEmails(),
			templateEmail:              resourceTemplateEmail(),
			templateSms:                resourceTemplateSms(),
			threatInsightSettings:      resourceThreatInsightSettings(),
			trustedOrigin:              resourceTrustedOrigin(),
			userSchemaProperty:         resourceUserCustomSchemaProperty(),
			userBaseSchemaProperty:     resourceUserBaseSchemaProperty(),
			userType:                   resourceUserType(),
			userGroupMemberships:       resourceUserGroupMemberships(),
			userAdminRoles:             resourceUserAdminRoles(),
			userFactorQuestion:         resourceUserFactorQuestion(),

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
			"okta_app":                         dataSourceApp(),
			appGroupAssignments:                dataSourceAppGroupAssignments(),
			appSaml:                            dataSourceAppSaml(),
			appSignOnPolicy:                    dataSourceAppSignOnPolicy(),
			appOAuth:                           dataSourceAppOauth(),
			"okta_app_metadata_saml":           dataSourceAppMetadataSaml(),
			"okta_app_user_assignments":        dataSourceAppUserAssignments(),
			"okta_default_policies":            deprecatedPolicies,
			"okta_default_policy":              dataSourceDefaultPolicies(),
			"okta_everyone_group":              dataSourceEveryoneGroup(),
			behavior:                           dataSourceBehavior(),
			behaviors:                          dataSourceBehaviors(),
			oktaGroup:                          dataSourceGroup(),
			oktaGroups:                         dataSourceGroups(),
			"okta_idp_metadata_saml":           dataSourceIdpMetadataSaml(),
			idpSaml:                            dataSourceIdpSaml(),
			idpOidc:                            dataSourceIdpOidc(),
			idpSocial:                          dataSourceIdpSocial(),
			"okta_policy":                      dataSourcePolicy(),
			authServerPolicy:                   dataSourceAuthServerPolicy(),
			"okta_user_profile_mapping_source": dataSourceUserProfileMappingSource(),
			oktaUser:                           dataSourceUser(),
			"okta_users":                       dataSourceUsers(),
			authServer:                         dataSourceAuthServer(),
			"okta_auth_server_scopes":          dataSourceAuthServerScopes(),
			userType:                           dataSourceUserType(),
			userSecurityQuestions:              dataSourceUserSecurityQuestions(),
			networkZone:                        dataSourceNetworkZone(),
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
		clientID:       d.Get("client_id").(string),
		privateKey:     d.Get("private_key").(string),
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
	if v := os.Getenv("OKTA_API_SCOPES"); v != "" && len(config.scopes) == 0 {
		config.scopes = strings.Split(v, ",")
	}
	if err := config.loadAndValidate(ctx); err != nil {
		return nil, diag.Errorf("[ERROR] invalid configuration: %v", err)
	}
	return &config, nil
}

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
