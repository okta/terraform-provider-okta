package governance

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ datasource.DataSource = &requestConditionDataSource{}

func newRequestConditionDataSource() datasource.DataSource {
	return &requestConditionDataSource{}
}

func (d *requestConditionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_request_condition"
}

func (d *requestConditionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

func (d *requestConditionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required: true,
			},
			"resource_id": schema.StringAttribute{
				Required: true,
			},
			"created": schema.StringAttribute{
				Computed: true,
			},
			"created_by": schema.StringAttribute{
				Computed: true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"last_updated_by": schema.StringAttribute{
				Computed: true,
			},
			"status": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Computed: true,
			},
			"priority": schema.Int32Attribute{
				Computed: true,
			},
		},
		Blocks: map[string]schema.Block{
			"access_scope_settings": schema.SingleNestedBlock{
				Blocks: map[string]schema.Block{
					"ids": schema.ListNestedBlock{
						Description: "List of groups/entitlement bundles.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Computed:    true,
									Description: "The group/entitlement bundle ID.",
								},
							},
						},
					},
				},
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Computed: true,
					},
				},
			},
			"requester_settings": schema.SingleNestedBlock{
				Blocks: map[string]schema.Block{
					"ids": schema.ListNestedBlock{
						Description: "List of teams/groups ids.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Computed:    true,
									Description: "The group/team ID.",
								},
							},
						},
					},
				},
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Computed: true,
					},
				},
			},
		},
	}
}

type requestConditionDataSource struct {
	*config.Config
}
type requestConditionsDataSourceModel struct {
	Id                  types.String `tfsdk:"id"`
	ResourceId          types.String `tfsdk:"resource_id"`
	Created             types.String `tfsdk:"created"`
	CreatedBy           types.String `tfsdk:"created_by"`
	LastUpdated         types.String `tfsdk:"last_updated"`
	LastUpdatedBy       types.String `tfsdk:"last_updated_by"`
	Status              types.String `tfsdk:"status"`
	Name                types.String `tfsdk:"name"`
	Priority            types.Int32  `tfsdk:"priority"`
	AccessScopeSettings *Settings    `tfsdk:"access_scope_settings"`
	RequesterSettings   *Settings    `tfsdk:"requester_settings"`
}

func (d *requestConditionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data requestConditionsDataSourceModel

	// Read Terraform configuration Data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	readRequestConditionResp, _, err := d.OktaGovernanceClient.OktaIGSDKClient().RequestConditionsAPI.GetResourceRequestConditionV2(ctx, data.ResourceId.ValueString(), data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Request conditions",
			"Could not read Request conditions, unexpected error: "+err.Error(),
		)
		return
	}
	// Example Data value setting
	data.Id = types.StringValue(readRequestConditionResp.GetId())
	data.Name = types.StringValue(readRequestConditionResp.GetName())
	data.Priority = types.Int32Value(readRequestConditionResp.GetPriority())
	data.Created = types.StringValue(readRequestConditionResp.GetCreated().Format(time.RFC3339))
	data.CreatedBy = types.StringValue(readRequestConditionResp.GetCreatedBy())
	data.LastUpdated = types.StringValue(readRequestConditionResp.GetLastUpdated().Format(time.RFC3339))
	data.LastUpdatedBy = types.StringValue(readRequestConditionResp.GetLastUpdatedBy())
	data.Status = types.StringValue(string(readRequestConditionResp.GetStatus()))
	data.AccessScopeSettings, _ = setAccessScopeSettings(readRequestConditionResp.GetAccessScopeSettings())
	data.RequesterSettings, _ = setRequesterSettings(readRequestConditionResp.GetRequesterSettings())

	// Save Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
