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
	resp.TypeName = req.ProviderTypeName + "_agent_pool"
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
	NotifyAdmins  types.Bool             `tfsdk:"notify_admins"`
	Reason        types.String           `tfsdk:"reason"`
	Schedule      *UpdateSchedule        `tfsdk:"schedule"`
	SortOrder     types.Int64            `tfsdk:"sort_order"`
	TargetVersion types.String           `tfsdk:"target_version"`
	Status        types.String           `tfsdk:"status"`
	Agents        []AgentDataSourceModel `tfsdk:"agents"`
}

type AgentDataSourceModel struct {
	ID    types.String `tfsdk:"id"`
	Agent Agent        `tfsdk:"agent"`
}

func (d *agentPoolDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Okta Agent Pool Update. Agent pool updates allow you to schedule and manage updates for agent pools.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The unique identifier of the agent pool update.",
			},
			"pool_id": schema.StringAttribute{
				Required:    true,
				Description: "The unique identifier of the agent pool to update.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The name of the agent pool update.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
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
				Computed:    true,
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
				Computed:    true,
				Description: "Whether to send notifications when the update completes.",
			},
		},
		Blocks: map[string]schema.Block{
			"schedule": schema.SingleNestedBlock{
				Description: "The schedule configuration for the agent pool update.",
				Attributes: map[string]schema.Attribute{
					"cron": schema.StringAttribute{
						Computed: true,
					},
					"delay": schema.Int64Attribute{
						Computed: true,
					},
					"duration": schema.Int64Attribute{
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
						"last_connection": schema.StringAttribute{
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
	data.Reason = types.StringValue(getAgentPoolUpdateResp.GetReason())
	data.AgentType = types.StringValue(getAgentPoolUpdateResp.GetAgentType())
	data.SortOrder = types.Int64Value(int64(getAgentPoolUpdateResp.GetSortOrder()))
	data.TargetVersion = types.StringValue(getAgentPoolUpdateResp.GetTargetVersion())
	data.Status = types.StringValue(getAgentPoolUpdateResp.GetStatus())
	var agents []AgentDataSourceModel
	for _, agentItem := range getAgentPoolUpdateResp.GetAgents() {
		agent := AgentDataSourceModel{
			ID: types.StringValue(agentItem.GetId()),
		}
		agent.Agent = Agent{}
		agent.Agent.IsHidden = types.BoolValue(agentItem.GetIsHidden())
		agent.Agent.IsLatestGAedVersion = types.BoolValue(agentItem.GetIsLatestGAedVersion())
		agent.Agent.LastConnection = types.StringValue(agentItem.GetLastConnection().Format(time.RFC3339))
		agent.Agent.Name = types.StringValue(agentItem.GetName())
		agent.Agent.OperationalStatus = types.StringValue(agentItem.GetOperationalStatus())
		agent.Agent.PoolId = types.StringValue(agentItem.GetPoolId())
		agent.Agent.Type = types.StringValue(agentItem.GetType())
		agent.Agent.UpdateMessage = types.StringValue(agentItem.GetUpdateMessage())
		agent.Agent.UpdateStatus = types.StringValue(agentItem.GetUpdateStatus())
		agent.Agent.Version = types.StringValue(agentItem.GetVersion())

		agents = append(agents, agent)
	}
}
