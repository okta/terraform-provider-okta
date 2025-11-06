package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	v5okta "github.com/okta/okta-sdk-golang/v5/okta"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ resource.Resource                = &devicesResource{}
	_ resource.ResourceWithConfigure   = &devicesResource{}
	_ resource.ResourceWithImportState = &devicesResource{}
)

func newDevicesResource() resource.Resource {
	return &devicesResource{}
}

type devicesResource struct {
	*config.Config
}

func (r *devicesResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *devicesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *devicesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device"
}

type deviceResourceModel struct {
	Id           types.String `tfsdk:"id"`
	Status       types.String `tfsdk:"status"`
	ResourceType types.String `tfsdk:"resource_type"`
	Action       types.String `tfsdk:"action"`
}

func (r *devicesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the device.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the device.",
			},
			"resource_type": schema.StringAttribute{
				Computed:    true,
				Description: "The resource type of the device.",
			},
			"action": schema.StringAttribute{
				Optional:    true,
				Description: "The action of the device.",
			},
		},
	}
}

func (r *devicesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	resp.Diagnostics.AddError(
		"Create Not Supported",
		"This resource cannot be created via Terraform. Please import it or let Terraform read it from the existing system.",
	)
}

func (r *devicesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data deviceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	getDeviceResp, _, err := r.OktaIDaaSClient.OktaSDKClientV5().DeviceAPI.GetDevice(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Device",
			"Could not read Device with Id "+data.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(applyDevicesToState(getDeviceResp, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *devicesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state deviceResourceModel

	// Read the state from the Terraform configuration
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Action.ValueStringPointer() != nil && data.Action.ValueString() == "ACTIVE" {
		_, err := r.OktaIDaaSClient.OktaSDKClientV5().DeviceAPI.ActivateDevice(ctx, state.Id.ValueString()).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Activating Device",
				"Could not activate device with Id "+state.Id.ValueString()+": "+err.Error(),
			)
			return
		}
	} else if data.Action.ValueStringPointer() != nil && data.Action.ValueString() == "DEACTIVATED" {
		_, err := r.OktaIDaaSClient.OktaSDKClientV5().DeviceAPI.DeactivateDevice(ctx, state.Id.ValueString()).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Deactivating Device",
				"Could not deactivate device with Id "+state.Id.ValueString()+": "+err.Error(),
			)
			return
		}
	} else if data.Action.ValueStringPointer() != nil && data.Action.ValueString() == "SUSPENDED" {
		_, err := r.OktaIDaaSClient.OktaSDKClientV5().DeviceAPI.SuspendDevice(ctx, state.Id.ValueString()).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Suspending Device",
				"Could not suspend device with Id "+state.Id.ValueString()+": "+err.Error(),
			)
			return
		}
	} else if data.Action.ValueStringPointer() != nil && data.Action.ValueString() == "UNSUSPEND" {
		_, err := r.OktaIDaaSClient.OktaSDKClientV5().DeviceAPI.UnsuspendDevice(ctx, state.Id.ValueString()).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error unsuspending the Device",
				"Could not unsuspend device with Id "+state.Id.ValueString()+": "+err.Error(),
			)
			return
		}
	}

	getDeviceResp, _, err := r.OktaIDaaSClient.OktaSDKClientV5().DeviceAPI.GetDevice(ctx, state.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Device",
			"Could not read Device with Id "+data.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(applyDevicesToState(getDeviceResp, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *devicesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data deviceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.OktaIDaaSClient.OktaSDKClientV5().DeviceAPI.DeleteDevice(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Device",
			"Could not delete device with Id "+data.Id.ValueString()+": "+err.Error(),
		)
		return
	}
}

func applyDevicesToState(resp *v5okta.Device, s *deviceResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	s.Id = types.StringValue(resp.GetId())
	s.Status = types.StringValue(resp.GetStatus())
	s.ResourceType = types.StringValue(resp.GetResourceType())
	return diags
}
