package idaas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

var _ datasource.DataSource = &authenticatorWebauthnCustomAAGUIDsDataSource{}

func newAuthenticatorWebauthnCustomAAGUIDsDataSource() datasource.DataSource {
	return &authenticatorWebauthnCustomAAGUIDsDataSource{}
}

type authenticatorWebauthnCustomAAGUIDsDataSource struct {
	*config.Config
}

type customAAGUIDsDataSourceModel struct {
	ID              types.String             `tfsdk:"id"`
	AuthenticatorID types.String             `tfsdk:"authenticator_id"`
	CustomAAGUIDs   []customAAGUIDItemModel  `tfsdk:"custom_aaguids"`
}

type customAAGUIDItemModel struct {
	AAGUID                       types.String                           `tfsdk:"aaguid"`
	Name                         types.String                           `tfsdk:"name"`
	AuthenticatorCharacteristics *authenticatorCharacteristicsDataModel `tfsdk:"authenticator_characteristics"`
}

type authenticatorCharacteristicsDataModel struct {
	FipsCompliant     types.Bool `tfsdk:"fips_compliant"`
	HardwareProtected types.Bool `tfsdk:"hardware_protected"`
	PlatformAttached  types.Bool `tfsdk:"platform_attached"`
}

func (d *authenticatorWebauthnCustomAAGUIDsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authenticator_webauthn_custom_aaguids"
}

func (d *authenticatorWebauthnCustomAAGUIDsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all custom AAGUIDs for a WebAuthn authenticator.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of this data source, set to the authenticator ID.",
			},
			"authenticator_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the WebAuthn authenticator.",
			},
			"custom_aaguids": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of custom AAGUIDs configured for this authenticator.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"aaguid": schema.StringAttribute{
							Computed:    true,
							Description: "The AAGUID identifier.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The product name associated with the AAGUID.",
						},
						"authenticator_characteristics": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "Properties of the custom AAGUID authenticator.",
							Attributes: map[string]schema.Attribute{
								"fips_compliant": schema.BoolAttribute{
									Computed:    true,
									Description: "Indicates whether the authenticator meets FIPS compliance requirements.",
								},
								"hardware_protected": schema.BoolAttribute{
									Computed:    true,
									Description: "Indicates whether the authenticator stores the private key on a hardware component.",
								},
								"platform_attached": schema.BoolAttribute{
									Computed:    true,
									Description: "Indicates whether the AAGUID is built into the authenticator or is external.",
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *authenticatorWebauthnCustomAAGUIDsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

func (d *authenticatorWebauthnCustomAAGUIDsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.Config != nil && fwproviderIsClassicOrg(ctx, d.Config) {
		resp.Diagnostics.Append(frameworkOIEOnlyFeatureError("data-sources", resources.OktaIDaaSAuthenticatorWebauthnCustomAAGUIDs)...)
		return
	}

	var state customAAGUIDsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	authenticatorID := state.AuthenticatorID.ValueString()
	client := d.Config.OktaIDaaSClient.OktaSDKClientV6()

	aaguids, _, err := client.AuthenticatorAPI.ListAllCustomAAGUIDs(ctx, authenticatorID).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error listing custom AAGUIDs",
			fmt.Sprintf("Could not list custom AAGUIDs for authenticator %s: %s", authenticatorID, err.Error()),
		)
		return
	}

	items := make([]customAAGUIDItemModel, 0, len(aaguids))
	for _, a := range aaguids {
		item := customAAGUIDItemModel{
			AAGUID: types.StringPointerValue(a.Aaguid),
			Name:   types.StringPointerValue(a.Name),
		}

		if a.AuthenticatorCharacteristics != nil {
			item.AuthenticatorCharacteristics = &authenticatorCharacteristicsDataModel{
				FipsCompliant:     types.BoolPointerValue(a.AuthenticatorCharacteristics.FipsCompliant),
				HardwareProtected: types.BoolPointerValue(a.AuthenticatorCharacteristics.HardwareProtected),
				PlatformAttached:  types.BoolPointerValue(a.AuthenticatorCharacteristics.PlatformAttached),
			}
		}

		items = append(items, item)
	}

	state.CustomAAGUIDs = items
	state.ID = types.StringValue(authenticatorID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
