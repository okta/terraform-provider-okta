package idaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/okta/terraform-provider-okta/okta/config"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &defaultSigninPageDataSource{}
	_ datasource.DataSourceWithConfigure = &defaultSigninPageDataSource{}
)

func newDefaultSigninPageDataSource() datasource.DataSource {
	return &defaultSigninPageDataSource{}
}

type defaultSigninPageDataSource struct {
	*config.Config
}

func (d *defaultSigninPageDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_signin_page"
}

func (d *defaultSigninPageDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	newSchema := dataSourceSignInSchema
	newSchema.Description = "Retrieve the default signin page of a brand"
	resp.Schema = newSchema
}

func (d *defaultSigninPageDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	config, ok := req.ProviderData.(*config.Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Datasource Configure Type",
			fmt.Sprintf("Expected *Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.Config = config
}

func (d *defaultSigninPageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data signinPageModel

	resp.Diagnostics.Append(resp.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	defaultSigninPage, _, err := d.OktaIDaaSClient.OktaSDKClientV3().CustomizationAPI.GetDefaultSignInPage(ctx, data.BrandID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving default signin page",
			fmt.Sprintf("Error returned: %s", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(mapSignInPageToState(ctx, defaultSigninPage, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
