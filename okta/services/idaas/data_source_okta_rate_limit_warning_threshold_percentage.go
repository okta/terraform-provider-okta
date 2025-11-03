package idaas

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	tfpath "github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ datasource.DataSource = &rateLimitWarningThresholdPercentageDataSource{}

func newRateLimitWarningThresholdPercentageDataSource() datasource.DataSource {
	return &rateLimitWarningThresholdPercentageDataSource{}
}

func (d *rateLimitWarningThresholdPercentageDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

type rateLimitWarningThresholdPercentageDataSource struct {
	*config.Config
}

func (d *rateLimitWarningThresholdPercentageDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rate_limit_warning_threshold_percentage"
}

func (d *rateLimitWarningThresholdPercentageDataSource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, tfpath.Root("id"), req, resp)
}

func (d *rateLimitWarningThresholdPercentageDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"warning_threshold": schema.Int32Attribute{
				Computed: true,
			},
		},
	}
}

type rateLimitAdminWarningThresholdDataSourceDataSourceModel struct {
	Id               types.String `tfsdk:"id"`
	WarningThreshold types.Int32  `tfsdk:"warning_threshold"`
}

func (d *rateLimitWarningThresholdPercentageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data rateLimitAdminWarningThresholdDataSourceDataSourceModel

	// Read Terraform prior state Data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	getRateLimitWarningThresholdResp, _, err := d.OktaIDaaSClient.OktaSDKClientV5().RateLimitSettingsAPI.GetRateLimitSettingsWarningThreshold(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to read rate limit warning threshold percentage",
			err.Error(),
		)
		return
	}

	data.Id = types.StringValue("rate_limiting_warning_threshold_percentage")
	data.WarningThreshold = types.Int32Value(getRateLimitWarningThresholdResp.GetWarningThreshold())

	// Save updated Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
