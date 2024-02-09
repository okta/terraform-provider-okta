package okta

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/okta/okta-sdk-golang/v4/okta"
)

func NewBrandResource() resource.Resource {
	return &brandResource{}
}

type brandResource struct {
	*Config
}

type brandResourceModel struct {
	BrandID                         types.String `tfsdk:"brand_id"`
	ID                              types.String `tfsdk:"id"`
	Name                            types.String `tfsdk:"name"`
	IsDefault                       types.Bool   `tfsdk:"is_default"`
	EmailDomainID                   types.String `tfsdk:"email_domain_id"`
	Locale                          types.String `tfsdk:"locale"`
	AgreeToCustomPrivacyPolicy      types.Bool   `tfsdk:"agree_to_custom_privacy_policy"`
	CustomPrivacyPolicyURL          types.String `tfsdk:"custom_privacy_policy_url"`
	RemovePoweredByOkta             types.Bool   `tfsdk:"remove_powered_by_okta"`
	DefaultAppAppInstanceID         types.String `tfsdk:"default_app_app_instance_id"`
	DefaultAppAppLinkName           types.String `tfsdk:"default_app_app_link_name"`
	DefaultAppClassicApplicationURI types.String `tfsdk:"default_app_classic_application_uri"`
	Links                           types.String `tfsdk:"links"`
}

func (r *brandResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_brand"
}

func (r *brandResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages brand",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Brand id",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Brand name",
				Required:    true,
			},
			"is_default": schema.BoolAttribute{
				Description: "Is this the default brand",
				Computed:    true,
			},
			"email_domain_id": schema.StringAttribute{
				Description: "Email Domain ID tied to this brand",
				Computed:    true,
			},
			"locale": schema.StringAttribute{
				Description: "The language specified as an IETF BCP 47 language tag",
				Optional:    true,
			},
			"agree_to_custom_privacy_policy": schema.BoolAttribute{
				Description: "Consent for updating the custom privacy policy URL.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.Bool{
					boolvalidator.AlsoRequires(path.Expressions{
						path.MatchRoot("custom_privacy_policy_url"),
					}...),
				},
			},
			"custom_privacy_policy_url": schema.StringAttribute{
				Description: "Custom privacy policy URL",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.Expressions{
						path.MatchRoot("agree_to_custom_privacy_policy"),
					}...),
				},
			},
			"remove_powered_by_okta": schema.BoolAttribute{
				Description: `Removes "Powered by Okta" from the Okta-hosted sign-in page and "© 2021 Okta, Inc." from the Okta End-User Dashboard`,
				Optional:    true,
				Computed:    true,
			},
			"default_app_app_instance_id": schema.StringAttribute{
				Description: "Default app app instance id",
				Optional:    true,
			},
			"default_app_app_link_name": schema.StringAttribute{
				Description: "Default app app link name",
				Optional:    true,
			},
			"default_app_classic_application_uri": schema.StringAttribute{
				Description: "Default app classic application uri",
				Optional:    true,
			},
			"links": schema.StringAttribute{
				Description: "Link relations for this object - JSON HAL - Discoverable resources related to the brand",
				Computed:    true,
			},

			"brand_id": schema.StringAttribute{
				Description:        "Brand ID - Note: Okta API for brands only reads and updates therefore the okta_brand resource needs to act as a quasi data source. Do this by setting brand_id.",
				Optional:           true,
				Computed:           true,
				DeprecationMessage: "Okta has fully support brand creation, this attribute is a no op and will be removed",
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *brandResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = resourceConfiguration(req, resp)
}

func (r *brandResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state brandResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	createReqBody, err := buildCreateBrandRequest(state)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to build brand request",
			err.Error(),
		)
		return
	}

	createdBrand, _, err := r.oktaSDKClientV3.CustomizationAPI.CreateBrand(ctx).CreateBrandRequest(createReqBody).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to create brand",
			err.Error(),
		)
		return
	}

	updateReqBody, err := buildUpdateBrandRequest(state)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to build brand request",
			err.Error(),
		)
		return
	}

	updatedBrand, _, err := r.oktaSDKClientV3.CustomizationAPI.ReplaceBrand(ctx, createdBrand.GetId()).Brand(updateReqBody).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to update brand",
			err.Error(),
		)
		return
	}

	brand, _, err := r.oktaSDKClientV3.CustomizationAPI.GetBrand(ctx, updatedBrand.GetId()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to read brand",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(mapBrandToState(brand, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *brandResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state brandResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var brandWithEmbedded *okta.BrandWithEmbedded
	var err error

	if state.BrandID.ValueString() == "default" {
		brands, _, err := r.oktaSDKClientV3.CustomizationAPI.ListBrands(ctx).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"failed to get list brand",
				err.Error(),
			)
			return
		}

		for _, brand := range brands {
			if *brand.IsDefault {
				brandWithEmbedded = &brand
			}
		}
	} else {
		brandWithEmbedded, _, err = r.oktaSDKClientV3.CustomizationAPI.GetBrand(ctx, state.ID.ValueString()).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"failed to read brand",
				err.Error(),
			)
			return
		}
	}

	resp.Diagnostics.Append(mapBrandToState(brandWithEmbedded, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *brandResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state brandResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.oktaSDKClientV3.CustomizationAPI.DeleteBrand(ctx, state.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to delete brand",
			err.Error(),
		)
		return
	}
}

func (r *brandResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state brandResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqBody, err := buildUpdateBrandRequest(state)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to build brand request",
			err.Error(),
		)
		return
	}

	updatedBrand, _, err := r.oktaSDKClientV3.CustomizationAPI.ReplaceBrand(ctx, state.ID.ValueString()).Brand(reqBody).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to update brand",
			err.Error(),
		)
		return
	}

	brand, _, err := r.oktaSDKClientV3.CustomizationAPI.GetBrand(ctx, updatedBrand.GetId()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to read brand",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(mapBrandToState(brand, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *brandResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func buildCreateBrandRequest(model brandResourceModel) (okta.CreateBrandRequest, error) {
	return okta.CreateBrandRequest{
		Name: model.Name.ValueString(),
	}, nil
}

func buildUpdateBrandRequest(model brandResourceModel) (okta.BrandRequest, error) {
	defaultApp := &okta.DefaultApp{}
	if !model.DefaultAppAppInstanceID.IsNull() && model.DefaultAppAppInstanceID.ValueString() != "" {
		defaultApp.AppInstanceId = model.DefaultAppAppInstanceID.ValueStringPointer()
	}
	if !model.DefaultAppAppLinkName.IsNull() && model.DefaultAppAppLinkName.ValueString() != "" {
		defaultApp.AppLinkName = model.DefaultAppAppLinkName.ValueStringPointer()
	}
	if !model.DefaultAppClassicApplicationURI.IsNull() && model.DefaultAppClassicApplicationURI.ValueString() != "" {
		defaultApp.ClassicApplicationUri = model.DefaultAppClassicApplicationURI.ValueStringPointer()
	}
	return okta.BrandRequest{
		Name:                       model.Name.ValueStringPointer(),
		AgreeToCustomPrivacyPolicy: model.AgreeToCustomPrivacyPolicy.ValueBoolPointer(),
		CustomPrivacyPolicyUrl:     model.CustomPrivacyPolicyURL.ValueStringPointer(),
		DefaultApp:                 defaultApp,
		Locale:                     model.Locale.ValueStringPointer(),
		RemovePoweredByOkta:        model.RemovePoweredByOkta.ValueBoolPointer(),
	}, nil
}

func mapBrandToState(data *okta.BrandWithEmbedded, state *brandResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	state.ID = types.StringPointerValue(data.Id)
	state.BrandID = types.StringPointerValue(data.Id)
	state.Name = types.StringPointerValue(data.Name)
	state.IsDefault = types.BoolPointerValue(data.IsDefault)
	state.EmailDomainID = types.StringPointerValue(data.EmailDomainId)
	state.Locale = types.StringPointerValue(data.Locale)
	state.AgreeToCustomPrivacyPolicy = types.BoolPointerValue(data.AgreeToCustomPrivacyPolicy)
	state.CustomPrivacyPolicyURL = types.StringPointerValue(data.CustomPrivacyPolicyUrl)
	state.RemovePoweredByOkta = types.BoolPointerValue(data.RemovePoweredByOkta)
	state.DefaultAppAppInstanceID = types.StringPointerValue(data.DefaultApp.AppInstanceId)
	state.DefaultAppAppLinkName = types.StringPointerValue(data.DefaultApp.AppLinkName)
	state.DefaultAppClassicApplicationURI = types.StringPointerValue(data.DefaultApp.ClassicApplicationUri)
	links, _ := json.Marshal(data.GetLinks())
	state.Links = types.StringValue(string(links))
	return diags
}
