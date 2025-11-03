package idaas

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	v5okta "github.com/okta/okta-sdk-golang/v5/okta"
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
}

type Agent struct {
	IsHidden          types.Bool   `tfsdk:"is_hidden"`
	IsLatestGAVersion types.Bool   `tfsdk:"is_latest_ga_version"`
	LastConnection    types.String `tfsdk:"last_connection"`
	Name              types.String `tfsdk:"name"`
	OperationalStatus types.String `tfsdk:"operational_status"`
	PoolId            types.String `tfsdk:"pool_id"`
	Type              types.String `tfsdk:"type"`
	UpdateMessage     types.String `tfsdk:"update_message"`
	UpdateStatus      types.String `tfsdk:"update_status"`
	Version           types.String `tfsdk:"version"`
}

type Schedule struct {
	Cron        types.String `tfsdk:"cron"`
	Delay       types.Int64  `tfsdk:"delay"`
	Duration    types.Int64  `tfsdk:"duration"`
	LastUpdated types.String `tfsdk:"last_updated"`
	Timezone    types.String `tfsdk:"timezone"`
}

type agentPoolUpdateResourceModel struct {
	Id            types.String `tfsdk:"id"`
	PoolId        types.String `tfsdk:"pool_id"`
	Name          types.String `tfsdk:"name"`
	AgentType     types.String `tfsdk:"agent_type"`
	Enabled       types.Bool   `tfsdk:"enabled"`
	NotifyAdmins  types.Bool   `tfsdk:"notify_admins"`
	Reason        types.String `tfsdk:"reason"`
	Schedule      *Schedule    `tfsdk:"schedule"`
	SortOrder     types.Int64  `tfsdk:"sort_order"`
	TargetVersion types.String `tfsdk:"target_version"`
	Status        types.String `tfsdk:"status"`
	Agents        []Agent      `tfsdk:"agents"`
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
			"completed_date": schema.StringAttribute{
				Computed:    true,
				Description: "The date and time when the update was completed, in RFC3339 format.",
			},
			"scheduled_date": schema.StringAttribute{
				Optional:    true,
				Description: "The date and time when the update should be executed, in RFC3339 format.",
			},
			"created_date": schema.StringAttribute{
				Computed:    true,
				Description: "The date and time when the update was created, in RFC3339 format.",
			},
			"last_updated": schema.StringAttribute{
				Computed:    true,
				Description: "The date and time when the update was last modified, in RFC3339 format.",
			},
			"notify_on_completion": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to send notifications when the update completes.",
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

	createReq := r.buildCreateRequest(plan)

	createResp, _, err := r.OktaIDaaSClient.OktaSDKClientV5().AgentPoolsAPI.CreateAgentPoolsUpdate(
		ctx,
		plan.PoolId.ValueString(),
	).AgentPoolUpdate(createReq).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating agent pool update",
			fmt.Sprintf("Could not create agent pool update: %s", err.Error()),
		)
		return
	}

	r.mapResponseToState(createResp, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *agentPoolUpdateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state agentPoolUpdateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	getResp, _, err := r.OktaIDaaSClient.OktaSDKClientV5().AgentPoolsAPI.GetAgentPoolsUpdateInstance(
		ctx,
		state.PoolId.ValueString(),
		state.Id.ValueString(),
	).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading agent pool update",
			fmt.Sprintf("Could not read agent pool update %s: %s", state.Id.ValueString(), err.Error()),
		)
		return
	}

	r.mapResponseToState(getResp, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *agentPoolUpdateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state agentPoolUpdateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := r.buildUpdateRequest(plan)

	updateResp, _, err := r.OktaIDaaSClient.OktaSDKClientV5().AgentPoolsAPI.UpdateAgentPoolsUpdate(
		ctx,
		state.PoolId.ValueString(),
		state.Id.ValueString(),
	).AgentPoolUpdate(updateReq).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating agent pool update",
			fmt.Sprintf("Could not update agent pool update %s: %s", state.Id.ValueString(), err.Error()),
		)
		return
	}

	r.mapResponseToState(updateResp, &plan)
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
		state.PoolId.ValueString(),
		state.Id.ValueString(),
	).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting agent pool update",
			fmt.Sprintf("Could not delete agent pool update %s: %s", state.Id.ValueString(), err.Error()),
		)
		return
	}
}

func (r *agentPoolUpdateResource) buildCreateRequest(plan agentPoolUpdateResourceModel) v5okta.AgentPoolUpdate {
	req := v5okta.AgentPoolUpdate{}

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
		schedule := v5okta.AutoUpdateSchedule{}

		if !plan.Schedule.Cron.IsNull() {
			schedule.Cron = plan.Schedule.Cron.ValueStringPointer()
		}

		if !plan.Schedule.Delay.IsNull() {
			delay := int32(plan.Schedule.Delay.ValueInt64())
			schedule.Delay = &delay
		}

		if !plan.Schedule.Duration.IsNull() {
			duration := int32(plan.Schedule.Duration.ValueInt64())
			schedule.Duration = &duration
		}

		if !plan.Schedule.Timezone.IsNull() {
			schedule.Timezone = plan.Schedule.Timezone.ValueStringPointer()
		}

		req.Schedule = &schedule
	}

	// Handle Agents - convert from plan to API format
	if len(plan.Agents) > 0 {
		agents := make([]v5okta.Agent, len(plan.Agents))
		for i, agent := range plan.Agents {
			apiAgent := v5okta.Agent{}
			if !agent.Name.IsNull() {
				apiAgent.Name = agent.Name.ValueStringPointer()
			}
			if !agent.Type.IsNull() {
				apiAgent.Type = agent.Type.ValueStringPointer()
			}
			if !agent.PoolId.IsNull() {
				apiAgent.PoolId = agent.PoolId.ValueStringPointer()
			}
			if !agent.Version.IsNull() {
				apiAgent.Version = agent.Version.ValueStringPointer()
			}
			if !agent.UpdateStatus.IsNull() {
				apiAgent.UpdateStatus = agent.UpdateStatus.ValueStringPointer()
			}
			if !agent.UpdateMessage.IsNull() {
				apiAgent.UpdateMessage = agent.UpdateMessage.ValueStringPointer()
			}
			if !agent.OperationalStatus.IsNull() {
				apiAgent.OperationalStatus = agent.OperationalStatus.ValueStringPointer()
			}
			if !agent.IsLatestGAVersion.IsNull() {
				apiAgent.IsLatestGAedVersion = agent.IsLatestGAVersion.ValueBoolPointer()
			}
			if !agent.IsHidden.IsNull() {
				apiAgent.IsHidden = agent.IsHidden.ValueBoolPointer()
			}
			if !agent.LastConnection.IsNull() && !agent.LastConnection.IsUnknown() {
				t, err := time.Parse(time.RFC3339, agent.LastConnection.ValueString())
				if err == nil {
					apiAgent.LastConnection = &t
				}
			}

			agents[i] = apiAgent
		}
		req.Agents = agents
	}

	return req
}

func (r *agentPoolUpdateResource) buildUpdateRequest(plan agentPoolUpdateResourceModel) v5okta.AgentPoolUpdate {
	// Same as create request for this resource
	return r.buildCreateRequest(plan)
}

func (r *agentPoolUpdateResource) mapResponseToState(resp *v5okta.AgentPoolUpdate, state *agentPoolUpdateResourceModel) {
	if resp.Id != nil {
		state.Id = types.StringValue(*resp.Id)
	}

	if resp.Name != nil {
		state.Name = types.StringValue(*resp.Name)
	} else {
		state.Name = types.StringNull()
	}

	if resp.AgentType != nil {
		state.AgentType = types.StringValue(*resp.AgentType)
	} else {
		state.AgentType = types.StringNull()
	}

	if resp.Enabled != nil {
		state.Enabled = types.BoolValue(*resp.Enabled)
	} else {
		state.Enabled = types.BoolNull()
	}

	if resp.NotifyAdmin != nil {
		state.NotifyAdmins = types.BoolValue(*resp.NotifyAdmin)
	} else {
		state.NotifyAdmins = types.BoolNull()
	}

	if resp.Reason != nil {
		state.Reason = types.StringValue(*resp.Reason)
	} else {
		state.Reason = types.StringNull()
	}

	if resp.SortOrder != nil {
		state.SortOrder = types.Int64Value(int64(*resp.SortOrder))
	} else {
		state.SortOrder = types.Int64Null()
	}

	if resp.TargetVersion != nil {
		state.TargetVersion = types.StringValue(*resp.TargetVersion)
	} else {
		state.TargetVersion = types.StringNull()
	}

	if resp.Status != nil {
		state.Status = types.StringValue(*resp.Status)
	} else {
		state.Status = types.StringNull()
	}

	// Handle Schedule
	if resp.Schedule != nil {
		schedule := &Schedule{}

		if resp.Schedule.Cron != nil {
			schedule.Cron = types.StringValue(*resp.Schedule.Cron)
		} else {
			schedule.Cron = types.StringNull()
		}

		if resp.Schedule.Delay != nil {
			schedule.Delay = types.Int64Value(int64(*resp.Schedule.Delay))
		} else {
			schedule.Delay = types.Int64Null()
		}

		if resp.Schedule.Duration != nil {
			schedule.Duration = types.Int64Value(int64(*resp.Schedule.Duration))
		} else {
			schedule.Duration = types.Int64Null()
		}

		if resp.Schedule.LastUpdated != nil {
			schedule.LastUpdated = types.StringValue(resp.Schedule.LastUpdated.Format(time.RFC3339))
		} else {
			schedule.LastUpdated = types.StringNull()
		}

		if resp.Schedule.Timezone != nil {
			schedule.Timezone = types.StringValue(*resp.Schedule.Timezone)
		} else {
			schedule.Timezone = types.StringNull()
		}

		state.Schedule = schedule
	} else {
		state.Schedule = nil
	}

	// Handle Agents
	if len(resp.Agents) > 0 {
		agents := make([]Agent, len(resp.Agents))
		for i, agent := range resp.Agents {
			tfAgent := Agent{}

			if agent.Id != nil {
				tfAgent.Id = types.StringValue(*agent.Id)
			} else {
				tfAgent.Id = types.StringNull()
			}

			if agent.Name != nil {
				tfAgent.Name = types.StringValue(*agent.Name)
			} else {
				tfAgent.Name = types.StringNull()
			}

			if agent.Status != nil {
				tfAgent.Status = types.StringValue(*agent.Status)
			} else {
				tfAgent.Status = types.StringNull()
			}

			if agent.Type != nil {
				tfAgent.Type = types.StringValue(*agent.Type)
			} else {
				tfAgent.Type = types.StringNull()
			}

			agents[i] = tfAgent
		}
		state.Agents = agents
	} else {
		state.Agents = []Agent{}
	}
}
