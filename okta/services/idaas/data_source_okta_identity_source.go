package idaas

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var _ datasource.DataSource = &identitySourceDataSource{}

func newIdentitySourceDataSource() datasource.DataSource {
	return &identitySourceDataSource{}
}

type identitySourceDataSource struct {
	*config.Config
}

type identitySourceDataSourceModel struct {
	Id               types.String `tfsdk:"id"`
	Created          types.String `tfsdk:"created"`
	LastUpdated      types.String `tfsdk:"last_updated"`
	IdentitySourceId types.String `tfsdk:"identity_source_id"` // Indicates the minimum required SKU to manage the campaign. Values can be `BASIC` and `PREMIUM`.
	ImportType       types.String `tfsdk:"import_type"`
	Status           types.String `tfsdk:"status"`
}

func (d *identitySourceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_identity_source"
}

func (d *identitySourceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the identity source session.",
			},
			"created": schema.StringAttribute{
				Computed: true,
			},
			"identity_source_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the custom identity source for which the session is created.",
			},
			"import_type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of import. All imports are `INCREMENTAL` imports.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The current status of the identity source session.",
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *identitySourceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data identitySourceDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	getIdentitySourceResp, _, err := d.OktaIDaaSClient.OktaSDKClientV5().IdentitySourceAPI.GetIdentitySourceSession(ctx, data.IdentitySourceId.ValueString(), data.Id.ValueString()).Execute()
	if err != nil {
		return
	}

	data.Id = types.StringValue(getIdentitySourceResp.GetId())
	data.IdentitySourceId = types.StringValue(getIdentitySourceResp.GetIdentitySourceId())
	data.Created = types.StringValue(getIdentitySourceResp.GetCreated().String())
	data.LastUpdated = types.StringValue(getIdentitySourceResp.GetLastUpdated().String())
	data.ImportType = types.StringValue(getIdentitySourceResp.GetImportType())
	data.Status = types.StringValue(getIdentitySourceResp.GetStatus())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
