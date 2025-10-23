package idaas

import (
	"context"
	v5okta "github.com/okta/okta-sdk-golang/v5/okta"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ resource.Resource                = &userLockoutSettings{}
	_ resource.ResourceWithConfigure   = &userLockoutSettings{}
	_ resource.ResourceWithImportState = &userLockoutSettings{}
)

func newUserLockoutSettingsResource() resource.Resource {
	return &userLockoutSettings{}
}

type userLockoutSettings struct {
	*config.Config
}

type userLockoutSettingsResourceModel struct {
	Id                                         types.String `tfsdk:"id"`
	PreventBruteForceLockoutFromUnknownDevices types.Bool   `tfsdk:"prevent_brute_force_lockout_from_unknown_devices"`
}

func (r *userLockoutSettings) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_lockout_settings"
}

func (r *userLockoutSettings) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The id property of an entitlement.",
			},
			"prevent_brute_force_lockout_from_unknown_devices": schema.BoolAttribute{
				Required:    true,
				Description: "Prevents brute-force lockout from unknown devices for the password authenticator.",
			},
		},
	}
}

func (r *userLockoutSettings) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data userLockoutSettingsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic
	createUserLockoutSettingsResp, _, err := r.OktaIDaaSClient.OktaSDKClientV5().AttackProtectionAPI.ReplaceUserLockoutSettings(ctx).LockoutSettings(createLockoutSettings(data)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating user lockout settings",
			"Could not create user lockout settings, unexpected error: "+err.Error(),
		)
		return
	}

	data.Id = types.StringValue("default")
	data.PreventBruteForceLockoutFromUnknownDevices = types.BoolValue(createUserLockoutSettingsResp.GetPreventBruteForceLockoutFromUnknownDevices())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func createLockoutSettings(data userLockoutSettingsResourceModel) v5okta.UserLockoutSettings {
	return v5okta.UserLockoutSettings{
		PreventBruteForceLockoutFromUnknownDevices: data.PreventBruteForceLockoutFromUnknownDevices.ValueBoolPointer(),
	}
}

func (r *userLockoutSettings) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data userLockoutSettingsResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	readUserLockoutSettingsResp, _, err := r.OktaIDaaSClient.OktaSDKClientV5().AttackProtectionAPI.GetUserLockoutSettings(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading user lockout settings",
			"Could not read user lockout settings, unexpected error: "+err.Error(),
		)
		return
	}

	data.PreventBruteForceLockoutFromUnknownDevices = types.BoolValue(readUserLockoutSettingsResp[0].GetPreventBruteForceLockoutFromUnknownDevices())
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *userLockoutSettings) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data userLockoutSettingsResourceModel
	var state userLockoutSettingsResourceModel
	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Id = state.Id
	// Update API call logic
	replaceUserLockoutSettingsResp, _, err := r.OktaIDaaSClient.OktaSDKClientV5().AttackProtectionAPI.ReplaceUserLockoutSettings(ctx).LockoutSettings(createLockoutSettings(data)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating user lockout settings",
			"An error occurred while updating the user lockout settings: "+err.Error(),
		)
		return
	}
	data.PreventBruteForceLockoutFromUnknownDevices = types.BoolValue(replaceUserLockoutSettingsResp.GetPreventBruteForceLockoutFromUnknownDevices())
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *userLockoutSettings) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddError(
		"Delete Not Supported",
		"The resource cannot be deleted via Terraform.",
	)
}

func (r *userLockoutSettings) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *userLockoutSettings) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(request, response)
}
