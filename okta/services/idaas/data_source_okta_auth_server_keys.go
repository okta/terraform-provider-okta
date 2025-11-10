package idaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ datasource.DataSource              = &authServerKeysDataSource{}
	_ datasource.DataSourceWithConfigure = &authServerKeysDataSource{}
)

type authServerKeysDataSource struct {
	*config.Config
}

type authServerKeysDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	KeyID        types.String `tfsdk:"key_id"`
	AuthServerID types.String `tfsdk:"auth_server_id"`
	Alg          types.String `tfsdk:"alg"`
	E            types.String `tfsdk:"e"`
	Kid          types.String `tfsdk:"kid"`
	N            types.String `tfsdk:"n"`
	Status       types.String `tfsdk:"status"`
	Use          types.String `tfsdk:"use"`
}

func newAuthServerKeysDataSource() datasource.DataSource {
	return &authServerKeysDataSource{}
}

func (d *authServerKeysDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_auth_server_keys"
}

func (d *authServerKeysDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

func (d *authServerKeysDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The internal identifier for this data source, required by Terraform to track state. This field does not exist in the Okta API response.",
			},
			"key_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the certificate key.",
			},
			"auth_server_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the authorization server.",
			},
			"alg": schema.StringAttribute{
				Computed:    true,
				Description: "The algorithm used with the Key. Valid value: RS256.",
			},
			"e": schema.StringAttribute{
				Computed:    true,
				Description: "RSA key value (public exponent) for Key binding.",
			},
			"kid": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier for the key.",
			},
			"n": schema.StringAttribute{
				Computed:    true,
				Description: "RSA modulus value that is used by both the public and private keys and provides a link between them.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "An ACTIVE Key is used to sign tokens issued by the authorization server. Supported values: ACTIVE, NEXT, or EXPIRED. A NEXT Key is the next Key that the authorization server uses to sign tokens when Keys are rotated. The NEXT Key might not be listed if it hasn't been generated. An EXPIRED Key is the previous Key that the authorization server used to sign tokens. The EXPIRED Key might not be listed if no Key has expired or the expired Key was deleted.",
			},
			"use": schema.StringAttribute{
				Computed:    true,
				Description: "Acceptable use of the key. Valid value: sig.",
			},
		},
	}
}

func (d *authServerKeysDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data authServerKeysDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic - List keys to find our specific client
	authServerJSONWebKey, _, err := getOktaV6ClientFromMetadata(d.Config).AuthorizationServerKeysAPI.GetAuthorizationServerKey(ctx, data.AuthServerID.ValueString(), data.KeyID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to read authorization server keys",
			err.Error(),
		)
		return
	}
	data.KeyID = types.StringValue(authServerJSONWebKey.GetKid())
	data.Alg = types.StringValue(authServerJSONWebKey.GetAlg())
	data.E = types.StringValue(authServerJSONWebKey.GetE())
	data.Kid = types.StringValue(authServerJSONWebKey.GetKid())
	data.N = types.StringValue(authServerJSONWebKey.GetN())
	data.Status = types.StringValue(authServerJSONWebKey.GetStatus())
	data.Use = types.StringValue(authServerJSONWebKey.GetUse())
	data.ID = types.StringValue(fmt.Sprintf("%s-%s", data.AuthServerID.ValueString(), data.KeyID.ValueString()))
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
