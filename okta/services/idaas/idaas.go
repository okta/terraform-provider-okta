package idaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	fwdiag "github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/okta/okta-sdk-golang/v4/okta"
	oktav5sdk "github.com/okta/okta-sdk-golang/v5/okta"
	v6okta "github.com/okta/okta-sdk-golang/v6/okta"
	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/okta/internal/mutexkv"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/sdk"
)

// oktaMutexKV is a global MutexKV for use within this plugin
var oktaMutexKV = mutexkv.NewMutexKV()

func listAppsV5(ctx context.Context, config *config.Config, filters *AppFilters, limit int64) ([]oktav5sdk.ListApplications200ResponseInner, error) {
	req := config.OktaIDaaSClient.OktaSDKClientV5().ApplicationAPI.ListApplications(ctx).Limit(int32(limit))
	if filters != nil {
		req = req.Filter(filters.Status)
		req = req.Q(filters.GetQ())
	}
	apps, resp, err := req.Execute()
	if err != nil {
		return nil, err
	}
	for resp.HasNextPage() {
		var nextApps []oktav5sdk.ListApplications200ResponseInner
		resp, err = resp.Next(&nextApps)
		if err != nil {
			return nil, err
		}
		apps = append(apps, nextApps...)
	}
	return apps, nil
}

func getOktaClientFromMetadata(meta interface{}) *sdk.Client {
	return meta.(*config.Config).OktaIDaaSClient.OktaSDKClientV2()
}

func getOktaV3ClientFromMetadata(meta interface{}) *okta.APIClient {
	return meta.(*config.Config).OktaIDaaSClient.OktaSDKClientV3()
}

func getOktaV5ClientFromMetadata(meta interface{}) *oktav5sdk.APIClient {
	return meta.(*config.Config).OktaIDaaSClient.OktaSDKClientV5()
}

func getOktaV6ClientFromMetadata(meta interface{}) *v6okta.APIClient {
	return meta.(*config.Config).OktaIDaaSClient.OktaSDKClientV6()
}

func getAPISupplementFromMetadata(meta interface{}) *sdk.APISupplement {
	return meta.(*config.Config).OktaIDaaSClient.OktaSDKSupplementClient()
}

func logger(meta interface{}) hclog.Logger {
	return meta.(*config.Config).Logger
}

func getRequestExecutor(m interface{}) *sdk.RequestExecutor {
	return getOktaClientFromMetadata(m).GetRequestExecutor()
}

func fwproviderIsClassicOrg(ctx context.Context, config *config.Config) bool {
	return config.IsClassicOrg(ctx)
}

func providerIsClassicOrg(ctx context.Context, m interface{}) bool {
	if config, ok := m.(*config.Config); ok && config.IsClassicOrg(ctx) {
		return true
	}
	return false
}

func FWProviderResources() []func() resource.Resource {
	return []func() resource.Resource{
		newAppAccessPolicyAssignmentResource,
		newAppOAuthRoleAssignmentResource,
		newTrustedServerResource,
		newBrandResource,
		newLogStreamResource,
		newPolicyDeviceAssuranceAndroidResource,
		newPolicyDeviceAssuranceChromeOSResource,
		newPolicyDeviceAssuranceIOSResource,
		newPolicyDeviceAssuranceMacOSResource,
		newPolicyDeviceAssuranceWindowsResource,
		newCustomizedSigninResource,
		newPreviewSigninResource,
		newGroupOwnerResource,
		newAppSignOnPolicyResource,
		newEmailTemplateSettingsResource,
		newFeaturesResource,
		newRealmResource,
		newRealmAssignmentResource,
		newRateLimitResource,
		newRateLimitAdminNotificationSettingsResource,
		newRateLimitWarningThresholdPercentageResource,
		newPrincipalRateLimitsResource,
		newSecurityEventsProviderResource,
		newDevicesResource,
		newAppFeaturesResource,
		newPushProvidersResource,
		newHookKeyResource,
		newAPIServiceIntegrationResource,
		newAPITokenResource,
		newAppTokenResource,
		newAppConnectionsResource,
		newAgentPoolUpdateResource,
		newUISchemaResource,
	}
}

func FWProviderDataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		newAuthServerClientsDataSource,
		newAuthServerKeysDataSource,
		newOrgMetadataDataSource,
		newDefaultSigninPageDataSource,
		newLogStreamDataSource,
		newAppsDataSource,
		newUserTypeDataSource,
		newDeviceAssurancePolicyDataSource,
		newFeaturesDataSource,
		newRealmDataSource,
		newRealmAssignmentDataSource,
		newRateLimitAdminNotificationSettingsDataSource,
		newRateLimitWarningThresholdPercentageDataSource,
		newPrincipalRateLimitsDataSource,
		newSecurityEventsProviderDataSource,
		newDeviceDataSource,
		newAppFeaturesDataSource,
		newPushProviderDataSource,
		newHookKeyDataSource,
		newAPIServiceIntegrationDataSource,
		newAPITokenDataSource,
		newAppTokenDataSource,
		newAppConnectionsDataSource,
		newAgentPoolDataSource,
		newUISchemaDataSource,
	}
}

func ProviderResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		resources.OktaIDaaSAdminRoleCustom:               resourceAdminRoleCustom(),
		resources.OktaIDaaSAdminRoleCustomAssignments:    resourceAdminRoleCustomAssignments(),
		resources.OktaIDaaSAdminRoleTargets:              resourceAdminRoleTargets(),
		resources.OktaIDaaSAppAutoLogin:                  resourceAppAutoLogin(),
		resources.OktaIDaaSAppBasicAuth:                  resourceAppBasicAuth(),
		resources.OktaIDaaSAppBookmark:                   resourceAppBookmark(),
		resources.OktaIDaaSAppGroupAssignment:            resourceAppGroupAssignment(),
		resources.OktaIDaaSAppGroupAssignments:           resourceAppGroupAssignments(),
		resources.OktaIDaaSAppOAuth:                      resourceAppOAuth(),
		resources.OktaIDaaSAppOAuthAPIScope:              resourceAppOAuthAPIScope(),
		resources.OktaIDaaSAppOAuthPostLogoutRedirectURI: resourceAppOAuthPostLogoutRedirectURI(),
		resources.OktaIDaaSAppOAuthRedirectURI:           resourceAppOAuthRedirectURI(),
		resources.OktaIDaaSAppSaml:                       resourceAppSaml(),
		resources.OktaIDaaSAppSamlAppSettings:            resourceAppSamlAppSettings(),
		resources.OktaIDaaSAppSecurePasswordStore:        resourceAppSecurePasswordStore(),
		resources.OktaIDaaSAppSharedCredentials:          resourceAppSharedCredentials(),
		resources.OktaIDaaSAppSignOnPolicyRule:           resourceAppSignOnPolicyRule(),
		resources.OktaIDaaSAppSwa:                        resourceAppSwa(),
		resources.OktaIDaaSAppThreeField:                 resourceAppThreeField(),
		resources.OktaIDaaSAppUser:                       resourceAppUser(),
		resources.OktaIDaaSAppUserBaseSchemaProperty:     resourceAppUserBaseSchemaProperty(),
		resources.OktaIDaaSAppUserSchemaProperty:         resourceAppUserSchemaProperty(),
		resources.OktaIDaaSAuthenticator:                 resourceAuthenticator(),
		resources.OktaIDaaSAuthServer:                    resourceAuthServer(),
		resources.OktaIDaaSAuthServerClaim:               resourceAuthServerClaim(),
		resources.OktaIDaaSAuthServerClaimDefault:        resourceAuthServerClaimDefault(),
		resources.OktaIDaaSAuthServerDefault:             resourceAuthServerDefault(),
		resources.OktaIDaaSAuthServerPolicy:              resourceAuthServerPolicy(),
		resources.OktaIDaaSAuthServerPolicyRule:          resourceAuthServerPolicyRule(),
		resources.OktaIDaaSAuthServerScope:               resourceAuthServerScope(),
		resources.OktaIDaaSBehavior:                      resourceBehavior(),
		resources.OktaIDaaSCaptcha:                       resourceCaptcha(),
		resources.OktaIDaaSCaptchaOrgWideSettings:        resourceCaptchaOrgWideSettings(),
		resources.OktaIDaaSDomain:                        resourceDomain(),
		resources.OktaIDaaSDomainCertificate:             resourceDomainCertificate(),
		resources.OktaIDaaSDomainVerification:            resourceDomainVerification(),
		resources.OktaIDaaSEmailCustomization:            resourceEmailCustomization(),
		resources.OktaIDaaSEmailDomain:                   resourceEmailDomain(),
		resources.OktaIDaaSEmailDomainVerification:       resourceEmailDomainVerification(),
		resources.OktaIDaaSEmailSender:                   resourceEmailSender(),
		resources.OktaIDaaSEmailSenderVerification:       resourceEmailSenderVerification(),
		resources.OktaIDaaSEmailSMTPServer:               resourceEmailSMTP(),
		resources.OktaIDaaSEventHook:                     resourceEventHook(),
		resources.OktaIDaaSEventHookVerification:         resourceEventHookVerification(),
		resources.OktaIDaaSFactor:                        resourceFactor(),
		resources.OktaIDaaSFactorTotp:                    resourceFactorTOTP(),
		resources.OktaIDaaSGroup:                         resourceGroup(),
		resources.OktaIDaaSGroupMemberships:              resourceGroupMemberships(),
		resources.OktaIDaaSGroupRole:                     resourceGroupRole(),
		resources.OktaIDaaSGroupRule:                     resourceGroupRule(),
		resources.OktaIDaaSGroupSchemaProperty:           resourceGroupCustomSchemaProperty(),
		resources.OktaIDaaSIdpOidc:                       resourceIdpOidc(),
		resources.OktaIDaaSIdpSaml:                       resourceIdpSaml(),
		resources.OktaIDaaSIdpSamlKey:                    resourceIdpSigningKey(),
		resources.OktaIDaaSIdpSocial:                     resourceIdpSocial(),
		resources.OktaIDaaSInlineHook:                    resourceInlineHook(),
		resources.OktaIDaaSLinkDefinition:                resourceLinkDefinition(),
		resources.OktaIDaaSLinkValue:                     resourceLinkValue(),
		resources.OktaIDaaSNetworkZone:                   resourceNetworkZone(),
		resources.OktaIDaaSOrgConfiguration:              resourceOrgConfiguration(),
		resources.OktaIDaaSOrgSupport:                    resourceOrgSupport(),
		resources.OktaIDaaSPolicyMfa:                     resourcePolicyMfa(),
		resources.OktaIDaaSPolicyMfaDefault:              resourcePolicyMfaDefault(),
		resources.OktaIDaaSPolicyPassword:                resourcePolicyPassword(),
		resources.OktaIDaaSPolicyPasswordDefault:         resourcePolicyPasswordDefault(),
		resources.OktaIDaaSPolicyProfileEnrollment:       resourcePolicyProfileEnrollment(),
		resources.OktaIDaaSPolicyProfileEnrollmentApps:   resourcePolicyProfileEnrollmentApps(),
		resources.OktaIDaaSPolicyRuleIdpDiscovery:        resourcePolicyRuleIdpDiscovery(),
		resources.OktaIDaaSPolicyRuleMfa:                 resourcePolicyMfaRule(),
		resources.OktaIDaaSPolicyRulePassword:            resourcePolicyPasswordRule(),
		resources.OktaIDaaSPolicyRuleProfileEnrollment:   resourcePolicyProfileEnrollmentRule(),
		resources.OktaIDaaSPolicyRuleSignOn:              resourcePolicySignOnRule(),
		resources.OktaIDaaSPolicySignOn:                  resourcePolicySignOn(),
		resources.OktaIDaaSProfileMapping:                resourceProfileMapping(),
		//resources.OktaIDaaSRateLimiting:                  resourceRateLimiting(),
		resources.OktaIDaaSResourceSet:                resourceResourceSet(),
		resources.OktaIDaaSRoleSubscription:           resourceRoleSubscription(),
		resources.OktaIDaaSSecurityNotificationEmails: resourceSecurityNotificationEmails(),
		resources.OktaIDaaSTemplateSms:                resourceTemplateSms(),
		resources.OktaIDaaSTheme:                      resourceTheme(),
		resources.OktaIDaaSThreatInsightSettings:      resourceThreatInsightSettings(),
		resources.OktaIDaaSTrustedOrigin:              resourceTrustedOrigin(),
		resources.OktaIDaaSUser:                       resourceUser(),
		resources.OktaIDaaSUserAdminRoles:             resourceUserAdminRoles(),
		resources.OktaIDaaSUserBaseSchemaProperty:     resourceUserBaseSchemaProperty(),
		resources.OktaIDaaSUserFactorQuestion:         resourceUserFactorQuestion(),
		resources.OktaIDaaSUserGroupMemberships:       resourceUserGroupMemberships(),
		resources.OktaIDaaSUserSchemaProperty:         resourceUserCustomSchemaProperty(),
		resources.OktaIDaaSUserType:                   resourceUserType(),
	}
}

func ProviderDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		resources.OktaIDaaSApp:                      dataSourceApp(),
		resources.OktaIDaaSAppGroupAssignments:      dataSourceAppGroupAssignments(),
		resources.OktaIDaaSAppMetadataSaml:          dataSourceAppMetadataSaml(),
		resources.OktaIDaaSAppOAuth:                 dataSourceAppOauth(),
		resources.OktaIDaaSAppSaml:                  dataSourceAppSaml(),
		resources.OktaIDaaSAppSignOnPolicy:          dataSourceAppSignOnPolicy(),
		resources.OktaIDaaSAppUserAssignments:       dataSourceAppUserAssignments(),
		resources.OktaIDaaSAuthenticator:            dataSourceAuthenticator(),
		resources.OktaIDaaSAuthServer:               dataSourceAuthServer(),
		resources.OktaIDaaSAuthServerClaim:          dataSourceAuthServerClaim(),
		resources.OktaIDaaSAuthServerClaims:         dataSourceAuthServerClaims(),
		resources.OktaIDaaSAuthServerPolicy:         dataSourceAuthServerPolicy(),
		resources.OktaIDaaSAuthServerScopes:         dataSourceAuthServerScopes(),
		resources.OktaIDaaSBehavior:                 dataSourceBehavior(),
		resources.OktaIDaaSBehaviors:                dataSourceBehaviors(),
		resources.OktaIDaaSBrand:                    dataSourceBrand(),
		resources.OktaIDaaSBrands:                   dataSourceBrands(),
		resources.OktaIDaaSDomain:                   dataSourceDomain(),
		resources.OktaIDaaSEmailCustomization:       dataSourceEmailCustomization(),
		resources.OktaIDaaSEmailCustomizations:      dataSourceEmailCustomizations(),
		resources.OktaIDaaSEmailSMTPServer:          dataSourceEmailSMTPServers(),
		resources.OktaIDaaSEmailTemplate:            dataSourceEmailTemplate(),
		resources.OktaIDaaSEmailTemplates:           dataSourceEmailTemplates(),
		resources.OktaIDaaSDefaultPolicy:            dataSourceDefaultPolicy(),
		resources.OktaIDaaSGroup:                    dataSourceGroup(),
		resources.OktaIDaaSGroupEveryone:            dataSourceEveryoneGroup(),
		resources.OktaIDaaSGroupRule:                dataSourceGroupRule(),
		resources.OktaIDaaSGroups:                   dataSourceGroups(),
		resources.OktaIDaaSIdpMetadataSaml:          dataSourceIdpMetadataSaml(),
		resources.OktaIDaaSIdpOidc:                  dataSourceIdpOidc(),
		resources.OktaIDaaSIdpSaml:                  dataSourceIdpSaml(),
		resources.OktaIDaaSIdpSocial:                dataSourceIdpSocial(),
		resources.OktaIDaaSNetworkZone:              dataSourceNetworkZone(),
		resources.OktaIDaaSPolicy:                   dataSourcePolicy(),
		resources.OktaIDaaSRoleSubscription:         dataSourceRoleSubscription(),
		resources.OktaIDaaSTheme:                    dataSourceTheme(),
		resources.OktaIDaaSThemes:                   dataSourceThemes(),
		resources.OktaIDaaSTrustedOrigins:           dataSourceTrustedOrigins(),
		resources.OktaIDaaSUser:                     dataSourceUser(),
		resources.OktaIDaaSUserProfileMappingSource: dataSourceUserProfileMappingSource(),
		resources.OktaIDaaSUsers:                    dataSourceUsers(),
		resources.OktaIDaaSUserSecurityQuestions:    dataSourceUserSecurityQuestions(),
	}
}

func stringIsJSON(i interface{}, k cty.Path) diag.Diagnostics {
	v, ok := i.(string)
	if !ok {
		return diag.Errorf("expected type of %s to be string", k)
	}
	if v == "" {
		return diag.Errorf("expected %q JSON to not be empty, got %v", k, i)
	}
	if _, err := structure.NormalizeJsonString(v); err != nil {
		return diag.Errorf("%q contains an invalid JSON: %s", k, err)
	}
	return nil
}

// doNotRetry helper function to flag if provider should be using backoff.Retry
func doNotRetry(m interface{}, err error) bool {
	return m.(*config.Config).TimeOperations.DoNotRetry(err)
}

func datasourceOIEOnlyFeatureError(name string) diag.Diagnostics {
	return oieOnlyFeatureError("data-sources", name)
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

func dataSourceConfiguration(req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) *config.Config {
	if req.ProviderData == nil {
		return nil
	}

	config, ok := req.ProviderData.(*config.Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *config.Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return nil
	}

	return config
}

func resourceConfiguration(req resource.ConfigureRequest, resp *resource.ConfigureResponse) *config.Config {
	if req.ProviderData == nil {
		return nil
	}

	p, ok := req.ProviderData.(*config.Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *config.Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return nil
	}

	return p
}

func frameworkResourceOIEOnlyFeatureError(name string) fwdiag.Diagnostics {
	return frameworkOIEOnlyFeatureError("resources", name)
}

func frameworkOIEOnlyFeatureError(kind, name string) fwdiag.Diagnostics {
	url := fmt.Sprintf("https://registry.terraform.io/providers/okta/okta/latest/docs/%s/%s", kind, string(name[5:]))
	if kind == "resources" {
		kind = "resource"
	}
	if kind == "data-sources" {
		kind = "datasource"
	}
	var diags fwdiag.Diagnostics
	diags.AddError(fmt.Sprintf("%q is a %s for OIE Orgs only", name, kind), fmt.Sprintf(", see %s", url))
	return diags
}
