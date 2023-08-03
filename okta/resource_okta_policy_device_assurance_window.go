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
	_ resource.Resource              = &policyDeviceAssuranceWindowsResource{}
	_ resource.ResourceWithConfigure = &policyDeviceAssuranceWindowsResource{}
	// _ resource.ResourceWithImportState = &policyDeviceAssuranceResource{}
)

func NewPolicyDeviceAssuranceWindowsResource() resource.Resource {
	return &policyDeviceAssuranceWindowsResource{}
}

type policyDeviceAssuranceWindowsResource struct {
	*Config
}

type policyDeviceAssuranceWindowsResourceModel struct {
	ID                    types.String   `tfsdk:"id"`
	Name                  types.String   `tfsdk:"name"`
	Platform              types.String   `tfsdk:"platform"`
	DiskEncryptionType    []types.String `tfsdk:"disk_encryption_type"`
	OsVersion             types.String   `tfsdk:"os_version"`
	SecureHardwarePresent types.Bool     `tfsdk:"secure_hardware_present"`
	ScreenLockType        []types.String `tfsdk:"screenlock_type"`
	CreateDate            types.String   `tfsdk:"created_date"`
	CreateBy              types.String   `tfsdk:"created_by"`
	LastUpdate            types.String   `tfsdk:"last_update"`
	LastUpdatedBy         types.String   `tfsdk:"last_updated_by"`
	// // TODU no access to feature request
	// ThirdPartySignalProviders thirdPartySignalProviders `tfsdk:"third_party_signal_providers"`
}

func (r *policyDeviceAssuranceWindowsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_device_assurance_windows"
}

func (r *policyDeviceAssuranceWindowsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			// // TODU no access to feature request
			// "third_party_signal_providers": schema.ObjectAttribute{
			// 	Description: "Settings for third-party signal providers. Required for ChromeOS platform, optional for others",
			// 	Optional:    true,
			// 	AttributeTypes: map[string]attr.Type{
			// 		// TODU chromeOS only
			// 		"allow_screen_lock":                 types.BoolType,
			// 		"browser_version":                   types.StringType,
			// 		"builtin_dns_client_enabled":        types.BoolType,
			// 		"chrome_remote_desktop_app_blocked": types.BoolType,
			// 		// TODU window only
			// 		"crowd_strike_agent_id": types.StringType,
			// 		// TODU window only
			// 		"crowd_strike_customer_id":             types.StringType,
			// 		"device_enrollement_domain":            types.StringType,
			// 		"disk_encrypted":                       types.BoolType,
			// 		"key_trust_level":                      types.StringType,
			// 		"os_firewall":                          types.BoolType,
			// 		"os_version":                           types.StringType,
			// 		"password_proctection_warning_trigger": types.StringType,
			// 		"realtime_url_check_mode":              types.BoolType,
			// 		"safe_browsing_protection_level":       types.StringType,
			// 		"screen_lock_secured":                  types.BoolType,
			// 		// TODU window only
			// 		"secure_boot_enabled":    types.BoolType,
			// 		"site_isolation_enabled": types.BoolType,
			// 		// TODU window only
			// 		"third_party_blocking_enabled": types.BoolType,
			// 		// TODU window only
			// 		"window_machine_domain": types.StringType,
			// 		// TODU window only
			// 		"window_user_domain": types.StringType,
			// 	},
			// },
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

	deviceAssurance, _, err := r.v3Client.DeviceAssuranceApi.CreateDeviceAssurancePolicy(ctx).DeviceAssurance(reqBody).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to create device assurance",
			err.Error(),
		)
		return
	}

	// TODU need to do additional read?
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

	deviceAssurance, _, err := r.v3Client.DeviceAssuranceApi.GetDeviceAssurancePolicy(ctx, state.ID.ValueString()).Execute()
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

	_, err := r.v3Client.DeviceAssuranceApi.DeleteDeviceAssurancePolicy(ctx, state.ID.ValueString()).Execute()
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

	deviceAssurance, _, err := r.v3Client.DeviceAssuranceApi.ReplaceDeviceAssurancePolicy(ctx, state.ID.ValueString()).DeviceAssurance(reqBody).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to create device assurance",
			err.Error(),
		)
		return
	}

	// TODU need to do additional read?
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
	var windows = &okta.DeviceAssuranceWindowsPlatform{}
	windows.SetName(model.Name.ValueString())
	windows.SetPlatform(okta.PLATFORM_WINDOWS)
	if len(model.DiskEncryptionType) > 0 {
		diskEncryptionType := make([]okta.DiskEncryptionType, 0)
		for _, det := range model.DiskEncryptionType {
			v, err := okta.NewDiskEncryptionTypeFromValue(det.ValueString())
			if err != nil {
				return okta.ListDeviceAssurancePolicies200ResponseInner{DeviceAssuranceWindowsPlatform: windows}, err
			}
			diskEncryptionType = append(diskEncryptionType, *v)
		}
		windows.DiskEncryptionType = &okta.DeviceAssuranceAndroidPlatformAllOfDiskEncryptionType{Include: diskEncryptionType}
	}
	if !model.OsVersion.IsNull() {
		windows.OsVersion = &okta.OSVersion{Minimum: model.OsVersion.ValueStringPointer()}
	}
	if len(model.ScreenLockType) > 0 {
		screenlockType := make([]okta.ScreenLockType, 0)
		for _, det := range model.ScreenLockType {
			v, err := okta.NewScreenLockTypeFromValue(det.ValueString())
			if err != nil {
				return okta.ListDeviceAssurancePolicies200ResponseInner{DeviceAssuranceWindowsPlatform: windows}, err
			}
			screenlockType = append(screenlockType, *v)
		}
		windows.ScreenLockType = &okta.DeviceAssuranceAndroidPlatformAllOfScreenLockType{Include: screenlockType}
	}
	windows.SecureHardwarePresent = model.SecureHardwarePresent.ValueBoolPointer()
	return okta.ListDeviceAssurancePolicies200ResponseInner{DeviceAssuranceWindowsPlatform: windows}, nil
}

// Map response body to schema
func mapDeviceAssuranceWindowsToState(data *okta.ListDeviceAssurancePolicies200ResponseInner, state *policyDeviceAssuranceWindowsResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	state.ID = types.StringValue(data.DeviceAssuranceWindowsPlatform.GetId())
	state.Name = types.StringValue(data.DeviceAssuranceWindowsPlatform.GetName())
	state.Platform = types.StringValue(string(data.DeviceAssuranceWindowsPlatform.GetPlatform()))

	if _, ok := data.DeviceAssuranceWindowsPlatform.GetSecureHardwarePresentOk(); ok {
		state.SecureHardwarePresent = types.BoolValue(data.DeviceAssuranceWindowsPlatform.GetSecureHardwarePresent())
	}
	if _, ok := data.DeviceAssuranceWindowsPlatform.GetOsVersionOk(); ok {
		state.OsVersion = types.StringValue(data.DeviceAssuranceWindowsPlatform.OsVersion.GetMinimum())
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

	state.CreateDate = types.StringValue(string(data.DeviceAssuranceWindowsPlatform.GetCreatedDate()))
	state.CreateBy = types.StringValue(string(data.DeviceAssuranceWindowsPlatform.GetCreatedBy()))
	state.LastUpdate = types.StringValue(string(data.DeviceAssuranceWindowsPlatform.GetLastUpdate()))
	state.LastUpdatedBy = types.StringValue(string(data.DeviceAssuranceWindowsPlatform.GetLastUpdatedBy()))
	return diags
}
