package idaas

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	tfpath "github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	v6okta "github.com/okta/okta-sdk-golang/v6/okta"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ resource.Resource                = &appGroupPushMappingResource{}
	_ resource.ResourceWithConfigure   = &appGroupPushMappingResource{}
	_ resource.ResourceWithImportState = &appGroupPushMappingResource{}
)

type appGroupPushMappingResource struct {
	*config.Config
}

func newAppGroupPushMappingResource() resource.Resource {
	return &appGroupPushMappingResource{}
}

type appGroupPushMappingResourceModel struct {
	ID              types.String `tfsdk:"id"`
	AppID           types.String `tfsdk:"app_id"`
	SourceGroupID   types.String `tfsdk:"source_group_id"`
	TargetGroupID   types.String `tfsdk:"target_group_id"`
	TargetGroupName types.String `tfsdk:"target_group_name"`
	Status          types.String `tfsdk:"status"`
	Created         types.String `tfsdk:"created"`
	LastPush        types.String `tfsdk:"last_push"`
	LastUpdated     types.String `tfsdk:"last_updated"`
	ErrorSummary    types.String `tfsdk:"error_summary"`
}

func (r *appGroupPushMappingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_group_push_mapping"
}

func (r *appGroupPushMappingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *appGroupPushMappingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Group Push Mapping for an Okta Application. This resource allows you to push Okta groups to a target application that supports provisioning.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier of the group push mapping.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"app_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the application to push groups to. The application must have provisioning enabled.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"source_group_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the Okta group to push to the application.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"target_group_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The ID of an existing target group in the application to link to. Either target_group_id or target_group_name must be specified, but not both.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"target_group_name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of a new target group to create in the application. Either target_group_id or target_group_name must be specified, but not both.",
			},
			"status": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The status of the group push mapping. Valid values are `ACTIVE` or `INACTIVE`. Defaults to `ACTIVE`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the group push mapping was created.",
			},
			"last_push": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp of the last push operation.",
			},
			"last_updated": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the group push mapping was last modified.",
			},
			"error_summary": schema.StringAttribute{
				Computed:    true,
				Description: "Error message if the latest push operation failed.",
			},
		},
	}
}

func (r *appGroupPushMappingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			"Import ID must be in the format: app_id/mapping_id",
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, tfpath.Root("app_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, tfpath.Root("id"), parts[1])...)
}

func (r *appGroupPushMappingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan appGroupPushMappingResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := v6okta.CreateGroupPushMappingRequest{
		SourceGroupId: plan.SourceGroupID.ValueString(),
	}

	if !plan.TargetGroupID.IsNull() && !plan.TargetGroupID.IsUnknown() {
		createReq.TargetGroupId = plan.TargetGroupID.ValueStringPointer()
	}

	if !plan.TargetGroupName.IsNull() && !plan.TargetGroupName.IsUnknown() {
		createReq.TargetGroupName = plan.TargetGroupName.ValueStringPointer()
	}

	if !plan.Status.IsNull() && !plan.Status.IsUnknown() {
		createReq.Status = plan.Status.ValueStringPointer()
	}

	mapping, _, err := r.OktaIDaaSClient.OktaSDKClientV6().GroupPushMappingAPI.CreateGroupPushMapping(ctx, plan.AppID.ValueString()).Body(createReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Group Push Mapping",
			fmt.Sprintf("Could not create group push mapping for app %s: %s", plan.AppID.ValueString(), err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(applyGroupPushMappingToState(mapping, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *appGroupPushMappingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state appGroupPushMappingResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	mapping, httpResp, err := r.OktaIDaaSClient.OktaSDKClientV6().GroupPushMappingAPI.GetGroupPushMapping(ctx, state.AppID.ValueString(), state.ID.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Group Push Mapping",
			fmt.Sprintf("Could not read group push mapping %s for app %s: %s", state.ID.ValueString(), state.AppID.ValueString(), err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(applyGroupPushMappingToState(mapping, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *appGroupPushMappingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan appGroupPushMappingResourceModel
	var state appGroupPushMappingResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := v6okta.UpdateGroupPushMappingRequest{}

	if !plan.Status.IsNull() && !plan.Status.IsUnknown() {
		updateReq.Status = plan.Status.ValueString()
	}

	mapping, _, err := r.OktaIDaaSClient.OktaSDKClientV6().GroupPushMappingAPI.UpdateGroupPushMapping(ctx, state.AppID.ValueString(), state.ID.ValueString()).Body(updateReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Group Push Mapping",
			fmt.Sprintf("Could not update group push mapping %s for app %s: %s", state.ID.ValueString(), state.AppID.ValueString(), err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(applyGroupPushMappingToState(mapping, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *appGroupPushMappingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state appGroupPushMappingResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.OktaIDaaSClient.OktaSDKClientV6().GroupPushMappingAPI.DeleteGroupPushMapping(ctx, state.AppID.ValueString(), state.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Group Push Mapping",
			fmt.Sprintf("Could not delete group push mapping %s for app %s: %s", state.ID.ValueString(), state.AppID.ValueString(), err.Error()),
		)
		return
	}
}

func applyGroupPushMappingToState(mapping *v6okta.GroupPushMapping, model *appGroupPushMappingResourceModel) (diags diag.Diagnostics) {
	if mapping.Id != nil {
		model.ID = types.StringValue(*mapping.Id)
	}

	if mapping.SourceGroupId != nil {
		model.SourceGroupID = types.StringValue(*mapping.SourceGroupId)
	}

	if mapping.TargetGroupId != nil {
		model.TargetGroupID = types.StringValue(*mapping.TargetGroupId)
	}

	if mapping.Status != nil {
		model.Status = types.StringValue(*mapping.Status)
	}

	if mapping.Created != nil {
		model.Created = types.StringValue(mapping.Created.String())
	}

	if mapping.LastPush != nil {
		model.LastPush = types.StringValue(mapping.LastPush.String())
	}

	if mapping.LastUpdated != nil {
		model.LastUpdated = types.StringValue(mapping.LastUpdated.String())
	}

	if mapping.ErrorSummary != nil {
		model.ErrorSummary = types.StringValue(*mapping.ErrorSummary)
	}

	return
}
