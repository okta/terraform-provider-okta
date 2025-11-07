package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/okta-sdk-golang/v5/okta"
	"github.com/okta/terraform-provider-okta/okta/config"
)

type appConnectionsDataSource struct {
	config *config.Config
}

type ProfileDataModel struct {
	AuthScheme types.String `tfsdk:"auth_scheme"`
}

type appConnectionsDataModel struct {
	ID         types.String      `tfsdk:"id"`
	Profile    *ProfileDataModel `tfsdk:"profile"`
	Status     types.String      `tfsdk:"status"`
	AuthScheme types.String      `tfsdk:"auth_scheme"`
	BaseURL    types.String      `tfsdk:"base_url"`
}

func newAppConnectionsDataSource() datasource.DataSource {
	return &appConnectionsDataSource{}
}

func (r *appConnectionsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_connection"
}

func (r *appConnectionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	r.config = dataSourceConfiguration(req, resp)
}

func (r *appConnectionsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The application ID.",
				Required:    true,
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Provisioning connection status.",
			},
			"auth_scheme": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "A token is used to authenticate with the app. This property is only returned for the TOKEN authentication scheme.",
			},
			"base_url": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "The base URL for the provisioning connection.",
			},
		},
		Blocks: map[string]schema.Block{
			"profile": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"auth_scheme": schema.StringAttribute{
						Computed:    true,
						Description: "Defines the method of authentication",
					},
				},
			},
		},
		Description: "Previews the SSO SAML metadata for an application.",
	}
}

func (r *appConnectionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data appConnectionsDataModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readProvisionConnection, _, err := r.config.OktaIDaaSClient.OktaSDKClientV5().ApplicationConnectionsAPI.GetDefaultProvisioningConnectionForApplication(ctx, data.ID.ValueString()).Execute()
	if err != nil {
		return
	}

	resp.Diagnostics.Append(mapProvisionConnectionToState(readProvisionConnection, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func mapProvisionConnectionToState(connResp *okta.ProvisioningConnectionResponse, state *appConnectionsDataModel) diag.Diagnostics {
	var diags diag.Diagnostics
	state.BaseURL = types.StringValue(connResp.GetBaseUrl())
	state.Status = types.StringValue(connResp.GetStatus())
	state.AuthScheme = types.StringValue(connResp.GetAuthScheme())
	state.Profile = &ProfileDataModel{}
	if profile, ok := connResp.GetProfileOk(); ok {
		state.Profile.AuthScheme = types.StringValue(profile.GetAuthScheme())
	}
	return diags
}
