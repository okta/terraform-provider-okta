package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	v6okta "github.com/okta/okta-sdk-golang/v6/okta"
	"github.com/okta/terraform-provider-okta/okta/config"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &policyDeviceAssuranceChromeOSResource{}
	_ resource.ResourceWithConfigure   = &policyDeviceAssuranceChromeOSResource{}
	_ resource.ResourceWithImportState = &policyDeviceAssuranceChromeOSResource{}
)

func newPolicyDeviceAssuranceChromeOSResource() resource.Resource {
	return &policyDeviceAssuranceChromeOSResource{}
}

type policyDeviceAssuranceChromeOSResource struct {
	*config.Config
}

type policyDeviceAssuranceChromeOSResourceModel struct {
	ID                                   types.String      `tfsdk:"id"`
	Name                                 types.String      `tfsdk:"name"`
	Platform                             types.String      `tfsdk:"platform"`
	GracePeriod                          *gracePeriodModel `tfsdk:"grace_period"`
	DisplayRemediationMode               types.String      `tfsdk:"display_remediation_mode"`
	CreateDate                           types.String      `tfsdk:"created_date"`
	CreateBy                             types.String      `tfsdk:"created_by"`
	LastUpdate                           types.String      `tfsdk:"last_update"`
	LastUpdatedBy                        types.String      `tfsdk:"last_updated_by"`
	TpspAllowScreenLock                  types.Bool        `tfsdk:"tpsp_allow_screen_lock"`
	TpspBrowserVersion                   types.String      `tfsdk:"tpsp_browser_version"`
	TpspBuiltInDNSClientEnabled          types.Bool        `tfsdk:"tpsp_builtin_dns_client_enabled"`
	TpspChromeRemoteDesktopAppBlocked    types.Bool        `tfsdk:"tpsp_chrome_remote_desktop_app_blocked"`
	TpspDeviceEnrollmentDomain           types.String      `tfsdk:"tpsp_device_enrollment_domain"`
	TpspDiskEncrypted                    types.Bool        `tfsdk:"tpsp_disk_encrypted"`
	TpspKeyTrustLevel                    types.String      `tfsdk:"tpsp_key_trust_level"`
	TpspOsFirewall                       types.Bool        `tfsdk:"tpsp_os_firewall"`
	TpspOsVersion                        types.String      `tfsdk:"tpsp_os_version"`
	TpspPasswordProtectionWarningTrigger types.String      `tfsdk:"tpsp_password_proctection_warning_trigger"`
	TpspRealtimeURLCheckMode             types.Bool        `tfsdk:"tpsp_realtime_url_check_mode"`
	TpspSafeBrowsingProtectionLevel      types.String      `tfsdk:"tpsp_safe_browsing_protection_level"`
	TpspScreenLockSecured                types.Bool        `tfsdk:"tpsp_screen_lock_secured"`
	TpspSiteIsolationEnabled             types.Bool        `tfsdk:"tpsp_site_isolation_enabled"`
}

func (r *policyDeviceAssuranceChromeOSResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_device_assurance_chromeos"
}

func (r *policyDeviceAssuranceChromeOSResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a device assurance policy for chromeos.",
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
			"tpsp_allow_screen_lock": schema.BoolAttribute{
				Description: "Third party signal provider allow screen lock",
				Optional:    true,
			},
			"tpsp_browser_version": schema.StringAttribute{
				Description: "Third party signal provider minimum browser version",
				Optional:    true,
			},
			"tpsp_builtin_dns_client_enabled": schema.BoolAttribute{
				Description: "Third party signal provider builtin dns client enabled",
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
			"display_remediation_mode": schema.StringAttribute{
				Description: "Display remediation mode for non-compliant devices (Early Access feature): HIDE or SHOW.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("HIDE", "SHOW"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
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
		Blocks: map[string]schema.Block{
			"grace_period": schema.SingleNestedBlock{
				Description: "Grace period configuration for the device assurance policy (Early Access feature).",
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description: "Grace period type: BY_DATE_TIME or BY_DURATION.",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOf("BY_DATE_TIME", "BY_DURATION"),
						},
					},
					"expiry": schema.StringAttribute{
						Description: "Grace period expiry. ISO 8601 datetime (e.g. 2024-12-01T00:00:00.000Z) for BY_DATE_TIME, or ISO 8601 duration (e.g. P7D, P30D, 1-180 days) for BY_DURATION.",
						Optional:    true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *policyDeviceAssuranceChromeOSResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
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

	deviceAssurance, _, err := r.OktaIDaaSClient.OktaSDKClientV6().DeviceAssuranceAPI.CreateDeviceAssurancePolicy(ctx).DeviceAssurance(reqBody).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to create device assurance",
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

func (r *policyDeviceAssuranceChromeOSResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state policyDeviceAssuranceChromeOSResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deviceAssurance, _, err := r.OktaIDaaSClient.OktaSDKClientV6().DeviceAssuranceAPI.GetDeviceAssurancePolicy(ctx, state.ID.ValueString()).Execute()
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

	_, err := r.OktaIDaaSClient.OktaSDKClientV6().DeviceAssuranceAPI.DeleteDeviceAssurancePolicy(ctx, state.ID.ValueString()).Execute()
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

	deviceAssurance, _, err := r.OktaIDaaSClient.OktaSDKClientV6().DeviceAssuranceAPI.ReplaceDeviceAssurancePolicy(ctx, state.ID.ValueString()).DeviceAssurance(reqBody).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to update device assurance",
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

func buildDeviceAssuranceChromeOSPolicyRequest(model policyDeviceAssuranceChromeOSResourceModel) (v6okta.ListDeviceAssurancePolicies200ResponseInner, error) {
	chromeOS := &v6okta.DeviceAssuranceChromeOSPlatform{}
	chromeOS.SetName(model.Name.ValueString())
	chromeOS.SetPlatform("CHROMEOS")

	var thirdPartySignalProviders v6okta.DeviceAssuranceChromeOSPlatformAllOfThirdPartySignalProviders
	var dtc v6okta.DTCChromeOS
	dtc.AllowScreenLock = model.TpspAllowScreenLock.ValueBoolPointer()
	if !model.TpspBrowserVersion.IsNull() {
		dtc.BrowserVersion = &v6okta.ChromeBrowserVersion{Minimum: model.TpspBrowserVersion.ValueStringPointer()}
	}
	dtc.BuiltInDnsClientEnabled = model.TpspBuiltInDNSClientEnabled.ValueBoolPointer()
	dtc.ChromeRemoteDesktopAppBlocked = model.TpspChromeRemoteDesktopAppBlocked.ValueBoolPointer()
	dtc.DeviceEnrollmentDomain = model.TpspDeviceEnrollmentDomain.ValueStringPointer()
	dtc.DiskEncrypted = model.TpspDiskEncrypted.ValueBoolPointer()
	if !model.TpspKeyTrustLevel.IsNull() {
		dtc.KeyTrustLevel = model.TpspKeyTrustLevel.ValueStringPointer()
	}
	dtc.OsFirewall = model.TpspOsFirewall.ValueBoolPointer()
	if !model.TpspOsVersion.IsNull() {
		dtc.OsVersion = &v6okta.OSVersionFourComponents{Minimum: model.TpspOsVersion.ValueStringPointer()}
	}
	if !model.TpspPasswordProtectionWarningTrigger.IsNull() {
		dtc.PasswordProtectionWarningTrigger = model.TpspPasswordProtectionWarningTrigger.ValueStringPointer()
	}
	dtc.RealtimeUrlCheckMode = model.TpspRealtimeURLCheckMode.ValueBoolPointer()
	if !model.TpspSafeBrowsingProtectionLevel.IsNull() {
		dtc.SafeBrowsingProtectionLevel = model.TpspSafeBrowsingProtectionLevel.ValueStringPointer()
	}
	dtc.ScreenLockSecured = model.TpspScreenLockSecured.ValueBoolPointer()
	dtc.SiteIsolationEnabled = model.TpspSiteIsolationEnabled.ValueBoolPointer()
	thirdPartySignalProviders.SetDtc(dtc)
	chromeOS.SetThirdPartySignalProviders(thirdPartySignalProviders)

	if model.GracePeriod != nil {
		gp := v6okta.NewGracePeriod()
		gp.SetType(model.GracePeriod.Type.ValueString())
		gp.SetExpiry(v6okta.StringAsGracePeriodExpiry(model.GracePeriod.Expiry.ValueStringPointer()))
		chromeOS.SetGracePeriod(*gp)
	}
	if !model.DisplayRemediationMode.IsNull() {
		chromeOS.SetDisplayRemediationMode(model.DisplayRemediationMode.ValueString())
	}

	return v6okta.ListDeviceAssurancePolicies200ResponseInner{DeviceAssuranceChromeOSPlatform: chromeOS}, nil
}

// Map response body to schema
func mapDeviceAssuranceChromeOSToState(data *v6okta.ListDeviceAssurancePolicies200ResponseInner, state *policyDeviceAssuranceChromeOSResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	if data.DeviceAssuranceChromeOSPlatform == nil {
		diags.AddError("Empty response", "ChromeOS object")
		return diags
	}

	state.ID = types.StringPointerValue(data.DeviceAssuranceChromeOSPlatform.Id)
	state.Name = types.StringPointerValue(data.DeviceAssuranceChromeOSPlatform.Name)
	state.Platform = types.StringPointerValue(data.DeviceAssuranceChromeOSPlatform.Platform)

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
			state.TpspKeyTrustLevel = types.StringPointerValue(data.DeviceAssuranceChromeOSPlatform.ThirdPartySignalProviders.Dtc.KeyTrustLevel)
		}
		state.TpspOsFirewall = types.BoolPointerValue(data.DeviceAssuranceChromeOSPlatform.ThirdPartySignalProviders.Dtc.OsFirewall)
		if _, ok := data.DeviceAssuranceChromeOSPlatform.ThirdPartySignalProviders.Dtc.GetOsVersionOk(); ok {
			state.TpspOsVersion = types.StringPointerValue(data.DeviceAssuranceChromeOSPlatform.ThirdPartySignalProviders.Dtc.OsVersion.Minimum)
		}
		state.TpspPasswordProtectionWarningTrigger = types.StringPointerValue(data.DeviceAssuranceChromeOSPlatform.ThirdPartySignalProviders.Dtc.PasswordProtectionWarningTrigger)
		state.TpspRealtimeURLCheckMode = types.BoolPointerValue(data.DeviceAssuranceChromeOSPlatform.ThirdPartySignalProviders.Dtc.RealtimeUrlCheckMode)
		state.TpspSafeBrowsingProtectionLevel = types.StringPointerValue(data.DeviceAssuranceChromeOSPlatform.ThirdPartySignalProviders.Dtc.SafeBrowsingProtectionLevel)
		state.TpspScreenLockSecured = types.BoolPointerValue(data.DeviceAssuranceChromeOSPlatform.ThirdPartySignalProviders.Dtc.ScreenLockSecured)
		state.TpspSiteIsolationEnabled = types.BoolPointerValue(data.DeviceAssuranceChromeOSPlatform.ThirdPartySignalProviders.Dtc.SiteIsolationEnabled)
	}

	if gp, ok := data.DeviceAssuranceChromeOSPlatform.GetGracePeriodOk(); ok && gp != nil {
		state.GracePeriod = &gracePeriodModel{
			Type: types.StringPointerValue(gp.Type),
		}
		if gp.Expiry != nil {
			if s := gp.Expiry.String; s != nil {
				state.GracePeriod.Expiry = types.StringPointerValue(s)
			} else if t := gp.Expiry.TimeTime; t != nil {
				state.GracePeriod.Expiry = types.StringValue(t.Format("2006-01-02T15:04:05.000Z07:00"))
			} else {
				state.GracePeriod.Expiry = types.StringNull()
			}
		} else {
			state.GracePeriod.Expiry = types.StringNull()
		}
	}
	state.DisplayRemediationMode = types.StringPointerValue(data.DeviceAssuranceChromeOSPlatform.DisplayRemediationMode)

	state.CreateDate = types.StringPointerValue(data.DeviceAssuranceChromeOSPlatform.CreatedDate)
	state.CreateBy = types.StringPointerValue(data.DeviceAssuranceChromeOSPlatform.CreatedBy)
	state.LastUpdate = types.StringPointerValue(data.DeviceAssuranceChromeOSPlatform.LastUpdate)
	state.LastUpdatedBy = types.StringPointerValue(data.DeviceAssuranceChromeOSPlatform.LastUpdatedBy)
	return diags
}

func (r *policyDeviceAssuranceChromeOSResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
