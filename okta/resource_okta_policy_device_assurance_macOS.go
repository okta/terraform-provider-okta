package okta

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/okta-sdk-golang/v3/okta"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource = &policyDeviceAssuranceMacOSResource{}
	// _ resource.ResourceWithConfigure   = &policyDeviceAssuranceResource{}
	// _ resource.ResourceWithImportState = &policyDeviceAssuranceResource{}
)

func NewPolicyDeviceAssuranceMacOSResource() resource.Resource {
	return &policyDeviceAssuranceMacOSResource{}
}

type policyDeviceAssuranceMacOSResource struct {
	v3Client *okta.APIClient
}

type policyDeviceAssuranceMacOSResourceModel struct {
	ID                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	Platform              types.String `tfsdk:"platform"`
	DiskEncryptionType    types.List   `tfsdk:"disk_encryption_type"`
	JailBreak             types.Bool   `tfsdk:"jailbreak"`
	OsVersion             types.String `tfsdk:"os_version"`
	SecureHardwarePresent types.Bool   `tfsdk:"secure_hardware_present"`
	ScreenLockType        types.List   `tfsdk:"screenlock_type"`
	// TODU
	ThirdPartySignalProviders thirdPartySignalProviders `tfsdk:"third_party_signal_providers"`
	CreateDate                types.String              `tfsdk:"created_date"`
	CreateBy                  types.String              `tfsdk:"created_by"`
	LastUpdate                types.String              `tfsdk:"last_update"`
	LastUpdatedBy             types.String              `tfsdk:"last_updated_by"`
}

type thirdPartySignalProviders struct {
	AllowScreenLock                   types.Bool   `tfsdk:"allow_screen_lock"`
	BrowserVersion                    types.String `tfsdk:"browser_version, omitempty"`
	BuiltInDNSClientEnabled           types.Bool   `tfsdk:"builtin_dns_client_enabled"`
	ChromeRemoteDesktopAppBlocked     types.Bool   `tfsdk:"chrome_remote_desktop_app_blocked"`
	CrowdStrikeAgentID                types.String `tfsdk:"crowd_strike_agent_id"`
	CrowdStrikeCustomerID             types.String `tfsdk:"crowd_strike_customer_id"`
	DeviceEnrollementDomain           types.String `tfsdk:"device_enrollement_domain"`
	DiskEncrypted                     types.Bool   `tfsdk:"disk_encrypted"`
	KeyTrustLevel                     types.String `tfsdk:"key_trust_level"`
	OsFirewall                        types.Bool   `tfsdk:"os_firewall"`
	OsVersion                         types.String `tfsdk:"os_version"`
	PasswordProctectionWarningTrigger types.String `tfsdk:"password_proctection_warning_trigger"`
	RealtimeURLCheckMode              types.Bool   `tfsdk:"realtime_url_check_mode"`
	SafeBrowsingProtectionLevel       types.String `tfsdk:"safe_browsing_protection_level"`
	ScreenLockSecured                 types.Bool   `tfsdk:"screen_lock_secured"`
	SecureBootEnabled                 types.Bool   `tfsdk:"secure_boot_enabled"`
	SiteIsolationEnabled              types.Bool   `tfsdk:"site_isolation_enabled"`
	ThirdPartyBlockingEnabled         types.Bool   `tfsdk:"third_party_blocking_enabled"`
	WindowMachineDomain               types.String `tfsdk:"window_machine_domain"`
	WindowUserDomain                  types.String `tfsdk:"window_user_domain"`
}

func (r *policyDeviceAssuranceMacOSResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_device_assurance_macOS"
}

// TODU different requirement for request and response?
// TODU validation
func (r *policyDeviceAssuranceMacOSResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages device assurance on policy",
		Attributes: map[string]schema.Attribute{
			// TODU needed?
			"id": schema.StringAttribute{
				Description: "Policy assurance id",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					// TODU
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Policy device assurance name",
				Required:    true,
			},
			"platform": schema.StringAttribute{
				Description: "Policy device assurance platform, can be ANDROID, CHROMEOS, IOS, MACOS or WINDOWS",
				Required:    true,
			},
			"disk_encryption_type": schema.ListAttribute{
				Description: "List of disk encryption type, can be ALL_INTERNAL_VOLUMES, FULL, or USER",
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.AtLeastOneOf(path.Expressions{
						path.MatchRoot("jail_break"),
						path.MatchRoot("os_version"),
						path.MatchRoot("secure_hardware_present"),
						path.MatchRoot("screenlock_type"),
					}...),
				},
			},
			"jail_break": schema.BoolAttribute{
				Description: "The device jailbreak. Only for android and iOS platform",
				Optional:    true,
				Validators: []validator.Bool{
					boolvalidator.AtLeastOneOf(path.Expressions{
						path.MatchRoot("disk_encryption_type"),
						path.MatchRoot("os_version"),
						path.MatchRoot("secure_hardware_present"),
						path.MatchRoot("screenlock_type"),
					}...),
				},
			},
			"os_version": schema.StringAttribute{
				Description: "The device os version",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.AtLeastOneOf(path.Expressions{
						path.MatchRoot("disk_encryption_type"),
						path.MatchRoot("jail_break"),
						path.MatchRoot("secure_hardware_present"),
						path.MatchRoot("screenlock_type"),
					}...),
				},
			},
			"secure_hardware_present": schema.BoolAttribute{
				Description: "Indicates if the device constains a secure hardware functionality",
				Optional:    true,
				Validators: []validator.Bool{
					boolvalidator.AtLeastOneOf(path.Expressions{
						path.MatchRoot("disk_encryption_type"),
						path.MatchRoot("jail_break"),
						path.MatchRoot("os_version"),
						path.MatchRoot("screenlock_type"),
					}...),
				},
			},
			"screenlock_type": schema.ListAttribute{
				Description: "List of screenlock type",
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.AtLeastOneOf(path.Expressions{
						path.MatchRoot("disk_encryption_type"),
						path.MatchRoot("jail_break"),
						path.MatchRoot("os_version"),
						path.MatchRoot("secure_hardware_present"),
					}...),
				},
			},
			"third_party_signal_providers": schema.ObjectAttribute{
				Description: "Settings for third-party signal providers. Required for ChromeOS platform, optional for others",
				Optional:    true,
				AttributeTypes: map[string]attr.Type{
					// TODU chromeOS only
					"allow_screen_lock":                 types.BoolType,
					"browser_version":                   types.StringType,
					"builtin_dns_client_enabled":        types.BoolType,
					"chrome_remote_desktop_app_blocked": types.BoolType,
					// TODU window only
					"crowd_strike_agent_id": types.StringType,
					// TODU window only
					"crowd_strike_customer_id":             types.StringType,
					"device_enrollement_domain":            types.StringType,
					"disk_encrypted":                       types.BoolType,
					"key_trust_level":                      types.StringType,
					"os_firewall":                          types.BoolType,
					"os_version":                           types.StringType,
					"password_proctection_warning_trigger": types.StringType,
					"realtime_url_check_mode":              types.BoolType,
					"safe_browsing_protection_level":       types.StringType,
					"screen_lock_secured":                  types.BoolType,
					// TODU window only
					"secure_boot_enabled":    types.BoolType,
					"site_isolation_enabled": types.BoolType,
					// TODU window only
					"third_party_blocking_enabled": types.BoolType,
					// TODU window only
					"window_machine_domain": types.StringType,
					// TODU window only
					"window_user_domain": types.StringType,
				},
			},
			"created_date": schema.StringAttribute{
				Description: "Created date",
				Computed:    true,
			},
			"created_by": schema.StringAttribute{
				Description: "Created by",
				Computed:    true,
			},
			"last_update": schema.StringAttribute{
				Description: "Last update",
				Computed:    true,
			},
			"last_updated_by": schema.StringAttribute{
				Description: "Last updated by",
				Computed:    true,
			},
		},
	}
}

// TODU
func (r *policyDeviceAssuranceMacOSResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state policyDeviceAssuranceMacOSResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// reqBody, err := buildDeviceAssurancePolicyRequest(state)
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"failed to build device assurance request",
	// 		err.Error(),
	// 	)
	// 	return
	// }
	reqBody, diag := buildDeviceAssuranceMacOSPolicyRequest(state)
	resp.Diagnostics.Append(diag)

	deviceAssurance, _, err := r.v3Client.DeviceAssuranceApi.CreateDeviceAssurancePolicy(ctx).DeviceAssurance(reqBody).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to create device assurance",
			err.Error(),
		)
		return
	}
	// TODU need to do additional read?
	resp.Diagnostics.Append(mapDeviceAssuranceMacOSToState(deviceAssurance, state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// TODU
func (r *policyDeviceAssuranceMacOSResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state policyDeviceAssuranceMacOSResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	_, err := r.v3Client.DeviceAssuranceApi.DeleteDeviceAssurancePolicy(ctx, state.ID.String()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to delete device assurance",
			err.Error(),
		)
		return
	}
}

// TODU
func (r *policyDeviceAssuranceMacOSResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state policyDeviceAssuranceMacOSResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	deviceAssurance, _, err := r.v3Client.DeviceAssuranceApi.GetDeviceAssurancePolicy(ctx, state.ID.String()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to read device assurance",
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(mapDeviceAssuranceMacOSToState(deviceAssurance, state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// TODU
func (r *policyDeviceAssuranceMacOSResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state policyDeviceAssuranceMacOSResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// reqBody, err := buildDeviceAssurancePolicyRequest(state)
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"failed to build device assurance request",
	// 		err.Error(),
	// 	)
	// 	return
	// }
	reqBody, diag := buildDeviceAssuranceMacOSPolicyRequest(state)
	resp.Diagnostics.Append(diag)

	deviceAssurance, _, err := r.v3Client.DeviceAssuranceApi.ReplaceDeviceAssurancePolicy(ctx, state.ID.String()).DeviceAssurance(reqBody).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to create device assurance",
			err.Error(),
		)
		return
	}
	// TODU need to do additional read?
	resp.Diagnostics.Append(mapDeviceAssuranceMacOSToState(deviceAssurance, state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// TODU
// func buildDeviceAssurancePolicyRequest(model policyDeviceAssuranceResourceModel) (okta.ListDeviceAssurancePolicies200ResponseInner, error) {
// 	var android = &okta.DeviceAssuranceAndroidPlatform{}
// 	var iOS = &okta.DeviceAssuranceIOSPlatform{}
// 	var chromeOS = &okta.DeviceAssuranceChromeOSPlatform{}
// 	var macOS = &okta.DeviceAssuranceMacOSPlatform{}
// 	var windows = &okta.DeviceAssuranceWindowsPlatform{}
// 	switch model.Platform.ValueString() {
// 	case string(okta.PLATFORM_ANDROID):
// 		android.SetName(model.Name.ValueString())
// 		android.SetPlatform(okta.Platform(model.Platform.ValueString()))
// 	case string(okta.PLATFORM_IOS):
// 		iOS.SetName(model.Name.String())
// 		iOS.SetPlatform(okta.Platform(model.Platform.String()))
// 	case string(okta.PLATFORM_CHROMEOS):
// 		chromeOS.SetName(model.Name.String())
// 		chromeOS.SetPlatform(okta.Platform(model.Platform.String()))
// 		tsp := okta.DeviceAssuranceChromeOSPlatformAllOfThirdPartySignalProviders{}
// 		tsp.Dtc.SetAllowScreenLock(model.ThirdPartySignalProviders.AllowScreenLock.ValueBool())
// 		tsp.Dtc.SetBrowserVersion(okta.ChromeBrowserVersion{Minimum: model.ThirdPartySignalProviders.BrowserVersion.ValueStringPointer()})
// 		tsp.Dtc.SetBuiltInDnsClientEnabled(model.ThirdPartySignalProviders.BuiltInDNSClientEnabled.ValueBool())
// 		tsp.Dtc.SetChromeRemoteDesktopAppBlocked(model.ThirdPartySignalProviders.ChromeRemoteDesktopAppBlocked.ValueBool())
// 		tsp.Dtc.SetDeviceEnrollmentDomain(model.ThirdPartySignalProviders.DeviceEnrollementDomain.ValueString())
// 		tsp.Dtc.SetDiskEnrypted(model.ThirdPartySignalProviders.DiskEncrypted.ValueBool())
// 		tsp.Dtc.SetKeyTrustLevel(okta.KeyTrustLevelOSMode(model.ThirdPartySignalProviders.KeyTrustLevel.ValueString()))
// 		tsp.Dtc.SetOsFirewall(model.ThirdPartySignalProviders.OsFirewall.ValueBool())
// 		tsp.Dtc.SetOsVersion(okta.OSVersion{Minimum: model.ThirdPartySignalProviders.OsVersion.ValueStringPointer()})
// 		tsp.Dtc.SetPasswordProtectionWarningTrigger(okta.PasswordProtectionWarningTrigger(model.ThirdPartySignalProviders.PasswordProctectionWarningTrigger.ValueString()))
// 		tsp.Dtc.SetRealtimeUrlCheckMode(model.ThirdPartySignalProviders.RealtimeURLCheckMode.ValueBool())
// 		tsp.Dtc.SetSafeBrowsingProtectionLevel(okta.SafeBrowsingProtectionLevel(model.ThirdPartySignalProviders.SafeBrowsingProtectionLevel.ValueString()))
// 		tsp.Dtc.SetScreenLockSecured(model.ThirdPartySignalProviders.ScreenLockSecured.ValueBool())
// 		tsp.Dtc.SetSiteIsolationEnabled(model.ThirdPartySignalProviders.SiteIsolationEnabled.ValueBool())
// 		chromeOS.SetThirdPartySignalProviders(tsp)
// 	case string(okta.PLATFORM_MACOS):
// 		macOS.SetName(model.Name.String())
// 		macOS.SetPlatform(okta.Platform(model.Platform.String()))
// 	case string(okta.PLATFORM_WINDOWS):
// 		windows.SetName(model.Name.String())
// 		windows.SetPlatform(okta.Platform(model.Platform.String()))
// 	default:
// 		return okta.ListDeviceAssurancePolicies200ResponseInner{}, errors.New("unidentified platform")
// 	}
// 	return okta.ListDeviceAssurancePolicies200ResponseInner{
// 		DeviceAssuranceAndroidPlatform:  android,
// 		DeviceAssuranceIOSPlatform:      iOS,
// 		DeviceAssuranceChromeOSPlatform: chromeOS,
// 		DeviceAssuranceMacOSPlatform:    macOS,
// 		DeviceAssuranceWindowsPlatform:  windows,
// 	}, nil
// }

func buildDeviceAssuranceMacOSPolicyRequest(model policyDeviceAssuranceMacOSResourceModel) (okta.ListDeviceAssurancePolicies200ResponseInner, diag.Diagnostic) {
	var android = &okta.DeviceAssuranceAndroidPlatform{}
	var iOS = &okta.DeviceAssuranceIOSPlatform{}
	var chromeOS = &okta.DeviceAssuranceChromeOSPlatform{}
	var macOS = &okta.DeviceAssuranceMacOSPlatform{}
	var windows = &okta.DeviceAssuranceWindowsPlatform{}
	switch model.Platform.ValueString() {
	case string(okta.PLATFORM_ANDROID):
		android.SetName(model.Name.ValueString())
		android.SetPlatform(okta.Platform(model.Platform.ValueString()))
	case string(okta.PLATFORM_IOS):
		iOS.SetName(model.Name.String())
		iOS.SetPlatform(okta.Platform(model.Platform.String()))
	case string(okta.PLATFORM_CHROMEOS):
		chromeOS.SetName(model.Name.String())
		chromeOS.SetPlatform(okta.Platform(model.Platform.String()))
		tsp := okta.DeviceAssuranceChromeOSPlatformAllOfThirdPartySignalProviders{}
		tsp.Dtc.SetAllowScreenLock(model.ThirdPartySignalProviders.AllowScreenLock.ValueBool())
		tsp.Dtc.SetBrowserVersion(okta.ChromeBrowserVersion{Minimum: model.ThirdPartySignalProviders.BrowserVersion.ValueStringPointer()})
		tsp.Dtc.SetBuiltInDnsClientEnabled(model.ThirdPartySignalProviders.BuiltInDNSClientEnabled.ValueBool())
		tsp.Dtc.SetChromeRemoteDesktopAppBlocked(model.ThirdPartySignalProviders.ChromeRemoteDesktopAppBlocked.ValueBool())
		tsp.Dtc.SetDeviceEnrollmentDomain(model.ThirdPartySignalProviders.DeviceEnrollementDomain.ValueString())
		tsp.Dtc.SetDiskEnrypted(model.ThirdPartySignalProviders.DiskEncrypted.ValueBool())
		tsp.Dtc.SetKeyTrustLevel(okta.KeyTrustLevelOSMode(model.ThirdPartySignalProviders.KeyTrustLevel.ValueString()))
		tsp.Dtc.SetOsFirewall(model.ThirdPartySignalProviders.OsFirewall.ValueBool())
		tsp.Dtc.SetOsVersion(okta.OSVersion{Minimum: model.ThirdPartySignalProviders.OsVersion.ValueStringPointer()})
		tsp.Dtc.SetPasswordProtectionWarningTrigger(okta.PasswordProtectionWarningTrigger(model.ThirdPartySignalProviders.PasswordProctectionWarningTrigger.ValueString()))
		tsp.Dtc.SetRealtimeUrlCheckMode(model.ThirdPartySignalProviders.RealtimeURLCheckMode.ValueBool())
		tsp.Dtc.SetSafeBrowsingProtectionLevel(okta.SafeBrowsingProtectionLevel(model.ThirdPartySignalProviders.SafeBrowsingProtectionLevel.ValueString()))
		tsp.Dtc.SetScreenLockSecured(model.ThirdPartySignalProviders.ScreenLockSecured.ValueBool())
		tsp.Dtc.SetSiteIsolationEnabled(model.ThirdPartySignalProviders.SiteIsolationEnabled.ValueBool())
		chromeOS.SetThirdPartySignalProviders(tsp)
	case string(okta.PLATFORM_MACOS):
		macOS.SetName(model.Name.String())
		macOS.SetPlatform(okta.Platform(model.Platform.String()))
	case string(okta.PLATFORM_WINDOWS):
		windows.SetName(model.Name.String())
		windows.SetPlatform(okta.Platform(model.Platform.String()))
	default:
		return okta.ListDeviceAssurancePolicies200ResponseInner{}, diag.NewErrorDiagnostic("unidentified platform ", model.Platform.ValueString())
	}
	return okta.ListDeviceAssurancePolicies200ResponseInner{
		DeviceAssuranceAndroidPlatform:  android,
		DeviceAssuranceIOSPlatform:      iOS,
		DeviceAssuranceChromeOSPlatform: chromeOS,
		DeviceAssuranceMacOSPlatform:    macOS,
		DeviceAssuranceWindowsPlatform:  windows,
	}, nil
}

// Map response body to schema
func mapDeviceAssuranceMacOSToState(data *okta.ListDeviceAssurancePolicies200ResponseInner, state policyDeviceAssuranceMacOSResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}
