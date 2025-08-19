package governance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ datasource.DataSource = (*requestSettingResourceDataSource)(nil)

func newRequestSettingResourceDataSource() datasource.DataSource {
	return &requestSettingResourceDataSource{}
}

type requestSettingResourceDataSource struct {
	*config.Config
}

func (d *requestSettingResourceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

type requestSettingResourceDataSourceModel struct {
	Id                          types.String                 `tfsdk:"id"`
	ValidAccessDurationSettings *validAccessDurationSettings `tfsdk:"valid_access_duration_settings"`
	ValidAccessScopeSettings    []supportedTypes             `tfsdk:"valid_access_scope_settings"`
	ValidRequesterSettings      []supportedTypes             `tfsdk:"valid_requester_settings"`
	RequestOnBehalfOfSettings   *requestOnBehalfOfSettings   `tfsdk:"request_on_behalf_of_settings"`
	RiskSettings                *riskSettings                `tfsdk:"risk_settings"`
}

func (d *requestSettingResourceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_request_setting_resource"
}

func (d *requestSettingResourceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The id of the resource in Okta ID format.",
			},
		},
		Blocks: map[string]schema.Block{
			"valid_access_duration_settings": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"maximum_days": schema.Float32Attribute{
						Computed:    true,
						Description: "The maximum value allowed for a request condition or risk setting.",
					},
					"maximum_weeks": schema.Float32Attribute{
						Computed:    true,
						Description: "The maximum value allowed for a request condition or risk setting.",
					},
					"maximum_hours": schema.Float32Attribute{
						Computed:    true,
						Description: "The maximum value allowed for a request condition or risk setting.",
					},
					"required": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether accessDurationSetting must be included in the request conditions or risk settings for the specified resource.",
					},
				},
				Blocks: map[string]schema.Block{
					"supported_types": schema.SetNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Computed: true,
									Validators: []validator.String{
										stringvalidator.OneOf("ADMIN_FIXED_DURATION", "REQUESTER_SPECIFIED_DURATION"),
									},
								},
							},
						},
						Description: "Access duration settings that are eligible to be added to a request condition or risk settings for the specified resource.",
					},
				},
			},
			"valid_access_scope_settings": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Computed: true,
							Validators: []validator.String{
								stringvalidator.OneOf(
									"ENTITLEMENT_BUNDLES",
									"GROUPS",
									"RESOURCE_DEFAULT",
								),
							},
						},
					},
				},
				Description: "Access scope settings eligible to be added to a request condition.",
			},
			"valid_requester_settings": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Computed: true,
							Validators: []validator.String{
								stringvalidator.OneOf(
									"EVERYONE",
									"GROUPS",
									"TEAMS",
								),
							},
						},
					},
				},
				Description: "Access scope settings eligible to be added to a request condition.",
			},
			"request_on_behalf_of_settings": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"allowed": schema.BoolAttribute{
						Computed: true,
					},
				},
				Blocks: map[string]schema.Block{
					"only_for": schema.SetNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Computed: true,
								},
							},
						},
					},
				},
				Description: "Specifies if and for whom a requester may request the resource for.",
			},
			"risk_settings": schema.SingleNestedBlock{
				Blocks: map[string]schema.Block{
					"default_setting": schema.SingleNestedBlock{
						Attributes: map[string]schema.Attribute{
							"request_submission_type": schema.StringAttribute{
								Computed: true,
							},
							"error": schema.ListAttribute{
								ElementType: types.StringType,
								Computed:    true,
							},
							"approval_sequence_id": schema.StringAttribute{
								Computed: true,
							},
						},
						Blocks: map[string]schema.Block{
							"access_duration_settings": schema.SingleNestedBlock{
								Attributes: map[string]schema.Attribute{
									"type": schema.StringAttribute{
										Computed: true,
									},
									"duration": schema.StringAttribute{
										Computed: true,
									},
								},
							},
						},
						Description: "Default risk settings that are valid for an access request when a risk has been detected for the resource and requesting user.",
					},
				},
				Description: "Risk settings that are valid for an access request when a risk has been detected for the resource and requesting user",
			},
		},
	}
}

func (d *requestSettingResourceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data requestSettingResourceDataSourceModel

	// Read Terraform configuration Data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	reqSettingsResp, _, err := d.OktaGovernanceClient.OktaGovernanceSDKClient().RequestSettingsAPI.GetRequestSettingsV2(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Request Settings",
			"Could not read request settings, unexpected error: "+err.Error(),
		)
		return
	}

	data.ValidAccessScopeSettings = setValidAccessSettings(reqSettingsResp.ValidAccessScopeSettings)
	data.ValidAccessDurationSettings = setValidAccessDurationSettings(reqSettingsResp.ValidAccessDurationSettings)
	data.ValidRequesterSettings = setValidRequesterSettings(reqSettingsResp.ValidRequesterSettings)
	data.RequestOnBehalfOfSettings = setRequesterOnBehalfSettings(reqSettingsResp.RequestOnBehalfOfSettings)
	data.RiskSettings = setRiskSettings(reqSettingsResp.RiskSettings)

	// Save Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
