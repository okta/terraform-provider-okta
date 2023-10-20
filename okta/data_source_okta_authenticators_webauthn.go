package okta

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewAuthenticatorWebauthnDataSource() datasource.DataSource {
	return &authenticatorWebauthnDatasource{}
}

type authenticatorWebauthnDatasource struct {
	*Config
}

// TODU
type AuthenticatorWebauthnDatasourceModel struct {
	ID          types.String `tfsdk:"id"`
	AuthDevices types.List   `tfsdk:"auth_devices"`
}

// TODU
type AuthenticatorWebauthn struct {
	AAGUID    types.String `tfsdk:"aaguid"`
	ModelName types.String `tfsdk:"model_name"`
}

// TODU
func (d *authenticatorWebauthnDatasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authenticators_webauthn"
}

// TODU
func (d *authenticatorWebauthnDatasource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier",
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"auth_devices": schema.ListNestedBlock{
				Description: "",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"aaguid": schema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"model_name": schema.StringAttribute{
							Description: "",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the datasource.
func (d *authenticatorWebauthnDatasource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	p, ok := req.ProviderData.(*Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Datasource Configure Type",
			fmt.Sprintf("Expected *Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.Config = p
}

// TODU
func (d *authenticatorWebauthnDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state AuthenticatorWebauthnDatasourceModel

	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// TODU
	webauthnCatalog, err := d.Config.oktaSDKsupplementClient.ListWebauthnCatalog(ctx, "webauthn", d.Config.orgName, d.Config.domain, d.Config.apiToken, d.Config.oktaSDKClientV3.GetConfig().HTTPClient)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving webauthn catalog",
			fmt.Sprintf("Error returned: %s", err.Error()),
		)
		return
	}
	catalog := make([]AuthenticatorWebauthn, 0)
	for _, webauthnDevice := range webauthnCatalog {
		var authenticatorWebauthn AuthenticatorWebauthn
		authenticatorWebauthn.AAGUID = types.StringValue(webauthnDevice.AAGUID)
		authenticatorWebauthn.ModelName = types.StringValue(webauthnDevice.ModelName)
		catalog = append(catalog, authenticatorWebauthn)
	}
	authDeviceValue, diags := types.ListValueFrom(ctx, state.AuthDevices.ElementType(ctx), catalog)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	state.AuthDevices = authDeviceValue
	// TODU
	state.ID = types.StringValue("a")
	// Save data into state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
