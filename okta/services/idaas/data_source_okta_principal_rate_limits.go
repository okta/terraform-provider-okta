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
				Computed:    true,
				Description: "The unique identifier of the principal. This is the ID of the API token or OAuth 2.0 app.",
			},
			"principal_type": schema.StringAttribute{
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						"OAUTH_CLIENT",
						"SSWS_TOKEN",
					}...),
				},
				Description: "The type of principal, either an API token or an OAuth 2.0 app.",
			},
			"default_concurrency_percentage": schema.Int32Attribute{
				Computed:    true,
				Description: "The default percentage of a given concurrency limit threshold that the owning principal can consume.",
			},
			"default_percentage": schema.Int32Attribute{
				Computed:    true,
				Description: "The default percentage of a given rate limit threshold that the owning principal can consume.",
			},
			"created_by": schema.StringAttribute{
				Computed:    true,
				Description: "The Okta user ID of the user who created the principle rate limit entity.",
			},
			"created_date": schema.StringAttribute{
				Computed:    true,
				Description: "The date and time the principle rate limit entity was created.",
			},
			"last_update": schema.StringAttribute{
				Computed:    true,
				Description: "The date and time the principle rate limit entity was last updated.",
			},
			"last_updated_by": schema.StringAttribute{
				Computed:    true,
				Description: "The Okta user ID of the user who last updated the principle rate limit entity.",
			},
			"org_id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier of the Okta org.",
			},
		},
	}
}

type principalRateLimitsDataSourceModel struct {
	Id                           types.String `tfsdk:"id"`
	PrincipalId                  types.String `tfsdk:"principal_id"`
	PrincipalType                types.String `tfsdk:"principal_type"`
	DefaultConcurrencyPercentage types.Int32  `tfsdk:"default_concurrency_percentage"`
	DefaultPercentage            types.Int32  `tfsdk:"default_percentage"`
	CreatedBy                    types.String `tfsdk:"created_by"`
	CreatedDate                  types.String `tfsdk:"created_date"`
	LastUpdate                   types.String `tfsdk:"last_update"`
	LastUpdatedBy                types.String `tfsdk:"last_updated_by"`
	OrgId                        types.String `tfsdk:"org_id"`
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
		resp.Diagnostics.AddError(
			"failed to read principal rate limit",
			err.Error(),
		)
		return
	}

	data.Id = types.StringValue(getPrincipalRateSettingsResp.GetId())
	data.PrincipalId = types.StringValue(getPrincipalRateSettingsResp.GetPrincipalId())
	data.PrincipalType = types.StringValue(getPrincipalRateSettingsResp.GetPrincipalType())
	data.DefaultConcurrencyPercentage = types.Int32Value(getPrincipalRateSettingsResp.GetDefaultConcurrencyPercentage())
	data.DefaultPercentage = types.Int32Value(getPrincipalRateSettingsResp.GetDefaultPercentage())
	data.CreatedBy = types.StringValue(getPrincipalRateSettingsResp.GetCreatedBy())
	data.CreatedDate = types.StringValue(getPrincipalRateSettingsResp.GetCreatedDate().Format(time.RFC3339))
	data.LastUpdate = types.StringValue(getPrincipalRateSettingsResp.GetLastUpdate().Format(time.RFC3339))
	data.LastUpdatedBy = types.StringValue(getPrincipalRateSettingsResp.GetLastUpdatedBy())
	data.OrgId = types.StringValue(getPrincipalRateSettingsResp.GetOrgId())

	// Save updated Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
