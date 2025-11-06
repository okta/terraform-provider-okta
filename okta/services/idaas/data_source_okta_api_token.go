package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ datasource.DataSource = &apiTokenDataSource{}

func newAPITokenDataSource() datasource.DataSource {
	return &apiTokenDataSource{}
}

type apiTokenDataSource struct {
	*config.Config
}

func (d *apiTokenDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_token"
}

func (d *apiTokenDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

func (d *apiTokenDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the API token.",
			},
			"user_id": schema.StringAttribute{
				Computed:    true,
				Description: "The userId of the user who created the API Token.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the API token.",
			},
			"created": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the API token was created.",
			},
			"client_name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the API token client",
			},
		},
		Blocks: map[string]schema.Block{
			"network": schema.SingleNestedBlock{
				Description: "The Network Condition of the API Token.",
				Attributes: map[string]schema.Attribute{
					"connection": schema.StringAttribute{
						Computed:    true,
						Description: "The connection type of the Network Condition.",
					},
					"exclude": schema.ListAttribute{
						Optional:    true,
						Computed:    true,
						Description: "The IP address the excluded zone.",
						ElementType: types.StringType,
					},
					"include": schema.ListAttribute{
						Optional:    true,
						Computed:    true,
						Description: "The IP address the included zone.",
						ElementType: types.StringType,
					},
				},
			},
		},
	}
}

func (d *apiTokenDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data apiTokenResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	getAPITokenResp, _, err := d.OktaIDaaSClient.OktaSDKClientV5().ApiTokenAPI.GetApiToken(ctx, data.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error in getting API token",
			err.Error(),
		)
		return
	}
	mapAPITokenToState(ctx, getAPITokenResp, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
