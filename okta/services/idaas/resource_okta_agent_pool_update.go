package idaas

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	v6okta "github.com/okta/okta-sdk-golang/v6/okta"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ resource.Resource                = &agentPoolUpdateResource{}
	_ resource.ResourceWithConfigure   = &agentPoolUpdateResource{}
	_ resource.ResourceWithImportState = &agentPoolUpdateResource{}
)

func newAgentPoolUpdateResource() resource.Resource {
	return &agentPoolUpdateResource{}
}

type agentPoolUpdateResource struct {
	*config.Config
}

func (r *agentPoolUpdateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *agentPoolUpdateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_agent_pool_update"
}

func (r *agentPoolUpdateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid ID",
			"Expected format: request_id/entry_id",
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("pool_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}

type Agent struct {
	ID     types.String `tfsdk:"id"`
	PoolId types.String `tfsdk:"pool_id"`
}

type UpdateSchedule struct {
	Cron        types.String `tfsdk:"cron"`
	Delay       types.Int32  `tfsdk:"delay"`
	Duration    types.Int32  `tfsdk:"duration"`
	LastUpdated types.String `tfsdk:"last_updated"`
	Timezone    types.String `tfsdk:"timezone"`
}

type agentPoolUpdateResourceModel struct {
	ID                 types.String    `tfsdk:"id"`
	PoolID             types.String    `tfsdk:"pool_id"`
	Name               types.String    `tfsdk:"name"`
	AgentType          types.String    `tfsdk:"agent_type"`
	Enabled            types.Bool      `tfsdk:"enabled"`
	NotifyAdmins       types.Bool      `tfsdk:"notify_admins"`
	Reason             types.String    `tfsdk:"reason"`
	Schedule           *UpdateSchedule `tfsdk:"schedule"`
	SortOrder          types.Int64     `tfsdk:"sort_order"`
	TargetVersion      types.String    `tfsdk:"target_version"`
	Status             types.String    `tfsdk:"status"`
	Agents             []Agent         `tfsdk:"agents"`
	Description        types.String    `tfsdk:"description"`
	NotifyOnCompletion types.Bool      `tfsdk:"notify_on_completion"`
}

func (r *agentPoolUpdateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Okta Agent Pool Update. Agent pool updates allow you to schedule and manage updates for agent pools.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier of the agent pool update.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"pool_id": schema.StringAttribute{
				Required:    true,
				Description: "The unique identifier of the agent pool to update.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"agent_type": schema.StringAttribute{
				Optional:    true,
				Description: "Agent types that are being monitored (e.g. AD, LDAP, IWA, RADIUS).",
			},
			"enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Indicates if auto-update is enabled for the agent pool.",
			},
			"notify_admins": schema.BoolAttribute{
				Optional:    true,
				Description: "Indicates if the admin is notified about the update.",
			},
			"sort_order": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Specifies the sort order.",
			},
			"target_version": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The agent version to update to.",
			},
			"reason": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Reason for the update.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the agent pool update.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the agent pool update.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the agent pool update (e.g., SCHEDULED, IN_PROGRESS, COMPLETED, FAILED).",
			},
			"notify_on_completion": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to send notifications when the update completes.",
			},
		},
		Blocks: map[string]schema.Block{
			"schedule": schema.SingleNestedBlock{
				Description: "The schedule configuration for the agent pool update.",
				Attributes: map[string]schema.Attribute{
					"cron": schema.StringAttribute{
						Optional:    true,
						Description: "The schedule of the update in cron format.",
					},
					"delay": schema.Int32Attribute{
						Optional:    true,
						Description: "Delay in days.",
					},
					"duration": schema.Int32Attribute{
						Optional:    true,
						Description: "Duration in minutes.",
					},
					"last_updated": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Timestamp when the update finished (only for a successful or failed update, not for a cancelled update). Null is returned if the job hasn't finished once yet.",
					},
					"timezone": schema.StringAttribute{
						Optional:    true,
						Description: "Timezone of where the scheduled job takes place.",
					},
				},
			},
			"agents": schema.ListNestedBlock{
				Description: "The agents associated with the agent pool update.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "The unique identifier of the agent.",
						},
						"pool_id": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Pool ID.",
						},
					},
				},
			},
		},
	}
}

func (r *agentPoolUpdateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan agentPoolUpdateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := buildCreateRequest(plan)

	createAgentPoolUpdateResp, _, err := r.OktaIDaaSClient.OktaSDKClientV6().AgentPoolsAPI.CreateAgentPoolsUpdate(
		ctx,
		plan.PoolID.ValueString(),
	).AgentPoolUpdate(createReq).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating agent pool update",
			fmt.Sprintf("Could not create agent pool update: %s", err.Error()),
		)
		return
	}

	mapResponseToState(createAgentPoolUpdateResp, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *agentPoolUpdateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state agentPoolUpdateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	getAgentPoolsUpdateResp, _, err := r.OktaIDaaSClient.OktaSDKClientV6().AgentPoolsAPI.GetAgentPoolsUpdateInstance(
		ctx,
		state.PoolID.ValueString(),
		state.ID.ValueString(),
	).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading agent pool update",
			fmt.Sprintf("Could not read agent pool update %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	mapResponseToState(getAgentPoolsUpdateResp, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *agentPoolUpdateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state agentPoolUpdateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := buildUpdateRequest(plan)

	updateAgentPoolsUpdateResp, _, err := r.OktaIDaaSClient.OktaSDKClientV6().AgentPoolsAPI.UpdateAgentPoolsUpdate(
		ctx,
		state.PoolID.ValueString(),
		state.ID.ValueString(),
	).AgentPoolUpdate(updateReq).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating agent pool update",
			fmt.Sprintf("Could not update agent pool update %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	mapResponseToState(updateAgentPoolsUpdateResp, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *agentPoolUpdateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state agentPoolUpdateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.OktaIDaaSClient.OktaSDKClientV5().AgentPoolsAPI.DeleteAgentPoolsUpdate(
		ctx,
		state.PoolID.ValueString(),
		state.ID.ValueString(),
	).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting agent pool update",
			fmt.Sprintf("Could not delete agent pool update %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}
}

func buildCreateRequest(plan agentPoolUpdateResourceModel) v6okta.AgentPoolUpdate {
	req := v6okta.AgentPoolUpdate{}

	if !plan.Name.IsNull() {
		req.Name = plan.Name.ValueStringPointer()
	}

	if !plan.AgentType.IsNull() {
		req.AgentType = plan.AgentType.ValueStringPointer()
	}

	if !plan.Enabled.IsNull() {
		req.Enabled = plan.Enabled.ValueBoolPointer()
	}

	if !plan.NotifyAdmins.IsNull() {
		req.NotifyAdmin = plan.NotifyAdmins.ValueBoolPointer()
	}

	if !plan.Reason.IsNull() {
		req.Reason = plan.Reason.ValueStringPointer()
	}

	if !plan.SortOrder.IsNull() {
		sortOrder := int32(plan.SortOrder.ValueInt64())
		req.SortOrder = &sortOrder
	}

	if !plan.TargetVersion.IsNull() {
		req.TargetVersion = plan.TargetVersion.ValueStringPointer()
	}

	// Handle Schedule
	if plan.Schedule != nil {
		schedule := v6okta.AutoUpdateSchedule{}

		if !plan.Schedule.Cron.IsNull() {
			schedule.Cron = plan.Schedule.Cron.ValueStringPointer()
		}

		if !plan.Schedule.Delay.IsNull() {
			delay := plan.Schedule.Delay.ValueInt32()
			schedule.Delay = &delay
		}

		if !plan.Schedule.Duration.IsNull() {
			duration := plan.Schedule.Duration.ValueInt32()
			schedule.Duration = &duration
		}

		if !plan.Schedule.Timezone.IsNull() {
			schedule.Timezone = plan.Schedule.Timezone.ValueStringPointer()
		}

		req.Schedule = &schedule
	}

	// Handle Agents - convert from plan to API format
	if len(plan.Agents) > 0 {
		agents := make([]v6okta.Agent, len(plan.Agents))
		for i, agent := range plan.Agents {
			apiAgent := v6okta.Agent{}
			apiAgent.Id = agent.ID.ValueStringPointer()
			apiAgent.PoolId = agent.PoolId.ValueStringPointer()
			agents[i] = apiAgent
		}
		req.Agents = agents
	}

	return req
}

func buildUpdateRequest(plan agentPoolUpdateResourceModel) v6okta.AgentPoolUpdate {
	// Same as create request for this resource
	return buildCreateRequest(plan)
}

func mapResponseToState(resp *v6okta.AgentPoolUpdate, state *agentPoolUpdateResourceModel) {
	state.ID = types.StringValue(resp.GetId())
	state.Name = types.StringValue(resp.GetName())
	state.AgentType = types.StringValue(resp.GetAgentType())
	state.Enabled = types.BoolValue(resp.GetEnabled())
	state.NotifyAdmins = types.BoolValue(resp.GetNotifyAdmin())
	state.Reason = types.StringValue(resp.GetReason())
	state.SortOrder = types.Int64Value(int64(resp.GetSortOrder()))
	state.TargetVersion = types.StringValue(resp.GetTargetVersion())
	state.Status = types.StringValue(resp.GetStatus())

	// Handle Schedule
	s := resp.Schedule
	if s != nil {
		schedule := &UpdateSchedule{}
		schedule.Cron = types.StringValue(s.GetCron())
		schedule.Delay = types.Int32Value((s.GetDelay()))
		schedule.Duration = types.Int32Value((s.GetDuration()))
		schedule.LastUpdated = types.StringValue(s.GetLastUpdated().Format(time.RFC3339))
		schedule.Timezone = types.StringValue(s.GetTimezone())
		state.Schedule = schedule
	}

	// Handle Agents
	if len(resp.Agents) > 0 {
		agents := make([]Agent, len(resp.Agents))
		for i, agent := range resp.Agents {
			tfAgent := Agent{}
			tfAgent.ID = types.StringValue(agent.GetId())
			tfAgent.PoolId = types.StringValue(agent.GetPoolId())
			agents[i] = tfAgent
		}
		state.Agents = agents
	}
}
