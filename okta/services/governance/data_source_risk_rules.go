package governance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ datasource.DataSource = &riskRulesDataSource{}

func newRiskRulesDataSource() datasource.DataSource {
	return &riskRulesDataSource{}
}

type riskRulesDataSource struct {
	*config.Config
}

type riskRulesDataSourceModel struct {
	Id               types.String            `tfsdk:"id"`
	Name             types.String            `tfsdk:"name"`
	Description      types.String            `tfsdk:"description"`
	Type             types.String            `tfsdk:"type"`
	Status           types.String            `tfsdk:"status"`
	Resources        []riskRuleResourceModel `tfsdk:"resources"`
	LastUpdated      types.String            `tfsdk:"last_updated"`
	Created          types.String            `tfsdk:"created"`
	CreatedBy        types.String            `tfsdk:"created_by"`
	LastUpdatedBy    types.String            `tfsdk:"last_updated_by"`
	ConflictCriteria *conflictCriteriaModel  `tfsdk:"conflict_criteria"`
	Notes            types.String            `tfsdk:"notes"`
}

func (d *riskRulesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_risk_rules"
}

func (d *riskRulesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

func (d *riskRulesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Schema for a Separation of Duties Policy",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required: true,
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The name of the policy.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description of the policy.",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "Type of the policy, e.g., SEPARATION_OF_DUTIES.",
			},
			"created": schema.StringAttribute{
				Computed: true,
			},
			"created_by": schema.StringAttribute{
				Computed: true,
			},
			"last_updated_by": schema.StringAttribute{
				Computed: true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"notes": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"status": schema.StringAttribute{
				Computed: true,
			},
		},
		Blocks: map[string]schema.Block{
			"resources": schema.ListNestedBlock{
				Description: "List of resources this policy applies to.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"resource_orn": schema.StringAttribute{
							Computed:    true,
							Description: "Object Resource Name of the resource.",
						},
					},
				},
			},
			"conflict_criteria": schema.SingleNestedBlock{
				Description: "Defines the logical conditions for policy conflict.",
				Blocks: map[string]schema.Block{
					"and": schema.ListNestedBlock{
						Description: "AND group of conditions.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Computed:    true,
									Description: "Name of the rule in the AND group.",
								},
								"attribute": schema.StringAttribute{
									Computed:    true,
									Description: "Attribute path to check for conflict.",
								},
								"operation": schema.StringAttribute{
									Computed:    true,
									Description: "Operation type, e.g., CONTAINS_ONE, CONTAINS_ALL.",
								},
							},
							Blocks: map[string]schema.Block{
								"value": schema.SingleNestedBlock{
									Description: "Values to match for the condition.",
									Attributes: map[string]schema.Attribute{
										"type": schema.StringAttribute{
											Computed:    true,
											Description: "Type of the value, e.g., ENTITLEMENTS.",
										},
									},
									Blocks: map[string]schema.Block{
										"value": schema.ListNestedBlock{
											Description: "List of entitlements and their values.",
											NestedObject: schema.NestedBlockObject{
												Attributes: map[string]schema.Attribute{
													"id": schema.StringAttribute{
														Computed: true,
													},
													"name": schema.StringAttribute{
														Optional: true,
														Computed: true,
													},
													"external_value": schema.StringAttribute{
														Optional: true,
														Computed: true,
													},
												},
												Blocks: map[string]schema.Block{
													"values": schema.ListNestedBlock{
														NestedObject: schema.NestedBlockObject{
															Attributes: map[string]schema.Attribute{
																"id": schema.StringAttribute{
																	Computed: true,
																},
																"name": schema.StringAttribute{
																	Optional: true,
																	Computed: true,
																},
																"external_value": schema.StringAttribute{
																	Optional: true,
																	Computed: true,
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *riskRulesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data riskRulesDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	getRiskRuleResp, _, err := d.OktaGovernanceClient.OktaIGSDKClient().RiskRulesAPI.GetRiskRule(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		return
	}

	// Example data value setting
	data.Id = types.StringValue(getRiskRuleResp.Id)
	data.Name = types.StringValue(getRiskRuleResp.Name)
	data.Description = types.StringPointerValue(getRiskRuleResp.Description)
	data.Status = types.StringValue(getRiskRuleResp.Status)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
