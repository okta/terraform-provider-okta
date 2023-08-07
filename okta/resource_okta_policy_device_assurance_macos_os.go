package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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
	_ resource.Resource              = &policyDeviceAssuranceMacOSResource{}
	_ resource.ResourceWithConfigure = &policyDeviceAssuranceMacOSResource{}
	// _ resource.ResourceWithImportState = &policyDeviceAssuranceResource{}
)

func NewPolicyDeviceAssuranceMacOSResource() resource.Resource {
	return &policyDeviceAssuranceMacOSResource{}
}

type policyDeviceAssuranceMacOSResource struct {
	*Config
}

type policyDeviceAssuranceMacOSResourceModel struct {
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
	TpspDeviceEnrollmentDomain           types.String   `tfsdk:"tpsp_device_enrollment_domain"`
	TpspDiskEncrypted                    types.Bool     `tfsdk:"tpsp_disk_encrypted"`
	TpspKeyTrustLevel                    types.String   `tfsdk:"tpsp_key_trust_level"`
	TpspOsFirewall                       types.Bool     `tfsdk:"tpsp_os_firewall"`
	TpspOsVersion                        types.String   `tfsdk:"tpsp_os_version"`
	TpspPasswordProtectionWarningTrigger types.String   `tfsdk:"tpsp_password_proctection_warning_trigger"`
	TpspRealtimeURLCheckMode             types.Bool     `tfsdk:"tpsp_realtime_url_check_mode"`
	TpspSafeBrowsingProtectionLevel      types.String   `tfsdk:"tpsp_safe_browsing_protection_level"`
	TpspScreenLockSecured                types.Bool     `tfsdk:"tpsp_screen_lock_secured"`
	TpspSiteIsolationEnabled             types.Bool     `tfsdk:"tpsp_site_isolation_enabled"`
}

func (r *policyDeviceAssuranceMacOSResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_device_assurance_macos"
}

func (r *policyDeviceAssuranceMacOSResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			// use set to avoid order change as v3 does not have diff suppress func
			"disk_encryption_type": schema.SetAttribute{
				Description: "List of disk encryption type, can be ALL_INTERNAL_VOLUMES",
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.AtLeastOneOf(path.Expressions{
						path.MatchRoot("os_version"),
						path.MatchRoot("secure_hardware_present"),
						path.MatchRoot("screenlock_type"),
					}...),
				},
			},
			"os_version": schema.StringAttribute{
				Description: "The device os minimum version",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.AtLeastOneOf(path.Expressions{
						path.MatchRoot("disk_encryption_type"),
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
						path.MatchRoot("os_version"),
						path.MatchRoot("screenlock_type"),
					}...),
				},
			},
			"screenlock_type": schema.SetAttribute{
				Description: "List of screenlock type, can be BIOMETRIC or BIOMETRIC, PASSCODE",
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.AtLeastOneOf(path.Expressions{
						path.MatchRoot("disk_encryption_type"),
						path.MatchRoot("os_version"),
						path.MatchRoot("secure_hardware_present"),
					}...),
				},
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
func (r *policyDeviceAssuranceMacOSResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *policyDeviceAssuranceMacOSResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state policyDeviceAssuranceMacOSResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqBody, err := buildDeviceAssuranceMacOSPolicyRequest(state)
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

	resp.Diagnostics.Append(mapDeviceAssuranceMacOSToState(deviceAssurance, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *policyDeviceAssuranceMacOSResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state policyDeviceAssuranceMacOSResourceModel
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

	resp.Diagnostics.Append(mapDeviceAssuranceMacOSToState(deviceAssurance, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *policyDeviceAssuranceMacOSResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state policyDeviceAssuranceMacOSResourceModel
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

func (r *policyDeviceAssuranceMacOSResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state policyDeviceAssuranceMacOSResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqBody, err := buildDeviceAssuranceMacOSPolicyRequest(state)
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

	resp.Diagnostics.Append(mapDeviceAssuranceMacOSToState(deviceAssurance, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func buildDeviceAssuranceMacOSPolicyRequest(model policyDeviceAssuranceMacOSResourceModel) (okta.ListDeviceAssurancePolicies200ResponseInner, error) {
	var macos = &okta.DeviceAssuranceMacOSPlatform{}
	macos.SetName(model.Name.ValueString())
	macos.SetPlatform(okta.PLATFORM_MACOS)
	if len(model.DiskEncryptionType) > 0 {
		diskEncryptionType := make([]okta.DiskEncryptionType, 0)
		for _, det := range model.DiskEncryptionType {
			v, err := okta.NewDiskEncryptionTypeFromValue(det.ValueString())
			if err != nil {
				return okta.ListDeviceAssurancePolicies200ResponseInner{DeviceAssuranceMacOSPlatform: macos}, err
			}
			diskEncryptionType = append(diskEncryptionType, *v)
		}
		macos.DiskEncryptionType = &okta.DeviceAssuranceAndroidPlatformAllOfDiskEncryptionType{Include: diskEncryptionType}
	}
	if !model.OsVersion.IsNull() {
		macos.OsVersion = &okta.OSVersion{Minimum: model.OsVersion.ValueStringPointer()}
	}
	if len(model.ScreenLockType) > 0 {
		screenlockType := make([]okta.ScreenLockType, 0)
		for _, det := range model.ScreenLockType {
			v, err := okta.NewScreenLockTypeFromValue(det.ValueString())
			if err != nil {
				return okta.ListDeviceAssurancePolicies200ResponseInner{DeviceAssuranceMacOSPlatform: macos}, err
			}
			screenlockType = append(screenlockType, *v)
		}
		macos.ScreenLockType = &okta.DeviceAssuranceAndroidPlatformAllOfScreenLockType{Include: screenlockType}
	}
	macos.SecureHardwarePresent = model.SecureHardwarePresent.ValueBoolPointer()

	if model.ThirdPartySignalProviders.ValueBool() {
		var thirdPartySignalProviders okta.DeviceAssuranceMacOSPlatformAllOfThirdPartySignalProviders
		var dtc okta.DTCMacOS
		if !model.TpspBrowserVersion.IsNull() {
			dtc.BrowserVersion = &okta.ChromeBrowserVersion{Minimum: model.TpspBrowserVersion.ValueStringPointer()}
		}
		dtc.BuiltInDnsClientEnabled = model.TpspBuiltInDNSClientEnabled.ValueBoolPointer()
		dtc.ChromeRemoteDesktopAppBlocked = model.TpspChromeRemoteDesktopAppBlocked.ValueBoolPointer()
		dtc.DeviceEnrollmentDomain = model.TpspDeviceEnrollmentDomain.ValueStringPointer()
		dtc.DiskEncrypted = model.TpspDiskEncrypted.ValueBoolPointer()
		if !model.TpspKeyTrustLevel.IsNull() {
			v, err := okta.NewKeyTrustLevelBrowserKeyFromValue(model.TpspKeyTrustLevel.ValueString())
			if err != nil {
				return okta.ListDeviceAssurancePolicies200ResponseInner{DeviceAssuranceMacOSPlatform: macos}, err
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
				return okta.ListDeviceAssurancePolicies200ResponseInner{DeviceAssuranceMacOSPlatform: macos}, err
			}
			dtc.PasswordProtectionWarningTrigger = v
		}
		dtc.RealtimeUrlCheckMode = model.TpspRealtimeURLCheckMode.ValueBoolPointer()
		if !model.TpspSafeBrowsingProtectionLevel.IsNull() {
			v, err := okta.NewSafeBrowsingProtectionLevelFromValue(model.TpspSafeBrowsingProtectionLevel.ValueString())
			if err != nil {
				return okta.ListDeviceAssurancePolicies200ResponseInner{DeviceAssuranceMacOSPlatform: macos}, err
			}
			dtc.SafeBrowsingProtectionLevel = v
		}
		dtc.ScreenLockSecured = model.TpspScreenLockSecured.ValueBoolPointer()
		dtc.SiteIsolationEnabled = model.TpspSiteIsolationEnabled.ValueBoolPointer()
		thirdPartySignalProviders.SetDtc(dtc)
		macos.SetThirdPartySignalProviders(thirdPartySignalProviders)
	}

	return okta.ListDeviceAssurancePolicies200ResponseInner{DeviceAssuranceMacOSPlatform: macos}, nil
}

// Map response body to schema
func mapDeviceAssuranceMacOSToState(data *okta.ListDeviceAssurancePolicies200ResponseInner, state *policyDeviceAssuranceMacOSResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	if data.DeviceAssuranceMacOSPlatform == nil {
		diags.AddError("Empty response", "MacOS object")
		return diags
	}

	state.ID = types.StringPointerValue(data.DeviceAssuranceMacOSPlatform.Id)
	state.Name = types.StringPointerValue(data.DeviceAssuranceMacOSPlatform.Name)
	state.Platform = types.StringPointerValue((*string)(data.DeviceAssuranceMacOSPlatform.Platform))

	state.SecureHardwarePresent = types.BoolPointerValue(data.DeviceAssuranceMacOSPlatform.SecureHardwarePresent)
	if _, ok := data.DeviceAssuranceMacOSPlatform.GetOsVersionOk(); ok {
		state.OsVersion = types.StringPointerValue(data.DeviceAssuranceMacOSPlatform.OsVersion.Minimum)
	}
	if _, ok := data.DeviceAssuranceMacOSPlatform.DiskEncryptionType.GetIncludeOk(); ok {
		diskEncryptionType := make([]types.String, 0)
		for _, det := range data.DeviceAssuranceMacOSPlatform.DiskEncryptionType.GetInclude() {
			diskEncryptionType = append(diskEncryptionType, types.StringValue(string(det)))
		}
		state.DiskEncryptionType = diskEncryptionType
	}
	if _, ok := data.DeviceAssuranceMacOSPlatform.ScreenLockType.GetIncludeOk(); ok {
		screenLockType := make([]types.String, 0)
		for _, slt := range data.DeviceAssuranceMacOSPlatform.ScreenLockType.GetInclude() {
			screenLockType = append(screenLockType, types.StringValue(string(slt)))
		}
		state.ScreenLockType = screenLockType
	}

	if _, ok := data.DeviceAssuranceMacOSPlatform.ThirdPartySignalProviders.GetDtcOk(); ok {
		if _, ok := data.DeviceAssuranceMacOSPlatform.ThirdPartySignalProviders.Dtc.GetBrowserVersionOk(); ok {
			state.TpspBrowserVersion = types.StringPointerValue(data.DeviceAssuranceMacOSPlatform.ThirdPartySignalProviders.Dtc.BrowserVersion.Minimum)
		}
		state.TpspBuiltInDNSClientEnabled = types.BoolPointerValue(data.DeviceAssuranceMacOSPlatform.ThirdPartySignalProviders.Dtc.BuiltInDnsClientEnabled)
		state.TpspChromeRemoteDesktopAppBlocked = types.BoolPointerValue(data.DeviceAssuranceMacOSPlatform.ThirdPartySignalProviders.Dtc.ChromeRemoteDesktopAppBlocked)
		state.TpspDeviceEnrollmentDomain = types.StringPointerValue(data.DeviceAssuranceMacOSPlatform.ThirdPartySignalProviders.Dtc.DeviceEnrollmentDomain)
		state.TpspDiskEncrypted = types.BoolPointerValue(data.DeviceAssuranceMacOSPlatform.ThirdPartySignalProviders.Dtc.DiskEncrypted)
		if _, ok := data.DeviceAssuranceMacOSPlatform.ThirdPartySignalProviders.Dtc.GetKeyTrustLevelOk(); ok {
			state.TpspKeyTrustLevel = types.StringPointerValue((*string)(data.DeviceAssuranceMacOSPlatform.ThirdPartySignalProviders.Dtc.KeyTrustLevel))
		}
		state.TpspOsFirewall = types.BoolPointerValue(data.DeviceAssuranceMacOSPlatform.ThirdPartySignalProviders.Dtc.OsFirewall)
		if _, ok := data.DeviceAssuranceMacOSPlatform.ThirdPartySignalProviders.Dtc.GetOsVersionOk(); ok {
			state.TpspOsVersion = types.StringPointerValue(data.DeviceAssuranceMacOSPlatform.ThirdPartySignalProviders.Dtc.OsVersion.Minimum)
		}
		state.TpspPasswordProtectionWarningTrigger = types.StringPointerValue((*string)(data.DeviceAssuranceMacOSPlatform.ThirdPartySignalProviders.Dtc.PasswordProtectionWarningTrigger))
		state.TpspRealtimeURLCheckMode = types.BoolPointerValue(data.DeviceAssuranceMacOSPlatform.ThirdPartySignalProviders.Dtc.RealtimeUrlCheckMode)
		state.TpspSafeBrowsingProtectionLevel = types.StringPointerValue((*string)(data.DeviceAssuranceMacOSPlatform.ThirdPartySignalProviders.Dtc.SafeBrowsingProtectionLevel))
		state.TpspScreenLockSecured = types.BoolPointerValue(data.DeviceAssuranceMacOSPlatform.ThirdPartySignalProviders.Dtc.ScreenLockSecured)
		state.TpspSiteIsolationEnabled = types.BoolPointerValue(data.DeviceAssuranceMacOSPlatform.ThirdPartySignalProviders.Dtc.SiteIsolationEnabled)
	}

	state.CreateDate = types.StringPointerValue(data.DeviceAssuranceMacOSPlatform.CreatedDate)
	state.CreateBy = types.StringPointerValue(data.DeviceAssuranceMacOSPlatform.CreatedBy)
	state.LastUpdate = types.StringPointerValue(data.DeviceAssuranceMacOSPlatform.LastUpdate)
	state.LastUpdatedBy = types.StringPointerValue(data.DeviceAssuranceMacOSPlatform.LastUpdatedBy)
	return diags
}
