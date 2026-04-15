package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

var _ datasource.DataSource = &authenticatorMethodWebauthnDataSource{}

func newAuthenticatorMethodWebauthnDataSource() datasource.DataSource {
	return &authenticatorMethodWebauthnDataSource{}
}

type authenticatorMethodWebauthnDataSource struct {
	*config.Config
}

type aaguidGroupDataSourceModel struct {
	Name    types.String   `tfsdk:"name"`
	Aaguids []types.String `tfsdk:"aaguids"`
}

type authenticatorMethodWebauthnDataSourceModel struct {
	ID                             types.String                  `tfsdk:"id"`
	AuthenticatorID                types.String                  `tfsdk:"authenticator_id"`
	Status                         types.String                  `tfsdk:"status"`
	UserVerification               types.String                  `tfsdk:"user_verification"`
	UserVerificationForVerify      types.String                  `tfsdk:"user_verification_for_verify"`
	Attachment                     types.String                  `tfsdk:"attachment"`
	EnableAutofillUI               types.Bool                    `tfsdk:"enable_autofill_ui"`
	ResidentKeyRequirement         types.String                  `tfsdk:"resident_key_requirement"`
	ShowSignInWithAPasskeyButton   types.Bool                    `tfsdk:"show_sign_in_with_a_passkey_button"`
	CertBasedAttestationValidation types.Bool                    `tfsdk:"cert_based_attestation_validation"`
	HardwareProtected              types.Bool                    `tfsdk:"hardware_protected"`
	FipsCompliant                  types.Bool                    `tfsdk:"fips_compliant"`
	AllowSyncablePasskeys          types.Bool                    `tfsdk:"allow_syncable_passkeys"`
	AaguidGroups                   []aaguidGroupDataSourceModel  `tfsdk:"aaguid_groups"`
}

func (d *authenticatorMethodWebauthnDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authenticator_method_webauthn"
}

func (d *authenticatorMethodWebauthnDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Reads the WebAuthn authenticator method settings from Okta.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The resource ID.",
			},
			"authenticator_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the authenticator.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the authenticator method.",
			},
			"user_verification": schema.StringAttribute{
				Computed:    true,
				Description: "User verification setting for enrollment.",
			},
			"user_verification_for_verify": schema.StringAttribute{
				Computed:    true,
				Description: "User verification setting for verification (sign-in).",
			},
			"attachment": schema.StringAttribute{
				Computed:    true,
				Description: "Method attachment.",
			},
			"enable_autofill_ui": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the passkeys autofill UI is enabled.",
			},
			"resident_key_requirement": schema.StringAttribute{
				Computed:    true,
				Description: "Resident key requirement setting.",
			},
			"show_sign_in_with_a_passkey_button": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the Sign in with a Passkey button is shown.",
			},
			"cert_based_attestation_validation": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether certificate-based attestation validation is enabled.",
			},
			"hardware_protected": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the authenticator requires hardware-protected key storage.",
			},
			"fips_compliant": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the authenticator is required to be FIPS compliant.",
			},
			"allow_syncable_passkeys": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether syncable passkeys are allowed.",
			},
		},
		Blocks: map[string]schema.Block{
			"aaguid_groups": schema.ListNestedBlock{
				Description: "The AAGUID groups available to the WebAuthn authenticator.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "A name to identify the AAGUID group.",
						},
						"aaguids": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "A list of FIDO2 AAGUIDs in the group.",
						},
					},
				},
			},
		},
	}
}

func (d *authenticatorMethodWebauthnDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.Config = dataSourceConfiguration(req, resp)
}

func (d *authenticatorMethodWebauthnDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.Config != nil && fwproviderIsClassicOrg(ctx, d.Config) {
		resp.Diagnostics.Append(frameworkOIEOnlyFeatureError("data-sources", resources.OktaIDaaSAuthenticatorMethodWebauthn)...)
		return
	}

	var data authenticatorMethodWebauthnDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	authenticatorID := data.AuthenticatorID.ValueString()

	result, _, err := d.OktaIDaaSClient.OktaSDKClientV6().AuthenticatorAPI.GetAuthenticatorMethod(ctx, authenticatorID, "webauthn").Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading WebAuthn Authenticator Method",
			"Could not read WebAuthn authenticator method for authenticator "+authenticatorID+": "+err.Error(),
		)
		return
	}

	webauthn := result.AuthenticatorMethodWebAuthn
	if webauthn == nil {
		resp.Diagnostics.AddError(
			"Error Reading WebAuthn Authenticator Method",
			"API response did not contain a WebAuthn authenticator method for authenticator "+authenticatorID,
		)
		return
	}

	data.ID = data.AuthenticatorID
	data.Status = types.StringPointerValue(webauthn.Status)

	settings := webauthn.Settings
	if settings != nil {
		data.UserVerification = types.StringPointerValue(settings.UserVerification)
		data.UserVerificationForVerify = types.StringPointerValue(settings.UserVerificationForVerify)
		data.Attachment = types.StringPointerValue(settings.Attachment)
		data.EnableAutofillUI = types.BoolPointerValue(settings.EnableAutofillUI)
		data.ResidentKeyRequirement = types.StringPointerValue(settings.ResidentKeyRequirement)
		data.ShowSignInWithAPasskeyButton = types.BoolPointerValue(settings.ShowSignInWithAPasskeyButton)
		data.CertBasedAttestationValidation = types.BoolPointerValue(settings.CertBasedAttestationValidation)
		data.HardwareProtected = types.BoolPointerValue(settings.HardwareProtected)
		data.FipsCompliant = types.BoolPointerValue(settings.FipsCompliant)
		data.AllowSyncablePasskeys = types.BoolPointerValue(settings.AllowSyncablePasskeys)

		groups := make([]aaguidGroupDataSourceModel, 0, len(settings.AaguidGroups))
		for _, g := range settings.AaguidGroups {
			aaguidValues := make([]types.String, 0, len(g.Aaguids))
			for _, a := range g.Aaguids {
				aaguidValues = append(aaguidValues, types.StringValue(a))
			}
			groups = append(groups, aaguidGroupDataSourceModel{
				Name:    types.StringPointerValue(g.Name),
				Aaguids: aaguidValues,
			})
		}
		data.AaguidGroups = groups
	} else {
		data.UserVerification = types.StringNull()
		data.UserVerificationForVerify = types.StringNull()
		data.Attachment = types.StringNull()
		data.EnableAutofillUI = types.BoolNull()
		data.ResidentKeyRequirement = types.StringNull()
		data.ShowSignInWithAPasskeyButton = types.BoolNull()
		data.CertBasedAttestationValidation = types.BoolNull()
		data.HardwareProtected = types.BoolNull()
		data.FipsCompliant = types.BoolNull()
		data.AllowSyncablePasskeys = types.BoolNull()
		data.AaguidGroups = nil
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
