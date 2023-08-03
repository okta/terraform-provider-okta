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
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Platform      types.String `tfsdk:"platform"`
	CreateDate    types.String `tfsdk:"created_date"`
	CreateBy      types.String `tfsdk:"created_by"`
	LastUpdate    types.String `tfsdk:"last_update"`
	LastUpdatedBy types.String `tfsdk:"last_updated_by"`
	// TODU no feature access
	ThirdPartySignalProviders thirdPartySignalProvidersChromeOS `tfsdk:"third_party_signal_providers"`
}

type thirdPartySignalProvidersChromeOS struct {
	AllowScreenLock                   types.Bool   `tfsdk:"allow_screen_lock"`
	BrowserVersion                    types.String `tfsdk:"browser_version, omitempty"`
	BuiltInDNSClientEnabled           types.Bool   `tfsdk:"builtin_dns_client_enabled"`
	ChromeRemoteDesktopAppBlocked     types.Bool   `tfsdk:"chrome_remote_desktop_app_blocked"`
	DeviceEnrollmentDomain            types.String `tfsdk:"device_enrollment_domain"`
	DiskEncrypted                     types.Bool   `tfsdk:"disk_encrypted"`
	KeyTrustLevel                     types.String `tfsdk:"key_trust_level"`
	OsFirewall                        types.Bool   `tfsdk:"os_firewall"`
	OsVersion                         types.String `tfsdk:"os_version"`
	PasswordProctectionWarningTrigger types.String `tfsdk:"password_proctection_warning_trigger"`
	RealtimeURLCheckMode              types.Bool   `tfsdk:"realtime_url_check_mode"`
	SafeBrowsingProtectionLevel       types.String `tfsdk:"safe_browsing_protection_level"`
	ScreenLockSecured                 types.Bool   `tfsdk:"screen_lock_secured"`
	SiteIsolationEnabled              types.Bool   `tfsdk:"site_isolation_enabled"`
}

func (r *policyDeviceAssuranceChromeOSResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_device_assurance_chromeOS"
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
			// // TODU no access to feature request
			// "third_party_signal_providers": schema.ObjectAttribute{
			// 	Description: "Settings for third-party signal providers. Required for ChromeOS platform, optional for others",
			// 	Required:    true,
			// 	AttributeTypes: map[string]attr.Type{
			// 		"allow_screen_lock":                    types.BoolType,
			// 		"browser_version":                      types.StringType,
			// 		"builtin_dns_client_enabled":           types.BoolType,
			// 		"chrome_remote_desktop_app_blocked":    types.BoolType,
			// 		"device_enrollement_domain":            types.StringType,
			// 		"disk_encrypted":                       types.BoolType,
			// 		"key_trust_level":                      types.StringType,
			// 		"os_firewall":                          types.BoolType,
			// 		"os_version":                           types.StringType,
			// 		"password_proctection_warning_trigger": types.StringType,
			// 		"realtime_url_check_mode":              types.BoolType,
			// 		"safe_browsing_protection_level":       types.StringType,
			// 		"screen_lock_secured":                  types.BoolType,
			// 		"site_isolation_enabled":               types.BoolType,
			// 	},
			// },
			// "created_date": schema.StringAttribute{
			// 	Description: "Created date",
			// 	Computed:    true,
			// },
			// "created_by": schema.StringAttribute{
			// 	Description: "Created by",
			// 	Computed:    true,
			// },
			// "last_update": schema.StringAttribute{
			// 	Description: "Last update",
			// 	Computed:    true,
			// },
			// "last_updated_by": schema.StringAttribute{
			// 	Description: "Last updated by",
			// 	Computed:    true,
			// },
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

	return okta.ListDeviceAssurancePolicies200ResponseInner{DeviceAssuranceChromeOSPlatform: chromeOS}, nil
}

// Map response body to schema
func mapDeviceAssuranceChromeOSToState(data *okta.ListDeviceAssurancePolicies200ResponseInner, state *policyDeviceAssuranceChromeOSResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	state.ID = types.StringValue(data.DeviceAssuranceChromeOSPlatform.GetId())
	state.Name = types.StringValue(data.DeviceAssuranceChromeOSPlatform.GetName())
	state.Platform = types.StringValue(string(data.DeviceAssuranceChromeOSPlatform.GetPlatform()))

	state.CreateDate = types.StringValue(string(data.DeviceAssuranceChromeOSPlatform.GetCreatedDate()))
	state.CreateBy = types.StringValue(string(data.DeviceAssuranceChromeOSPlatform.GetCreatedBy()))
	state.LastUpdate = types.StringValue(string(data.DeviceAssuranceChromeOSPlatform.GetLastUpdate()))
	state.LastUpdatedBy = types.StringValue(string(data.DeviceAssuranceChromeOSPlatform.GetLastUpdatedBy()))
	return diags
}
