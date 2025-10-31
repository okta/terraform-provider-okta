package idaas

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
	"time"
)

var _ datasource.DataSource = &appTokenDataSource{}

func newAppTokenDataSource() datasource.DataSource {
	return &appTokenDataSource{}
}

type appTokenDataSource struct {
	*config.Config
}

type appTokenDataSourceModel struct {
	Id        types.String `tfsdk:"id"`
	ClientId  types.String `tfsdk:"client_id"`
	UserId    types.String `tfsdk:"user_id"`
	Status    types.String `tfsdk:"status"`
	Created   types.String `tfsdk:"created"`
	ExpiresAt types.String `tfsdk:"expires_at"`
	Scopes    types.List   `tfsdk:"scopes"`
	Issuer    types.String `tfsdk:"issuer"`
}

func (d *appTokenDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The unique Okta ID of this key record",
			},
			"client_id": schema.StringAttribute{
				Required:    true,
				Description: "The unique Okta ID of the application associated with this token.",
			},
			"user_id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique Okta ID of the user associated with this token.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the token.",
			},
			"created": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the object was created.",
			},
			"expires_at": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the object was expired.",
			},
			"scopes": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "The scope names attached to the Token.",
			},
			"issuer": schema.StringAttribute{
				Computed:    true,
				Description: "The complete URL of the authorization server that issued the Token.",
			},
		},
	}
}

func (d *appTokenDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_token"
}

func (d *appTokenDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

func (d *appTokenDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data appTokenDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	getAppTokenRes, _, err := d.OktaIDaaSClient.OktaSDKClientV5().ApplicationTokensAPI.GetOAuth2TokenForApplication(ctx, data.ClientId.ValueString(), data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error reading application token", "Could not read application token, unexpected error: "+err.Error())
		return
	}

	data.Id = types.StringValue(getAppTokenRes.GetId())
	data.ClientId = types.StringValue(getAppTokenRes.GetClientId())
	data.UserId = types.StringValue(getAppTokenRes.GetUserId())
	data.Status = types.StringValue(getAppTokenRes.GetStatus())
	data.Created = types.StringValue(getAppTokenRes.GetCreated().Format(time.RFC3339))
	data.ExpiresAt = types.StringValue(getAppTokenRes.GetExpiresAt().Format(time.RFC3339))
	scopes := make([]attr.Value, len(getAppTokenRes.GetScopes()))
	for i, scope := range getAppTokenRes.GetScopes() {
		scopes[i] = types.StringValue(scope)
	}
	s, diags := types.ListValue(types.StringType, scopes)
	resp.Diagnostics.Append(diags...)
	data.Scopes = s
	data.Issuer = types.StringValue(getAppTokenRes.GetIssuer())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...) // Save data into Terraform state
}
