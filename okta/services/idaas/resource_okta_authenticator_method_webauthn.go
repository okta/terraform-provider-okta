package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	v6okta "github.com/okta/okta-sdk-golang/v6/okta"
	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/okta/resources"
)

type authenticatorMethodWebauthnResourceModel struct {
	ID                             types.String                  `tfsdk:"id"`
	AuthenticatorID                types.String                  `tfsdk:"authenticator_id"`
	UserVerification               types.String                  `tfsdk:"user_verification"`
	UserVerificationForVerify      types.String                  `tfsdk:"user_verification_for_verify"`
	Attachment                     types.String                  `tfsdk:"attachment"`
	AaguidGroups                   []aaguidGroupModel            `tfsdk:"aaguid_group"`
	EnableAutofillUI               types.Bool                    `tfsdk:"enable_autofill_ui"`
	ResidentKeyRequirement         types.String                  `tfsdk:"resident_key_requirement"`
	ShowSignInWithAPasskeyButton   types.Bool                    `tfsdk:"show_sign_in_with_a_passkey_button"`
	CertBasedAttestationValidation types.Bool                    `tfsdk:"cert_based_attestation_validation"`
	HardwareProtected              types.Bool                    `tfsdk:"hardware_protected"`
	FipsCompliant                  types.Bool                    `tfsdk:"fips_compliant"`
	AllowSyncablePasskeys          types.Bool                    `tfsdk:"allow_syncable_passkeys"`
	Status                         types.String                  `tfsdk:"status"`
}

type aaguidGroupModel struct {
	Name    types.String `tfsdk:"name"`
	Aaguids []types.String `tfsdk:"aaguids"`
}

type authenticatorMethodWebauthnResource struct {
	config *config.Config
}

func newAuthenticatorMethodWebauthnResource() resource.Resource {
	return &authenticatorMethodWebauthnResource{}
}

func (r *authenticatorMethodWebauthnResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authenticator_method_webauthn"
}

func (r *authenticatorMethodWebauthnResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages WebAuthn authenticator method settings in Okta.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The authenticator ID (same as authenticator_id).",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"authenticator_id": schema.StringAttribute{
				Description: "The ID of the authenticator.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"user_verification": schema.StringAttribute{
				Description: "User verification setting for enrollment. Valid values: DISCOURAGED, PREFERRED, REQUIRED.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_verification_for_verify": schema.StringAttribute{
				Description: "User verification setting for verification (sign-in).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"attachment": schema.StringAttribute{
				Description: "Method attachment setting.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enable_autofill_ui": schema.BoolAttribute{
				Description: "Enables the passkeys autofill UI.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"resident_key_requirement": schema.StringAttribute{
				Description: "Resident key requirement setting.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"show_sign_in_with_a_passkey_button": schema.BoolAttribute{
				Description: "Whether the Sign in with a Passkey button is shown on the Sign-In Widget.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"cert_based_attestation_validation": schema.BoolAttribute{
				Description: "Whether certificate-based attestation validation is enabled.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"hardware_protected": schema.BoolAttribute{
				Description: "Whether the authenticator is required to store the private key on a hardware component.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"fips_compliant": schema.BoolAttribute{
				Description: "Whether the authenticator is required to be FIPS compliant.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"allow_syncable_passkeys": schema.BoolAttribute{
				Description: "Whether syncable passkeys are allowed.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				Description: "The status of the authenticator method.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"aaguid_group": schema.ListNestedBlock{
				Description: "AAGUID groups for the WebAuthn authenticator.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "A name to identify the group of AAGUIDs.",
							Required:    true,
						},
						"aaguids": schema.ListAttribute{
							Description: "A list of FIDO2 AAGUIDs in the group.",
							Required:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

func (r *authenticatorMethodWebauthnResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.config = resourceConfiguration(req, resp)
}

func (r *authenticatorMethodWebauthnResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data authenticatorMethodWebauthnResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if r.config != nil && fwproviderIsClassicOrg(ctx, r.config) {
		resp.Diagnostics.Append(frameworkResourceOIEOnlyFeatureError(resources.OktaIDaaSAuthenticatorMethodWebauthn)...)
		return
	}

	authenticatorID := data.AuthenticatorID.ValueString()
	client := r.config.OktaIDaaSClient.OktaSDKClientV6()

	// Read the current method state
	current, _, err := client.AuthenticatorAPI.GetAuthenticatorMethod(ctx, authenticatorID, "webauthn").Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error reading WebAuthn authenticator method", err.Error())
		return
	}

	if current.AuthenticatorMethodWebAuthn == nil {
		resp.Diagnostics.AddError("Error reading WebAuthn authenticator method", "API response did not contain a WebAuthn method")
		return
	}

	// Build the replacement body from plan data, preserving current status
	webauthnMethod := buildWebAuthnMethodFromModel(&data, current.AuthenticatorMethodWebAuthn)
	body := v6okta.AuthenticatorMethodWebAuthnAsListAuthenticatorMethods200ResponseInner(webauthnMethod)

	result, _, err := client.AuthenticatorAPI.ReplaceAuthenticatorMethod(ctx, authenticatorID, "webauthn").ListAuthenticatorMethods200ResponseInner(body).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error creating WebAuthn authenticator method", err.Error())
		return
	}

	if result.AuthenticatorMethodWebAuthn == nil {
		resp.Diagnostics.AddError("Error creating WebAuthn authenticator method", "API response did not contain a WebAuthn method")
		return
	}

	resp.Diagnostics.Append(mapWebAuthnMethodToState(result.AuthenticatorMethodWebAuthn, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *authenticatorMethodWebauthnResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state authenticatorMethodWebauthnResourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	authenticatorID := state.AuthenticatorID.ValueString()
	client := r.config.OktaIDaaSClient.OktaSDKClientV6()

	result, _, err := client.AuthenticatorAPI.GetAuthenticatorMethod(ctx, authenticatorID, "webauthn").Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error reading WebAuthn authenticator method", err.Error())
		return
	}

	if result.AuthenticatorMethodWebAuthn == nil {
		resp.Diagnostics.AddError("Error reading WebAuthn authenticator method", "API response did not contain a WebAuthn method")
		return
	}

	resp.Diagnostics.Append(mapWebAuthnMethodToState(result.AuthenticatorMethodWebAuthn, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *authenticatorMethodWebauthnResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data authenticatorMethodWebauthnResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if r.config != nil && fwproviderIsClassicOrg(ctx, r.config) {
		resp.Diagnostics.Append(frameworkResourceOIEOnlyFeatureError(resources.OktaIDaaSAuthenticatorMethodWebauthn)...)
		return
	}

	authenticatorID := data.AuthenticatorID.ValueString()
	client := r.config.OktaIDaaSClient.OktaSDKClientV6()

	// Read the current method state to preserve status and other server-managed fields
	current, _, err := client.AuthenticatorAPI.GetAuthenticatorMethod(ctx, authenticatorID, "webauthn").Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error reading WebAuthn authenticator method", err.Error())
		return
	}

	if current.AuthenticatorMethodWebAuthn == nil {
		resp.Diagnostics.AddError("Error reading WebAuthn authenticator method", "API response did not contain a WebAuthn method")
		return
	}

	webauthnMethod := buildWebAuthnMethodFromModel(&data, current.AuthenticatorMethodWebAuthn)
	body := v6okta.AuthenticatorMethodWebAuthnAsListAuthenticatorMethods200ResponseInner(webauthnMethod)

	result, _, err := client.AuthenticatorAPI.ReplaceAuthenticatorMethod(ctx, authenticatorID, "webauthn").ListAuthenticatorMethods200ResponseInner(body).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error updating WebAuthn authenticator method", err.Error())
		return
	}

	if result.AuthenticatorMethodWebAuthn == nil {
		resp.Diagnostics.AddError("Error updating WebAuthn authenticator method", "API response did not contain a WebAuthn method")
		return
	}

	resp.Diagnostics.Append(mapWebAuthnMethodToState(result.AuthenticatorMethodWebAuthn, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *authenticatorMethodWebauthnResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state authenticatorMethodWebauthnResourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	authenticatorID := state.AuthenticatorID.ValueString()
	client := r.config.OktaIDaaSClient.OktaSDKClientV6()

	// Methods can't be deleted — reset to default settings
	defaultSettings := v6okta.NewAuthenticatorMethodWebAuthnAllOfSettings()
	defaultSettings.AaguidGroups = []v6okta.AAGUIDGroupObject{}

	webauthnMethod := v6okta.NewAuthenticatorMethodWebAuthn()
	webauthnMethod.Settings = defaultSettings

	body := v6okta.AuthenticatorMethodWebAuthnAsListAuthenticatorMethods200ResponseInner(webauthnMethod)

	_, _, err := client.AuthenticatorAPI.ReplaceAuthenticatorMethod(ctx, authenticatorID, "webauthn").ListAuthenticatorMethods200ResponseInner(body).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error resetting WebAuthn authenticator method to defaults", err.Error())
		return
	}
}

func (r *authenticatorMethodWebauthnResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("authenticator_id"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

// buildWebAuthnMethodFromModel constructs an AuthenticatorMethodWebAuthn from the Terraform model.
func buildWebAuthnMethodFromModel(data *authenticatorMethodWebauthnResourceModel, current *v6okta.AuthenticatorMethodWebAuthn) *v6okta.AuthenticatorMethodWebAuthn {
	method := v6okta.NewAuthenticatorMethodWebAuthn()

	// Preserve the current status and type
	method.Status = current.Status
	method.Type = current.Type

	settings := v6okta.NewAuthenticatorMethodWebAuthnAllOfSettings()

	if !data.UserVerification.IsNull() && !data.UserVerification.IsUnknown() {
		settings.UserVerification = data.UserVerification.ValueStringPointer()
	} else if current.Settings != nil {
		settings.UserVerification = current.Settings.UserVerification
	}

	if !data.UserVerificationForVerify.IsNull() && !data.UserVerificationForVerify.IsUnknown() {
		settings.UserVerificationForVerify = data.UserVerificationForVerify.ValueStringPointer()
	} else if current.Settings != nil {
		settings.UserVerificationForVerify = current.Settings.UserVerificationForVerify
	}

	if !data.Attachment.IsNull() && !data.Attachment.IsUnknown() {
		settings.Attachment = data.Attachment.ValueStringPointer()
	} else if current.Settings != nil {
		settings.Attachment = current.Settings.Attachment
	}

	if !data.EnableAutofillUI.IsNull() && !data.EnableAutofillUI.IsUnknown() {
		settings.EnableAutofillUI = data.EnableAutofillUI.ValueBoolPointer()
	} else if current.Settings != nil {
		settings.EnableAutofillUI = current.Settings.EnableAutofillUI
	}

	if !data.ResidentKeyRequirement.IsNull() && !data.ResidentKeyRequirement.IsUnknown() {
		settings.ResidentKeyRequirement = data.ResidentKeyRequirement.ValueStringPointer()
	} else if current.Settings != nil {
		settings.ResidentKeyRequirement = current.Settings.ResidentKeyRequirement
	}

	if !data.ShowSignInWithAPasskeyButton.IsNull() && !data.ShowSignInWithAPasskeyButton.IsUnknown() {
		settings.ShowSignInWithAPasskeyButton = data.ShowSignInWithAPasskeyButton.ValueBoolPointer()
	} else if current.Settings != nil {
		settings.ShowSignInWithAPasskeyButton = current.Settings.ShowSignInWithAPasskeyButton
	}

	if !data.CertBasedAttestationValidation.IsNull() && !data.CertBasedAttestationValidation.IsUnknown() {
		settings.CertBasedAttestationValidation = data.CertBasedAttestationValidation.ValueBoolPointer()
	} else if current.Settings != nil {
		settings.CertBasedAttestationValidation = current.Settings.CertBasedAttestationValidation
	}

	if !data.HardwareProtected.IsNull() && !data.HardwareProtected.IsUnknown() {
		settings.HardwareProtected = data.HardwareProtected.ValueBoolPointer()
	} else if current.Settings != nil {
		settings.HardwareProtected = current.Settings.HardwareProtected
	}

	if !data.FipsCompliant.IsNull() && !data.FipsCompliant.IsUnknown() {
		settings.FipsCompliant = data.FipsCompliant.ValueBoolPointer()
	} else if current.Settings != nil {
		settings.FipsCompliant = current.Settings.FipsCompliant
	}

	if !data.AllowSyncablePasskeys.IsNull() && !data.AllowSyncablePasskeys.IsUnknown() {
		settings.AllowSyncablePasskeys = data.AllowSyncablePasskeys.ValueBoolPointer()
	} else if current.Settings != nil {
		settings.AllowSyncablePasskeys = current.Settings.AllowSyncablePasskeys
	}

	// AAGUID groups
	if data.AaguidGroups != nil {
		groups := make([]v6okta.AAGUIDGroupObject, len(data.AaguidGroups))
		for i, g := range data.AaguidGroups {
			aaguids := make([]string, len(g.Aaguids))
			for j, a := range g.Aaguids {
				aaguids[j] = a.ValueString()
			}
			groups[i] = v6okta.AAGUIDGroupObject{
				Name:    g.Name.ValueStringPointer(),
				Aaguids: aaguids,
			}
		}
		settings.AaguidGroups = groups
	} else if current.Settings != nil {
		settings.AaguidGroups = current.Settings.AaguidGroups
	}

	// Preserve RpId from current if it exists (read-only / server-managed)
	if current.Settings != nil {
		settings.RpId = current.Settings.RpId
	}

	method.Settings = settings
	return method
}

// mapWebAuthnMethodToState maps the API response to the Terraform state model.
func mapWebAuthnMethodToState(method *v6okta.AuthenticatorMethodWebAuthn, state *authenticatorMethodWebauthnResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	state.ID = state.AuthenticatorID
	state.Status = types.StringPointerValue(method.Status)

	if method.Settings != nil {
		s := method.Settings
		state.UserVerification = types.StringPointerValue(s.UserVerification)
		state.UserVerificationForVerify = types.StringPointerValue(s.UserVerificationForVerify)
		state.Attachment = types.StringPointerValue(s.Attachment)
		state.EnableAutofillUI = types.BoolPointerValue(s.EnableAutofillUI)
		state.ResidentKeyRequirement = types.StringPointerValue(s.ResidentKeyRequirement)
		state.ShowSignInWithAPasskeyButton = types.BoolPointerValue(s.ShowSignInWithAPasskeyButton)
		state.CertBasedAttestationValidation = types.BoolPointerValue(s.CertBasedAttestationValidation)
		state.HardwareProtected = types.BoolPointerValue(s.HardwareProtected)
		state.FipsCompliant = types.BoolPointerValue(s.FipsCompliant)
		state.AllowSyncablePasskeys = types.BoolPointerValue(s.AllowSyncablePasskeys)

		if len(s.AaguidGroups) > 0 {
			groups := make([]aaguidGroupModel, len(s.AaguidGroups))
			for i, g := range s.AaguidGroups {
				aaguids := make([]types.String, len(g.Aaguids))
				for j, a := range g.Aaguids {
					aaguids[j] = types.StringValue(a)
				}
				groups[i] = aaguidGroupModel{
					Name:    types.StringPointerValue(g.Name),
					Aaguids: aaguids,
				}
			}
			state.AaguidGroups = groups
		} else {
			state.AaguidGroups = nil
		}
	}

	return diags
}
