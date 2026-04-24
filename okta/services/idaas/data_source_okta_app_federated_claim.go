package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ datasource.DataSource = (*appFederatedClaimDataSource)(nil)

func newAppFederatedClaimDataSource() datasource.DataSource {
	return &appFederatedClaimDataSource{}
}

type appFederatedClaimDataSource struct {
	*config.Config
}

func (d *appFederatedClaimDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_federated_claim"
}

func (d *appFederatedClaimDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

func (d *appFederatedClaimDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "`id` used to specify the app feature ID. Its a combination of `app_id` and `name` separated by a forward slash (/).",
			},
			"app_id": schema.StringAttribute{
				Required:    true,
				Description: "`app_id` used to specify the app ID.",
			},
			"expression": schema.StringAttribute{
				Computed:    true,
				Description: "The Okta Expression Language expression to be evaluated at runtime.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the claim to be used in the produced token.",
			},
		},
	}
}

func (d *appFederatedClaimDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data appFederatedClaimModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	getAppFederatedClaimResp, _, err := d.OktaIDaaSClient.OktaSDKClientV6().ApplicationSSOFederatedClaimsAPI.GetFederatedClaim(ctx, data.AppID.ValueString(), data.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading app federated claim",
			"Could not read app federated claim, unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(updateReadAppFederatedClaimToState(&data, getAppFederatedClaimResp)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
