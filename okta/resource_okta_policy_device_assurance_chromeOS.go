package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/okta-sdk-golang/v3/okta"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &policyDeviceAssuranceChromeOSResource{}
	_ resource.ResourceWithConfigure = &policyDeviceAssuranceChromeOSResource{}
	// _ resource.ResourceWithImportState = &policyDeviceAssuranceResource{}
)

func NewPolicyDeviceAssuranceChromeOSResource() resource.Resource {
	return &policyDeviceAssuranceChromeOSResource{}
}

type policyDeviceAssuranceChromeOSResource struct {
	*Config
}

type policyDeviceAssuranceChromeOSResourceModel struct {
	ID                                   types.String `tfsdk:"id"`
	Name                                 types.String `tfsdk:"name"`
	Platform                             types.String `tfsdk:"platform"`
	CreateDate                           types.String `tfsdk:"created_date"`
	CreateBy                             types.String `tfsdk:"created_by"`
	LastUpdate                           types.String `tfsdk:"last_update"`
	LastUpdatedBy                        types.String `tfsdk:"last_updated_by"`
	TpspAllowScreenLock                  types.Bool   `tfsdk:"tpsp_allow_screen_lock"`
	TpspBrowserVersion                   types.String `tfsdk:"tpsp_browser_version"`
	TpspBuiltInDNSClientEnabled          types.Bool   `tfsdk:"tpsp_builtin_dns_client_enabled"`
	TpspChromeRemoteDesktopAppBlocked    types.Bool   `tfsdk:"tpsp_chrome_remote_desktop_app_blocked"`
	TpspDeviceEnrollmentDomain           types.String `tfsdk:"tpsp_device_enrollment_domain"`
	TpspDiskEncrypted                    types.Bool   `tfsdk:"tpsp_disk_encrypted"`
	TpspKeyTrustLevel                    types.String `tfsdk:"tpsp_key_trust_level"`
	TpspOsFirewall                       types.Bool   `tfsdk:"tpsp_os_firewall"`
	TpspOsVersion                        types.String `tfsdk:"tpsp_os_version"`
	TpspPasswordProtectionWarningTrigger types.String `tfsdk:"tpsp_password_proctection_warning_trigger"`
	TpspRealtimeURLCheckMode             types.Bool   `tfsdk:"tpsp_realtime_url_check_mode"`
	TpspSafeBrowsingProtectionLevel      types.String `tfsdk:"tpsp_safe_browsing_protection_level"`
	TpspScreenLockSecured                types.Bool   `tfsdk:"tpsp_screen_lock_secured"`
	TpspSiteIsolationEnabled             types.Bool   `tfsdk:"tpsp_site_isolation_enabled"`
}

func (r *policyDeviceAssuranceChromeOSResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_device_assurance_chromeos"
}

func (r *policyDeviceAssuranceChromeOSResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages device assurance on policy",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Policy assurance id",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Policy device assurance name",
				Required:    true,
			},
			"platform": schema.StringAttribute{
				Description: "Policy device assurance platform",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tpsp_allow_screen_lock": schema.BoolAttribute{
				Description: "Third party signal provider allow screen lock",
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
			"tpsp_site_isolation_enabled": schema.BoolAttribute{
				Description: "Third party signal provider site isolation enabled",
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
func (r *policyDeviceAssuranceChromeOSResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	p, ok := req.ProviderData.(*Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.Config = p
}

func (r *policyDeviceAssuranceChromeOSResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state policyDeviceAssuranceChromeOSResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqBody, err := buildDeviceAssuranceChromeOSPolicyRequest(state)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to build device assurance request",
			err.Error(),
		)
		return
	}

	deviceAssurance, _, err := r.v3Client.DeviceAssuranceApi.CreateDeviceAssurancePolicy(ctx).DeviceAssurance(reqBody).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to create device assurance",
			err.Error(),
		)
		return
	}

	// TODU need to do additional read?
	resp.Diagnostics.Append(mapDeviceAssuranceChromeOSToState(deviceAssurance, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *policyDeviceAssuranceChromeOSResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state policyDeviceAssuranceChromeOSResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deviceAssurance, _, err := r.v3Client.DeviceAssuranceApi.GetDeviceAssurancePolicy(ctx, state.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to read device assurance",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(mapDeviceAssuranceChromeOSToState(deviceAssurance, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *policyDeviceAssuranceChromeOSResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state policyDeviceAssuranceChromeOSResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.v3Client.DeviceAssuranceApi.DeleteDeviceAssurancePolicy(ctx, state.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to delete device assurance",
			err.Error(),
		)
		return
	}
}

func (r *policyDeviceAssuranceChromeOSResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state policyDeviceAssuranceChromeOSResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqBody, err := buildDeviceAssuranceChromeOSPolicyRequest(state)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to build device assurance request",
			err.Error(),
		)
		return
	}

	deviceAssurance, _, err := r.v3Client.DeviceAssuranceApi.ReplaceDeviceAssurancePolicy(ctx, state.ID.ValueString()).DeviceAssurance(reqBody).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to create device assurance",
			err.Error(),
		)
		return
	}

	// TODU need to do additional read?
	resp.Diagnostics.Append(mapDeviceAssuranceChromeOSToState(deviceAssurance, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func buildDeviceAssuranceChromeOSPolicyRequest(model policyDeviceAssuranceChromeOSResourceModel) (okta.ListDeviceAssurancePolicies200ResponseInner, error) {
	var chromeOS = &okta.DeviceAssuranceChromeOSPlatform{}
	chromeOS.SetName(model.Name.ValueString())
	chromeOS.SetPlatform(okta.PLATFORM_CHROMEOS)

	var thirdPartySignalProviders okta.DeviceAssuranceChromeOSPlatformAllOfThirdPartySignalProviders
	var dtc okta.DTCChromeOS
	dtc.AllowScreenLock = model.TpspAllowScreenLock.ValueBoolPointer()
	if !model.TpspBrowserVersion.IsNull() {
		dtc.BrowserVersion = &okta.ChromeBrowserVersion{Minimum: model.TpspBrowserVersion.ValueStringPointer()}
	}
	dtc.BuiltInDnsClientEnabled = model.TpspBuiltInDNSClientEnabled.ValueBoolPointer()
	dtc.ChromeRemoteDesktopAppBlocked = model.TpspChromeRemoteDesktopAppBlocked.ValueBoolPointer()
	dtc.DeviceEnrollmentDomain = model.TpspDeviceEnrollmentDomain.ValueStringPointer()
	dtc.DiskEncrypted = model.TpspDiskEncrypted.ValueBoolPointer()
	if !model.TpspKeyTrustLevel.IsNull() {
		v, err := okta.NewKeyTrustLevelOSModeFromValue(model.TpspKeyTrustLevel.ValueString())
		if err != nil {
			return okta.ListDeviceAssurancePolicies200ResponseInner{DeviceAssuranceChromeOSPlatform: chromeOS}, err
		}
		dtc.KeyTrustLevel = v
	}
	dtc.OsFirewall = model.TpspOsFirewall.ValueBoolPointer()
	if !model.TpspOsVersion.IsNull() {
		dtc.OsVersion = &okta.OSVersion{Minimum: model.TpspOsVersion.ValueStringPointer()}
	}
	if !model.TpspPasswordProtectionWarningTrigger.IsNull() {
		v, err := okta.NewPasswordProtectionWarningTriggerFromValue(model.TpspPasswordProtectionWarningTrigger.ValueString())
		if err != nil {
			return okta.ListDeviceAssurancePolicies200ResponseInner{DeviceAssuranceChromeOSPlatform: chromeOS}, err
		}
		dtc.PasswordProtectionWarningTrigger = v
	}
	dtc.RealtimeUrlCheckMode = model.TpspRealtimeURLCheckMode.ValueBoolPointer()
	if !model.TpspSafeBrowsingProtectionLevel.IsNull() {
		v, err := okta.NewSafeBrowsingProtectionLevelFromValue(model.TpspSafeBrowsingProtectionLevel.ValueString())
		if err != nil {
			return okta.ListDeviceAssurancePolicies200ResponseInner{DeviceAssuranceChromeOSPlatform: chromeOS}, err
		}
		dtc.SafeBrowsingProtectionLevel = v
	}
	dtc.ScreenLockSecured = model.TpspScreenLockSecured.ValueBoolPointer()
	dtc.SiteIsolationEnabled = model.TpspSiteIsolationEnabled.ValueBoolPointer()
	thirdPartySignalProviders.SetDtc(dtc)
	chromeOS.SetThirdPartySignalProviders(thirdPartySignalProviders)

	return okta.ListDeviceAssurancePolicies200ResponseInner{DeviceAssuranceChromeOSPlatform: chromeOS}, nil
}

// Map response body to schema
func mapDeviceAssuranceChromeOSToState(data *okta.ListDeviceAssurancePolicies200ResponseInner, state *policyDeviceAssuranceChromeOSResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	if data.DeviceAssuranceChromeOSPlatform == nil {
		diags.AddError("Empty response", "ChromeOS object")
		return diags
	}

	state.ID = types.StringPointerValue(data.DeviceAssuranceChromeOSPlatform.Id)
	state.Name = types.StringPointerValue(data.DeviceAssuranceChromeOSPlatform.Name)
	state.Platform = types.StringPointerValue((*string)(data.DeviceAssuranceChromeOSPlatform.Platform))

	if _, ok := data.DeviceAssuranceChromeOSPlatform.ThirdPartySignalProviders.GetDtcOk(); ok {
		state.TpspAllowScreenLock = types.BoolPointerValue(data.DeviceAssuranceChromeOSPlatform.ThirdPartySignalProviders.Dtc.AllowScreenLock)
		if _, ok := data.DeviceAssuranceChromeOSPlatform.ThirdPartySignalProviders.Dtc.GetBrowserVersionOk(); ok {
			state.TpspBrowserVersion = types.StringPointerValue(data.DeviceAssuranceChromeOSPlatform.ThirdPartySignalProviders.Dtc.BrowserVersion.Minimum)
		}
		state.TpspBuiltInDNSClientEnabled = types.BoolPointerValue(data.DeviceAssuranceChromeOSPlatform.ThirdPartySignalProviders.Dtc.BuiltInDnsClientEnabled)
		state.TpspChromeRemoteDesktopAppBlocked = types.BoolPointerValue(data.DeviceAssuranceChromeOSPlatform.ThirdPartySignalProviders.Dtc.ChromeRemoteDesktopAppBlocked)
		state.TpspDeviceEnrollmentDomain = types.StringPointerValue(data.DeviceAssuranceChromeOSPlatform.ThirdPartySignalProviders.Dtc.DeviceEnrollmentDomain)
		state.TpspDiskEncrypted = types.BoolPointerValue(data.DeviceAssuranceChromeOSPlatform.ThirdPartySignalProviders.Dtc.DiskEncrypted)
		if _, ok := data.DeviceAssuranceChromeOSPlatform.ThirdPartySignalProviders.Dtc.GetKeyTrustLevelOk(); ok {
			state.TpspKeyTrustLevel = types.StringPointerValue((*string)(data.DeviceAssuranceChromeOSPlatform.ThirdPartySignalProviders.Dtc.KeyTrustLevel))
		}
		state.TpspOsFirewall = types.BoolPointerValue(data.DeviceAssuranceChromeOSPlatform.ThirdPartySignalProviders.Dtc.OsFirewall)
		if _, ok := data.DeviceAssuranceChromeOSPlatform.ThirdPartySignalProviders.Dtc.GetOsVersionOk(); ok {
			state.TpspOsVersion = types.StringPointerValue(data.DeviceAssuranceChromeOSPlatform.ThirdPartySignalProviders.Dtc.OsVersion.Minimum)
		}
		state.TpspPasswordProtectionWarningTrigger = types.StringPointerValue((*string)(data.DeviceAssuranceChromeOSPlatform.ThirdPartySignalProviders.Dtc.PasswordProtectionWarningTrigger))
		state.TpspRealtimeURLCheckMode = types.BoolPointerValue(data.DeviceAssuranceChromeOSPlatform.ThirdPartySignalProviders.Dtc.RealtimeUrlCheckMode)
		state.TpspSafeBrowsingProtectionLevel = types.StringPointerValue((*string)(data.DeviceAssuranceChromeOSPlatform.ThirdPartySignalProviders.Dtc.SafeBrowsingProtectionLevel))
		state.TpspScreenLockSecured = types.BoolPointerValue(data.DeviceAssuranceChromeOSPlatform.ThirdPartySignalProviders.Dtc.ScreenLockSecured)
		state.TpspSiteIsolationEnabled = types.BoolPointerValue(data.DeviceAssuranceChromeOSPlatform.ThirdPartySignalProviders.Dtc.SiteIsolationEnabled)
	}

	state.CreateDate = types.StringPointerValue(data.DeviceAssuranceChromeOSPlatform.CreatedDate)
	state.CreateBy = types.StringPointerValue(data.DeviceAssuranceChromeOSPlatform.CreatedBy)
	state.LastUpdate = types.StringPointerValue(data.DeviceAssuranceChromeOSPlatform.LastUpdate)
	state.LastUpdatedBy = types.StringPointerValue(data.DeviceAssuranceChromeOSPlatform.LastUpdatedBy)
	return diags
}
