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
	_ resource.Resource              = &policyDeviceAssuranceAndroidResource{}
	_ resource.ResourceWithConfigure = &policyDeviceAssuranceAndroidResource{}
	// _ resource.ResourceWithImportState = &policyDeviceAssuranceResource{}
)

func NewPolicyDeviceAssuranceAndroidResource() resource.Resource {
	return &policyDeviceAssuranceAndroidResource{}
}

type policyDeviceAssuranceAndroidResource struct {
	*Config
}

type policyDeviceAssuranceAndroidResourceModel struct {
	ID                    types.String   `tfsdk:"id"`
	Name                  types.String   `tfsdk:"name"`
	Platform              types.String   `tfsdk:"platform"`
	DiskEncryptionType    []types.String `tfsdk:"disk_encryption_type"`
	JailBreak             types.Bool     `tfsdk:"jailbreak"`
	OsVersion             types.String   `tfsdk:"os_version"`
	SecureHardwarePresent types.Bool     `tfsdk:"secure_hardware_present"`
	ScreenLockType        []types.String `tfsdk:"screenlock_type"`
	CreateDate            types.String   `tfsdk:"created_date"`
	CreateBy              types.String   `tfsdk:"created_by"`
	LastUpdate            types.String   `tfsdk:"last_update"`
	LastUpdatedBy         types.String   `tfsdk:"last_updated_by"`
}

func (r *policyDeviceAssuranceAndroidResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_device_assurance_android"
}

func (r *policyDeviceAssuranceAndroidResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
				Description: "List of disk encryption type, can be FULL, USER",
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.AtLeastOneOf(path.Expressions{
						path.MatchRoot("jailbreak"),
						path.MatchRoot("os_version"),
						path.MatchRoot("secure_hardware_present"),
						path.MatchRoot("screenlock_type"),
					}...),
				},
			},
			"jailbreak": schema.BoolAttribute{
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
				Description: "The device os minimum version",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.AtLeastOneOf(path.Expressions{
						path.MatchRoot("disk_encryption_type"),
						path.MatchRoot("jailbreak"),
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
						path.MatchRoot("jailbreak"),
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
						path.MatchRoot("jailbreak"),
						path.MatchRoot("os_version"),
						path.MatchRoot("secure_hardware_present"),
					}...),
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

// Configure adds the provider configured client to the resource.
func (r *policyDeviceAssuranceAndroidResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *policyDeviceAssuranceAndroidResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state policyDeviceAssuranceAndroidResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqBody, err := buildDeviceAssuranceAndroidPolicyRequest(state)
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
	resp.Diagnostics.Append(mapDeviceAssuranceAndroidToState(deviceAssurance, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *policyDeviceAssuranceAndroidResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state policyDeviceAssuranceAndroidResourceModel
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

	resp.Diagnostics.Append(mapDeviceAssuranceAndroidToState(deviceAssurance, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *policyDeviceAssuranceAndroidResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state policyDeviceAssuranceAndroidResourceModel
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

func (r *policyDeviceAssuranceAndroidResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state policyDeviceAssuranceAndroidResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqBody, err := buildDeviceAssuranceAndroidPolicyRequest(state)
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
	resp.Diagnostics.Append(mapDeviceAssuranceAndroidToState(deviceAssurance, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func buildDeviceAssuranceAndroidPolicyRequest(model policyDeviceAssuranceAndroidResourceModel) (okta.ListDeviceAssurancePolicies200ResponseInner, error) {
	var android = &okta.DeviceAssuranceAndroidPlatform{}
	android.SetName(model.Name.ValueString())
	android.SetPlatform(okta.PLATFORM_ANDROID)
	if len(model.DiskEncryptionType) > 0 {
		diskEncryptionType := make([]okta.DiskEncryptionType, 0)
		for _, det := range model.DiskEncryptionType {
			v, err := okta.NewDiskEncryptionTypeFromValue(det.ValueString())
			if err != nil {
				return okta.ListDeviceAssurancePolicies200ResponseInner{DeviceAssuranceAndroidPlatform: android}, err
			}
			diskEncryptionType = append(diskEncryptionType, *v)
		}
		android.DiskEncryptionType = &okta.DeviceAssuranceAndroidPlatformAllOfDiskEncryptionType{Include: diskEncryptionType}
	}
	android.Jailbreak = model.JailBreak.ValueBoolPointer()
	if !model.OsVersion.IsNull() {
		android.OsVersion = &okta.OSVersion{Minimum: model.OsVersion.ValueStringPointer()}
	}
	if len(model.ScreenLockType) > 0 {
		screenlockType := make([]okta.ScreenLockType, 0)
		for _, det := range model.ScreenLockType {
			v, err := okta.NewScreenLockTypeFromValue(det.ValueString())
			if err != nil {
				return okta.ListDeviceAssurancePolicies200ResponseInner{DeviceAssuranceAndroidPlatform: android}, err
			}
			screenlockType = append(screenlockType, *v)
		}
		android.ScreenLockType = &okta.DeviceAssuranceAndroidPlatformAllOfScreenLockType{Include: screenlockType}
	}
	android.SecureHardwarePresent = model.SecureHardwarePresent.ValueBoolPointer()
	return okta.ListDeviceAssurancePolicies200ResponseInner{DeviceAssuranceAndroidPlatform: android}, nil
}

// Map response body to schema
func mapDeviceAssuranceAndroidToState(data *okta.ListDeviceAssurancePolicies200ResponseInner, state *policyDeviceAssuranceAndroidResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	if data.DeviceAssuranceAndroidPlatform == nil {
		diags.AddError("Empty response", "Android object")
		return diags
	}
	state.ID = types.StringPointerValue(data.DeviceAssuranceAndroidPlatform.Id)
	state.Name = types.StringPointerValue(data.DeviceAssuranceAndroidPlatform.Name)
	state.Platform = types.StringPointerValue((*string)(data.DeviceAssuranceAndroidPlatform.Platform))

	state.JailBreak = types.BoolPointerValue(data.DeviceAssuranceAndroidPlatform.Jailbreak)
	state.SecureHardwarePresent = types.BoolPointerValue(data.DeviceAssuranceAndroidPlatform.SecureHardwarePresent)
	if _, ok := data.DeviceAssuranceAndroidPlatform.GetOsVersionOk(); ok {
		state.OsVersion = types.StringPointerValue(data.DeviceAssuranceAndroidPlatform.OsVersion.Minimum)
	}
	if _, ok := data.DeviceAssuranceAndroidPlatform.DiskEncryptionType.GetIncludeOk(); ok {
		diskEncryptionType := make([]types.String, 0)
		for _, det := range data.DeviceAssuranceAndroidPlatform.DiskEncryptionType.GetInclude() {
			diskEncryptionType = append(diskEncryptionType, types.StringValue(string(det)))
		}
		state.DiskEncryptionType = diskEncryptionType
	}
	if _, ok := data.DeviceAssuranceAndroidPlatform.ScreenLockType.GetIncludeOk(); ok {
		screenLockType := make([]types.String, 0)
		for _, slt := range data.DeviceAssuranceAndroidPlatform.ScreenLockType.GetInclude() {
			screenLockType = append(screenLockType, types.StringValue(string(slt)))
		}
		state.ScreenLockType = screenLockType
	}

	state.CreateDate = types.StringPointerValue(data.DeviceAssuranceAndroidPlatform.CreatedDate)
	state.CreateBy = types.StringPointerValue(data.DeviceAssuranceAndroidPlatform.CreatedBy)
	state.LastUpdate = types.StringPointerValue(data.DeviceAssuranceAndroidPlatform.LastUpdate)
	state.LastUpdatedBy = types.StringPointerValue(data.DeviceAssuranceAndroidPlatform.LastUpdatedBy)
	return diags
}
