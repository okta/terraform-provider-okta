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
	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/okta/internal/mutexkv"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/sdk"
)

// TODO some of these functions can be extracted to a utils or validators package, etc.

// oktaMutexKV is a global MutexKV for use within this plugin
var oktaMutexKV = mutexkv.NewMutexKV()

func listApps(ctx context.Context, config *config.Config, filters *AppFilters, limit int64) ([]oktav5sdk.ListApplications200ResponseInner, error) {
	req := config.OktaSDKClientV5.ApplicationAPI.ListApplications(ctx).Limit(int32(limit))
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

func GetOktaClientFromMetadata(meta interface{}) *sdk.Client {
	return meta.(*config.Config).OktaSDKClientV2
}

func GetOktaV3ClientFromMetadata(meta interface{}) *okta.APIClient {
	return meta.(*config.Config).OktaSDKClientV3
}

func GetOktaV5ClientFromMetadata(meta interface{}) *oktav5sdk.APIClient {
	return meta.(*config.Config).OktaSDKClientV5
}

func GetAPISupplementFromMetadata(meta interface{}) *sdk.APISupplement {
	return meta.(*config.Config).OktaSDKsupplementClient
}

func Logger(meta interface{}) hclog.Logger {
	return meta.(*config.Config).Logger
}

func GetRequestExecutor(m interface{}) *sdk.RequestExecutor {
	return GetOktaClientFromMetadata(m).GetRequestExecutor()
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
		NewAppAccessPolicyAssignmentResource,
		NewAppOAuthRoleAssignmentResource,
		NewTrustedServerResource,
		NewBrandResource,
		NewLogStreamResource,
		NewPolicyDeviceAssuranceAndroidResource,
		NewPolicyDeviceAssuranceChromeOSResource,
		NewPolicyDeviceAssuranceIOSResource,
		NewPolicyDeviceAssuranceMacOSResource,
		NewPolicyDeviceAssuranceWindowsResource,
		NewCustomizedSigninResource,
		NewPreviewSigninResource,
		NewGroupOwnerResource,
		NewAppSignOnPolicyResource,
		NewEmailTemplateSettingsResource,
	}
}

func FWProviderDataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewOrgMetadataDataSource,
		NewDefaultSigninPageDataSource,
		NewLogStreamDataSource,
		NewAppsDataSource,
		NewUserTypeDataSource,
		NewDeviceAssurancePolicyDataSource,
	}
}

func ProviderResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		resources.OktaIDaaSAdminRoleCustom:               ResourceAdminRoleCustom(),
		resources.OktaIDaaSAdminRoleCustomAssignments:    ResourceAdminRoleCustomAssignments(),
		resources.OktaIDaaSAdminRoleTargets:              ResourceAdminRoleTargets(),
		resources.OktaIDaaSAppAutoLogin:                  ResourceAppAutoLogin(),
		resources.OktaIDaaSAppBasicAuth:                  ResourceAppBasicAuth(),
		resources.OktaIDaaSAppBookmark:                   ResourceAppBookmark(),
		resources.OktaIDaaSAppGroupAssignment:            ResourceAppGroupAssignment(),
		resources.OktaIDaaSAppGroupAssignments:           ResourceAppGroupAssignments(),
		resources.OktaIDaaSAppOAuth:                      ResourceAppOAuth(),
		resources.OktaIDaaSAppOAuthAPIScope:              ResourceAppOAuthAPIScope(),
		resources.OktaIDaaSAppOAuthPostLogoutRedirectURI: ResourceAppOAuthPostLogoutRedirectURI(),
		resources.OktaIDaaSAppOAuthRedirectURI:           ResourceAppOAuthRedirectURI(),
		resources.OktaIDaaSAppSaml:                       ResourceAppSaml(),
		resources.OktaIDaaSAppSamlAppSettings:            ResourceAppSamlAppSettings(),
		resources.OktaIDaaSAppSecurePasswordStore:        ResourceAppSecurePasswordStore(),
		resources.OktaIDaaSAppSharedCredentials:          ResourceAppSharedCredentials(),
		resources.OktaIDaaSAppSignOnPolicyRule:           ResourceAppSignOnPolicyRule(),
		resources.OktaIDaaSAppSwa:                        ResourceAppSwa(),
		resources.OktaIDaaSAppThreeField:                 ResourceAppThreeField(),
		resources.OktaIDaaSAppUser:                       ResourceAppUser(),
		resources.OktaIDaaSAppUserBaseSchemaProperty:     ResourceAppUserBaseSchemaProperty(),
		resources.OktaIDaaSAppUserSchemaProperty:         ResourceAppUserSchemaProperty(),
		resources.OktaIDaaSAuthenticator:                 ResourceAuthenticator(),
		resources.OktaIDaaSAuthServer:                    ResourceAuthServer(),
		resources.OktaIDaaSAuthServerClaim:               ResourceAuthServerClaim(),
		resources.OktaIDaaSAuthServerClaimDefault:        ResourceAuthServerClaimDefault(),
		resources.OktaIDaaSAuthServerDefault:             ResourceAuthServerDefault(),
		resources.OktaIDaaSAuthServerPolicy:              ResourceAuthServerPolicy(),
		resources.OktaIDaaSAuthServerPolicyRule:          ResourceAuthServerPolicyRule(),
		resources.OktaIDaaSAuthServerScope:               ResourceAuthServerScope(),
		resources.OktaIDaaSBehavior:                      ResourceBehavior(),
		resources.OktaIDaaSCaptcha:                       ResourceCaptcha(),
		resources.OktaIDaaSCaptchaOrgWideSettings:        ResourceCaptchaOrgWideSettings(),
		resources.OktaIDaaSDomain:                        ResourceDomain(),
		resources.OktaIDaaSDomainCertificate:             ResourceDomainCertificate(),
		resources.OktaIDaaSDomainVerification:            ResourceDomainVerification(),
		resources.OktaIDaaSEmailCustomization:            ResourceEmailCustomization(),
		resources.OktaIDaaSEmailDomain:                   ResourceEmailDomain(),
		resources.OktaIDaaSEmailDomainVerification:       ResourceEmailDomainVerification(),
		resources.OktaIDaaSEmailSender:                   ResourceEmailSender(),
		resources.OktaIDaaSEmailSenderVerification:       ResourceEmailSenderVerification(),
		resources.OktaIDaaSEventHook:                     ResourceEventHook(),
		resources.OktaIDaaSEventHookVerification:         ResourceEventHookVerification(),
		resources.OktaIDaaSFactor:                        ResourceFactor(),
		resources.OktaIDaaSFactorTotp:                    ResourceFactorTOTP(),
		resources.OktaIDaaSGroup:                         ResourceGroup(),
		resources.OktaIDaaSGroupMemberships:              ResourceGroupMemberships(),
		resources.OktaIDaaSGroupRole:                     ResourceGroupRole(),
		resources.OktaIDaaSGroupRule:                     ResourceGroupRule(),
		resources.OktaIDaaSGroupSchemaProperty:           ResourceGroupCustomSchemaProperty(),
		resources.OktaIDaaSIdpOidc:                       ResourceIdpOidc(),
		resources.OktaIDaaSIdpSaml:                       ResourceIdpSaml(),
		resources.OktaIDaaSIdpSamlKey:                    ResourceIdpSigningKey(),
		resources.OktaIDaaSIdpSocial:                     ResourceIdpSocial(),
		resources.OktaIDaaSInlineHook:                    ResourceInlineHook(),
		resources.OktaIDaaSLinkDefinition:                ResourceLinkDefinition(),
		resources.OktaIDaaSLinkValue:                     ResourceLinkValue(),
		resources.OktaIDaaSNetworkZone:                   ResourceNetworkZone(),
		resources.OktaIDaaSOrgConfiguration:              ResourceOrgConfiguration(),
		resources.OktaIDaaSOrgSupport:                    ResourceOrgSupport(),
		resources.OktaIDaaSPolicyMfa:                     ResourcePolicyMfa(),
		resources.OktaIDaaSPolicyMfaDefault:              ResourcePolicyMfaDefault(),
		resources.OktaIDaaSPolicyPassword:                ResourcePolicyPassword(),
		resources.OktaIDaaSPolicyPasswordDefault:         ResourcePolicyPasswordDefault(),
		resources.OktaIDaaSPolicyProfileEnrollment:       ResourcePolicyProfileEnrollment(),
		resources.OktaIDaaSPolicyProfileEnrollmentApps:   ResourcePolicyProfileEnrollmentApps(),
		resources.OktaIDaaSPolicyRuleIdpDiscovery:        ResourcePolicyRuleIdpDiscovery(),
		resources.OktaIDaaSPolicyRuleMfa:                 ResourcePolicyMfaRule(),
		resources.OktaIDaaSPolicyRulePassword:            ResourcePolicyPasswordRule(),
		resources.OktaIDaaSPolicyRuleProfileEnrollment:   ResourcePolicyProfileEnrollmentRule(),
		resources.OktaIDaaSPolicyRuleSignOn:              ResourcePolicySignOnRule(),
		resources.OktaIDaaSPolicySignOn:                  ResourcePolicySignOn(),
		resources.OktaIDaaSProfileMapping:                ResourceProfileMapping(),
		resources.OktaIDaaSRateLimiting:                  ResourceRateLimiting(),
		resources.OktaIDaaSResourceSet:                   ResourceResourceSet(),
		resources.OktaIDaaSRoleSubscription:              ResourceRoleSubscription(),
		resources.OktaIDaaSSecurityNotificationEmails:    ResourceSecurityNotificationEmails(),
		resources.OktaIDaaSTemplateSms:                   ResourceTemplateSms(),
		resources.OktaIDaaSTheme:                         ResourceTheme(),
		resources.OktaIDaaSThreatInsightSettings:         ResourceThreatInsightSettings(),
		resources.OktaIDaaSTrustedOrigin:                 ResourceTrustedOrigin(),
		resources.OktaIDaaSUser:                          ResourceUser(),
		resources.OktaIDaaSUserAdminRoles:                ResourceUserAdminRoles(),
		resources.OktaIDaaSUserBaseSchemaProperty:        ResourceUserBaseSchemaProperty(),
		resources.OktaIDaaSUserFactorQuestion:            ResourceUserFactorQuestion(),
		resources.OktaIDaaSUserGroupMemberships:          ResourceUserGroupMemberships(),
		resources.OktaIDaaSUserSchemaProperty:            ResourceUserCustomSchemaProperty(),
		resources.OktaIDaaSUserType:                      ResourceUserType(),
	}
}

func ProviderDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		resources.OktaIDaaSApp:                      DataSourceApp(),
		resources.OktaIDaaSAppGroupAssignments:      DataSourceAppGroupAssignments(),
		resources.OktaIDaaSAppMetadataSaml:          DataSourceAppMetadataSaml(),
		resources.OktaIDaaSAppOAuth:                 DataSourceAppOauth(),
		resources.OktaIDaaSAppSaml:                  DataSourceAppSaml(),
		resources.OktaIDaaSAppSignOnPolicy:          DataSourceAppSignOnPolicy(),
		resources.OktaIDaaSAppUserAssignments:       DataSourceAppUserAssignments(),
		resources.OktaIDaaSAuthenticator:            DataSourceAuthenticator(),
		resources.OktaIDaaSAuthServer:               DataSourceAuthServer(),
		resources.OktaIDaaSAuthServerClaim:          DataSourceAuthServerClaim(),
		resources.OktaIDaaSAuthServerClaims:         DataSourceAuthServerClaims(),
		resources.OktaIDaaSAuthServerPolicy:         DataSourceAuthServerPolicy(),
		resources.OktaIDaaSAuthServerScopes:         DataSourceAuthServerScopes(),
		resources.OktaIDaaSBehavior:                 DataSourceBehavior(),
		resources.OktaIDaaSBehaviors:                DataSourceBehaviors(),
		resources.OktaIDaaSBrand:                    DataSourceBrand(),
		resources.OktaIDaaSBrands:                   DataSourceBrands(),
		resources.OktaIDaaSDomain:                   DataSourceDomain(),
		resources.OktaIDaaSEmailCustomization:       DataSourceEmailCustomization(),
		resources.OktaIDaaSEmailCustomizations:      DataSourceEmailCustomizations(),
		resources.OktaIDaaSEmailTemplate:            DataSourceEmailTemplate(),
		resources.OktaIDaaSEmailTemplates:           DataSourceEmailTemplates(),
		resources.OktaIDaaSDefaultPolicy:            DataSourceDefaultPolicy(),
		resources.OktaIDaaSGroup:                    DataSourceGroup(),
		resources.OktaIDaaSGroupEveryone:            DataSourceEveryoneGroup(),
		resources.OktaIDaaSGroupRule:                DataSourceGroupRule(),
		resources.OktaIDaaSGroups:                   DataSourceGroups(),
		resources.OktaIDaaSIdpMetadataSaml:          DataSourceIdpMetadataSaml(),
		resources.OktaIDaaSIdpOidc:                  DataSourceIdpOidc(),
		resources.OktaIDaaSIdpSaml:                  DataSourceIdpSaml(),
		resources.OktaIDaaSIdpSocial:                DataSourceIdpSocial(),
		resources.OktaIDaaSNetworkZone:              DataSourceNetworkZone(),
		resources.OktaIDaaSPolicy:                   DataSourcePolicy(),
		resources.OktaIDaaSRoleSubscription:         DataSourceRoleSubscription(),
		resources.OktaIDaaSTheme:                    DataSourceTheme(),
		resources.OktaIDaaSThemes:                   DataSourceThemes(),
		resources.OktaIDaaSTrustedOrigins:           DataSourceTrustedOrigins(),
		resources.OktaIDaaSUser:                     DataSourceUser(),
		resources.OktaIDaaSUserProfileMappingSource: DataSourceUserProfileMappingSource(),
		resources.OktaIDaaSUsers:                    DataSourceUsers(),
		resources.OktaIDaaSUserSecurityQuestions:    DataSourceUserSecurityQuestions(),
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
