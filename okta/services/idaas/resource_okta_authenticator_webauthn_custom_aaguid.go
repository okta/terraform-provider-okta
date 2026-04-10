package idaas

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	v6okta "github.com/okta/okta-sdk-golang/v6/okta"
	"github.com/okta/terraform-provider-okta/okta/config"
	"github.com/okta/terraform-provider-okta/okta/resources"
	"github.com/okta/terraform-provider-okta/okta/utils"
)

type authenticatorWebauthnCustomAAGUIDResourceModel struct {
	ID                          types.String `tfsdk:"id"`
	AuthenticatorID             types.String `tfsdk:"authenticator_id"`
	AAGUID                      types.String `tfsdk:"aaguid"`
	Name                        types.String `tfsdk:"name"`
	AuthenticatorCharacteristics types.Object `tfsdk:"authenticator_characteristics"`
	AttestationRootCertificates  types.List   `tfsdk:"attestation_root_certificate"`
}

type authenticatorCharacteristicsModel struct {
	FIPSCompliant     types.Bool `tfsdk:"fips_compliant"`
	HardwareProtected types.Bool `tfsdk:"hardware_protected"`
	PlatformAttached  types.Bool `tfsdk:"platform_attached"`
}

type attestationRootCertificateModel struct {
	X5C     types.String `tfsdk:"x5c"`
	X5TS256 types.String `tfsdk:"x5t_s256"`
	Issuer  types.String `tfsdk:"issuer"`
	Expiry  types.String `tfsdk:"expiry"`
}

type authenticatorWebauthnCustomAAGUIDResource struct {
	config *config.Config
}

func newAuthenticatorWebauthnCustomAAGUIDResource() resource.Resource {
	return &authenticatorWebauthnCustomAAGUIDResource{}
}

func (r *authenticatorWebauthnCustomAAGUIDResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authenticator_webauthn_custom_aaguid"
}

func (r *authenticatorWebauthnCustomAAGUIDResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a custom AAGUID on a WebAuthn authenticator in Okta.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The AAGUID value, used as the resource ID.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"authenticator_id": schema.StringAttribute{
				Description: "The ID of the WebAuthn authenticator.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"aaguid": schema.StringAttribute{
				Description: "The AAGUID (Authenticator Attestation Global Unique Identifier) string.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The product name associated with the AAGUID.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"authenticator_characteristics": schema.SingleNestedAttribute{
				Description: "Characteristics of the authenticator.",
				Optional:    true,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"fips_compliant": schema.BoolAttribute{
						Description: "Whether the authenticator is FIPS compliant.",
						Optional:    true,
						Computed:    true,
					},
					"hardware_protected": schema.BoolAttribute{
						Description: "Whether the authenticator keys are hardware protected.",
						Optional:    true,
						Computed:    true,
					},
					"platform_attached": schema.BoolAttribute{
						Description: "Whether the authenticator is a platform authenticator.",
						Optional:    true,
						Computed:    true,
					},
				},
			},
			"attestation_root_certificate": schema.ListNestedAttribute{
				Description: "List of attestation root certificates.",
				Optional:    true,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"x5c": schema.StringAttribute{
							Description: "Base64-encoded X.509 certificate.",
							Required:    true,
						},
						"x5t_s256": schema.StringAttribute{
							Description: "SHA-256 thumbprint of the certificate.",
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"issuer": schema.StringAttribute{
							Description: "Certificate issuer.",
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"expiry": schema.StringAttribute{
							Description: "Certificate expiry.",
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
					},
				},
			},
		},
	}
}

func (r *authenticatorWebauthnCustomAAGUIDResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.config = resourceConfiguration(req, resp)
}

func (r *authenticatorWebauthnCustomAAGUIDResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.config != nil && fwproviderIsClassicOrg(ctx, r.config) {
		resp.Diagnostics.Append(frameworkResourceOIEOnlyFeatureError(resources.OktaIDaaSAuthenticatorWebauthnCustomAAGUID)...)
		return
	}

	var data authenticatorWebauthnCustomAAGUIDResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := v6okta.CustomAAGUIDCreateRequestObject{
		Aaguid: data.AAGUID.ValueStringPointer(),
	}

	resp.Diagnostics.Append(setCreateRequestCharacteristics(ctx, &data, &body)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(setCreateRequestCertificates(ctx, &data, &body)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, _, err := r.config.OktaIDaaSClient.OktaSDKClientV6().AuthenticatorAPI.CreateCustomAAGUID(ctx, data.AuthenticatorID.ValueString()).CustomAAGUIDCreateRequestObject(body).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error creating custom AAGUID", err.Error())
		return
	}

	resp.Diagnostics.Append(mapCustomAAGUIDToState(ctx, result, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *authenticatorWebauthnCustomAAGUIDResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state authenticatorWebauthnCustomAAGUIDResourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, apiResp, err := r.config.OktaIDaaSClient.OktaSDKClientV6().AuthenticatorAPI.GetCustomAAGUID(ctx, state.AuthenticatorID.ValueString(), state.AAGUID.ValueString()).Execute()
	if err != nil {
		if utils.SuppressErrorOn404_V6(apiResp, err) == nil {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading custom AAGUID", err.Error())
		return
	}

	resp.Diagnostics.Append(mapCustomAAGUIDToState(ctx, result, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *authenticatorWebauthnCustomAAGUIDResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.config != nil && fwproviderIsClassicOrg(ctx, r.config) {
		resp.Diagnostics.Append(frameworkResourceOIEOnlyFeatureError(resources.OktaIDaaSAuthenticatorWebauthnCustomAAGUID)...)
		return
	}

	var data authenticatorWebauthnCustomAAGUIDResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := v6okta.CustomAAGUIDUpdateRequestObject{
		Name: data.Name.ValueStringPointer(),
	}

	resp.Diagnostics.Append(setUpdateRequestCharacteristics(ctx, &data, &body)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(setUpdateRequestCertificates(ctx, &data, &body)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, _, err := r.config.OktaIDaaSClient.OktaSDKClientV6().AuthenticatorAPI.ReplaceCustomAAGUID(ctx, data.AuthenticatorID.ValueString(), data.AAGUID.ValueString()).CustomAAGUIDUpdateRequestObject(body).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error updating custom AAGUID", err.Error())
		return
	}

	resp.Diagnostics.Append(mapCustomAAGUIDToState(ctx, result, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *authenticatorWebauthnCustomAAGUIDResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.config != nil && fwproviderIsClassicOrg(ctx, r.config) {
		resp.Diagnostics.Append(frameworkResourceOIEOnlyFeatureError(resources.OktaIDaaSAuthenticatorWebauthnCustomAAGUID)...)
		return
	}

	var state authenticatorWebauthnCustomAAGUIDResourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.config.OktaIDaaSClient.OktaSDKClientV6().AuthenticatorAPI.DeleteCustomAAGUID(ctx, state.AuthenticatorID.ValueString(), state.AAGUID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error deleting custom AAGUID", err.Error())
		return
	}
}

func (r *authenticatorWebauthnCustomAAGUIDResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importID := req.ID
	if importID == "" {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Import ID cannot be empty. Expected format: authenticator_id/aaguid",
		)
		return
	}

	parts := strings.Split(importID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID format",
			fmt.Sprintf("Expected import ID format 'authenticator_id/aaguid', got '%s'", importID),
		)
		return
	}

	authenticatorID := parts[0]
	aaguid := parts[1]

	if authenticatorID == "" || aaguid == "" {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Both authenticator_id and aaguid must be provided in import ID",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("authenticator_id"), authenticatorID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("aaguid"), aaguid)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), aaguid)...)
}

func setCreateRequestCharacteristics(ctx context.Context, data *authenticatorWebauthnCustomAAGUIDResourceModel, body *v6okta.CustomAAGUIDCreateRequestObject) diag.Diagnostics {
	var diags diag.Diagnostics
	if data.AuthenticatorCharacteristics.IsNull() || data.AuthenticatorCharacteristics.IsUnknown() {
		return diags
	}

	var chars authenticatorCharacteristicsModel
	diags.Append(data.AuthenticatorCharacteristics.As(ctx, &chars, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return diags
	}

	body.AuthenticatorCharacteristics = &v6okta.AAGUIDAuthenticatorCharacteristics{
		FipsCompliant:     chars.FIPSCompliant.ValueBoolPointer(),
		HardwareProtected: chars.HardwareProtected.ValueBoolPointer(),
		PlatformAttached:  chars.PlatformAttached.ValueBoolPointer(),
	}
	return diags
}

func setUpdateRequestCharacteristics(ctx context.Context, data *authenticatorWebauthnCustomAAGUIDResourceModel, body *v6okta.CustomAAGUIDUpdateRequestObject) diag.Diagnostics {
	var diags diag.Diagnostics
	if data.AuthenticatorCharacteristics.IsNull() || data.AuthenticatorCharacteristics.IsUnknown() {
		return diags
	}

	var chars authenticatorCharacteristicsModel
	diags.Append(data.AuthenticatorCharacteristics.As(ctx, &chars, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return diags
	}

	body.AuthenticatorCharacteristics = &v6okta.AAGUIDAuthenticatorCharacteristics{
		FipsCompliant:     chars.FIPSCompliant.ValueBoolPointer(),
		HardwareProtected: chars.HardwareProtected.ValueBoolPointer(),
		PlatformAttached:  chars.PlatformAttached.ValueBoolPointer(),
	}
	return diags
}

func setCreateRequestCertificates(ctx context.Context, data *authenticatorWebauthnCustomAAGUIDResourceModel, body *v6okta.CustomAAGUIDCreateRequestObject) diag.Diagnostics {
	var diags diag.Diagnostics
	if data.AttestationRootCertificates.IsNull() || data.AttestationRootCertificates.IsUnknown() {
		return diags
	}

	var certs []attestationRootCertificateModel
	diags.Append(data.AttestationRootCertificates.ElementsAs(ctx, &certs, false)...)
	if diags.HasError() {
		return diags
	}

	for _, cert := range certs {
		body.AttestationRootCertificates = append(body.AttestationRootCertificates, v6okta.AttestationRootCertificatesRequestInner{
			X5c: cert.X5C.ValueStringPointer(),
		})
	}
	return diags
}

func setUpdateRequestCertificates(ctx context.Context, data *authenticatorWebauthnCustomAAGUIDResourceModel, body *v6okta.CustomAAGUIDUpdateRequestObject) diag.Diagnostics {
	var diags diag.Diagnostics
	if data.AttestationRootCertificates.IsNull() || data.AttestationRootCertificates.IsUnknown() {
		return diags
	}

	var certs []attestationRootCertificateModel
	diags.Append(data.AttestationRootCertificates.ElementsAs(ctx, &certs, false)...)
	if diags.HasError() {
		return diags
	}

	for _, cert := range certs {
		body.AttestationRootCertificates = append(body.AttestationRootCertificates, v6okta.AttestationRootCertificatesRequestInner{
			X5c: cert.X5C.ValueStringPointer(),
		})
	}
	return diags
}

func mapCustomAAGUIDToState(ctx context.Context, result *v6okta.CustomAAGUIDResponseObject, data *authenticatorWebauthnCustomAAGUIDResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	data.ID = types.StringPointerValue(result.Aaguid)
	data.AAGUID = types.StringPointerValue(result.Aaguid)
	data.Name = types.StringPointerValue(result.Name)

	if result.AuthenticatorCharacteristics != nil {
		chars := authenticatorCharacteristicsModel{
			FIPSCompliant:     types.BoolPointerValue(result.AuthenticatorCharacteristics.FipsCompliant),
			HardwareProtected: types.BoolPointerValue(result.AuthenticatorCharacteristics.HardwareProtected),
			PlatformAttached:  types.BoolPointerValue(result.AuthenticatorCharacteristics.PlatformAttached),
		}
		charsObj, d := types.ObjectValueFrom(ctx, authenticatorCharacteristicsAttrTypes, chars)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}
		data.AuthenticatorCharacteristics = charsObj
	}

	if result.AttestationRootCertificates != nil {
		var certModels []attestationRootCertificateModel
		for _, cert := range result.AttestationRootCertificates {
			certModels = append(certModels, attestationRootCertificateModel{
				X5C:     types.StringPointerValue(cert.X5c),
				X5TS256: types.StringPointerValue(cert.X5tS256),
				Issuer:  types.StringPointerValue(cert.Iss),
				Expiry:  types.StringPointerValue(cert.Exp),
			})
		}
		certList, d := types.ListValueFrom(ctx, attestationRootCertificateObjectType, certModels)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}
		data.AttestationRootCertificates = certList
	}

	return diags
}

var authenticatorCharacteristicsAttrTypes = map[string]attr.Type{
	"fips_compliant":     types.BoolType,
	"hardware_protected": types.BoolType,
	"platform_attached":  types.BoolType,
}

var attestationRootCertificateObjectType = types.ObjectType{
	AttrTypes: attestationRootCertificateAttrTypes,
}

var attestationRootCertificateAttrTypes = map[string]attr.Type{
	"x5c":      types.StringType,
	"x5t_s256": types.StringType,
	"issuer":   types.StringType,
	"expiry":   types.StringType,
}
