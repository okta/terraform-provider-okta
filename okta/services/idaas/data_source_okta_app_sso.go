package idaas

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

type appSSODataSource struct {
	config *config.Config
}

type appSSODataModel struct {
	ID       types.String `tfsdk:"id"`
	Metadata types.String `tfsdk:"metadata"`
}

func newAppSSODataSource() datasource.DataSource {
	return &appSSODataSource{}
}

func (r *appSSODataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_sso"
}

func (r *appSSODataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	r.config = dataSourceConfiguration(req, resp)
}

func (r *appSSODataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Realm Assignment ID.",
				Required:    true,
			},
			"metadata": schema.StringAttribute{
				Computed: true,
			},
		},
		Description: "Previews the SSO SAML metadata for an application.",
	}
}

func (r *appSSODataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state appSSODataModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	samlMetadata, _, err := r.config.OktaIDaaSClient.OktaSDKClientV5().ApplicationSSOAPI.PreviewSAMLmetadataForApplication(ctx, state.ID.ValueString()).Execute()
	if err != nil {
		return
	}

	resp.Diagnostics.Append(mapSAMLMetadataToState(samlMetadata, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func mapSAMLMetadataToState(metadata string, state *appSSODataModel) diag.Diagnostics {
	var diags diag.Diagnostics
	state.Metadata = types.StringValue(metadata)
	return diags
}
