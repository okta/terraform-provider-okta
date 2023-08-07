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
	_ resource.Resource              = &policyDeviceAssuranceIOSResource{}
	_ resource.ResourceWithConfigure = &policyDeviceAssuranceIOSResource{}
	// _ resource.ResourceWithImportState = &policyDeviceAssuranceResource{}
)

func NewPolicyDeviceAssuranceIOSResource() resource.Resource {
	return &policyDeviceAssuranceIOSResource{}
}

type policyDeviceAssuranceIOSResource struct {
	*Config
}

type policyDeviceAssuranceIOSResourceModel struct {
	ID             types.String   `tfsdk:"id"`
	Name           types.String   `tfsdk:"name"`
	Platform       types.String   `tfsdk:"platform"`
	JailBreak      types.Bool     `tfsdk:"jailbreak"`
	OsVersion      types.String   `tfsdk:"os_version"`
	ScreenLockType []types.String `tfsdk:"screenlock_type"`
	CreateDate     types.String   `tfsdk:"created_date"`
	CreateBy       types.String   `tfsdk:"created_by"`
	LastUpdate     types.String   `tfsdk:"last_update"`
	LastUpdatedBy  types.String   `tfsdk:"last_updated_by"`
}

func (r *policyDeviceAssuranceIOSResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_device_assurance_ios"
}

func (r *policyDeviceAssuranceIOSResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			},
			"jailbreak": schema.BoolAttribute{
				Description: "The device jailbreak. Only for android and iOS platform",
				Optional:    true,
				Validators: []validator.Bool{
					boolvalidator.AtLeastOneOf(path.Expressions{
						path.MatchRoot("os_version"),
						path.MatchRoot("screenlock_type"),
					}...),
				},
			},
			"os_version": schema.StringAttribute{
				Description: "The device os minimum version",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.AtLeastOneOf(path.Expressions{
						path.MatchRoot("jailbreak"),
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
						path.MatchRoot("jailbreak"),
						path.MatchRoot("os_version"),
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
func (r *policyDeviceAssuranceIOSResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *policyDeviceAssuranceIOSResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state policyDeviceAssuranceIOSResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqBody, err := buildDeviceAssuranceIOSPolicyRequest(state)
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

	resp.Diagnostics.Append(mapDeviceAssuranceIOSToState(deviceAssurance, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *policyDeviceAssuranceIOSResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state policyDeviceAssuranceIOSResourceModel
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

	resp.Diagnostics.Append(mapDeviceAssuranceIOSToState(deviceAssurance, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *policyDeviceAssuranceIOSResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state policyDeviceAssuranceIOSResourceModel
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

func (r *policyDeviceAssuranceIOSResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state policyDeviceAssuranceIOSResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqBody, err := buildDeviceAssuranceIOSPolicyRequest(state)
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

	resp.Diagnostics.Append(mapDeviceAssuranceIOSToState(deviceAssurance, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func buildDeviceAssuranceIOSPolicyRequest(model policyDeviceAssuranceIOSResourceModel) (okta.ListDeviceAssurancePolicies200ResponseInner, error) {
	var iOS = &okta.DeviceAssuranceIOSPlatform{}
	iOS.SetName(model.Name.ValueString())
	iOS.SetPlatform(okta.PLATFORM_IOS)

	iOS.Jailbreak = model.JailBreak.ValueBoolPointer()
	if !model.OsVersion.IsNull() {
		iOS.OsVersion = &okta.OSVersion{Minimum: model.OsVersion.ValueStringPointer()}
	}
	if len(model.ScreenLockType) > 0 {
		screenlockType := make([]okta.ScreenLockType, 0)
		for _, det := range model.ScreenLockType {
			v, err := okta.NewScreenLockTypeFromValue(det.ValueString())
			if err != nil {
				return okta.ListDeviceAssurancePolicies200ResponseInner{DeviceAssuranceIOSPlatform: iOS}, err
			}
			screenlockType = append(screenlockType, *v)
		}
		iOS.ScreenLockType = &okta.DeviceAssuranceAndroidPlatformAllOfScreenLockType{Include: screenlockType}
	}
	return okta.ListDeviceAssurancePolicies200ResponseInner{DeviceAssuranceIOSPlatform: iOS}, nil
}

// Map response body to schema
func mapDeviceAssuranceIOSToState(data *okta.ListDeviceAssurancePolicies200ResponseInner, state *policyDeviceAssuranceIOSResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	if data.DeviceAssuranceIOSPlatform == nil {
		diags.AddError("Empty response", "iOS object")
		return diags
	}

	state.ID = types.StringPointerValue(data.DeviceAssuranceIOSPlatform.Id)
	state.Name = types.StringPointerValue(data.DeviceAssuranceIOSPlatform.Name)
	state.Platform = types.StringPointerValue((*string)(data.DeviceAssuranceIOSPlatform.Platform))

	state.JailBreak = types.BoolPointerValue(data.DeviceAssuranceIOSPlatform.Jailbreak)
	if _, ok := data.DeviceAssuranceIOSPlatform.GetOsVersionOk(); ok {
		state.OsVersion = types.StringPointerValue(data.DeviceAssuranceIOSPlatform.OsVersion.Minimum)
	}
	if _, ok := data.DeviceAssuranceIOSPlatform.ScreenLockType.GetIncludeOk(); ok {
		screenLockType := make([]types.String, 0)
		for _, slt := range data.DeviceAssuranceIOSPlatform.ScreenLockType.GetInclude() {
			screenLockType = append(screenLockType, types.StringValue(string(slt)))
		}
		state.ScreenLockType = screenLockType
	}

	state.CreateDate = types.StringPointerValue(data.DeviceAssuranceIOSPlatform.CreatedDate)
	state.CreateBy = types.StringPointerValue(data.DeviceAssuranceIOSPlatform.CreatedBy)
	state.LastUpdate = types.StringPointerValue(data.DeviceAssuranceIOSPlatform.LastUpdate)
	state.LastUpdatedBy = types.StringPointerValue(data.DeviceAssuranceIOSPlatform.LastUpdatedBy)
	return diags
}
