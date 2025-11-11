package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	v5okta "github.com/okta/okta-sdk-golang/v5/okta"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ resource.Resource                = &securityEventsProviderResource{}
	_ resource.ResourceWithConfigure   = &securityEventsProviderResource{}
	_ resource.ResourceWithImportState = &securityEventsProviderResource{}
)

func newSecurityEventsProviderResource() resource.Resource {
	return &securityEventsProviderResource{}
}

type securityEventsProviderResource struct {
	*config.Config
}

// The overall model for the Terraform resource
type securityEventsProviderModel struct {
	Id        types.String   `tfsdk:"id"`
	Name      types.String   `tfsdk:"name"`
	Type      types.String   `tfsdk:"type"`
	IsEnabled types.String   `tfsdk:"is_enabled"`
	Settings  *settingsModel `tfsdk:"settings"` // This is a SingleNestedAttribute in the schema
	Status    types.String   `tfsdk:"status"`
}

// The model for the required 'settings' attribute, which handles the one-of logic
type settingsModel struct {
	// Both are Optional, and validators in the schema enforce that exactly one group is provided.

	// Fields for "Provider with well-known URL setting"
	WellKnownUrl types.String `tfsdk:"well_known_url"`

	// Fields for "Provider with issuer and JWKS settings"
	Issuer  types.String `tfsdk:"issuer"`
	JwksUrl types.String `tfsdk:"jwks_url"`
}

func (r *securityEventsProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// 1. Read the state from the plan (Terraform configuration)
	var data securityEventsProviderModel
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createSSFReceiverResp, _, err := r.OktaIDaaSClient.OktaSDKClientV5().SSFReceiverAPI.CreateSecurityEventsProviderInstance(ctx).Instance(createSecurityEventsProviderRequest(data)).Execute()
	if err != nil {
		return
	}

	var activateResp *v5okta.SecurityEventsProviderResponse
	if data.IsEnabled.ValueStringPointer() != nil && data.IsEnabled.ValueString() == "ACTIVE" {
		activateResp, _, err = r.OktaIDaaSClient.OktaSDKClientV5().SSFReceiverAPI.ActivateSecurityEventsProviderInstance(ctx, createSSFReceiverResp.GetId()).Execute()
		data.IsEnabled = types.StringValue("ACTIVE")
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Activating Security Events Provider",
				"Could not activate Security Events Provider ID "+data.Id.ValueString()+": "+err.Error(),
			)
			return
		}
	} else if data.IsEnabled.ValueStringPointer() != nil && data.IsEnabled.ValueString() == "INACTIVE" {
		activateResp, _, err = r.OktaIDaaSClient.OktaSDKClientV5().SSFReceiverAPI.DeactivateSecurityEventsProviderInstance(ctx, createSSFReceiverResp.GetId()).Execute()
		data.IsEnabled = types.StringValue("INACTIVE")
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Deactivating Security Events Provider",
				"Could not deactivate Security Events Provider ID "+data.Id.ValueString()+": "+err.Error(),
			)
			return
		}
	}

	resp.Diagnostics.Append(applySecurityEventsProviderToState(activateResp, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *securityEventsProviderResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var data securityEventsProviderModel

	// Read the state from the Terraform configuration
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	getSecurityEventsProviderResp, _, err := r.OktaIDaaSClient.OktaSDKClientV5().SSFReceiverAPI.GetSecurityEventsProviderInstance(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		response.Diagnostics.AddError(
			"Error Reading Security Events Provider",
			"Could not read Security Events Provider ID "+data.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	response.Diagnostics.Append(applySecurityEventsProviderToState(getSecurityEventsProviderResp, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	// Save updated state
	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *securityEventsProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state securityEventsProviderModel

	// Read the state from the Terraform configuration
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, _, err := r.OktaIDaaSClient.OktaSDKClientV5().SSFReceiverAPI.ReplaceSecurityEventsProviderInstance(ctx, state.Id.ValueString()).Instance(createSecurityEventsProviderRequest(data)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Security Events Provider",
			"Could not update Security Events Provider ID "+data.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	var activateResp *v5okta.SecurityEventsProviderResponse
	if data.IsEnabled.ValueStringPointer() != nil && data.IsEnabled.ValueString() == "ACTIVE" {
		activateResp, _, err = r.OktaIDaaSClient.OktaSDKClientV5().SSFReceiverAPI.ActivateSecurityEventsProviderInstance(ctx, state.Id.ValueString()).Execute()
		data.IsEnabled = types.StringValue("ACTIVE")
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Activating Security Events Provider",
				"Could not activate Security Events Provider ID "+data.Id.ValueString()+": "+err.Error(),
			)
			return
		}
	} else if data.IsEnabled.ValueStringPointer() != nil && data.IsEnabled.ValueString() == "INACTIVE" {
		activateResp, _, err = r.OktaIDaaSClient.OktaSDKClientV5().SSFReceiverAPI.DeactivateSecurityEventsProviderInstance(ctx, state.Id.ValueString()).Execute()
		data.IsEnabled = types.StringValue("INACTIVE")
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Deactivating Security Events Provider",
				"Could not deactivate Security Events Provider ID "+data.Id.ValueString()+": "+err.Error(),
			)
			return
		}
	}

	resp.Diagnostics.Append(applySecurityEventsProviderToState(activateResp, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *securityEventsProviderResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var data securityEventsProviderModel

	// Read the state from the Terraform configuration
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	_, err := r.OktaIDaaSClient.OktaSDKClientV5().SSFReceiverAPI.DeleteSecurityEventsProviderInstance(ctx, data.Id.ValueString()).Execute()
	if err != nil {
		response.Diagnostics.AddError(
			"Error Deleting Security Events Provider",
			"Could not delete Security Events Provider ID "+data.Id.ValueString()+": "+err.Error(),
		)
		return
	}
}

func applySecurityEventsProviderToState(resp *v5okta.SecurityEventsProviderResponse, s *securityEventsProviderModel) diag.Diagnostics {
	var diags diag.Diagnostics
	s.Id = types.StringValue(resp.GetId())
	s.Name = types.StringValue(resp.GetName())
	s.Status = types.StringValue(resp.GetStatus())
	s.IsEnabled = types.StringValue(resp.GetStatus())
	s.Type = types.StringValue(resp.GetType())
	s.IsEnabled = types.StringValue(resp.GetStatus())
	settings := resp.GetSettings()
	s.Settings = &settingsModel{}
	if settings.HasWellKnownUrl() {
		s.Settings.WellKnownUrl = types.StringValue(settings.GetWellKnownUrl())
	} else if settings.HasJwksUrl() && settings.HasIssuer() {
		s.Settings.Issuer = types.StringValue(settings.GetIssuer())
		s.Settings.JwksUrl = types.StringValue(settings.GetJwksUrl())
	}
	return diags
}

func createSecurityEventsProviderRequest(model securityEventsProviderModel) v5okta.SecurityEventsProviderRequest {
	securityEventsProviderReq := v5okta.SecurityEventsProviderRequest{
		Name: model.Name.ValueString(),
		Type: model.Type.ValueString(),
	}

	securityEventsProviderRequestSettings := v5okta.SecurityEventsProviderRequestSettings{}

	if model.Settings.WellKnownUrl.ValueStringPointer() != nil {
		securityEventsProviderRequestSettings.SecurityEventsProviderSettingsSSFCompliant = &v5okta.SecurityEventsProviderSettingsSSFCompliant{}
		securityEventsProviderRequestSettings.SecurityEventsProviderSettingsSSFCompliant.WellKnownUrl = model.Settings.WellKnownUrl.ValueString()
	} else if model.Settings.Issuer.ValueStringPointer() != nil {
		securityEventsProviderRequestSettings.SecurityEventsProviderSettingsNonSSFCompliant = &v5okta.SecurityEventsProviderSettingsNonSSFCompliant{}
		securityEventsProviderRequestSettings.SecurityEventsProviderSettingsNonSSFCompliant.Issuer = model.Settings.Issuer.ValueString()
		securityEventsProviderRequestSettings.SecurityEventsProviderSettingsNonSSFCompliant.JwksUrl = model.Settings.JwksUrl.ValueString()
	}

	securityEventsProviderReq.Settings = securityEventsProviderRequestSettings
	return securityEventsProviderReq
}

func (r *securityEventsProviderResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
}

func (r *securityEventsProviderResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(request, response)
}

func (r *securityEventsProviderResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_security_events_provider"
}

func (r *securityEventsProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a Security Events Provider instance for signal ingestion.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier of this instance.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the Security Events Provider instance.",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "The application type of the Security Events Provider.",
			},
			"is_enabled": schema.StringAttribute{
				Required:    true,
				Description: "Whether or not the Security Events Provider is enabled.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Indicates whether the Security Events Provider is active or not.",
			},
		},
		Blocks: map[string]schema.Block{
			// The required 'settings' block that houses the "one-of" logic
			"settings": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					// --- Well-Known URL Setting ---
					"well_known_url": schema.StringAttribute{
						Optional:    true,
						Description: "The published well-known URL of the Security Events Provider (the SSF transmitter).",
						Validators: []validator.String{
							stringvalidator.LengthAtMost(1000),
						},
					},

					// --- Issuer and JWKS Settings ---
					"issuer": schema.StringAttribute{
						Optional:    true,
						Description: "Issuer URL. Use with jwks_url",
						Validators: []validator.String{
							stringvalidator.LengthAtMost(700),
						},
					},
					"jwks_url": schema.StringAttribute{
						Optional:    true,
						Description: "The public URL where the JWKS public key is uploaded. Use with issuer.",
						Validators: []validator.String{
							stringvalidator.LengthAtMost(1000),
						},
					},
				},
				Description: "Information about the Security Events Provider for signal ingestion.",
			},
		},
	}
}
