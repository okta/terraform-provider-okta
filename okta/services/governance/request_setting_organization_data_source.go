package governance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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

type requestSettingOrganizationDataSourceModel struct {
	LongTimePastProvisioned   types.Bool   `tfsdk:"long_time_past_provisioned"`
	ProvisioningStatus        types.String `tfsdk:"provisioning_status"`
	RequestExperiences        types.List   `tfsdk:"request_experiences"`
	SubprocessorsAcknowledged types.Bool   `tfsdk:"subprocessors_acknowledged"`
}

func (d *requestSettingOrganizationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_request_setting_organization"
}

func (d *requestSettingOrganizationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"long_time_past_provisioned": schema.StringAttribute{
				Computed: true,
			},
			"provisioning_status": schema.StringAttribute{
				Computed: true,
			},
			"request_experiences": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"subprocessors_acknowledged": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *requestSettingOrganizationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data requestSettingOrganizationDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	orgSettingsResp, _, err := d.OktaGovernanceClient.OktaIGSDKClient().RequestSettingsAPI.GetOrgRequestSettingsV2(ctx).Execute()
	if err != nil {
		return
	}

	// Example data value setting
	data.SubprocessorsAcknowledged = types.BoolValue(orgSettingsResp.SubprocessorsAcknowledged)
	data.ProvisioningStatus = types.StringValue(string(orgSettingsResp.ProvisioningStatus))
	data.LongTimePastProvisioned = types.BoolValue(orgSettingsResp.LongTimePastProvisioned)
	var reqExpVals []attr.Value
	for _, reqExp := range orgSettingsResp.GetRequestExperiences() {
		reqExpVals = append(reqExpVals, types.StringValue(string(reqExp)))
	}

	listVal, _ := types.ListValue(types.StringType, reqExpVals)
	data.RequestExperiences = listVal

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
