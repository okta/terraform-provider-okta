package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ datasource.DataSource = &userRiskDataSource{}

func newUserRiskDataSource() datasource.DataSource {
	return &userRiskDataSource{}
}

func (d *userRiskDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

type userRiskDataSource struct {
	*config.Config
}

type userRiskDataSourceModel struct {
	ID        types.String `tfsdk:"id"`
	UserID    types.String `tfsdk:"user_id"`
	RiskLevel types.String `tfsdk:"risk_level"`
	Reason    types.String `tfsdk:"reason"`
}

func (d *userRiskDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_risk"
}

func (d *userRiskDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Gets a user's risk level in Okta.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of this resource (same as user_id).",
			},
			"user_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the user.",
			},
			"risk_level": schema.StringAttribute{
				Computed:    true,
				Description: "Risk level of the user. Possible values: `HIGH`, `LOW`, `NONE`. `NONE` indicates no risk level has been set.",
			},
			"reason": schema.StringAttribute{
				Computed:    true,
				Description: "Reason for the risk level, if available.",
			},
		},
	}
}

func (d *userRiskDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data userRiskDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	userId := data.UserID.ValueString()

	d.Logger.Info("reading user risk data source", "user_id", userId)

	riskResp, _, err := d.OktaIDaaSClient.OktaSDKClientV6().UserRiskAPI.GetUserRisk(ctx, userId).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading user risk",
			"Could not read user risk for user "+userId+": "+err.Error(),
		)
		return
	}

	var riskLevel, reason string
	if riskResp.UserRiskLevelExists != nil {
		riskLevel = riskResp.UserRiskLevelExists.GetRiskLevel()
		reason = riskResp.UserRiskLevelExists.GetReason()
	} else if riskResp.UserRiskLevelNone != nil {
		riskLevel = "NONE"
		reason = ""
	}

	data.ID = types.StringValue(userId)
	data.UserID = types.StringValue(userId)
	data.RiskLevel = types.StringValue(riskLevel)
	data.Reason = types.StringValue(reason)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
