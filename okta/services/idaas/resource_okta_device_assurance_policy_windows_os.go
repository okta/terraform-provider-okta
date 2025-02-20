package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/okta-sdk-golang/v4/okta"
	"github.com/okta/terraform-provider-okta/okta/config"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &policyDeviceAssuranceWindowsResource{}
	_ resource.ResourceWithConfigure   = &policyDeviceAssuranceWindowsResource{}
	_ resource.ResourceWithImportState = &policyDeviceAssuranceWindowsResource{}
)

func NewPolicyDeviceAssuranceWindowsResource() resource.Resource {
	return &policyDeviceAssuranceWindowsResource{}
}

type policyDeviceAssuranceWindowsResource struct {
	*config.Config
}

type policyDeviceAssuranceWindowsResourceModel struct {
	ID                                   types.String   `tfsdk:"id"`
	Name                                 types.String   `tfsdk:"name"`
	Platform                             types.String   `tfsdk:"platform"`
	DiskEncryptionType                   []types.String `tfsdk:"disk_encryption_type"`
	OsVersion                            types.String   `tfsdk:"os_version"`
	SecureHardwarePresent                types.Bool     `tfsdk:"secure_hardware_present"`
	ScreenLockType                       []types.String `tfsdk:"screenlock_type"`
	CreateDate                           types.String   `tfsdk:"created_date"`
	CreateBy                             types.String   `tfsdk:"created_by"`
	LastUpdate                           types.String   `tfsdk:"last_update"`
	LastUpdatedBy                        types.String   `tfsdk:"last_updated_by"`
	ThirdPartySignalProviders            types.Bool     `tfsdk:"third_party_signal_providers"`
	TpspBrowserVersion                   types.String   `tfsdk:"tpsp_browser_version"`
	TpspBuiltInDNSClientEnabled          types.Bool     `tfsdk:"tpsp_builtin_dns_client_enabled"`
	TpspChromeRemoteDesktopAppBlocked    types.Bool     `tfsdk:"tpsp_chrome_remote_desktop_app_blocked"`
	TpspCrowdStrikeAgentID               types.String   `tfsdk:"tpsp_crowd_strike_agent_id"`
	TpspCrowdStrikeCustomerID            types.String   `tfsdk:"tpsp_crowd_strike_customer_id"`
	TpspDeviceEnrollmentDomain           types.String   `tfsdk:"tpsp_device_enrollment_domain"`
	TpspDiskEncrypted                    types.Bool     `tfsdk:"tpsp_disk_encrypted"`
	TpspKeyTrustLevel                    types.String   `tfsdk:"tpsp_key_trust_level"`
	TpspOsFirewall                       types.Bool     `tfsdk:"tpsp_os_firewall"`
	TpspOsVersion                        types.String   `tfsdk:"tpsp_os_version"`
	TpspPasswordProtectionWarningTrigger types.String   `tfsdk:"tpsp_password_proctection_warning_trigger"`
	TpspRealtimeURLCheckMode             types.Bool     `tfsdk:"tpsp_realtime_url_check_mode"`
	TpspSafeBrowsingProtectionLevel      types.String   `tfsdk:"tpsp_safe_browsing_protection_level"`
	TpspScreenLockSecured                types.Bool     `tfsdk:"tpsp_screen_lock_secured"`
	TpspSecureBootEnabled                types.Bool     `tfsdk:"tpsp_secure_boot_enabled"`
	TpspSiteIsolationEnabled             types.Bool     `tfsdk:"tpsp_site_isolation_enabled"`
	TpspThirdPartyBlockingEnabled        types.Bool     `tfsdk:"tpsp_third_party_blocking_enabled"`
	TpspWindowsMachineDomain             types.String   `tfsdk:"tpsp_windows_machine_domain"`
	TpspWindowsUserDomain                types.String   `tfsdk:"tpsp_windows_user_domain"`
}

func (r *policyDeviceAssuranceWindowsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_device_assurance_windows"
}

func (r *policyDeviceAssuranceWindowsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a device assurance policy for windows.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Policy assurance id",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the device assurance policy.",
				Required:    true,
			},
			"platform": schema.StringAttribute{
				Description: "Policy device assurance platform",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			// use set to avoid order change as v3 does not have diff suppress func
			"disk_encryption_type": schema.SetAttribute{
				Description: "List of disk encryption type, can be `ALL_INTERNAL_VOLUMES`",
				Optional:    true,
				ElementType: types.StringType,
			},
			"os_version": schema.StringAttribute{
				Description: "Minimum os version of the device in the device assurance policy.",
				Optional:    true,
			},
			"secure_hardware_present": schema.BoolAttribute{
				Description: "Is the device secure with hardware in the device assurance policy.",
				Optional:    true,
			},
			"screenlock_type": schema.SetAttribute{
				Description: "List of screenlock type, can be `BIOMETRIC` or `BIOMETRIC, PASSCODE`",
				Optional:    true,
				ElementType: types.StringType,
			},
			"third_party_signal_providers": schema.BoolAttribute{
				Description: "Check to include third party signal provider",
				Optional:    true,
			},
			"tpsp_browser_version": schema.StringAttribute{
				Description: "Third party signal provider minimum browser version",
				Optional:    true,
			},
			"tpsp_builtin_dns_client_enabled": schema.BoolAttribute{
				Description: "Third party signal provider builtin dns client enable",
				Optional:    true,
			},
			"tpsp_chrome_remote_desktop_app_blocked": schema.BoolAttribute{
				Description: "Third party signal provider chrome remote desktop app blocked",
				Optional:    true,
			},
			"tpsp_crowd_strike_agent_id": schema.StringAttribute{
				Description: "Third party signal provider crowdstrike agent id",
				Optional:    true,
			},
			"tpsp_crowd_strike_customer_id": schema.StringAttribute{
				Description: "Third party signal provider crowdstrike user id",
				Optional:    true,
			},
			"tpsp_device_enrollment_domain": schema.StringAttribute{
				Description: "Third party signal provider device enrollment domain",
				Optional:    true,
			},
			"tpsp_disk_encrypted": schema.BoolAttribute{
				Description: "Third party signal provider disk encrypted",
				Optional:    true,
			},
			"tpsp_key_trust_level": schema.StringAttribute{
				Description: "Third party signal provider key trust level",
				Optional:    true,
			},
			"tpsp_os_firewall": schema.BoolAttribute{
				Description: "Third party signal provider os firewall",
				Optional:    true,
			},
			"tpsp_os_version": schema.StringAttribute{
				Description: "Third party signal provider minimum os version",
				Optional:    true,
			},
			"tpsp_password_proctection_warning_trigger": schema.StringAttribute{
				Description: "Third party signal provider password protection warning trigger",
				Optional:    true,
			},
			"tpsp_realtime_url_check_mode": schema.BoolAttribute{
				Description: "Third party signal provider realtime url check mode",
				Optional:    true,
			},
			"tpsp_safe_browsing_protection_level": schema.StringAttribute{
				Description: "Third party signal provider safe browsing protection level",
				Optional:    true,
			},
			"tpsp_screen_lock_secured": schema.BoolAttribute{
				Description: "Third party signal provider screen lock secure",
				Optional:    true,
			},
			"tpsp_secure_boot_enabled": schema.BoolAttribute{
				Description: "Third party signal provider secure boot enabled",
				Optional:    true,
			},
			"tpsp_site_isolation_enabled": schema.BoolAttribute{
				Description: "Third party signal provider site isolation enabled",
				Optional:    true,
			},
			"tpsp_third_party_blocking_enabled": schema.BoolAttribute{
				Description: "Third party signal provider third party blocking enabled",
				Optional:    true,
			},
			"tpsp_windows_machine_domain": schema.StringAttribute{
				Description: "Third party signal provider windows machine domain",
				Optional:    true,
			},
			"tpsp_windows_user_domain": schema.StringAttribute{
				Description: "Third party signal provider windows user domain",
				Optional:    true,
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

// Configure adds the provider configured client to the resource.
func (r *policyDeviceAssuranceWindowsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *policyDeviceAssuranceWindowsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state policyDeviceAssuranceWindowsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqBody, err := buildDeviceAssuranceWindowsPolicyRequest(state)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to build device assurance request",
			err.Error(),
		)
		return
	}

	deviceAssurance, _, err := r.OktaSDKClientV3.DeviceAssuranceAPI.CreateDeviceAssurancePolicy(ctx).DeviceAssurance(reqBody).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to create device assurance",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(mapDeviceAssuranceWindowsToState(deviceAssurance, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *policyDeviceAssuranceWindowsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state policyDeviceAssuranceWindowsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deviceAssurance, _, err := r.OktaSDKClientV3.DeviceAssuranceAPI.GetDeviceAssurancePolicy(ctx, state.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to read device assurance",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(mapDeviceAssuranceWindowsToState(deviceAssurance, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *policyDeviceAssuranceWindowsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state policyDeviceAssuranceWindowsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.OktaSDKClientV3.DeviceAssuranceAPI.DeleteDeviceAssurancePolicy(ctx, state.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to delete device assurance",
			err.Error(),
		)
		return
	}
}

func (r *policyDeviceAssuranceWindowsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state policyDeviceAssuranceWindowsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqBody, err := buildDeviceAssuranceWindowsPolicyRequest(state)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to build device assurance request",
			err.Error(),
		)
		return
	}

	deviceAssurance, _, err := r.OktaSDKClientV3.DeviceAssuranceAPI.ReplaceDeviceAssurancePolicy(ctx, state.ID.ValueString()).DeviceAssurance(reqBody).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to update device assurance",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(mapDeviceAssuranceWindowsToState(deviceAssurance, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func buildDeviceAssuranceWindowsPolicyRequest(model policyDeviceAssuranceWindowsResourceModel) (okta.ListDeviceAssurancePolicies200ResponseInner, error) {
	windows := &okta.DeviceAssuranceWindowsPlatform{}
	windows.SetName(model.Name.ValueString())
	windows.SetPlatform("WINDOWS")
	if len(model.DiskEncryptionType) > 0 {
		diskEncryptionType := make([]string, 0)
		for _, det := range model.DiskEncryptionType {
			diskEncryptionType = append(diskEncryptionType, det.ValueString())
		}
		windows.DiskEncryptionType = &okta.DeviceAssuranceMacOSPlatformAllOfDiskEncryptionType{Include: diskEncryptionType}
	}
	if !model.OsVersion.IsNull() {
		windows.OsVersion = &okta.OSVersionFourComponents{Minimum: model.OsVersion.ValueStringPointer()}
	}
	if len(model.ScreenLockType) > 0 {
		screenlockType := make([]string, 0)
		for _, det := range model.ScreenLockType {
			screenlockType = append(screenlockType, det.ValueString())
		}
		windows.ScreenLockType = &okta.DeviceAssuranceAndroidPlatformAllOfScreenLockType{Include: screenlockType}
	}
	windows.SecureHardwarePresent = model.SecureHardwarePresent.ValueBoolPointer()

	if model.ThirdPartySignalProviders.ValueBool() {
		var thirdPartySignalProviders okta.DeviceAssuranceWindowsPlatformAllOfThirdPartySignalProviders
		var dtc okta.DTCWindows
		if !model.TpspBrowserVersion.IsNull() {
			dtc.BrowserVersion = &okta.ChromeBrowserVersion{Minimum: model.TpspBrowserVersion.ValueStringPointer()}
		}
		dtc.BuiltInDnsClientEnabled = model.TpspBuiltInDNSClientEnabled.ValueBoolPointer()
		dtc.ChromeRemoteDesktopAppBlocked = model.TpspChromeRemoteDesktopAppBlocked.ValueBoolPointer()
		dtc.CrowdStrikeAgentId = model.TpspCrowdStrikeAgentID.ValueStringPointer()
		dtc.CrowdStrikeCustomerId = model.TpspCrowdStrikeCustomerID.ValueStringPointer()
		dtc.DeviceEnrollmentDomain = model.TpspDeviceEnrollmentDomain.ValueStringPointer()
		dtc.DiskEncrypted = model.TpspDiskEncrypted.ValueBoolPointer()
		if !model.TpspKeyTrustLevel.IsNull() {
			dtc.KeyTrustLevel = model.TpspKeyTrustLevel.ValueStringPointer()
		}
		dtc.OsFirewall = model.TpspOsFirewall.ValueBoolPointer()
		if !model.TpspOsVersion.IsNull() {
			dtc.OsVersion = &okta.OSVersionFourComponents{Minimum: model.TpspOsVersion.ValueStringPointer()}
		}
		if !model.TpspPasswordProtectionWarningTrigger.IsNull() {
			dtc.PasswordProtectionWarningTrigger = model.TpspPasswordProtectionWarningTrigger.ValueStringPointer()
		}
		dtc.RealtimeUrlCheckMode = model.TpspRealtimeURLCheckMode.ValueBoolPointer()

		if !model.TpspSafeBrowsingProtectionLevel.IsNull() {
			dtc.SafeBrowsingProtectionLevel = model.TpspSafeBrowsingProtectionLevel.ValueStringPointer()
		}
		dtc.ScreenLockSecured = model.TpspScreenLockSecured.ValueBoolPointer()
		dtc.SecureBootEnabled = model.TpspSecureBootEnabled.ValueBoolPointer()
		dtc.SiteIsolationEnabled = model.TpspSiteIsolationEnabled.ValueBoolPointer()
		dtc.ThirdPartyBlockingEnabled = model.TpspThirdPartyBlockingEnabled.ValueBoolPointer()
		dtc.WindowsMachineDomain = model.TpspWindowsMachineDomain.ValueStringPointer()
		dtc.WindowsUserDomain = model.TpspWindowsUserDomain.ValueStringPointer()
		thirdPartySignalProviders.SetDtc(dtc)
		windows.SetThirdPartySignalProviders(thirdPartySignalProviders)
	}

	return okta.ListDeviceAssurancePolicies200ResponseInner{DeviceAssuranceWindowsPlatform: windows}, nil
}

// Map response body to schema
func mapDeviceAssuranceWindowsToState(data *okta.ListDeviceAssurancePolicies200ResponseInner, state *policyDeviceAssuranceWindowsResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	if data.DeviceAssuranceWindowsPlatform == nil {
		diags.AddError("Empty response", "Windows object")
		return diags
	}

	state.ID = types.StringPointerValue(data.DeviceAssuranceWindowsPlatform.Id)
	state.Name = types.StringPointerValue(data.DeviceAssuranceWindowsPlatform.Name)
	state.Platform = types.StringPointerValue((*string)(data.DeviceAssuranceWindowsPlatform.Platform))

	state.SecureHardwarePresent = types.BoolPointerValue(data.DeviceAssuranceWindowsPlatform.SecureHardwarePresent)
	if _, ok := data.DeviceAssuranceWindowsPlatform.GetOsVersionOk(); ok {
		state.OsVersion = types.StringPointerValue(data.DeviceAssuranceWindowsPlatform.OsVersion.Minimum)
	}
	if _, ok := data.DeviceAssuranceWindowsPlatform.DiskEncryptionType.GetIncludeOk(); ok {
		diskEncryptionType := make([]types.String, 0)
		for _, det := range data.DeviceAssuranceWindowsPlatform.DiskEncryptionType.GetInclude() {
			diskEncryptionType = append(diskEncryptionType, types.StringValue(string(det)))
		}
		state.DiskEncryptionType = diskEncryptionType
	}
	if _, ok := data.DeviceAssuranceWindowsPlatform.ScreenLockType.GetIncludeOk(); ok {
		screenLockType := make([]types.String, 0)
		for _, slt := range data.DeviceAssuranceWindowsPlatform.ScreenLockType.GetInclude() {
			screenLockType = append(screenLockType, types.StringValue(string(slt)))
		}
		state.ScreenLockType = screenLockType
	}

	if _, ok := data.DeviceAssuranceWindowsPlatform.ThirdPartySignalProviders.GetDtcOk(); ok {
		if _, ok := data.DeviceAssuranceWindowsPlatform.ThirdPartySignalProviders.Dtc.GetBrowserVersionOk(); ok {
			state.TpspBrowserVersion = types.StringPointerValue(data.DeviceAssuranceWindowsPlatform.ThirdPartySignalProviders.Dtc.BrowserVersion.Minimum)
		}
		state.TpspBuiltInDNSClientEnabled = types.BoolPointerValue(data.DeviceAssuranceWindowsPlatform.ThirdPartySignalProviders.Dtc.BuiltInDnsClientEnabled)
		state.TpspChromeRemoteDesktopAppBlocked = types.BoolPointerValue(data.DeviceAssuranceWindowsPlatform.ThirdPartySignalProviders.Dtc.ChromeRemoteDesktopAppBlocked)
		state.TpspCrowdStrikeAgentID = types.StringPointerValue(data.DeviceAssuranceWindowsPlatform.ThirdPartySignalProviders.Dtc.CrowdStrikeAgentId)
		state.TpspCrowdStrikeCustomerID = types.StringPointerValue(data.DeviceAssuranceWindowsPlatform.ThirdPartySignalProviders.Dtc.CrowdStrikeCustomerId)
		state.TpspDeviceEnrollmentDomain = types.StringPointerValue(data.DeviceAssuranceWindowsPlatform.ThirdPartySignalProviders.Dtc.DeviceEnrollmentDomain)
		state.TpspDiskEncrypted = types.BoolPointerValue(data.DeviceAssuranceWindowsPlatform.ThirdPartySignalProviders.Dtc.DiskEncrypted)
		if _, ok := data.DeviceAssuranceWindowsPlatform.ThirdPartySignalProviders.Dtc.GetKeyTrustLevelOk(); ok {
			state.TpspKeyTrustLevel = types.StringPointerValue((*string)(data.DeviceAssuranceWindowsPlatform.ThirdPartySignalProviders.Dtc.KeyTrustLevel))
		}
		state.TpspOsFirewall = types.BoolPointerValue(data.DeviceAssuranceWindowsPlatform.ThirdPartySignalProviders.Dtc.OsFirewall)
		if _, ok := data.DeviceAssuranceWindowsPlatform.ThirdPartySignalProviders.Dtc.GetOsVersionOk(); ok {
			state.TpspOsVersion = types.StringPointerValue(data.DeviceAssuranceWindowsPlatform.ThirdPartySignalProviders.Dtc.OsVersion.Minimum)
		}
		state.TpspPasswordProtectionWarningTrigger = types.StringPointerValue((*string)(data.DeviceAssuranceWindowsPlatform.ThirdPartySignalProviders.Dtc.PasswordProtectionWarningTrigger))
		state.TpspRealtimeURLCheckMode = types.BoolPointerValue(data.DeviceAssuranceWindowsPlatform.ThirdPartySignalProviders.Dtc.RealtimeUrlCheckMode)
		state.TpspSafeBrowsingProtectionLevel = types.StringPointerValue((*string)(data.DeviceAssuranceWindowsPlatform.ThirdPartySignalProviders.Dtc.SafeBrowsingProtectionLevel))
		state.TpspScreenLockSecured = types.BoolPointerValue(data.DeviceAssuranceWindowsPlatform.ThirdPartySignalProviders.Dtc.ScreenLockSecured)
		state.TpspSecureBootEnabled = types.BoolPointerValue(data.DeviceAssuranceWindowsPlatform.ThirdPartySignalProviders.Dtc.SecureBootEnabled)
		state.TpspSiteIsolationEnabled = types.BoolPointerValue(data.DeviceAssuranceWindowsPlatform.ThirdPartySignalProviders.Dtc.SiteIsolationEnabled)
		state.TpspThirdPartyBlockingEnabled = types.BoolPointerValue(data.DeviceAssuranceWindowsPlatform.ThirdPartySignalProviders.Dtc.ThirdPartyBlockingEnabled)
		state.TpspWindowsMachineDomain = types.StringPointerValue(data.DeviceAssuranceWindowsPlatform.ThirdPartySignalProviders.Dtc.WindowsMachineDomain)
		state.TpspWindowsUserDomain = types.StringPointerValue(data.DeviceAssuranceWindowsPlatform.ThirdPartySignalProviders.Dtc.WindowsUserDomain)
	}

	state.CreateDate = types.StringPointerValue(data.DeviceAssuranceWindowsPlatform.CreatedDate)
	state.CreateBy = types.StringPointerValue(data.DeviceAssuranceWindowsPlatform.CreatedBy)
	state.LastUpdate = types.StringPointerValue(data.DeviceAssuranceWindowsPlatform.LastUpdate)
	state.LastUpdatedBy = types.StringPointerValue(data.DeviceAssuranceWindowsPlatform.LastUpdatedBy)
	return diags
}

func (r *policyDeviceAssuranceWindowsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
