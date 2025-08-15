package governance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ datasource.DataSource = (*requestSettingOrganizationDataSource)(nil)

func newRequestSettingOrganizationDataSource() datasource.DataSource {
	return &requestSettingOrganizationDataSource{}
}

type requestSettingOrganizationDataSource struct {
	*config.Config
}

func (d *requestSettingOrganizationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

type requestSettingOrganizationDataSourceModel struct {
	LongTimePastProvisioned   types.Bool   `tfsdk:"long_time_past_provisioned"`
	ProvisioningStatus        types.String `tfsdk:"provisioning_status"`
	RequestExperiences        []experience `tfsdk:"request_experiences"`
	SubprocessorsAcknowledged types.Bool   `tfsdk:"subprocessors_acknowledged"`
}

func (d *requestSettingOrganizationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_request_setting_organization"
}

func (d *requestSettingOrganizationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"long_time_past_provisioned": schema.BoolAttribute{
				Computed: true,
			},
			"provisioning_status": schema.StringAttribute{
				Computed: true,
			},
			"subprocessors_acknowledged": schema.BoolAttribute{
				Computed: true,
			},
		},
		Blocks: map[string]schema.Block{
			"request_experiences": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"experience_type": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *requestSettingOrganizationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data requestSettingOrganizationDataSourceModel

	// Read Terraform configuration Data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	orgSettingsResp, _, err := d.OktaGovernanceClient.OktaIGSDKClient().RequestSettingsAPI.GetOrgRequestSettingsV2(ctx).Execute()
	if err != nil {
		return
	}

	// Example Data value setting
	data.SubprocessorsAcknowledged = types.BoolValue(orgSettingsResp.SubprocessorsAcknowledged)
	data.ProvisioningStatus = types.StringValue(string(orgSettingsResp.ProvisioningStatus))
	data.LongTimePastProvisioned = types.BoolValue(orgSettingsResp.LongTimePastProvisioned)
	var experiences []experience
	for _, exp := range orgSettingsResp.GetRequestExperiences() {
		experiences = append(experiences, experience{
			ExperienceType: string(exp),
		})
	}

	data.RequestExperiences = experiences

	// Save Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
