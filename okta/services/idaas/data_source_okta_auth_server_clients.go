package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ datasource.DataSource              = &authServerClientsDataSource{}
	_ datasource.DataSourceWithConfigure = &authServerClientsDataSource{}
)

type authServerClientsDataSource struct {
	*config.Config
}

type authServerClientsDataSourceModel struct {
	Id           types.String `tfsdk:"id"`
	AuthServerId types.String `tfsdk:"auth_server_id"`
	ClientID     types.String `tfsdk:"client_id"`
	Created      types.String `tfsdk:"created"`
	ExpiresAt    types.String `tfsdk:"expires_at"`
	Issuer       types.String `tfsdk:"issuer"`
	LastUpdated  types.String `tfsdk:"last_updated"`
	Scopes       types.List   `tfsdk:"scopes"`
	Status       types.String `tfsdk:"status"`
	UserId       types.String `tfsdk:"user_id"`
}

func newAuthServerClientsDataSource() datasource.DataSource {
	return &authServerClientsDataSource{}
}

func (d *authServerClientsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_auth_server_clients"
}

func (d *authServerClientsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

func (d *authServerClientsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the token.",
			},
			"auth_server_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the authorization server.",
			},
			"client_id": schema.StringAttribute{
				Required:    true,
				Description: "The client ID of the app.",
			},
			"created": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the object was created.",
			},
			"expires_at": schema.StringAttribute{
				Computed:    true,
				Description: "Expiration time of the OAuth 2.0 Token.",
			},
			"issuer": schema.StringAttribute{
				Computed:    true,
				Description: "The complete URL of the authorization server that issued the Token",
			},
			"last_updated": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the object was last updated.",
			},
			"scopes": schema.ListAttribute{
				Computed:    true,
				Description: "The scope names attached to the Token.",
				ElementType: types.StringType,
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status URI of the OAuth 2.0 application.",
			},
			"user_id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the user associated with the Token.",
			},
		},
	}
}

func (d *authServerClientsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Retrieve a refresh token for a client
	var data authServerClientsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic - List clients to find our specific client
	OAuth2RefreshToken, _, err := getOktaV5ClientFromMetadata(d.Config).AuthorizationServerClientsAPI.GetRefreshTokenForAuthorizationServerAndClient(ctx, data.AuthServerId.ValueString(), data.ClientID.ValueString(), data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to read authorization server clients",
			err.Error(),
		)
		return
	}
	data.ClientID = types.StringValue(OAuth2RefreshToken.GetClientId())
	data.Created = types.StringValue(OAuth2RefreshToken.GetCreated().String())
	data.ExpiresAt = types.StringValue(OAuth2RefreshToken.GetExpiresAt().String())
	data.Id = types.StringValue(OAuth2RefreshToken.GetId()) // Token ID
	data.Issuer = types.StringValue(OAuth2RefreshToken.GetIssuer())
	data.LastUpdated = types.StringValue(OAuth2RefreshToken.GetLastUpdated().String())
	scopes, diags := types.ListValueFrom(ctx, types.StringType, OAuth2RefreshToken.GetScopes())
	if diags.HasError() {
		return
	}
	data.Scopes = scopes
	data.Status = types.StringValue(OAuth2RefreshToken.GetStatus())
	data.UserId = types.StringValue(OAuth2RefreshToken.GetUserId())

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
