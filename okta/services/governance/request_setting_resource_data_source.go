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

type requestSettingResourceDataSourceModel struct {
	ResourceId                  types.String                `tfsdk:"resource_id"`
	ValidAccessDurationSettings validAccessDurationSettings `tfsdk:"valid_access_duration_settings"`
	ValidAccessScopeSettings    []supportedTypes            `tfsdk:"valid_access_scope_settings"`
	ValidRequesterSettings      []supportedTypes            `tfsdk:"valid_requester_settings"`
}

func (d *requestSettingResourceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_request_setting_resource"
}

func (d *requestSettingResourceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"resource_id": schema.StringAttribute{
				Required: true,
			},
		},
		Blocks: map[string]schema.Block{
			"valid_access_duration_settings": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"maximum_days": schema.Float32Attribute{
						Computed: true,
					},
					"maximum_weeks": schema.Float32Attribute{
						Computed: true,
					},
					"maximum_hours": schema.Float32Attribute{
						Computed: true,
					},
					"required": schema.BoolAttribute{
						Computed: true,
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
					},
				},
			},
			"valid_access_scope_settings": schema.ListNestedBlock{
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
			"valid_requester_settings": schema.ListNestedBlock{
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
		},
	}
}

func (d *requestSettingResourceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data requestSettingResourceDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	reqSettingsResp, _, err := d.OktaGovernanceClient.OktaIGSDKClient().RequestSettingsAPI.GetRequestSettingsV2(ctx, data.ResourceId.ValueString()).Execute()
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

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
