package idaas

import (
	"context"
	v5okta "github.com/okta/okta-sdk-golang/v5/okta"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ resource.Resource                = &authenticatorSettings{}
	_ resource.ResourceWithConfigure   = &authenticatorSettings{}
	_ resource.ResourceWithImportState = &authenticatorSettings{}
)

func newAuthenticatorSettingsResource() resource.Resource {
	return &authenticatorSettings{}
}

type authenticatorSettings struct {
	*config.Config
}

type authenticatorSettingsResourceModel struct {
	Id                                   types.String `tfsdk:"id"`
	VerifyKnowledgeSecondWhen2faRequired types.Bool   `tfsdk:"verify_knowledge_second_when_2fa_required"`
}

func (r *authenticatorSettings) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_authenticator_settings"
}

func (r *authenticatorSettings) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The id property of an entitlement.",
			},
			"verify_knowledge_second_when_2fa_required": schema.BoolAttribute{
				Required:    true,
				Description: "If true, requires users to verify a possession factor before verifying a knowledge factor when the assurance requires two-factor authentication (2FA).",
			},
		},
	}
}

func (r *authenticatorSettings) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data authenticatorSettingsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic
	createAuthSettingsResp, _, err := r.OktaIDaaSClient.OktaSDKClientV5().AttackProtectionAPI.ReplaceAuthenticatorSettings(ctx).AuthenticatorSettings(createAuthenticatorSettingReq(data)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating AuthSettings",
			"Could not create AuthSettings, unexpected error: "+err.Error(),
		)
		return
	}

	data.Id = types.StringValue("default")
	data.VerifyKnowledgeSecondWhen2faRequired = types.BoolValue(createAuthSettingsResp.GetVerifyKnowledgeSecondWhen2faRequired())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func createAuthenticatorSettingReq(data authenticatorSettingsResourceModel) v5okta.AttackProtectionAuthenticatorSettings {
	return v5okta.AttackProtectionAuthenticatorSettings{
		VerifyKnowledgeSecondWhen2faRequired: data.VerifyKnowledgeSecondWhen2faRequired.ValueBoolPointer(),
	}
}

func (r *authenticatorSettings) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data authenticatorSettingsResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	readAuthSettingsResp, _, err := r.OktaIDaaSClient.OktaSDKClientV5().AttackProtectionAPI.GetAuthenticatorSettings(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading AuthSetting",
			"Could not create AuthSetting, unexpected error: "+err.Error(),
		)
		return
	}

	data.VerifyKnowledgeSecondWhen2faRequired = types.BoolValue(readAuthSettingsResp[0].GetVerifyKnowledgeSecondWhen2faRequired())
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *authenticatorSettings) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data authenticatorSettingsResourceModel
	var state authenticatorSettingsResourceModel
	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Id = state.Id
	// Update API call logic
	replaceAuthSettingsResp, _, err := r.OktaIDaaSClient.OktaSDKClientV5().AttackProtectionAPI.ReplaceAuthenticatorSettings(ctx).AuthenticatorSettings(createAuthenticatorSettingReq(data)).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating AuthSetting",
			"An error occurred while updating the AuthSetting: "+err.Error(),
		)
		return
	}
	data.VerifyKnowledgeSecondWhen2faRequired = types.BoolValue(replaceAuthSettingsResp.GetVerifyKnowledgeSecondWhen2faRequired())
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *authenticatorSettings) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddError(
		"Delete Not Supported",
		"The resource cannot be deleted via Terraform.",
	)
}

func (r *authenticatorSettings) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *authenticatorSettings) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(request, response)
}
