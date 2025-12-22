package idaas

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ datasource.DataSource = &agentPoolDataSource{}

func newAgentPoolDataSource() datasource.DataSource {
	return &agentPoolDataSource{}
}

type agentPoolDataSource struct {
	*config.Config
}

func (d *agentPoolDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_agent_pool_update"
}

func (d *agentPoolDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

type AgentPoolUpdateDataSourceModel struct {
	ID            types.String           `tfsdk:"id"`
	PoolId        types.String           `tfsdk:"pool_id"`
	Name          types.String           `tfsdk:"name"`
	AgentType     types.String           `tfsdk:"agent_type"`
	Enabled       types.Bool             `tfsdk:"enabled"`
	NotifyAdmin   types.Bool             `tfsdk:"notify_admin"`
	Reason        types.String           `tfsdk:"reason"`
	SortOrder     types.Int64            `tfsdk:"sort_order"`
	Status        types.String           `tfsdk:"status"`
	TargetVersion types.String           `tfsdk:"target_version"`
	Schedule      *UpdateSchedule        `tfsdk:"schedule"`
	Agents        []AgentDataSourceModel `tfsdk:"agents"`
}

type AgentDataSourceModel struct {
	ID                  types.String `tfsdk:"id"`
	IsHidden            types.Bool   `tfsdk:"is_hidden"`
	IsLatestGAedVersion types.Bool   `tfsdk:"is_latest_gaed_version"`
	LastConnection      types.Int64  `tfsdk:"last_connection"`
	Name                types.String `tfsdk:"name"`
	OperationalStatus   types.String `tfsdk:"operational_status"`
	PoolId              types.String `tfsdk:"pool_id"`
	Type                types.String `tfsdk:"type"`
	UpdateMessage       types.String `tfsdk:"update_message"`
	UpdateStatus        types.String `tfsdk:"update_status"`
	Version             types.String `tfsdk:"version"`
}

func (d *agentPoolDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves an Okta Agent Pool Update. Agent pool updates allow you to schedule and manage updates for agent pools.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The unique identifier of the agent pool update.",
			},
			"pool_id": schema.StringAttribute{
				Required:    true,
				Description: "The unique identifier of the agent pool.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the agent pool update.",
			},
			"agent_type": schema.StringAttribute{
				Computed:    true,
				Description: "Agent types that are being monitored.",
			},
			"enabled": schema.BoolAttribute{
				Computed:    true,
				Description: "Indicates if auto-update is enabled for the agent pool.",
			},
			"notify_admin": schema.BoolAttribute{
				Computed:    true,
				Description: "Indicates if the admin is notified about the update.",
			},
			"reason": schema.StringAttribute{
				Computed:    true,
				Description: "Reason for the update.",
			},
			"sort_order": schema.Int64Attribute{
				Computed:    true,
				Description: "Specifies the sort order.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Overall state for the auto-update job from the admin perspective.",
			},
			"target_version": schema.StringAttribute{
				Computed:    true,
				Description: "The agent version to update to.",
			},
		},
		Blocks: map[string]schema.Block{
			"schedule": schema.SingleNestedBlock{
				Description: "The schedule configuration for the agent pool update.",
				Attributes: map[string]schema.Attribute{
					"cron": schema.StringAttribute{
						Computed: true,
					},
					"delay": schema.Int32Attribute{
						Computed: true,
					},
					"duration": schema.Int32Attribute{
						Computed: true,
					},
					"last_updated": schema.StringAttribute{
						Computed: true,
					},
					"timezone": schema.StringAttribute{
						Computed: true,
					},
				},
			},
			"agents": schema.SetNestedBlock{
				Description: "The agents associated with the agent pool update.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"is_hidden": schema.BoolAttribute{
							Computed: true,
						},
						"is_latest_gaed_version": schema.BoolAttribute{
							Computed: true,
						},
						"last_connection": schema.Int64Attribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"operational_status": schema.StringAttribute{
							Computed: true,
						},
						"pool_id": schema.StringAttribute{
							Computed: true,
						},
						"type": schema.StringAttribute{
							Computed: true,
						},
						"update_message": schema.StringAttribute{
							Computed: true,
						},
						"update_status": schema.StringAttribute{
							Computed: true,
						},
						"version": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *agentPoolDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AgentPoolUpdateDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	getAgentPoolUpdateResp, _, err := d.OktaIDaaSClient.OktaSDKClientV6().AgentPoolsAPI.GetAgentPoolsUpdateInstance(ctx, data.PoolId.ValueString(), data.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading agent pool update",
			"Could not read agent pool update, unexpected error: "+err.Error(),
		)
		return
	}

	data.Name = types.StringValue(getAgentPoolUpdateResp.GetName())
	data.Enabled = types.BoolValue(getAgentPoolUpdateResp.GetEnabled())
	data.NotifyAdmin = types.BoolValue(getAgentPoolUpdateResp.GetNotifyAdmin())
	data.Reason = types.StringValue(getAgentPoolUpdateResp.GetReason())
	data.AgentType = types.StringValue(getAgentPoolUpdateResp.GetAgentType())
	data.SortOrder = types.Int64Value(int64(getAgentPoolUpdateResp.GetSortOrder()))
	data.TargetVersion = types.StringValue(getAgentPoolUpdateResp.GetTargetVersion())
	data.Status = types.StringValue(getAgentPoolUpdateResp.GetStatus())
	data.Schedule = &UpdateSchedule{}
	if getAgentPoolUpdateResp.Schedule != nil {
		data.Schedule.Delay = types.Int32Value(getAgentPoolUpdateResp.Schedule.GetDelay())
		data.Schedule.Duration = types.Int32Value(getAgentPoolUpdateResp.Schedule.GetDuration())
		data.Schedule.Cron = types.StringValue(getAgentPoolUpdateResp.Schedule.GetCron())
		data.Schedule.Timezone = types.StringValue(getAgentPoolUpdateResp.Schedule.GetTimezone())
		data.Schedule.LastUpdated = types.StringValue(getAgentPoolUpdateResp.Schedule.GetLastUpdated().Format(time.RFC3339))
	}

	var agents []AgentDataSourceModel
	for _, agentItem := range getAgentPoolUpdateResp.GetAgents() {
		agent := AgentDataSourceModel{
			ID: types.StringValue(agentItem.GetId()),
		}
		agent.PoolId = types.StringValue(agentItem.GetPoolId())
		agent.Name = types.StringValue(agentItem.GetName())
		agent.Type = types.StringValue(agentItem.GetType())
		agent.Version = types.StringValue(agentItem.GetVersion())
		agent.OperationalStatus = types.StringValue(agentItem.GetOperationalStatus())
		agent.UpdateStatus = types.StringValue(agentItem.GetUpdateStatus())
		agent.UpdateMessage = types.StringValue(agentItem.GetUpdateMessage())
		agent.LastConnection = types.Int64Value(agentItem.GetLastConnection())
		agent.IsHidden = types.BoolValue(agentItem.GetIsHidden())
		agent.IsLatestGAedVersion = types.BoolValue(agentItem.GetIsLatestGAedVersion())
		agents = append(agents, agent)
	}
	data.Agents = agents

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
