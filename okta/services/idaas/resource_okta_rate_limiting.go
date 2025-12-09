package idaas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	tfpath "github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	v5okta "github.com/okta/okta-sdk-golang/v5/okta"
	"github.com/okta/terraform-provider-okta/okta/config"
)

var (
	_ resource.Resource                = &rateLimitResource{}
	_ resource.ResourceWithConfigure   = &rateLimitResource{}
	_ resource.ResourceWithImportState = &rateLimitResource{}
)

var _ resource.Resource = &rateLimitResource{}

type rateLimitResource struct {
	*config.Config
}

type useCaseModeOverrides struct {
	LoginPage       types.String `tfsdk:"login_page"`
	OAuth2Authorize types.String `tfsdk:"oauth2_authorize"`
	OieAppIntent    types.String `tfsdk:"oie_app_intent"`
}

type rateLimitResourceModel struct {
	Id                   types.String          `tfsdk:"id"`
	DefaultMode          types.String          `tfsdk:"default_mode"`
	UseCaseModeOverrides *useCaseModeOverrides `tfsdk:"use_case_mode_overrides"`
}

const (
	DISABLE = "DISABLE"
	ENFORCE = "ENFORCE"
	PREVIEW = "PREVIEW"
)

func newRateLimitResource() resource.Resource {
	return &rateLimitResource{}
}

func (r *rateLimitResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rate_limiting"
}

func (r *rateLimitResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *rateLimitResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, tfpath.Root("id"), req, resp)
}

func (r *rateLimitResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"default_mode": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						DISABLE,
						ENFORCE,
						PREVIEW,
					}...),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"use_case_mode_overrides": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"login_page": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf([]string{
								DISABLE,
								ENFORCE,
								PREVIEW,
							}...),
						},
					},
					"oauth2_authorize": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf([]string{
								DISABLE,
								ENFORCE,
								PREVIEW,
							}...),
						},
					},
					"oie_app_intent": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf([]string{
								DISABLE,
								ENFORCE,
								PREVIEW,
							}...),
						},
					},
				},
				Description: "A map of Per-Client Rate Limit Use Case to the applicable PerClientRateLimitMode.Overrides the defaultMode property for the specified use cases.",
			},
		},
	}
}

func (r *rateLimitResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data rateLimitResourceModel

	// Read Terraform plan Data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic
	clientRateLimitSettingsResp, _, err := r.OktaIDaaSClient.OktaSDKClientV5().RateLimitSettingsAPI.ReplaceRateLimitSettingsPerClient(ctx).PerClientRateLimitSettings(buildPerClientRateLimitSettings(data)).Execute()
	if err != nil {
		return
	}

	// Example Data value setting
	data.Id = types.StringValue("rate_limiting")
	data.DefaultMode = types.StringValue(clientRateLimitSettingsResp.GetDefaultMode())
	overrides := clientRateLimitSettingsResp.GetUseCaseModeOverrides()
	modeOverrides := &useCaseModeOverrides{}
	_, ok := overrides.GetLOGIN_PAGEOk()
	if ok {
		modeOverrides.LoginPage = types.StringValue(overrides.GetLOGIN_PAGE())
	}
	_, ok = overrides.GetOIE_APP_INTENTOk()
	if ok {
		modeOverrides.OieAppIntent = types.StringValue(overrides.GetOIE_APP_INTENT())
	}
	_, ok = overrides.GetOAUTH2AUTHORIZEOk()
	if ok {
		modeOverrides.OAuth2Authorize = types.StringValue(overrides.GetOAUTH2AUTHORIZE())
	}
	data.UseCaseModeOverrides = modeOverrides

	// Save Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *rateLimitResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data rateLimitResourceModel

	// Read Terraform prior state Data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	getRateLimitResp, _, err := r.OktaIDaaSClient.OktaSDKClientV5().RateLimitSettingsAPI.GetRateLimitSettingsPerClient(ctx).Execute()
	if err != nil {
		return
	}

	data.DefaultMode = types.StringValue(getRateLimitResp.GetDefaultMode())
	overrides := getRateLimitResp.GetUseCaseModeOverrides()
	modeOverrides := &useCaseModeOverrides{}
	_, ok := overrides.GetLOGIN_PAGEOk()
	if ok {
		modeOverrides.LoginPage = types.StringValue(overrides.GetLOGIN_PAGE())
	}
	_, ok = overrides.GetOIE_APP_INTENTOk()
	if ok {
		modeOverrides.OieAppIntent = types.StringValue(overrides.GetOIE_APP_INTENT())
	}
	_, ok = overrides.GetOAUTH2AUTHORIZEOk()
	if ok {
		modeOverrides.OAuth2Authorize = types.StringValue(overrides.GetOAUTH2AUTHORIZE())
	}
	data.UseCaseModeOverrides = modeOverrides
	data.Id = types.StringValue("rate_limiting")
	// Save updated Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *rateLimitResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data rateLimitResourceModel

	// Read Terraform plan Data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic
	clientRateLimitSettingsResp, _, err := r.OktaIDaaSClient.OktaSDKClientV5().RateLimitSettingsAPI.ReplaceRateLimitSettingsPerClient(ctx).PerClientRateLimitSettings(buildPerClientRateLimitSettings(data)).Execute()
	if err != nil {
		return
	}

	// Example Data value setting
	data.Id = types.StringValue("rate_limiting")
	data.DefaultMode = types.StringValue(clientRateLimitSettingsResp.GetDefaultMode())
	overrides := clientRateLimitSettingsResp.GetUseCaseModeOverrides()
	modeOverrides := &useCaseModeOverrides{}
	_, ok := overrides.GetLOGIN_PAGEOk()
	if ok {
		modeOverrides.LoginPage = types.StringValue(overrides.GetLOGIN_PAGE())
		data.UseCaseModeOverrides = modeOverrides
	}
	_, ok = overrides.GetOIE_APP_INTENTOk()
	if ok {
		modeOverrides.OieAppIntent = types.StringValue(overrides.GetOIE_APP_INTENT())
		data.UseCaseModeOverrides = modeOverrides
	}
	_, ok = overrides.GetOAUTH2AUTHORIZEOk()
	if ok {
		modeOverrides.OAuth2Authorize = types.StringValue(overrides.GetOAUTH2AUTHORIZE())
		data.UseCaseModeOverrides = modeOverrides
	}

	// Save updated Data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *rateLimitResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning(
		"Delete Not Supported",
		"This resource cannot be deleted via Terraform.",
	)
}

func buildPerClientRateLimitSettings(data rateLimitResourceModel) v5okta.PerClientRateLimitSettings {
	rateLimitSettings := v5okta.PerClientRateLimitSettings{
		DefaultMode: data.DefaultMode.ValueString(),
	}
	useCaseOverrides := &v5okta.PerClientRateLimitSettingsUseCaseModeOverrides{}
	if data.UseCaseModeOverrides != nil {
		if !data.UseCaseModeOverrides.LoginPage.IsNull() {
			useCaseOverrides.LOGIN_PAGE = data.UseCaseModeOverrides.LoginPage.ValueStringPointer()
		}
		if !data.UseCaseModeOverrides.OAuth2Authorize.IsNull() {
			useCaseOverrides.OAUTH2AUTHORIZE = data.UseCaseModeOverrides.OAuth2Authorize.ValueStringPointer()
		}

		if !data.UseCaseModeOverrides.OieAppIntent.IsNull() {
			useCaseOverrides.OIE_APP_INTENT = data.UseCaseModeOverrides.OieAppIntent.ValueStringPointer()
		}
	}

	if data.UseCaseModeOverrides != nil && useCaseOverrides != nil {
		rateLimitSettings.SetUseCaseModeOverrides(*useCaseOverrides)
	}
	return rateLimitSettings
}
