package idaas

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	tfpath "github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ datasource.DataSource = &principalRateLimitsDataSource{}

func newPrincipalRateLimitsDataSource() datasource.DataSource {
	return &principalRateLimitsDataSource{}
}

func (d *principalRateLimitsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

type principalRateLimitsDataSource struct {
	*config.Config
}

func (d *principalRateLimitsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_principal_rate_limits"
}

func (d *principalRateLimitsDataSource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, tfpath.Root("id"), req, resp)
}

func (d *principalRateLimitsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required: true,
			},
			"principal_id": schema.StringAttribute{
				Computed: true,
			},
			"principal_type": schema.StringAttribute{
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						"OAUTH_CLIENT",
						"SSWS_TOKEN",
					}...),
				},
			},
			"default_concurrency_percentage": schema.Int32Attribute{
				Computed: true,
			},
			"default_percentage": schema.Int32Attribute{
				Computed: true,
			},
			"created_by": schema.StringAttribute{
				Computed: true,
			},
			"created_date": schema.StringAttribute{
				Computed: true,
			},
			"last_update": schema.StringAttribute{
				Computed: true,
			},
			"last_updated_by": schema.StringAttribute{
				Computed: true,
			},
			"org_id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

type principalRateLimitsDataSourceModel struct {
	Id                           types.String `tfsdk:"id"`
	principalId                  types.String `tfsdk:"principal_id"`
	principalType                types.String `tfsdk:"principal_type"`
	defaultConcurrencyPercentage types.Int32  `tfsdk:"default_concurrency_percentage"`
	defaultPercentage            types.Int32  `tfsdk:"default_percentage"`
	createdBy                    types.String `tfsdk:"created_by"`
	createdDate                  types.String `tfsdk:"created_date"`
	lastUpdate                   types.String `tfsdk:"last_update"`
	lastUpdatedBy                types.String `tfsdk:"last_updated_by"`
	orgId                        types.String `tfsdk:"org_id"`
}

func (d *principalRateLimitsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data principalRateLimitsDataSourceModel

	// Read Terraform prior state Data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	getPrincipalRateSettingsResp, _, err := d.OktaIDaaSClient.OktaSDKClientV5().PrincipalRateLimitAPI.GetPrincipalRateLimitEntity(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		return
	}

	data.Id = types.StringValue(getPrincipalRateSettingsResp.GetId())
	data.principalId = types.StringValue(getPrincipalRateSettingsResp.GetPrincipalId())
	data.principalType = types.StringValue(getPrincipalRateSettingsResp.GetPrincipalType())
	data.defaultConcurrencyPercentage = types.Int32Value(getPrincipalRateSettingsResp.GetDefaultConcurrencyPercentage())
	data.defaultPercentage = types.Int32Value(getPrincipalRateSettingsResp.GetDefaultPercentage())
	data.createdBy = types.StringValue(getPrincipalRateSettingsResp.GetCreatedBy())
	data.createdDate = types.StringValue(getPrincipalRateSettingsResp.GetCreatedDate().Format(time.RFC3339))
	data.lastUpdate = types.StringValue(getPrincipalRateSettingsResp.GetLastUpdate().Format(time.RFC3339))
	data.lastUpdatedBy = types.StringValue(getPrincipalRateSettingsResp.GetLastUpdatedBy())
	data.orgId = types.StringValue(getPrincipalRateSettingsResp.GetOrgId())

	// Save updated Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
